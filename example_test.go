package pooligo

import (
	"fmt"
	"sync"
	"time"
)

var (
	exampleWG sync.WaitGroup
)

// delayprinter is the struct used to perform asynchonous jobs on.
// This will be used as an example of using a worker pool
type delayprinter struct {
	id    int
	delay time.Duration
}

// Action is the function called by the workers from the pool
func (d *delayprinter) Action() {
	time.Sleep(d.delay)
	fmt.Printf("Performed action on ID: %d\n", d.id)

	exampleWG.Done()
}

func ExampleNewPool() {
	// Create two delayprinters
	d1 := &delayprinter{1, time.Second * 2}
	d2 := &delayprinter{2, time.Second}

	// Create a worker pool
	pool := NewPool(10, 5)

	exampleWG.Add(2)

	// Add the jobs to the worker pool
	pool.Add(d1)
	pool.Add(d2)

	// Wait for the goroutines to complete
	exampleWG.Wait()

	// Output:
	// Performed action on ID: 2
	// Performed action on ID: 1
}
