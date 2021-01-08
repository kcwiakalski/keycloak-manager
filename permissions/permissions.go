package permissions

import (
	"context"
	"keycloak-tools/access"

	"github.com/Nerzal/gocloak/v7"
)

type PermissionService struct {
	client gocloak.GoCloak
	ctx    context.Context
	token  string
}

func New(ctx *access.KeycloakContext) *PermissionService {
	return &PermissionService{
		client: ctx.Client,
		ctx:    ctx.Ctx,
		token:  ctx.Token.AccessToken,
	}
}

func (s *PermissionService) AddPermission(clientId string, permission gocloak.PermissionRepresentation) error {
	_, err := s.client.CreatePermission(s.ctx, s.token, "products", clientId, permission)
	return err
}
