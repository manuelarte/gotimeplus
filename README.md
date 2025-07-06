# 🕐 GoTime Plus(+)

[![Go](https://github.com/manuelarte/gotimeplus/actions/workflows/go.yml/badge.svg)](https://github.com/manuelarte/gotimeplus/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/manuelarte/gotimeplus)](https://goreportcard.com/report/github.com/manuelarte/gotimeplus)
![coverage](https://raw.githubusercontent.com/manuelarte/gotimeplus/badges/.badges/main/coverage.svg)
![version](https://img.shields.io/github/v/release/manuelarte/gotimeplus)

GoTime Plus is a Go library that adds some missing functionality to the standard `time.Time` Go package.

- [🕐 GoTime Plus(+)](#-gotime-plus)
    * [⬇️ How to use it](#-how-to-use-it)
    * [🚀 Features](#-features)
        + [LocalDate](#localdate)
        + [LocalTime](#localtime)
        + [LocalDateTime](#localdatetime)
        + [TimePeriod](#timeperiod)
    * [📂 Examples](#-examples)

## ⬇️ How to use it

```bash
go get github.com/manuelarte/gotime@latest
``` 

## 🚀 Features

GoTime Plus contains the following features:

### LocalDate

Same concept as java [LocalDate][javaLocalDate], this struct represents a date, often viewed as year-month-day.
This struct does not represent a time or time-zone. Instead, it is a description of the date, as used for birthdays. 
It cannot represent a `time.Time` without additional information such a time-zone.

e.g.:

```go
goBirthdate := localdate.New(2009, time.November, 10)
```

### LocalTime

Same concept as java [LocalTime][javaLocalTime]. This struct represents a time without a time-zone, such as 10:15:30.
This struct does not represent a time or time-zone. Instead, it is a description of the local time, as seen on a wall clock.
It cannot represent a `time.Time` without additional information such as date and a time-zone.

e.g.:

```go
lunchTime := localtime.New(12, 0, 0, 0)
```

### LocalDateTime

Same concept as java [LocalDateTime][javaLocalDateTime]. This struct represents a date-time, such as 2007-12-03T10:15:30.
This struct does not represent a time-zone. Instead, it is a description of the date plus a time.
It cannot represent a `time.Time` without additional information such as a time-zone.

e.g.:

```go
newYear2025 := localdatetime.New(localdate.New(2025, 1, 1), localtime.New(0, 0, 0, 0))
```

### TimePeriod

Create a `TimePeriod` instance by specifying a start time and an end time:

> tp, err := NewTimePeriod(startTime, endTime)

+ `startTime`: The beginning of the time period. Use `time.Time{}` for no lower limit.
+ `endTime`: The end of the time period. Use `time.Time{}` for no upper limit.

Returns:

+ `tp`: The resulting `TimePeriod`.
+ `err`: An error if the inputs are invalid.

The `TimePeriod` is built based on the overlapping period between the two dates.

```bash
Input Times
time1 ____|________...
time2 _________|___...
Resulting Time Period
tp    ____|‾‾‾‾|___...
```

> [!WARNING]  
> Passing a zero value for `startTime` or `endTime` indicates an unbounded period on that side.

---

The struct also provides a function `Overlaps`. This method checks whether two time periods overlap.
If so, returns the overlapping period, e.g.:

```bash
Input Time Periods
tp1 ____|‾‾‾‾‾‾‾‾‾‾‾‾‾‾...
tp2 _________|‾‾‾‾‾‾|__...
Resulting Overlap
tp  ____|‾‾‾‾|_________...
```

## 📂 Examples

Refer to the [examples](./examples) directory for usage examples.

[javaLocalDate]: https://docs.oracle.com/javase/8/docs/api/java/time/LocalDate.html
[javaLocalTime]: https://docs.oracle.com/javase/8/docs/api/java/time/LocalTime.html
[javaLocalDateTime]: https://docs.oracle.com/javase/8/docs/api/java/time/LocalDateTime.html
