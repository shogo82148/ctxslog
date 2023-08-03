[![test](https://github.com/shogo82148/ctxslog/actions/workflows/test.yml/badge.svg)](https://github.com/shogo82148/ctxslog/actions/workflows/test.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/shogo82148/ctxslog.svg)](https://pkg.go.dev/github.com/shogo82148/ctxslog)

# ctxslog

```go
handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{ReplaceAttr: replace})
slog.SetDefault(slog.New(ctxslog.New(handler)))

ctx := context.Background()
// associate the key my_context and the value foo-bar with the context.
ctx = ctxslog.With(ctx, "my_context", "foo-bar")

slog.InfoContext(ctx, "hello", "count", 42)
slog.InfoContext(ctx, "world")
// Output:
// time=2023-08-03T18:10:20.424+09:00 level=INFO msg=hello count=42 my_context=foo-bar
// time=2023-08-03T18:10:20.424+09:00 level=INFO msg=world my_context=foo-bar
```
