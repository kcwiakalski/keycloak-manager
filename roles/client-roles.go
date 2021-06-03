package roles

import (
	"context"
	"keycloak-manager/access"
	"keycloak-manager/modules"
	"keycloak-manager/tools"

	"github.com/Nerzal/gocloak/v7"
	"github.com/rs/zerolog/log"
)

type clientRoleService struct {
	client gocloak.GoCloak
	ctx    context.Context
	token  string
	realm  string
}

// implementation of modules.ConfigurationHandler.Apply method
func (s *clientRoleService) Apply(keycloakConfig *modules.ClientChangeContext) error {
	var finalError error
	clientId := *keycloakConfig.Client.ID
	for _, role := range keycloakConfig.Changes.Roles {
		if role.Op == "ADD" {
			err := s.addRole(clientId, &role.RoleSpec)
			if err != nil {
				finalError = err
			}
		}
		if role.Op == "UPD" {
			err := s.updateRole(clientId, &role.RoleSpec)
			if err != nil {
				finalError = err
			}
		}
		if role.Op == "DEL" {
			s.deleteRole(clientId, &role.RoleSpec)
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
	var existingRoles []*gocloak.Role
	if keycloakConfig.ClientOp.Op == "NONE" {
		var err error
		existingRoles, err = s.getRoles(*keycloakConfig.ClientOp.ClientSpec.ID)
		if err != nil {
			return err
		}
	}
	configuredRoles := keycloakConfig.Declaration.Roles
	var configuredRolesMap map[string]gocloak.Role = make(map[string]gocloak.Role)
	for _, inputRole := range configuredRoles {
		configuredRolesMap[*inputRole.Name] = inputRole
	}
	for _, existingRole := range existingRoles {
		name := *existingRole.Name
		// configuredRole, found := inputRoles[name]
		configuredRole, found := configuredRolesMap[name]
		if found {
			if !roleEquals(existingRole, &configuredRole) {
				log.Info().Str("name", name).Msg("Client role update detected, update op required")
				configuredRole.ID = existingRole.ID
				ops = append(ops, modules.RolesOp{
					Op:       "UPD",
					RoleSpec: configuredRole,
				})
			}
			delete(configuredRolesMap, name)
		} else {
			log.Info().Str("name", name).Msg("Deprecated/Old role detected, delete op required")
			ops = append(ops, modules.RolesOp{
				Op:       "DEL",
				RoleSpec: *existingRole,
			})
		}
	}
	for key := range configuredRolesMap {
		role := configuredRolesMap[key]
		log.Info().Str("name", *role.Name).Str("key", key).Msg("New client role detected, add op required")
		ops = append(ops, modules.RolesOp{
			Op:       "ADD",
			RoleSpec: role,
		})
	}
	opsConfig.Roles = ops
	return nil
}

func InitializeService(keycloak *access.KeycloakContext, applyModules map[string]modules.ConfigurationHandler, diffHandlers map[string]modules.DiffHandler) {
	service := &clientRoleService{
		client: keycloak.Client,
		ctx:    keycloak.Ctx,
		token:  keycloak.Token.AccessToken,
		realm:  keycloak.Realm,
	}
	applyModules["client-roles"] = service
	diffHandlers["client-roles"] = service
}

// simple wrapper for keycloak service
func (s *clientRoleService) addRole(clientId string, role *gocloak.Role) error {
	_, err := s.client.CreateClientRole(s.ctx, s.token, s.realm, clientId, *role)
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
	err := s.client.DeleteClientRole(s.ctx, s.token, s.realm, clientId, *role.Name)
	if err != nil {
		log.Err(err).Str("name", *role.Name).Msg("Cannot remove role")
		return err
	} else {
		log.Info().Str("name", *role.Name).Msg("Role removed")
	}
	return nil
}

//update client role - Simple wrapper for keycloak service
func (s *clientRoleService) updateRole(clientId string, role *gocloak.Role) error {
	err := s.client.UpdateRole(s.ctx, s.token, s.realm, clientId, *role)
	if err != nil {
		log.Err(err).Str("name", *role.Name).Msg("Cannot update role")
		return err
	} else {
		log.Info().Str("name", *role.Name).Msg("Role updated")
	}
	return nil
}

// Simple wrapper for keycloak service
func (s *clientRoleService) getRoles(clientId string) ([]*gocloak.Role, error) {
	roles, err := s.client.GetClientRoles(s.ctx, s.token, s.realm, clientId)
	if err != nil {
		log.Err(err).Str("client", clientId).Msg("Fetching client roles failed")
	}
	return roles, err

}

// helper methods
func roleEquals(first *gocloak.Role, second *gocloak.Role) bool {
	if tools.ObjectComparableInDepth(first, second) {
		if tools.StringEquals(first.Description, second.Description) {
			return true
		}
	}
	return false
}
