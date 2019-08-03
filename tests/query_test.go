package tests

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/gernest/hiro/models"
	"github.com/gernest/hiro/query"
	"github.com/gernest/hiro/testutil"
	uuid "github.com/satori/go.uuid"
)

func RunQueryTest(t *testing.T, ctx *testutil.Context) {
	db := ctx.DB
	t.Run("accounts", func(ts *testing.T) {
		id := accountsTest(ts, db)
		ts.Run("tokens", func(tst *testing.T) {
			tokenTest(tst, db, id)
		})
		ts.Run("qrcode", func(tst *testing.T) {
			qrTest(tst, db, id)
		})
		t.Run("collections", func(ts *testing.T) {
			collectionTest(ts, db, id)
		})
	})
}

func accountsTest(t *testing.T, db *query.SQL) uuid.UUID {
	now := time.Now()
	a := &models.Account{
		UUID:      uuid.NewV4(),
		Email:     "mail@example.com",
		Password:  "pass",
		CreatedAt: now,
		UpdatedAt: now,
	}
	ctx := context.Background()

	err := db.CreateAccount(ctx, a)
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.GetAccount(ctx, a.Email)
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.GetAccount(ctx, a.Email)
	if err != nil {
		t.Fatal(err)
	}
	return a.UUID
}

func tokenTest(t *testing.T, db *query.SQL, id uuid.UUID) {
	tk := models.DefaultLoginToken(id)
	ctx := context.Background()
	err := db.CreateToken(ctx, tk)
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.GetToken(ctx, tk.UUID)
	if err != nil {
		t.Fatal(err)
	}
}

func qrTest(t *testing.T, db *query.SQL, id uuid.UUID) {
	now := time.Now()
	c := &models.QR{
		UUID:      uuid.NewV4(),
		UserID:    id,
		URL:       "/some/code",
		ImageData: []byte("data"),
		CreatedAt: now,
		UpdatedAt: now,
	}
	ctx := context.Background()
	err := db.CreateQR(ctx, c)
	if err != nil {
		t.Fatal(err)
	}
	g, err := db.GetQR(ctx, c.UUID)
	if err != nil {
		t.Fatal(err)
	}
	if !uuid.Equal(g.UUID, c.UUID) {
		t.Errorf("expected %v got %v", c.UUID, g.UUID)
	}

	g.Context = []*models.ContextItem{
		{
			Kind:  models.ItemString,
			Key:   "name",
			Value: "test",
		},
	}

	{
		err = db.UpdateQR(ctx, g)
		if err != nil {
			t.Fatal(err)
		}

		up, err := db.GetQR(ctx, g.UUID)
		if err != nil {
			t.Fatal(err)
		}
		if len(up.Context) != 1 {
			t.Errorf("expected %d got %d", 1, len(up.Context))
		}
		if up.Context[0].Kind != g.Context[0].Kind {
			t.Errorf("expected %s got %s", up.Context[0].Kind, g.Context[0].Kind)
		}
		if up.Context[0].Key != g.Context[0].Key {
			t.Errorf("expected %s got %s", up.Context[0].Key, g.Context[0].Key)
		}
		if up.Context[0].Value.(string) != g.Context[0].Value.(string) {
			t.Errorf("expected %v got %v", up.Context[0].Value, g.Context[0].Value)
		}
	}

	list, err := db.ListQR(ctx, id)
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 1 {
		t.Errorf("expected 1 qrcode got %d", len(list))
	}

	err = db.DeleteQR(ctx, c.UUID)
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.GetQR(ctx, c.UUID)
	if err != sql.ErrNoRows {
		t.Errorf("expected %v got %v", sql.ErrNoRows, err)
	}
}

func collectionTest(t *testing.T, db *query.SQL, usr uuid.UUID) {
	names := []string{"q", "b", "c"}
	var list []*models.Collection
	ctx := context.Background()
	for _, n := range names {
		now := time.Now()
		c, err := db.CreateCollection(ctx, &models.Collection{
			UUID:      uuid.NewV4(),
			Name:      n,
			AccountID: usr,
			CreatedAt: now,
			UpdatedAt: now,
		})
		if err != nil {
			t.Fatal(err)
		}
		list = append(list, c)
	}

	for _, v := range list {
		c, err := db.GetCollection(ctx, v.UUID)
		if err != nil {
			t.Fatal(err)
		}
		if c.Name != v.Name {
			t.Errorf("expected %s got %s", v.Name, c.Name)
		}
	}
	cs, err := db.ListCollections(ctx, usr)
	if err != nil {
		t.Fatal(err)
	}
	if len(cs) != len(list) {
		t.Errorf("expected %d got %d", len(list), len(cs))
	}
	for _, v := range list {
		err := db.DeleteCollection(ctx, v.UUID)
		if err != nil {
			t.Fatal(err)
		}
	}
	cs, err = db.ListCollections(ctx, usr)
	if err != nil {
		t.Fatal(err)
	}
	if len(cs) != 0 {
		t.Errorf("expected 0 items got %d", len(cs))
	}
}

func TestUpdateGroup(t *testing.T) {
	q, args, ok := query.UpdateGroups(uuid.UUID{}, "one", "two")
	if !ok {
		t.Fatal("expected to be true")
	}
	if len(args) != 3 {
		t.Errorf("expected 3 got %d", len(args))
	}
	e := "INSERT INTO qr_collections VALUES ($1,$2),($1,$3);"
	if q != e {
		t.Errorf("expected %s got %s", e, q)
	}
}
