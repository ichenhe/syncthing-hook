package extevent

import (
	"SyncthingHook/stclient"
	"github.com/syncthing/syncthing/lib/events"
	"path/filepath"
	"strings"
)

type LocalFolderContentChangedDetectorOptions struct {
	StFolderId string
	Path       string
	Since      int
}

type LocalFolderContentChangedDetector struct {
	st     *stclient.Syncthing
	option LocalFolderContentChangedDetectorOptions

	lastDispatchTime int64 // use atomic instead of mutex
}

func NewLocalFolderContentChangedDetector(syncthing *stclient.Syncthing, stFolderId string, modOption func(option *LocalFolderContentChangedDetectorOptions)) (EventDetector, error) {
	option := LocalFolderContentChangedDetectorOptions{
		StFolderId: stFolderId,
		Path:       "/",
		Since:      0,
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

func checkLocalFolderContentChangedDetectorOption(opt *LocalFolderContentChangedDetectorOptions) *IllegalArgumentError {
	if len(opt.StFolderId) == 0 {
		return newIllegalArgumentError("'StFolderId' can not be empty")
	} else if len(opt.Path) == 0 {
		return newIllegalArgumentError("'Path' can not be empty")
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
	for e := range upstreamCh {
		if !d.matchEvent(&e) {
			continue
		}

		newEvent := NewEvent(LocalFolderContentChangeDetected, e.Time, e.Data)
		dispatcher(newEvent)
	}
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

// matchPath tests whether given eventPath matches the pattern specified by LocalFolderContentChangedDetectorOptions.Path.
// eventPath must start with '/' and represents a folder.
func (d *LocalFolderContentChangedDetector) matchPath(eventPath string) bool {
	rel, err := filepath.Rel(d.option.Path, eventPath)
	return err == nil && !strings.HasPrefix(rel, "..")
}
