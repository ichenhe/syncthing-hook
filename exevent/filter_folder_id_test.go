package exevent

import (
	"github.com/ichenhe/syncthing-hook/domain"
	"github.com/ichenhe/syncthing-hook/mocker"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"testing"
)

func TestDefaultFolderIdExtractor(t *testing.T) {
	tests := []struct {
		name                  string
		data                  map[string]any
		wantOk                bool
		wantExtractedFolderId string
	}{
		{name: "success", data: map[string]any{"folder": "id"}, wantOk: true, wantExtractedFolderId: "id"},
		{name: "no folder", data: map[string]any{}, wantOk: false, wantExtractedFolderId: ""},
		{name: "not a string value", data: map[string]any{"folder": 123}, wantOk: false, wantExtractedFolderId: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ev := &domain.Event{Data: tt.data}
			extracted, err := defaultFolderIdExtractor(ev)
			if tt.wantOk {
				require.Nil(t, err)
				require.EqualValues(t, tt.wantExtractedFolderId, extracted)
			} else {
				require.NotNil(t, err)
			}
		})
	}
}

func TestFolderIdFilter_Handle(t *testing.T) {
	tests := []struct {
		name         string
		folderId     string
		pattern      string
		wantCallNext bool
	}{
		{name: "match folder id", folderId: "folder-id", pattern: "folder-id", wantCallNext: true},
		{name: "mismatch folder id", folderId: "folder-id", pattern: "folder_id", wantCallNext: false},
		{name: "no folder id", folderId: "", pattern: "folder-id", wantCallNext: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := map[string]any{"folder": tt.folderId}
			if tt.folderId == "" {
				delete(data, "folder")
			}
			ev := &domain.Event{Data: data}
			next := &mocker.MockEventHandler{}
			next.EXPECT().Handle(ev).Return()

			filter := NewFolderIdFilter(tt.pattern, zap.NewNop().Sugar())
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
