// Package localdatetime provides LocalDateTime, storing a date + time, timezone independent.
// Same concept as https://docs.oracle.com/javase/8/docs/api/java/time/LocalDateTime.html.
package localdatetime

import (
	"time"

	"github.com/manuelarte/gotimeplus/localdate"
	"github.com/manuelarte/gotimeplus/localtime"
)

var _ LocalDateTime = new(localDateTime)

type (
	LocalDateTime interface {
		// After reports whether the LocalDateTime is after the given other LocalDateTime.
		After(other LocalDateTime) bool
		// Before reports whether the LocalDateTime is before the given other LocalDateTime.
		Before(other LocalDateTime) bool
		// Equal reports whether the LocalDateTime is equal to the given other LocalDateTime.
		Equal(other LocalDateTime) bool
		// ToTime converts the LocalDateTime to a time.Time in the provided location.
		ToTime(loc *time.Location) time.Time
	}

	localDateTime struct {
		ld localdate.LocalDate
		lt localtime.LocalTime
	}
)

// NewLocalDateTime New LocalDateTime from localDate and localTime.
func NewLocalDateTime(ld localdate.LocalDate, lt localtime.LocalTime) LocalDateTime {
	return &localDateTime{
		ld: ld,
		lt: lt,
	}
}

// FromTime converts time.Time to LocalDate.
func FromTime(t time.Time) LocalDateTime {
	return NewLocalDateTime(
		localdate.FromTime(t),
		localtime.NewLocalTime(t.Hour(), t.Minute(), t.Second(), t.Nanosecond()),
	)
}

func (ldt localDateTime) After(other LocalDateTime) bool {
	return ldt.ToTime(time.UTC).After(other.ToTime(time.UTC))
}

func (ldt localDateTime) Before(other LocalDateTime) bool {
	return ldt.ToTime(time.UTC).Before(other.ToTime(time.UTC))
}

func (ldt localDateTime) Equal(other LocalDateTime) bool {
	return ldt.ToTime(time.UTC).Equal(other.ToTime(time.UTC))
}

func (ldt localDateTime) ToTime(loc *time.Location) time.Time {
	return time.Date(ldt.ld.Year(), ldt.ld.Month(), ldt.ld.Day(),
		ldt.lt.Hour(), ldt.lt.Min(), ldt.lt.Sec(), ldt.lt.Nanosecond(), loc)
}
