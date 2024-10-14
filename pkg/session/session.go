package session

import (
	"errors"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
)

// ErrStoreNotFound indicates that the session store was not present in the context
var (
	ErrStoreNotFound = errors.New("session store not found")

	SessionKey = "session"
)

// Get returns a session
func Get(ctx echo.Context, name string) (*sessions.Session, error) {
	s := ctx.Get(SessionKey)
	if s == nil {
		return nil, ErrStoreNotFound
	}
	store := s.(sessions.Store)
	return store.Get(ctx.Request(), name)
}

// Store sets the session storage in the context
func Store(ctx echo.Context, store sessions.Store) {
	ctx.Set(SessionKey, store)
}
