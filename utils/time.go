package utils

import (
	"time"
)

// StartOfDay returns the start of the day (00:00:00) for a given time
func StartOfDay(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}

// EndOfDay returns the end of the day (23:59:59.999999999) for a given time
func EndOfDay(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 23, 59, 59, 999999999, t.Location())
}

// StartOfDayUTC returns the start of the day in UTC
func StartOfDayUTC(t time.Time) time.Time {
	return StartOfDay(t.UTC())
}

// EndOfDayUTC returns the end of the day in UTC
func EndOfDayUTC(t time.Time) time.Time {
	return EndOfDay(t.UTC())
}

// IsToday checks if the given time is today
func IsToday(t time.Time) bool {
	now := time.Now().UTC()
	tUTC := t.UTC()
	return StartOfDayUTC(now).Equal(StartOfDayUTC(tUTC))
}

// GetDateString returns date in YYYY-MM-DD format
func GetDateString(t time.Time) string {
	return t.Format("2006-01-02")
}

// ParseDateString parses a date string in YYYY-MM-DD format
func ParseDateString(s string) (time.Time, error) {
	return time.Parse("2006-01-02", s)
}

// GetPastDates returns a list of dates from startDate to endDate (inclusive)
func GetPastDates(startDate, endDate time.Time) []time.Time {
	var dates []time.Time
	current := StartOfDayUTC(startDate)
	end := StartOfDayUTC(endDate)

	for !current.After(end) {
		dates = append(dates, current)
		current = current.AddDate(0, 0, 1)
	}

	return dates
}

// GetYesterday returns yesterday's date at start of day
func GetYesterday() time.Time {
	return StartOfDayUTC(time.Now().UTC().AddDate(0, 0, -1))
}

// NowUTC returns current time in UTC
func NowUTC() time.Time {
	return time.Now().UTC()
}
