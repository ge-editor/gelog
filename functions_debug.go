// keychord_debug.go
//go:build debug
// +build debug

package gelog

import (
	"context"
	"fmt"
	"log/slog"
	"runtime"
)

func Info(msg string, args ...any) {
	slog.Info(msg, args...)
}

func Warn(msg string, args ...any) {
	slog.Warn(msg, args...)
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

// 呼び出し元の情報を整形して返すヘルパー関数
func CallerInfo() string {
	// skip: 0=Caller自身, 1=callerInfoの呼び出し元, 2=さらにその呼び出し元
	pc, file, line, ok := runtime.Caller(2)
	if !ok {
		return "unknown:0"
	}
	// パス全体が長い場合はファイル名だけに絞ることも可能
	// shortFile := filepath.Base(file)
	return fmt.Sprintf("%s:%d (%s)", file, line, runtime.FuncForPC(pc).Name())
}
