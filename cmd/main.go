package main

import (
	"flag"
	"os"
	"os/signal"
	"rest_server/pkg/middleware"
	"syscall"

	"github.com/sirupsen/logrus"
)

func main() {
	cfgPath := flag.String("cfg", "", "cfg file path")
	flag.Parse()

	if *cfgPath == "" {
		logrus.Fatalln("Config for server is not set")
	}

	cfg, err := middleware.ParseConfig(cfgPath)
	if err != nil {
		logrus.WithError(err).Panic("Failed to parse config")
	}

	if err := cfg.ValidateConfig(); err != nil {
		logrus.WithError(err).Panic("Failed to validate config")
	}

	logrus.Infof("Config: %+v", cfg)

	server, err := middleware.NewServer(cfg)
	if err != nil {
		logrus.WithError(err).Panic("Failed to initialize server")
	}
	go server.Run()

	// Handle SIGINT and SIGTERM.
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	s := <-ch
	logrus.Printf("Got signal: %v", s)

	server.Stop()
	close(ch)
}
