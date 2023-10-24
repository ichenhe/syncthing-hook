package extevent

import (
	"github.com/stretchr/testify/assert"
	"github.com/syncthing/syncthing/lib/events"
	"testing"
)

func TestMatchPath(t *testing.T) {
	newDetectorCore := func(pathPattern string) *LocalFolderContentChangedDetector {
		return &LocalFolderContentChangedDetector{
			st:     nil,
			option: LocalFolderContentChangedDetectorOptions{Path: pathPattern},
		}
	}

	tests := []struct {
		name    string
		pattern string
		path    string
		match   bool
	}{
		{name: "match_sameDir", pattern: "/test", path: "/test", match: true},
		{name: "match_subDir", pattern: "/test", path: "/test/sub", match: true},
		{name: "match_root", pattern: "/", path: "/", match: true},
		{name: "mismatch_parentDir", pattern: "/d1/d2", path: "/d1", match: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			core := newDetectorCore(tt.pattern)
			assert.EqualValues(t, tt.match, core.matchPath(tt.path))
		})
	}
}

func TestMatchEvent(t *testing.T) {
	newDetectorCore := func(modOpt func(option *LocalFolderContentChangedDetectorOptions)) *LocalFolderContentChangedDetector {
		opt := &LocalFolderContentChangedDetectorOptions{}
		modOpt(opt)
		return &LocalFolderContentChangedDetector{
			st:     nil,
			option: *opt,
		}
	}
	newEvent := func(folder string, t string, path string) *events.Event {
		return &events.Event{Data: map[string]any{"folder": folder, "type": t, "path": path}}
	}

	tests := []struct {
		name          string
		matchFolderId string
		matchPath     string
		event         *events.Event
		match         bool
	}{
		{name: "match file", matchFolderId: "folder-id", matchPath: "/abc", event: newEvent("folder-id", "file", "abc/def/a.png"), match: true},
		{name: "match dir", matchFolderId: "folder-id", matchPath: "/abc", event: newEvent("folder-id", "dir", "abc/def"), match: true},
		{name: "mismatch folder id", matchFolderId: "folder-id", matchPath: "/", event: newEvent("xxx", "file", "t"), match: false},
		{name: "mismatch path", matchFolderId: "folder-id", matchPath: "/sub", event: newEvent("folder-id", "file", "/sub2/a.png"), match: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			core := newDetectorCore(func(option *LocalFolderContentChangedDetectorOptions) {
				option.StFolderId = tt.matchFolderId
				option.Path = tt.matchPath
			})
			assert.EqualValues(t, tt.match, core.matchEvent(tt.event))
		})
	}
}
