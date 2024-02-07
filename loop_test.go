package loop

import (
	"context"
	"testing"
	"time"
)

func TestMainLoop(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	Main(ctx, Options{Delay: 1 * time.Millisecond, Stats: true}, func(ctx context.Context) {
		cancel()
	})
}
