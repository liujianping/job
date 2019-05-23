package version

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersion(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			"default",
			"dev, commit none, built at unknown",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
	Info("v0.0.1", "now", "timestamp")
	assert.Equal(t, "v0.0.1, commit now, built at timestamp", String())
}
