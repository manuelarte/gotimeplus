package main

import (
	"fmt"
	"time"

	"github.com/manuelarte/gotimeplus/timeperiod"
)

const timeFormat = "2006-January-02"

func main() {
	startTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	endTime := startTime.Add(31 * 24 * time.Hour)

	timePeriod, err := timeperiod.New(&startTime, &endTime)
	if err != nil {
		panic(err)
	}
	fmt.Printf("TimePeriod1:\t%s\t - %s\n", timePeriod.StartTime().Format(timeFormat), timePeriod.EndTime().Format(timeFormat))

	startTime2 := startTime.Add(14 * 24 * time.Hour)
	endTime2 := endTime.Add(14 * 24 * time.Hour)
	timePeriod2, err := timeperiod.New(&startTime2, &endTime2)
	if err != nil {
		panic(err)
	}
	fmt.Printf("TimePeriod2:\t%s\t - %s\n", timePeriod2.StartTime().Format(timeFormat), timePeriod2.EndTime().Format(timeFormat))

	overlapPeriod, ok := timePeriod.Overlaps(timePeriod2)
	if ok {
		fmt.Printf("OverlapPeriod:\t%s\t - %s\n", overlapPeriod.StartTime().Format(timeFormat), overlapPeriod.EndTime().Format(timeFormat))
	}

	fmt.Printf("%+v\n", timeperiod.Infinite)
}
