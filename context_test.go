package ctxslog

import (
	"bytes"
	"context"
	"log/slog"
	"strings"
	"testing"
)

func TestWith(t *testing.T) {
	buf := &bytes.Buffer{}
	parent := slog.NewTextHandler(buf, nil)
	child := New(parent)
	logger := slog.New(child)

	ctx := context.Background()
	ctx = With(ctx, "hello", 1)
	ctx = With(ctx, "world", 2)
	ctx = With(ctx)
	logger.InfoContext(ctx, "hello", "count", 42)
	if !strings.HasSuffix(buf.String(), " level=INFO msg=hello count=42 hello=1 world=2\n") {
		t.Errorf("unexpected output: %q", buf.String())
	}
}

func TestWithAttrs(t *testing.T) {
	buf := &bytes.Buffer{}
	parent := slog.NewTextHandler(buf, nil)
	child := New(parent)
	logger := slog.New(child)

	ctx := context.Background()
	ctx = WithAttrs(ctx)
	ctx = WithAttrs(ctx, slog.Int("hello", 1), slog.Int("world", 2))
	logger.InfoContext(ctx, "hello", "count", 42)
	if !strings.HasSuffix(buf.String(), " level=INFO msg=hello count=42 hello=1 world=2\n") {
		t.Errorf("unexpected output: %q", buf.String())
	}
}

func TestHandlerWithAttrs(t *testing.T) {
	buf := &bytes.Buffer{}
	parent := slog.NewTextHandler(buf, nil)
	child := New(parent).WithAttrs([]slog.Attr{slog.Int("hello", 1)})
	logger := slog.New(child)

	ctx := context.Background()
	ctx = WithAttrs(ctx, slog.Int("world", 2))
	logger.InfoContext(ctx, "hello", "count", 42)
	if !strings.HasSuffix(buf.String(), " level=INFO msg=hello hello=1 count=42 world=2\n") {
		t.Errorf("unexpected output: %q", buf.String())
	}
}

func TestHandlerWithGroup(t *testing.T) {
	buf := &bytes.Buffer{}
	parent := slog.NewTextHandler(buf, nil)
	child := New(parent).WithGroup("my_group")
	logger := slog.New(child)

	ctx := context.Background()
	ctx = WithAttrs(ctx, slog.Int("hello", 1), slog.Int("world", 2))
	logger.InfoContext(ctx, "hello", "count", 42)
	if !strings.HasSuffix(buf.String(), " level=INFO msg=hello my_group.count=42 my_group.hello=1 my_group.world=2\n") {
		t.Errorf("unexpected output: %q", buf.String())
	}
}
