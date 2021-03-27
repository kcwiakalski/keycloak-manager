package roles

import (
	"context"
	"keycloak-manager/model"
	"keycloak-manager/modules"

	"github.com/Nerzal/gocloak/v7"
	"github.com/rs/zerolog/log"
)

type clientRoleService struct {
	client gocloak.GoCloak
	ctx    context.Context
	token  string
}

var service *clientRoleService

// implementation of modules.ConfigurationHandler.Apply method
func (s *clientRoleService) Apply(keycloakConfig *modules.ClientChangeContext) error {
	var finalError error
	clientId := *keycloakConfig.Client.ID
	for _, role := range keycloakConfig.Changes.Roles {
		if role.Op == "ADD" {
			err := service.addRole(clientId, &role.RoleSpec)
			if err != nil {
				finalError = err
			}
		}
		if role.Op == "DEL" {
			service.deleteRole(clientId, &role.RoleSpec)
		}
	}
	return finalError
}

func (s *clientRoleService) Order() int {
	return 1
}

// implementation of modules.DiffHandler.Diff method
func (s *clientRoleService) Diff(keycloakConfig *modules.ClientDiffContext, opsConfig *modules.ClientChanges) error {
	var ops []modules.RolesOp = make([]modules.RolesOp, 0)
	var roles []*gocloak.Role
	if keycloakConfig.ClientOp.Op == "NONE" {
		var err error
		roles, err = s.getRoles(*keycloakConfig.ClientOp.ClientSpec.ID)
		if err != nil {
			return err
		}
	}
	for _, role := range roles {
		log.Info().Msg(*role.Name)
	}
	x0 := keycloakConfig.Declaration.Roles
	var inputRoles map[string]gocloak.Role = make(map[string]gocloak.Role)
	for _, inputRole := range x0 {
		inputRoles[*inputRole.Name] = inputRole
	}
	for _, role := range roles {
		name := *role.Name
		_, found := inputRoles[name]
		if found {
			delete(inputRoles, name)
		} else {
			log.Info().Str("name", name).Msg("Deprecated/Old Scope detected, delete op required")
			ops = append(ops, modules.RolesOp{
				Op:       "DEL",
				RoleSpec: *role,
			})
		}
	}
	for key := range inputRoles {
		role := inputRoles[key]
		log.Info().Str("name", *role.Name).Str("key", key).Msg("New client role detected, add op required")
		ops = append(ops, modules.RolesOp{
			Op:       "ADD",
			RoleSpec: role,
		})
	}
	opsConfig.Roles = ops
	return nil
}

func init() {
	ctx := modules.Keycloak
	service = &clientRoleService{
		client: ctx.Client,
		ctx:    ctx.Ctx,
		token:  ctx.Token.AccessToken,
	}
	modules.Modules["client-roles"] = service
	modules.DiffModules["client-roles"] = service
}

// simple wrapper for keycloak service
func (s *clientRoleService) addRole(clientId string, role *gocloak.Role) error {
	_, err := s.client.CreateClientRole(s.ctx, s.token, model.CLI.Realm, clientId, *role)
	if err != nil {
		log.Err(err).Str("name", *role.Name).Msg("Cannot create role")
		return err
	} else {
		log.Info().Str("name", *role.Name).Msg("Role created")
	}
	return nil
}

//deleteScope - Simple wrapper for keycloak service
func (s *clientRoleService) deleteRole(clientId string, role *gocloak.Role) error {
	err := s.client.DeleteClientRole(s.ctx, s.token, model.CLI.Realm, clientId, *role.Name)
	if err != nil {
		log.Err(err).Str("name", *role.Name).Msg("Cannot remove role")
		return err
	} else {
		log.Info().Str("name", *role.Name).Msg("Role removed")
	}
	return nil
}

// Simple wrapper for keycloak service
func (s *clientRoleService) getRoles(clientId string) ([]*gocloak.Role, error) {
	roles, err := s.client.GetClientRoles(s.ctx, s.token, model.CLI.Realm, clientId)
	if err != nil {
		log.Err(err).Str("client", clientId).Msg("Fetching client roles failed")
	}
	return roles, err

}
