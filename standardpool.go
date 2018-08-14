package pooligo

import (
	"context"
	"sync/atomic"
)

// pool struct representing a standard worker-pool
type pool struct {
	cancelFunc context.CancelFunc
	queue      chan<- job
	closed     atomic.Value
}

// SetClosed marks the pool as closed
// This will prevent additional jobs being added.  It will not force quit any
// currently running jobs.
func (p *pool) SetClosed() {
	p.closed.Store(true)
}

// Add is used to add a job to the worker-pool
func (p *pool) Add(j job) {
	if p.closed.Load() != true {
		p.queue <- j
	}
}

// Close will stop any workers doing any jobs added after calling Close
// Any jobs being processed will be completed
func (p *pool) Close() {
	p.cancelFunc()
}

// NewPool creates a standard worker pool
func NewPool(workerCount, queueSize int) Pooli {
	// Create a context
	ctx, cancel := context.WithCancel(context.Background())

	// Create the queue
	q := make(chan job, queueSize)

	// Create the pool
	p := &pool{
		cancelFunc: cancel,
		queue:      q,
	}
	p.closed.Store(false)

	// Create workers
	for i := 0; i < queueSize; i++ {
		runWorker(ctx, q, p)
	}

	return p
}
