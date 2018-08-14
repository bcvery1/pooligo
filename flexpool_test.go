package pooligo

import (
	"testing"
)

func TestNewFlexPool(t *testing.T) {
	if err := testPoolCreation(NewFlexPool(workerPoolSmall, queueSizeSmall)); err != nil {
		t.Errorf(err.Error())
	}
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
