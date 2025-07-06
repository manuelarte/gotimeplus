// Package localtime provides LocalTime, storing a time, timezone independent.
// Same concept as https://docs.oracle.com/javase/8/docs/api/java/time/LocalTime.html.
package localtime

import (
	"time"

	"github.com/manuelarte/gotimeplus/localdate"
)

var _ LocalTime = new(localTime)

type (
	LocalTime interface {
		// After reports whether the LocalDate is after the given other LocalDate.
		After(other LocalTime) bool
		// Before reports whether the LocalDate is before the given other LocalDate.
		Before(other LocalTime) bool
		// Equal reports whether the LocalDate is equal to the given other LocalDate.
		Equal(other LocalTime) bool
		Hour() int
		Min() int
		Nanosecond() int
		Sec() int
		// ToTime converts the LocalTime to a time.Time provided with a LocalDate and a location.
		ToTime(localDate localdate.LocalDate, loc *time.Location) time.Time
	}

	localTime struct {
		hour, min, sec, nsec int
	}
)

// New LocalTime from hours, minutes, seconds and nanoseconds.
func New(hour, minutes, sec, nsec int) LocalTime {
	return &localTime{
		hour: hour,
		min:  minutes,
		sec:  sec,
		nsec: nsec,
	}
}

func (lt localTime) After(other LocalTime) bool {
	ld := localdate.New(2009, time.November, 10)

	return lt.ToTime(ld, time.UTC).After(other.ToTime(ld, time.UTC))
}

func (lt localTime) Before(other LocalTime) bool {
	ld := localdate.New(2009, time.November, 10)

	return lt.ToTime(ld, time.UTC).Before(other.ToTime(ld, time.UTC))
}

func (lt localTime) Equal(other LocalTime) bool {
	ld := localdate.New(2009, time.November, 10)

	return lt.ToTime(ld, time.UTC).Equal(other.ToTime(ld, time.UTC))
}

func (lt localTime) Hour() int {
	return lt.hour
}

func (lt localTime) Min() int {
	return lt.min
}

func (lt localTime) Nanosecond() int {
	return lt.nsec
}

func (lt localTime) Sec() int {
	return lt.sec
}

func (lt localTime) ToTime(ld localdate.LocalDate, loc *time.Location) time.Time {
	return time.Date(ld.Year(), ld.Month(), ld.Day(), lt.Hour(), lt.Min(), lt.Sec(), lt.Nanosecond(), loc)
}
