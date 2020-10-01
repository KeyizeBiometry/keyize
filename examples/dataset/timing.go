package main

import (
	"fmt"
	"time"
)

var startedAt time.Time
var timerName string

func EndTime() {
	finished := time.Now()

	fmt.Printf("> %s completed in %dms\n", timerName, finished.Sub(startedAt).Milliseconds())
}

func StartTime(name string) {
	timerName = name
	startedAt = time.Now()
}
