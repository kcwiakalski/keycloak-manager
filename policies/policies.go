package policies

import (
	"context"
	"keycloak-tools/modules"

	"github.com/Nerzal/gocloak/v7"
	"github.com/rs/zerolog/log"
)

type PolicyService struct {
	client gocloak.GoCloak
	ctx    context.Context
	token  string
}

var service *PolicyService

// TODO remove groupsService := New(keycloak)
func (s *PolicyService) Apply(keycloakConfig *modules.ConfigurationContext) error {
	var finalError error
	clientId := *keycloakConfig.Client.ID
	for _, policy := range keycloakConfig.Config.ClientConfig.Policies {
		err := service.CreatePolicy(clientId, &policy.PolicySpec)
		if err != nil {
			finalError = err
		}
	}
	return finalError
}

func (s *PolicyService) Order() int {
	return 3
}

func init() {
	ctx := modules.Keycloak
	service = &PolicyService{
		client: ctx.Client,
		ctx:    ctx.Ctx,
		token:  ctx.Token.AccessToken,
	}
	modules.Modules["policies"] = service
}

func (s *PolicyService) CreatePolicy(clientId string, policy *gocloak.PolicyRepresentation) error {
	_, err := s.client.CreatePolicy(s.ctx, s.token, "products", clientId, *policy)
	if err != nil {
		log.Err(err).Str("name", *policy.Name).Msg("Cannot create policy")
		return err
	} else {
		log.Info().Str("name", *policy.Name).Msg("policy created")
	}
	return nil
}
