package cmd

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestOutput(t *testing.T) {
	cmd := RootCmd()
	viper.BindPFlags(cmd.Flags())
	cmd.Flags().Set("config", "../etc/job.yaml")
	cmd.Flags().Set("output", "true")
	assert.Nil(t, Main(cmd, []string{}))
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

	cmd := RootCmd()
	viper.BindPFlags(cmd.Flags())
	cmd.Flags().Set("config", "../etc/job.yaml")
	cmd.Flags().Set("report", "true")
	assert.Nil(t, Main(cmd, []string{}))
}

func TestVersion(t *testing.T) {
	cmd := RootCmd()
	viper.BindPFlags(cmd.Flags())
	cmd.Flags().Set("version", "true")
	assert.Nil(t, Main(cmd, []string{}))
}

func TestReport(t *testing.T) {
	cmd := RootCmd()
	viper.BindPFlags(cmd.Flags())
	cmd.Flags().Set("report", "true")
	assert.Nil(t, Main(cmd, []string{"echo", "hello"}))
}

func TestVerbose(t *testing.T) {
	cmd := RootCmd()
	viper.BindPFlags(cmd.Flags())
	cmd.Flags().Set("verbose", "true")
	assert.Nil(t, Main(cmd, []string{"echox", "hello"}))
}
