package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/moq77111113/chmoly-santas/pkg/services"
)

type Handler interface {
	Routes(g *echo.Group)

	Init(*services.Container) error
}

var handlers []Handler

func Register(h Handler) {
	handlers = append(handlers, h)
}

func GetHandlers() []Handler {
	return handlers
}
