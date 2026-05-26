package gelog

import (
	"os"
	"sync"
)

type SizeLimitedWriter struct {
	file    *os.File
	maxSize int64
	mu      sync.Mutex
}

func NewSizeLimitedWriter(path string, maxSizeMB int) (*SizeLimitedWriter, error) {
	f, err := os.OpenFile(
		path,
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0644,
	)
	if err != nil {
		return nil, err
	}

	return &SizeLimitedWriter{
		file:    f,
		maxSize: int64(maxSizeMB) * 1024 * 1024,
	}, nil
}

func (w *SizeLimitedWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	info, err := w.file.Stat()
	if err != nil {
		return 0, err
	}

	// サイズ超過ならリセット
	if info.Size()+int64(len(p)) > w.maxSize {
		if err := w.file.Truncate(0); err != nil {
			return 0, err
		}
		if _, err := w.file.Seek(0, 0); err != nil {
			return 0, err
		}
	}

	return w.file.Write(p)
}
