package eventhandler

import (
	"SyncthingHook/extevent"
	"go.uber.org/zap"
)

type Handler interface {
	SetNext(next Handler)
	GetNext() Handler
	// Handle processes the event. It should call next handler in most cases.
	Handle(event *extevent.Event)
	// Destroy will be called when the handler will not be used anymore.
	// Typically, it should cancel any ongoing coroutines or release other resources.
	//
	// Do not call next handler.
	Destroy()
}

// wrapLogger adds 'tag' field to logger
func wrapLogger(logger *zap.SugaredLogger, name string) *zap.SugaredLogger {
	return logger.With(zap.String("tag", name))
}

// baseHandler provides convenience method to call next handler.
type baseHandler struct {
	next Handler
}

func (h *baseHandler) SetNext(next Handler) {
	h.next = next
}

func (h *baseHandler) GetNext() Handler {
	return h.next
}

func (h *baseHandler) Destroy() {
}

func (h *baseHandler) callNext(event *extevent.Event) {
	if h.next != nil {
		h.next.Handle(event)
	}
}
