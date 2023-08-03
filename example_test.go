package ctxslog_test

import (
	"context"
	"log/slog"
	"os"

	"github.com/shogo82148/ctxslog"
)

func Example() {
	// it's for testing.
	replace := func(groups []string, a slog.Attr) slog.Attr {
		// Remove time.
		if a.Key == slog.TimeKey && len(groups) == 0 {
			return slog.Attr{}
		}
		return a
	}

	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{ReplaceAttr: replace})
	slog.SetDefault(slog.New(ctxslog.New(handler)))

	ctx := context.Background()
	ctx = ctxslog.With(ctx, "my_context", "foo-bar")

	slog.InfoContext(ctx, "hello", "count", 42)
	slog.InfoContext(ctx, "world")

	// Output:
	// level=INFO msg=hello count=42 my_context=foo-bar
	// level=INFO msg=world my_context=foo-bar
}
