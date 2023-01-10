package main

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"time"

	"github.com/Samar2170/portfolio-manager/account"
	"github.com/Samar2170/portfolio-manager/utils"
)

var logger *log.Logger

func startLogger() {
	t := time.Now()
	file, err := os.OpenFile(fmt.Sprintf("logs/Workerpool_logs_%d-%d-%d", t.Day(), int(t.Month()), t.Year()), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}
	logger = log.New(file, "[Workerpool]-", log.Ldate|log.Ltime|log.Lshortfile)
	logger.Println("----------------Workerpool Logger Started---------------")
}

type Job interface {
	Do() error
}

func NewJobRecord(j Job, success bool, errorString string) account.JobRecord {

	refJob := reflect.TypeOf(j)
	fields := reflect.VisibleFields(refJob)
	fieldString := ""
	for _, field := range fields {
		fieldString += field.Name + ","
	}
	structMap := utils.Inspect(j)
	args := ""
	for _, v := range structMap {
		args += v
	}
	jr := account.JobRecord{
		Name:    refJob.Name(),
		Fields:  fieldString,
		Args:    args,
		Success: success,
		Error:   errorString,
	}
	jr.Create()
	return jr
}

type Worker struct {
	WorkerPool chan chan Job
	JobChannel chan Job
	quit       chan bool
}

type Supervisor struct {
	MaxWorkers uint
	WorkerPool chan chan Job
	quit       chan bool
}

func NewWorker(workerPool chan chan Job) Worker {
	return Worker{
		WorkerPool: workerPool,
		JobChannel: make(chan Job),
		quit:       make(chan bool),
	}
}

func (w Worker) Start() {
	go func() {
		for {
			w.WorkerPool <- w.JobChannel
			// take the Job channel and put in in workerpool
			select {
			case job := <-w.JobChannel:

				err := job.Do()
				if err != nil {

					jr := NewJobRecord(job, false, err.Error())
					logger.Printf("Executing Job Name %s with fields %s", jr.Name, jr.Fields)
				} else {
					jr := NewJobRecord(job, true, "")
					logger.Printf("Executing Job Name %s with fields %s", jr.Name, jr.Fields)
				}

			case <-w.quit:
				return
			}

		}

	}()
}

func (w Worker) Stop() {
	go func() {
		w.quit <- true
	}()
}
func NewSuperVisor(maxWorkers uint) Supervisor {
	startLogger()
	pool := make(chan chan Job, maxWorkers)
	return Supervisor{
		MaxWorkers: maxWorkers,
		WorkerPool: pool,
		quit:       make(chan bool),
	}
}

var JobQueue chan Job

func (s Supervisor) Run() {
	for i := 0; i < int(s.MaxWorkers); i++ {
		woker := NewWorker(s.WorkerPool)
		woker.Start()
	}
	go s.dispatch()
}

func (s Supervisor) dispatch() {
	for {
		select {
		case job := <-JobQueue:
			go func(job Job) {
				JobChannel := <-s.WorkerPool
				JobChannel <- job
			}(job)
		case <-s.quit:
			return
		}
	}
}

func (s Supervisor) Stop() {
	go func() {
		s.quit <- true
	}()
}

func RunSuperVisor() {
	JobQueue = make(chan Job)
	supervisor := NewSuperVisor(4)
	supervisor.Run()
}
