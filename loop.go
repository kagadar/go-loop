package loop

import (
	"context"
	"math"
	"runtime"
	"time"
)

type Main interface {
	EnableStats(bool)
	Stats() Stats
	SetDelay(time.Duration)
	Run(context.Context, func(context.Context))
}

type main struct {
	delay         time.Duration
	loop          *time.Ticker
	stats         Stats
	statsUpdateFn func(context.Context, *Stats, time.Time)
}

func noop(context.Context, *Stats, time.Time) {}

func updateStats(ctx context.Context, stats *Stats, tick time.Time) {
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

func New(delay time.Duration) Main {
	return &main{
		delay: delay,
		stats: Stats{
			start:    time.Now(),
			lastTick: time.Now(),
			maxTick:  math.MinInt64,
			minTick:  math.MaxInt64,
		},
		statsUpdateFn: noop,
	}
}

func (m *main) EnableStats(b bool) {
	if b {
		m.statsUpdateFn = updateStats
	} else {
		m.statsUpdateFn = noop
	}
}

func (m *main) Stats() Stats {
	return m.stats
}

func (m *main) SetDelay(d time.Duration) {
	if m.loop != nil {
		m.loop.Reset(d)
	}
	m.delay = d
}

func (m *main) Run(ctx context.Context, fn func(context.Context)) {
	cleanup := configureTimer()
	defer cleanup()
	m.loop = time.NewTicker(m.delay)
	defer m.loop.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		m.statsUpdateFn(ctx, &m.stats, <-m.loop.C)
		fn(ctx)
	}
}
