package apis

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/moq77111113/chmoly-santas/internal/core"
)

func Init(app *core.App) (*echo.Echo, error) {
	e := echo.New()
	e.Debug = true

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
	},
	))

	api := e.Group("/api")
	bindSSEAPi(*app, api)
	bindTestApi(*app, api)
	bindGroupApi(*app, api)

	api.Any("/*", func(c echo.Context) error {
		return echo.ErrNotFound
	})

	return e, nil

}
