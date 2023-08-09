package ctxslog

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"runtime"
	"testing"
)

var optsRemoveTime = &slog.HandlerOptions{
	ReplaceAttr: removeTime,
}

// removeTime removes time attribute for stable testing.
func removeTime(groups []string, a slog.Attr) slog.Attr {
	// Remove time.
	if a.Key == slog.TimeKey && len(groups) == 0 {
		return slog.Attr{}
	}
	return a
}

func TestString(t *testing.T) {
	ctx := context.Background()
	ctx = With(ctx, "hello", 1)
	ctx = With(ctx, "world", 2)
	ctx = With(ctx)
	if got, want := fmt.Sprint(ctx), "context.Background.WithValue(type *ctxslog.ctxKey, val hello=1).WithValue(type *ctxslog.ctxKey, val world=2)"; got != want {
		t.Errorf("unexpected output: got %q, want %q", got, want)
	}
}

func TestWith(t *testing.T) {
	buf := &bytes.Buffer{}
	parent := slog.NewTextHandler(buf, optsRemoveTime)
	child := New(parent)
	logger := slog.New(child)

	ctx := context.Background()
	ctx = With(ctx, "hello", 1)
	ctx = With(ctx, "world", 2)
	ctx = With(ctx)
	logger.InfoContext(ctx, "hello", "count", 42)
	if buf.String() != "level=INFO msg=hello count=42 hello=1 world=2\n" {
		t.Errorf("unexpected output: %q", buf.String())
	}
}

func TestWithAttrs(t *testing.T) {
	buf := &bytes.Buffer{}
	parent := slog.NewTextHandler(buf, optsRemoveTime)
	child := New(parent)
	logger := slog.New(child)

	ctx := context.Background()
	ctx = WithAttrs(ctx)
	ctx = WithAttrs(ctx, slog.Int("hello", 1), slog.Int("world", 2))
	logger.InfoContext(ctx, "hello", "count", 42)
	if buf.String() != "level=INFO msg=hello count=42 hello=1 world=2\n" {
		t.Errorf("unexpected output: %q", buf.String())
	}
}

func TestHandlerWithAttrs(t *testing.T) {
	buf := &bytes.Buffer{}
	parent := slog.NewTextHandler(buf, optsRemoveTime)
	child := New(parent).WithAttrs([]slog.Attr{slog.Int("hello", 1)})
	logger := slog.New(child)

	ctx := context.Background()
	ctx = WithAttrs(ctx, slog.Int("world", 2))
	logger.InfoContext(ctx, "hello", "count", 42)
	if buf.String() != "level=INFO msg=hello hello=1 count=42 world=2\n" {
		t.Errorf("unexpected output: %q", buf.String())
	}
}

func TestHandlerWithGroup(t *testing.T) {
	buf := &bytes.Buffer{}
	parent := slog.NewTextHandler(buf, optsRemoveTime)
	child := New(parent).WithGroup("my_group")
	logger := slog.New(child)

	ctx := context.Background()
	ctx = WithAttrs(ctx, slog.Int("hello", 1), slog.Int("world", 2))
	logger.InfoContext(ctx, "hello", "count", 42)
	if buf.String() != "level=INFO msg=hello my_group.count=42 my_group.hello=1 my_group.world=2\n" {
		t.Errorf("unexpected output: %q", buf.String())
	}
}

func TestWithLevel(t *testing.T) {
	buf := &bytes.Buffer{}
	parent := slog.NewTextHandler(buf, optsRemoveTime)
	child := New(parent)
	logger := slog.New(child)

	ctx := context.Background()
	logger.InfoContext(ctx, "this log should be logged")
	logger.DebugContext(ctx, "this log should not be logged")

	ctx = WithLevel(ctx, slog.LevelDebug) // enable debug level
	logger.DebugContext(ctx, "this log should be logged in this context")

	got := buf.String()
	want := `level=INFO msg="this log should be logged"
level=DEBUG msg="this log should be logged in this context"
`
	if got != want {
		t.Errorf("unexpected output: got %q, want %q", got, want)
	}
}

func BenchmarkWith(b *testing.B) {
	ctx := context.Background()
	for i := 0; i < b.N; i++ {
		ctx := ctx
		for j := 0; j < 128; j++ {
			ctx = With(ctx, fmt.Sprintf("hello%d", j), j)
		}
		runtime.KeepAlive(ctx)
	}
}

func BenchmarkWithAttrs(b *testing.B) {
	ctx := context.Background()
	for i := 0; i < b.N; i++ {
		ctx := ctx
		for j := 0; j < 128; j++ {
			ctx = WithAttrs(ctx, slog.Int(fmt.Sprintf("hello%d", j), j))
		}
		runtime.KeepAlive(ctx)
	}
}

func BenchmarkLog(b *testing.B) {
	ctx := context.Background()
	for j := 0; j < 128; j++ {
		ctx = WithAttrs(ctx, slog.Int(fmt.Sprintf("hello%d", j), j))
	}

	parent := slog.NewTextHandler(io.Discard, nil)
	child := New(parent)
	logger := slog.New(child)
	for i := 0; i < b.N; i++ {
		logger.InfoContext(ctx, "hello")
	}
}
