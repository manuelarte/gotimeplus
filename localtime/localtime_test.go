package localtime

import (
	"testing"
	"time"

	"github.com/manuelarte/gotimeplus/localdate"
)

func TestBefore(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		a, b LocalTime
		want bool
	}{
		"A is before B": {
			a:    New(10, 0, 0, 0),
			b:    New(11, 0, 0, 0),
			want: true,
		},
		"A is equal to B": {
			a:    New(15, 30, 45, 123),
			b:    New(15, 30, 45, 123),
			want: false,
		},
		"A is after B": {
			a:    New(23, 59, 59, 999999),
			b:    New(22, 0, 0, 0),
			want: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := test.a.Before(test.b)
			if got != test.want {
				t.Errorf("Before = %v, want %v", got, test.want)
			}
		})
	}
}

func TestAfter(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		a, b   LocalTime
		expect bool
	}{
		"A is after B": {
			a:      New(20, 0, 0, 0),
			b:      New(19, 59, 59, 999),
			expect: true,
		},
		"A is equal to B": {
			a:      New(5, 5, 5, 5),
			b:      New(5, 5, 5, 5),
			expect: false,
		},
		"A is before B": {
			a:      New(0, 0, 1, 0),
			b:      New(0, 0, 2, 0),
			expect: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := test.a.After(test.b)
			if got != test.expect {
				t.Errorf("After = %v, want %v", got, test.expect)
			}
		})
	}
}

func TestEqual(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		a, b LocalTime
		want bool
	}{
		"Equal": {
			a:    New(12, 0, 0, 0),
			b:    New(12, 0, 0, 0),
			want: true,
		},
		"Different hour": {
			a:    New(10, 0, 0, 0),
			b:    New(11, 0, 0, 0),
			want: false,
		},
		"Different nanosecond": {
			a:    New(12, 0, 0, 1),
			b:    New(12, 0, 0, 0),
			want: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := test.a.Equal(test.b)
			if got != test.want {
				t.Errorf("Equal = %v, want %v", got, test.want)
			}
		})
	}
}

func TestToTime(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		localTime LocalTime
		localDate localdate.LocalDate
		location  *time.Location
		want      time.Time
	}{
		"UTC": {
			localTime: New(12, 34, 56, 789),
			localDate: localdate.New(2000, time.January, 1),
			location:  time.UTC,
			want:      time.Date(2000, time.January, 1, 12, 34, 56, 789, time.UTC), // fixed dummy date
		},
		"Offset zone": {
			localTime: New(1, 2, 3, 4),
			localDate: localdate.New(2000, time.January, 1),
			location:  time.FixedZone("Test", -5*60*60),
			want:      time.Date(2000, 1, 1, 1, 2, 3, 4, time.FixedZone("Test", -5*60*60)),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := test.localTime.ToTime(test.localDate, test.location)
			if !got.Equal(test.want) {
				t.Errorf("ToTime = %v, want %v", got, test.want)
			}
		})
	}
}
