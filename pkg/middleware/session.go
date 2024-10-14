package middleware

import (
	"net/http"

	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"github.com/moq77111113/chmoly-santas/pkg/session"
)

func Session(store sessions.Store) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			defer context.Clear(c.Request())
			session.Store(c, store)
			return next(c)
		}
	}
}

func CookieStore(key string) sessions.Store {
	s := sessions.NewCookieStore([]byte(key))
	s.Options = &sessions.Options{
		HttpOnly: true,
		MaxAge:   24 * 60 * 60,
		SameSite: http.SameSiteStrictMode,
	}
	return s
}
