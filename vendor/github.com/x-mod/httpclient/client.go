package httpclient

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/x-mod/errors"
)

//Client struct
type Client struct {
	config *config
	m      sync.Once
	con    *http.Client
}

//Opt for client
type Opt func(*config)

//Request opt
func Request(builder *RequestBuilder) Opt {
	return func(cf *config) {
		cf.request = builder
	}
}

//Transport opt
func Transport(transport http.RoundTripper) Opt {
	return func(cf *config) {
		if cf.client == nil {
			cf.client = &http.Client{
				Transport: &http.Transport{},
			}
		}
		cf.client.Transport = transport
	}
}

//Timeout opt
func Timeout(duration time.Duration) Opt {
	return func(cf *config) {
		if cf.client == nil {
			cf.client = &http.Client{
				Transport: &http.Transport{},
			}
		}
		cf.client.Timeout = duration
	}
}

//Connection opt
func Connection(client *http.Client) Opt {
	return func(cf *config) {
		cf.client = client
	}
}

//Response opt
func Response(processor ResponseProcessor) Opt {
	return func(cf *config) {
		cf.response = processor
	}
}

//New client
func New(opts ...Opt) *Client {
	cf := &config{}
	for _, opt := range opts {
		opt(cf)
	}
	return &Client{config: cf}
}

func (c *Client) get() *http.Client {
	if c.config.client != nil {
		return c.config.client
	}
	return &http.Client{
		Transport: &http.Transport{},
	}
}

//Do client
func (c *Client) Do(ctx context.Context) (*http.Response, error) {
	if c.config.request == nil {
		return nil, errors.New("request required")
	}

	req, err := c.config.request.Build()
	if err != nil {
		return nil, err
	}
	return c.get().Do(req.WithContext(ctx))
}

//Execute client
func (c *Client) Execute(ctx context.Context) error {
	rsp, err := c.Do(ctx)
	if err != nil {
		return err
	}

	if c.config.response != nil {
		return c.config.response.Process(rsp)
	}
	return defaultProcess(rsp)
}
