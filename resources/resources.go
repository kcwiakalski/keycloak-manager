package resources

import (
	"context"
	"keycloak-manager/access"
	"keycloak-manager/modules"

	"github.com/Nerzal/gocloak/v7"
	"github.com/rs/zerolog/log"
)

type resourceService struct {
	client gocloak.GoCloak
	ctx    context.Context
	token  string
	realm  string
}

func InitializeService(keycloak *access.KeycloakContext, applyHandlers map[string]modules.ConfigurationHandler, diffHandlers map[string]modules.DiffHandler) {
	service := &resourceService{
		client: keycloak.Client,
		ctx:    keycloak.Ctx,
		token:  keycloak.Token.AccessToken,
		realm:  keycloak.Realm,
	}
	applyHandlers["resources"] = service
	diffHandlers["resources"] = service
}

func (s *resourceService) addResource(clientId string, resource gocloak.ResourceRepresentation) error {
	_, err := s.client.CreateResource(s.ctx, s.token, s.realm, clientId, resource)
	if err != nil {
		log.Err(err).Str("name", *resource.Name).Msg("Cannot create resource")
		return err
	} else {
		log.Info().Str("name", *resource.Name).Msg("Resource created")
	}
	return nil
}
func (s *resourceService) getResources(clientName string) ([]*gocloak.ResourceRepresentation, error) {
	params := gocloak.GetResourceParams{}
	resources, err := s.client.GetResources(s.ctx, s.token, s.realm, clientName, params)
	if err != nil {
		return nil, err
	}
	return resources, nil
}
func (s *resourceService) deleteResource(clientId string, resource gocloak.ResourceRepresentation) error {
	err := s.client.DeleteResource(s.ctx, s.token, s.realm, clientId, *resource.ID)
	if err != nil {
		log.Err(err).Str("resourceName", *resource.Name).Msg("Cannot delete deprecated resource")
		return err
	} else {
		log.Info().Str("name", *resource.Name).Msg("Resource deleted")
	}
	return nil
}
func (s *resourceService) updateResource(clientId string, resource gocloak.ResourceRepresentation) error {
	err := s.client.UpdateResource(s.ctx, s.token, s.realm, clientId, resource)
	if err != nil {
		log.Err(err).Str("resourceName", *resource.Name).Msg("Cannot update resource")
		return err
	} else {
		log.Info().Str("name", *resource.Name).Msg("Resource updated")
	}
	return nil
}

func resourceEquals(first *gocloak.ResourceRepresentation, second *gocloak.ResourceRepresentation) bool {
	if *first.DisplayName != *second.DisplayName {
		return false
	}
	if len(*first.Scopes) != len(*second.Scopes) {
		return false
	}
	var firstScopes = make([]string, 0)
	for _, scope := range *first.Scopes {
		firstScopes = append(firstScopes, *scope.Name)
	}
	var remainingScopes = make([]string, len(firstScopes))
	copy(remainingScopes, firstScopes)
	var secondeScopes = make([]string, 0)
	for _, scope := range *second.Scopes {
		secondeScopes = append(secondeScopes, *scope.Name)
	}

	for i, firstScope := range firstScopes {
		for j, otherScope := range secondeScopes {
			if firstScope == otherScope {
				secondeScopes[j] = secondeScopes[len(secondeScopes)-1]
				secondeScopes = secondeScopes[:len(secondeScopes)-1]
				remainingScopes[i] = ""
				break
			}
		}
	}
	remaining := false
	for _, scope := range remainingScopes {
		if scope != "" {
			remaining = true
		}
	}
	if len(secondeScopes) > 0 || remaining {
		return false
	}
	return true
}

func (s *resourceService) Diff(keycloakConfig *modules.ClientDiffContext, opsConfig *modules.ClientChanges) error {
	var ops []modules.ResourcesOp = make([]modules.ResourcesOp, 0)
	var existingResources []*gocloak.ResourceRepresentation
	if keycloakConfig.ClientOp.Op == "NONE" {
		var err error
		existingResources, err = s.getResources(*keycloakConfig.ClientOp.ClientSpec.ID)
		if err != nil {
			return err
		}
	}
	configuredResources := keycloakConfig.Declaration.Resources
	var expectedResources map[string]gocloak.ResourceRepresentation = make(map[string]gocloak.ResourceRepresentation)
	for _, configuredResource := range configuredResources {
		expectedResources[*configuredResource.Name] = configuredResource
	}
	for _, resource := range existingResources {
		name := *resource.Name
		expectedResource, found := expectedResources[name]
		if found {
			if !resourceEquals(resource, &expectedResource) {
				log.Info().Str("name", name).Msg("Resource changed detected, update op required")
				finalResource := expectedResource
				finalResource.ID = resource.ID
				ops = append(ops, modules.ResourcesOp{
					Op:           "UPD",
					ResourceSpec: finalResource,
				})
			}
			delete(expectedResources, name)
		} else {
			log.Info().Str("name", name).Msg("Deprecated/Old Resource detected, delete op required")
			ops = append(ops, modules.ResourcesOp{
				Op:           "DEL",
				ResourceSpec: *resource,
			})
		}
	}
	for key := range expectedResources {
		resource := expectedResources[key]
		log.Info().Str("name", *resource.Name).Str("key", key).Msg("New resource detected, add op required")
		ops = append(ops, modules.ResourcesOp{
			Op:           "ADD",
			ResourceSpec: resource,
		})
	}
	opsConfig.Resources = ops
	return nil
}

func (s *resourceService) Apply(keycloakConfig *modules.ClientChangeContext) error {
	var finalError error
	clientId := *keycloakConfig.Client.ID
	for _, resource := range keycloakConfig.Changes.Resources {
		if resource.Op == "ADD" {
			err := s.addResource(clientId, resource.ResourceSpec)
			if err != nil {
				finalError = err
			}
		} else if resource.Op == "DEL" {
			err := s.deleteResource(clientId, resource.ResourceSpec)
			if err != nil {
				finalError = err
			}
		} else if resource.Op == "UPD" {
			err := s.updateResource(clientId, resource.ResourceSpec)
			if err != nil {
				finalError = err
			}
		}
	}
	return finalError
}

func (s *resourceService) Order() int {
	return 3
}
