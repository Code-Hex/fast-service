package logger

import (
	"context"
	"fmt"
	"strings"

	"github.com/blendle/zapdriver"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type loggerKey struct{}

var contextKey = &loggerKey{}

// Extract takes the call-scoped Logger from context.
func Extract(ctx context.Context) *zap.Logger {
	l, ok := ctx.Value(contextKey).(*zap.Logger)
	if !ok || l == nil {
		return NewDiscard()
	}
	return l
}

// ToContext adds the zap.Logger into context.
func ToContext(ctx context.Context, l *zap.Logger) context.Context {
	return context.WithValue(ctx, contextKey, l)
}

// Debug logs a message at DebugLevel.
func Debug(ctx context.Context, msg string, fields ...zap.Field) {
	Extract(ctx).Debug(msg, fields...)
}

// Info logs a message at InfoLevel.
func Info(ctx context.Context, msg string, fields ...zap.Field) {
	Extract(ctx).Info(msg, fields...)
}

// Warn logs a message at WarnLevel.
func Warn(ctx context.Context, msg string, fields ...zap.Field) {
	Extract(ctx).Warn(msg, fields...)
}

// Error logs a message at ErrorLevel.
func Error(ctx context.Context, msg string, fields ...zap.Field) {
	Extract(ctx).Error(msg, fields...)
}

// New creates a new zap logger with the given log level.
func New(level string) (*zap.Logger, error) {
	l, err := logLevel(level)
	if err != nil {
		return nil, err
	}

	config := zapdriver.NewProductionConfig()
	config.Level = zap.NewAtomicLevelAt(l)
	config.DisableStacktrace = true
	return config.Build()
}

// NewDiscard creates logger which output to ioutil.Discard.
// This can be used for testing.
func NewDiscard() *zap.Logger {
	return zap.NewNop()
}

func logLevel(level string) (zapcore.Level, error) {
	level = strings.ToUpper(level)
	var l zapcore.Level
	switch level {
	case "DEBUG":
		l = zapcore.DebugLevel
	case "INFO":
		l = zapcore.InfoLevel
	case "ERROR":
		l = zapcore.ErrorLevel
	default:
		return l, fmt.Errorf("invalid loglevel: %s", level)
	}
	return l, nil
}
