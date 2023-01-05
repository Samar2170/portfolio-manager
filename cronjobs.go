package main

import (
	"log"
	"time"

	"github.com/go-co-op/gocron"
)

const (
	UpdateNextIPDateTiming = "12:01:00"
)

func StartCronServer() {
	s := gocron.NewScheduler(time.UTC)
	s.Every(2).Second().Do(func() {
		log.Println("Working buddy")
	})
}

// run job for each fd in bulk, store status at the task level. Bitwise storage
