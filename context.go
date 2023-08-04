package ctxslog

import (
	"bytes"
	"context"
	"log/slog"
	"strings"
	"time"
)

type ctxKey struct{ name string }

func (key *ctxKey) String() string {
	return key.name
}

var key = &ctxKey{"ctxslog"}

type mergedAttrs struct {
	parent *mergedAttrs
	args   []any
	attrs  []slog.Attr
}

var handlerOptions = &slog.HandlerOptions{
	ReplaceAttr: replaceAttr,
}

func replaceAttr(groups []string, a slog.Attr) slog.Attr {
	if len(groups) != 0 {
		return a
	}

	// Remove time, level and msg.
	if a.Key == slog.TimeKey || a.Key == slog.LevelKey || a.Key == slog.MessageKey {
		return slog.Attr{}
	}

	return a
}

func (attrs *mergedAttrs) String() string {
	buf := &bytes.Buffer{}
	handler := slog.NewTextHandler(buf, handlerOptions)
	var pc uintptr
	record := slog.NewRecord(time.Time{}, slog.LevelInfo, "msg", pc)
	if len(attrs.attrs) != 0 {
		record.AddAttrs(attrs.attrs...)
	}
	if len(attrs.args) != 0 {
		record.Add(attrs.args...)
	}
	handler.Handle(context.Background(), record)
	return strings.TrimSpace(buf.String())
}

// WithAttrs is a more efficient version of [With] that accepts only [log/slog.Attrs].
func WithAttrs(ctx context.Context, attrs ...slog.Attr) context.Context {
	if len(attrs) == 0 {
		return ctx
	}
	value := &mergedAttrs{
		parent: contextAttrs(ctx),
		attrs:  attrs,
	}
	return context.WithValue(ctx, key, value)
}

// With returns a new context with the given attributes.
// The attributes are added into the log record.
func With(ctx context.Context, args ...any) context.Context {
	if len(args) == 0 {
		return ctx
	}
	value := &mergedAttrs{
		parent: contextAttrs(ctx),
		args:   args,
	}
	return context.WithValue(ctx, key, value)
}

func contextAttrs(ctx context.Context) *mergedAttrs {
	attrs := ctx.Value(key)
	if attrs == nil {
		return nil
	}
	return attrs.(*mergedAttrs)
}

func (attrs *mergedAttrs) addToRecord(record *slog.Record) {
	if attrs == nil {
		return
	}
	if attrs.parent != nil {
		attrs.parent.addToRecord(record)
	}
	if len(attrs.attrs) != 0 {
		record.AddAttrs(attrs.attrs...)
	}
	if len(attrs.args) != 0 {
		record.Add(attrs.args...)
	}
}

// New returns a new slog.Handler that injects the attributes from the context.
func New(parent slog.Handler) slog.Handler {
	return &wrapper{
		handler: inject,
		parent:  parent,
	}
}

func inject(ctx context.Context, parent func(ctx context.Context, record slog.Record) error, record slog.Record) error {
	attrs := contextAttrs(ctx)
	newRecord := record.Clone()
	attrs.addToRecord(&newRecord)
	return parent(ctx, newRecord)
}
