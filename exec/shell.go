package exec

import (
	"context"
	"io/ioutil"

	"github.com/liujianping/job/config"
	"github.com/x-mod/errors"
	"github.com/x-mod/routine"
)

//ShellCommand struct
type ShellCommand struct {
	cmd *config.Command
}

//NewShellCommand new
func NewShellCommand(cmd *config.Command) routine.Executor {
	return &ShellCommand{cmd: cmd}
}

//Execute of ShellCommand
func (sh *ShellCommand) Execute(ctx context.Context) error {
	if sh.cmd.Shell == nil {
		return errors.New("command shell required")
	}
	if len(sh.cmd.Shell.Name) == 0 {
		return errors.New("command required")
	}

	cmd := routine.Command(sh.cmd.Shell.Name, sh.cmd.Shell.Args...)
	if sh.cmd.Retry > 0 {
		cmd = routine.Retry(sh.cmd.Retry, cmd)
	}
	if sh.cmd.Timeout > 0 {
		cmd = routine.Timeout(sh.cmd.Timeout, cmd)
	}
	for _, kv := range sh.cmd.Shell.Envs {
		ctx = routine.WithEnviron(ctx, kv.Name, kv.Value)
	}
	if !sh.cmd.Stdout {
		ctx = routine.WithStdout(ctx, ioutil.Discard)
	}
	return cmd.Execute(ctx)
}
