package exec

import (
	"time"

	"github.com/liujianping/job/pb"
	"github.com/x-mod/routine"
)

//Option for JOB
type Option func(*options)

type options struct {
	//options for runner
	report     chan *routine.Result
	pipe       bool
	guarantee  bool
	concurrent int
	jobTimeout time.Duration
	cmdTimeout time.Duration
	//options for job
	job *pb.Job
}

func PipeLineMode(pipe bool) Option {
	return func(opts *options) {
		opts.pipe = pipe
	}
}

func Guarantee(g bool) Option {
	return func(opts *options) {
		opts.guarantee = g
	}
}

func Report(report chan *routine.Result) Option {
	return func(opts *options) {
		opts.report = report
	}
}

func Concurrent(c int) Option {
	return func(opts *options) {
		opts.concurrent = c
	}
}

//JobName Optioons
func JobName(name string) Option {
	return func(opts *options) {
		opts.job.Name = name
	}
}
func JobCommand(command string) Option {
	return func(opts *options) {
		opts.job.Command = command
	}
}
func JobArgs(args ...string) Option {
	return func(opts *options) {
		opts.job.Argments = args
	}
}
func JobEnv(key, value string) Option {
	return func(opts *options) {
		opts.job.Envs[key] = value
	}
}

func JobPayload(payload []byte) Option {
	return func(opts *options) {
		opts.job.Payload = payload
	}
}

func JobRetryTimes(retry int32) Option {
	return func(opts *options) {
		opts.job.RetryTimes = retry
	}
}

func JobRepeatTimes(retry int32) Option {
	return func(opts *options) {
		opts.job.RepeatTimes = retry
	}
}
func JobRepeatInterval(d time.Duration) Option {
	return func(opts *options) {
		opts.job.RepeatInterval = d.String()
	}
}

func JobCrontab(plan string) Option {
	return func(opts *options) {
		opts.job.Crontab = plan
	}
}

func CmdTimeout(duration time.Duration) Option {
	return func(opts *options) {
		opts.cmdTimeout = duration
	}
}

func JobTimeout(duration time.Duration) Option {
	return func(opts *options) {
		opts.jobTimeout = duration
	}
}
