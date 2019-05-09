package cmd

import (
	"context"
	"sort"

	"github.com/liujianping/job/config"
	"github.com/liujianping/job/exec"
	"github.com/x-mod/routine"
)

//JOBs type
type JOBs struct {
	jds    []*config.JD
	report *exec.Reporter
}

//NewJOBs new JOBs
func NewJOBs(jds []*config.JD, report *exec.Reporter) *JOBs {
	return &JOBs{
		jds:    jds,
		report: report,
	}
}

//Len of JOBs
func (jobs JOBs) Len() int {
	return len(jobs.jds)
}

//Less Cmp of JOBs
func (jobs JOBs) Less(i, j int) bool {
	return jobs.jds[i].Order.Weight < jobs.jds[j].Order.Weight
}

//Swap of JOBs
func (jobs JOBs) Swap(i, j int) {
	jobs.jds[i], jobs.jds[j] = jobs.jds[j], jobs.jds[i]
}

//Sort of JOBs
func (jobs JOBs) Sort() {
	sort.Sort(jobs)
}

//Execute impl executor
func (jobs JOBs) Execute(ctx context.Context) error {
	jmap := make(map[string]chan error, jobs.Len())
	var tail chan error
	for _, jd := range jobs.jds {
		if len(jd.Order.Precondition) == 0 {
			job := exec.NewJob(jd, jobs.report)
			ch := routine.Go(ctx, job)
			jmap[job.String()] = ch
			tail = ch
			if jd.Order.Wait {
				<-ch
			}
		}
	}
	for _, jd := range jobs.jds {
		for _, pre := range jd.Order.Precondition {
			if ch, ok := jmap[pre]; ok {
				<-ch
			}
			job := exec.NewJob(jd, jobs.report)
			ch := routine.Go(ctx, job)
			jmap[job.String()] = ch
			tail = ch
			if jd.Order.Wait {
				<-ch
			}
		}
	}
	return <-tail
}
