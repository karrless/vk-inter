// Package logger provides logger
package logger

import (
	"context"

	"go.uber.org/zap"
)

type KeyString string

const LoggerKey KeyString = "logger"

// Logger interface
type Logger interface {
	Debug(msg string, fields ...zap.Field)
	Info(msg string, fields ...zap.Field)
	Fatal(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
}

// logger
type logger struct {
	logger *zap.Logger
}

// Debug message
func (l *logger) Debug(msg string, fields ...zap.Field) {
	l.logger.Debug(msg, fields...)
}

// Info message
func (l *logger) Info(msg string, fields ...zap.Field) {
	l.logger.Info(msg, fields...)
}

// Fatal message
func (l *logger) Fatal(msg string, fields ...zap.Field) {
	l.logger.Fatal(msg, fields...)
}

// Warn message
func (l *logger) Warn(msg string, fields ...zap.Field) {
	l.logger.Warn(msg, fields...)
}

// Get new logger
//
// if debug, then Logger.logger = zap.NewDevelopment()
//
// else Logger.logger = zap.NewProduction()
func New(debug bool) Logger {
	var zapLogger *zap.Logger
	var err error
	if debug {
		zapLogger, err = zap.NewDevelopment()
	} else {
		zapLogger, err = zap.NewProduction()
	}
	if err != nil {
		panic(err)
	}
	defer zapLogger.Sync()

	return &logger{logger: zapLogger}
}

func FromContext(ctx context.Context) Logger {
	return ctx.Value(LoggerKey).(Logger)
}
