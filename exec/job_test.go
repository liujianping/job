package exec

import (
	"context"
	"testing"

	"github.com/liujianping/job/config"
	"github.com/stretchr/testify/assert"
)

func TestJob_Execute(t *testing.T) {
	cmd := config.CommandJD()
	cmd.Command.Shell.Name = "echo"
	cmd.Command.Shell.Args = []string{"hello", "world"}

	job := NewJob(cmd, nil)
	assert.NotNil(t, job)
	err := job.Execute(context.TODO())
	assert.Nil(t, err)

	http := config.HTTPCommandJD()
	http.Command.HTTP.Request.URL = "https://github.com"
	http.Command.HTTP.Request.Method = "GET"

	job2 := NewJob(http, nil)
	assert.NotNil(t, job2)
	err2 := job2.Execute(context.TODO())
	assert.Nil(t, err2)
}
