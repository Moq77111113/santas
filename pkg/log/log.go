package log

import (
	"log/slog"

	"github.com/labstack/echo/v4"
)

func Set(ctx echo.Context, logger *slog.Logger) {
	ctx.Set("logger", logger)
}

func Ctx(ctx echo.Context) *slog.Logger {
	if logger, ok := ctx.Get("logger").(*slog.Logger); ok {
		return logger
	}
	return slog.Default()
}
