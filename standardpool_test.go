package pooligo

import (
	"testing"
	"time"
)

func TestNewPool(t *testing.T) {
	counter = 0

	// Create pool and test job
	p := NewPool(workerPoolSmall, queueSizeSmall)

	// Check the size of the pool
	if p.Size() != workerPoolSmall {
		t.Errorf("Expect pool size of %d, got: %d", workerPoolSmall, p.Size())
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
		t.Errorf("Incorrect number of jobs ran.  Expected %d, got %d", jobCount, counter)
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
		t.Errorf("Workers running after pool closed.  Ran %d jobs", counter-jobCount)
	}

	// Clear the waitgroup addition not done after closing
	wg.Done()
}

// Check how fast a small work/job queue pool runs
func BenchmarkNewPool(b *testing.B) {
	p := NewPool(workerPoolSmall, queueSizeSmall)

	j := testJob{}
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		p.Add(j)
	}
	wg.Wait()
}

// Check how fast a large work/job queue pool runs
func BenchmarkNewPoolLarge(b *testing.B) {
	p := NewPool(workerPoolLarge, queueSizeLarge)

	j := testJob{}
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		p.Add(j)
	}
	wg.Wait()
}
