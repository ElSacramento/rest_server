package middleware

import (
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

// Router implements http.Handler interface
type Router struct {
	routes map[string]route
}

type route struct {
	handler http.Handler
	pattern string
	method string
}

// Add route with method GET
func (r *Router) Get(pattern string, h http.HandlerFunc) {
	r.routes[pattern] = route{handler: h, pattern: pattern, method: "GET"}
}

// Add route with method POST
func (r *Router) Post(pattern string, h http.HandlerFunc) {
	r.routes[pattern] = route{handler: h, pattern: pattern, method: "POST"}
}

func (r *Router) ServeHTTP(w http.ResponseWriter, request *http.Request) {
	method := strings.ToUpper(request.Method)
	foundRoute, exist := r.routes[request.URL.Path]
	if !exist {
		log.Infof("url not found: %s", request.URL.Path)
		http.NotFound(w, request)
		return
	}
	if foundRoute.method != method {
		log.Infof("wrong http method, expected %s, got %s", foundRoute.method, method)
		http.Error(w, http.ErrBodyNotAllowed.Error(), http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	foundRoute.handler.ServeHTTP(w, request)
}
