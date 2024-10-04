package handlers

import (
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/moq77111113/chmoly-santas/pkg/services"
)

// Setup the router with registered handlers
func Bootstrap(c *services.Container) error {

	c.Web.Use(middleware.Logger())
	c.Web.Use(middleware.Recover())
	c.Web.Pre(middleware.RemoveTrailingSlashWithConfig(middleware.TrailingSlashConfig{
		Skipper: func(c echo.Context) bool {
			return !strings.HasPrefix(c.Request().URL.Path, "/api")
		},
	}))

	g := c.Web.Group("")

	g.Use(
		middleware.Recover(),
		middleware.Secure(),
		middleware.RequestID(),
		middleware.Gzip(),
		middleware.TimeoutWithConfig(middleware.TimeoutConfig{
			Timeout: (time.Second * 10),
		}),
	)

	a := g.Group("/api")

	for _, h := range GetHandlers() {
		if err := h.Init(c); err != nil {
			return err
		}
		if h.IsApi() {
			h.Routes(a)
		} else {
			h.Routes(g)
		}
	}

	return nil

}
