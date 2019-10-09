package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

type HTTPRecorder struct {
	w      http.ResponseWriter
	req    *http.Request
	status int
	body   []byte
	start  time.Time
}

func (r *HTTPRecorder) Header() http.Header {
	return r.w.Header()
}

func (r *HTTPRecorder) Write(b []byte) (int, error) {
	r.body = b
	return r.w.Write(b)
}

func (r *HTTPRecorder) WriteHeader(statusCode int) {
	r.status = statusCode
	r.w.WriteHeader(statusCode)
}

func (r *HTTPRecorder) LogRequest() {
	agent := r.req.Header.Get("User-Agent")
	conn := r.req.Header.Get("Connection")
	headersSet := fmt.Sprintf("Host: %q Connection: %q User-Agent: %q", r.req.Host, conn, agent)
	logrus.Infof("%s %s %s %s", r.req.Method, r.req.RequestURI, r.req.Proto, headersSet)
}

func (r *HTTPRecorder) LogResponse() {
	// log response body if error
	if r.status > http.StatusOK {
		logrus.Errorf("[%d] %s %q %s", r.status, r.req.RequestURI, string(r.body), time.Since(r.start))
		return
	}
	logrus.Infof("[%d] %s %s", r.status, r.req.RequestURI, time.Since(r.start))
}
