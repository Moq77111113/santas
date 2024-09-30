package apis

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/moq77111113/chmoly-santas/internal/core"
)

func bindSSEAPi(app core.App, rg *echo.Group) {
	api := connectApi{app: app}

	rg.GET("/sse", api.join)

}

type connectApi struct {
	app core.App
}

func (api *connectApi) join(c echo.Context) error {

	ch := make(chan string)

	api.app.Notifier.Add <- ch
	defer func() {
		api.app.Notifier.Remove <- ch
		close(ch)
	}()

	c.Response().Header().Set("Content-Type", "text/event-stream")
	c.Response().Writer.WriteHeader(http.StatusOK)

	if _, err := c.Response().Writer.Write([]byte(": ping\n\n")); err != nil {
		return err
	}
	c.Response().Flush()
	ticker := time.NewTicker(30 * time.Second)

	defer ticker.Stop()

	for {
		select {
		case msg := <-ch:
			if _, err := c.Response().Writer.Write([]byte("data: " + msg + "\n\n")); err != nil {
				return err
			}
			c.Response().Flush()
		case <-ticker.C:
			if _, err := c.Response().Writer.Write([]byte(": ping\n\n")); err != nil {
				return err
			}
			c.Response().Flush()
		case <-c.Request().Context().Done():
			return nil
		}
	}

}
