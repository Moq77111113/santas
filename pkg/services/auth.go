package services

import (
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

func (a *AuthClient) GetAuthenticatedUser(ctx echo.Context) (*ent.Member, error) {
	s, err := session.Get(ctx, authSessionName)
	if err != nil {
		return nil, err
	}

	if s.Values[authSessionKey] == nil {
		return nil, err
	}

	m, err := a.orm.Member.Get(ctx.Request().Context(), s.Values[authSessionKey].(int))

	if err != nil {
		return nil, err
	}

	return m, nil
}

func (a *AuthClient) SetAuthenticatedUser(ctx echo.Context, m *ent.Member) error {
	s, err := session.Get(ctx, authSessionName)
	if err != nil {
		return err
	}

	s.Values[authSessionKey] = m.ID

	return s.Save(ctx.Request(), ctx.Response())
}
