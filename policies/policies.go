package policies

import (
	"context"
	"keycloak-manager/model"
	"keycloak-manager/modules"

	"github.com/Nerzal/gocloak/v7"
	"github.com/rs/zerolog/log"
)

type policyService struct {
	client gocloak.GoCloak
	ctx    context.Context
	token  string
}

var service *policyService

func (s *policyService) Apply(keycloakConfig *modules.ClientChangeContext) error {
	var finalError error
	clientId := *keycloakConfig.Client.ID
	for _, policy := range keycloakConfig.Changes.Policies {
		if policy.Op == "ADD" {

			err := service.CreatePolicy(clientId, &policy.PolicySpec)
			if err != nil {
				finalError = err
			}
		} else if policy.Op == "DEL" {
			err := service.deletePolicy(clientId, &policy.PolicySpec)
			if err != nil {
				finalError = err
			}
		}
	}
	return finalError
}

func (s *policyService) Order() int {
	return 4
}

func init() {
	ctx := modules.Keycloak
	service = &policyService{
		client: ctx.Client,
		ctx:    ctx.Ctx,
		token:  ctx.Token.AccessToken,
	}
	modules.Modules["policies"] = service
	modules.DiffModules["policies"] = service
}

func (s *policyService) CreatePolicy(clientId string, policy *gocloak.PolicyRepresentation) error {
	_, err := s.client.CreatePolicy(s.ctx, s.token, model.CLI.Realm, clientId, *policy)
	if err != nil {
		log.Err(err).Str("name", *policy.Name).Msg("Cannot create policy")
		return err
	} else {
		log.Info().Str("name", *policy.Name).Msg("policy created")
	}
	return nil
}

func (s *policyService) deletePolicy(clientId string, policy *gocloak.PolicyRepresentation) error {
	err := s.client.DeletePolicy(s.ctx, s.token, model.CLI.Realm, clientId, *policy.ID)
	if err != nil {
		log.Err(err).Str("name", *policy.Name).Msg("Cannot remove policy")
		return err
	} else {
		log.Info().Str("name", *policy.Name).Msg("policy removed")
	}
	return nil
}
func (s *policyService) Diff(keycloakConfig *modules.ClientDiffContext, opsConfig *modules.ClientChanges) error {
	var ops []modules.PoliciesOp = make([]modules.PoliciesOp, 0)
	var existingPolicies []*gocloak.PolicyRepresentation
	if keycloakConfig.ClientOp.Op == "NONE" {
		var err error
		existingPolicies, err = s.getPoliciesForClient(*keycloakConfig.ClientOp.ClientSpec.ID)
		if err != nil {
			return err
		}
	}
	expectedPolicies := keycloakConfig.Declaration.Policies
	var expectedPoliciesMap map[string]gocloak.PolicyRepresentation = make(map[string]gocloak.PolicyRepresentation)
	for _, expectedPolicy := range expectedPolicies {
		expectedPoliciesMap[*expectedPolicy.Name] = expectedPolicy
	}
	for _, policy := range existingPolicies {
		name := *policy.Name
		_, found := expectedPoliciesMap[name]
		if found {
			delete(expectedPoliciesMap, name)
		} else {
			log.Info().Str("name", name).Msg("Deprecated/Old Policy detected, delete op required")
			ops = append(ops, modules.PoliciesOp{
				Op:         "DEL",
				PolicySpec: *policy,
			})
		}
	}
	for key := range expectedPoliciesMap {
		policy := expectedPoliciesMap[key]
		log.Info().Str("name", *policy.Name).Str("key", key).Msg("New resource detected, add op required")
		ops = append(ops, modules.PoliciesOp{
			Op:         "ADD",
			PolicySpec: policy,
		})
	}
	opsConfig.Policies = ops
	return nil
}
func (s *policyService) getPoliciesForClient(clientName string) ([]*gocloak.PolicyRepresentation, error) {
	noPerms := false
	params := gocloak.GetPolicyParams{
		Permission: &noPerms,
	}
	policies, err := s.client.GetPolicies(s.ctx, s.token, model.CLI.Realm, clientName, params)
	if err != nil {
		return nil, err
	}
	return policies, nil
}
