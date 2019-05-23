package exec

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/x-mod/routine"
)

func TestReporter_Execute(t *testing.T) {
	report := NewReporter(10)
	assert.NotNil(t, report)

	ch := report.Report()
	assert.NotNil(t, ch)

	for i := 0; i < 10; i++ {
		ch <- &routine.Result{
			Code: i,
		}
	}
	report.Stop()

	assert.Nil(t, report.Execute(context.TODO()))
	report.Finalize()
}
