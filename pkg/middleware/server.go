package middleware

import (
	log "github.com/sirupsen/logrus"
	"net/http"
	"rest_server/pkg/user"
)

type server struct {
	addr string
	db *database
	router *Router
}

func NewServer() (*server, error) {
	db := &database{user: "postgres", password: "pwd", baseURL: "db:5432/postgres"}
	if err := db.createConnection(); err != nil {
		return nil, err
	}

	s := &server{addr: ":8080", router: &Router{routes: make(map[string]route)}, db: db}
	s.routes()

	return s, nil
}

func (s *server) Run() {
	log.Printf("Start listen to: %s", s.addr)
	go func() {
		log.Fatalln("Fatal: http:", http.ListenAndServe(s.addr, s.router))
	}()
}

func (s *server) Stop() error {
	log.Printf("Stop listen to: %s", s.addr)
	if err := s.db.closeConnection(); err != nil {
		return err
	}
	return nil
}

func (s *server) routes() {
	s.router.Get("/user", user.Get)
	s.router.Post("/user/add", user.Add)
}
