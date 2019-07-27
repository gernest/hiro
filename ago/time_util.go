package ago

import (
	"strconv"
	"time"
)


// Ago formats duration into a more human friendly way. This only format in
// terms of day/hours/minutes/seconds.
//
// duration less that 30 seconds will return "just now". The verbs go like "just
// now" , " a minute ago", "10 seconds ago", "3 days ago"...etc.
func Ago(duration time.Duration) string {
	seconds := int64(duration.Seconds()) % 60
	minutes := int64(duration.Minutes()) % 60
	hours := int64(duration.Hours()) % 24
	days := int64(duration/(24*time.Hour)) % 365 % 7
	switch {
	case days > 0:
		if days == 1 {
			return "a day ago"
		}
		return strconv.FormatInt(days, 10) + " days ago"
	case hours > 0:
		if hours == 1 {
			return "an hour ago"
		}
		return strconv.FormatInt(hours, 10) + " hours ago"
	case minutes > 0:
		if minutes == 1 {
			return "a minute ago"
		}
		return strconv.FormatInt(minutes, 10) + " minutes ago"
	default:
		if seconds < 30 {
			return "just now"
		}
		return strconv.FormatInt(seconds, 10) + " seconds ago"
	}
}
