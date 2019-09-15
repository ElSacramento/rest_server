package main

import (
	log "github.com/sirupsen/logrus"
	"rest_server/pkg/middleware"
)

func main() {
	server, err := middleware.NewServer()
	if err != nil {
		log.WithError(err).Panic("Failed to initialize server")
	}
	server.Run()
}
