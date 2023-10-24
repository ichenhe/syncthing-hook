package clock

import "time"

type TimeProvider interface {
	NowUnixMilli() int64
}

// DefaultTimeProvider just uses builtin time package.
type DefaultTimeProvider struct {
}

func (d *DefaultTimeProvider) NowUnixMilli() int64 {
	return time.Now().UnixMilli()
}
