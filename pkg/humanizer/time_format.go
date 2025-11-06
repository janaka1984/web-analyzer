package humanizer

import (
	"time"
)

// FormatDate formats a time.Time into a human-readable string.
func FormatDate(t time.Time) string {
	if t.IsZero() {
		return "-"
	}
	return t.Format("01 Jan 2025, 15:04")
}
