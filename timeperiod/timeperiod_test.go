package timeperiod

import (
	"errors"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestGetDuration(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		timePeriod TimePeriod
		want       time.Duration
	}{
		"One hour": {
			timePeriod: Must(
				ptr(time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC)),
				ptr(time.Date(2022, time.January, 1, 13, 0, 0, 0, time.UTC)),
			),
			want: 1 * time.Hour,
		},
		"One day": {
			timePeriod: Must(
				ptr(time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC)),
				ptr(time.Date(2022, time.January, 2, 12, 0, 0, 0, time.UTC)),
			),
			want: 24 * time.Hour,
		},
		"One month Jan 2022": {
			timePeriod: Must(
				yearMonthDay(2022, time.January, 1),
				yearMonthDay(2022, time.February, 1),
			),
			want: 24 * time.Hour * 31,
		},
		"Less than one hour": {
			timePeriod: Must(
				ptr(time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC)),
				ptr(time.Date(2022, time.January, 1, 12, 59, 0, 0, time.UTC)),
			),
			want: 59 * time.Minute,
		},
		"No end, max duration": {
			timePeriod: Must(
				ptr(time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC)),
				nil,
			),
			want: 1<<63 - 1,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := test.timePeriod.Duration()

			if test.want != got {
				t.Errorf("Duration() = %v, wantOk %v", got, test.want)
			}
		})
	}
}

func TestDoesIntersect(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		basePeriod    TimePeriod
		comparePeriod TimePeriod
		want          bool
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
			want: true,
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
			want: true,
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
			want: true,
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
			want: true,
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
			want: true,
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
			want: false,
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
			want: false,
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
			want: false,
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
			want: false,
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
			want: true,
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
			want: true,
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
			want: true,
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
			want: false,
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
			want: false,
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
			want: true,
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
			want: true,
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
			want: true,
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
			want: false,
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
			want: true,
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
			want: true,
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
			want: true,
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
			want: true,
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
			want: false,
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
			want: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			_, got := test.basePeriod.Overlaps(test.comparePeriod)
			if test.want != got {
				t.Errorf("Overlaps = %v, wantOk %v", got, test.want)
			}
		})
	}
}

func TestIntersect(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		basePeriod    TimePeriod
		comparePeriod TimePeriod
		wantPeriod    TimePeriod
		wantOk        bool
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
			wantPeriod: Must(
				yearMonthDay(2023, time.February, 1),
				yearMonthDay(2023, time.March, 1),
			),
			wantOk: true,
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
			wantPeriod: Must(
				yearMonthDay(2023, time.February, 1),
				yearMonthDay(2023, time.March, 1),
			),
			wantOk: true,
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
			wantPeriod: Must(
				yearMonthDay(2023, time.February, 10),
				yearMonthDay(2023, time.February, 20),
			),
			wantOk: true,
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
			wantPeriod: Must(
				yearMonthDay(2023, time.February, 15),
				yearMonthDay(2023, time.March, 1),
			),
			wantOk: true,
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
			wantPeriod: Must(
				yearMonthDay(2023, time.February, 1),
				yearMonthDay(2023, time.February, 15),
			),
			wantOk: true,
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
			wantPeriod: Must(nil, nil),
			wantOk:     false,
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
			wantPeriod: Must(nil, nil),
			wantOk:     false,
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
			wantPeriod: Must(nil, nil),
			wantOk:     false,
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
			wantPeriod: Must(nil, nil),
			wantOk:     false,
		},
		"Base Period no limits and compare period with limits": {
			basePeriod: Infinite,
			comparePeriod: Must(
				yearMonthDay(2023, time.January, 1),
				yearMonthDay(2023, time.February, 1),
			),
			wantPeriod: Must(
				yearMonthDay(2023, time.January, 1),
				yearMonthDay(2023, time.February, 1),
			),
			wantOk: true,
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
			wantPeriod: Must(
				yearMonthDay(2023, time.March, 1),
				yearMonthDay(2023, time.April, 1),
			),
			wantOk: true,
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
			wantPeriod: Must(
				yearMonthDay(2023, time.February, 1),
				yearMonthDay(2023, time.March, 1),
			),
			wantOk: true,
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
			wantPeriod: Must(
				yearMonthDay(2023, time.February, 1),
				yearMonthDay(2023, time.March, 1),
			),
			wantOk: true,
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
			wantPeriod: Must(nil, nil),
			wantOk:     false,
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
			wantPeriod: Must(nil, nil),
			wantOk:     false,
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
			wantPeriod: Must(
				yearMonthDay(2023, time.March, 1),
				nil,
			),
			wantOk: true,
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
			wantPeriod: Must(
				yearMonthDay(2023, time.February, 1),
				nil,
			),
			wantOk: true,
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
			wantPeriod: Must(
				yearMonthDay(2023, time.February, 1),
				nil,
			),
			wantOk: true,
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
			wantPeriod: Must(
				yearMonthDay(2023, time.February, 1),
				yearMonthDay(2023, time.March, 1),
			),
			wantOk: true,
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
			wantPeriod: Must(
				yearMonthDay(2023, time.February, 1),
				yearMonthDay(2023, time.March, 1),
			),
			wantOk: true,
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
			wantPeriod: Must(
				yearMonthDay(2023, time.February, 15),
				yearMonthDay(2023, time.March, 1),
			),
			wantOk: true,
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
			wantPeriod: Must(
				yearMonthDay(2023, time.February, 1),
				yearMonthDay(2023, time.March, 1),
			),
			wantOk: true,
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
			wantPeriod: Must(nil, nil),
			wantOk:     false,
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
			wantPeriod: Must(
				yearMonthDay(2023, time.February, 1),
				yearMonthDay(2023, time.February, 15),
			),
			wantOk: true,
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
			wantPeriod: Must(nil, nil),
			wantOk:     false,
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
			wantPeriod: Must(nil, nil),
			wantOk:     false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			gotPeriod, gotOk := test.basePeriod.Overlaps(test.comparePeriod)
			if gotOk != test.wantOk {
				t.Fatalf("Expected: %v, Actual: %v", test.wantOk, gotOk)
			}

			if diff := cmp.Diff(test.wantPeriod, gotPeriod, cmp.AllowUnexported(startTimeEndTimePeriod{})); diff != "" {
				t.Errorf("Overlaps() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestNew(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		startTime *time.Time
		endTime   *time.Time
		want      error
	}{
		"endTime after startTime": {
			startTime: ptr(time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC)),
			endTime:   ptr(time.Date(2022, time.January, 1, 13, 0, 0, 0, time.UTC)),
			want:      nil,
		},
		"endTime before startTime": {
			startTime: ptr(time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC)),
			endTime:   ptr(time.Date(2022, time.January, 1, 10, 0, 0, 0, time.UTC)),
			want:      ErrEndTimeBeforeStartTime,
		},
		"endTime equal to startTime": {
			startTime: ptr(time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC)),
			endTime:   ptr(time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC)),
			want:      nil,
		},
		"No endTime": {
			startTime: ptr(time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC)),
			endTime:   nil,
			want:      nil,
		},
		"No startTime": {
			startTime: nil,
			endTime:   ptr(time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC)),
			want:      nil,
		},
		"No Start nor endTime": {
			startTime: nil,
			endTime:   nil,
			want:      nil,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			_, err := New(test.startTime, test.endTime)
			if !errors.Is(test.want, err) {
				t.Errorf("New() = %v, wantOk: %v", err, test.want)
			}
		})
	}
}

func TestMust(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		startTime *time.Time
		endTime   *time.Time
		isPanic   bool
	}{
		"endTime after startTime": {
			startTime: ptr(time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC)),
			endTime:   ptr(time.Date(2022, time.January, 1, 13, 0, 0, 0, time.UTC)),
			isPanic:   false,
		},
		"endTime before startTime": {
			startTime: ptr(time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC)),
			endTime:   ptr(time.Date(2022, time.January, 1, 10, 0, 0, 0, time.UTC)),
			isPanic:   true,
		},
		"endTime equal to startTime": {
			startTime: ptr(time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC)),
			endTime:   ptr(time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC)),
			isPanic:   false,
		},
		"No endTime": {
			startTime: ptr(time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC)),
			endTime:   nil,
			isPanic:   false,
		},
		"No startTime": {
			startTime: nil,
			endTime:   ptr(time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC)),
			isPanic:   false,
		},
		"No Start nor endTime": {
			startTime: nil,
			endTime:   nil,
			isPanic:   false,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			defer func() {
				if recover() != nil {
					if !test.isPanic {
						t.Errorf("Panic not wantOk")
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
	if diff := cmp.Diff(startTime, timePeriod.StartTime()); diff != "" {
		t.Errorf("timePeriod.StartTime() (-wantOk +got):\n%s", diff)
	}

	if diff := cmp.Diff(endTime, timePeriod.EndTime()); diff != "" {
		t.Errorf("timePeriod.EndTime() (-wantOk +got):\n%s", diff)
	}
}

func TestStartTimeInputReassignedDoesNotAffect(t *testing.T) {
	t.Parallel()

	startTime := yearMonthDay(2026, time.January, 1)
	timePeriod := Must(startTime, yearMonthDay(2026, time.February, 1))
	startTime = ptr(startTime.Add(60 * time.Hour))

	if cmp.Equal(startTime, timePeriod.StartTime()) {
		t.Errorf("original input has been modified, and then it should not be reflected in timePeriod.StartTime")
	}
}

func yearMonthDay(year int, month time.Month, day int) *time.Time {
	return ptr(time.Date(year, month, day, 0, 0, 0, 0, time.UTC))
}

func ptr[T any](t T) *T {
	return &t
}
