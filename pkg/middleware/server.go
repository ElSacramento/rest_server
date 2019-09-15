package middleware

import (
	log "github.com/sirupsen/logrus"
	"net/http"
	"rest_server/pkg/user"
)

type server struct {
	addr string
	//db *sql.DB
	router *Router
}

func NewServer() (*server, error) {
	s := &server{addr: "127.0.0.1:8080", router: &Router{routes: make(map[string]route)}}
	return s, nil
}

func (s *server) Run() {
	s.routes()
	log.Printf("Start listen to: %s", s.addr)
	log.Fatalln("Fatal: http:", http.ListenAndServe(s.addr, s.router))
}

func (s *server) routes() {
	s.router.Get("/user", user.Get)
	s.router.Post("/user/add", user.Add)
}
