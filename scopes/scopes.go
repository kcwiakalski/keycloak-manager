package scopes

import (
	"context"
	"keycloak-manager/access"
	"keycloak-manager/modules"

	"github.com/Nerzal/gocloak/v7"
	"github.com/rs/zerolog/log"
)

type scopeService struct {
	client gocloak.GoCloak
	ctx    context.Context
	token  string
	realm  string
}

// implementation of modules.ConfigurationHandler.Apply method
func (s *scopeService) Apply(keycloakConfig *modules.ClientChangeContext) error {
	var finalError error
	clientId := *keycloakConfig.Client.ID
	for _, scope := range keycloakConfig.Changes.Scopes {
		if scope.Op == "ADD" {
			err := s.addScope(clientId, &scope.ScopeSpec)
			if err != nil {
				finalError = err
			}
		}
		if scope.Op == "DEL" {
			s.deleteScope(clientId, &scope.ScopeSpec)
		}
		if scope.Op == "UPD" {
			s.updateScope(clientId, &scope.ScopeSpec)
		}
	}
	return finalError
}

func (s *scopeService) Order() int {
	return 2
}

func scopeEquals(first *gocloak.ScopeRepresentation, second *gocloak.ScopeRepresentation) bool {
	if first == nil && second == nil {
		return true
	}
	if (first == nil && second != nil) || (first != nil && second == nil) {
		return false
	}
	if first.DisplayName == nil && second.DisplayName == nil {
		return true
	}
	if (first.DisplayName == nil && second.DisplayName != nil) || (first.DisplayName != nil && second.DisplayName == nil) {
		return false
	}
	if *first.DisplayName == *second.DisplayName {
		return true
	}
	return false
}

// implementation of modules.DiffHandler.Diff method
func (s *scopeService) Diff(keycloakConfig *modules.ClientDiffContext, opsConfig *modules.ClientChanges) error {
	var ops []modules.ScopesOp = make([]modules.ScopesOp, 0)
	var existingScopes []*gocloak.ScopeRepresentation
	if keycloakConfig.ClientOp.Op == "NONE" {
		var err error
		existingScopes, err = s.getScopes(*keycloakConfig.ClientOp.ClientSpec.ID)
		if err != nil {
			return err
		}
	}
	x0 := keycloakConfig.Declaration.Scopes
	var configuredScopes map[string]gocloak.ScopeRepresentation = make(map[string]gocloak.ScopeRepresentation)
	for _, inputScope := range x0 {
		configuredScopes[*inputScope.Name] = inputScope
	}
	for _, existingScope := range existingScopes {
		name := *existingScope.Name
		configuredScope, found := configuredScopes[name]
		if found {
			if !scopeEquals(&configuredScope, existingScope) {
				log.Info().Str("name", name).Msg("Scope update detected, update op required")
				configuredScope.ID = existingScope.ID
				ops = append(ops, modules.ScopesOp{
					Op:        "UPD",
					ScopeSpec: *&configuredScope,
				})
			}
			delete(configuredScopes, name)
		} else {
			log.Info().Str("name", name).Msg("Deprecated/Old Scope detected, delete op required")
			ops = append(ops, modules.ScopesOp{
				Op:        "DEL",
				ScopeSpec: *existingScope,
			})
		}
	}
	for key := range configuredScopes {
		scope := configuredScopes[key]
		log.Info().Str("name", *scope.Name).Str("key", key).Msg("New scope detected, add op required")
		ops = append(ops, modules.ScopesOp{
			Op:        "ADD",
			ScopeSpec: scope,
		})
	}
	opsConfig.Scopes = ops
	return nil
}

func InitializeService(keycloak *access.KeycloakContext, applyModules map[string]modules.ConfigurationHandler, diffHandlers map[string]modules.DiffHandler) {
	ctx := keycloak
	service := &scopeService{
		client: ctx.Client,
		ctx:    ctx.Ctx,
		token:  ctx.Token.AccessToken,
		realm:  ctx.Realm,
	}
	applyModules["scopes"] = service
	diffHandlers["scopes"] = service
}

// simple wrapper for keycloak service
func (s *scopeService) addScope(clientId string, scope *gocloak.ScopeRepresentation) error {
	_, err := s.client.CreateScope(s.ctx, s.token, s.realm, clientId, *scope)
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
	err := s.client.DeleteScope(s.ctx, s.token, s.realm, clientId, *scope.ID)
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
	scopes, err := s.client.GetScopes(s.ctx, s.token, s.realm, clientId, params)
	if err != nil {
		log.Err(err).Str("client", clientId).Msg("Fetching client scopes failed")
	}
	return scopes, err
}

//updateScope - Simple wrapper for keycloak service
func (s *scopeService) updateScope(clientId string, scope *gocloak.ScopeRepresentation) error {
	err := s.client.UpdateScope(s.ctx, s.token, s.realm, clientId, *scope)
	if err != nil {
		log.Err(err).Str("name", *scope.Name).Msg("Cannot update scope")
		return err
	} else {
		log.Info().Str("name", *scope.Name).Msg("Scope updated")
	}
	return nil
}
