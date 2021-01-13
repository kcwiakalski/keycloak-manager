package scopes

import (
	"context"
	"keycloak-tools/modules"

	"github.com/Nerzal/gocloak/v7"
	"github.com/rs/zerolog/log"
)

type scopeService struct {
	client gocloak.GoCloak
	ctx    context.Context
	token  string
}

var service *scopeService

func (s *scopeService) Apply(keycloakConfig *modules.ConfigurationContext) error {
	var finalError error
	clientId := *keycloakConfig.Client.ID
	for _, scope := range keycloakConfig.Config.ClientConfig.Scopes {
		err := service.addScope(clientId, &scope.ScopeSpec)
		if err != nil {
			finalError = err
		}
	}
	return finalError
}

func (s *scopeService) Order() int {
	return 1
}

func init() {
	ctx := modules.Keycloak
	service = &scopeService{
		client: ctx.Client,
		ctx:    ctx.Ctx,
		token:  ctx.Token.AccessToken,
	}
	modules.Modules["scopes"] = service
}

func (s *scopeService) addScope(clientId string, scope *gocloak.ScopeRepresentation) error {
	_, err := s.client.CreateScope(s.ctx, s.token, "products", clientId, *scope)
	if err != nil {
		log.Err(err).Str("name", *scope.Name).Msg("Cannot create scope")
		return err
	} else {
		log.Info().Str("name", *scope.Name).Msg("Scope created")
	}
	return nil
}
