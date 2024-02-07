# Go Loop

Provides a main loop that executes using a configurable ticker.
If the provided context is cancelled, the loop will stop.

## Getting Started

### Installing

```sh
go get github.com/kagadar/go-loop
```

### Usage

```go
import (
    "context"
    "os"
    "os/signal"
    "time"

    "github.com/kagadar/go-loop"
)

func main() {
    ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
    defer stop()
    l := loop.New(16 * time.Millisecond)
    l.Run(ctx, func(ctx context.Context) {
        // Do something
    })
}
```
