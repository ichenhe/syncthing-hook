package exevent

import (
	"github.com/ichenhe/syncthing-hook/domain"
	"github.com/ichenhe/syncthing-hook/mocker"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"testing"
	"time"
)

func TestToleranceFilter_Handle(t *testing.T) {
	tests := []struct {
		name       string
		eventTypes []bool // true means this event should be processed. the length represents total number of events
		intervalMs []int  // the length must equal to len(eventTypes)-1
		tolerance  int64
	}{
		{name: "no tolerance", eventTypes: []bool{true, true, true, true, true}, intervalMs: []int{0, 0, 0, 0}, tolerance: 0},
		{name: "override previous event", eventTypes: []bool{false, false, false, true}, intervalMs: []int{20, 20, 20}, tolerance: 30},
		{name: "mix", eventTypes: []bool{true, false, false, true}, intervalMs: []int{50, 20, 20}, tolerance: 30},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, len(tt.eventTypes)-1, len(tt.intervalMs), "len(intervalMs) must equal to len(eventTypes)-1")

			ev := &domain.Event{}
			ev2 := &domain.Event{Type: domain.LocalFolderContentChangeDetected}
			next := &mocker.MockEventHandler{}
			next.EXPECT().Handle(ev).Return()
			next.On("Handle", ev2).Return()
			h := NewToleranceFilter(tt.tolerance, zap.NewNop().Sugar())
			h.SetNext(next)
			callNextTimes := 0
			for i, static := range tt.eventTypes {
				if static {
					callNextTimes++
					h.Handle(ev)
				} else {
					h.Handle(ev2)
				}
				if i < len(tt.intervalMs) && tt.intervalMs[i] > 0 {
					time.Sleep(time.Millisecond * time.Duration(tt.intervalMs[i]))
				}
			}
			time.Sleep(time.Millisecond * (time.Duration(tt.tolerance) + 100))

			next.AssertNextHandlerCalled(t, ev, callNextTimes)
			next.AssertNextHandlerNotCalled(t, ev2)
		})
	}
}

func TestToleranceFilter_Destroy(t *testing.T) {
	tests := []struct {
		name             string
		timerRemainingMs int
	}{
		{name: "no timer", timerRemainingMs: -1},
		{name: "timer already expired", timerRemainingMs: 0},
		{name: "timer running", timerRemainingMs: 5000},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var timer *time.Timer
			if tt.timerRemainingMs == -1 {
				timer = nil
			} else if tt.timerRemainingMs == 0 {
				timer = time.NewTimer(0 * time.Nanosecond)
			} else {
				timer = time.NewTimer(time.Millisecond * time.Duration(tt.timerRemainingMs))
			}
			h := NewToleranceFilter(0, zap.NewNop().Sugar())
			h.timer = timer
			h.Destroy()
			require.Nil(t, h.timer)
			if timer != nil {
				require.False(t, timer.Stop(), "timer must have stopped")
				require.Empty(t, timer.C)
			}
		})
	}
}
