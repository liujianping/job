package cmd

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestOutput(t *testing.T) {
	outCmd := RootCmd()
	outCmd.Flags().Set("config", "../etc/job.yaml")
	outCmd.Flags().Set("output", "true")
	assert.Nil(t, Main(outCmd, []string{}))
}

func TestMain(t *testing.T) {
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

	mainCmd := RootCmd()
	mainCmd.Flags().Set("config", "../etc/job.yaml")
	mainCmd.Flags().Set("report", "true")
	assert.Nil(t, Main(mainCmd, []string{}))
}

func TestVersion(t *testing.T) {
	verCmd := RootCmd()
	verCmd.Flags().Set("version", "true")
	assert.Nil(t, Main(verCmd, []string{}))
}

func TestReport(t *testing.T) {
	rptCmd := RootCmd()
	rptCmd.Flags().Set("report", "true")
	assert.Nil(t, Main(rptCmd, []string{"echo", "hello"}))
}

func TestVerbose(t *testing.T) {
	verbCmd := RootCmd()
	verbCmd.Flags().Set("verbose", "true")
	assert.Nil(t, Main(verbCmd, []string{"echox", "hello"}))
}
