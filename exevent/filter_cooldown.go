package exevent

import (
	"SyncthingHook/domain"
	"SyncthingHook/utils/clock"
	"go.uber.org/zap"
)

// options ----------------------------------------------------------------------------------

type CoolDownFilterOptions struct {
	timeProvider domain.TimeProvider
}

type CoolDownFilterOption func(opt *CoolDownFilterOptions)

func CoolDownFilterTimeProvider(provider domain.TimeProvider) CoolDownFilterOption {
	return func(opt *CoolDownFilterOptions) {
		opt.timeProvider = provider
	}
}

// -----------------------------------------------------------------------------------------

type CoolDownFilter struct {
	baseHandler
	*CoolDownFilterOptions
	// coolDown indicates the maximum frequency in millisecond that this event can be triggered.
	// All events during coolDown period will be discarded. coolDown is calculated based on the moment receive and
	// deal with the event, not the moment it occurred.
	coolDown int64
	logger   *zap.SugaredLogger
	// the system time when the last event was processed and sent
	lastTime int64
}

func NewCoolDownFilter(cooldown int64, logger *zap.SugaredLogger, options ...CoolDownFilterOption) *CoolDownFilter {
	_logger := wrapLogger(logger, "CoolDownFilter")
	var _cooldown int64
	if cooldown < 0 {
		_logger.Warnf("'cooldown' must >= 0, use 0 instead of '%d'", cooldown)
		_cooldown = 0
	} else {
		_cooldown = cooldown
	}

	opt := &CoolDownFilterOptions{
		timeProvider: &clock.TimeProvider{},
	}
	for _, f := range options {
		f(opt)
	}
	return &CoolDownFilter{
		CoolDownFilterOptions: opt,
		coolDown:              _cooldown,
		logger:                _logger,
		lastTime:              0,
	}
}

func (h *CoolDownFilter) Handle(event *domain.Event) {
	if h.coolDown == 0 {
		h.callNext(event)
		return
	}
	sendTime := h.timeProvider.NowUnixMilli()
	if sendTime-h.lastTime < h.coolDown {
		h.logger.Debugln("still cooling down, discard event")
		return // still cooling down
	}
	h.lastTime = sendTime
	h.callNext(event)
}
