package eventhandler

import (
	"SyncthingHook/extevent"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
	"testing"
)

// newLogger creates zap logger with observer for testing, so that we can verify logs.
func newLogger() (*zap.SugaredLogger, *observer.ObservedLogs) {
	observedZapCore, observedLogs := observer.New(zap.DebugLevel)
	logger := zap.New(observedZapCore).Sugar()
	return logger, observedLogs
}

// createMockNextHandler creates a mock Handler that records the call of Handle method.
func createMockNextHandler(ev *extevent.Event) *MockHandler {
	next := &MockHandler{}
	next.On("Handle", ev).Return()
	return next
}

// assertNextHandlerCalled asserts the Handle method of handler was called.
func (next *MockHandler) assertNextHandlerCalled(t *testing.T, ev *extevent.Event, number int) {
	next.AssertCalled(t, "Handle", ev)
	next.AssertNumberOfCalls(t, "Handle", number)
}

func (next *MockHandler) assertNextHandlerNotCalled(t *testing.T, ev *extevent.Event) {
	next.AssertNotCalled(t, "Handle", ev)
}
