package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/moq77111113/chmoly-santas/config"
	"github.com/moq77111113/chmoly-santas/pkg/services"
	"github.com/moq77111113/chmoly-santas/ui"
)

type (
	UI struct {
		conf *config.Config
	}
)

func init() {
	Register(new(UI))
}

func (u *UI) Init(c *services.Container) error {
	u.conf = c.Config
	return nil
}

func (u *UI) Routes(g *echo.Group) {

	g.GET("/*", echo.StaticDirectoryHandler(ui.DistDirFs, false))

}
