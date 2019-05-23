package exec

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/liujianping/job/config"
	"github.com/stretchr/testify/assert"
)

func TestHTTPCommand_Execute(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.RequestURI {
			case "/get":
				w.WriteHeader(http.StatusOK)
				io.WriteString(w, `ok`)
				return
			case "/post":
				w.WriteHeader(http.StatusOK)
				io.WriteString(w, `ok`)
				return
			case "/json":
				w.WriteHeader(http.StatusOK)
				io.WriteString(w, `ok`)
				return
			case "/text":
				w.WriteHeader(http.StatusOK)
				io.WriteString(w, `ok`)
				return
			case "/xml":
				w.WriteHeader(http.StatusOK)
				io.WriteString(w, `ok`)
				return
			case "/timeout":
				time.Sleep(5 * time.Second)
				w.WriteHeader(http.StatusOK)
				io.WriteString(w, `ok`)
				return
			default:
				http.NotFound(w, r)
			}
		}),
	)
	defer ts.Close()

	jds, err := config.ParseJDs("../etc/http.yaml")
	assert.Nil(t, err)

	for _, jd := range jds {
		jd.Command.HTTP.Request.URL = strings.Replace(jd.Command.HTTP.Request.URL, "http://localhost:8080", ts.URL, -1)
		job := NewJob(jd, nil)

		if job.String() == "timeout" {
			assert.NotNil(t, job.Execute(context.TODO()), job.String())
		} else {
			assert.Nil(t, job.Execute(context.TODO()), job.String())
		}
	}
}
