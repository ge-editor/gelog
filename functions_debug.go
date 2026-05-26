// keychord_debug.go
//go:build debug
// +build debug

package gelog

import (
	"context"
	"log/slog"
)

func Info(msg string, args ...any) {
	slog.Info(msg, args...)
}

func Error(msg string, args ...any) {
	slog.Error(msg, args...)
}

func Debug(msg string, args ...any) {
	slog.Debug(msg, args...)
}

func Key(msg string, args ...any) {
	slog.Log(context.Background(), LevelKey,
		msg, args...,
	)
}
