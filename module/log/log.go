package log

import (
	"log/slog"
	"os"
)

func init() {
	var h slog.Handler
	h = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	})
	h = TraceIDHandler{h}
	slog.SetDefault(slog.New(h))
}
