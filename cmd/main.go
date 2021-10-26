package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/morzhanov/loki-grafana/internal/config"

	"github.com/morzhanov/loki-grafana/internal/logger"
	"github.com/morzhanov/loki-grafana/internal/rest"
)

func main() {
	c, err := config.NewConfig()
	if err != nil {
		log.Fatalf("error duting config initialization %s", err.Error())
	}
	l, err := logger.NewLogger(c.LokiClientURL)
	if err != nil {
		log.Fatal(fmt.Errorf("initialization error in step %s: %w", "logger", err))
	}
	l.Info("logger created")

	r := rest.NewREST(l, c.Version)
	l.Info("rest controller created")

	go r.Listen()
	l.Info("rest controller accepts a connections on the port :8080")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	log.Println("App successfully started!")
	<-quit
	log.Println("received os.Interrupt, exiting...")
}
