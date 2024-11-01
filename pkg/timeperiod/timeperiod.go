// Package timeperiod provides TimePeriod comparing different time periods
package timeperiod

import (
	"errors"
	"time"
)

var ErrEndTimeBeforeStartTime = errors.New("end time before start time")

// TimePeriod struct to track a Period of Time. It's composed of a StartTime and an EndTime
// If StartTime is zero then it means the beginning of time
// If EndTime is zero then it means no limit
type TimePeriod interface {
	GetStartTime() time.Time
	GetEndTime() time.Time
	// Overlaps Returns the overlap period between the two time periods, and the boolean wheter it overlaps or not
	Overlaps(other TimePeriod) (TimePeriod, bool)
	// GetDuration Get the duration
	GetDuration() time.Duration
}

// NewTimePeriod Creates new time period based on a start time and an end time
// Returns either the time period of an error is the end time is before the start time
func NewTimePeriod(startTime time.Time, endTime time.Time) (TimePeriod, error) {
	if (!startTime.IsZero() && !endTime.IsZero()) && endTime.Before(startTime) {
		return timePeriodImpl{}, ErrEndTimeBeforeStartTime
	}
	return timePeriodImpl{
		StartTime: startTime,
		EndTime:   endTime,
	}, nil
}

var _ TimePeriod = new(timePeriodImpl)

type timePeriodImpl struct {
	StartTime time.Time
	EndTime   time.Time
}

func (timePeriod timePeriodImpl) GetStartTime() time.Time {
	return timePeriod.StartTime
}

func (timePeriod timePeriodImpl) GetEndTime() time.Time {
	return timePeriod.EndTime
}

func (timePeriod timePeriodImpl) GetDuration() time.Duration {
	if timePeriod.StartTime.IsZero() || timePeriod.EndTime.IsZero() {
		// return maxDuration
		return 1<<63 - 1
	}
	return timePeriod.EndTime.Sub(timePeriod.StartTime)
}

func (timePeriod timePeriodImpl) doesIntersect(comparePeriod TimePeriod) bool {
	if timePeriod.EndTime.IsZero() && comparePeriod.GetEndTime().IsZero() {
		return true
	}
	if comparePeriod.GetEndTime().IsZero() && comparePeriod.GetStartTime().UTC().After(timePeriod.EndTime.UTC()) {
		return false
	}
	if !timePeriod.EndTime.IsZero() && (timePeriod.EndTime.UTC().Before(comparePeriod.GetStartTime().UTC()) || timePeriod.EndTime.UTC() == comparePeriod.GetStartTime().UTC()) {
		return false
	}
	if !comparePeriod.GetEndTime().IsZero() && (timePeriod.StartTime.UTC().After(comparePeriod.GetEndTime().UTC()) || timePeriod.StartTime.UTC() == comparePeriod.GetEndTime().UTC()) {
		return false
	}
	return true
}

func (timePeriod timePeriodImpl) intersect(comparePeriod TimePeriod) TimePeriod {
	if !timePeriod.doesIntersect(comparePeriod) {
		return timePeriodImpl{
			StartTime: time.Time{},
			EndTime:   time.Time{},
		}
	}
	intersectPeriod := timePeriod
	if timePeriod.StartTime.UTC().Before(comparePeriod.GetStartTime().UTC()) {
		intersectPeriod.StartTime = comparePeriod.GetStartTime()
	}
	if !comparePeriod.GetEndTime().IsZero() && !timePeriod.EndTime.IsZero() && timePeriod.EndTime.UTC().After(comparePeriod.GetEndTime().UTC()) {
		intersectPeriod.EndTime = comparePeriod.GetEndTime()
	}
	if timePeriod.EndTime.IsZero() {
		intersectPeriod.EndTime = comparePeriod.GetEndTime()
	}
	return intersectPeriod
}

func (timePeriod timePeriodImpl) Overlaps(other TimePeriod) (TimePeriod, bool) {
	if timePeriod.doesIntersect(other) {
		return timePeriod.intersect(other), true
	}
	return timePeriodImpl{}, false
}
