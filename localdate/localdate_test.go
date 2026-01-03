package localdate

import (
	"testing"
	"time"
)

func TestAfter(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		a, b   LocalDate
		expect bool
	}{
		"A is after B": {
			a:      New(2025, time.December, 31),
			b:      New(2024, time.December, 31),
			expect: true,
		},
		"A is same as B": {
			a:      New(2024, time.July, 5),
			b:      New(2024, time.July, 5),
			expect: false,
		},
		"A is before B": {
			a:      New(2023, time.January, 1),
			b:      New(2024, time.January, 1),
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

func TestBefore(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		a, b   LocalDate
		expect bool
	}{
		"A is before B": {
			a:      New(2023, time.January, 1),
			b:      New(2024, time.January, 1),
			expect: true,
		},
		"A is same as B": {
			a:      New(2024, time.July, 5),
			b:      New(2024, time.July, 5),
			expect: false,
		},
		"A is after B": {
			a:      New(2025, time.December, 31),
			b:      New(2024, time.December, 31),
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

func TestEqual(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		a, b   LocalDate
		expect bool
	}{
		"Same date": {
			a:      New(2024, time.July, 5),
			b:      New(2024, time.July, 5),
			expect: true,
		},
		"Different year": {
			a:      New(2023, time.July, 5),
			b:      New(2024, time.July, 5),
			expect: false,
		},
		"Different month": {
			a:      New(2024, time.June, 5),
			b:      New(2024, time.July, 5),
			expect: false,
		},
		"Different day": {
			a:      New(2024, time.July, 4),
			b:      New(2024, time.July, 5),
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
