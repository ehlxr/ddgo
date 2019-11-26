package pkg

import (
	"sync"
	"time"
)

type LimiterServer struct {
	interval time.Duration
	maxCount int
	sync.Mutex
	reqCount int
	time     time.Time
}

func NewLimiterServer(interval time.Duration, maxCount int) *LimiterServer {
	return &LimiterServer{
		interval: interval,
		maxCount: maxCount,
	}
}

func (limiter *LimiterServer) IsAvailable() bool {
	limiter.Lock()
	defer limiter.Unlock()
	now := time.Now()

	if limiter.time.IsZero() ||
		limiter.time.Add(limiter.interval).Before(now) {
		limiter.reqCount = 0
	}

	if limiter.reqCount < limiter.maxCount {
		limiter.reqCount += 1
		limiter.time = now
		return true
	}

	return false
}
