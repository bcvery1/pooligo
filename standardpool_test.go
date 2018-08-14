package pooligo

import (
	"testing"
)

func TestNewPool(t *testing.T) {
	if err := testPoolCreation(NewPool(workerPoolSmall, queueSizeSmall)); err != nil {
		t.Errorf(err.Error())
	}
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
