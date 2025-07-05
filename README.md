# 🕐 GoTime Plus(+)

[![Go](https://github.com/manuelarte/gotimeplus/actions/workflows/go.yml/badge.svg)](https://github.com/manuelarte/gotimeplus/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/manuelarte/gotimeplus)](https://goreportcard.com/report/github.com/manuelarte/gotimeplus)
![coverage](https://raw.githubusercontent.com/manuelarte/gotimeplus/badges/.badges/main/coverage.svg)
![version](https://img.shields.io/github/v/release/manuelarte/gotimeplus)

GoTime Plus is a Go library that adds some missing functionality to the standard `time.Time` Go package.

## ⬇️ How to use it

```bash
go get github.com/manuelarte/gotime@latest
``` 

## 🚀 Features

GoTime Plus contains the following features:

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
