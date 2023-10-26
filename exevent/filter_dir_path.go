package exevent

import (
	"errors"
	"fmt"
	"github.com/ichenhe/syncthing-hook/domain"
	"go.uber.org/zap"
	"path/filepath"
	"strings"
)

// options ----------------------------------------------------------------------------------

type DirPathFilterOptions struct {
	Logger   *zap.SugaredLogger
	BasePath string
	// extract path from the event, must start with '/'
	DirPathExtractor func(event *domain.Event) (string, error)
}

var defaultDirPathExtractor = func(event *domain.Event) (string, error) {
	data, ok := event.Data.(map[string]any)
	if !ok {
		return "", errors.New("unexpected 'Data' type")
	}
	folder, ex := data["path"]
	if !ex {
		return "", errors.New("key 'path' does not exist")
	}
	s, ok := folder.(string)
	if !ok {
		return "", errors.New("unexpected 'path' type")
	}
	if !strings.HasPrefix(s, "/") {
		s = "/" + s
	}
	if data["type"] == "file" {
		return filepath.Dir(s), nil
	} else if data["type"] == "dir" {
		return s, nil
	} else {
		return "", fmt.Errorf("unknown 'type': %s", data["type"])
	}
}

type DirPathFilterOption func(opt *DirPathFilterOptions)

func DirPathFilterWithPathExtractor(pathExtractor func(event *domain.Event) (string, error)) DirPathFilterOption {
	return func(opt *DirPathFilterOptions) {
		opt.DirPathExtractor = pathExtractor
	}
}

// -----------------------------------------------------------------------------------------

// DirPathFilter terminates all events except that the path matches expected base path (equals or subdir).
// Unable to retrieve the path or type from the event is considered a match failure.
type DirPathFilter struct {
	baseHandler
	*DirPathFilterOptions
}

var _ domain.EventHandler = (*DirPathFilter)(nil)

func NewDirPathFilter(basePath string, logger *zap.SugaredLogger, options ...DirPathFilterOption) *DirPathFilter {
	_logger := wrapLogger(logger, "CoolDownFilter")
	opt := &DirPathFilterOptions{
		Logger:           _logger,
		BasePath:         basePath,
		DirPathExtractor: defaultDirPathExtractor,
	}
	for _, f := range options {
		f(opt)
	}
	return &DirPathFilter{
		DirPathFilterOptions: opt,
	}
}

func (d *DirPathFilter) Handle(event *domain.Event) {
	path, err := d.DirPathExtractor(event)
	if err != nil {
		d.Logger.Warnf("filed to extract path from exevent: %s", err)
		return
	}
	if d.matchPath(path) {
		d.callNext(event)
	} else {
		d.Logger.Debugw("terminate event due to dirPath mismatching", zap.String("pattern", d.BasePath), zap.String("got", path))
	}
}

// matchPath tests whether given eventPath matches the pattern.
// eventPath must start with '/' and represents a folder.
func (d *DirPathFilter) matchPath(eventPath string) bool {
	rel, err := filepath.Rel(d.BasePath, eventPath)
	return err == nil && !strings.HasPrefix(rel, "..")
}
