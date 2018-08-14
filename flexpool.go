package pooligo

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"
)

// Flexpool is an extension to the standard worker pool.  It can be grown and
// shrunk after creation
type Flexpool struct {
	cancelFunc context.CancelFunc
	queue      chan job
	closed     atomic.Value

	// ctx carries the context for this pool.  Not ideal practise to store the
	// context in a struct, but needed to abstract growing the worker pool.
	ctx context.Context

	// workers holds a slice of individual context cancel functions, one assigned
	// to each worker in this pool.  All contexts are children of the pool context
	workers []context.CancelFunc
}

func (p *Flexpool) setClosed() {
	p.closed.Store(true)
}

// Add is used to add a job to the worker-pool.
func (p *Flexpool) Add(j job) {
	if p.closed.Load() != true {
		p.queue <- j
	}
}

// Close will stop any workers doing any jobs added after calling Close.
// Any jobs being processed will be completed.
func (p *Flexpool) Close() {
	p.cancelFunc()
}

// Size returns the current size of the worker pool.
func (p *Flexpool) Size() int {
	return len(p.workers)
}

// Grow will increase the worker pool by the number provided.
func (p *Flexpool) Grow(i uint) {
	for ; i > 0; i-- {
		p.addWorker()
	}
}

// Shrink will decrease the worker pool by the number provided.  Will return an
// error if the number provided is greater than the worker pool size.
// Shrink will cancel the pool context if the number provided is the same size
// as the worker pool.
func (p *Flexpool) Shrink(i int) error {
	if i > len(p.workers) {
		return fmt.Errorf("Cannot remove %d workers, worker pool has only  %d workers", i, len(p.workers))
	}
	if i < 1 {
		return errors.New("Cannot shrink by less than 1 worker")
	}
	if i == len(p.workers) {
		p.Close()
		return nil
	}

	for ; i > 0; i-- {
		// Call the cancel function
		p.workers[0]()
		// Remove the cancel function representing that worker
		p.workers = p.workers[1:]
	}
	return nil
}

// addWorker will add a worker to the pool.  This is used so code is not duplicated
// in this file.
func (p *Flexpool) addWorker() {
	workerCtx, workerCancel := context.WithCancel(p.ctx)
	runWorker(workerCtx, p.queue, p)
	p.workers = append(p.workers, workerCancel)
}

// NewFlexPool creates a standard worker pool
func NewFlexPool(workerCount, queueSize int) *Flexpool {
	// Create a context
	ctx, cancel := context.WithCancel(context.Background())

	// Create the queue
	q := make(chan job, queueSize)

	// Create the pool
	p := &Flexpool{
		cancelFunc: cancel,
		queue:      q,
		workers:    []context.CancelFunc{},
		ctx:        ctx,
	}
	p.closed.Store(false)

	// Create workers
	for i := 0; i < workerCount; i++ {
		p.addWorker()
	}

	return p
}
