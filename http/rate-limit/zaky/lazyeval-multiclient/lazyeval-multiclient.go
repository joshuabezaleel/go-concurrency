package lazyevalmulticlient

import (
	"sync"
	"time"
)

type throttleKeyFunc func(key string) bool

// throttle limits the request to at most N requests / second *for each key*
// this function returns a func(key string) bool which indicates
// subsequent request can be safely made or not for that certain key.
func throttleKey(n int) throttleKeyFunc {
	var lock sync.Mutex
	lastRequest := make(map[string]time.Time)
	counter := make(map[string]int)
	return func(key string) bool {
		lock.Lock()
		defer lock.Unlock()
		now := time.Now()
		t, okt := lastRequest[key]
		c, okc := counter[key]
		if !okt || !okc || now.Sub(t) > time.Second {
			lastRequest[key] = now
			counter[key] = 1
			return true
		}
		c++
		counter[key] = c
		return c <= n
	}
}
