package ctxslog

import (
	"context"
	"log/slog"
)

var _ slog.Handler = (*wrapper)(nil)

type wrapper struct {
	handler func(ctx context.Context, parent func(ctx context.Context, record slog.Record) error, record slog.Record) error
	parent  slog.Handler
}

func (w *wrapper) Handle(ctx context.Context, record slog.Record) error {
	return w.handler(ctx, w.parent.Handle, record)
}

func (w *wrapper) Enabled(ctx context.Context, level slog.Level) bool {
	return w.parent.Enabled(ctx, level)
}

func (w *wrapper) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &wrapper{
		parent: w.parent.WithAttrs(attrs),
	}
}

func (w *wrapper) WithGroup(name string) slog.Handler {
	return &wrapper{
		parent: w.parent.WithGroup(name),
	}
}
