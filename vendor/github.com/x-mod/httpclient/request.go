package httpclient

import (
	"io"
	"net/http"

	"github.com/x-mod/errors"
)

//RequestBuilder struct
type RequestBuilder struct {
	config  *requestConfig
	request *http.Request
}

//ReqOpt opt
type ReqOpt func(*requestConfig)

//Method opt
func Method(method string) ReqOpt {
	return func(cf *requestConfig) {
		cf.Method = method
	}
}

//URL opt
func URL(url string) ReqOpt {
	return func(cf *requestConfig) {
		cf.URL = url
	}
}

//Query opt
func Query(name string, value string) ReqOpt {
	return func(cf *requestConfig) {
		cf.Queries[name] = value
	}
}

//Header opt
func Header(name string, value string) ReqOpt {
	return func(cf *requestConfig) {
		cf.Headers[name] = value
	}
}

//Cookie opt
func Cookie(cookie *http.Cookie) ReqOpt {
	return func(cf *requestConfig) {
		if cookie != nil {
			cf.Cookies = append(cf.Cookies, cookie)
		}
	}
}

//BasicAuth opt
func BasicAuth(username string, password string) ReqOpt {
	return func(cf *requestConfig) {
		cf.Auth = &authConfig{
			username: username,
			password: password,
		}
	}
}

//Fragment opt
func Fragment(name string) ReqOpt {
	return func(cf *requestConfig) {
		cf.Fragment = name
	}
}

//Content opt
func Content(opts ...BodyOpt) ReqOpt {
	return func(cf *requestConfig) {
		body := &bodyConfig{}
		for _, opt := range opts {
			opt(body)
		}
		cf.Content = &Body{config: body}
	}
}

//NewRequestBuilder new
func NewRequestBuilder(opts ...ReqOpt) *RequestBuilder {
	config := &requestConfig{
		Headers: make(map[string]string),
		Queries: make(map[string]string),
		Cookies: []*http.Cookie{},
	}
	for _, opt := range opts {
		opt(config)
	}
	return &RequestBuilder{config: config}
}

func (req *RequestBuilder) makeRequest() (*http.Request, error) {
	if len(req.config.URL) == 0 {
		return nil, errors.New("url required")
	}
	//body
	var body io.Reader
	if req.config.Content != nil {
		rd, err := req.config.Content.Get()
		if err != nil {
			return nil, err
		}
		body = rd
	}
	rr, err := http.NewRequest(req.config.Method, req.config.URL, body)
	if err != nil {
		return nil, err
	}
	// queries
	if len(req.config.Queries) > 0 {
		q := rr.URL.Query()
		for k, v := range req.config.Queries {
			q.Add(k, v)
		}
		rr.URL.RawQuery = q.Encode()
	}
	// fragment
	if len(req.config.Fragment) > 0 {
		rr.URL.Fragment = req.config.Fragment
	}
	// content-type
	if req.config.Content != nil {
		rr.Header.Set("Content-Type", req.config.Content.ContentType())
	}
	// headers, can replace content-type
	for k, v := range req.config.Headers {
		rr.Header.Set(k, v)
	}
	// cookies
	for _, v := range req.config.Cookies {
		rr.AddCookie(v)
	}
	// auth
	if req.config.Auth != nil {
		rr.SetBasicAuth(req.config.Auth.username, req.config.Auth.password)
	}
	req.request = rr
	return rr, nil
}

//Get http.Request
func (req *RequestBuilder) Get() (*http.Request, error) {
	if req.request == nil {
		return req.makeRequest()
	}
	return req.request, nil
}

//Clear RequestBuilder
func (req *RequestBuilder) Clear() {
	req.request = nil
}
