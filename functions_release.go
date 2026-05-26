// keychord_release.go
//go:build !debug
// +build !debug

package gelog

func Info(msg string, args ...any)  {}
func Error(msg string, args ...any) {}
func Debug(msg string, args ...any) {}
func Key(msg string, args ...any)   {}
