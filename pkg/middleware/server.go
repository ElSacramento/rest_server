package middleware

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"rest_server/pkg/database"
	"rest_server/pkg/user"
)

type server struct {
	host   string
	port   string
	addr   string
	router *Router

	db database.DataStore

	user *user.Service
}

func NewServer(cfg *Config) (*server, error) {
	db := &database.PostgresDB{
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		Name:     cfg.Database.Name,
	}
	if err := db.CreateConnection(); err != nil {
		return nil, err
	}

	s := &server{
		host:   cfg.Server.Host,
		port:   cfg.Server.Port,
		addr:   cfg.Server.Host + ":" + cfg.Server.Port,
		router: &Router{routes: make(map[string]route)},
		db:     db,
		user:   &user.Service{DB: db},
	}
	s.routes()

	return s, nil
}

func (s *server) Run() {
	logrus.Printf("Start listen to: %s", s.addr)
	go func() {
		logrus.Fatalln("Fatal: http:", http.ListenAndServe(s.addr, s.router))
	}()
}

func (s *server) Stop() error {
	logrus.Printf("Stop listen to: %s", s.addr)
	if err := s.db.CloseConnection(); err != nil {
		return err
	}
	return nil
}

func (s *server) routes() {
	s.router.Get("/user", s.user.Get)
	s.router.Post("/user/add", s.user.Add)
}
