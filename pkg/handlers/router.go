package handlers

import (
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	echoMw "github.com/labstack/echo/v4/middleware"

	"github.com/moq77111113/chmoly-santas/pkg/middleware"
	"github.com/moq77111113/chmoly-santas/pkg/services"
)

// Setup the router with registered handlers
func Bootstrap(c *services.Container) error {

	// c.Web.Use(echoMw.Logger())
	c.Web.Use(echoMw.Recover())
	c.Web.Pre(echoMw.RemoveTrailingSlashWithConfig(echoMw.TrailingSlashConfig{
		Skipper: func(c echo.Context) bool {
			return !strings.HasPrefix(c.Request().URL.Path, "/api")
		},
	}))

	g := c.Web.Group("")

	g.Use(
		echoMw.Recover(),
		echoMw.Secure(),
		echoMw.RequestID(),
		echoMw.Gzip(),
		middleware.Session(middleware.CookieStore(c.Config.App.EncryptionKey)),
		middleware.LoadUser(c.Auth),
		echoMw.TimeoutWithConfig(echoMw.TimeoutConfig{
			Timeout: (time.Second * 10),
			Skipper: func(c echo.Context) bool {
				return strings.Contains(c.Request().URL.Path, "events")
			},
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
