// Package yearmonth provides YearMonth concept.
// Same concept as https://docs.oracle.com/javase/8/docs/api/java/time/YearMonth.html.
package yearmonth

import "time"

type (
	// YearMonth interface to track a period of time.
	YearMonth interface {
		// AddMonths Returns a new YearMonth with the specified number of months added.
		AddMonths(int) YearMonth
		// AddYears Returns a new YearMonth with the specified number of years added.
		AddYears(int) YearMonth
		// After Returns true if this YearMonth is after the specified YearMonth.
		After(YearMonth) bool
		// Before Returns true if this YearMonth is before the specified YearMonth.
		Before(YearMonth) bool
		// Compare compares this YearMonth with the specified YearMonth.
		Compare(YearMonth) int
		// Equal returns true if this YearMonth equals the specified YearMonth.
		Equal(YearMonth) bool
		// Year Get the year and month.
		Year() int
		// Month Get the month.
		Month() time.Month
	}

	yearMonth struct {
		year  int
		month time.Month
	}
)

// New Creates a new YearMonth.
func New(year int, month time.Month) YearMonth {
	return yearMonth{
		year:  year,
		month: month,
	}
}

// AddMonths Returns a new YearMonth with the specified number of months added.
func (ym yearMonth) AddMonths(n int) YearMonth {
	totalMonths := ym.year*12 + int(ym.month) - 1 + n
	year := totalMonths / 12

	month := totalMonths % 12
	if month < 0 {
		year--
		month += 12
	}

	return New(year, time.Month(month+1))
}

// AddYears Returns a new YearMonth with the specified number of years added.
func (ym yearMonth) AddYears(n int) YearMonth { return ym.AddMonths(12 * n) }

// After Returns true if this YearMonth is after the specified YearMonth.
func (ym yearMonth) After(o YearMonth) bool { return ym.Compare(o) > 0 }

// Before Returns true if this YearMonth is before the specified YearMonth.
func (ym yearMonth) Before(o YearMonth) bool { return ym.Compare(o) < 0 }

// Compare compares this YearMonth with the specified YearMonth.
func (ym yearMonth) Compare(o YearMonth) int {
	switch {
	case ym.Year() != o.Year():
		if ym.Year() < o.Year() {
			return -1
		}

		return 1
	case ym.Month() != o.Month():
		if ym.Month() < o.Month() {
			return -1
		}

		return 1
	}

	return 0
}

func (ym yearMonth) Equal(o YearMonth) bool { return ym.Compare(o) == 0 }

// Month Returns the month of this period.
func (ym yearMonth) Month() time.Month {
	return ym.month
}

// Year Returns the year of this period.
func (ym yearMonth) Year() int {
	return ym.year
}
