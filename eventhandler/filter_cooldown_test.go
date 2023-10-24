package eventhandler

import (
	"SyncthingHook/extevent"
	"SyncthingHook/utils/clock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCoolDownFilter_Handle(t *testing.T) {
	tests := []struct {
		name               string
		coolDown           int64
		msAfterLastSending int64
		shouldCallNext     bool
	}{
		{name: "no cool down", coolDown: 0, msAfterLastSending: 0, shouldCallNext: true},
		{name: "cooling down", coolDown: 500, msAfterLastSending: 480, shouldCallNext: false},
		{name: "ready", coolDown: 500, msAfterLastSending: 510, shouldCallNext: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ev := &extevent.Event{}
			next := createMockNextHandler(ev)
			logger, _ := newLogger()
			timeProvider := &clock.MockTimeProvider{}
			h := NewCoolDownFilter(tt.coolDown, logger, CoolDownFilterTimeProvider(timeProvider))
			h.SetNext(next)
			h.lastTime = 1000
			newSendTime := 1000 + tt.msAfterLastSending
			timeProvider.On("NowUnixMilli").Return(newSendTime)
			h.Handle(ev)

			if tt.shouldCallNext {
				next.assertNextHandlerCalled(t, ev, 1)
				if tt.coolDown != 0 {
					// should update the last send time
					require.Equal(t, newSendTime, h.lastTime)
				}
			} else {
				next.assertNextHandlerNotCalled(t, ev)
			}
		})
	}
}
