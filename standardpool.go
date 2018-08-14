package pooligo

import (
	"context"
)

// NewPool creates a standard worker pool.  This pool will run, processing jobs
// asynchronously until the application exits.  If you need to stop the pool
// based on a context.Context being closed, use `NewCtxPool()`
func NewPool(workerCount, queueSize int) Pooli {
	return NewCtxPool(context.Background(), workerCount, queueSize)
}
