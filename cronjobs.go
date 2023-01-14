package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Samar2170/portfolio-manager/securities"
	"github.com/go-co-op/gocron"
)

var cronLogger *log.Logger

const (
	UpdateNextIPDateTiming = "12:01:00"
)

func StartCronServer() {
	t := time.Now()
	s := gocron.NewScheduler(time.UTC)

	file, err := os.OpenFile(fmt.Sprintf("logs/Cron_logs_%d-%d-%d", t.Day(), int(t.Month()), t.Year()), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}
	cronLogger = log.New(file, "[Cron]-", log.Ldate|log.Ltime|log.Lshortfile)
	cronLogger.Println("Cron Server Working buddy")

	s.Every(1).Day().At("00:01:00").Do(func() {
		cronLogger.Println("Running Update Next IP Date Job")
		securities.UpdateNextIPDatesFDs()
	})
}

// run job for each fd in bulk, store status at the task level. Bitwise storage
