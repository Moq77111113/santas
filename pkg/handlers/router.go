package handlers

import (
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	echoMw "github.com/labstack/echo/v4/middleware"

	"github.com/moq77111113/chmoly-santas/pkg/middleware"
	"github.com/moq77111113/chmoly-santas/pkg/services"
)

// Setup the router with registered handlers
func Bootstrap(c *services.Container) error {

	c.Web.Pre(echoMw.RemoveTrailingSlashWithConfig(echoMw.TrailingSlashConfig{
		Skipper: func(c echo.Context) bool {
			return !strings.HasPrefix(c.Request().URL.Path, "/api")
		},
		RedirectCode: http.StatusMovedPermanently,
	}))

	g := c.Web.Group("")

	g.Use(
		echoMw.Recover(),
		echoMw.Secure(),
		echoMw.RequestID(),
	)

	a := g.Group("/api")

	a.Use(
		middleware.Logger(),
		middleware.LogRequest(),
		echoMw.Gzip(),
		middleware.Session(sessions.NewCookieStore([]byte(c.Config.App.EncryptionKey))),
		middleware.LoadUser(c.Auth),
		echoMw.TimeoutWithConfig(echoMw.TimeoutConfig{
			Timeout: (time.Second * 10),
			Skipper: func(c echo.Context) bool {
				return strings.Contains(c.Request().URL.Path, "events")
			},
		}),
	)

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
