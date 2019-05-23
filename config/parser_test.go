package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseJDs(t *testing.T) {
	jds, err := ParseJDs("../etc/job.yaml")
	assert.Nil(t, err)
	assert.Equal(t, 7, len(jds))
}
