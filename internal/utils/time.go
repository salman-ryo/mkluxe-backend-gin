package utils

import "time"

func FormatISO(t time.Time) string {
	return t.Format(time.RFC3339)
}

func StartOfDay(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, t.Location())
}
