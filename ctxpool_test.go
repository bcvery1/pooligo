package pooligo

import (
	"context"
	"testing"
)

func TestNewCtxPool(t *testing.T) {
	if err := testPoolCreation(NewCtxPool(context.Background(), workerPoolSmall, queueSizeSmall)); err != nil {
		t.Errorf(err.Error())
	}
}

// Check how fast a small work/job queue pool runs
func BenchmarkNewCtxPool(b *testing.B) {
	p := NewCtxPool(context.Background(), workerPoolSmall, queueSizeSmall)

	j := testJob{}
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		p.Add(j)
	}
	wg.Wait()
}

// Check how fast a large work/job queue pool runs
func BenchmarkNewCtxPoolLarge(b *testing.B) {
	p := NewCtxPool(context.Background(), workerPoolLarge, queueSizeLarge)

	j := testJob{}
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		p.Add(j)
	}
	wg.Wait()
}
