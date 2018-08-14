package pooligo

import (
	"context"
	"sync/atomic"
)

// CtxPool struct representing a standard worker-pool but with a parent context.
// If you do not need to explicitly pass in a context, use `Pool` (from `NewPool()`)
type CtxPool struct {
	cancelFunc  context.CancelFunc
	queue       chan<- job
	closed      atomic.Value
	workerCount int
}

func (p *CtxPool) setClosed() {
	p.closed.Store(true)
}

// Add is used to add a job to the worker-pool
func (p *CtxPool) Add(j job) {
	if p.closed.Load() != true {
		p.queue <- j
	}
}

// Close will stop any workers doing any jobs added after calling Close
// Any jobs being processed will be completed
func (p *CtxPool) Close() {
	p.cancelFunc()
}

// Size returns the current size of the worker pool.
func (p *CtxPool) Size() int {
	return p.workerCount
}

// NewCtxPool creates a standard worker pool with a parent context
func NewCtxPool(parentCtx context.Context, workerCount, queueSize int) *CtxPool {
	// Create a context
	ctx, cancel := context.WithCancel(parentCtx)

	// Create the queue
	q := make(chan job, queueSize)

	// Create the pool
	p := &CtxPool{
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
