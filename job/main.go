package job

import (
	"context"
	"fmt"
	"syscall"

	"github.com/spf13/viper"
	"github.com/x-mod/cmd"
	"github.com/x-mod/routine"
)

func Main(c *cmd.Command, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("job command required")
	}
	cmdOpts := []routine.CommandOpt{}
	for index, argument := range args {
		if index >= 1 {
			cmdOpts = append(cmdOpts, routine.ARG(argument))
		}
	}
	command := routine.Executor(routine.Command(args[0], cmdOpts...))
	if viper.GetDuration("cmd-timeout") > 0 {
		command = routine.Timeout(viper.GetDuration("cmd-timeout"), command)
	}
	if viper.GetInt("retry") > 0 {
		command = routine.Retry(viper.GetInt("retry"), command)
	}
	if viper.GetInt("repeat-times") > 0 {
		command = routine.Repeat(
			viper.GetInt("repeat-times"),
			viper.GetDuration("repeat-interval"),
			command)
	}
	if viper.GetInt("concurrent") > 0 {
		command = routine.Concurrent(viper.GetInt("concurrent"), command)
	}
	if viper.GetDuration("job-timeout") > 0 {
		command = routine.Timeout(viper.GetDuration("job-timeout"), command)
	}
	if len(viper.GetString("schedule")) > 0 {
		command = routine.Crontab(viper.GetString("schedule"), command)
	}
	ctx, cancel := context.WithCancel(context.TODO())
	return routine.Main(
		ctx,
		command,
		routine.Signal(syscall.SIGINT, routine.SigHandler(func() {
			cancel()
		})),
	)
}
