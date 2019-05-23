package cmd

import (
	"context"
	"testing"

	"github.com/liujianping/job/config"
	"github.com/stretchr/testify/assert"
)

func TestJOBs(t *testing.T) {
	job1 := config.CommandJD()
	job1.Name = "one"
	job1.Command.Shell.Name = "echo"
	job1.Order.Weight = 1
	job1.Order.Wait = true

	job2 := config.CommandJD()
	job2.Name = "two"
	job2.Command.Shell.Name = "echo"
	job2.Order.Weight = 2
	job1.Order.Wait = true
	job2.Order.Precondition = []string{"three"}

	job3 := config.CommandJD()
	job3.Name = "three"
	job3.Command.Shell.Name = "echo"
	job3.Order.Weight = 3

	jobs := []*config.JD{job2, job3, job1}

	JBs := NewJOBs(jobs, nil)
	assert.NotNil(t, JBs)
	JBs.Sort()

	err := JBs.Execute(context.TODO())
	assert.Nil(t, err)
}

func TestCmd(t *testing.T) {
	jds, err := config.ParseJDs("../etc/job.yaml")
	assert.Nil(t, err)
	assert.Equal(t, 7, len(jds))

	JBs := NewJOBs(jds, nil)
	assert.NotNil(t, JBs)
	JBs.Sort()
	assert.Nil(t, JBs.Execute(context.TODO()))
}
