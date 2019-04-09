package exec

import (
	"context"
	"strings"
	"time"

	"github.com/x-mod/errors"

	"github.com/x-mod/routine"

	"github.com/liujianping/job/pb"
	"github.com/satori/go.uuid"
)

//Job struct
type Job struct {
	opts *options
	cmd  string
}

//NewJob new job
func NewJob(opts ...Option) *Job {
	defaults := options{
		job: &pb.Job{
			Uuid: uuid.NewV4().String(),
		},
	}
	for _, opt := range opts {
		opt(&defaults)
	}

	return &Job{opts: &defaults}
}

func (j *Job) String() string {
	return j.cmd
}

//Execute implement executor
func (j *Job) Execute(ctx context.Context) error {
	jd := j.opts.job
	if len(jd.Command) == 0 {
		return errors.New("cmd required")
	}
	cmds := []string{jd.Command}
	cmds = append(cmds, jd.Argments...)
	j.cmd = strings.Join(cmds, " ")
	exec := routine.Command(jd.Command, jd.Argments...)
	if j.opts.cmdTimeout > 0 {
		exec = routine.Timeout(j.opts.cmdTimeout, exec)
	}
	if jd.RetryTimes > 0 {
		exec = routine.Retry(int(jd.RetryTimes), exec)
	}
	if j.opts.report != nil {
		exec = routine.Report(j.opts.report, exec)
	}
	if j.opts.guarantee {
		exec = routine.Guarantee(exec)
	}
	if j.opts.concurrent > 0 {
		exec = routine.Concurrent(j.opts.concurrent, exec)
	}
	if jd.RepeatTimes != 1 {
		d, err := time.ParseDuration(jd.RepeatInterval)
		if err != nil {
			return err
		}
		exec = routine.Repeat(int(jd.RepeatTimes), d, exec)
	}
	//* * * * *
	if len(jd.Crontab) > 8 {
		exec = routine.Crontab(jd.Crontab, exec)
	}
	if j.opts.jobTimeout > 0 {
		exec = routine.Deadline(time.Now().Add(j.opts.jobTimeout), exec)
	}
	return exec.Execute(ctx)
}
