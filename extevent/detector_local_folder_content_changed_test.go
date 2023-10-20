package extevent

import (
	"github.com/stretchr/testify/assert"
	"github.com/syncthing/syncthing/lib/events"
	"testing"
	"time"
)

func TestMatchPath(t *testing.T) {
	newDetectorCore := func(pathPattern string) *LocalFolderContentChangedDetector {
		return &LocalFolderContentChangedDetector{
			st:     nil,
			option: LocalFolderContentChangedDetectorOption{Path: pathPattern},
		}
	}

	t.Run("match_sameDir", func(t *testing.T) {
		core := newDetectorCore("/test")
		assert.True(t, core.matchPath("/test"))
	})

	t.Run("match_subDir", func(t *testing.T) {
		core := newDetectorCore("/test")
		assert.True(t, core.matchPath("/test/sub"))
	})

	t.Run("match_root", func(t *testing.T) {
		core := newDetectorCore("/")
		assert.True(t, core.matchPath("/"))
	})

	t.Run("unmatch_parentDir", func(t *testing.T) {
		core := newDetectorCore("/d1/d2")
		assert.False(t, core.matchPath("/d1"))
	})
}

func TestMatchEvent(t *testing.T) {
	newDetectorCore := func(modOpt func(option *LocalFolderContentChangedDetectorOption)) *LocalFolderContentChangedDetector {
		opt := &LocalFolderContentChangedDetectorOption{}
		modOpt(opt)
		return &LocalFolderContentChangedDetector{
			st:     nil,
			option: *opt,
		}
	}
	newEvent := func(folder string, t string, path string) *events.Event {
		return &events.Event{Data: map[string]any{"folder": folder, "type": t, "path": path}}
	}

	t.Run("matchFile", func(t *testing.T) {
		core := newDetectorCore(func(option *LocalFolderContentChangedDetectorOption) {
			option.StFolderId = "folder-id"
			option.Path = "/abc"
		})
		assert.True(t, core.matchEvent(newEvent("folder-id", "file", "abc/def/a.png")))
	})

	t.Run("matchDir", func(t *testing.T) {
		core := newDetectorCore(func(option *LocalFolderContentChangedDetectorOption) {
			option.StFolderId = "folder-id"
			option.Path = "/abc"
		})
		assert.True(t, core.matchEvent(newEvent("folder-id", "dir", "abc/def")))
	})

	t.Run("folderIdNotMatch", func(t *testing.T) {
		core := newDetectorCore(func(option *LocalFolderContentChangedDetectorOption) {
			option.StFolderId = "folder-id"
			option.Path = "/"
		})
		assert.False(t, core.matchEvent(newEvent("xxx", "file", "t")))
		assert.True(t, core.matchPath("/t"))
	})

	t.Run("pathNotMatch", func(t *testing.T) {
		core := newDetectorCore(func(option *LocalFolderContentChangedDetectorOption) {
			option.StFolderId = "folder-id"
			option.Path = "/sub"
		})
		assert.False(t, core.matchEvent(newEvent("folder-id", "file", "/sub2/a.png")))
	})
}

func TestTolerance(t *testing.T) {
	newDetectorCore := func(tolerance int64) *LocalFolderContentChangedDetector {
		return &LocalFolderContentChangedDetector{
			st:     nil,
			option: LocalFolderContentChangedDetectorOption{Path: "/", StFolderId: "f", Tolerance: tolerance, Cooldown: 0},
		}
	}

	t.Run("noTolerance", func(t *testing.T) {
		const NUM = 50
		var count int
		dispatcher := func(event *Event) {
			count++
		}
		for i := 0; i < NUM; i++ {
			newDetectorCore(0).handleUpstreamEvent(&events.Event{}, nil, dispatcher)
		}
		assert.Equal(t, NUM, count)
	})

	t.Run("hasTolerance", func(t *testing.T) {
		results := make([]int, 0)
		dispatcher := func(event *Event) {
			results = append(results, event.Data.(int))
		}
		core := newDetectorCore(50)

		timer := core.handleUpstreamEvent(&events.Event{Data: 1}, nil, dispatcher) // will be overridden
		timer = core.handleUpstreamEvent(&events.Event{Data: 2}, timer, dispatcher)
		time.Sleep(time.Millisecond * time.Duration(60))
		timer = core.handleUpstreamEvent(&events.Event{Data: 3}, timer, dispatcher)
		time.Sleep(time.Millisecond * time.Duration(60))

		assert.Len(t, results, 2)
		assert.ElementsMatch(t, results, []int{2, 3})
	})
}

func TestCooldown(t *testing.T) {
	newDetectorCore := func(cooldown int64) *LocalFolderContentChangedDetector {
		return &LocalFolderContentChangedDetector{
			st:     nil,
			option: LocalFolderContentChangedDetectorOption{Path: "/", StFolderId: "f", Tolerance: 0, Cooldown: cooldown},
		}
	}

	t.Run("noCooldown", func(t *testing.T) {
		const NUM = 10
		var count int
		dispatcher := func(event *Event) {
			count++
		}
		var timer *time.Timer
		for i := 0; i < NUM; i++ {
			timer = newDetectorCore(0).handleUpstreamEvent(&events.Event{}, timer, dispatcher)
		}
		assert.Equal(t, NUM, count)
	})

	t.Run("hasCooldown", func(t *testing.T) {
		results := make([]int, 0)
		dispatcher := func(event *Event) {
			results = append(results, event.Data.(int))
		}
		core := newDetectorCore(50)
		timer := core.handleUpstreamEvent(&events.Event{Data: 1}, nil, dispatcher)
		timer = core.handleUpstreamEvent(&events.Event{Data: 2}, timer, dispatcher) // will be discarded
		time.Sleep(time.Millisecond * time.Duration(60))
		timer = core.handleUpstreamEvent(&events.Event{Data: 3}, timer, dispatcher)
		time.Sleep(time.Millisecond * time.Duration(60))

		assert.Len(t, results, 2)
		assert.ElementsMatch(t, results, []int{1, 3})
	})
}
