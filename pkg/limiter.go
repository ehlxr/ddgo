package pkg

import (
	"sync"
	"time"
)

type LimiterServer struct {
	interval time.Duration
	maxCount int
	sync.Mutex
	reqCount  int
	startTime time.Time
}

func NewLimiterServer(i time.Duration, c int) *LimiterServer {
	return &LimiterServer{
		interval: i,
		maxCount: c,
	}
}

func (limiter *LimiterServer) IsAvailable() bool {
	limiter.Lock()
	defer limiter.Unlock()

	if limiter.startTime.IsZero() ||
		limiter.startTime.Add(limiter.interval).Before(time.Now()) {
		limiter.reqCount = 1
		limiter.startTime = time.Now()

		return true
	}

	if limiter.reqCount < limiter.maxCount {
		limiter.reqCount += 1

		return true
	}

	return false
}
