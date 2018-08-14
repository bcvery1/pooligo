package pooligo

import (
	"context"
	"sync/atomic"
)

// Pool struct representing a standard worker-pool
type Pool struct {
	cancelFunc  context.CancelFunc
	queue       chan<- job
	closed      atomic.Value
	workerCount int
}

func (p *Pool) setClosed() {
	p.closed.Store(true)
}

// Add is used to add a job to the worker-pool
func (p *Pool) Add(j job) {
	if p.closed.Load() != true {
		p.queue <- j
	}
}

// Close will stop any workers doing any jobs added after calling Close
// Any jobs being processed will be completed
func (p *Pool) Close() {
	p.cancelFunc()
}

// Size returns the current size of the worker pool.
func (p *Pool) Size() int {
	return p.workerCount
}

// NewPool creates a standard worker pool
func NewPool(workerCount, queueSize int) *Pool {
	// Create a context
	ctx, cancel := context.WithCancel(context.Background())

	// Create the queue
	q := make(chan job, queueSize)

	// Create the pool
	p := &Pool{
		cancelFunc:  cancel,
		queue:       q,
		workerCount: workerCount,
	}
	p.closed.Store(false)

	// Create workers
	for i := 0; i < queueSize; i++ {
		runWorker(ctx, q, p)
	}

	return p
}
