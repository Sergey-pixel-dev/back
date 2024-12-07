package ratelimiter

import (
	"sync"
	"time"
)

type RateLimiter struct {
	mu       sync.Mutex
	requests map[string]int
}

func NewRateLimiter() *RateLimiter {
	rl := &RateLimiter{
		requests: make(map[string]int),
	}
	go rl.Clear()
	return rl
}

func (rl *RateLimiter) Allow(apiKey string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.requests[apiKey] >= 5 {
		return false
	}

	rl.requests[apiKey]++
	return true
}

func (rl *RateLimiter) Reset() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	for k := range rl.requests {
		rl.requests[k] = 0
	}
}
func (rl *RateLimiter) Clear() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		<-ticker.C
		rl.Reset()
	}
}
