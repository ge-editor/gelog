package gelog

import (
	"log/slog"
	"runtime"
)

func InitLogger(logFilePath string, maxSizeMB int) error {

	writer, err := NewSizeLimitedWriter(logFilePath, maxSizeMB)
	if err != nil {
		return err
	}

	fileHandler := &PlainHandler{w: writer}

	var handler slog.Handler
	if runtime.GOOS == "windows" {
		handler = fileHandler
	} else {
		handler = NewMultiHandler(
			fileHandler,
		)
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)

	return nil
}
