package main

import (
	"reflect"

	"github.com/Samar2170/portfolio-manager/account"
	"github.com/Samar2170/portfolio-manager/utils"
)

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
	structMap := utils.Inspect(&j)
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
					NewJobRecord(job, false, err.Error())
				} else {
					NewJobRecord(job, true, "")
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
	pool := make(chan chan Job, maxWorkers)
	return Supervisor{
		MaxWorkers: maxWorkers,
		WorkerPool: pool,
		quit:       make(chan bool),
	}
}

var JobQueue chan Job

func (s Supervisor) Run() {
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
