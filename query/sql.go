package query

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	valid "github.com/asaskevich/govalidator"
	"github.com/gernest/hiro/models"
	"github.com/jmoiron/sqlx"
	manager "github.com/ory/ladon/manager/sql"
	//postgres
	_ "github.com/lib/pq"
	uuid "github.com/satori/go.uuid"
)


// SQL implements Query interface.
type SQL struct {
	db   *sql.DB
	conn string
}

// New returns *SQL instance with db.
func New(driver, conn string) (*SQL, error) {
	db, err := sql.Open(driver, conn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return &SQL{db: db, conn: conn}, nil
}

// DB returns the underlying *sql.DB instance.
func (s *SQL) DB() *sql.DB {
	return s.db
}

// CreateQR inserts the qr code to the database.
func (s *SQL) CreateQR(ctx context.Context, qr *models.QR, groups ...string) error {
	q := `INSERT into qr (
				uuid,account_id,name,
				uri,image_url,data,should_redirect,
				context,redirect_url,
				created_at,updated_at)
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11);`

	if qr.Context == nil {
		qr.Context = []*models.ContextItem{}
	}
	data, err := json.Marshal(qr.Context)
	if err != nil {
		return err
	}
	_, err = s.db.ExecContext(ctx, q,
		qr.UUID, qr.UserID, qr.Name, qr.URL, qr.ImageURL, qr.ImageData,
		qr.ShouldRedirect, data, qr.RedirectURL,
		qr.CreatedAt, qr.UpdatedAt,
	)
	if err != nil {
		return err
	}
	if groups != nil {
		if query, args, ok := updateGroups(qr.UUID, groups...); ok {
			_, err = s.db.ExecContext(ctx, query, args...)
			return err
		}
	}
	return nil
}

func updateGroups(code uuid.UUID, groups ...string) (string, []interface{}, bool) {
	if len(groups) > 0 {
		q := "INSERT INTO qr_collections VALUES "
		i := 2
		f := "($1,$%d)"
		args := []interface{}{code}
		for _, v := range groups {
			if i > 2 {
				q += "," + fmt.Sprintf(f, i)
			} else {
				q += fmt.Sprintf(f, i)
			}
			i++
			args = append(args, v)
		}
		return q + ";", args, true
	}
	return "", nil, false
}

func (s *SQL) getGroupsByQR(ctx context.Context, id uuid.UUID) ([]*models.Collection, error) {
	q := `select c.uuid, c.name,c.color,
			c.created_at,c.updated_at from collections as c
  			inner join qr_collections
			on c.uuid=qr_collections.collections_id 
			where qr_collections.qr_id=$1 ;`
	rows, err := s.db.QueryContext(ctx, q, id)
	if err != nil {
		return nil, err
	}
	var groups []*models.Collection
	defer rows.Close()
	for rows.Next() {
		g := &models.Collection{}
		err = rows.Scan(&g.UUID, &g.Name, &g.Color, &g.CreatedAt, &g.UpdatedAt)
		if err != nil {
			return nil, err
		}
		groups = append(groups, g)
	}
	return groups, nil

}

// GetQR queries for qrcode image data stored in the database as raw bytes.
func (s *SQL) GetQR(ctx context.Context, id uuid.UUID) (*models.QR, error) {
	q := `select uuid,
			account_id,name,
			uri,
			image_url,
			data,
			should_redirect,
			context,
			redirect_url,
			created_at,
			updated_at
			from qr where uuid=$1 limit 1;`
	o := &models.QR{}
	var data json.RawMessage
	err := s.db.QueryRowContext(ctx, q, id).Scan(
		&o.UUID, &o.UserID, &o.Name, &o.URL, &o.ImageURL, &o.ImageData, &o.ShouldRedirect,
		&data, &o.RedirectURL,
		&o.CreatedAt, &o.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	if data != nil {
		o.Context = []*models.ContextItem{}
		err = json.Unmarshal(data, &o.Context)
		if err != nil {
			return nil, err
		}
	}
	return o, nil
}

// UpdateQR updates the database qr info.
func (s *SQL) UpdateQR(ctx context.Context, q *models.QR) error {
	ss := `
	UPDATE qr
	SET (name,should_redirect,
		 redirect_url,
		 context,updated_at)=($2,
				   $3,
				   $4,
				   $5,$6)
	WHERE uuid=$1
	`
	var data []byte
	var err error
	if q.Context != nil {
		data, err = json.Marshal(q.Context)
		if err != nil {
			return err
		}
	}
	_, err = s.db.ExecContext(ctx, ss, q.UUID, q.Name,
		q.ShouldRedirect, q.RedirectURL, data, time.Now())
	return err
}

// CreateToken creates a role and assigns the role to the given token.
func (s *SQL) CreateToken(ctx context.Context, tk *models.Token) error {
	q := `
	INSERT INTO tokens (uuid,issuer,subject,expires_at,not_before,created_at,updated_at)
	VALUES ($1,
			$2,
			$3,
			$4,
			$5,
			$6,
			$7);
	`
	_, err := s.db.ExecContext(ctx, q,
		tk.UUID, tk.Issuer, tk.Subject, tk.ExpiresOn,
		tk.NotBefore, tk.CreatedAt, tk.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

// CreateAccount saves the account request into the database. The password is
// hashed with bcrypt before storing.
func (s *SQL) CreateAccount(ctx context.Context, a *models.Account) error {
	q := `

	INSERT INTO accounts (uuid,name,email, password, created_at,updated_at)
	VALUES ($1,
			$2,
			$3,
			$4,
			$5,
			$6)`
	_, err := s.db.ExecContext(ctx, q, a.UUID, a.Name,
		a.Email, a.Password, a.CreatedAt, a.UpdatedAt)
	return err
}

// GetAccount finds account by name/email address.
func (s *SQL) GetAccount(ctx context.Context, name string) (*models.Account, error) {
	q := `
SELECT uuid,
       name,
       email,
	   password,
       created_at,
       updated_at
FROM accounts
WHERE name=$1`
	if valid.IsEmail(name) {
		q = `
SELECT uuid,
       name,
       email,
	   password,
       created_at,
       updated_at
FROM accounts
WHERE email=$1
		`
	}
	a := &models.Account{}
	err := s.db.QueryRowContext(ctx, q, name).Scan(
		&a.UUID, &a.Name, &a.Email, &a.Password,
		&a.CreatedAt, &a.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return a, nil
}

// GetToken returns *models.Token with the given id.
func (s *SQL) GetToken(ctx context.Context, id uuid.UUID) (*models.Token, error) {
	q := `
SELECT uuid,
       issuer,
       subject,
       expires_at,
       not_before,
       created_at,
       updated_at
FROM tokens
WHERE uuid=$1;

	`
	tk := &models.Token{}
	err := s.db.QueryRowContext(ctx, q, id).Scan(
		&tk.UUID, &tk.Issuer, &tk.Subject, &tk.ExpiresOn, &tk.NotBefore,
		&tk.CreatedAt, &tk.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	a, err := s.GetAccountByUUID(ctx, tk.Issuer)
	if err != nil {
		return nil, err
	}
	tk.IssuerAccount = a
	b, err := s.GetAccountByUUID(ctx, tk.Subject)
	if err != nil {
		return nil, err
	}
	tk.SubjectAccount = b
	return tk, nil
}

func (s *SQL) GetAccountByUUID(ctx context.Context, id uuid.UUID) (*models.Account, error) {
	q := `
SELECT uuid,
       name,
       email,
       password,
       created_at,
       updated_at
FROM accounts
WHERE uuid=$1`
	a := &models.Account{}
	err := s.db.QueryRowContext(ctx, q, id).Scan(
		&a.UUID, &a.Name, &a.Email, &a.Password, &a.CreatedAt, &a.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func (s *SQL) ListQR(ctx context.Context,
	id uuid.UUID, opts ...*models.ListOptions) ([]*models.QR, error) {
	q := `
SELECT uuid,
       account_id,name,
	   uri,
	   image_url,
       data,
       should_redirect,
       context,
	   redirect_url,
       created_at,
       updated_at
FROM qr WHERE account_id=$1 ORDER BY updated_at DESC %s;
	`
	q = fmt.Sprintf(q, addListOpts(opts...))
	rows, err := s.db.QueryContext(ctx, q, id)
	if err != nil {
		return nil, err
	}
	var out []*models.QR
	defer rows.Close()

	for rows.Next() {
		o := &models.QR{}
		var data json.RawMessage
		err := rows.Scan(
			&o.UUID, &o.UserID, &o.Name, &o.URL, &o.ImageURL, &o.ImageData, &o.ShouldRedirect,
			&data, &o.RedirectURL,
			&o.CreatedAt, &o.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		if data != nil {
			o.Context = []*models.ContextItem{}
			err = json.Unmarshal(data, &o.Context)
			if err != nil {
				return nil, err
			}
		}
		out = append(out, o)
	}
	for i := range out {
		groups, err := s.getGroupsByQR(ctx, out[i].UUID)
		if err != nil {
			if err != sql.ErrNoRows {
				return nil, err
			}
		}
		if groups != nil {
			out[i].Groups = append(out[i].Groups, groups...)
		}
	}
	return out, nil
}

// TotalCodes returns the total number of qrcodes belonging to the user.
func (s *SQL) TotalCodes(ctx context.Context, userID uuid.UUID) int64 {
	var count int64
	q := `select count(*) from qr where account_id=$1;`
	s.db.QueryRowContext(ctx, q, userID).Scan(&count)
	return count
}

// DeleteQR deletes qrcode with id from the database.
func (s *SQL) DeleteQR(ctx context.Context, id uuid.UUID) error {
	q := `delete from qr where uuid=$1`
	_, err := s.db.ExecContext(ctx, q, id)
	return err
}

// Close closes the underlying connection.
func (s *SQL) Close() error {
	return s.db.Close()
}

const schema = `
CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE  IF NOT EXISTS accounts(
	id serial,
	uuid uuid UNIQUE,
	name varchar(255) UNIQUE NOT NULL,
	email citext UNIQUE,
	password text,
	created_at timestamptz not null default(now() at time zone 'utc'),
	updated_at timestamptz not null default(now() at time zone 'utc')
);
CREATE TABLE  IF NOT EXISTS tokens(
	id serial,
	uuid uuid UNIQUE NOT NULL,
	issuer uuid REFERENCES accounts(uuid) ON DELETE CASCADE,
	subject uuid REFERENCES accounts(uuid),
	expires_at timestamptz not null,
	not_before timestamptz not null,
	created_at timestamptz not null default(now() at time zone 'utc'),
	updated_at timestamptz not null default(now() at time zone 'utc')
);

CREATE TABLE  IF NOT EXISTS qr(
	id serial,
	uuid uuid PRIMARY KEY,
	account_id uuid REFERENCES accounts(uuid) ON DELETE CASCADE,
	name text,
	uri text,
	image_url text,
	data bytea,
	should_redirect boolean not null default false,
	context json,
	redirect_url text,
	user_group integer not null default 1,
	permissions  integer not null default 500,
	created_at timestamptz not null default(now() at time zone 'utc'),
	updated_at timestamptz not null default(now() at time zone 'utc')
);

CREATE TABLE IF NOT EXISTS collections (
	id serial,
	uuid uuid PRIMARY KEY,
	account_id uuid REFERENCES accounts(uuid) ON DELETE CASCADE,
	name text,
	color text,
	created_at timestamptz not null default(now() at time zone 'utc'),
	updated_at timestamptz not null default(now() at time zone 'utc')
);

CREATE TABLE IF NOT EXISTS qr_collections (
	qr_id uuid REFERENCES qr(uuid) ON DELETE CASCADE,
	collections_id uuid REFERENCES collections(uuid) ON DELETE CASCADE,
	CONSTRAINT qr_collections_pkey PRIMARY KEY (qr_id, collections_id) 
);
`

const down = `
DROP TABLE IF EXISTS tokens;
DROP TABLE IF EXISTS qr_collections;
DROP TABLE IF EXISTS qr;
DROP TABLE IF EXISTS collections;
DROP TABLE IF EXISTS accounts;

DROP TABLE IF EXISTS ladon_policy_resource;
DROP TABLE IF EXISTS ladon_policy_subject;
DROP TABLE IF EXISTS ladon_policy_permission;
DROP TABLE IF EXISTS ladon_policy_resource_rel;
DROP TABLE IF EXISTS ladon_policy_action_rel;
DROP TABLE IF EXISTS ladon_policy_subject_rel;

DROP TABLE IF EXISTS ladon_action;
DROP TABLE IF EXISTS ladon_subject;
DROP TABLE IF EXISTS ladon_resource;
DROP TABLE IF EXISTS ladon_policy ;
DROP TABLE IF EXISTS gorp_migrations ;

`

// Up runs migrations up for the backend. This also includes the migrations for
// IF EXISTS ladon, the access control manager.
func (s *SQL) Up(ctx context.Context) error {
	if _, err := s.db.ExecContext(ctx, schema); err != nil {
		return err
	}
	db := sqlx.NewDb(s.db, "postgres")
	mngr := manager.NewSQLManager(db, nil)
	_, err := mngr.CreateSchemas("", "")
	return err
}

// Down drops all tables associate's with the backend.
func (s *SQL) Down(ctx context.Context) error {
	_, err := s.db.ExecContext(ctx, down)
	return err
}

// CreateCollection inserts the collection c into the database.
func (s *SQL) CreateCollection(ctx context.Context, c *models.Collection) (*models.Collection, error) {
	q := `INSERT INTO collections (uuid,account_id,name,color,created_at,updated_at)
		VALUES ($1,$2,$3,$4,$5,$6);`
	_, err := s.db.ExecContext(ctx, q, c.UUID, c.AccountID,
		c.Name, c.Color, c.CreatedAt, c.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return c, nil
}

// ListCollections query for all collections belonging to the user with id. If
// opts is supplied then LIMIT and OFFSEt on the query will be applied.
func (s *SQL) ListCollections(ctx context.Context, id uuid.UUID, opts ...*models.ListOptions) ([]*models.Collection, error) {
	q := `SELECT uuid,account_id,name,color,created_at,updated_at FROM collections 
		WHERE account_id=$1 ORDER BY updated_at DESC %s;`
	q = fmt.Sprintf(q, addListOpts(opts...))
	var c []*models.Collection
	rows, err := s.db.QueryContext(ctx, q, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		v := &models.Collection{}
		err = rows.Scan(&v.UUID, &v.AccountID, &v.Name,
			&v.Color, &v.CreatedAt, &v.UpdatedAt)
		if err != nil {
			return nil, err
		}
		c = append(c, v)
	}
	return c, nil
}

func addListOpts(opts ...*models.ListOptions) string {
	o := &models.ListOptions{}
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Limit == 0 {
		o.Limit = 10
	}
	return fmt.Sprintf(" OFFSET %d LIMIT %d", o.Offset, o.Limit)
}

func (s *SQL) GetCollection(ctx context.Context, id uuid.UUID) (*models.Collection, error) {
	q := `SELECT uuid,name,color,created_at,updated_at FROM collections WHERE uuid=$1;`
	c := &models.Collection{}
	err := s.db.QueryRowContext(ctx, q, id).Scan(
		&c.UUID, &c.Name, &c.Color, &c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (s *SQL) DeleteCollection(ctx context.Context, id uuid.UUID) error {
	q := `DELETE  FROM collections WHERE uuid=$1;`
	_, err := s.db.ExecContext(ctx, q, id)
	return err
}

// AssignCollection assigns the qrcode with uuid qr to the collection with given
// id.
func (s *SQL) AssignCollection(ctx context.Context, id uuid.UUID, qr uuid.UUID) error {
	q := `INSERT INTO qr_collections (collections_id,qr_id) VALUES ($1,$2);`
	_, err := s.db.ExecContext(ctx, q, id, qr)
	return err
}

// DeAssignCollection removes many to many relationship between the qrcode and
// the group.
func (s *SQL) DeAssignCollection(ctx context.Context, id uuid.UUID, qr uuid.UUID) error {
	q := `DELETE  FROM qr_collections WHERE collections_id=$1 AND qr_id=$2;`
	_, err := s.db.ExecContext(ctx, q, id, qr)
	return err
}

// UpdateCollections updates collection details (mainly collection name and
// color).
func (s *SQL) UpdateCollections(ctx context.Context, c *models.Collection) (*models.Collection, error) {
	q := `UPDATE collections SET (name,updated_at)=($1,$2) WHERE uuid=$3`
	now := time.Now()
	_, err := s.db.ExecContext(ctx, q, c.Name, now, c.UUID)
	if err != nil {
		return nil, err
	}
	c.UpdatedAt = now
	return c, err
}
