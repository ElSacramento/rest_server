package middleware

import (
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/sirupsen/logrus"
)

type ResponseRecorder struct {
	w      http.ResponseWriter
	status int
	body   []byte
}

func (r ResponseRecorder) Header() http.Header {
	return r.w.Header()
}

func (r *ResponseRecorder) Write(b []byte) (int, error) {
	r.body = b
	return r.w.Write(b)
}

func (r *ResponseRecorder) WriteHeader(statusCode int) {
	r.status = statusCode
	r.w.WriteHeader(statusCode)
}

func logHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		d, err := httputil.DumpRequest(r, false)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		start := time.Now()
		logrus.Printf("%s", string(d))

		rec := ResponseRecorder{w: w, status: 200}
		fn(&rec, r)

		// log response body if error
		if rec.status > http.StatusOK {
			logrus.Printf("[%d] %s %q %s", rec.status, r.RequestURI, string(rec.body), time.Since(start))
			return
		}
		logrus.Printf("[%d] %s %s", rec.status, r.RequestURI, time.Since(start))
	}
}
