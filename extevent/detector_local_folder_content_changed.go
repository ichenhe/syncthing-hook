package extevent

import (
	"SyncthingHook/stclient"
	"github.com/syncthing/syncthing/lib/events"
	"path/filepath"
	"strings"
	"sync/atomic"
	"time"
)

type LocalFolderContentChangedDetectorOption struct {
	StFolderId string
	Path       string
	Since      int
	// Tolerance indicates how many milliseconds to wait before triggering a new event.
	// Usually there will be multiple file changes under one folder in a short period of time, and we do not want to
	// trigger this event frequently. Tolerance is calculated based on the moment an upstream event is received, not
	// the moment it occurs.
	Tolerance int64
	// Cooldown indicates the maximum frequency in millisecond that this event can be triggered. All triggering during
	// Cooldown period will be discarded. Cooldown is calculated based on the moment dispatch event to downstream, not
	// the moment it occurs.
	Cooldown int64
}

type LocalFolderContentChangedDetector struct {
	st     *stclient.Syncthing
	option LocalFolderContentChangedDetectorOption

	lastDispatchTime int64 // use atomic instead of mutex
}

func NewLocalFolderContentChangedDetector(syncthing *stclient.Syncthing, stFolderId string, modOption func(option *LocalFolderContentChangedDetectorOption)) (EventDetector, error) {
	option := LocalFolderContentChangedDetectorOption{
		StFolderId: stFolderId,
		Path:       "/",
		Since:      0,
		Tolerance:  1000,
		Cooldown:   500,
	}
	if modOption != nil {
		modOption(&option)
	}
	if err := checkLocalFolderContentChangedDetectorOption(&option); err != nil {
		return nil, err
	}
	return newDetector(syncthing, &LocalFolderContentChangedDetector{
		st:     syncthing,
		option: option,
	}), nil
}

func checkLocalFolderContentChangedDetectorOption(opt *LocalFolderContentChangedDetectorOption) *IllegalArgumentError {
	if len(opt.StFolderId) == 0 {
		return newIllegalArgumentError("'StFolderId' can not be empty")
	} else if len(opt.Path) == 0 {
		return newIllegalArgumentError("'Path' can not be empty")
	} else if opt.Tolerance < 0 {
		return newIllegalArgumentError("'Tolerance' must >0")
	} else if opt.Cooldown < 0 {
		return newIllegalArgumentError("'Cooldown' must >0")
	}
	return nil
}

func (d *LocalFolderContentChangedDetector) unsubscribeUpstream(st *stclient.Syncthing, upstreamCh upstream) {
	st.UnsubscribeEvent(upstreamCh)
}

func (d *LocalFolderContentChangedDetector) subscribeUpstream(st *stclient.Syncthing) upstream {
	// register upstream subscription
	return st.SubscribeEvent(eventTypes(events.LocalChangeDetected), d.option.Since)
}

func (d *LocalFolderContentChangedDetector) handleUpstream(upstreamCh upstream, dispatcher downstreamDispatcher) {
	var timer *time.Timer
	for e := range upstreamCh {
		if !d.matchEvent(&e) {
			continue
		}
		timer = d.handleUpstreamEvent(&e, timer, dispatcher)
	}
	// No downstream subscriber anymore, let's cancel the pending event.
	// However, it may result in  the subsequent subscribers not receive the event which triggered by previous
	// st event and should be dispatched to downstream at the moment.
	// If we don't cancel the timer here, there is a risk of routine leakage if Tolerance is too long.
	stopTimerIfNeeded(timer)
}

func (d *LocalFolderContentChangedDetector) matchEvent(e *events.Event) bool {
	data := e.Data.(map[string]any)
	if d.option.StFolderId != data["folder"] {
		return false
	}
	if data["type"] == "file" {
		return d.matchPath(filepath.Dir("/" + data["path"].(string)))
	} else if data["type"] == "dir" {
		return d.matchPath("/" + data["path"].(string))
	} else {
		return false
	}
}

// matchPath tests whether given eventPath matches the pattern specified by LocalFolderContentChangedDetectorOption.Path.
// eventPath must start with '/' and represents a folder.
func (d *LocalFolderContentChangedDetector) matchPath(eventPath string) bool {
	rel, err := filepath.Rel(d.option.Path, eventPath)
	return err == nil && !strings.HasPrefix(rel, "..")
}

// handleUpstreamEvent deals with upstream event. timer is used to cancel the previous pending event.
func (d *LocalFolderContentChangedDetector) handleUpstreamEvent(e *events.Event, timer *time.Timer, dispatcher downstreamDispatcher) *time.Timer {
	newEvent := NewEvent(LocalFolderContentChangeDetected, e.Time, e.Data)
	if d.option.Tolerance == 0 {
		d.dispatchEventIfNotCoolingDown(newEvent, dispatcher)
		return timer
	} else {
		if !d.checkCoolingDown(time.Now().UnixMilli() + d.option.Tolerance) {
			return timer // will be still cooling down even after Tolerance waiting, ignore this event
		}
		// place a new pending dispatch
		stopTimerIfNeeded(timer)
		return time.AfterFunc(time.Millisecond*time.Duration(d.option.Tolerance), func() {
			d.dispatchEventIfNotCoolingDown(newEvent, dispatcher)
		})
	}
}

// checkCoolingDown returns true if it is free to dispatch at expectedDispatchTime.
func (d *LocalFolderContentChangedDetector) checkCoolingDown(expectedDispatchTime int64) bool {
	return d.option.Cooldown == 0 || expectedDispatchTime-atomic.LoadInt64(&d.lastDispatchTime) > d.option.Cooldown
}

func (d *LocalFolderContentChangedDetector) dispatchEventIfNotCoolingDown(newEvent *Event, dispatcher downstreamDispatcher) {
	sendTime := time.Now().UnixMilli()
	if d.option.Cooldown != 0 {
		for {
			lastTime := atomic.LoadInt64(&d.lastDispatchTime)
			if sendTime-lastTime < d.option.Cooldown {
				return // still cooling down
			}
			if atomic.CompareAndSwapInt64(&d.lastDispatchTime, lastTime, sendTime) {
				break
			}
		}
	}
	dispatcher(newEvent)
}
