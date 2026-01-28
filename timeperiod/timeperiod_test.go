package timeperiod

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestGetDuration(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		timePeriod TimePeriod
		expected   time.Duration
	}{
		"One hour": {
			timePeriod: Must(
				ptr(time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC)),
				ptr(time.Date(2022, time.January, 1, 13, 0, 0, 0, time.UTC)),
			),
			expected: 1 * time.Hour,
		},
		"One day": {
			timePeriod: Must(
				ptr(time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC)),
				ptr(time.Date(2022, time.January, 2, 12, 0, 0, 0, time.UTC)),
			),
			expected: 24 * time.Hour,
		},
		"One month Jan 2022": {
			timePeriod: Must(
				yearMonthDay(2022, time.January, 1),
				yearMonthDay(2022, time.February, 1),
			),
			expected: 24 * time.Hour * 31,
		},
		"Less than one hour": {
			timePeriod: Must(
				ptr(time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC)),
				ptr(time.Date(2022, time.January, 1, 12, 59, 0, 0, time.UTC)),
			),
			expected: 59 * time.Minute,
		},
		"No end, max duration": {
			timePeriod: Must(
				ptr(time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC)),
				nil,
			),
			expected: 1<<63 - 1,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			actual := test.timePeriod.Duration()

			if test.expected != actual {
				t.Errorf("\nExpected: %v\nActual: %v", test.expected, actual)
			}
		})
	}
}

func TestDoesIntersect(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		basePeriod     TimePeriod
		comparePeriod  TimePeriod
		expectedResult bool
	}{
		"Base Period is exactly the same as Compare Period": {
			basePeriod: Must(
				yearMonthDay(2023, time.February, 1),
				yearMonthDay(2023, time.March, 1),
			),
			comparePeriod: Must(
				yearMonthDay(2023, time.February, 1),
				yearMonthDay(2023, time.March, 1),
			),
			expectedResult: true,
		},
		"Base Period falls inside Compare Period": {
			basePeriod: Must(
				yearMonthDay(2023, time.February, 1),
				yearMonthDay(2023, time.March, 1),
			),
			comparePeriod: Must(
				yearMonthDay(2023, time.January, 1),
				yearMonthDay(2023, time.April, 1),
			),
			expectedResult: true,
		},
		"Base Period contains Compare Period": {
			basePeriod: Must(
				yearMonthDay(2023, time.February, 1),
				yearMonthDay(2023, time.March, 1),
			),
			comparePeriod: Must(
				yearMonthDay(2023, time.February, 10),
				yearMonthDay(2023, time.February, 20),
			),
			expectedResult: true,
		},
		"Base Period overlaps first part of Compare Period": {
			basePeriod: Must(
				yearMonthDay(2023, time.February, 1),
				yearMonthDay(2023, time.March, 1),
			),
			comparePeriod: Must(
				yearMonthDay(2023, time.February, 15),
				yearMonthDay(2023, time.March, 15),
			),
			expectedResult: true,
		},
		"Base Period overlaps last part of Compare Period": {
			basePeriod: Must(
				yearMonthDay(2023, time.February, 1),
				yearMonthDay(2023, time.March, 1),
			),
			comparePeriod: Must(
				yearMonthDay(2023, time.January, 15),
				yearMonthDay(2023, time.February, 15),
			),
			expectedResult: true,
		},
		"Base Period is before Compare Period": {
			basePeriod: Must(
				yearMonthDay(2023, time.February, 1),
				yearMonthDay(2023, time.March, 1),
			),
			comparePeriod: Must(
				yearMonthDay(2023, time.April, 1),
				yearMonthDay(2023, time.May, 1),
			),
			expectedResult: false,
		},
		"Base Period is after Compare Period": {
			basePeriod: Must(
				yearMonthDay(2023, time.February, 1),
				yearMonthDay(2023, time.March, 1),
			),
			comparePeriod: Must(
				yearMonthDay(2023, time.January, 1),
				yearMonthDay(2023, time.January, 20),
			),
			expectedResult: false,
		},
		"Base Period ends on start of Compare Period": {
			basePeriod: Must(
				yearMonthDay(2023, time.February, 1),
				yearMonthDay(2023, time.March, 1),
			),
			comparePeriod: Must(
				yearMonthDay(2023, time.March, 1),
				yearMonthDay(2023, time.April, 1),
			),
			expectedResult: false,
		},
		"Base Period starts on end of Compare Period": {
			basePeriod: Must(
				yearMonthDay(2023, time.February, 1),
				yearMonthDay(2023, time.March, 1),
			),
			comparePeriod: Must(
				yearMonthDay(2023, time.January, 1),
				yearMonthDay(2023, time.February, 1),
			),
			expectedResult: false,
		},
		"Base Period has no end time and starts before Compare Period": {
			basePeriod: Must(
				yearMonthDay(2023, time.February, 1),
				nil,
			),
			comparePeriod: Must(
				yearMonthDay(2023, time.March, 1),
				yearMonthDay(2023, time.April, 1),
			),
			expectedResult: true,
		},
		"Base Period has no end time and starts on Compare Period start": {
			basePeriod: Must(
				yearMonthDay(2023, time.February, 1),
				nil,
			),
			comparePeriod: Must(
				yearMonthDay(2023, time.February, 1),
				yearMonthDay(2023, time.March, 1),
			),
			expectedResult: true,
		},
		"Base Period has no end time and starts inside Compare Period": {
			basePeriod: Must(
				yearMonthDay(2023, time.February, 1),
				nil,
			),
			comparePeriod: Must(
				yearMonthDay(2023, time.January, 1),
				yearMonthDay(2023, time.March, 1),
			),
			expectedResult: true,
		},
		"Base Period has no end time and starts on Compare Period end": {
			basePeriod: Must(
				yearMonthDay(2023, time.February, 1),
				nil,
			),
			comparePeriod: Must(
				yearMonthDay(2023, time.January, 1),
				yearMonthDay(2023, time.February, 1),
			),
			expectedResult: false,
		},
		"Base Period has no end time and starts after Compare Period end": {
			basePeriod: Must(
				yearMonthDay(2023, time.March, 1),
				nil,
			),
			comparePeriod: Must(
				yearMonthDay(2023, time.January, 1),
				yearMonthDay(2023, time.February, 1),
			),
			expectedResult: false,
		},
		"Base Period and Compare Period have no end times and Base starts before Compare start": {
			basePeriod: Must(
				yearMonthDay(2023, time.February, 1),
				nil,
			),
			comparePeriod: Must(
				yearMonthDay(2023, time.March, 1),
				nil,
			),
			expectedResult: true,
		},
		"Base Period and Compare Period have no end times and Base starts on Compare start": {
			basePeriod: Must(
				yearMonthDay(2023, time.February, 1),
				nil,
			),
			comparePeriod: Must(
				yearMonthDay(2023, time.February, 1),
				nil,
			),
			expectedResult: true,
		},
		"Base Period and Compare Period have no end times and Base starts after Compare start": {
			basePeriod: Must(
				yearMonthDay(2023, time.February, 1),
				nil,
			),
			comparePeriod: Must(
				yearMonthDay(2023, time.January, 1),
				nil,
			),
			expectedResult: true,
		},
		"Compare Period has no start time and ends before Base Period": {
			basePeriod: Must(
				yearMonthDay(2023, time.February, 1),
				yearMonthDay(2023, time.March, 1),
			),
			comparePeriod: Must(
				nil,
				yearMonthDay(2023, time.January, 1),
			),
			expectedResult: false,
		},
		"Compare Period has no start time and ends between Base Period": {
			basePeriod: Must(
				yearMonthDay(2023, time.February, 1),
				yearMonthDay(2023, time.March, 1),
			),
			comparePeriod: Must(
				nil,
				yearMonthDay(2023, time.February, 15),
			),
			expectedResult: true,
		},
		"Compare Period has no end time and starts before Base Period": {
			basePeriod: Must(
				yearMonthDay(2023, time.February, 1),
				yearMonthDay(2023, time.March, 1),
			),
			comparePeriod: Must(
				yearMonthDay(2023, time.January, 1),
				nil,
			),
			expectedResult: true,
		},
		"Compare Period has no end time and starts on Base Period start": {
			basePeriod: Must(
				yearMonthDay(2023, time.February, 1),
				yearMonthDay(2023, time.March, 1),
			),
			comparePeriod: Must(
				yearMonthDay(2023, time.February, 1),
				nil,
			),
			expectedResult: true,
		},
		"Compare Period has no end time and starts inside Base Period": {
			basePeriod: Must(
				yearMonthDay(2023, time.February, 1),
				yearMonthDay(2023, time.March, 1),
			),
			comparePeriod: Must(
				yearMonthDay(2023, time.February, 15),
				nil,
			),
			expectedResult: true,
		},
		"Compare Period has no end time and starts on Base Period end": {
			basePeriod: Must(
				yearMonthDay(2023, time.February, 1),
				yearMonthDay(2023, time.March, 1),
			),
			comparePeriod: Must(
				yearMonthDay(2023, time.March, 1),
				nil,
			),
			expectedResult: false,
		},
		"Compare Period has no end time and starts after Base Period end": {
			basePeriod: Must(
				yearMonthDay(2023, time.February, 1),
				yearMonthDay(2023, time.March, 1),
			),
			comparePeriod: Must(
				yearMonthDay(2023, time.April, 1),
				nil,
			),
			expectedResult: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			_, actualResult := test.basePeriod.Overlaps(test.comparePeriod)
			if test.expectedResult != actualResult {
				t.Errorf("Test %v: Expected %v but got %v", name, test.expectedResult, actualResult)
			}
		})
	}
}

func TestIntersect(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		basePeriod     TimePeriod
		comparePeriod  TimePeriod
		expectedPeriod TimePeriod
		expectedOk     bool
	}{
		"Base Period is exactly the same as Compare Period": {
			basePeriod: Must(
				yearMonthDay(2023, time.February, 1),
				yearMonthDay(2023, time.March, 1),
			),
			comparePeriod: Must(
				yearMonthDay(2023, time.February, 1),
				yearMonthDay(2023, time.March, 1),
			),
			expectedPeriod: Must(
				yearMonthDay(2023, time.February, 1),
				yearMonthDay(2023, time.March, 1),
			),
			expectedOk: true,
		},
		"Base Period falls inside Compare Period": {
			basePeriod: Must(
				yearMonthDay(2023, time.February, 1),
				yearMonthDay(2023, time.March, 1),
			),
			comparePeriod: Must(
				yearMonthDay(2023, time.January, 1),
				yearMonthDay(2023, time.April, 1),
			),
			expectedPeriod: Must(
				yearMonthDay(2023, time.February, 1),
				yearMonthDay(2023, time.March, 1),
			),
			expectedOk: true,
		},
		"Base Period contains Compare Period": {
			basePeriod: Must(
				yearMonthDay(2023, time.February, 1),
				yearMonthDay(2023, time.March, 1),
			),
			comparePeriod: Must(
				yearMonthDay(2023, time.February, 10),
				yearMonthDay(2023, time.February, 20),
			),
			expectedPeriod: Must(
				yearMonthDay(2023, time.February, 10),
				yearMonthDay(2023, time.February, 20),
			),
			expectedOk: true,
		},
		"Base Period overlaps first part of Compare Period": {
			basePeriod: Must(
				yearMonthDay(2023, time.February, 1),
				yearMonthDay(2023, time.March, 1),
			),
			comparePeriod: Must(
				yearMonthDay(2023, time.February, 15),
				yearMonthDay(2023, time.March, 15),
			),
			expectedPeriod: Must(
				yearMonthDay(2023, time.February, 15),
				yearMonthDay(2023, time.March, 1),
			),
			expectedOk: true,
		},
		"Base Period overlaps last part of Compare Period": {
			basePeriod: Must(
				yearMonthDay(2023, time.February, 1),
				yearMonthDay(2023, time.March, 1),
			),
			comparePeriod: Must(
				yearMonthDay(2023, time.January, 15),
				yearMonthDay(2023, time.February, 15),
			),
			expectedPeriod: Must(
				yearMonthDay(2023, time.February, 1),
				yearMonthDay(2023, time.February, 15),
			),
			expectedOk: true,
		},
		"Base Period is before Compare Period": {
			basePeriod: Must(
				yearMonthDay(2023, time.February, 1),
				yearMonthDay(2023, time.March, 1),
			),
			comparePeriod: Must(
				yearMonthDay(2023, time.April, 1),
				yearMonthDay(2023, 5, 1),
			),
			expectedPeriod: Must(nil, nil),
			expectedOk:     false,
		},
		"Base Period is after Compare Period": {
			basePeriod: Must(
				yearMonthDay(2023, time.February, 1),
				yearMonthDay(2023, time.March, 1),
			),
			comparePeriod: Must(
				yearMonthDay(2023, time.January, 1),
				yearMonthDay(2023, time.January, 20),
			),
			expectedPeriod: Must(nil, nil),
			expectedOk:     false,
		},
		"Base Period ends on start of Compare Period": {
			basePeriod: Must(
				yearMonthDay(2023, time.February, 1),
				yearMonthDay(2023, time.March, 1),
			),
			comparePeriod: Must(
				yearMonthDay(2023, time.March, 1),
				yearMonthDay(2023, time.April, 1),
			),
			expectedPeriod: Must(nil, nil),
			expectedOk:     false,
		},
		"Base Period starts on end of Compare Period": {
			basePeriod: Must(
				yearMonthDay(2023, time.February, 1),
				yearMonthDay(2023, time.March, 1),
			),
			comparePeriod: Must(
				yearMonthDay(2023, time.January, 1),
				yearMonthDay(2023, time.February, 1),
			),
			expectedPeriod: Must(nil, nil),
			expectedOk:     false,
		},
		"Base Period no limits and compare period with limits": {
			basePeriod: Infinite,
			comparePeriod: Must(
				yearMonthDay(2023, time.January, 1),
				yearMonthDay(2023, time.February, 1),
			),
			expectedPeriod: Must(
				yearMonthDay(2023, time.January, 1),
				yearMonthDay(2023, time.February, 1),
			),
			expectedOk: true,
		},
		"Base Period has no end time and starts before Compare Period": {
			basePeriod: Must(
				yearMonthDay(2023, time.February, 1),
				nil,
			),
			comparePeriod: Must(
				yearMonthDay(2023, time.March, 1),
				yearMonthDay(2023, time.April, 1),
			),
			expectedPeriod: Must(
				yearMonthDay(2023, time.March, 1),
				yearMonthDay(2023, time.April, 1),
			),
			expectedOk: true,
		},
		"Base Period has no end time and starts on Compare Period start": {
			basePeriod: Must(
				yearMonthDay(2023, time.February, 1),
				nil,
			),
			comparePeriod: Must(
				yearMonthDay(2023, time.February, 1),
				yearMonthDay(2023, time.March, 1),
			),
			expectedPeriod: Must(
				yearMonthDay(2023, time.February, 1),
				yearMonthDay(2023, time.March, 1),
			),
			expectedOk: true,
		},
		"Base Period has no end time and starts inside Compare Period": {
			basePeriod: Must(
				yearMonthDay(2023, time.February, 1),
				nil,
			),
			comparePeriod: Must(
				yearMonthDay(2023, time.January, 1),
				yearMonthDay(2023, time.March, 1),
			),
			expectedPeriod: Must(
				yearMonthDay(2023, time.February, 1),
				yearMonthDay(2023, time.March, 1),
			),
			expectedOk: true,
		},
		"Base Period has no end time and starts on Compare Period end": {
			basePeriod: Must(
				yearMonthDay(2023, time.February, 1),
				nil,
			),
			comparePeriod: Must(
				yearMonthDay(2023, time.January, 1),
				yearMonthDay(2023, time.February, 1),
			),
			expectedPeriod: Must(nil, nil),
			expectedOk:     false,
		},
		"Base Period has no end time and starts after Compare Period end": {
			basePeriod: Must(
				yearMonthDay(2023, time.March, 1),
				nil,
			),
			comparePeriod: Must(
				yearMonthDay(2023, time.January, 1),
				yearMonthDay(2023, time.February, 1),
			),
			expectedPeriod: Must(nil, nil),
			expectedOk:     false,
		},
		"Base Period and Compare Period have no end times and Base starts before Compare start": {
			basePeriod: Must(
				yearMonthDay(2023, time.February, 1),
				nil,
			),
			comparePeriod: Must(
				yearMonthDay(2023, time.March, 1),
				nil,
			),
			expectedPeriod: Must(
				yearMonthDay(2023, time.March, 1),
				nil,
			),
			expectedOk: true,
		},
		"Base Period and Compare Period have no end times and Base starts on Compare start": {
			basePeriod: Must(
				yearMonthDay(2023, time.February, 1),
				nil,
			),
			comparePeriod: Must(
				yearMonthDay(2023, time.February, 1),
				nil,
			),
			expectedPeriod: Must(
				yearMonthDay(2023, time.February, 1),
				nil,
			),
			expectedOk: true,
		},
		"Base Period and Compare Period have no end times and Base starts after Compare start": {
			basePeriod: Must(
				yearMonthDay(2023, time.February, 1),
				nil,
			),
			comparePeriod: Must(
				yearMonthDay(2023, time.January, 1),
				nil,
			),
			expectedPeriod: Must(
				yearMonthDay(2023, time.February, 1),
				nil,
			),
			expectedOk: true,
		},
		"Compare Period has no end time and starts before Base Period": {
			basePeriod: Must(
				yearMonthDay(2023, time.February, 1),
				yearMonthDay(2023, time.March, 1),
			),
			comparePeriod: Must(
				yearMonthDay(2023, time.January, 1),
				nil,
			),
			expectedPeriod: Must(
				yearMonthDay(2023, time.February, 1),
				yearMonthDay(2023, time.March, 1),
			),
			expectedOk: true,
		},
		"Compare Period has no end time and starts on Base Period start": {
			basePeriod: Must(
				yearMonthDay(2023, time.February, 1),
				yearMonthDay(2023, time.March, 1),
			),
			comparePeriod: Must(
				yearMonthDay(2023, time.February, 1),
				nil,
			),
			expectedPeriod: Must(
				yearMonthDay(2023, time.February, 1),
				yearMonthDay(2023, time.March, 1),
			),
			expectedOk: true,
		},
		"Compare Period has no end time and starts inside Base Period": {
			basePeriod: Must(
				yearMonthDay(2023, time.February, 1),
				yearMonthDay(2023, time.March, 1),
			),
			comparePeriod: Must(
				yearMonthDay(2023, time.February, 15),
				nil,
			),
			expectedPeriod: Must(
				yearMonthDay(2023, time.February, 15),
				yearMonthDay(2023, time.March, 1),
			),
			expectedOk: true,
		},
		"Compare Period has no start time and ends on Base Period end": {
			basePeriod: Must(
				yearMonthDay(2023, time.February, 1),
				yearMonthDay(2023, time.March, 1),
			),
			comparePeriod: Must(
				nil,
				yearMonthDay(2023, time.March, 1),
			),
			expectedPeriod: Must(
				yearMonthDay(2023, time.February, 1),
				yearMonthDay(2023, time.March, 1),
			),
			expectedOk: true,
		},
		"Compare Period has no start time and ends before Base Period": {
			basePeriod: Must(
				yearMonthDay(2023, time.February, 1),
				yearMonthDay(2023, time.March, 1),
			),
			comparePeriod: Must(
				nil,
				yearMonthDay(2023, time.January, 1),
			),
			expectedPeriod: Must(nil, nil),
			expectedOk:     false,
		},
		"Compare Period has no start time and ends during Base Period": {
			basePeriod: Must(
				yearMonthDay(2023, time.February, 1),
				yearMonthDay(2023, time.March, 1),
			),
			comparePeriod: Must(
				nil,
				yearMonthDay(2023, time.February, 15),
			),
			expectedPeriod: Must(
				yearMonthDay(2023, time.February, 1),
				yearMonthDay(2023, time.February, 15),
			),
			expectedOk: true,
		},
		"Compare Period has no end time and starts on Base Period end": {
			basePeriod: Must(
				yearMonthDay(2023, time.February, 1),
				yearMonthDay(2023, time.March, 1),
			),
			comparePeriod: Must(
				yearMonthDay(2023, time.March, 1),
				nil,
			),
			expectedPeriod: Must(nil, nil),
			expectedOk:     false,
		},
		"Compare Period has no end time and starts after Base Period end": {
			basePeriod: Must(
				yearMonthDay(2023, time.February, 1),
				yearMonthDay(2023, time.March, 1),
			),
			comparePeriod: Must(
				yearMonthDay(2023, time.April, 1),
				nil,
			),
			expectedPeriod: Must(nil, nil),
			expectedOk:     false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			actualResult, ok := test.basePeriod.Overlaps(test.comparePeriod)
			if ok != test.expectedOk {
				t.Fatalf("Expected: %v, Actual: %v", test.expectedOk, ok)
			}

			if !cmp.Equal(test.expectedPeriod, actualResult, cmp.AllowUnexported(startTimeEndTimePeriod{})) {
				t.Errorf("Diff: %s", cmp.Diff(test.expectedPeriod, actualResult, cmp.AllowUnexported(startTimeEndTimePeriod{})))
			}
		})
	}
}

func TestNew(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		startTime   *time.Time
		endTime     *time.Time
		expectedErr error
	}{
		"endTime after startTime": {
			startTime:   ptr(time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC)),
			endTime:     ptr(time.Date(2022, time.January, 1, 13, 0, 0, 0, time.UTC)),
			expectedErr: nil,
		},
		"endTime before startTime": {
			startTime:   ptr(time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC)),
			endTime:     ptr(time.Date(2022, time.January, 1, 10, 0, 0, 0, time.UTC)),
			expectedErr: ErrEndTimeBeforeStartTime,
		},
		"endTime equal to startTime": {
			startTime:   ptr(time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC)),
			endTime:     ptr(time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC)),
			expectedErr: nil,
		},
		"No endTime": {
			startTime:   ptr(time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC)),
			endTime:     nil,
			expectedErr: nil,
		},
		"No startTime": {
			startTime:   nil,
			endTime:     ptr(time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC)),
			expectedErr: nil,
		},
		"No Start nor endTime": {
			startTime:   nil,
			endTime:     nil,
			expectedErr: nil,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			_, err := New(test.startTime, test.endTime)
			if !errors.Is(test.expectedErr, err) {
				t.Errorf("Expected: %v, Actual: %v", test.expectedErr, err)
			}
		})
	}
}

func TestMust(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		startTime     *time.Time
		endTime       *time.Time
		expectedPanic bool
	}{
		"endTime after startTime": {
			startTime:     ptr(time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC)),
			endTime:       ptr(time.Date(2022, time.January, 1, 13, 0, 0, 0, time.UTC)),
			expectedPanic: false,
		},
		"endTime before startTime": {
			startTime:     ptr(time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC)),
			endTime:       ptr(time.Date(2022, time.January, 1, 10, 0, 0, 0, time.UTC)),
			expectedPanic: true,
		},
		"endTime equal to startTime": {
			startTime:     ptr(time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC)),
			endTime:       ptr(time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC)),
			expectedPanic: false,
		},
		"No endTime": {
			startTime:     ptr(time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC)),
			endTime:       nil,
			expectedPanic: false,
		},
		"No startTime": {
			startTime:     nil,
			endTime:       ptr(time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC)),
			expectedPanic: false,
		},
		"No Start nor endTime": {
			startTime:     nil,
			endTime:       nil,
			expectedPanic: false,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			defer func() {
				if recover() != nil {
					if !test.expectedPanic {
						t.Errorf("Panic not expected")
					}
				}
			}()

			t.Parallel()

			_ = Must(test.startTime, test.endTime)
		})
	}
}

func TestInputDatesSameValueAsGetters(t *testing.T) {
	t.Parallel()

	startTime := yearMonthDay(2026, time.January, 1)
	endTime := yearMonthDay(2026, time.February, 1)

	timePeriod := Must(startTime, endTime)
	if !cmp.Equal(startTime, timePeriod.StartTime()) {
		t.Errorf("timePeriod.StartTime should be have the same value as original input")
	}

	if !cmp.Equal(endTime, timePeriod.EndTime()) {
		t.Errorf("timePeriod.EndTime should be have the same value as original input")
	}
}

func TestStartTimeInputReassignedDoesNotAffect(t *testing.T) {
	t.Parallel()

	startTime := yearMonthDay(2026, time.January, 1)
	timePeriod := Must(startTime, yearMonthDay(2026, time.February, 1))
	startTime = ptr(startTime.Add(60 * time.Hour))

	if reflect.DeepEqual(startTime, timePeriod.StartTime()) {
		t.Errorf("original input has been modified, and then it should not be reflected in timePeriod.StartTime")
	}
}

func yearMonthDay(year int, month time.Month, day int) *time.Time {
	return ptr(time.Date(year, month, day, 0, 0, 0, 0, time.UTC))
}

func ptr[T any](t T) *T {
	return &t
}
