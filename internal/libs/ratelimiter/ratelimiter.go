package ratelimiter

import (
	"sync"
	"time"
)

type RateLimiter struct {
	limit    int
	time     int
	mu       sync.Mutex
	requests map[string]int
}

func NewRateLimiter(limit, time int) *RateLimiter {
	rl := &RateLimiter{
		requests: make(map[string]int),
		limit:    limit,
		time:     time,
	}
	go rl.Clear()
	return rl
}

func (rl *RateLimiter) Allow(apiKey string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.requests[apiKey] >= rl.limit {
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
	ticker := time.NewTicker(time.Duration(rl.time) * time.Millisecond)
	defer ticker.Stop()

	for {
		<-ticker.C
		rl.Reset()
	}
}
