[![Go](https://github.com/manuelarte/GoTime/actions/workflows/go.yml/badge.svg)](https://github.com/manuelarte/GoTime/actions/workflows/go.yml)
![coverage](https://raw.githubusercontent.com/manuelarte/GoTime/badges/.badges/main/coverage.svg)
# 🕐 GoTime

GoTime is a Go library for working with time periods, enabling easy creation and overlap calculation. GoTime simplifies the creation and manipulation of time periods, allowing you to easily define, compare, and compute overlaps between time intervals.

## 📝 How to install it

> go get github.com/manuelarte/GoTime

## ✏️ Introduction

GoTime contains the following utility struct

### TimePeriod

Create a `TimePeriod` instance by specifying a start time and an end time:

> tp, err := NewTimePeriod(startTime, endTime)

+ `startTime`: The beginning of the time period. Use `time.Time{}` for no lower limit.
+ `endTime`: The end of the time period. Use `time.Time{}` for no upper limit.
Returns:
+ `tp`: The resulting TimePeriod.
+ `err`: An error if the inputs are invalid.

The time period is built based on the overlapping period between the two dates.

```
Input Times
time1 ____|________...
time2 _________|___...
Resulting Time Period
tp    ____|‾‾‾‾|___...
```

Passing a zero value for `startTime` or `endTime` indicates an unbounded period on that side.

---

The struct also provides a function `Overlaps`. This method checks whether two time periods overlap and, if so, returns the overlapping period.

e.g.

```
Input Time Periods
tp1 ____|‾‾‾‾‾‾‾‾‾‾‾‾‾‾...
tp2 _________|‾‾‾‾‾‾|__...
Resulting Overlap
tp  ____|‾‾‾‾|_________...
```

## Full Example

```go
start1 := time.Date(2024, 12, 1, 0, 0, 0, 0, time.UTC)
end1 := time.Date(2024, 12, 10, 0, 0, 0, 0, time.UTC)
tp1, _ := NewTimePeriod(start1, end1)

start2 := time.Date(2024, 12, 5, 0, 0, 0, 0, time.UTC)
end2 := time.Date(2024, 12, 15, 0, 0, 0, 0, time.UTC)
tp2, _ := NewTimePeriod(start2, end2)

overlapExists, overlapPeriod := tp1.Overlaps(tp2)
fmt.Println("Overlap Exists:", overlapExists)
fmt.Println("Overlap Period:", overlapPeriod)
```

## 📂 Examples

Refer to the [examples](./examples) directory for usage examples.

## 📜 License
This project is licensed under the MIT License.
