package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCommandJD(t *testing.T) {
	jd1 := CommandJD()
	assert.NotNil(t, jd1)
	jd2 := HTTPCommandJD()
	assert.NotNil(t, jd2)

	opts := []Option{}
	opt := Name("name")
	assert.NotNil(t, opt)
	opts = append(opts, opt)

	opt = Metadata("key", "val")
	assert.NotNil(t, opt)
	opts = append(opts, opt)

	opt = CommandName("echo")
	assert.NotNil(t, opt)
	opts = append(opts, opt)

	opt = CommandArgs("aa", "bb", "cc")
	assert.NotNil(t, opt)
	opts = append(opts, opt)

	opt = CommandEnv("key", "val")
	assert.NotNil(t, opt)
	opts = append(opts, opt)

	opt = CommandRetry(3)
	assert.NotNil(t, opt)
	opts = append(opts, opt)

	opt = CommandStdoutDiscard(true)
	assert.NotNil(t, opt)
	opts = append(opts, opt)

	opt = CommandTimeout(time.Second)
	assert.NotNil(t, opt)
	opts = append(opts, opt)

	opt = RepeatTimes(10)
	assert.NotNil(t, opt)
	opts = append(opts, opt)

	opt = RepeatInterval(time.Second)
	assert.NotNil(t, opt)
	opts = append(opts, opt)

	opt = Concurrent(10)
	assert.NotNil(t, opt)
	opts = append(opts, opt)

	opt = Crontab("* * * * *")
	assert.NotNil(t, opt)
	opts = append(opts, opt)

	for _, opt := range opts {
		opt(jd1)
	}
	assert.Equal(t, "name", jd1.Name)
	assert.Equal(t, map[string]interface{}{"key": "val"}, jd1.Metadata)
	assert.Equal(t, "echo", jd1.Command.Shell.Name)
	assert.Equal(t, []string{"aa", "bb", "cc"}, jd1.Command.Shell.Args)
	assert.Equal(t, []KV{{"key", "val"}}, jd1.Command.Shell.Envs)
	assert.Equal(t, 3, jd1.Command.Retry)
	assert.Equal(t, true, !jd1.Command.Stdout)
	assert.Equal(t, time.Second, jd1.Command.Timeout)
	assert.Equal(t, 10, jd1.Repeat.Times)
	assert.Equal(t, time.Second, jd1.Repeat.Interval)
	assert.Equal(t, "* * * * *", jd1.Crontab)
}
