package pooligo

import (
	"context"
)

// Pooli an interface to pool types.  All pools returned by this module will be
// implement Pooli.
type Pooli interface {
	// Add will add a job to the worker pool.  This will be process asychronously.
	Add(job)

	// Close will stop any new jobs being added and/or processed.  Close will not
	// force stop any currently running jobs, unless specific pool types allow
	// and specify.
	Close()

	// setClosed is called by `runWorker` when the pool context is cancelled.  This
	// can be done by the calling code with 'Pooli.Close()'
	setClosed()

	// Size returns the current size of the worker pool.
	Size() int
}

type job interface {
	Action()
}

// runWorker will create a worker on the queue provided
// The worker will stop accepting jobs after the context is cancelled
func runWorker(ctx context.Context, jobQueue <-chan job, p Pooli) {
	go func(q <-chan job) {
		for {
			select {
			case j := <-jobQueue:
				j.Action()
			case <-ctx.Done():
				p.setClosed()
				return
			}
		}
	}(jobQueue)
}
