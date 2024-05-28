package ctx

import (
	"context"

	"github.com/VCell/AgiWhisper/module/log"
	"github.com/labstack/echo"
)

func GetCtxFromEcho(ctx echo.Context) context.Context {
	traceID := ctx.Get(log.KEY_TRACE_ID)
	result := context.WithValue(ctx.Request().Context(), log.TraceIDContextKey{}, traceID)
	return result
}
