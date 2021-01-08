package resources

import (
	"context"
	"keycloak-tools/access"

	"github.com/Nerzal/gocloak/v7"
)

type ResourceService struct {
	client gocloak.GoCloak
	ctx    context.Context
	token  string
}

func New(ctx *access.KeycloakContext) *ResourceService {
	return &ResourceService{
		client: ctx.Client,
		ctx:    ctx.Ctx,
		token:  ctx.Token.AccessToken,
	}
}

func (s *ResourceService) AddResource(clientId string, resource gocloak.ResourceRepresentation) error {
	_, err := s.client.CreateResource(s.ctx, s.token, "products", clientId, resource)
	return err
}
