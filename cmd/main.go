package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/morzhanov/go-elk-example/internal/rest"

	"github.com/morzhanov/go-elk-example/internal/metrics"

	"github.com/morzhanov/go-elk-example/internal/es"
	"github.com/morzhanov/go-elk-example/internal/generator"

	"github.com/morzhanov/go-elk-example/internal/logger"

	"go.uber.org/zap"

	"github.com/morzhanov/go-elk-example/internal/config"
)

func failOnError(l *zap.Logger, step string, err error) {
	if err != nil {
		l.Fatal("initialization error", zap.Error(err), zap.String("step", step))
	}
}

func main() {
	l, err := logger.NewLogger()
	if err != nil {
		log.Fatal(fmt.Errorf("initialization error in step %s: %w", "logger", err))
	}
	l.Info("logger created")
	c, err := config.NewConfig()
	failOnError(l, "config", err)
	l.Info("config created")
	m := metrics.NewMetricsCollector()
	l.Info("metrics collector created")

	esearch, err := es.NewES(c.ESuri, c.ESindex, l)
	failOnError(l, "elastic search client", err)
	l.Info("elastic search client created")
	g := generator.NewGenerator(esearch, l, m)
	l.Info("generator created")
	r := rest.NewREST(esearch, l, m)
	l.Info("rest controller created")

	go g.Generate()
	l.Info("documents generator started...")
	go r.Listen()
	l.Info("rest controller accepts a connections on the port :8080")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	log.Println("App successfully started!")
	<-quit
	log.Println("received os.Interrupt, exiting...")
}
