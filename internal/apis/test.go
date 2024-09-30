package apis

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/moq77111113/chmoly-santas/internal/core"
)

func bindTestApi(app core.App, rg *echo.Group) {
	api := testApi{app: app}

	group := rg.Group("/test")
	group.GET("", api.test)
}

type testApi struct {
	app core.App
}

var i = 0

func (api *testApi) test(c echo.Context) error {
	i++
	api.app.Notifier.Notify(fmt.Sprintf("Test %d", i))
	c.Response().Writer.WriteHeader(http.StatusOK)
	return nil
}
