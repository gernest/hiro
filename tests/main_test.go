package tests

import (
	"context"
	"os"
	"testing"

	"github.com/gernest/hiro/query"
	"github.com/gernest/hiro/testutil"
)

type TestCase struct {
	Name string
	Case func(*testing.T, *testutil.Context)
}

func Case(name string, fn func(*testing.T, *testutil.Context)) TestCase {
	return TestCase{Name: name, Case: fn}
}

func RunCase(t *testing.T, db *query.SQL, cs TestCase) {
	err := db.Up(context.Background())
	if err != nil {
		t.Errorf("initializing test context :%v", err)
		return
	}
	ctx, err := testutil.New(db)
	if err != nil {
		t.Errorf("initializing test context :%v", err)
		return
	}

	defer func() {
		db.Down(context.Background())
	}()
	if cs.Case != nil {
		t.Run(cs.Name, func(ts *testing.T) {
			cs.Case(ts, ctx)
		})
	}
}

func TestIntegration(t *testing.T) {
	conn := os.Getenv("HIRO_DB_CONN")
	db, err := query.New("postgres", conn)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	Run(t, db,
		Case("query", RunQueryTest),
		Case("qrcode", RunQRcodeTest),
		Case("collections", RunCollectionsTest),
		Case("accounts", RunAccountsTes),
	)
}

func Run(t *testing.T, db *query.SQL, cases ...TestCase) {
	for _, cs := range cases {
		RunCase(t, db, cs)
	}
}
