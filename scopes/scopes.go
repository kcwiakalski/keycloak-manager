package scopes

import (
	"context"
	"keycloak-tools/access"

	"github.com/Nerzal/gocloak/v7"
)

type ScopeService struct {
	client gocloak.GoCloak
	ctx    context.Context
	token  string
}

func New(ctx *access.KeycloakContext) *ScopeService {
	return &ScopeService{
		client: ctx.Client,
		ctx:    ctx.Ctx,
		token:  ctx.Token.AccessToken,
	}
}

func (s *ScopeService) AddScope(clientId string, scope *gocloak.ScopeRepresentation) error {
	_, err := s.client.CreateScope(s.ctx, s.token, "products", clientId, *scope)
	if err != nil {
		return err
	}
	return nil
}
