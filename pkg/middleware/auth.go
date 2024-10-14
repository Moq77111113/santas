package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/moq77111113/chmoly-santas/pkg/services"
)

const (
	AuthSessionKey = "me"
)

func LoadUser(auth *services.AuthClient) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			u, err := auth.GetAuthenticatedUser(c)
			switch err.(type) {
			case nil:
				c.Set(AuthSessionKey, u)
			default:
				c.Logger().Warn(err)
			}
			return next(c)
		}
	}
}

func WithAuthentication(
	onError func(c echo.Context) error,
) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			u := c.Get(AuthSessionKey)
			if u == nil {
				return onError(c)
			}
			return next(c)
		}
	}
}

func WithoutAuthentication(onError func(c echo.Context) error) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			u := c.Get(AuthSessionKey)
			if u != nil {
				return onError(c)
			}
			return next(c)
		}
	}
}
