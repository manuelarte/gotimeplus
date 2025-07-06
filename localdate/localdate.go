// Package localdatetime provides LocalDate, storing a date, timezone independent.
// Same concept as https://docs.oracle.com/javase/8/docs/api/java/time/LocalDate.html.
package localdate

import "time"

var _ LocalDate = new(localDate)

type (
	LocalDate interface {
		// After reports whether the LocalDate is after the given other LocalDate.
		After(other LocalDate) bool
		// Before reports whether the LocalDate is before the given other LocalDate.
		Before(other LocalDate) bool
		Day() int
		// Equal reports whether the LocalDate is equal to the given other LocalDate.
		Equal(other LocalDate) bool
		Month() time.Month
		// ToTime converts the LocalDate to a time.Time at midnight in the provided location.
		ToTime(loc *time.Location) time.Time
		Year() int
	}

	localDate struct {
		year  int
		month time.Month
		day   int
	}
)

// New LocalDate from year, month and day.
func New(year int, month time.Month, day int) LocalDate {
	return &localDate{
		year:  year,
		month: month,
		day:   day,
	}
}

// FromTime converts time.Time to LocalDate.
func FromTime(t time.Time) LocalDate {
	return New(t.Year(), t.Month(), t.Day())
}

func (ld localDate) After(other LocalDate) bool {
	return ld.ToTime(time.UTC).After(other.ToTime(time.UTC))
}

func (ld localDate) Day() int {
	return ld.day
}

func (ld localDate) Before(other LocalDate) bool {
	return ld.ToTime(time.UTC).Before(other.ToTime(time.UTC))
}

func (ld localDate) Equal(other LocalDate) bool {
	return ld.ToTime(time.UTC).Equal(other.ToTime(time.UTC))
}

func (ld localDate) Month() time.Month {
	return ld.month
}

func (ld localDate) ToTime(loc *time.Location) time.Time {
	return time.Date(ld.year, ld.month, ld.day, 0, 0, 0, 0, loc)
}

func (ld localDate) Year() int {
	return ld.year
}
