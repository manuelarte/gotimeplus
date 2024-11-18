[![Go](https://github.com/manuelarte/GoTime/actions/workflows/go.yml/badge.svg)](https://github.com/manuelarte/GoTime/actions/workflows/go.yml)
![coverage](https://raw.githubusercontent.com/manuelarte/GoTime/badges/.badges/main/coverage.svg)
# 🕐 GoTime 🕐

## 📝 How to install it

> go get github.com/manuelarte/GoTime

## ✏️ Introduction

GoTime contains the following utility struct

### TimePeriod

Construct a time period based on start time and end time.

> tp, err := NewTimePeriod(startTime, endTime)

The time period is built based on the overlapping period between the two dates.

```
t1 ____|________
t2 _________|
tp ____|‾‾‾‾|___
```

If start time or end time are zero, it means no limit.  

The struct also provides a function `Overlaps` to check whether two time periods overlaps, and what is the overlapping period

e.g.

```
tp1 ____|‾‾‾‾‾‾‾‾‾‾‾‾‾‾
tp2 _________|‾‾‾‾‾‾|__
tp  ____|‾‾‾‾|_________
```

For more information check the [examples](./examples)
