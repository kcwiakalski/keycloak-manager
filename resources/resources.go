package resources

import (
	"context"
	"keycloak-tools/modules"

	"github.com/Nerzal/gocloak/v7"
	"github.com/rs/zerolog/log"
)

type resourceService struct {
	client gocloak.GoCloak
	ctx    context.Context
	token  string
}

var service *resourceService

// TODO remove groupsService := New(keycloak)
func (s *resourceService) Apply(keycloakConfig *modules.ConfigurationContext) error {
	var finalError error
	clientId := *keycloakConfig.Client.ID
	for _, resource := range keycloakConfig.Config.ClientConfig.Resources {
		err := service.addResource(clientId, resource.ResourceSpec)
		if err != nil {
			finalError = err
		}
	}
	return finalError
}

func (s *resourceService) Order() int {
	return 2
}

func init() {
	ctx := modules.Keycloak
	service = &resourceService{
		client: ctx.Client,
		ctx:    ctx.Ctx,
		token:  ctx.Token.AccessToken,
	}
	modules.Modules["resources"] = service
}

func (s *resourceService) addResource(clientId string, resource gocloak.ResourceRepresentation) error {
	_, err := s.client.CreateResource(s.ctx, s.token, "products", clientId, resource)
	if err != nil {
		log.Err(err).Str("name", *resource.Name).Msg("Cannot create resource")
		return err
	} else {
		log.Info().Str("name", *resource.Name).Msg("Resource created")
	}
	return nil
}
