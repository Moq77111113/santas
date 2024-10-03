package handlers

import (
	"time"

	"github.com/labstack/echo/v4/middleware"

	"github.com/moq77111113/chmoly-santas/pkg/services"
)

// Setup the router with registered handlers
func Bootstrap(c *services.Container) error {

	c.Web.Use(middleware.Logger())
	c.Web.Use(middleware.Recover())

	g := c.Web.Group("")

	g.Use(
		middleware.Recover(),
		middleware.Secure(),
		middleware.RequestID(),
		middleware.Gzip(),
		middleware.TimeoutWithConfig(middleware.TimeoutConfig{
			Timeout: (time.Second * 10),
		}),
	// middleware.CSRFWithConfig(middleware.CSRFConfig{
	// 	TokenLookup: "form:_csrf"})
	)

	for _, h := range GetHandlers() {
		if err := h.Init(c); err != nil {
			return err
		}
		h.Routes(g)
	}

	return nil

}
