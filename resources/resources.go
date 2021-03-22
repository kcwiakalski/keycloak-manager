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

func init() {
	ctx := modules.Keycloak
	service = &resourceService{
		client: ctx.Client,
		ctx:    ctx.Ctx,
		token:  ctx.Token.AccessToken,
	}
	modules.Modules["resources"] = service
	modules.DiffModules["resources"] = service
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
func (s *resourceService) getResources(clientName string) ([]*gocloak.ResourceRepresentation, error) {
	params := gocloak.GetResourceParams{}
	resources, err := s.client.GetResources(s.ctx, s.token, modules.REALM_NAME, clientName, params)
	if err != nil {
		return nil, err
	}
	return resources, nil
}
func (s *resourceService) deleteResource(clientId string, resource gocloak.ResourceRepresentation) error {
	err := s.client.DeleteResource(s.ctx, s.token, modules.REALM_NAME, clientId, *resource.ID)
	if err != nil {
		log.Err(err).Str("resourceName", *resource.Name).Msg("Cannot delete deprecated resource")
		return err
	}
	return nil
}
func (s *resourceService) Diff(keycloakConfig *modules.ClientDiffContext, opsConfig *modules.ClientChanges) error {
	var ops []modules.ResourcesOp = make([]modules.ResourcesOp, 0)
	var resources []*gocloak.ResourceRepresentation
	if keycloakConfig.ClientOp.Op == "NONE" {
		var err error
		resources, err = s.getResources(*keycloakConfig.ClientOp.ClientSpec.ID)
		if err != nil {
			return err
		}
	}
	x0 := keycloakConfig.Declaration.Resources
	var inputResources map[string]gocloak.ResourceRepresentation = make(map[string]gocloak.ResourceRepresentation)
	for _, inputResource := range x0 {
		inputResources[*inputResource.Name] = inputResource
	}
	for _, resource := range resources {
		name := *resource.Name
		_, found := inputResources[name]
		if found {
			delete(inputResources, name)
		} else {
			log.Info().Str("name", name).Msg("Deprecated/Old Resource detected, delete op required")
			ops = append(ops, modules.ResourcesOp{
				Op:           "DEL",
				ResourceSpec: *resource,
			})
		}
	}
	for key := range inputResources {
		resource := inputResources[key]
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
			err := service.addResource(clientId, resource.ResourceSpec)
			if err != nil {
				finalError = err
			}
		} else if resource.Op == "DEL" {
			err := service.deleteResource(clientId, resource.ResourceSpec)
			if err != nil {
				finalError = err
			}
		}
	}
	return finalError
}

func (s *resourceService) Order() int {
	return 2
}
