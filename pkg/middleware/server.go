package middleware

import (
	"context"
	"net/http"
	"rest_server/pkg/database"
	"rest_server/pkg/user"
	"time"

	"github.com/sirupsen/logrus"
)

type server struct {
	host   string
	port   string
	addr   string
	router *Router
	server *http.Server

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
	s.server = &http.Server{Addr: s.addr, Handler: s.router}
	err := s.server.ListenAndServe()
	switch err {
	case http.ErrServerClosed:
		logrus.Print("Server shut down")
	default:
		if err := s.db.CloseConnection(); err != nil {
			logrus.Error("Problem with closing connection: ", err)
		}
		logrus.Fatalln("Fatal: http:", err)
	}
}

func (s *server) Stop() {
	hasError := false
	s.server.SetKeepAlivesEnabled(false)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := s.server.Shutdown(ctx); err != nil {
		hasError = true
		logrus.Error("Could not gracefully shutdown the server: ", err)
	}
	logrus.Printf("Stop listen to: %s", s.addr)
	if err := s.db.CloseConnection(); err != nil {
		hasError = true
		logrus.Error("Problem with closing connection: ", err)
	}
	if hasError {
		panic("Failed to correct stopped server")
	}
}

func (s *server) routes() {
	s.router.Get("/user", s.user.Get)
	s.router.Post("/user/add", s.user.Add)
}
