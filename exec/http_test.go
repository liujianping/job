package exec

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHTTPCommand_Execute(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.RequestURI {
			case "/head":
				if r.Header.Get("X-HEAD") != "x-head-value" {
					http.Error(w, "head not equal", http.StatusBadRequest)
					return
				}
				w.WriteHeader(http.StatusOK)
				io.WriteString(w, `ok`)
				return
			case "/auth":
				if user, pass, ok := r.BasicAuth(); ok {
					if user == "jay" && pass == "123" {
						w.WriteHeader(http.StatusOK)
						io.WriteString(w, `ok`)
						return
					}
				}
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			case "/error":
				http.Error(w, "error", http.StatusBadRequest)
				return
			case "/sleep":
				time.Sleep(time.Second)
				w.WriteHeader(http.StatusOK)
				io.WriteString(w, `sleeped`)
				return
			case "/ping":
				w.WriteHeader(http.StatusOK)
				io.WriteString(w, `ok`)
				return
			default:
				http.NotFound(w, r)
			}
			if r.Header.Get("X-HEAD") != "x-head-value" {
				http.Error(w, "head not equal", http.StatusBadRequest)
				return
			}
			w.WriteHeader(http.StatusOK)
			io.WriteString(w, `ok`)
		}),
	)
	defer ts.Close()
}
