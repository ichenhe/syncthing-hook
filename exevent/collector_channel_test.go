package exevent

import (
	"github.com/ichenhe/syncthing-hook/domain"
	"github.com/ichenhe/syncthing-hook/mocker"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func Test_applyHandlers(t *testing.T) {
	tests := []struct {
		name       string
		handlerNum int
	}{
		{name: "1 handler", handlerNum: 1},
		{name: "10 handlers", handlerNum: 10},
		{name: "no handler", handlerNum: 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nextMap := make(map[domain.EventHandler]domain.EventHandler)
			ev := &domain.Event{}
			handlers := make([]domain.EventHandler, tt.handlerNum)
			for i := 0; i < tt.handlerNum; i++ {
				h := mocker.NewMockEventHandler(t)
				h.EXPECT().SetNext(mock.Anything).Run(func(next domain.EventHandler) {
					nextMap[h] = next.(domain.EventHandler)
				})
				h.EXPECT().GetNext().RunAndReturn(func() domain.EventHandler {
					return nextMap[h]
				})
				h.EXPECT().Handle(ev).Run(func(event *domain.Event) {
					if n := h.GetNext(); n != nil {
						n.Handle(event)
					}
				})
				h.EXPECT().Destroy().Return()
				handlers[i] = h
			}
			eventCh := make(chan *domain.Event)
			processedCh := applyHandlers(eventCh, handlers...)
			if tt.handlerNum == 0 {
				require.EqualValues(t, eventCh, processedCh, "if no handler is given, the original ch should be returned")
				return
			}
			go func() {
				eventCh <- ev
				close(eventCh)
			}()
			done := make(chan struct{}, 1)
			go func() {
				for e := range processedCh {
					require.Equal(t, e, ev)
				}
				done <- struct{}{}
			}()
			wait := time.Millisecond * 500
			select {
			case <-done:
				return
			case <-time.After(wait):
				t.Errorf("returned collector channel is not closed after %v", wait)
				t.FailNow()
			}
		})
	}
}

func TestChannelCollector_Handle(t *testing.T) {
	ev := &domain.Event{}
	next := mocker.NewMockEventHandler(t)
	next.EXPECT().Handle(ev).Return() // should call next handler
	h := NewChannelCollector()
	h.SetNext(next)
	done := make(chan bool, 1)
	go func() {
		result := false
		for e := range h.GetCh() {
			if e == ev {
				result = true
			}
		}
		done <- result
	}()
	h.Handle(ev)
	h.Destroy()
	wait := time.Millisecond * 500
	select {
	case ok := <-done:
		require.True(t, ok, "the channel should receive the event")
	case <-time.After(wait):
		t.Errorf("returned collector channel is not closed after %v", wait)
		t.FailNow()
	}
}

func TestChannelCollector_Destroy(t *testing.T) {
	h := NewChannelCollector()
	ch := h.GetCh()
	h.Destroy()
	require.Nil(t, h.ch)
	select {
	case _, ok := <-ch:
		require.False(t, ok, "the channel should have been closed")
	case <-time.After(time.Millisecond * 200):
		t.Error("the channel is not closed")
		t.FailNow()
	}
}
