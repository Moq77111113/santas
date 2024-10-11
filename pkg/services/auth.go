package services

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/moq77111113/chmoly-santas/config"
	"github.com/moq77111113/chmoly-santas/ent"
	"github.com/moq77111113/chmoly-santas/pkg/session"
)

const (
	authSessionName = "ua"

	authSessionKey = "uid"
)

type AuthClient struct {
	config *config.Config
	orm    *ent.Client
}

func NewAuthClient(config *config.Config, orm *ent.Client) *AuthClient {
	return &AuthClient{config: config, orm: orm}
}

func (a *AuthClient) GetAuthenticatedUserId(ctx echo.Context) (int, error) {
	s, err := session.Get(ctx, authSessionName)
	if err != nil {
		return -1, err
	}

	if s.Values[authSessionKey] == nil {
		return -1, fmt.Errorf("no authenticated user")
	}
	return s.Values[authSessionKey].(int), nil
}
func (a *AuthClient) GetAuthenticatedUser(ctx echo.Context) (*ent.Member, error) {
	if uid, err := a.GetAuthenticatedUserId(ctx); err == nil {
		return a.orm.Member.Get(ctx.Request().Context(), uid)
	}
	return nil, fmt.Errorf("no authenticated user")
}

func (a *AuthClient) SetAuthenticatedUser(ctx echo.Context, m *ent.Member) error {
	s, err := session.Get(ctx, authSessionName)
	if err != nil {

		return err
	}

	s.Values[authSessionKey] = m.ID

	return s.Save(ctx.Request(), ctx.Response())
}
