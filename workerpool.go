package main

type Job interface {
	Do()
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
				job.Do()
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
