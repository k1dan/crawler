package workerpool

import (
	"context"
	"sync"
)

type WorkerPool struct {
	workerCount int
	jobs        chan Job
	results     chan Result
	Done        chan struct{}
}

// New func initializes new WorkerPool instance
func New(wCount int) WorkerPool {
	return WorkerPool{
		workerCount: wCount,
		jobs:        make(chan Job, wCount),
		results:     make(chan Result, wCount),
		Done:        make(chan struct{}),
	}
}

// Run func starts worker pool
func (wp WorkerPool) Run(ctx context.Context) {
	var wg sync.WaitGroup

	for i := 0; i < wp.workerCount; i++ {
		wg.Add(1)

		go worker(ctx, &wg, i, wp.jobs, wp.results)
	}

	wg.Wait()
	close(wp.Done)
	close(wp.results)
}

// Results func is getter for results filed
func (wp WorkerPool) Results() <-chan Result {
	return wp.results
}

// AddJob func adds new job to job channel
func (wp WorkerPool) AddJob(job Job) {
	wp.jobs <- job
}

// CloseJobsChan closes jobs channel
func (wp WorkerPool) CloseJobsChan() {
	close(wp.jobs)
}

// worker func monitor job channel, executes job and send result to result chan
// either it stops monitoring if signal to ctx.Done chan is sent
func worker(ctx context.Context, wg *sync.WaitGroup, i int, jobs <-chan Job, results chan<- Result) {
	defer wg.Done()
	for {
		select {
		case job, ok := <-jobs:
			if !ok {
				return
			}
			results <- job.execute(ctx)
		case <-ctx.Done():
			results <- Result{
				Err: ctx.Err(),
			}
			return
		}
	}
}
