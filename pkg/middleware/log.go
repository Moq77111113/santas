package middleware

import (
	"fmt"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/color"

	"github.com/moq77111113/chmoly-santas/pkg/log"
)

func Logger() echo.MiddlewareFunc {

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			id := c.Response().Header().Get(echo.HeaderXRequestID)
			logger := log.Ctx(c).With("request_id", id)

			log.Set(c, logger)
			return next(c)
		}
	}
}

func LogRequest() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			res := c.Response()

			start := time.Now()
			if err := next(c); err != nil {
				c.Error(err)
			}
			duration := time.Since(start)

			sub := log.Ctx(c).With(
				"ip", c.RealIP(),
				"host", req.Host,
				"referer", req.Referer(),
				"status", res.Status,
				"latency", duration,
			)

			msg := fmt.Sprintf("%s %s", req.Method, req.URL.RequestURI())

			if res.Status >= 500 {
				sub.Error(color.Red(msg))
			} else if res.Status >= 400 {
				sub.Warn(color.Yellow(msg))
			} else {
				sub.Info(color.Green(msg))
			}

			return nil
		}
	}
}
