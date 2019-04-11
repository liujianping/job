package exec

import (
	"context"

	"github.com/x-mod/routine"

	"github.com/liujianping/job/config"
)

type HTTPExecutor struct {
	Http *config.HTTP
}

func NewHTTPExecutor(h *config.HTTP) routine.Executor {
	return &HTTPExecutor{
		Http: h,
	}
}

func (h *HTTPExecutor) Execute(ctx context.Context) error {
	return nil
}
