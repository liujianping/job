package exec

import (
	"context"
	"testing"
	"time"

	"github.com/liujianping/job/config"
	"github.com/stretchr/testify/assert"
)

func TestJob_Execute(t *testing.T) {
	cmd := config.CommandJD()
	opts := []config.Option{}
	opt := config.Name("name")
	assert.NotNil(t, opt)
	opts = append(opts, opt)

	opt = config.Metadata("key", "val")
	assert.NotNil(t, opt)
	opts = append(opts, opt)

	opt = config.CommandName("echo")
	assert.NotNil(t, opt)
	opts = append(opts, opt)

	opt = config.CommandArgs("aa", "bb", "cc")
	assert.NotNil(t, opt)
	opts = append(opts, opt)

	opt = config.CommandEnv("key", "val")
	assert.NotNil(t, opt)
	opts = append(opts, opt)

	opt = config.CommandRetry(3)
	assert.NotNil(t, opt)
	opts = append(opts, opt)

	opt = config.CommandStdoutDiscard(true)
	assert.NotNil(t, opt)
	opts = append(opts, opt)

	opt = config.CommandTimeout(time.Second)
	assert.NotNil(t, opt)
	opts = append(opts, opt)

	for _, o := range opts {
		o(cmd)
	}

	job := NewJob(cmd, nil)
	assert.NotNil(t, job)
	err := job.Execute(context.TODO())
	assert.Nil(t, err)
}
