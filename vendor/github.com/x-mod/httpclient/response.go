package httpclient

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"

	json "github.com/json-iterator/go"
	"github.com/x-mod/errors"
)

func defaultProcess(rsp *http.Response) error {
	if rsp.StatusCode == http.StatusOK {
		return nil
	}
	return errors.CodeError(code(rsp.StatusCode))
}

//http code
type code int

func (c code) Value() int32 {
	return int32(c)
}

func (c code) String() string {
	return http.StatusText(int(c))
}

//ResponseProcessor interface
type ResponseProcessor interface {
	Process(*http.Response) error
}

//DumpResponse struct
type DumpResponse struct {
}

//NewDumpResponse new
func NewDumpResponse() *DumpResponse {
	return &DumpResponse{}
}

//Process of DumpResponse
func (d *DumpResponse) Process(rsp *http.Response) error {
	defer rsp.Body.Close()
	io.Copy(os.Stdout, rsp.Body)
	if rsp.StatusCode == http.StatusOK {
		return nil
	}
	return errors.CodeError(code(rsp.StatusCode))
}

//DiscardResponse struct
type DiscardResponse struct {
}

//NewDiscardResponse new
func NewDiscardResponse() *DiscardResponse {
	return &DiscardResponse{}
}

//Process of DiscardResponse
func (d *DiscardResponse) Process(rsp *http.Response) error {
	defer rsp.Body.Close()
	io.Copy(ioutil.Discard, rsp.Body)
	if rsp.StatusCode == http.StatusOK {
		return nil
	}
	return errors.CodeError(code(rsp.StatusCode))
}

//CompareResponse struct
type CompareResponse struct {
	status int
	body   *Body
}

//NewCompareResponse new
func NewCompareResponse(status int, opts ...BodyOpt) *CompareResponse {
	config := &bodyConfig{}
	for _, opt := range opts {
		opt(config)
	}
	return &CompareResponse{status: status, body: &Body{config: config}}
}

//Process of CompareResponse
func (cmp *CompareResponse) Process(rsp *http.Response) error {
	if rsp.StatusCode != cmp.status {
		return errors.Errorf("response compare failed: want (%d), get (%d)", cmp.status, rsp.StatusCode)
	}
	defer rsp.Body.Close()
	if cmp.body != nil {
		switch cmp.body.config.bodyType {
		case "text":
			b, err := ioutil.ReadAll(rsp.Body)
			if err != nil {
				return errors.Annotate(err, "response body read failed")
			}
			if !reflect.DeepEqual(cmp.body.config.bodyObject, string(b)) {
				return errors.Errorf("text compare failed: want (%v), get (%s)", cmp.body.config.bodyObject, string(b))
			}
		case "binary":
			b, err := ioutil.ReadAll(rsp.Body)
			if err != nil {
				return errors.Annotate(err, "response body read failed")
			}
			if !reflect.DeepEqual(cmp.body.config.bodyObject, b) {
				return errors.Errorf("binary compare failed: want (%v), get (%s)", cmp.body.config.bodyObject, string(b))
			}
		case "json":
			var obj interface{}
			if err := json.NewDecoder(rsp.Body).Decode(&obj); err != nil {
				return errors.Annotate(err, "response body json decode failed")
			}
			if !reflect.DeepEqual(cmp.body.config.bodyObject, obj) {
				return errors.Errorf("json compare failed: want (%v), get (%v)", cmp.body.config.bodyObject, obj)
			}
		case "pb":
		case "pbjson":
		case "xml":
		case "form":
		}
	}
	return nil
}
