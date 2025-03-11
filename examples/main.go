package main

import (
	"fmt"
	"time"

	"github.com/manuelarte/gotime/pkg/timeperiod"
)

const timeFormat = "2006-January-02"

func main() {
	startTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	endTime := startTime.Add(31 * 24 * time.Hour)

	timePeriod, err := timeperiod.NewTimePeriod(startTime, endTime)
	if err != nil {
		panic(err)
	}
	fmt.Printf("TimePeriod1:\t%s\t - %s\n", timePeriod.GetStartTime().Format(timeFormat), timePeriod.GetEndTime().Format(timeFormat))

	startTime2 := startTime.Add(14 * 24 * time.Hour)
	endTime2 := endTime.Add(14 * 24 * time.Hour)
	timePeriod2, err := timeperiod.NewTimePeriod(startTime2, endTime2)
	if err != nil {
		panic(err)
	}
	fmt.Printf("TimePeriod2:\t%s\t - %s\n", timePeriod2.GetStartTime().Format(timeFormat), timePeriod2.GetEndTime().Format(timeFormat))

	overlapPeriod, ok := timePeriod.Overlaps(timePeriod2)
	if ok {
		fmt.Printf("OverlapPeriod:\t%s\t - %s\n", overlapPeriod.GetStartTime().Format(timeFormat), overlapPeriod.GetEndTime().Format(timeFormat))
	}
}
