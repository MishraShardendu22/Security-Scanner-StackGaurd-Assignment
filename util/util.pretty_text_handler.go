package util

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"time"
)

type PrettyTextHandler struct {
	h slog.Handler
}

func NewPrettyTextHandler(w io.Writer, opts *slog.HandlerOptions) slog.Handler {
	return &PrettyTextHandler{h: slog.NewTextHandler(w, opts)}
}

func (h *PrettyTextHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.h.Enabled(ctx, level)
}

func (h *PrettyTextHandler) Handle(ctx context.Context, r slog.Record) error {
	t := r.Time.Format(time.RFC3339)
	lvl := r.Level.String()
	msg := r.Message
	fmt.Fprintf(os.Stdout, "\n--- LOG ENTRY ---\nTime: %s\nLevel: %s\nMessage: %s\n", t, lvl, msg)
	r.Attrs(func(a slog.Attr) bool {
		fmt.Fprintf(os.Stdout, "%s: %v\n", a.Key, a.Value)
		return true
	})
	fmt.Fprintln(os.Stdout, "------------------")
	return nil
}

func (h *PrettyTextHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &PrettyTextHandler{h: h.h.WithAttrs(attrs)}
}

func (h *PrettyTextHandler) WithGroup(name string) slog.Handler {
	return &PrettyTextHandler{h: h.h.WithGroup(name)}
}
