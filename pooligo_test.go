package pooligo

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
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

func testPoolCreation(p Pooli) error {
	counter = 0

	// Check the size of the pool
	if p.Size() != workerPoolSmall {
		return fmt.Errorf("Expect pool size of %d, got: %d", workerPoolSmall, p.Size())
	}

	// Loop adding all jobs
	wg.Add(jobCount)
	tj := testJob{}
	for i := 0; i < jobCount; i++ {
		p.Add(tj)
	}

	// Allow time for workers to finish
	wg.Wait()

	// Check the counter has been increased sufficiently
	if counter != jobCount {
		return fmt.Errorf("Incorrect number of jobs ran.  Expected %d, got %d", jobCount, counter)
	}

	// Close the context
	p.Close()
	// Allow time for the channel to close
	time.Sleep(500 * time.Millisecond)

	// Add another job to the pool
	// Add to the wait group so we don't get a negative counter
	wg.Add(1)
	p.Add(tj)

	// Allow time for workers to finish
	// Note: they shouldn't be working anyway - this pause if for if it is broken
	time.Sleep(500 * time.Millisecond)

	if counter != jobCount {
		return fmt.Errorf("Workers running after pool closed.  Ran %d jobs", counter-jobCount)
	}

	// Clear the waitgroup addition not done after closing
	wg.Done()

	return nil
}
