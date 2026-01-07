// Package timeperiod provides TimePeriod comparing different time periods.
// Same concept as https://docs.oracle.com/javase/8/docs/api/java/time/Period.html.
package timeperiod

import (
	"errors"
	"time"
)

var ErrEndTimeBeforeStartTime = errors.New("end time before start time")

var (
	_ TimePeriod = new(startTimeEndTimePeriod)

	// Infinite time period constant.
	//nolint:gochecknoglobals // global to improve readability.
	Infinite TimePeriod = startTimeEndTimePeriod{}
)

type (

	// TimePeriod to track a Period of Time. It's composed of a StartTime and an EndTime
	// If StartTime is nil, then it means the beginning of time.
	// If EndTime is nil, then it means the end of time.
	TimePeriod interface {
		StartTime() *time.Time
		EndTime() *time.Time

		// Duration Returns the duration of this period.
		Duration() time.Duration
		// Overlaps Returns the overlap period between the two time periods, and whether it overlaps or not.
		Overlaps(other TimePeriod) (TimePeriod, bool)
	}

	startTimeEndTimePeriod struct {
		startTime *time.Time
		endTime   *time.Time
	}
)

// New Creates a new time period based on a start time and an end time
// Returns either the time period of an error is the end time is before the start time.
func New(startTime, endTime *time.Time) (TimePeriod, error) {
	if (startTime != nil && endTime != nil) && endTime.Before(*startTime) {
		return startTimeEndTimePeriod{}, ErrEndTimeBeforeStartTime
	}

	return startTimeEndTimePeriod{
		startTime: startTime,
		endTime:   endTime,
	}, nil
}

// Must Create a new time period based on a start time, and an end time
// Panics if end time is before the start time.
func Must(startTime, endTime *time.Time) TimePeriod {
	period, err := New(startTime, endTime)
	if err != nil {
		panic(err)
	}

	return period
}

// Duration Returns the duration of this period.
func (tp startTimeEndTimePeriod) Duration() time.Duration {
	if tp.startTime == nil || tp.endTime == nil {
		// return maxDuration
		return 1<<63 - 1
	}

	return tp.endTime.Sub(*tp.startTime)
}

func (tp startTimeEndTimePeriod) EndTime() *time.Time {
	return tp.endTime
}

// Overlaps Returns the overlap period between the two time periods, and whether it overlaps or not.
func (tp startTimeEndTimePeriod) Overlaps(other TimePeriod) (TimePeriod, bool) {
	if tp.doesIntersect(other) {
		return tp.intersect(other), true
	}

	return startTimeEndTimePeriod{}, false
}

func (tp startTimeEndTimePeriod) StartTime() *time.Time {
	return tp.startTime
}

func (tp startTimeEndTimePeriod) doesIntersect(comparePeriod TimePeriod) bool {
	// Condition 1: tp.start < comparePeriod.end
	// True if comparePeriod.end is nil (infinity) or tp.start is nil (-infinity)
	if comparePeriod.EndTime() != nil && tp.startTime != nil {
		if !tp.startTime.Before(*comparePeriod.EndTime()) {
			return false
		}
	}

	// Condition 2: comparePeriod.start < tp.end
	// True if tp.end is nil (infinity) or comparePeriod.start is nil (-infinity)
	if tp.endTime != nil && comparePeriod.StartTime() != nil {
		if !comparePeriod.StartTime().Before(*tp.endTime) {
			return false
		}
	}

	return true
}

func (tp startTimeEndTimePeriod) intersect(comparePeriod TimePeriod) TimePeriod {
	start := comparePeriod.StartTime()
	switch {
	case tp.startTime == nil:
		start = comparePeriod.StartTime()
	case comparePeriod.StartTime() == nil:
		start = tp.startTime
	case tp.startTime.After(*comparePeriod.StartTime()):
		start = tp.startTime
	}

	end := comparePeriod.EndTime()
	switch {
	case tp.endTime == nil:
		end = comparePeriod.EndTime()
	case comparePeriod.EndTime() == nil:
		end = tp.endTime
	case tp.endTime.Before(*comparePeriod.EndTime()):
		end = tp.endTime
	}

	return startTimeEndTimePeriod{
		startTime: start,
		endTime:   end,
	}
}
