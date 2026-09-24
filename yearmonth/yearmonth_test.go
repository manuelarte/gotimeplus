package yearmonth

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestNew(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		year  int
		month time.Month
	}{
		"January 2024": {
			year:  2024,
			month: time.January,
		},
		"December 1999": {
			year:  1999,
			month: time.December,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := New(test.year, test.month)
			if got.Year() != test.year {
				t.Errorf("Year = %v, want %v", got.Year(), test.year)
			}

			if got.Month() != test.month {
				t.Errorf("Month = %v, want %v", got.Month(), test.month)
			}
		})
	}
}

func TestAddMonths(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		original yearMonth
		add      int
		want     yearMonth
	}{
		"within year": {
			original: yearMonth{
				year:  2024,
				month: time.March,
			},
			add: 4,
			want: yearMonth{
				year:  2024,
				month: time.July,
			},
		},
		"crosses into next year": {
			original: yearMonth{
				year:  2024,
				month: time.November,
			},
			add: 3,
			want: yearMonth{
				year:  2025,
				month: time.February,
			},
		},
		"crosses into previous year": {
			original: yearMonth{
				year:  2024,
				month: time.February,
			},
			add: -3,
			want: yearMonth{
				year:  2023,
				month: time.November,
			},
		},
		"zero months": {
			original: yearMonth{
				year:  2024,
				month: time.June,
			},
			add: 0,
			want: yearMonth{
				year:  2024,
				month: time.June,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := test.original.AddMonths(test.add)
			if cmp.Diff(got, test.want, cmp.AllowUnexported(yearMonth{})) != "" {
				t.Errorf("AddMonths() mismatch (-want +got):\n%s", cmp.Diff(test.want, got))
			}
		})
	}
}

func TestAddYears(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		original yearMonth
		add      int
		want     yearMonth
	}{
		"adds years": {
			original: yearMonth{
				year:  2024,
				month: time.March,
			},
			add: 3,
			want: yearMonth{
				year:  2027,
				month: time.March,
			},
		},
		"subtracts years": {
			original: yearMonth{
				year:  2024,
				month: time.March,
			},
			add: -2,
			want: yearMonth{
				year:  2022,
				month: time.March,
			},
		},
		"zero years": {
			original: yearMonth{
				year:  2024,
				month: time.March,
			},
			add: 0,
			want: yearMonth{
				year:  2024,
				month: time.March,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := test.original.AddYears(test.add)
			if cmp.Diff(got, test.want, cmp.AllowUnexported(yearMonth{})) != "" {
				t.Errorf("AddYears() mismatch (-want +got):\n%s", cmp.Diff(test.want, got))
			}
		})
	}
}

func TestAfter(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		a, b YearMonth
		want bool
	}{
		"later year": {
			a: yearMonth{
				year:  2025,
				month: time.January,
			},
			b: yearMonth{
				year:  2024,
				month: time.December,
			},
			want: true,
		},
		"later month in same year": {
			a: yearMonth{
				year:  2024,
				month: time.December,
			},
			b: yearMonth{
				year:  2024,
				month: time.January,
			},
			want: true,
		},
		"equal": {
			a: yearMonth{
				year:  2024,
				month: time.June,
			},
			b: yearMonth{
				year:  2024,
				month: time.June,
			},
			want: false,
		},
		"before": {
			a: yearMonth{
				year:  2024,
				month: time.January,
			},
			b: yearMonth{
				year:  2024,
				month: time.December,
			},
			want: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if got := test.a.After(test.b); got != test.want {
				t.Errorf("After = %v, want %v", got, test.want)
			}
		})
	}
}

func TestBefore(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		a, b YearMonth
		want bool
	}{
		"earlier year": {
			a: yearMonth{
				year:  2023,
				month: time.December,
			},
			b: yearMonth{
				year:  2024,
				month: time.January,
			},
			want: true,
		},
		"earlier month in same year": {
			a: yearMonth{
				year:  2024,
				month: time.January,
			},
			b: yearMonth{
				year:  2024,
				month: time.December,
			},
			want: true,
		},
		"equal": {
			a: yearMonth{
				year:  2024,
				month: time.June,
			},
			b: yearMonth{
				year:  2024,
				month: time.June,
			},
			want: false,
		},
		"after": {
			a: yearMonth{
				year:  2024,
				month: time.December,
			},
			b: yearMonth{
				year:  2024,
				month: time.January,
			},
			want: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if got := test.a.Before(test.b); got != test.want {
				t.Errorf("Before = %v, want %v", got, test.want)
			}
		})
	}
}

func TestCompare(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		a, b YearMonth
		want int
	}{
		"earlier year": {
			a: yearMonth{
				year:  2023,
				month: time.December,
			},
			b: yearMonth{
				year:  2024,
				month: time.January,
			},
			want: -1,
		},
		"earlier month": {
			a: yearMonth{
				year:  2024,
				month: time.January,
			},
			b: yearMonth{
				year:  2024,
				month: time.December,
			},
			want: -1,
		},
		"equal": {
			a: yearMonth{
				year:  2024,
				month: time.June,
			},
			b: yearMonth{
				year:  2024,
				month: time.June,
			},
			want: 0,
		},
		"later month": {
			a: yearMonth{
				year:  2024,
				month: time.December,
			},
			b: yearMonth{
				year:  2024,
				month: time.January,
			},
			want: 1,
		},
		"later year": {
			a: yearMonth{
				year:  2025,
				month: time.January,
			},
			b:    New(2024, time.December),
			want: 1,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if got := test.a.Compare(test.b); got != test.want {
				t.Errorf("Compare = %v, want %v", got, test.want)
			}
		})
	}
}

func TestEqual(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		a, b YearMonth
		want bool
	}{
		"same year and month": {
			a: yearMonth{
				year:  2024,
				month: time.June,
			},
			b: yearMonth{
				year:  2024,
				month: time.June,
			},
			want: true,
		},
		"different year": {
			a: yearMonth{
				year:  2023,
				month: time.June,
			},
			b: yearMonth{
				year:  2024,
				month: time.June,
			},
			want: false,
		},
		"different month": {
			a: yearMonth{
				year:  2024,
				month: time.May,
			},
			b:    New(2024, time.June),
			want: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if got := test.a.Equal(test.b); got != test.want {
				t.Errorf("Equal = %v, want %v", got, test.want)
			}
		})
	}
}
