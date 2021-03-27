package scopes

import (
	"context"
	"keycloak-manager/model"
	"keycloak-manager/modules"

	"github.com/Nerzal/gocloak/v7"
	"github.com/rs/zerolog/log"
)

type scopeService struct {
	client gocloak.GoCloak
	ctx    context.Context
	token  string
}

var service *scopeService

// implementation of modules.ConfigurationHandler.Apply method
func (s *scopeService) Apply(keycloakConfig *modules.ClientChangeContext) error {
	var finalError error
	clientId := *keycloakConfig.Client.ID
	for _, scope := range keycloakConfig.Changes.Scopes {
		if scope.Op == "ADD" {
			err := service.addScope(clientId, &scope.ScopeSpec)
			if err != nil {
				finalError = err
			}
		}
		if scope.Op == "DEL" {
			service.deleteScope(clientId, &scope.ScopeSpec)
		}
	}
	return finalError
}

func (s *scopeService) Order() int {
	return 2
}

// implementation of modules.DiffHandler.Diff method
func (s *scopeService) Diff(keycloakConfig *modules.ClientDiffContext, opsConfig *modules.ClientChanges) error {
	var ops []modules.ScopesOp = make([]modules.ScopesOp, 0)
	var scopes []*gocloak.ScopeRepresentation
	if keycloakConfig.ClientOp.Op == "NONE" {
		var err error
		scopes, err = s.getScopes(*keycloakConfig.ClientOp.ClientSpec.ID)
		if err != nil {
			return err
		}
	}
	x0 := keycloakConfig.Declaration.Scopes
	var inputScopes map[string]gocloak.ScopeRepresentation = make(map[string]gocloak.ScopeRepresentation)
	for _, inputScope := range x0 {
		inputScopes[*inputScope.Name] = inputScope
	}
	for _, scope := range scopes {
		name := *scope.Name
		_, found := inputScopes[name]
		if found {
			delete(inputScopes, name)
		} else {
			log.Info().Str("name", name).Msg("Deprecated/Old Scope detected, delete op required")
			ops = append(ops, modules.ScopesOp{
				Op:        "DEL",
				ScopeSpec: *scope,
			})
		}
	}
	for key := range inputScopes {
		scope := inputScopes[key]
		log.Info().Str("name", *scope.Name).Str("key", key).Msg("New scope detected, add op required")
		ops = append(ops, modules.ScopesOp{
			Op:        "ADD",
			ScopeSpec: scope,
		})
	}
	opsConfig.Scopes = ops
	return nil
}

func init() {
	ctx := modules.Keycloak
	service = &scopeService{
		client: ctx.Client,
		ctx:    ctx.Ctx,
		token:  ctx.Token.AccessToken,
	}
	modules.Modules["scopes"] = service
	modules.DiffModules["scopes"] = service
}

// simple wrapper for keycloak service
func (s *scopeService) addScope(clientId string, scope *gocloak.ScopeRepresentation) error {
	_, err := s.client.CreateScope(s.ctx, s.token, model.CLI.Realm, clientId, *scope)
	if err != nil {
		log.Err(err).Str("name", *scope.Name).Msg("Cannot create scope")
		return err
	} else {
		log.Info().Str("name", *scope.Name).Msg("Scope created")
	}
	return nil
}

//deleteScope - Simple wrapper for keycloak service
func (s *scopeService) deleteScope(clientId string, scope *gocloak.ScopeRepresentation) error {
	err := s.client.DeleteScope(s.ctx, s.token, model.CLI.Realm, clientId, *scope.ID)
	if err != nil {
		log.Err(err).Str("name", *scope.Name).Msg("Cannot remove scope")
		return err
	} else {
		log.Info().Str("name", *scope.Name).Msg("Scope removed")
	}
	return nil
}

// Simple wrapper for keycloak service
func (s *scopeService) getScopes(clientId string) ([]*gocloak.ScopeRepresentation, error) {
	deep := false
	max := 200
	params := gocloak.GetScopeParams{
		Deep: &deep,
		Max:  &max,
	}
	scopes, err := s.client.GetScopes(s.ctx, s.token, model.CLI.Realm, clientId, params)
	if err != nil {
		log.Err(err).Str("client", clientId).Msg("Fetching client scopes failed")
	}
	return scopes, err

}
