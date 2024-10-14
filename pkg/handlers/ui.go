package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/moq77111113/chmoly-santas/config"
	"github.com/moq77111113/chmoly-santas/pkg/services"
	"github.com/moq77111113/chmoly-santas/ui"
)

type (

	// UI serves the ui pages
	UI struct {
		conf *config.Config
	}
)

func init() {
	Register(new(UI))
}

func (h *UI) Init(c *services.Container) error {
	h.conf = c.Config
	return nil
}

func (u *UI) Routes(g *echo.Group) {

	g.Use((middleware.StaticWithConfig(middleware.StaticConfig{
		HTML5:      true,
		Filesystem: http.FS(ui.DistDirFs),
		Index:      "index.html",
	})))
}

func (u *UI) IsApi() bool {
	return false
}
