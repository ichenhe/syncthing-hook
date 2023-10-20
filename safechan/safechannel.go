package safechan

import (
	"sync"
)

// SafeChannel can send or be closed safely. Do not close the inner channel directly.
type SafeChannel[T any] struct {
	c      chan T
	closed bool
	mutex  sync.Mutex
}

// OutCh returns a native read-only channel.
func (c *SafeChannel[T]) OutCh() <-chan T {
	return c.c
}

func NewSafeChannel[T any]() *SafeChannel[T] {
	return &SafeChannel[T]{c: make(chan T)}
}

func (c *SafeChannel[T]) SafeClose() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if !c.closed {
		close(c.c)
		c.closed = true
	}
}

// SafeSend sends given value safely which means there's no panic even if channel is closed. Return false if it has been
// closed.
func (c *SafeChannel[T]) SafeSend(value T) (ok bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if c.closed {
		return false
	}
	c.c <- value
	return true
}
