package mocker

import (
	"SyncthingHook/domain"
	"testing"
)

// AssertNextHandlerCalled asserts the Handle method of handler was called.
func (next *MockEventHandler) AssertNextHandlerCalled(t *testing.T, ev *domain.Event, number int) {
	next.AssertCalled(t, "Handle", ev)
	next.AssertNumberOfCalls(t, "Handle", number)
}

func (next *MockEventHandler) AssertNextHandlerNotCalled(t *testing.T, ev *domain.Event) {
	next.AssertNotCalled(t, "Handle", ev)
}
