package httpclient

import (
	"net/http"

	"github.com/facebookgo/httpcontrol"
)

//NewHTTPTransport new http transport
func NewHTTPTransport(opts ...RoundTripperOpt) http.RoundTripper {
	defaultConfig := &transportConfig{
		maxIdleConnsPerHost: 0,
		retry:               0,
		timeout:             0,
	}
	for _, o := range opts {
		o(defaultConfig)
	}

	return &httpcontrol.Transport{
		MaxTries:            uint(defaultConfig.retry),
		MaxIdleConnsPerHost: defaultConfig.maxIdleConnsPerHost,
	}
}

//RoundTripperOpt option
type RoundTripperOpt func(*transportConfig)

//MaxIdleConnections opt
func MaxIdleConnections(max int) RoundTripperOpt {
	return func(cf *transportConfig) {
		cf.maxIdleConnsPerHost = max
	}
}

//Retry opt
func Retry(retry int) RoundTripperOpt {
	return func(cf *transportConfig) {
		cf.retry = retry
	}
}
