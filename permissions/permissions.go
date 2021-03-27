package permissions

import (
	"context"
	"keycloak-manager/model"
	"keycloak-manager/modules"

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
	modules.DiffModules["permissions"] = service
}

func (s *permissionService) Apply(keycloakConfig *modules.ClientChangeContext) error {
	var finalError error
	clientId := *keycloakConfig.Client.ID
	for _, perm := range keycloakConfig.Changes.Permissions {
		if perm.Op == "ADD" {
			err := service.AddPermission(clientId, perm.PermSpec)
			if err != nil {
				finalError = err
			}
		} else if perm.Op == "DEL" {
			err := service.DeletePermission(clientId, perm.PermSpec)
			if err != nil {
				finalError = err
			}
		}
	}
	return finalError
}

func (s *permissionService) Order() int {
	return 4
}

func (s *permissionService) AddPermission(clientId string, permission gocloak.PermissionRepresentation) error {
	_, err := s.client.CreatePermission(s.ctx, s.token, model.CLI.Realm, clientId, permission)
	if err != nil {
		log.Err(err).Str("name", *permission.Name).Msg("Cannot create permission")
		return err
	} else {
		log.Info().Str("name", *permission.Name).Msg("Permission created")
	}
	return nil
}
func (s *permissionService) DeletePermission(clientId string, permission gocloak.PermissionRepresentation) error {
	err := s.client.DeletePermission(s.ctx, s.token, model.CLI.Realm, clientId, *permission.ID)
	if err != nil {
		log.Err(err).Str("name", *permission.Name).Msg("Cannot delete permission")
		return err
	} else {
		log.Info().Str("name", *permission.Name).Msg("Permission deleted")
	}
	return nil
}

func (s *permissionService) getPermissions(clientName string) ([]*gocloak.PermissionRepresentation, error) {
	params := gocloak.GetPermissionParams{}
	permissions, err := s.client.GetPermissions(s.ctx, s.token, model.CLI.Realm, clientName, params)
	if err != nil {
		return nil, err
	}
	return permissions, nil
}

func (s *permissionService) Diff(keycloakConfig *modules.ClientDiffContext, opsConfig *modules.ClientChanges) error {
	var ops []modules.PermissionsOp = make([]modules.PermissionsOp, 0)
	var existingPerms []*gocloak.PermissionRepresentation
	if keycloakConfig.ClientOp.Op == "NONE" {
		var err error
		existingPerms, err = s.getPermissions(*keycloakConfig.ClientOp.ClientSpec.ID)
		if err != nil {
			return err
		}
	}
	expectedPerms := keycloakConfig.Declaration.Permissions
	var expectedPermsMap map[string]gocloak.PermissionRepresentation = make(map[string]gocloak.PermissionRepresentation)
	for _, expectedPerm := range expectedPerms {
		expectedPermsMap[*expectedPerm.Name] = expectedPerm
	}
	for _, existingPerm := range existingPerms {
		name := *existingPerm.Name
		_, found := expectedPermsMap[name]
		if found {
			delete(expectedPermsMap, name)
		} else {
			log.Info().Str("name", name).Msg("Deprecated/Old Permission detected, delete op required")
			ops = append(ops, modules.PermissionsOp{
				Op:       "DEL",
				PermSpec: *existingPerm,
			})
		}
	}
	for key := range expectedPermsMap {
		perm := expectedPermsMap[key]
		log.Info().Str("name", *perm.Name).Str("key", key).Msg("New permission detected, add op required")
		ops = append(ops, modules.PermissionsOp{
			Op:       "ADD",
			PermSpec: perm,
		})
	}
	opsConfig.Permissions = ops
	return nil
}
