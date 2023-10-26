package exevent

import (
	"SyncthingHook/domain"
	"SyncthingHook/mocker"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"testing"
)

func TestDefaultDirPathExtractor(t *testing.T) {
	tests := []struct {
		name              string
		data              map[string]any
		wantOk            bool
		wantExtractedPath string
	}{
		{name: "from file", data: map[string]any{"type": "file", "path": "abc/def/a.png"}, wantOk: true, wantExtractedPath: "/abc/def"},
		{name: "from dir", data: map[string]any{"type": "dir", "path": "abc/d"}, wantOk: true, wantExtractedPath: "/abc/d"},
		{name: "no path", data: map[string]any{"type": "file"}, wantOk: false, wantExtractedPath: ""},
		{name: "not a string value", data: map[string]any{"type": "file", "path": 123}, wantOk: false, wantExtractedPath: ""},
		{name: "no type", data: map[string]any{"path": "abc/d"}, wantOk: false, wantExtractedPath: ""},
		{name: "not a valid type", data: map[string]any{"type": "?", "path": "abc/d"}, wantOk: false, wantExtractedPath: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ev := &domain.Event{Data: tt.data}
			extracted, err := defaultDirPathExtractor(ev)
			if tt.wantOk {
				require.Nil(t, err)
				require.EqualValues(t, tt.wantExtractedPath, extracted)
			} else {
				require.NotNil(t, err)
			}
		})
	}
}

func TestDirPathFilter_Handle(t *testing.T) {
	tests := []struct {
		name         string
		path         string
		pattern      string
		wantCallNext bool
	}{
		{name: "match path", path: "abc/def", pattern: "/abc", wantCallNext: true},
		{name: "mismatch path", path: "abc", pattern: "/abc/def", wantCallNext: false},
		{name: "no path", path: "", pattern: "/", wantCallNext: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := map[string]any{"type": "dir", "path": tt.path}
			if tt.path == "" {
				delete(data, "path")
			}

			ev := &domain.Event{Data: data}
			next := &mocker.MockEventHandler{}
			next.EXPECT().Handle(ev).Return()

			filter := NewDirPathFilter(tt.pattern, zap.NewNop().Sugar())
			filter.SetNext(next)
			filter.Handle(ev)

			if tt.wantCallNext {
				next.AssertNextHandlerCalled(t, ev, 1)
			} else {
				next.AssertNextHandlerNotCalled(t, ev)
			}
		})
	}
}

func TestDirPathFilter_matchPath(t *testing.T) {
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
			filter := NewDirPathFilter(tt.pattern, zap.NewNop().Sugar())
			require.EqualValues(t, tt.match, filter.matchPath(tt.path))
		})
	}
}
