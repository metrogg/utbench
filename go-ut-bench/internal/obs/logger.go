package obs

import (
	"context"
	"log/slog"
	"os"
)

type Logger struct {
	base *slog.Logger
}

func NewLogger(verbose bool) *Logger {
	level := slog.LevelInfo
	if verbose {
		level = slog.LevelDebug
	}
	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: level})
	return &Logger{base: slog.New(handler)}
}

func (l *Logger) Info(msg string, attrs ...any) {
	l.base.Info(msg, attrs...)
}

func (l *Logger) Debug(msg string, attrs ...any) {
	l.base.Debug(msg, attrs...)
}

func (l *Logger) Warn(msg string, attrs ...any) {
	l.base.Warn(msg, attrs...)
}

func (l *Logger) Error(msg string, attrs ...any) {
	l.base.Error(msg, attrs...)
}

func (l *Logger) With(attrs ...any) *Logger {
	return &Logger{base: l.base.With(attrs...)}
}

func (l *Logger) WithContext(_ context.Context) *Logger {
	return l
}
