package log

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/labstack/echo"
)

const (
	KEY_TRACE_ID = "traceid"
)

type TraceIDHandler struct {
	slog.Handler
}

type TraceIDContextKey struct{}

func (h TraceIDHandler) Handle(ctx context.Context, r slog.Record) error {
	if traceID, ok := ctx.Value(TraceIDContextKey{}).(string); ok {
		r.Add(KEY_TRACE_ID, slog.StringValue(traceID))
	}

	return h.Handler.Handle(ctx, r)
}

// Middleware to generate and inject traceID
func TraceIDMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		traceID := c.Request().Header.Get("X-Trace-ID")
		if traceID == "" {
			traceID = c.QueryParam(KEY_TRACE_ID)
		}
		if traceID == "" {
			traceID = uuid.New().String()
		}

		c.Response().Header().Set("X-Trace-ID", traceID)
		c.Set(KEY_TRACE_ID, traceID)
		return next(c)
	}
}
