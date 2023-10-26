package exevent

import (
	"SyncthingHook/domain"
	"go.uber.org/zap"
)

// wrapLogger adds 'tag' field to Logger
func wrapLogger(logger *zap.SugaredLogger, name string) *zap.SugaredLogger {
	return logger.With(zap.String("tag", name))
}

// baseHandler provides convenience method to call next handler safely.
//
// Note: this does not fully implement domain.EventHandler interface. Typically, the real implementation should contain
// this as an embed field.
type baseHandler struct {
	next domain.EventHandler
}

func (h *baseHandler) SetNext(next domain.EventHandler) {
	h.next = next
}

func (h *baseHandler) GetNext() domain.EventHandler {
	return h.next
}

func (h *baseHandler) Destroy() {
}

func (h *baseHandler) callNext(event *domain.Event) {
	if h.next != nil {
		h.next.Handle(event)
	}
}
