package utils

import "time"

func StartOfWeek(t time.Time) time.Time {
	year, week := t.ISOWeek()
	return time.Date(year, 0, 0, 0, 0, 0, 0, time.UTC).AddDate(0, 0, 7*(week-1))
}
