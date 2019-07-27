package ago

import (
	"testing"
	"time"
)

func TestAgo(t *testing.T) {
	sample := []struct {
		duration time.Duration
		text     string
	}{
		{
			20 * time.Second, "just now",
		},
		{
			45 * time.Second, "45 seconds ago",
		},
		{
			60 * time.Second, "a minute ago",
		},
		{
			120 * time.Second, "2 minutes ago",
		},
		{
			60 * 60 * time.Second, "an hour ago",
		},
		{
			60 * 60 * 2 * time.Second, "2 hours ago",
		},
		{
			60 * 60 * 24 * time.Second, "a day ago",
		},
		{
			60 * 60 * 24 * 2 * time.Second, "2 days ago",
		},
	}
	for _, v := range sample {
		got := Ago(v.duration)
		if got != v.text {
			t.Errorf("expected %s got %s", v.text, got)
		}
	}
}
