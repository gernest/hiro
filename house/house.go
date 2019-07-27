package house

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/gernest/hiro/bus"
	//clickhouse driver
	_ "github.com/mailru/go-clickhouse"
)

// Client is a wrapper around clickhouse client offering methods for storing and
// querying clickhouse database.
type Client struct {
	db *sql.DB
}

// New opens a database connection to the clickhouse server.
func New(conn string) (*Client, error) {
	db, err := sql.Open("clickhouse", conn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &Client{db: db}, nil
}

func (c *Client) Up(ctx context.Context) error {
	q := `
CREATE TABLE IF NOT EXISTS scan_stats(
	status UInt64,
	qr_id String,
	user_id String,
	stamp Datetime
) ENGINE = Memory;
`
	_, err := c.db.ExecContext(ctx, q)
	return err
}

// SaveScanEvent stores ScanEvent.
func (c *Client) SaveScanEvent(ctx context.Context, s *bus.ScanEvent) error {
	q := `INSERT INTO scan_stats (qr_id,status,user_id,stamp) VALUES (?,?,?,?);`
	stamp, err := time.Parse(time.RFC3339, s.Timestamp)
	if err != nil {
		return err
	}
	_, err = c.db.ExecContext(ctx, q, s.ID, s.Status, s.User, stamp)
	return err
}

// SaveScanBatch writes the scan events in a single query.
func (c *Client) SaveScanBatch(ctx context.Context, s []*bus.ScanEvent) error {
	q, args, err := batch(s)
	if err != nil {
		return err
	}
	_, err = c.db.ExecContext(ctx, q, args...)
	return err
}

func batch(s []*bus.ScanEvent) (string, []interface{}, error) {
	query := "INSERT INTO scan_stats (qr_id,status,user_id,stamp) VALUES "
	values := "($%d,$%d,$%d,$%d)"
	var out string
	var args []interface{}
	i := 0
	for _, event := range s {
		if out != "" {
			out += ","
		}
		out += fmt.Sprintf(values, i+1, i+2, i+3, i+4)
		i += 4
		stamp, err := time.Parse(time.RFC3339, event.Timestamp)
		if err != nil {
			return "", nil, err
		}
		args = append(args, event.ID, event.Status, event.User, stamp)
	}
	return query + out + ";", args, nil
}
