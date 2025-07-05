// Package timeperiod provides TimePeriod comparing different time periods
package timeperiod

import (
	"errors"
	"time"
)

var ErrEndTimeBeforeStartTime = errors.New("end time before start time")

var (
	_ TimePeriod = new(startTimeEndTimePeriod)
	//nolint:gochecknoglobals // Infinite time period constant.
	Infinite TimePeriod = startTimeEndTimePeriod{}
)

type (

	// TimePeriod to track a Period of Time. It's composed of a StartTime and an EndTime
	// If StartTime is zero then it means the beginning of time.
	// If EndTime is zero then it means end of time.
	TimePeriod interface {
		GetStartTime() time.Time
		GetEndTime() time.Time

		// GetDuration Get the duration.
		GetDuration() time.Duration
		Overlaps(other TimePeriod) (TimePeriod, bool)
	}

	startTimeEndTimePeriod struct {
		startTime time.Time
		endTime   time.Time
	}
)

// NewTimePeriod Creates new time period based on a start time and an end time
// Returns either the time period of an error is the end time is before the start time.
func NewTimePeriod(startTime, endTime time.Time) (TimePeriod, error) {
	if (!startTime.IsZero() && !endTime.IsZero()) && endTime.Before(startTime) {
		return startTimeEndTimePeriod{}, ErrEndTimeBeforeStartTime
	}

	return startTimeEndTimePeriod{
		startTime: startTime,
		endTime:   endTime,
	}, nil
}

// MustTimePeriod Creates new time period based on a start time and an end time
// Panics if end time is before the start time.
func MustTimePeriod(startTime, endTime time.Time) TimePeriod {
	period, err := NewTimePeriod(startTime, endTime)
	if err != nil {
		panic(err)
	}

	return period
}

func (tp startTimeEndTimePeriod) GetStartTime() time.Time {
	return tp.startTime
}

func (tp startTimeEndTimePeriod) GetEndTime() time.Time {
	return tp.endTime
}

func (tp startTimeEndTimePeriod) GetDuration() time.Duration {
	if tp.startTime.IsZero() || tp.endTime.IsZero() {
		// return maxDuration
		return 1<<63 - 1
	}

	return tp.endTime.Sub(tp.startTime)
}

// Overlaps Returns the overlap period between the two time periods, and the boolean whether it overlaps or not.
func (tp startTimeEndTimePeriod) Overlaps(other TimePeriod) (TimePeriod, bool) {
	if tp.doesIntersect(other) {
		return tp.intersect(other), true
	}

	return startTimeEndTimePeriod{startTime: time.Time{}, endTime: time.Time{}}, false
}

func (tp startTimeEndTimePeriod) doesIntersect(comparePeriod TimePeriod) bool {
	if tp.endTime.IsZero() && comparePeriod.GetEndTime().IsZero() {
		return true
	}

	if comparePeriod.GetEndTime().IsZero() && comparePeriod.GetStartTime().UTC().After(tp.endTime.UTC()) {
		return false
	}

	if !tp.endTime.IsZero() && (tp.endTime.UTC().Before(comparePeriod.GetStartTime().UTC()) ||
		tp.endTime.UTC().Equal(comparePeriod.GetStartTime().UTC())) {
		return false
	}

	if !comparePeriod.GetEndTime().IsZero() && (tp.startTime.UTC().After(comparePeriod.GetEndTime().UTC()) ||
		tp.startTime.UTC().Equal(comparePeriod.GetEndTime().UTC())) {
		return false
	}

	return true
}

func (tp startTimeEndTimePeriod) intersect(comparePeriod TimePeriod) TimePeriod {
	if !tp.doesIntersect(comparePeriod) {
		return startTimeEndTimePeriod{
			startTime: time.Time{},
			endTime:   time.Time{},
		}
	}

	intersectPeriod := tp

	if tp.startTime.UTC().Before(comparePeriod.GetStartTime().UTC()) {
		intersectPeriod.startTime = comparePeriod.GetStartTime()
	}

	if !comparePeriod.GetEndTime().IsZero() && !tp.endTime.IsZero() &&
		tp.endTime.UTC().After(comparePeriod.GetEndTime().UTC()) {
		intersectPeriod.endTime = comparePeriod.GetEndTime()
	}

	if tp.endTime.IsZero() {
		intersectPeriod.endTime = comparePeriod.GetEndTime()
	}

	return intersectPeriod
}
