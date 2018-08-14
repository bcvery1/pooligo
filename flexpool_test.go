package pooligo

import (
	"testing"
	"time"
)

func TestNewFlexPool(t *testing.T) {
	counter = 0

	// Create pool and test job
	p := NewFlexPool(workerPoolSmall, queueSizeSmall)

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

func TestFlexPoolSizeChange(t *testing.T) {
	p := NewFlexPool(workerPoolSmall, queueSizeSmall)

	// Check the initial size
	if p.Size() != workerPoolSmall {
		t.Errorf("Expect pool size of %d, got: %d", workerPoolSmall, p.Size())
	}

	// Grow by 1
	p.Grow(1)
	// Check the new size
	if p.Size() != workerPoolSmall+1 {
		t.Errorf("Expect pool size of %d, got: %d - After growing by 1", workerPoolSmall+1, p.Size())
	}

	// Check a negative shrink request fails
	if p.Shrink(-1) == nil {
		t.Errorf("Failed to error when shrinking by -1")
	}
	// Check the size hasn't changed
	if p.Size() != workerPoolSmall+1 {
		t.Errorf("Expect pool size of %d, got: %d - After negative shrink", workerPoolSmall+1, p.Size())
	}

	// Check pool cannot be shrunk by more than the pool size
	if p.Shrink(p.Size()+1) == nil {
		t.Errorf("Failed to error when shrinking by greater than pool size")
	}
	// Check the size hasn't changed
	if p.Size() != workerPoolSmall+1 {
		t.Errorf("Expect pool size of %d, got: %d - After too large shrink", workerPoolSmall+1, p.Size())
	}

	// Check we can shrink by 1
	if err := p.Shrink(1); err != nil {
		t.Errorf("Errored when shrinking by 1: %v", err)
	}
	// Check the new size
	if p.Size() != workerPoolSmall {
		t.Errorf("Expect pool size of %d, got: %d - After shrink by 1", workerPoolSmall, p.Size())
	}

	// Test shrinking by the size closes the pool
	if p.Shrink(p.Size()) != nil {
		t.Errorf("Failed to shrink by size of pool.  Should close pool ctx")
	}
}

// Check how fast a small work/job queue pool runs
func BenchmarkNewFlexPool(b *testing.B) {
	p := NewFlexPool(workerPoolSmall, queueSizeSmall)

	j := testJob{}
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		p.Add(j)
	}
	wg.Wait()
}

// Check how fast a large work/job queue pool runs
func BenchmarkNewFlexPoolLarge(b *testing.B) {
	p := NewFlexPool(workerPoolLarge, queueSizeLarge)

	j := testJob{}
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		p.Add(j)
	}
	wg.Wait()
}
