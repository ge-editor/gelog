package gelog

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"runtime"
	"strings"
)

func callerOutsideGelog() (string, int, bool) {
	_, file, line, ok := runtime.Caller(6) // 固定
	return file, line, ok
	/*
		for i := 3; i < 15; i++ {
			_, file, line, ok := runtime.Caller(i)
			if !ok {
				return "", 0, false
			}
			if !strings.Contains(file, "/gelog/") {
				return file, line, true
			}
		}
		return "", 0, false
	*/
}

func shortFile(path string) string {
	if i := strings.LastIndex(path, "/"); i >= 0 {
		return path[i+1:]
	}
	return path
}

type PlainHandler struct {
	w      io.Writer
	indent string
}

func NewPlainHandler(w io.Writer) *PlainHandler {
	return &PlainHandler{
		w:      w,
		indent: "    ",
	}
}

func (h *PlainHandler) Enabled(_ context.Context, _ slog.Level) bool {
	return true
}

const (
	LevelKey slog.Level = 2
)

func levelLabel(l slog.Level) string {
	switch l {
	case slog.LevelDebug:
		return "🐛[DEBUG]"
	case slog.LevelInfo:
		return "ℹ️[INFO]"
	case slog.LevelWarn:
		return "⚠️[WARN]"
	case slog.LevelError:
		return "❌[ERROR]"
	case LevelKey:
		return "🪷[KEY]"
	default:
		return "❓[" + strings.ToUpper(l.String()) + "]"
	}
}

func (h *PlainHandler) Handle(_ context.Context, r slog.Record) error {
	b := &strings.Builder{}

	/* 	fmt.Fprintf(b, "[%s] %s",
	   		strings.ToUpper(r.Level.String()),
	   		r.Message,
	   	)
	*/
	fmt.Fprintf(b, "%s %s",
		levelLabel(r.Level),
		r.Message,
	)

	first := true
	r.Attrs(func(a slog.Attr) bool {
		if first {
			// b.WriteString(" | ")
			b.WriteString(" ")
			first = false
		} else {
			b.WriteString(" ")
		}
		fmt.Fprintf(b, "%s=%v", a.Key, formatValue(a.Value))
		return true
	})

	// ★ ここを追加
	if file, line, ok := callerOutsideGelog(); ok {
		fmt.Fprintf(b, " (%s:%d)", shortFile(file), line)
	}

	_, err := fmt.Fprintln(h.w, b.String())
	return err
}

func formatValue(v slog.Value) any {
	switch v.Kind() {
	case slog.KindString:
		return fmt.Sprintf("%q", v.String())
	default:
		return v.Any()
	}
}

func (h *PlainHandler) WithAttrs(_ []slog.Attr) slog.Handler { return h }
func (h *PlainHandler) WithGroup(_ string) slog.Handler      { return h }

// --- 複数出力用ハンドラ ---
type MultiHandler struct{ handlers []slog.Handler }

func (h *MultiHandler) Enabled(ctx context.Context, level slog.Level) bool {
	for _, handler := range h.handlers {
		if handler.Enabled(ctx, level) {
			return true
		}
	}
	return false
}

func (h *MultiHandler) Handle(ctx context.Context, record slog.Record) error {
	var err error
	for _, handler := range h.handlers {
		if e := handler.Handle(ctx, record); e != nil {
			err = e
		}
	}
	return err
}

func (h *MultiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newHandlers := make([]slog.Handler, len(h.handlers))
	for i, handler := range h.handlers {
		newHandlers[i] = handler.WithAttrs(attrs)
	}
	return &MultiHandler{handlers: newHandlers}
}

func (h *MultiHandler) WithGroup(name string) slog.Handler {
	newHandlers := make([]slog.Handler, len(h.handlers))
	for i, handler := range h.handlers {
		newHandlers[i] = handler.WithGroup(name)
	}
	return &MultiHandler{handlers: newHandlers}
}

func NewMultiHandler(handlers ...slog.Handler) slog.Handler {
	return &MultiHandler{handlers: handlers}
}
