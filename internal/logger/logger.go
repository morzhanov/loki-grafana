package logger

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger() (*zap.Logger, error) {
	baseLogsPath := "./logs/"
	if err := os.MkdirAll(baseLogsPath, 0777); err != nil {
		return nil, err
	}

	// info
	infoFilePath := fmt.Sprintf("%sinfo.log", baseLogsPath)
	_, err := os.Create(infoFilePath)
	if err != nil {
		return nil, err
	}
	infoF, err := os.OpenFile(infoFilePath, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return nil, err
	}

	// error
	errFilePath := fmt.Sprintf("%serror.log", baseLogsPath)
	_, err = os.Create(errFilePath)
	if err != nil {
		return nil, err
	}
	errF, err := os.OpenFile(errFilePath, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return nil, err
	}

	fileEncoder := zapcore.NewJSONEncoder(zap.NewDevelopmentEncoderConfig())
	consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, zapcore.AddSync(infoF), zap.InfoLevel),
		zapcore.NewCore(fileEncoder, zapcore.AddSync(errF), zap.ErrorLevel),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zap.InfoLevel),
	)
	return zap.New(core), nil
}
