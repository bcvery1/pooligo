package pooligo

import (
	"sync"
	"sync/atomic"
)

const (
	workerPoolSmall = 10
	workerPoolLarge = 1024
	queueSizeSmall  = 5
	queueSizeLarge  = 80
	// How many jobs to run through the test pools
	jobCount = 10000
)

var (
	// counter is used to establish jobs have been actioned
	counter uint64
	wg      sync.WaitGroup
)

type testJob struct{}

// Action is so testJob fits the job interface
func (t testJob) Action() {
	defer wg.Done()
	atomic.AddUint64(&counter, 1)
}
