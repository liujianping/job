package httpclient

import (
	"net/http"
	"time"
)

type config struct {
	request   *RequestBuilder
	response  ResponseProcessor
	transport http.RoundTripper
}

type bodyConfig struct {
	bodyType   string
	bodyObject interface{}
}

type fileConfig struct {
	Field    string
	FileName string
}

type requestConfig struct {
	Method   string
	URL      string
	Headers  map[string]string
	Queries  map[string]string
	Fragment string
	Content  *Body
}

type responseConfig struct {
	StatusCode int
	Content    *Body
}

type transportConfig struct {
	maxIdleConnsPerHost int
	retry               int
	timeout             time.Duration
}
