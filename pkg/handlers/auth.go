package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/moq77111113/chmoly-santas/pkg/form"
	"github.com/moq77111113/chmoly-santas/pkg/middleware"
	"github.com/moq77111113/chmoly-santas/pkg/services"
)

type (
	Auth struct {
		Auth  *services.AuthClient
		Group *services.GroupRepo
	}

	register struct {
		Name string `form:"name" validate:"required"`
		form.Form
	}
)

func init() {
	Register(new(Auth))
}

func (h *Auth) Init(c *services.Container) error {
	h.Group = c.Repositories.Group
	h.Auth = c.Auth
	return nil
}

func (h *Auth) Routes(g *echo.Group) {

	auth := g.Group("/auth")

	auth.POST("/register", h.Register)
	wAuth := auth.Group("", middleware.WithAuthentication(func(c echo.Context) error {
		return c.JSON(401, "Unauthorized")
	}))
	wAuth.GET("/me", h.Me)

}

func (h *Auth) IsApi() bool {
	return true
}

func (h *Auth) Me(c echo.Context) error {
	me := c.Get(middleware.AuthSessionKey)

	if me == nil {
		return c.JSON(401, "Unauthorized")
	}
	return c.JSON(200, me)
}

func (h *Auth) Register(ctx echo.Context) error {

	var form register
	if err := form.BindAndValidate(ctx, &form); err != nil {
		ctx.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, form.Errors())
	}

	mm, err := h.Group.CreateMember(ctx.Request().Context(), form.Name)

	if err != nil {
		ctx.Logger().Error(err)
		return err
	}

	err = h.Auth.SetAuthenticatedUser(ctx, mm)

	if err != nil {
		ctx.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to register")
	}

	return ctx.JSON(http.StatusOK, nil)
}
