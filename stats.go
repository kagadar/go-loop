package loop

import (
	"fmt"
	"runtime"
	"time"
)

func fmtBytes(b uint64) string {
	if b < 1024 {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := 1024, 0
	for n := b / 1024; n >= 1024; n /= 1024 {
		div *= 1024
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(b)/float64(div), "KMGTPE"[exp])
}

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
func (s Stats) Duration() time.Duration {
	return time.Since(s.start)
}

// MaxTick returns the longest delay between two iterations of Main.
func (s Stats) MaxTick() time.Duration {
	return s.maxTick
}

// MinTick returns the shortest delay between two iterations of Main.
func (s Stats) MinTick() time.Duration {
	return s.minTick
}

// AvgTick returns the average delay between the iterations of Main.
func (s Stats) AvgTick() time.Duration {
	return time.Duration(int64(s.total) / s.ticks)
}

func (s Stats) HeapAllocFmt() string {
	return fmtBytes(s.HeapAlloc)
}
