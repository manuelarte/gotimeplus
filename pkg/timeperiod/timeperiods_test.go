package timeperiod

import (
	"errors"
	"testing"
	"time"
)

func TestTimePeriod_GetDuration(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		timePeriod TimePeriod
		expected   time.Duration
	}{
		"One hour": {
			timePeriod: TimePeriod{
				startTime: time.Date(2022, 1, 1, 12, 0, 0, 0, time.UTC),
				endTime:   time.Date(2022, 1, 1, 13, 0, 0, 0, time.UTC),
			},
			expected: 1 * time.Hour,
		},
		"One day": {
			timePeriod: TimePeriod{
				startTime: time.Date(2022, 1, 1, 12, 0, 0, 0, time.UTC),
				endTime:   time.Date(2022, 1, 2, 12, 0, 0, 0, time.UTC),
			},
			expected: 24 * time.Hour,
		},
		"No end, max duration": {
			timePeriod: TimePeriod{
				startTime: time.Date(2022, 1, 1, 12, 0, 0, 0, time.UTC),
				endTime:   time.Time{},
			},
			expected: 1<<63 - 1,
		},
		"Less than one hour": {
			timePeriod: TimePeriod{
				startTime: time.Date(2022, 1, 1, 12, 0, 0, 0, time.UTC),
				endTime:   time.Date(2022, 1, 1, 12, 59, 0, 0, time.UTC),
			},
			expected: 59 * time.Minute,
		},
	}
	for name, test := range tests {
		name, test := name, test
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			actual := test.timePeriod.GetDuration()

			if test.expected != actual {
				t.Errorf("\nExpected: %v\nActual: %v", test.expected, actual)
			}
		})
	}
}

func TestTimePeriod_DoesIntersect(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		basePeriod     TimePeriod
		comparePeriod  TimePeriod
		expectedResult bool
	}{
		"Base Period is exactly the same as Compare Period": {
			basePeriod: TimePeriod{
				startTime: yearMonthDay(2023, 2, 1),
				endTime:   yearMonthDay(2023, 3, 1),
			},
			comparePeriod: TimePeriod{
				startTime: yearMonthDay(2023, 2, 1),
				endTime:   yearMonthDay(2023, 3, 1),
			},
			expectedResult: true,
		},
		"Base Period falls inside Compare Period": {
			basePeriod: TimePeriod{
				startTime: yearMonthDay(2023, 2, 1),
				endTime:   yearMonthDay(2023, 3, 1),
			},
			comparePeriod: TimePeriod{
				startTime: yearMonthDay(2023, 1, 1),
				endTime:   yearMonthDay(2023, 4, 1),
			},
			expectedResult: true,
		},
		"Base Period contains Compare Period": {
			basePeriod: TimePeriod{
				startTime: yearMonthDay(2023, 2, 1),
				endTime:   yearMonthDay(2023, 3, 1),
			},
			comparePeriod: TimePeriod{
				startTime: yearMonthDay(2023, 2, 10),
				endTime:   yearMonthDay(2023, 2, 20),
			},
			expectedResult: true,
		},
		"Base Period overlaps first part of Compare Period": {
			basePeriod: TimePeriod{
				startTime: yearMonthDay(2023, 2, 1),
				endTime:   yearMonthDay(2023, 3, 1),
			},
			comparePeriod: TimePeriod{
				startTime: yearMonthDay(2023, 2, 15),
				endTime:   yearMonthDay(2023, 3, 15),
			},
			expectedResult: true,
		},
		"Base Period overlaps last part of Compare Period": {
			basePeriod: TimePeriod{
				startTime: yearMonthDay(2023, 2, 1),
				endTime:   yearMonthDay(2023, 3, 1),
			},
			comparePeriod: TimePeriod{
				startTime: yearMonthDay(2023, 1, 15),
				endTime:   yearMonthDay(2023, 2, 15),
			},
			expectedResult: true,
		},
		"Base Period is before Compare Period": {
			basePeriod: TimePeriod{
				startTime: yearMonthDay(2023, 2, 1),
				endTime:   yearMonthDay(2023, 3, 1),
			},
			comparePeriod: TimePeriod{
				startTime: yearMonthDay(2023, 4, 1),
				endTime:   yearMonthDay(2023, 5, 1),
			},
			expectedResult: false,
		},
		"Base Period is after Compare Period": {
			basePeriod: TimePeriod{
				startTime: yearMonthDay(2023, 2, 1),
				endTime:   yearMonthDay(2023, 3, 1),
			},
			comparePeriod: TimePeriod{
				startTime: yearMonthDay(2023, 1, 1),
				endTime:   yearMonthDay(2023, 1, 20),
			},
			expectedResult: false,
		},
		"Base Period ends on start of Compare Period": {
			basePeriod: TimePeriod{
				startTime: yearMonthDay(2023, 2, 1),
				endTime:   yearMonthDay(2023, 3, 1),
			},
			comparePeriod: TimePeriod{
				startTime: yearMonthDay(2023, 3, 1),
				endTime:   yearMonthDay(2023, 4, 1),
			},
			expectedResult: false,
		},
		"Base Period starts on end of Compare Period": {
			basePeriod: TimePeriod{
				startTime: yearMonthDay(2023, 2, 1),
				endTime:   yearMonthDay(2023, 3, 1),
			},
			comparePeriod: TimePeriod{
				startTime: yearMonthDay(2023, 1, 1),
				endTime:   yearMonthDay(2023, 2, 1),
			},
			expectedResult: false,
		},
		"Compare Period has no end time and starts before Base Period": {
			basePeriod: TimePeriod{
				startTime: yearMonthDay(2023, 2, 1),
				endTime:   yearMonthDay(2023, 3, 1),
			},
			comparePeriod: TimePeriod{
				startTime: yearMonthDay(2023, 1, 1),
				endTime:   time.Time{},
			},
			expectedResult: true,
		},
		"Compare Period has no end time and starts on Base Period start": {
			basePeriod: TimePeriod{
				startTime: yearMonthDay(2023, 2, 1),
				endTime:   yearMonthDay(2023, 3, 1),
			},
			comparePeriod: TimePeriod{
				startTime: yearMonthDay(2023, 2, 1),
				endTime:   time.Time{},
			},
			expectedResult: true,
		},
		"Compare Period has no end time and starts inside Base Period": {
			basePeriod: TimePeriod{
				startTime: yearMonthDay(2023, 2, 1),
				endTime:   yearMonthDay(2023, 3, 1),
			},
			comparePeriod: TimePeriod{
				startTime: yearMonthDay(2023, 2, 15),
				endTime:   time.Time{},
			},
			expectedResult: true,
		},
		"Compare Period has no end time and starts on Base Period end": {
			basePeriod: TimePeriod{
				startTime: yearMonthDay(2023, 2, 1),
				endTime:   yearMonthDay(2023, 3, 1),
			},
			comparePeriod: TimePeriod{
				startTime: yearMonthDay(2023, 3, 1),
				endTime:   time.Time{},
			},
			expectedResult: false,
		},
		"Compare Period has no end time and starts after Base Period end": {
			basePeriod: TimePeriod{
				startTime: yearMonthDay(2023, 2, 1),
				endTime:   yearMonthDay(2023, 3, 1),
			},
			comparePeriod: TimePeriod{
				startTime: yearMonthDay(2023, 4, 1),
				endTime:   time.Time{},
			},
			expectedResult: false,
		},
		"Base Period has no end time and starts before Compare Period": {
			basePeriod: TimePeriod{
				startTime: yearMonthDay(2023, 2, 1),
				endTime:   time.Time{},
			},
			comparePeriod: TimePeriod{
				startTime: yearMonthDay(2023, 3, 1),
				endTime:   yearMonthDay(2023, 4, 1),
			},
			expectedResult: true,
		},
		"Base Period has no end time and starts on Compare Period start": {
			basePeriod: TimePeriod{
				startTime: yearMonthDay(2023, 2, 1),
				endTime:   time.Time{},
			},
			comparePeriod: TimePeriod{
				startTime: yearMonthDay(2023, 2, 1),
				endTime:   yearMonthDay(2023, 3, 1),
			},
			expectedResult: true,
		},
		"Base Period has no end time and starts inside Compare Period": {
			basePeriod: TimePeriod{
				startTime: yearMonthDay(2023, 2, 1),
				endTime:   time.Time{},
			},
			comparePeriod: TimePeriod{
				startTime: yearMonthDay(2023, 1, 1),
				endTime:   yearMonthDay(2023, 3, 1),
			},
			expectedResult: true,
		},
		"Base Period has no end time and starts on Compare Period end": {
			basePeriod: TimePeriod{
				startTime: yearMonthDay(2023, 2, 1),
				endTime:   time.Time{},
			},
			comparePeriod: TimePeriod{
				startTime: yearMonthDay(2023, 1, 1),
				endTime:   yearMonthDay(2023, 2, 1),
			},
			expectedResult: false,
		},
		"Base Period has no end time and starts after Compare Period end": {
			basePeriod: TimePeriod{
				startTime: yearMonthDay(2023, 3, 1),
				endTime:   time.Time{},
			},
			comparePeriod: TimePeriod{
				startTime: yearMonthDay(2023, 1, 1),
				endTime:   yearMonthDay(2023, 2, 1),
			},
			expectedResult: false,
		},
		"Base Period and Compare Period have no end times and Base starts before Compare start": {
			basePeriod: TimePeriod{
				startTime: yearMonthDay(2023, 2, 1),
				endTime:   time.Time{},
			},
			comparePeriod: TimePeriod{
				startTime: yearMonthDay(2023, 3, 1),
				endTime:   time.Time{},
			},
			expectedResult: true,
		},
		"Base Period and Compare Period have no end times and Base starts on Compare start": {
			basePeriod: TimePeriod{
				startTime: yearMonthDay(2023, 2, 1),
				endTime:   time.Time{},
			},
			comparePeriod: TimePeriod{
				startTime: yearMonthDay(2023, 2, 1),
				endTime:   time.Time{},
			},
			expectedResult: true,
		},
		"Base Period and Compare Period have no end times and Base starts after Compare start": {
			basePeriod: TimePeriod{
				startTime: yearMonthDay(2023, 2, 1),
				endTime:   time.Time{},
			},
			comparePeriod: TimePeriod{
				startTime: yearMonthDay(2023, 1, 1),
				endTime:   time.Time{},
			},
			expectedResult: true,
		},
	}

	for name, test := range tests {
		name, test := name, test
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			actualResult := test.basePeriod.doesIntersect(test.comparePeriod)
			if test.expectedResult != actualResult {
				t.Errorf("Test %v: Expected %v but got %v", name, test.expectedResult, actualResult)
			}
		})
	}
}

func TestTimePeriod_Intersect(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		basePeriod     TimePeriod
		comparePeriod  TimePeriod
		expectedPeriod TimePeriod
	}{
		"Base Period is exactly the same as Compare Period": {
			basePeriod: TimePeriod{
				startTime: yearMonthDay(2023, 2, 1),
				endTime:   yearMonthDay(2023, 3, 1),
			},
			comparePeriod: TimePeriod{
				startTime: yearMonthDay(2023, 2, 1),
				endTime:   yearMonthDay(2023, 3, 1),
			},
			expectedPeriod: TimePeriod{
				startTime: yearMonthDay(2023, 2, 1),
				endTime:   yearMonthDay(2023, 3, 1),
			},
		},
		"Base Period falls inside Compare Period": {
			basePeriod: TimePeriod{
				startTime: yearMonthDay(2023, 2, 1),
				endTime:   yearMonthDay(2023, 3, 1),
			},
			comparePeriod: TimePeriod{
				startTime: yearMonthDay(2023, 1, 1),
				endTime:   yearMonthDay(2023, 4, 1),
			},
			expectedPeriod: TimePeriod{
				startTime: yearMonthDay(2023, 2, 1),
				endTime:   yearMonthDay(2023, 3, 1),
			},
		},
		"Base Period contains Compare Period": {
			basePeriod: TimePeriod{
				startTime: yearMonthDay(2023, 2, 1),
				endTime:   yearMonthDay(2023, 3, 1),
			},
			comparePeriod: TimePeriod{
				startTime: yearMonthDay(2023, 2, 10),
				endTime:   yearMonthDay(2023, 2, 20),
			},
			expectedPeriod: TimePeriod{
				startTime: yearMonthDay(2023, 2, 10),
				endTime:   yearMonthDay(2023, 2, 20),
			},
		},
		"Base Period overlaps first part of Compare Period": {
			basePeriod: TimePeriod{
				startTime: yearMonthDay(2023, 2, 1),
				endTime:   yearMonthDay(2023, 3, 1),
			},
			comparePeriod: TimePeriod{
				startTime: yearMonthDay(2023, 2, 15),
				endTime:   yearMonthDay(2023, 3, 15),
			},
			expectedPeriod: TimePeriod{
				startTime: yearMonthDay(2023, 2, 15),
				endTime:   yearMonthDay(2023, 3, 1),
			},
		},
		"Base Period overlaps last part of Compare Period": {
			basePeriod: TimePeriod{
				startTime: yearMonthDay(2023, 2, 1),
				endTime:   yearMonthDay(2023, 3, 1),
			},
			comparePeriod: TimePeriod{
				startTime: yearMonthDay(2023, 1, 15),
				endTime:   yearMonthDay(2023, 2, 15),
			},
			expectedPeriod: TimePeriod{
				startTime: yearMonthDay(2023, 2, 1),
				endTime:   yearMonthDay(2023, 2, 15),
			},
		},
		"Base Period is before Compare Period": {
			basePeriod: TimePeriod{
				startTime: yearMonthDay(2023, 2, 1),
				endTime:   yearMonthDay(2023, 3, 1),
			},
			comparePeriod: TimePeriod{
				startTime: yearMonthDay(2023, 4, 1),
				endTime:   yearMonthDay(2023, 5, 1),
			},
			expectedPeriod: TimePeriod{
				startTime: time.Time{},
				endTime:   time.Time{},
			},
		},
		"Base Period is after Compare Period": {
			basePeriod: TimePeriod{
				startTime: yearMonthDay(2023, 2, 1),
				endTime:   yearMonthDay(2023, 3, 1),
			},
			comparePeriod: TimePeriod{
				startTime: yearMonthDay(2023, 1, 1),
				endTime:   yearMonthDay(2023, 1, 20),
			},
			expectedPeriod: TimePeriod{
				startTime: time.Time{},
				endTime:   time.Time{},
			},
		},
		"Base Period ends on start of Compare Period": {
			basePeriod: TimePeriod{
				startTime: yearMonthDay(2023, 2, 1),
				endTime:   yearMonthDay(2023, 3, 1),
			},
			comparePeriod: TimePeriod{
				startTime: yearMonthDay(2023, 3, 1),
				endTime:   yearMonthDay(2023, 4, 1),
			},
			expectedPeriod: TimePeriod{
				startTime: time.Time{},
				endTime:   time.Time{},
			},
		},
		"Base Period starts on end of Compare Period": {
			basePeriod: TimePeriod{
				startTime: yearMonthDay(2023, 2, 1),
				endTime:   yearMonthDay(2023, 3, 1),
			},
			comparePeriod: TimePeriod{
				startTime: yearMonthDay(2023, 1, 1),
				endTime:   yearMonthDay(2023, 2, 1),
			},
			expectedPeriod: TimePeriod{
				startTime: time.Time{},
				endTime:   time.Time{},
			},
		},
		"Compare Period has no end time and starts before Base Period": {
			basePeriod: TimePeriod{
				startTime: yearMonthDay(2023, 2, 1),
				endTime:   yearMonthDay(2023, 3, 1),
			},
			comparePeriod: TimePeriod{
				startTime: yearMonthDay(2023, 1, 1),
				endTime:   time.Time{},
			},
			expectedPeriod: TimePeriod{
				startTime: yearMonthDay(2023, 2, 1),
				endTime:   yearMonthDay(2023, 3, 1),
			},
		},
		"Compare Period has no end time and starts on Base Period start": {
			basePeriod: TimePeriod{
				startTime: yearMonthDay(2023, 2, 1),
				endTime:   yearMonthDay(2023, 3, 1),
			},
			comparePeriod: TimePeriod{
				startTime: yearMonthDay(2023, 2, 1),
				endTime:   time.Time{},
			},
			expectedPeriod: TimePeriod{
				startTime: yearMonthDay(2023, 2, 1),
				endTime:   yearMonthDay(2023, 3, 1),
			},
		},
		"Compare Period has no end time and starts inside Base Period": {
			basePeriod: TimePeriod{
				startTime: yearMonthDay(2023, 2, 1),
				endTime:   yearMonthDay(2023, 3, 1),
			},
			comparePeriod: TimePeriod{
				startTime: yearMonthDay(2023, 2, 15),
				endTime:   time.Time{},
			},
			expectedPeriod: TimePeriod{
				startTime: yearMonthDay(2023, 2, 15),
				endTime:   yearMonthDay(2023, 3, 1),
			},
		},
		"Compare Period has no end time and starts on Base Period end": {
			basePeriod: TimePeriod{
				startTime: yearMonthDay(2023, 2, 1),
				endTime:   yearMonthDay(2023, 3, 1),
			},
			comparePeriod: TimePeriod{
				startTime: yearMonthDay(2023, 3, 1),
				endTime:   time.Time{},
			},
			expectedPeriod: TimePeriod{
				startTime: time.Time{},
				endTime:   time.Time{},
			},
		},
		"Compare Period has no end time and starts after Base Period end": {
			basePeriod: TimePeriod{
				startTime: yearMonthDay(2023, 2, 1),
				endTime:   yearMonthDay(2023, 3, 1),
			},
			comparePeriod: TimePeriod{
				startTime: yearMonthDay(2023, 4, 1),
				endTime:   time.Time{},
			},
			expectedPeriod: TimePeriod{
				startTime: time.Time{},
				endTime:   time.Time{},
			},
		},
		"Base Period has no end time and starts before Compare Period": {
			basePeriod: TimePeriod{
				startTime: yearMonthDay(2023, 2, 1),
				endTime:   time.Time{},
			},
			comparePeriod: TimePeriod{
				startTime: yearMonthDay(2023, 3, 1),
				endTime:   yearMonthDay(2023, 4, 1),
			},
			expectedPeriod: TimePeriod{
				startTime: yearMonthDay(2023, 3, 1),
				endTime:   yearMonthDay(2023, 4, 1),
			},
		},
		"Base Period has no end time and starts on Compare Period start": {
			basePeriod: TimePeriod{
				startTime: yearMonthDay(2023, 2, 1),
				endTime:   time.Time{},
			},
			comparePeriod: TimePeriod{
				startTime: yearMonthDay(2023, 2, 1),
				endTime:   yearMonthDay(2023, 3, 1),
			},
			expectedPeriod: TimePeriod{
				startTime: yearMonthDay(2023, 2, 1),
				endTime:   yearMonthDay(2023, 3, 1),
			},
		},
		"Base Period has no end time and starts inside Compare Period": {
			basePeriod: TimePeriod{
				startTime: yearMonthDay(2023, 2, 1),
				endTime:   time.Time{},
			},
			comparePeriod: TimePeriod{
				startTime: yearMonthDay(2023, 1, 1),
				endTime:   yearMonthDay(2023, 3, 1),
			},
			expectedPeriod: TimePeriod{
				startTime: yearMonthDay(2023, 2, 1),
				endTime:   yearMonthDay(2023, 3, 1),
			},
		},
		"Base Period has no end time and starts on Compare Period end": {
			basePeriod: TimePeriod{
				startTime: yearMonthDay(2023, 2, 1),
				endTime:   time.Time{},
			},
			comparePeriod: TimePeriod{
				startTime: yearMonthDay(2023, 1, 1),
				endTime:   yearMonthDay(2023, 2, 1),
			},
			expectedPeriod: TimePeriod{
				startTime: time.Time{},
				endTime:   time.Time{},
			},
		},
		"Base Period has no end time and starts after Compare Period end": {
			basePeriod: TimePeriod{
				startTime: yearMonthDay(2023, 3, 1),
				endTime:   time.Time{},
			},
			comparePeriod: TimePeriod{
				startTime: yearMonthDay(2023, 1, 1),
				endTime:   yearMonthDay(2023, 2, 1),
			},
			expectedPeriod: TimePeriod{
				startTime: time.Time{},
				endTime:   time.Time{},
			},
		},
		"Base Period and Compare Period have no end times and Base starts before Compare start": {
			basePeriod: TimePeriod{
				startTime: yearMonthDay(2023, 2, 1),
				endTime:   time.Time{},
			},
			comparePeriod: TimePeriod{
				startTime: yearMonthDay(2023, 3, 1),
				endTime:   time.Time{},
			},
			expectedPeriod: TimePeriod{
				startTime: yearMonthDay(2023, 3, 1),
				endTime:   time.Time{},
			},
		},
		"Base Period and Compare Period have no end times and Base starts on Compare start": {
			basePeriod: TimePeriod{
				startTime: yearMonthDay(2023, 2, 1),
				endTime:   time.Time{},
			},
			comparePeriod: TimePeriod{
				startTime: yearMonthDay(2023, 2, 1),
				endTime:   time.Time{},
			},
			expectedPeriod: TimePeriod{
				startTime: yearMonthDay(2023, 2, 1),
				endTime:   time.Time{},
			},
		},
		"Base Period and Compare Period have no end times and Base starts after Compare start": {
			basePeriod: TimePeriod{
				startTime: yearMonthDay(2023, 2, 1),
				endTime:   time.Time{},
			},
			comparePeriod: TimePeriod{
				startTime: yearMonthDay(2023, 1, 1),
				endTime:   time.Time{},
			},
			expectedPeriod: TimePeriod{
				startTime: yearMonthDay(2023, 2, 1),
				endTime:   time.Time{},
			},
		},
	}

	for name, test := range tests {
		name, test := name, test
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			actualResult := test.basePeriod.intersect(test.comparePeriod)
			if test.expectedPeriod != actualResult {
				t.Errorf("Expected: %v, Actual: %v", test.expectedPeriod, actualResult)
			}
		})
	}
}

func TestTimePeriod_NewTimePeriod(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		startTime   time.Time
		endTime     time.Time
		expectedErr error
	}{
		"endTime after startTime": {
			startTime:   time.Date(2022, 1, 1, 12, 0, 0, 0, time.UTC),
			endTime:     time.Date(2022, 1, 1, 13, 0, 0, 0, time.UTC),
			expectedErr: nil,
		},
		"endTime before startTime": {
			startTime:   time.Date(2022, 1, 1, 12, 0, 0, 0, time.UTC),
			endTime:     time.Date(2022, 1, 1, 10, 0, 0, 0, time.UTC),
			expectedErr: ErrEndTimeBeforeStartTime,
		},
		"endTime equal to startTime": {
			startTime:   time.Date(2022, 1, 1, 12, 0, 0, 0, time.UTC),
			endTime:     time.Date(2022, 1, 1, 12, 0, 0, 0, time.UTC),
			expectedErr: nil,
		},
		"No endTime": {
			startTime:   time.Date(2022, 1, 1, 12, 0, 0, 0, time.UTC),
			endTime:     time.Time{},
			expectedErr: nil,
		},
		"No startTime": {
			startTime:   time.Time{},
			endTime:     time.Date(2022, 1, 1, 12, 0, 0, 0, time.UTC),
			expectedErr: nil,
		},
		"No Start nor endTime": {
			startTime:   time.Time{},
			endTime:     time.Time{},
			expectedErr: nil,
		},
	}
	for name, test := range tests {
		name, test := name, test
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			_, err := NewTimePeriod(test.startTime, test.endTime)
			if !errors.Is(test.expectedErr, err) {
				t.Errorf("Expected: %v, Actual: %v", test.expectedErr, err)
			}
		})
	}
}

func yearMonthDay(year int, month time.Month, day int) time.Time {
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}
