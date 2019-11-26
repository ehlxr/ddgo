package pkg

import (
	"fmt"
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

func TestLimiter2(t *testing.T) {
	limiter := NewLimiterServer(10*time.Second, 10)

	for i := 0; i < 20; i++ {
		if limiter.IsAvailable() {
			fmt.Println("hello...", limiter.reqCount)
		} else {
			fmt.Println("limited")
		}
		time.Sleep(1 * time.Second)
	}
}
