package httpclient

import (
	"bytes"
	"encoding/xml"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/golang/protobuf/proto"
	json "github.com/json-iterator/go"
	"github.com/x-mod/errors"
)

var types = map[string]string{
	"html":       "text/html",
	"json":       "application/json",
	"pb":         "application/octet-stream",
	"pbjson":     "application/json",
	"xml":        "application/xml",
	"text":       "text/plain",
	"binary":     "application/octet-stream",
	"urlencoded": "application/x-www-form-urlencoded",
	"form":       "application/x-www-form-urlencoded",
	"form-data":  "application/x-www-form-urlencoded",
	"multipart":  "multipart/form-data",
}

//Body struct
type Body struct {
	config *bodyConfig
}

//BodyOpt type
type BodyOpt func(*bodyConfig)

//Text opt
func Text(txt string) BodyOpt {
	return func(cf *bodyConfig) {
		cf.bodyType = "text"
		cf.bodyObject = txt
	}
}

//Binary opt
func Binary(bytes []byte) BodyOpt {
	return func(cf *bodyConfig) {
		cf.bodyType = "binary"
		cf.bodyObject = bytes
	}
}

//JSON opt
func JSON(obj map[string]interface{}) BodyOpt {
	return func(cf *bodyConfig) {
		cf.bodyType = "json"
		cf.bodyObject = obj
	}
}

//PB opt
func PB(obj proto.Message) BodyOpt {
	return func(cf *bodyConfig) {
		cf.bodyType = "pb"
		cf.bodyObject = obj
	}
}

//PBJSON opt
func PBJSON(obj proto.Message) BodyOpt {
	return func(cf *bodyConfig) {
		cf.bodyType = "pbjson"
		cf.bodyObject = obj
	}
}

//XML opt
func XML(obj map[string]interface{}) BodyOpt {
	return func(cf *bodyConfig) {
		cf.bodyType = "xml"
		cf.bodyObject = obj
	}
}

//Form opt
func Form(obj url.Values) BodyOpt {
	return func(cf *bodyConfig) {
		cf.bodyType = "form"
		cf.bodyObject = obj
	}
}

//Get Body io.Reader
func (b *Body) Get() (io.Reader, error) {
	if b.config != nil {
		switch strings.ToLower(b.config.bodyType) {
		case "text":
			return bytes.NewBufferString(b.config.bodyObject.(string)), nil
		case "binary":
			return bytes.NewBuffer(b.config.bodyObject.([]byte)), nil
		case "json":
			byts, err := json.Marshal(b.config.bodyObject.(map[string]interface{}))
			if err != nil {
				return nil, errors.Annotate(err, "json marshal failed")
			}
			return bytes.NewBuffer(byts), nil
		case "pb":
			byts, err := proto.Marshal(b.config.bodyObject.(proto.Message))
			if err != nil {
				return nil, errors.Annotate(err, "pb marshal failed")
			}
			return bytes.NewBuffer(byts), nil
		case "pbjson":
			byts, err := proto.Marshal(b.config.bodyObject.(proto.Message))
			if err != nil {
				return nil, errors.Annotate(err, "pb marshal failed")
			}
			return bytes.NewBuffer(byts), nil
		case "xml":
			byts, err := xml.Marshal(b.config.bodyObject)
			if err != nil {
				return nil, errors.Annotate(err, "xml marshal failed")
			}
			return bytes.NewBuffer(byts), nil
		case "form":
			data := b.config.bodyObject.(url.Values).Encode()
			return strings.NewReader(data), nil
		}
	}
	return bytes.NewBuffer([]byte{}), nil
}

//ContentType Body Content-Type
func (b *Body) ContentType() string {
	if b.config != nil {
		if v, ok := types[strings.ToLower(b.config.bodyType)]; ok {
			return v
		}
	}
	return types["html"]
}

//RequestBuilder struct
type RequestBuilder struct {
	config *requestConfig
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
	}
	for _, opt := range opts {
		opt(config)
	}
	return &RequestBuilder{config: config}
}

//Build http.Request
func (req *RequestBuilder) Build() (*http.Request, error) {
	if req.config == nil {
		return nil, errors.New("request config required")
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
	return rr, nil
}
