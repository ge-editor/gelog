// keychord_debug.go
//go:build debug
// +build debug

package gelog

import "log/slog"

func Info(msg string, args ...any) {
	slog.Info(msg, args...)
}
