package exec

import (
	"context"
	"net/http"

	"github.com/x-mod/errors"

	"github.com/liujianping/job/config"
	"github.com/x-mod/httpclient"
	"github.com/x-mod/routine"
)

//HTTPCommand struct
type HTTPCommand struct {
	cmd *config.Command
}

//NewHTTPCommand new
func NewHTTPCommand(cmd *config.Command) routine.Executor {
	return &HTTPCommand{
		cmd: cmd,
	}
}

type _transport struct{}

//WithTransport context
func WithTransport(ctx context.Context, tr http.RoundTripper) context.Context {
	if ctx != nil {
		return context.WithValue(ctx, _transport{}, tr)
	}
	return context.WithValue(context.TODO(), _transport{}, tr)
}

//TransportFrom context
func TransportFrom(ctx context.Context) (http.RoundTripper, bool) {
	if ctx != nil {
		tr := ctx.Value(_transport{})
		if tr != nil {
			return tr.(http.RoundTripper), true
		}
	}
	return nil, false
}

//Execute of HTTPCommand
func (h *HTTPCommand) Execute(ctx context.Context) error {
	if h.cmd.HTTP == nil {
		return errors.New("command http required")
	}
	config := h.cmd.HTTP
	opts := []httpclient.ReqOpt{}
	opts = append(opts, httpclient.URL(config.Request.URL))
	opts = append(opts, httpclient.Method(config.Request.Method))
	for k, v := range config.Request.Queries {
		opts = append(opts, httpclient.Query(k, v))
	}
	for k, v := range config.Request.Headers {
		opts = append(opts, httpclient.Header(k, v))
	}
	if config.Request.Body != nil {
		if config.Request.Body.Text != nil {
			opts = append(opts, httpclient.Content(
				httpclient.Text(string(*config.Request.Body.Text)),
			))
		}
		if config.Request.Body.JSON != nil {
			opts = append(opts, httpclient.Content(
				httpclient.JSON(*config.Request.Body.JSON),
			))
		}
		if config.Request.Body.XML != nil {
			opts = append(opts, httpclient.Content(
				httpclient.XML(*config.Request.Body.XML),
			))
		}
		if config.Request.Body.Form != nil {
			opts = append(opts, httpclient.Content(
				httpclient.Form(*config.Request.Body.Form),
			))
		}
	}
	clientOpts := []httpclient.Opt{}
	clientOpts = append(clientOpts, httpclient.Request(
		httpclient.NewRequestBuilder(opts...),
	))
	if !h.cmd.Stdout {
		clientOpts = append(clientOpts, httpclient.Response(
			httpclient.NewDiscardResponse(),
		))
	} else {
		clientOpts = append(clientOpts, httpclient.Response(
			httpclient.NewDumpResponse(),
		))
	}
	if tr, ok := TransportFrom(ctx); ok {
		clientOpts = append(clientOpts, httpclient.Transport(tr))
	}
	if h.cmd.Timeout > 0 {
		clientOpts = append(clientOpts, httpclient.Timeout(h.cmd.Timeout))
	}
	client := httpclient.New(clientOpts...)
	return client.Execute(ctx)
}

//StatusCode int
type StatusCode int

//Value StatusCode
func (status StatusCode) Value() int32 {
	return int32(status)
}

//String StatusCode
func (status StatusCode) String() string {
	return http.StatusText(int(status))
}
