package hook

import (
	"fmt"
	"github.com/ichenhe/syncthing-hook/domain"
	"go.uber.org/zap"
	"sync"
)

type Manager struct {
	logger                *zap.SugaredLogger
	eventSource           domain.EventSource
	actionExecutorCreator domain.ActionExecutorCreator
	eventChannels         []domain.EventCh

	mutex sync.Mutex
}

var _ domain.HookManager = (*Manager)(nil)

func NewManager(eventSource domain.EventSource, executorCreator domain.ActionExecutorCreator, logger *zap.SugaredLogger) *Manager {
	return &Manager{
		logger:                logger,
		eventSource:           eventSource,
		actionExecutorCreator: executorCreator,
	}
}

func (m *Manager) RegisterHook(hook *domain.Hook, hookDef *domain.HookDefinition) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	evType := domain.UnmarshalEventType(hook.EventType)
	if evType == domain.UnknownEventType {
		return fmt.Errorf("failed to subscribe event: unknown event type: %s", hook.EventType)
	}
	eventCh, err := m.eventSource.Subscribe(evType, &hook.Parameter, nil)
	if err != nil {
		return fmt.Errorf("failed to subscribe exevent: %w", err)
	}

	executor, err := m.actionExecutorCreator.CreateExecutor(&hook.Action, hookDef)
	if err != nil {
		m.eventSource.Unsubscribe(eventCh)
		return fmt.Errorf("failed to create action executor: %w", err)
	}
	m.eventChannels = append(m.eventChannels, eventCh)
	go func() {
		for event := range eventCh {
			if executor != nil {
				executor.Handle(event)
			}
		}
		executor.Destroy()
	}()
	return nil
}

func (m *Manager) UnregisterAll() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	for _, eventCh := range m.eventChannels {
		m.eventSource.Unsubscribe(eventCh)
	}
	m.eventChannels = nil
}
