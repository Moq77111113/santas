package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/moq77111113/chmoly-santas/ent"
	"github.com/moq77111113/chmoly-santas/pkg/middleware"
)

// Helper function to get the current user from the context
func getCurrentUser(ctx echo.Context) (*ent.Member, error) {
	me := ctx.Get(middleware.AuthSessionKey)
	if me == nil {
		return nil, echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}
	return me.(*ent.Member), nil
}

func checkParamMw(path string, types ...string) echo.MiddlewareFunc {
	paramType := "int"
	if len(types) > 0 {
		paramType = types[0]
	}
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			param := ctx.Param(path)
			if param == "" {
				return echo.NewHTTPError(http.StatusBadRequest, "invalid "+path)
			}

			switch paramType {
			case "int":
				id, err := strconv.Atoi(param)
				if err != nil {
					return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("%v must be an integer", path))
				}
				ctx.Set(path, id)
			default: // unhandled other

			}

			return next(ctx)
		}
	}
}
