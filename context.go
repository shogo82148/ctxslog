package ctxslog

import (
	"context"
	"log/slog"
)

type ctxKey struct{ name string }

func (key ctxKey) String() string {
	return key.name
}

var key = &ctxKey{"ctxslog"}

const nAttrsInline = 5

type mergedAttrs struct {
	parent *mergedAttrs
	args   []any
	attrs  []slog.Attr
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

func (attrs *mergedAttrs) addToRecord(record slog.Record) {
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
	attrs.addToRecord(record)
	return parent(ctx, record)
}
