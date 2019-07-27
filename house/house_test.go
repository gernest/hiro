package house

import (
	"io/ioutil"
	"testing"
	"time"

	"github.com/gernest/hiro/bus"
)

func TestBatch(t *testing.T) {
	tm, err := time.Parse(time.ANSIC, "Mon Jan 2 15:04:05 2006")
	if err != nil {
		t.Fatal(err)
	}

	now := tm.Format(time.RFC3339)
	sample := []*bus.ScanEvent{
		{
			Timestamp: now,
			ID:        "one",
		},
		{
			Timestamp: now,
			ID:        "two",
		},
		{
			Timestamp: now,
			ID:        "three",
		},
	}
	expect := map[int]string{
		0: "one",
		4: "two",
		8: "three",
	}
	q, args, err := batch(sample)
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range []int{0, 4, 8} {
		got := args[v].(string)
		e := expect[v]
		if got != e {
			t.Errorf("expected %s got %s", e, got)
		}
	}
	eq := "INSERT INTO scan_stats (qr_id,status,user_id,stamp) VALUES ($1,$2,$3,$4),($5,$6,$7,$8),($9,$10,$11,$12);"
	if q != eq {
		ioutil.WriteFile("txt", []byte(q), 0600)
		t.Errorf("expected %s got %s", eq, q)
	}
}
