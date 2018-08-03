# pooligo
Package pooligo provides a basic worker-pool.  This can be used to restrict the number of threads created by an application while providing multi-threaded capabilities.

## Usage
```go
import (
  "github.com/bcvery1/pooligo"
)

const (
  fooCount = 1000
)

var (
  wg sync.WaitGroup
)

type Foo struct {}

// Provide an `Action` function for the structure
func (f *Foo) Action() {
  defer wg.Done()
  // ...
  // Code to be performed on f
  // ...
}

func main() {
  p := pooligo.NewPool(5, 10)

  for i := 0; i < fooCount; i++ {
    wg.Add(1)
    f := &Foo{}
    p.Add(f)
  }

  // Wait for pool to finish, or run indefinitely
  wg.Wait()
}
```
