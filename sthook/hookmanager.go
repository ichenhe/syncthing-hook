package sthook

import (
	"SyncthingHook/eventhandler"
	"SyncthingHook/extevent"
	"SyncthingHook/stclient"
	"errors"
	"fmt"
	"github.com/syncthing/syncthing/lib/events"
	"go.uber.org/zap"
	"strings"
)

type HookManager struct {
	st         *stclient.Syncthing
	logger     *zap.SugaredLogger
	stEventBus *stEventBus
	hooks      map[stEventBusExposedSubscriber]*Hook
}

func NewHookManager(st *stclient.Syncthing, logger *zap.Logger) *HookManager {
	return &HookManager{
		st:         st,
		logger:     logger.Sugar(),
		stEventBus: newStEventBus(st, logger.Sugar()),
		hooks:      make(map[stEventBusExposedSubscriber]*Hook),
	}
}

func (m *HookManager) RegisterHook(hook *Hook, hookIndex int) error {
	var subscriber stEventBusExposedSubscriber
	if strings.HasPrefix(hook.EventType, "st:") {
		stEventType := events.UnmarshalEventType(hook.EventType[3:])
		if stEventType == 0 {
			return newEventTypeError(fmt.Sprintf("unknown syncthing event type '%s'", stEventType))
		}
		subscriber = m.stEventBus.subscribe(stEventType)
	} else if strings.HasPrefix(hook.EventType, "ex:") {
		detector, err := m.newEventDetector(m.st, hook.EventType, hook.Parameter)
		if err != nil {
			return err
		}
		subscriber = detector.Subscribe()
	} else {
		return newEventTypeError(fmt.Sprintf("unknown event type classification '%s'", hook.EventType))
	}

	if subscriber != nil {
		headHandler := m.buildHandlerChain(
			m.createFilters(hook.Parameter, hook.Name, hookIndex),
			m.createExecutors(hook, hookIndex),
		)
		m.listeningSubscribeChannel(subscriber, headHandler)
	}
	return nil
}

// createFilters creates filter-type event handlers. Returns nil or empty if given params is nil or no filter needed.
func (m *HookManager) createFilters(params HookParameters, hookName string, hookIndex int) []eventhandler.Handler {
	if params == nil {
		return nil
	}
	logger := m.logger.With(zap.String("hookName", hookName), zap.Int("hookIndex", hookIndex))
	handlers := make([]eventhandler.Handler, 0)
	if params.containsKey("tolerance") {
		handlers = append(handlers, eventhandler.NewToleranceFilter(params.getInt64("tolerance", 0), logger))
	}
	if params.containsKey("cooldown") {
		handlers = append(handlers, eventhandler.NewCoolDownFilter(params.getInt64("cooldown", 0), logger))
	}
	return handlers
}

// createExecutors creates executor-type event handlers. Returns empty list if no valid action.
func (m *HookManager) createExecutors(hook *Hook, hookIndex int) []eventhandler.Handler {
	handlers := make([]eventhandler.Handler, 0)
	logger := m.logger.With(zap.String("hookName", hook.Name), zap.Int("hookIndex", hookIndex))
	switch hook.Action.Type {
	case "exec":
		handlers = append(handlers, eventhandler.NewExecExecutor(hook.Action.Cmd, logger))
	default:
		m.logger.Infow("unknown action type, ignore.", "actionType", hook.Action.Type, "hookName", hook.Name)
	}
	return handlers
}

// buildHandlerChain builds the given handlers into a chain and return the first handler. Returns nil if given handlers
// is nil or empty.
func (m *HookManager) buildHandlerChain(handlers ...[]eventhandler.Handler) eventhandler.Handler {
	if handlers == nil || len(handlers) == 0 {
		return nil
	}
	l := make([]eventhandler.Handler, 0, len(handlers))
	for _, hs := range handlers {
		for _, h := range hs {
			if h != nil {
				l = append(l, h)
			}
		}
	}
	for i, handler := range l {
		if i == len(l)-1 {
			continue
		}
		handler.SetNext(l[i+1])
	}
	return l[0]
}

func (m *HookManager) newEventDetector(st *stclient.Syncthing, extEventType extevent.EventType, parameters HookParameters) (extevent.EventDetector, error) {
	if st == nil {
		return nil, errors.New("'syncthing' can not be nil")
	}
	switch extEventType {
	case extevent.LocalFolderContentChangeDetected:
		return extevent.NewLocalFolderContentChangedDetector(st, parameters.getString("st-folder", ""),
			func(option *extevent.LocalFolderContentChangedDetectorOptions) {
				parameters.extractStringIfExist("path", &option.Path)
				parameters.extractIntIfExist("since", &option.Since)
			})
	default:
		return nil, newEventTypeError(fmt.Sprintf("unknown ext event type '%s'", extEventType))
	}
}

func (m *HookManager) listeningSubscribeChannel(subscriber stEventBusExposedSubscriber, handlerChain eventhandler.Handler) {
	for e := range subscriber {
		if handlerChain != nil {
			handlerChain.Handle(&e)
		}
	}
	// destroy
	_handler := handlerChain
	for _handler != nil {
		_handler.Destroy()
		_handler = _handler.GetNext()
	}
}
