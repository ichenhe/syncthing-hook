package eventhandler

import (
	"SyncthingHook/extevent"
	"errors"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

func setupMockExecutor(stupid bool) *MockcmdExecutor {
	exec := &MockcmdExecutor{}
	var ret error
	if stupid {
		ret = errors.New("for testing")
	}
	exec.On("exec", mock.AnythingOfType("string"), mock.Anything).Return(ret)
	return exec
}

func TestExecExecutor_Handle(t *testing.T) {
	tests := []struct {
		name        string
		cmd         []string
		wantIgnore  bool
		wantExecErr bool
	}{
		{name: "empty cmds", cmd: []string{}, wantIgnore: true},
		{name: "the first cmd is empty", cmd: []string{"", "ls"}, wantIgnore: true},
		{name: "exec successfully", cmd: []string{"ls"}, wantIgnore: false, wantExecErr: false},
		{name: "exec failed", cmd: []string{"xx"}, wantIgnore: false, wantExecErr: true},
	}
	ev := &extevent.Event{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmdExec := setupMockExecutor(tt.wantExecErr)
			logger, logObserver := newLogger()
			next := createMockNextHandler(ev)
			h := NewExecExecutor(tt.cmd, logger, ExecExecutorCmdExecutor(cmdExec))
			h.SetNext(next)
			h.Handle(ev)

			// build args
			_cmd := make([]any, len(tt.cmd))
			for i := 0; i < len(tt.cmd); i++ {
				_cmd[i] = tt.cmd[i]
			}
			// should call next
			next.assertNextHandlerCalled(t, ev, 1)

			if tt.wantIgnore {
				cmdExec.AssertNotCalled(t, "exec", _cmd...)
				return
			} else {
				cmdExec.AssertCalled(t, "exec", _cmd...)
			}
			require.EqualValues(t, 1, logObserver.Len(), "should only has 1 log")
			if tt.wantExecErr {
				require.Contains(t, logObserver.TakeAll()[0].Message, "failed")
			} else {
				require.Contains(t, logObserver.TakeAll()[0].Message, "successfully")
			}
		})
	}
}

func TestExecExecutor_execCmd(t *testing.T) {
	tests := []struct {
		name       string
		cmd        []string
		wantErr    bool
		calledArgs []interface{}
	}{
		{name: "len(cmd)=1", cmd: []string{"ls"}, wantErr: false},
		{name: "len(cmd)=2", cmd: []string{"ls", "-al"}, wantErr: false},
		{name: "exec error", cmd: []string{"ls", "-al"}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger, _ := newLogger()
			mockExecutor := setupMockExecutor(tt.wantErr)
			exec := NewExecExecutor(tt.cmd, logger, ExecExecutorCmdExecutor(mockExecutor))
			res := exec.execCmd()
			require.Equal(t, tt.wantErr, res != nil, "execCmd() error = %v, wantErr %v", res, tt.wantErr)

			_cmd := make([]any, len(tt.cmd))
			for i := 0; i < len(tt.cmd); i++ {
				_cmd[i] = tt.cmd[i]
			}
			mockExecutor.AssertCalled(t, "exec", _cmd...)
		})
	}
}
