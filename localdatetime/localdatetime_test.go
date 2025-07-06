package localdatetime

import (
	"testing"
	"time"

	"github.com/manuelarte/gotimeplus/localdate"
	"github.com/manuelarte/gotimeplus/localtime"
)

func TestLocalDateTime_NewFromEqual(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		localDateTime LocalDateTime
		other         LocalDateTime
		expected      bool
	}{
		"same local date time": {
			localDateTime: NewFrom(
				localdate.New(2000, 1, 1),
				localtime.New(0, 0, 0, 0),
			),
			other: NewFrom(
				localdate.New(2000, 1, 1),
				localtime.New(0, 0, 0, 0),
			),
			expected: true,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			actual := test.localDateTime.Equal(test.other)

			if test.expected != actual {
				t.Errorf("\nExpected: %v\nActual: %v", test.expected, actual)
			}
		})
	}
}

func TestLocalDateTime_Before(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		a, b   LocalDateTime
		expect bool
	}{
		"A is before B - different date": {
			a:      New(2023, time.December, 31, 23, 59, 59, 0),
			b:      New(2024, time.January, 1, 0, 0, 0, 0),
			expect: true,
		},
		"A is before B - same date, earlier time": {
			a:      New(2024, time.July, 5, 10, 0, 0, 0),
			b:      New(2024, time.July, 5, 11, 0, 0, 0),
			expect: true,
		},
		"A is equal to B": {
			a:      New(2024, time.July, 5, 15, 45, 30, 0),
			b:      New(2024, time.July, 5, 15, 45, 30, 0),
			expect: false,
		},
		"A is after B": {
			a:      New(2025, time.January, 1, 0, 0, 0, 0),
			b:      New(2024, time.December, 31, 23, 59, 59, 999),
			expect: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if got := test.a.Before(test.b); got != test.expect {
				t.Errorf("Before: expected %v, got %v", test.expect, got)
			}
		})
	}
}

func TestLocalDateTime_After(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		a, b   LocalDateTime
		expect bool
	}{
		"A is after B - different date": {
			a:      New(2024, time.December, 31, 23, 59, 59, 0),
			b:      New(2024, time.January, 1, 0, 0, 0, 0),
			expect: true,
		},
		"A is after B - same date, later time": {
			a:      New(2024, time.July, 5, 18, 30, 0, 0),
			b:      New(2024, time.July, 5, 18, 0, 0, 0),
			expect: true,
		},
		"A is equal to B": {
			a:      New(2024, time.July, 5, 15, 45, 30, 0),
			b:      New(2024, time.July, 5, 15, 45, 30, 0),
			expect: false,
		},
		"A is before B": {
			a:      New(2023, time.January, 1, 0, 0, 0, 0),
			b:      New(2024, time.January, 1, 0, 0, 0, 0),
			expect: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if got := test.a.After(test.b); got != test.expect {
				t.Errorf("After: expected %v, got %v", test.expect, got)
			}
		})
	}
}

func TestLocalDateTime_Equal(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		a, b   LocalDateTime
		expect bool
	}{
		"Equal datetime": {
			a:      New(2024, time.July, 5, 12, 0, 0, 0),
			b:      New(2024, time.July, 5, 12, 0, 0, 0),
			expect: true,
		},
		"Different nanosecond": {
			a:      New(2024, time.July, 5, 12, 0, 0, 1),
			b:      New(2024, time.July, 5, 12, 0, 0, 0),
			expect: false,
		},
		"Different second": {
			a:      New(2024, time.July, 5, 12, 0, 1, 0),
			b:      New(2024, time.July, 5, 12, 0, 0, 0),
			expect: false,
		},
		"Different date": {
			a:      New(2023, time.July, 5, 12, 0, 0, 0),
			b:      New(2024, time.July, 5, 12, 0, 0, 0),
			expect: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if got := test.a.Equal(test.b); got != test.expect {
				t.Errorf("Equal: expected %v, got %v", test.expect, got)
			}
		})
	}
}
