package main

import (
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"rest_server/pkg/middleware"
	"syscall"
)

func main() {
	server, err := middleware.NewServer()
	if err != nil {
		log.WithError(err).Panic("Failed to initialize server")
	}
	server.Run()

	// Handle SIGINT and SIGTERM.
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	s := <-ch
	log.Printf("Got signal: %v", s)

	if err := server.Stop(); err != nil {
		log.WithError(err).Panic("Failed to shutdown server")
	}
}
