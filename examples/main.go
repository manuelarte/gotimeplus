package main

import (
	"GoTime/pkg/timeperiod"
	"fmt"
	"time"
)

func main() {
	startTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	endTime := startTime.Add(2 * time.Hour)

	timePeriod, err := timeperiod.NewTimePeriod(startTime, endTime)
	if err != nil {
		panic(err)
	}
	fmt.Printf("TimePeriod created:\tStartTime: %s\tEndTime: %s\n", timePeriod.GetStartTime(), timePeriod.GetEndTime())

	startTime2 := startTime.Add(1 * time.Hour)
	endTime2 := endTime.Add(2 * time.Hour)
	timePeriod2, err := timeperiod.NewTimePeriod(startTime2, endTime2)
	if err != nil {
		panic(err)
	}
	fmt.Printf("TimePeriod2 created:\tStartTime: %s\tEndTime: %s\n", timePeriod2.GetStartTime(), timePeriod2.GetEndTime())

	overlapPeriod, ok := timePeriod.Overlaps(timePeriod2)
	if ok {
		fmt.Printf("OverlapPeriod created:\tStartTime: %s\tEndTime: %s\n", overlapPeriod.GetStartTime(), overlapPeriod.GetEndTime())
	}
}
