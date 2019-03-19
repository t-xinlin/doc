package workpool

import (
	log "github.com/sirupsen/logrus"
)

const maxWorkers = 10

type Dispatcher struct {
	// A pool of workers channels that are registered with the dispatcher
	WorkerPool chan chan Job
	MaxWorkers int
}

func NewDispatcher(maxWorkers int) *Dispatcher {
	pool := make(chan chan Job, maxWorkers)
	return &Dispatcher{WorkerPool: pool, MaxWorkers: MaxWorkers}
}

func (d *Dispatcher) Run() {
	// starting n number of workers
	//for i := 0; i < d.maxWorkers; i++ {
	for i := 0; i < d.MaxWorkers; i++ {
		//worker := NewWorker(d.pool)
		worker := NewWorker(d.WorkerPool)
		worker.Start()
	}

	go d.dispatch()
}

func (d *Dispatcher) dispatch() {
	for {
		select {
		case job := <-JobQueue:
			// a job request has been received
			go func(job Job) {
				// try to obtain a worker job channel that is available.
				// this will block until a worker is idle
				jobChannel := <-d.WorkerPool
				log.Printf("task run %+v\n", job)
				// dispatch the job to the worker job channel
				jobChannel <- job
			}(job)
		}
	}
}
