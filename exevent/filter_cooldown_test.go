package exevent

import (
	"github.com/ichenhe/syncthing-hook/domain"
	"github.com/ichenhe/syncthing-hook/mocker"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
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
			ev := &domain.Event{}
			next := &mocker.MockEventHandler{}
			next.EXPECT().Handle(ev).Return()
			timeProvider := &mocker.MockTimeProvider{}
			h := NewCoolDownFilter(tt.coolDown, zap.NewNop().Sugar(), CoolDownFilterTimeProvider(timeProvider))
			h.SetNext(next)
			h.lastTime = 1000
			newSendTime := 1000 + tt.msAfterLastSending
			timeProvider.On("NowUnixMilli").Return(newSendTime)
			h.Handle(ev)

			if tt.shouldCallNext {
				next.AssertNextHandlerCalled(t, ev, 1)
				if tt.coolDown != 0 {
					// should update the last send time
					require.Equal(t, newSendTime, h.lastTime)
				}
			} else {
				next.AssertNextHandlerNotCalled(t, ev)
			}
		})
	}
}
