package lazyevalsingleclient

import (
	"sync"
	"time"
)

type throttleFunc func() bool

func throttle(n int) throttleFunc {
	var lock sync.Mutex
	var checkpoint time.Time = time.Now()
	var counter = 0
	return func() bool {
		lock.Lock()
		defer lock.Unlock()
		now := time.Now()

		// it's been more than 1 second since the last checkpoint
		// then we reset checkpoint and counter
		if now.Sub(checkpoint) > time.Second {
			checkpoint = now
			counter = 1
			return true
		}
		counter++
		return counter <= n
	}
}
