package loop

import (
	"context"
	"math"
	"runtime"
	"time"
)

type statsKey struct{}

type Stats struct {
	runtime.MemStats
	start    time.Time
	lastTick time.Time
	maxTick  time.Duration
	minTick  time.Duration
	total    time.Duration
	ticks    int64
}

// Duration returns the duration that Main has been running for.
func (s *Stats) Duration() time.Duration {
	return time.Since(s.start)
}

// MaxTick returns the longest delay between two iterations of Main.
func (s *Stats) MaxTick() time.Duration {
	return s.maxTick
}

// MinTick returns the shortest delay between two iterations of Main.
func (s *Stats) MinTick() time.Duration {
	return s.minTick
}

// AvgTick returns the average delay between the iterations of Main.
func (s *Stats) AvgTick() time.Duration {
	return time.Duration(int64(s.total) / s.ticks)
}

// GetStats returns the Stats from the provided context.
// GetStats will panic if the context was not provided by a Main loop with the `Stats` option set.
func GetStats(ctx context.Context) *Stats {
	return ctx.Value(statsKey{}).(*Stats)
}

type Options struct {
	// The minimum time that must elapse between each iteration of Main.
	Delay time.Duration
	// Whether to calculate runtime stats and store them in the context used by Main.
	Stats bool
}

func noopStats(context.Context, time.Time) {}

func updateStats(ctx context.Context, tick time.Time) {
	stats := GetStats(ctx)
	stats.ticks++
	tickTime := tick.Sub(stats.lastTick)
	stats.total += tickTime
	if stats.maxTick <= tickTime {
		stats.maxTick = tickTime
	}
	if stats.minTick >= tickTime {
		stats.minTick = tickTime
	}
	stats.lastTick = tick
	runtime.ReadMemStats(&stats.MemStats)
}

func Main(ctx context.Context, opts Options, fn func(context.Context)) {
	stats := noopStats
	if opts.Stats {
		ctx = context.WithValue(ctx, statsKey{}, &Stats{
			start:    time.Now(),
			lastTick: time.Now(),
			maxTick:  math.MinInt64,
			minTick:  math.MaxInt64,
		})
		stats = updateStats
	}
	cleanup := configureTimer()
	defer cleanup()
	loop := time.NewTicker(opts.Delay)
	defer loop.Stop()
LOOP:
	for {
		select {
		case <-ctx.Done():
			loop.Stop()
			break LOOP
		default:
		}
		stats(ctx, <-loop.C)
		fn(ctx)
	}
}
