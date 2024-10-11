package handlers

import (
	"net/http"

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
