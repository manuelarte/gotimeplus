package timeperiod

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTimePeriod_GetDuration(t *testing.T) {
	// Arrange
	tests := map[string]struct {
		timePeriod    TimePeriod
		expectedHours time.Duration
	}{
		"One hour": {
			timePeriod: timePeriodImpl{
				StartTime: time.Date(2022, 1, 1, 12, 0, 0, 0, time.UTC),
				EndTime:   time.Date(2022, 1, 1, 13, 0, 0, 0, time.UTC),
			},
			expectedHours: 1 * time.Hour,
		},
		"One day": {
			timePeriod: timePeriodImpl{
				StartTime: time.Date(2022, 1, 1, 12, 0, 0, 0, time.UTC),
				EndTime:   time.Date(2022, 1, 2, 12, 0, 0, 0, time.UTC),
			},
			expectedHours: 24 * time.Hour,
		},
		"No end, max duration": {
			timePeriod: timePeriodImpl{
				StartTime: time.Date(2022, 1, 1, 12, 0, 0, 0, time.UTC),
			},
			expectedHours: 1<<63 - 1,
		},
		"Less than one hour": {
			timePeriod: timePeriodImpl{
				StartTime: time.Date(2022, 1, 1, 12, 0, 0, 0, time.UTC),
				EndTime:   time.Date(2022, 1, 1, 12, 59, 0, 0, time.UTC),
			},
			expectedHours: 59 * time.Minute,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			// Act
			actualHours := test.timePeriod.GetDuration()

			// Assert
			assert.Equal(t, test.expectedHours, actualHours)
		})
	}
}

func TestTimePeriod_DoesIntersect(t *testing.T) {
	// Arrange
	tests := map[string]struct {
		basePeriod     timePeriodImpl
		comparePeriod  TimePeriod
		expectedResult bool
	}{
		"Base Period is exactly the same as Compare Period": {
			basePeriod: timePeriodImpl{
				StartTime: time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC),
				EndTime:   time.Date(2023, 3, 1, 0, 0, 0, 0, time.UTC),
			},
			comparePeriod: timePeriodImpl{
				StartTime: time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC),
				EndTime:   time.Date(2023, 3, 1, 0, 0, 0, 0, time.UTC),
			},
			expectedResult: true,
		},
		"Base Period falls inside Compare Period": {
			basePeriod: timePeriodImpl{
				StartTime: time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC),
				EndTime:   time.Date(2023, 3, 1, 0, 0, 0, 0, time.UTC),
			},
			comparePeriod: timePeriodImpl{
				StartTime: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				EndTime:   time.Date(2023, 4, 1, 0, 0, 0, 0, time.UTC),
			},
			expectedResult: true,
		},
		"Base Period contains Compare Period": {
			basePeriod: timePeriodImpl{
				StartTime: time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC),
				EndTime:   time.Date(2023, 3, 1, 0, 0, 0, 0, time.UTC),
			},
			comparePeriod: timePeriodImpl{
				StartTime: time.Date(2023, 2, 10, 0, 0, 0, 0, time.UTC),
				EndTime:   time.Date(2023, 2, 20, 0, 0, 0, 0, time.UTC),
			},
			expectedResult: true,
		},
		"Base Period overlaps first part of Compare Period": {
			basePeriod: timePeriodImpl{
				StartTime: time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC),
				EndTime:   time.Date(2023, 3, 1, 0, 0, 0, 0, time.UTC),
			},
			comparePeriod: timePeriodImpl{
				StartTime: time.Date(2023, 2, 15, 0, 0, 0, 0, time.UTC),
				EndTime:   time.Date(2023, 3, 15, 0, 0, 0, 0, time.UTC),
			},
			expectedResult: true,
		},
		"Base Period overlaps last part of Compare Period": {
			basePeriod: timePeriodImpl{
				StartTime: time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC),
				EndTime:   time.Date(2023, 3, 1, 0, 0, 0, 0, time.UTC),
			},
			comparePeriod: timePeriodImpl{
				StartTime: time.Date(2023, 1, 15, 0, 0, 0, 0, time.UTC),
				EndTime:   time.Date(2023, 2, 15, 0, 0, 0, 0, time.UTC),
			},
			expectedResult: true,
		},
		"Base Period is before Compare Period": {
			basePeriod: timePeriodImpl{
				StartTime: time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC),
				EndTime:   time.Date(2023, 3, 1, 0, 0, 0, 0, time.UTC),
			},
			comparePeriod: timePeriodImpl{
				StartTime: time.Date(2023, 4, 1, 0, 0, 0, 0, time.UTC),
				EndTime:   (time.Date(2023, 5, 1, 0, 0, 0, 0, time.UTC)),
			},
			expectedResult: false,
		},
		"Base Period is after Compare Period": {
			basePeriod: timePeriodImpl{
				StartTime: time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC),
				EndTime:   (time.Date(2023, 3, 1, 0, 0, 0, 0, time.UTC)),
			},
			comparePeriod: timePeriodImpl{
				StartTime: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				EndTime:   (time.Date(2023, 1, 20, 0, 0, 0, 0, time.UTC)),
			},
			expectedResult: false,
		},
		"Base Period ends on start of Compare Period": {
			basePeriod: timePeriodImpl{
				StartTime: time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC),
				EndTime:   (time.Date(2023, 3, 1, 0, 0, 0, 0, time.UTC)),
			},
			comparePeriod: timePeriodImpl{
				StartTime: time.Date(2023, 3, 1, 0, 0, 0, 0, time.UTC),
				EndTime:   (time.Date(2023, 4, 1, 0, 0, 0, 0, time.UTC)),
			},
			expectedResult: false,
		},
		"Base Period starts on end of Compare Period": {
			basePeriod: timePeriodImpl{
				StartTime: time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC),
				EndTime:   (time.Date(2023, 3, 1, 0, 0, 0, 0, time.UTC)),
			},
			comparePeriod: timePeriodImpl{
				StartTime: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				EndTime:   (time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC)),
			},
			expectedResult: false,
		},
		"Compare Period has no end time and starts before Base Period": {
			basePeriod: timePeriodImpl{
				StartTime: time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC),
				EndTime:   (time.Date(2023, 3, 1, 0, 0, 0, 0, time.UTC)),
			},
			comparePeriod: timePeriodImpl{
				StartTime: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			expectedResult: true,
		},
		"Compare Period has no end time and starts on Base Period start": {
			basePeriod: timePeriodImpl{
				StartTime: time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC),
				EndTime:   (time.Date(2023, 3, 1, 0, 0, 0, 0, time.UTC)),
			},
			comparePeriod: timePeriodImpl{
				StartTime: time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC),
			},
			expectedResult: true,
		},
		"Compare Period has no end time and starts inside Base Period": {
			basePeriod: timePeriodImpl{
				StartTime: time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC),
				EndTime:   time.Date(2023, 3, 1, 0, 0, 0, 0, time.UTC),
			},
			comparePeriod: timePeriodImpl{
				StartTime: time.Date(2023, 2, 15, 0, 0, 0, 0, time.UTC),
			},
			expectedResult: true,
		},
		"Compare Period has no end time and starts on Base Period end": {
			basePeriod: timePeriodImpl{
				StartTime: time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC),
				EndTime:   time.Date(2023, 3, 1, 0, 0, 0, 0, time.UTC),
			},
			comparePeriod: timePeriodImpl{
				StartTime: time.Date(2023, 3, 1, 0, 0, 0, 0, time.UTC),
			},
			expectedResult: false,
		},
		"Compare Period has no end time and starts after Base Period end": {
			basePeriod: timePeriodImpl{
				StartTime: time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC),
				EndTime:   (time.Date(2023, 3, 1, 0, 0, 0, 0, time.UTC)),
			},
			comparePeriod: timePeriodImpl{
				StartTime: time.Date(2023, 4, 1, 0, 0, 0, 0, time.UTC),
			},
			expectedResult: false,
		},
		"Base Period has no end time and starts before Compare Period": {
			basePeriod: timePeriodImpl{
				StartTime: time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC),
			},
			comparePeriod: timePeriodImpl{
				StartTime: time.Date(2023, 3, 1, 0, 0, 0, 0, time.UTC),
				EndTime:   (time.Date(2023, 4, 1, 0, 0, 0, 0, time.UTC)),
			},
			expectedResult: true,
		},
		"Base Period has no end time and starts on Compare Period start": {
			basePeriod: timePeriodImpl{
				StartTime: time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC),
			},
			comparePeriod: timePeriodImpl{
				StartTime: time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC),
				EndTime:   time.Date(2023, 3, 1, 0, 0, 0, 0, time.UTC),
			},
			expectedResult: true,
		},
		"Base Period has no end time and starts inside Compare Period": {
			basePeriod: timePeriodImpl{
				StartTime: time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC),
			},
			comparePeriod: timePeriodImpl{
				StartTime: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				EndTime:   time.Date(2023, 3, 1, 0, 0, 0, 0, time.UTC),
			},
			expectedResult: true,
		},
		"Base Period has no end time and starts on Compare Period end": {
			basePeriod: timePeriodImpl{
				StartTime: time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC),
			},
			comparePeriod: timePeriodImpl{
				StartTime: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				EndTime:   (time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC)),
			},
			expectedResult: false,
		},
		"Base Period has no end time and starts after Compare Period end": {
			basePeriod: timePeriodImpl{
				StartTime: time.Date(2023, 3, 1, 0, 0, 0, 0, time.UTC),
			},
			comparePeriod: timePeriodImpl{
				StartTime: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				EndTime:   (time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC)),
			},
			expectedResult: false,
		},
		"Base Period and Compare Period have no end times and Base starts before Compare start": {
			basePeriod: timePeriodImpl{
				StartTime: time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC),
			},
			comparePeriod: timePeriodImpl{
				StartTime: time.Date(2023, 3, 1, 0, 0, 0, 0, time.UTC),
			},
			expectedResult: true,
		},
		"Base Period and Compare Period have no end times and Base starts on Compare start": {
			basePeriod: timePeriodImpl{
				StartTime: time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC),
			},
			comparePeriod: timePeriodImpl{
				StartTime: time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC),
			},
			expectedResult: true,
		},
		"Base Period and Compare Period have no end times and Base starts after Compare start": {
			basePeriod: timePeriodImpl{
				StartTime: time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC),
			},
			comparePeriod: timePeriodImpl{
				StartTime: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			expectedResult: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			// Act
			actualResult := test.basePeriod.doesIntersect(test.comparePeriod)

			// Assert
			assert.Equal(t, test.expectedResult, actualResult)
		})
	}
}

func TestTimePeriod_Intersect(t *testing.T) {
	// Arrange
	tests := map[string]struct {
		basePeriod     timePeriodImpl
		comparePeriod  TimePeriod
		expectedPeriod TimePeriod
	}{
		"Base Period is exactly the same as Compare Period": {
			basePeriod: timePeriodImpl{
				StartTime: time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC),
				EndTime:   time.Date(2023, 3, 1, 0, 0, 0, 0, time.UTC),
			},
			comparePeriod: timePeriodImpl{
				StartTime: time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC),
				EndTime:   (time.Date(2023, 3, 1, 0, 0, 0, 0, time.UTC)),
			},
			expectedPeriod: timePeriodImpl{
				StartTime: time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC),
				EndTime:   (time.Date(2023, 3, 1, 0, 0, 0, 0, time.UTC)),
			},
		},
		"Base Period falls inside Compare Period": {
			basePeriod: timePeriodImpl{
				StartTime: time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC),
				EndTime:   time.Date(2023, 3, 1, 0, 0, 0, 0, time.UTC),
			},
			comparePeriod: timePeriodImpl{
				StartTime: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				EndTime:   time.Date(2023, 4, 1, 0, 0, 0, 0, time.UTC),
			},
			expectedPeriod: timePeriodImpl{
				StartTime: time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC),
				EndTime:   time.Date(2023, 3, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		"Base Period contains Compare Period": {
			basePeriod: timePeriodImpl{
				StartTime: time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC),
				EndTime:   time.Date(2023, 3, 1, 0, 0, 0, 0, time.UTC),
			},
			comparePeriod: timePeriodImpl{
				StartTime: time.Date(2023, 2, 10, 0, 0, 0, 0, time.UTC),
				EndTime:   (time.Date(2023, 2, 20, 0, 0, 0, 0, time.UTC)),
			},
			expectedPeriod: timePeriodImpl{
				StartTime: time.Date(2023, 2, 10, 0, 0, 0, 0, time.UTC),
				EndTime:   time.Date(2023, 2, 20, 0, 0, 0, 0, time.UTC),
			},
		},
		"Base Period overlaps first part of Compare Period": {
			basePeriod: timePeriodImpl{
				StartTime: time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC),
				EndTime:   time.Date(2023, 3, 1, 0, 0, 0, 0, time.UTC),
			},
			comparePeriod: timePeriodImpl{
				StartTime: time.Date(2023, 2, 15, 0, 0, 0, 0, time.UTC),
				EndTime:   time.Date(2023, 3, 15, 0, 0, 0, 0, time.UTC),
			},
			expectedPeriod: timePeriodImpl{
				StartTime: time.Date(2023, 2, 15, 0, 0, 0, 0, time.UTC),
				EndTime:   time.Date(2023, 3, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		"Base Period overlaps last part of Compare Period": {
			basePeriod: timePeriodImpl{
				StartTime: time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC),
				EndTime:   time.Date(2023, 3, 1, 0, 0, 0, 0, time.UTC),
			},
			comparePeriod: timePeriodImpl{
				StartTime: time.Date(2023, 1, 15, 0, 0, 0, 0, time.UTC),
				EndTime:   time.Date(2023, 2, 15, 0, 0, 0, 0, time.UTC),
			},
			expectedPeriod: timePeriodImpl{
				StartTime: time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC),
				EndTime:   time.Date(2023, 2, 15, 0, 0, 0, 0, time.UTC),
			},
		},
		"Base Period is before Compare Period": {
			basePeriod: timePeriodImpl{
				StartTime: time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC),
				EndTime:   time.Date(2023, 3, 1, 0, 0, 0, 0, time.UTC),
			},
			comparePeriod: timePeriodImpl{
				StartTime: time.Date(2023, 4, 1, 0, 0, 0, 0, time.UTC),
				EndTime:   time.Date(2023, 5, 1, 0, 0, 0, 0, time.UTC),
			},
			expectedPeriod: timePeriodImpl{
				StartTime: time.Time{},
			},
		},
		"Base Period is after Compare Period": {
			basePeriod: timePeriodImpl{
				StartTime: time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC),
				EndTime:   time.Date(2023, 3, 1, 0, 0, 0, 0, time.UTC),
			},
			comparePeriod: timePeriodImpl{
				StartTime: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				EndTime:   time.Date(2023, 1, 20, 0, 0, 0, 0, time.UTC),
			},
			expectedPeriod: timePeriodImpl{
				StartTime: time.Time{},
			},
		},
		"Base Period ends on start of Compare Period": {
			basePeriod: timePeriodImpl{
				StartTime: time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC),
				EndTime:   time.Date(2023, 3, 1, 0, 0, 0, 0, time.UTC),
			},
			comparePeriod: timePeriodImpl{
				StartTime: time.Date(2023, 3, 1, 0, 0, 0, 0, time.UTC),
				EndTime:   time.Date(2023, 4, 1, 0, 0, 0, 0, time.UTC),
			},
			expectedPeriod: timePeriodImpl{
				StartTime: time.Time{},
			},
		},
		"Base Period starts on end of Compare Period": {
			basePeriod: timePeriodImpl{
				StartTime: time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC),
				EndTime:   time.Date(2023, 3, 1, 0, 0, 0, 0, time.UTC),
			},
			comparePeriod: timePeriodImpl{
				StartTime: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				EndTime:   time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC),
			},
			expectedPeriod: timePeriodImpl{
				StartTime: time.Time{},
			},
		},
		"Compare Period has no end time and starts before Base Period": {
			basePeriod: timePeriodImpl{
				StartTime: time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC),
				EndTime:   (time.Date(2023, 3, 1, 0, 0, 0, 0, time.UTC)),
			},
			comparePeriod: timePeriodImpl{
				StartTime: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			expectedPeriod: timePeriodImpl{
				StartTime: time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC),
				EndTime:   (time.Date(2023, 3, 1, 0, 0, 0, 0, time.UTC)),
			},
		},
		"Compare Period has no end time and starts on Base Period start": {
			basePeriod: timePeriodImpl{
				StartTime: time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC),
				EndTime:   (time.Date(2023, 3, 1, 0, 0, 0, 0, time.UTC)),
			},
			comparePeriod: timePeriodImpl{
				StartTime: time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC),
			},
			expectedPeriod: timePeriodImpl{
				StartTime: time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC),
				EndTime:   (time.Date(2023, 3, 1, 0, 0, 0, 0, time.UTC)),
			},
		},
		"Compare Period has no end time and starts inside Base Period": {
			basePeriod: timePeriodImpl{
				StartTime: time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC),
				EndTime:   (time.Date(2023, 3, 1, 0, 0, 0, 0, time.UTC)),
			},
			comparePeriod: timePeriodImpl{
				StartTime: time.Date(2023, 2, 15, 0, 0, 0, 0, time.UTC),
			},
			expectedPeriod: timePeriodImpl{
				StartTime: time.Date(2023, 2, 15, 0, 0, 0, 0, time.UTC),
				EndTime:   (time.Date(2023, 3, 1, 0, 0, 0, 0, time.UTC)),
			},
		},
		"Compare Period has no end time and starts on Base Period end": {
			basePeriod: timePeriodImpl{
				StartTime: time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC),
				EndTime:   (time.Date(2023, 3, 1, 0, 0, 0, 0, time.UTC)),
			},
			comparePeriod: timePeriodImpl{
				StartTime: time.Date(2023, 3, 1, 0, 0, 0, 0, time.UTC),
			},
			expectedPeriod: timePeriodImpl{
				StartTime: time.Time{},
			},
		},
		"Compare Period has no end time and starts after Base Period end": {
			basePeriod: timePeriodImpl{
				StartTime: time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC),
				EndTime:   (time.Date(2023, 3, 1, 0, 0, 0, 0, time.UTC)),
			},
			comparePeriod: timePeriodImpl{
				StartTime: time.Date(2023, 4, 1, 0, 0, 0, 0, time.UTC),
			},
			expectedPeriod: timePeriodImpl{
				StartTime: time.Time{},
			},
		},
		"Base Period has no end time and starts before Compare Period": {
			basePeriod: timePeriodImpl{
				StartTime: time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC),
			},
			comparePeriod: timePeriodImpl{
				StartTime: time.Date(2023, 3, 1, 0, 0, 0, 0, time.UTC),
				EndTime:   time.Date(2023, 4, 1, 0, 0, 0, 0, time.UTC),
			},
			expectedPeriod: timePeriodImpl{
				StartTime: time.Date(2023, 3, 1, 0, 0, 0, 0, time.UTC),
				EndTime:   time.Date(2023, 4, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		"Base Period has no end time and starts on Compare Period start": {
			basePeriod: timePeriodImpl{
				StartTime: time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC),
			},
			comparePeriod: timePeriodImpl{
				StartTime: time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC),
				EndTime:   (time.Date(2023, 3, 1, 0, 0, 0, 0, time.UTC)),
			},
			expectedPeriod: timePeriodImpl{
				StartTime: time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC),
				EndTime:   (time.Date(2023, 3, 1, 0, 0, 0, 0, time.UTC)),
			},
		},
		"Base Period has no end time and starts inside Compare Period": {
			basePeriod: timePeriodImpl{
				StartTime: time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC),
			},
			comparePeriod: timePeriodImpl{
				StartTime: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				EndTime:   time.Date(2023, 3, 1, 0, 0, 0, 0, time.UTC),
			},
			expectedPeriod: timePeriodImpl{
				StartTime: time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC),
				EndTime:   (time.Date(2023, 3, 1, 0, 0, 0, 0, time.UTC)),
			},
		},
		"Base Period has no end time and starts on Compare Period end": {
			basePeriod: timePeriodImpl{
				StartTime: time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC),
			},
			comparePeriod: timePeriodImpl{
				StartTime: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				EndTime:   (time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC)),
			},
			expectedPeriod: timePeriodImpl{
				StartTime: time.Time{},
			},
		},
		"Base Period has no end time and starts after Compare Period end": {
			basePeriod: timePeriodImpl{
				StartTime: time.Date(2023, 3, 1, 0, 0, 0, 0, time.UTC),
			},
			comparePeriod: timePeriodImpl{
				StartTime: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				EndTime:   (time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC)),
			},
			expectedPeriod: timePeriodImpl{
				StartTime: time.Time{},
			},
		},
		"Base Period and Compare Period have no end times and Base starts before Compare start": {
			basePeriod: timePeriodImpl{
				StartTime: time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC),
			},
			comparePeriod: timePeriodImpl{
				StartTime: time.Date(2023, 3, 1, 0, 0, 0, 0, time.UTC),
			},
			expectedPeriod: timePeriodImpl{
				StartTime: time.Date(2023, 3, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		"Base Period and Compare Period have no end times and Base starts on Compare start": {
			basePeriod: timePeriodImpl{
				StartTime: time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC),
			},
			comparePeriod: timePeriodImpl{
				StartTime: time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC),
			},
			expectedPeriod: timePeriodImpl{
				StartTime: time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		"Base Period and Compare Period have no end times and Base starts after Compare start": {
			basePeriod: timePeriodImpl{
				StartTime: time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC),
			},
			comparePeriod: timePeriodImpl{
				StartTime: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			expectedPeriod: timePeriodImpl{
				StartTime: time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC),
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			// Act
			actualResult := test.basePeriod.intersect(test.comparePeriod)

			// Assert
			assert.Equal(t, test.expectedPeriod, actualResult)
		})
	}
}

func TestTimePeriod_NewTimePeriod(t *testing.T) {
	// Arrange
	tests := map[string]struct {
		startTime   time.Time
		endTime     time.Time
		expectedErr error
	}{
		"EndTime after StartTime": {
			startTime:   time.Date(2022, 1, 1, 12, 0, 0, 0, time.UTC),
			endTime:     time.Date(2022, 1, 1, 13, 0, 0, 0, time.UTC),
			expectedErr: nil,
		},
		"EndTime before StartTime": {
			startTime:   time.Date(2022, 1, 1, 12, 0, 0, 0, time.UTC),
			endTime:     time.Date(2022, 1, 1, 10, 0, 0, 0, time.UTC),
			expectedErr: ErrEndTimeBeforeStartTime,
		},
		"EndTime equal to StartTime": {
			startTime:   time.Date(2022, 1, 1, 12, 0, 0, 0, time.UTC),
			endTime:     time.Date(2022, 1, 1, 12, 0, 0, 0, time.UTC),
			expectedErr: nil,
		},
		"No EndTime": {
			startTime:   time.Date(2022, 1, 1, 12, 0, 0, 0, time.UTC),
			endTime:     time.Time{},
			expectedErr: nil,
		},
		"No StartTime": {
			startTime:   time.Time{},
			endTime:     time.Date(2022, 1, 1, 12, 0, 0, 0, time.UTC),
			expectedErr: nil,
		},
		"No Start nor EndTime": {
			startTime:   time.Time{},
			endTime:     time.Time{},
			expectedErr: nil,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			// Act
			_, err := NewTimePeriod(test.startTime, test.endTime)

			// Assert
			assert.Equal(t, test.expectedErr, err)
		})
	}
}
