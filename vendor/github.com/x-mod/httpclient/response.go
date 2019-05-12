package httpclient

import (
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/x-mod/errors"
)

func defaultProcess(ctx context.Context, rsp *http.Response) error {
	if rsp.StatusCode == http.StatusOK {
		return nil
	}
	io.Copy(ioutil.Discard, rsp.Body)
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
	Process(context.Context, *http.Response) error
}

//ResponseProcessorFunc type
type ResponseProcessorFunc func(context.Context, *http.Response) error

//Process implemention of ResponseProcessor
func (f ResponseProcessorFunc) Process(ctx context.Context, rsp *http.Response) error {
	return f(ctx, rsp)
}

//DumpResponse struct
type DumpResponse struct {
	wr io.Writer
}

//DumpResponseOpt option
type DumpResponseOpt func(*DumpResponse)

//Output of DumpResponse
func Output(wr io.Writer) DumpResponseOpt {
	return func(d *DumpResponse) {
		d.wr = wr
	}
}

//NewDumpResponse new
func NewDumpResponse(opts ...DumpResponseOpt) *DumpResponse {
	dump := &DumpResponse{wr: os.Stdout}
	for _, opt := range opts {
		opt(dump)
	}
	return dump
}

//Process of DumpResponse
func (d *DumpResponse) Process(ctx context.Context, rsp *http.Response) error {
	defer rsp.Body.Close()
	if _, err := io.Copy(os.Stdout, rsp.Body); err != nil {
		return err
	}
	if rsp.StatusCode == http.StatusOK {
		return nil
	}
	return errors.CodeError(code(rsp.StatusCode))
}
