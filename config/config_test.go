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
	opts = append(opts, Name("name"))
	opts = append(opts, Metadata("key", "val"))
	opts = append(opts, CommandName("echo"))
	opts = append(opts, CommandArgs("aa", "bb", "cc"))
	opts = append(opts, CommandEnv("key", "val"))
	opts = append(opts, CommandRetry(3))
	opts = append(opts, CommandStdoutDiscard(true))
	opts = append(opts, CommandTimeout(time.Second))
	opts = append(opts, RepeatTimes(10))
	opts = append(opts, RepeatInterval(time.Second))
	opts = append(opts, Concurrent(10))
	opts = append(opts, Crontab("* * * * *"))
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
