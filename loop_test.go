package loop

import (
	"context"
	"testing"
	"time"
)

func TestMainLoop(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	loop := New(1 * time.Millisecond)
	loop.EnableStats(true)
	loop.Run(ctx, func(ctx context.Context) {
		cancel()
	})
}
