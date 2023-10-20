package sthook

import (
	"SyncthingHook/extevent"
	"SyncthingHook/stclient"
	"errors"
	"fmt"
	"github.com/syncthing/syncthing/lib/events"
	"go.uber.org/zap"
	"os/exec"
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

func (m *HookManager) RegisterHook(hook *Hook) error {
	if strings.HasPrefix(hook.EventType, "st:") {
		stEventType := events.UnmarshalEventType(hook.EventType[3:])
		if stEventType == 0 {
			return newEventTypeError(fmt.Sprintf("unknown syncthing event type '%s'", stEventType))
		}
		m.listeningSubscribeChannel(m.stEventBus.subscribe(stEventType), hook)
	} else if strings.HasPrefix(hook.EventType, "ex:") {
		detector, err := m.newEventDetector(m.st, hook.EventType, hook.Parameter)
		if err != nil {
			return err
		}
		m.listeningSubscribeChannel(detector.Subscribe(), hook)
	} else {
		return newEventTypeError(fmt.Sprintf("unknown event type classification '%s'", hook.EventType))
	}
	return nil
}

func (m *HookManager) newEventDetector(st *stclient.Syncthing, extEventType extevent.EventType, parameters HookParameters) (extevent.EventDetector, error) {
	if st == nil {
		return nil, errors.New("'syncthing' can not be nil")
	}
	switch extEventType {
	case extevent.LocalFolderContentChangeDetected:
		return extevent.NewLocalFolderContentChangedDetector(st, parameters.getString("st-folder", ""),
			func(option *extevent.LocalFolderContentChangedDetectorOption) {
				parameters.extractStringIfExist("path", &option.Path)
				parameters.extractIntIfExist("since", &option.Since)
				parameters.extractInt64IfExist("cooldown", &option.Cooldown)
				parameters.extractInt64IfExist("tolerance", &option.Tolerance)
			})
	default:
		return nil, newEventTypeError(fmt.Sprintf("unknown ext event type '%s'", extEventType))
	}
}

func (m *HookManager) listeningSubscribeChannel(subscriber stEventBusExposedSubscriber, hook *Hook) {
	for e := range subscriber {
		m.logger.Debug("%+v", e)
		m.executeAction(e, hook)
	}
}

func (m *HookManager) executeAction(event extevent.Event, hook *Hook) {
	if hook.Action.Type != "exec" {
		m.logger.Debugw("unknown action type, ignore.", "actionType", hook.Action.Type, "hookName", hook.Name)
		return
	}
	if len(hook.Action.Cmd) == 0 || len(hook.Action.Cmd[0]) == 0 {
		m.logger.Debugw("cmd is empty, ignore this exec action.", "hookName", hook.Name)
		return
	}
	if err := exec.Command(hook.Action.Cmd[0], hook.Action.Cmd[1:]...).Start(); err != nil {
		m.logger.Infow("filed to execute action: "+err.Error(), "cmd", hook.Action.Cmd)
	}
}
