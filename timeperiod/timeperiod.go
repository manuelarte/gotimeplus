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
		StartTime() time.Time
		EndTime() time.Time

		// Duration Get the duration.
		Duration() time.Duration
		Overlaps(other TimePeriod) (TimePeriod, bool)
	}

	startTimeEndTimePeriod struct {
		startTime time.Time
		endTime   time.Time
	}
)

// New Creates new time period based on a start time and an end time
// Returns either the time period of an error is the end time is before the start time.
func New(startTime, endTime time.Time) (TimePeriod, error) {
	if (!startTime.IsZero() && !endTime.IsZero()) && endTime.Before(startTime) {
		return startTimeEndTimePeriod{}, ErrEndTimeBeforeStartTime
	}

	return startTimeEndTimePeriod{
		startTime: startTime,
		endTime:   endTime,
	}, nil
}

// Must Creates new time period based on a start time and an end time
// Panics if end time is before the start time.
func Must(startTime, endTime time.Time) TimePeriod {
	period, err := New(startTime, endTime)
	if err != nil {
		panic(err)
	}

	return period
}

func (tp startTimeEndTimePeriod) Duration() time.Duration {
	if tp.startTime.IsZero() || tp.endTime.IsZero() {
		// return maxDuration
		return 1<<63 - 1
	}

	return tp.endTime.Sub(tp.startTime)
}

func (tp startTimeEndTimePeriod) EndTime() time.Time {
	return tp.endTime
}

// Overlaps Returns the overlap period between the two time periods, and the boolean whether it overlaps or not.
func (tp startTimeEndTimePeriod) Overlaps(other TimePeriod) (TimePeriod, bool) {
	if tp.doesIntersect(other) {
		return tp.intersect(other), true
	}

	return startTimeEndTimePeriod{startTime: time.Time{}, endTime: time.Time{}}, false
}

func (tp startTimeEndTimePeriod) StartTime() time.Time {
	return tp.startTime
}

func (tp startTimeEndTimePeriod) doesIntersect(comparePeriod TimePeriod) bool {
	if tp.endTime.IsZero() && comparePeriod.EndTime().IsZero() {
		return true
	}

	if comparePeriod.EndTime().IsZero() && comparePeriod.StartTime().UTC().After(tp.endTime.UTC()) {
		return false
	}

	if !tp.endTime.IsZero() && (tp.endTime.UTC().Before(comparePeriod.StartTime().UTC()) ||
		tp.endTime.UTC().Equal(comparePeriod.StartTime().UTC())) {
		return false
	}

	if !comparePeriod.EndTime().IsZero() && (tp.startTime.UTC().After(comparePeriod.EndTime().UTC()) ||
		tp.startTime.UTC().Equal(comparePeriod.EndTime().UTC())) {
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

	if tp.startTime.UTC().Before(comparePeriod.StartTime().UTC()) {
		intersectPeriod.startTime = comparePeriod.StartTime()
	}

	if !comparePeriod.EndTime().IsZero() && !tp.endTime.IsZero() &&
		tp.endTime.UTC().After(comparePeriod.EndTime().UTC()) {
		intersectPeriod.endTime = comparePeriod.EndTime()
	}

	if tp.endTime.IsZero() {
		intersectPeriod.endTime = comparePeriod.EndTime()
	}

	return intersectPeriod
}
