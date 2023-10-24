package eventhandler

import (
	"SyncthingHook/extevent"
	"go.uber.org/zap"
	"time"
)

type ToleranceFilter struct {
	baseHandler
	// tolerance indicates how many milliseconds to wait before triggering a new event.
	// Usually there will be multiple file changes under one folder in a short period of time, and we do not want to
	// trigger this event frequently. Tolerance is calculated based on the moment an upstream event is received, not
	// the moment it occurs.
	tolerance uint64
	logger    *zap.SugaredLogger
	timer     *time.Timer
}

func NewToleranceFilter(tolerance int64, logger *zap.SugaredLogger) *ToleranceFilter {
	_logger := wrapLogger(logger, "ToleranceFilter")
	var _tolerance uint64
	if tolerance < 0 {
		_logger.Warnf("'tolerance' must >= 0, use 0 instead of %d", tolerance)
		_tolerance = 0
	} else {
		_tolerance = uint64(tolerance)
	}
	return &ToleranceFilter{
		tolerance: _tolerance,
		logger:    _logger,
	}
}

func (h *ToleranceFilter) Handle(event *extevent.Event) {
	if h.tolerance == 0 {
		h.callNext(event)
		return
	}
	h.stopTimerIfNeeded()
	h.timer = time.AfterFunc(time.Millisecond*time.Duration(h.tolerance), func() {
		h.callNext(event)
	})
}

func (h *ToleranceFilter) Destroy() {
	// No downstream subscriber anymore, let's cancel the pending event.
	// However, it may result in  the subsequent subscribers not receive the event which triggered by previous
	// st event and should be dispatched to downstream at the moment.
	// If we don't cancel the timer here, there is a risk of routine leakage if Tolerance is too long.
	h.stopTimerIfNeeded()
	h.timer = nil
}

func (h *ToleranceFilter) stopTimerIfNeeded() {
	if h.timer != nil && !h.timer.Stop() {
		select {
		case <-h.timer.C:
		default:
		}
	}
}
