package gelog

import (
	"log/slog"
	// "os"
	"runtime"

	"gopkg.in/natefinch/lumberjack.v2"
)

func InitLogger(logFilePath string) {

	fileLogger := &lumberjack.Logger{
		Filename:   logFilePath,
		MaxSize:    3,
		MaxBackups: 3,
		MaxAge:     28,
		Compress:   true,
	}

	fileHandler := slog.NewTextHandler(fileLogger, &slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: true,
	})

	var handler slog.Handler
	if runtime.GOOS == "windows" {
		handler = fileHandler
	} else {
		handler = NewMultiHandler(
			// &PlainHandler{w: os.Stdout}, // Stdout 出力は停止
			fileHandler,
		)
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)
}
