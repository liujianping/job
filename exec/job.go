package exec

import (
	"context"
	"log"
	"time"

	"github.com/liujianping/job/config"
	"github.com/satori/go.uuid"
	"github.com/x-mod/errors"
	"github.com/x-mod/routine"
)

//Job struct instance of job
type Job struct {
	id       string
	jd       *config.JD
	reporter *Reporter
}

//NewJob new job
func NewJob(jd *config.JD, rpt *Reporter) *Job {
	return &Job{
		id:       uuid.NewV4().String(),
		jd:       jd,
		reporter: rpt,
	}
}

//ID string
func (j *Job) ID() string {
	return j.id
}

//String
func (j *Job) String() string {
	return j.jd.String()
}

//Execute implement executor
func (j *Job) Execute(ctx context.Context) error {
	log.Println("Job Executing ... ", j.String())
	var exec routine.Executor
	if j.jd.Command != nil {
		if len(j.jd.Command.Name) == 0 {
			return errors.New("cmd name required")
		}
		exec = routine.ExecutorFunc(func(ctx context.Context) error {
			cmd := routine.Command(j.jd.Command.Name, j.jd.Command.Args...)
			if j.jd.Command.Retry > 0 {
				cmd = routine.Retry(j.jd.Command.Retry, cmd)
			}
			if j.jd.Command.Timeout > 0 {
				cmd = routine.Timeout(j.jd.Command.Timeout, cmd)
			}
			for _, kv := range j.jd.Command.Envs {
				ctx = routine.WithEnviron(ctx, kv.Name, kv.Value)
			}
			return cmd.Execute(ctx)
		})
	}
	if j.jd.HTTP != nil {
		exec = routine.ExecutorFunc(func(ctx context.Context) error {
			return nil
		})
	}

	if j.reporter != nil {
		if j.jd.Report {
			exec = routine.Report(j.reporter.Report(), exec)
		}
	}
	if j.jd.Guarantee {
		exec = routine.Guarantee(exec)
	}
	if j.jd.Concurrent > 0 {
		exec = routine.Concurrent(j.jd.Concurrent, exec)
	}
	if j.jd.Repeat.Times != 1 {
		exec = routine.Repeat(j.jd.Repeat.Times, j.jd.Repeat.Interval, exec)
	}
	//* * * * *
	if len(j.jd.Crontab) > 8 {
		exec = routine.Crontab(j.jd.Crontab, exec)
	}
	if j.jd.Timeout > 0 {
		exec = routine.Deadline(time.Now().Add(j.jd.Timeout), exec)
	}
	return exec.Execute(ctx)
}
