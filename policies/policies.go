package policies

import (
	"context"
	"keycloak-tools/access"

	"github.com/Nerzal/gocloak/v7"
)

type PolicyService struct {
	client gocloak.GoCloak
	ctx    context.Context
	token  string
}

func New(ctx *access.KeycloakContext) *PolicyService {
	return &PolicyService{
		client: ctx.Client,
		ctx:    ctx.Ctx,
		token:  ctx.Token.AccessToken,
	}
}

func (s *PolicyService) CreatePolicy(clientId string, policy *gocloak.PolicyRepresentation) error {
	_, err := s.client.CreatePolicy(s.ctx, s.token, "products", clientId, *policy)
	return err
}
