package exevent

import (
	"github.com/ichenhe/syncthing-hook/domain"
	"go.uber.org/zap"
)

type ActionExecutorCreator struct {
	logger *zap.SugaredLogger
}

var _ domain.ActionExecutorCreator = (*ActionExecutorCreator)(nil)

func NewActionExecutorCreator(logger *zap.SugaredLogger) domain.ActionExecutorCreator {
	return &ActionExecutorCreator{
		logger: logger,
	}
}

// CreateExecutor creates executor handler based on given configuration.
//
// Returned error could be domain.ErrInvalidActionType.
func (c *ActionExecutorCreator) CreateExecutor(action *domain.HookAction, hookDef *domain.HookDefinition) (domain.EventHandler, error) {
	logger := hookDef.AddToLogger(c.logger)
	switch action.Type {
	case "exec":
		return NewExecExecutor(action.Cmd, logger), nil
	default:
		return nil, domain.ErrInvalidActionType
	}
}
