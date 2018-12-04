package pooligo

import (
	"context"
)

// Pooli an interface to pool types.  All pools returned by this module will be
// implement Pooli.
type Pooli interface {
	// Add will add a job to the worker pool.  This will be process asychronously.
	Add(Job)

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

// Job is the interface which is accepted by the `Pooli` interface, upon which
// work will be done
type Job interface {
	Action()
}

// runWorker will create a worker on the queue provided
// The worker will stop accepting jobs after the context is cancelled
func runWorker(ctx context.Context, jobQueue <-chan Job, p Pooli) {
	go func(q <-chan Job) {
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
