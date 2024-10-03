package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/moq77111113/chmoly-santas/pkg/services"
)

// A Handler handle http routes
type Handler interface {
	// Routes register routes to the router
	Routes(g *echo.Group)

	// Init is called once when the router bootstrap
	Init(*services.Container) error
}

var handlers []Handler

func Register(h Handler) {
	handlers = append(handlers, h)
}

func GetHandlers() []Handler {
	return handlers
}
