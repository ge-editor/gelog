package gelog

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"strings"
)

type PlainHandler struct{ w io.Writer }

func (h *PlainHandler) Enabled(_ context.Context, _ slog.Level) bool { return true }

func (h *PlainHandler) Handle(_ context.Context, r slog.Record) error {
	b := &strings.Builder{}

	// level を先頭に出す（例: INFO, ERROR, DEBUG）
	fmt.Fprintf(b, "[%s] %s", strings.ToUpper(r.Level.String()), r.Message)

	// key=value を追加
	r.Attrs(func(a slog.Attr) bool {
		b.WriteString(" ")
		b.WriteString(a.Key)
		b.WriteString("=")
		b.WriteString(fmt.Sprint(a.Value))
		return true
	})

	_, err := fmt.Fprintln(h.w, b.String())
	return err
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
