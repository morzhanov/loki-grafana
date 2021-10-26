package logger

import (
	"os"
	"time"

	"github.com/afiskon/promtail-client/promtail"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger(lokiUrl string) (*zap.Logger, error) {
	labels := "{source=go_app,job=go_app_logger}"
	cfg := promtail.ClientConfig{
		PushURL:            "http://localhost:3100/api/prom/push",
		Labels:             labels,
		BatchWait:          5 * time.Second,
		BatchEntriesNumber: 10000,
		SendLevel:          promtail.INFO,
		PrintLevel:         promtail.INFO,
	}
	loki, err := promtail.NewClientJson(cfg)
	if err != nil {
		return nil, err
	}

	consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zap.InfoLevel),
	)
	log := zap.New(core)
	log = log.WithOptions(zap.Hooks(func(entry zapcore.Entry) error {
		tstamp := time.Now().String()
		loki.Infof(`source = '%s', time = '%s', message = '%s'`, "go app", tstamp, entry.Message)
		time.Sleep(1 * time.Second)
		return nil
	}))
	return log, nil
}
