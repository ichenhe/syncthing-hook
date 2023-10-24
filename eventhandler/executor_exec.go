package eventhandler

import (
	"SyncthingHook/extevent"
	"go.uber.org/zap"
	"os/exec"
)

// options ----------------------------------------------------------------------------------

type ExecExecutorOptions struct {
	cmdExec cmdExecutor
}

type ExecExecutorOption func(opt *ExecExecutorOptions)

func ExecExecutorCmdExecutor(cmdExec cmdExecutor) ExecExecutorOption {
	return func(opt *ExecExecutorOptions) {
		opt.cmdExec = cmdExec
	}
}

// -----------------------------------------------------------------------------------------

type cmdExecutor interface {
	exec(name string, arg ...string) error
}

type cmdExecutorFunc func(name string, arg ...string) error

func (f cmdExecutorFunc) exec(name string, arg ...string) error {
	return f(name, arg...)
}

type ExecExecutor struct {
	baseHandler
	*ExecExecutorOptions
	cmd    []string
	logger *zap.SugaredLogger
}

func NewExecExecutor(cmd []string, logger *zap.SugaredLogger, options ...ExecExecutorOption) *ExecExecutor {
	opt := &ExecExecutorOptions{
		cmdExec: cmdExecutorFunc(func(name string, arg ...string) error {
			return exec.Command(name, arg...).Start()
		}),
	}
	for _, f := range options {
		f(opt)
	}

	return &ExecExecutor{
		ExecExecutorOptions: opt,
		cmd:                 cmd,
		logger:              wrapLogger(logger, "ExecExecutor"),
	}
}

func (h *ExecExecutor) Handle(event *extevent.Event) {
	if len(h.cmd) == 0 || len(h.cmd[0]) == 0 {
		h.logger.Debugw("cmd is empty, ignore this exec action.")
		h.callNext(event)
		return
	}
	if err := h.execCmd(); err != nil {
		h.logger.Infow("failed to execute action: "+err.Error(), "cmd", h.cmd)
	} else {
		h.logger.Debugw("execute cmd successfully", "cmd", h.cmd)
	}
	h.callNext(event)
}

func (h *ExecExecutor) execCmd() error {
	return h.cmdExec.exec(h.cmd[0], h.cmd[1:]...)
}
