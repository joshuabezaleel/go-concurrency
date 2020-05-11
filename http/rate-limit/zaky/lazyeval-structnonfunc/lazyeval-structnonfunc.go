package lazyevalstructnonfunc

import (
	"sync"
	"time"
)

type Throttler struct {
	lock       sync.Mutex
	checkpoint time.Time
	counter    int
	limit      int
}

func NewThrottler(limit int) *Throttler {
	return &Throttler{
		checkpoint: time.Now(),
		counter:    0,
		limit:      limit,
	}
}

func (t *Throttler) Allow(n int) bool {
	t.lock.Lock()
	defer t.lock.Unlock()
	now := time.Now()

	// it's been more than 1 second since the last checkpoint
	// then we reset checkpoint and counter
	if now.Sub(t.checkpoint) > time.Second {
		t.checkpoint = now
		t.counter = 1
		return true
	}
	t.counter++
	return t.counter <= n
}
