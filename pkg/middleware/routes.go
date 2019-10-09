package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

// Router implements http.Handler interface
type Router struct {
	routes map[string]route
}

type route struct {
	handler http.Handler
	pattern string
	method  string
}

// Add route with method GET
func (r *Router) Get(pattern string, h http.HandlerFunc) {
	r.routes[pattern] = route{handler: h, pattern: pattern, method: http.MethodGet}
}

// Add route with method POST
func (r *Router) Post(pattern string, h http.HandlerFunc) {
	r.routes[pattern] = route{handler: h, pattern: pattern, method: http.MethodPost}
}

func (r *Router) ServeHTTP(w http.ResponseWriter, request *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	rec := HTTPRecorder{w: w, req: request, status: 200, start: time.Now()}
	rec.LogRequest()

	method := request.Method
	foundRoute, exist := r.routes[request.URL.Path]
	if !exist {
		http.NotFound(&rec, request)
		rec.LogResponse()
		return
	}
	if foundRoute.method != method {
		logrus.Infof("wrong http method, expected %s, got %s", foundRoute.method, method)
		http.Error(&rec, "wrong request http method", http.StatusMethodNotAllowed)
		rec.LogResponse()
		return
	}

	foundRoute.handler.ServeHTTP(&rec, request.WithContext(ctx))
	rec.LogResponse()
}
