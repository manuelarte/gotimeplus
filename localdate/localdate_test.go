package localdate

import (
	"testing"
	"time"
)

func TestAfter(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		a, b LocalDate
		want bool
	}{
		"A is after B": {
			a:    New(2025, time.December, 31),
			b:    New(2024, time.December, 31),
			want: true,
		},
		"A is same as B": {
			a:    New(2024, time.July, 5),
			b:    New(2024, time.July, 5),
			want: false,
		},
		"A is before B": {
			a:    New(2023, time.January, 1),
			b:    New(2024, time.January, 1),
			want: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if got := test.a.After(test.b); got != test.want {
				t.Errorf("test.a.After(%v) = %v, got %v", test.b, got, test.want)
			}
		})
	}
}

func TestBefore(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		a, b LocalDate
		want bool
	}{
		"A is before B": {
			a:    New(2023, time.January, 1),
			b:    New(2024, time.January, 1),
			want: true,
		},
		"A is same as B": {
			a:    New(2024, time.July, 5),
			b:    New(2024, time.July, 5),
			want: false,
		},
		"A is after B": {
			a:    New(2025, time.December, 31),
			b:    New(2024, time.December, 31),
			want: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if got := test.a.Before(test.b); got != test.want {
				t.Errorf("test.a.Before(%v) = %v, got %v", test.b, got, test.want)
			}
		})
	}
}

func TestEqual(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		a, b LocalDate
		want bool
	}{
		"Same date": {
			a:    New(2024, time.July, 5),
			b:    New(2024, time.July, 5),
			want: true,
		},
		"Different year": {
			a:    New(2023, time.July, 5),
			b:    New(2024, time.July, 5),
			want: false,
		},
		"Different month": {
			a:    New(2024, time.June, 5),
			b:    New(2024, time.July, 5),
			want: false,
		},
		"Different day": {
			a:    New(2024, time.July, 4),
			b:    New(2024, time.July, 5),
			want: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if got := test.a.Equal(test.b); got != test.want {
				t.Errorf("test.a.Equal(%v) = %v, want %v", test.b, got, test.want)
			}
		})
	}
}
