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
	redirectURLFunc func(c echo.Context) string,
) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			u := c.Get(AuthSessionKey)
			if u == nil {
				redirectTo := redirectURLFunc(c)
				return c.Redirect(302, redirectTo)
			}
			return next(c)
		}
	}
}

func WithoutAuthentication(redirectURLFunc func(c echo.Context) string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			u := c.Get(AuthSessionKey)
			if u != nil {
				redirectTo := redirectURLFunc(c)
				return c.Redirect(302, redirectTo)
			}
			return next(c)
		}
	}
}
