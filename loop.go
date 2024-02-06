package loop

import (
	"context"
	"time"
)

func MainLoop(ctx context.Context, tickrate time.Duration, fn func(context.Context)) {
	cleanup := configureTimer()
	defer cleanup()
	loop := time.NewTicker(tickrate)
	defer loop.Stop()
MAIN_LOOP:
	for {
		select {
		case <-ctx.Done():
			loop.Stop()
			break MAIN_LOOP
		default:
		}
		fn(ctx)
	}
}
