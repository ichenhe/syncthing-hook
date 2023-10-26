package domain

import "errors"

// EventHandler implements the chain of responsibility pattern, one event will be processed by different handlers
// one by one. The handler may terminate event propagation (e.g. filter).
type EventHandler interface {
	SetNext(next EventHandler)
	GetNext() EventHandler
	// Handle processes the event. It should call next handler in most cases.
	Handle(event *Event)
	// Destroy will be called when the handler will not be used anymore.
	// Typically, it should cancel any ongoing coroutines or release other resources.
	//
	// Do not call next handler.
	Destroy()
}

var ErrInvalidActionType = errors.New("invalid action type")

type ActionExecutorCreator interface {
	CreateExecutor(action *HookAction, hookDef *HookDefinition) (EventHandler, error)
}

// BuildEventHandlerChain connects all given handlers into a chain, and returns the first one.
// If there is no handler, return nil.
// This function doesn't modify the tail handler.
func BuildEventHandlerChain(handlers ...EventHandler) EventHandler {
	if len(handlers) == 0 {
		return nil
	}
	for i := 0; i < len(handlers)-1; i++ {
		handlers[i].SetNext(handlers[i+1])
	}
	return handlers[0]
}

func DestroyHandlerChain(handlerChain EventHandler) {
	h := handlerChain
	for h != nil {
		h.Destroy()
		h = h.GetNext()
	}
}
