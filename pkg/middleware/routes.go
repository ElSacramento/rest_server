package middleware

import (
	"context"
	"github.com/sirupsen/logrus"
	"net/http"
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
	method := request.Method
	foundRoute, exist := r.routes[request.URL.Path]
	if !exist {
		logrus.Infof("url not found: %s", request.URL.Path)
		http.NotFound(w, request)
		return
	}
	if foundRoute.method != method {
		logrus.Infof("wrong http method, expected %s, got %s", foundRoute.method, method)
		http.Error(w, "wrong request http method", http.StatusMethodNotAllowed)
		return
	}
	foundRoute.handler.ServeHTTP(w, request.WithContext(context.Background()))
}
