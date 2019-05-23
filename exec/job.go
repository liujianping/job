package exec

import (
	"context"
	"time"

	"github.com/liujianping/job/config"
	uuid "github.com/satori/go.uuid"
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
	var exec routine.Executor
	if j.jd.Report {
		j.jd.Command.Stdout = false
	}
	if j.jd.Command.Shell != nil {
		exec = NewShellCommand(&j.jd.Command)
	}
	if j.jd.Command.HTTP != nil {
		exec = NewHTTPCommand(&j.jd.Command)
	}
	if j.jd.Guarantee {
		exec = routine.Guarantee(exec)
	}
	if j.jd.Report && j.reporter != nil {
		exec = routine.Report(j.reporter.Report(), exec)
	}
	if j.jd.Repeat.Times != 1 {
		exec = routine.Repeat(j.jd.Repeat.Times, j.jd.Repeat.Interval, exec)
	}
	//* * * * *
	if len(j.jd.Crontab) > 8 {
		exec = routine.Crontab(j.jd.Crontab, exec)
	}
	if j.jd.Concurrent > 0 {
		exec = routine.Concurrent(j.jd.Concurrent, exec)
	}
	if j.jd.Timeout > 0 {
		exec = routine.Deadline(time.Now().Add(j.jd.Timeout), exec)
	}
	return exec.Execute(ctx)
}
