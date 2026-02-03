package terminalpane

import "time"

type FrameLimiter struct {
	minInterval time.Duration
	last        time.Time
}

func NewFrameLimiter(maxFPS int, minInterval time.Duration) *FrameLimiter {
	interval := minInterval
	if maxFPS > 0 {
		fpsInterval := time.Second / time.Duration(maxFPS)
		if interval == 0 || fpsInterval > interval {
			interval = fpsInterval
		}
	}
	return &FrameLimiter{minInterval: interval}
}

func (f *FrameLimiter) Allow() bool {
	now := time.Now()
	if f.minInterval == 0 {
		f.last = now
		return true
	}
	if now.Sub(f.last) >= f.minInterval {
		f.last = now
		return true
	}
	return false
}
