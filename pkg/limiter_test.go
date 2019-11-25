package pkg

import (
	"testing"
	"time"
)

func TestLimiter(t *testing.T) {
	limiter := NewLimiterServer(1*time.Second, 5)

	for {
		if limiter.IsAvailable() {
			t.Log("hello...", limiter.reqCount)
		} else {
			return
		}
	}
}
