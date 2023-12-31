package clock

import (
	"github.com/ichenhe/syncthing-hook/domain"
	"time"
)

// TimeProvider just uses builtin time package.
type TimeProvider struct {
}

var _ domain.TimeProvider = (*TimeProvider)(nil)

func (d *TimeProvider) NowUnixMilli() int64 {
	return time.Now().UnixMilli()
}
