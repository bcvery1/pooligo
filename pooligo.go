package pooligo

import (
	"context"
)

// Pooli an interface to pool types.  All pools returned by this module will be
// implement Pooli.
type Pooli interface {
	Add(job)
	Close()
	SetClosed()
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
				p.SetClosed()
				return
			}
		}
	}(jobQueue)
}
