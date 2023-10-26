package exevent

import (
	"SyncthingHook/domain"
	"errors"
	"go.uber.org/zap"
)

// options ----------------------------------------------------------------------------------

type FolderIdFilterOptions struct {
	Logger            *zap.SugaredLogger
	FolderId          string
	FolderIdExtractor func(event *domain.Event) (string, error)
}

var defaultFolderIdExtractor = func(event *domain.Event) (string, error) {
	data, ok := event.Data.(map[string]any)
	if !ok {
		return "", errors.New("unexpected 'Data' type")
	}
	folder, ex := data["folder"]
	if !ex {
		return "", errors.New("key 'folder' does not exist")
	}
	if s, ok := folder.(string); !ok {
		return "", errors.New("unexpected 'folder' type")
	} else {
		return s, nil
	}
}

type FolderIdFilterOption func(opt *FolderIdFilterOptions)

func FolderIdFilterWithFolderIdExtractor(folderIdExtractor func(event *domain.Event) (string, error)) FolderIdFilterOption {
	return func(opt *FolderIdFilterOptions) {
		opt.FolderIdExtractor = folderIdExtractor
	}
}

// -----------------------------------------------------------------------------------------

// FolderIdFilter terminates all events except that the folder ID is the same as expected.
// Unable to retrieve the folderId from the event is considered a match failure.
type FolderIdFilter struct {
	baseHandler
	*FolderIdFilterOptions
}

var _ domain.EventHandler = (*FolderIdFilter)(nil)

func NewFolderIdFilter(folderId string, logger *zap.SugaredLogger, options ...FolderIdFilterOption) *FolderIdFilter {
	_logger := wrapLogger(logger, "CoolDownFilter")
	opt := &FolderIdFilterOptions{
		Logger:            _logger,
		FolderId:          folderId,
		FolderIdExtractor: defaultFolderIdExtractor,
	}
	for _, f := range options {
		f(opt)
	}
	return &FolderIdFilter{
		FolderIdFilterOptions: opt,
	}
}

func (f *FolderIdFilter) Handle(event *domain.Event) {
	folderId, err := f.FolderIdExtractor(event)
	if err != nil {
		f.Logger.Warnf("filed to extract folderId from exevent: %s", err)
		return
	}
	if folderId == f.FolderId {
		f.callNext(event)
	} else {
		f.Logger.Debugw("terminate event due to folderId mismatching", zap.String("expect", f.FolderId), zap.String("got", folderId))
	}
}
