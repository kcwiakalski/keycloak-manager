package permissions

import (
	"context"
	"keycloak-tools/modules"

	"github.com/Nerzal/gocloak/v7"
	"github.com/rs/zerolog/log"
)

type permissionService struct {
	client gocloak.GoCloak
	ctx    context.Context
	token  string
}

var service *permissionService

func init() {
	ctx := modules.Keycloak
	service = &permissionService{
		client: ctx.Client,
		ctx:    ctx.Ctx,
		token:  ctx.Token.AccessToken,
	}
	modules.Modules["permissions"] = service
}

func (s *permissionService) Apply(keycloakConfig *modules.ConfigurationContext) error {
	var finalError error
	clientId := *keycloakConfig.Client.ID
	for _, perm := range keycloakConfig.Config.ClientConfig.Permissions {
		err := service.AddPermission(clientId, perm.PermSpec)
		if err != nil {
			finalError = err
		}
	}
	return finalError
}

func (s *permissionService) Order() int {
	return 4
}

func (s *permissionService) AddPermission(clientId string, permission gocloak.PermissionRepresentation) error {
	_, err := s.client.CreatePermission(s.ctx, s.token, "products", clientId, permission)
	if err != nil {
		log.Err(err).Str("name", *permission.Name).Msg("Cannot create permission")
		return err
	} else {
		log.Info().Str("name", *permission.Name).Msg("Permission created")
	}
	return nil
}
