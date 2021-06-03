package policies

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"keycloak-manager/access"
	"keycloak-manager/modules"
	"keycloak-manager/tools"

	"github.com/Nerzal/gocloak/v7"
	"github.com/rs/zerolog/log"
)

type policyService struct {
	client gocloak.GoCloak
	ctx    context.Context
	token  string
	realm  string
}

func (s *policyService) Apply(keycloakConfig *modules.ClientChangeContext) error {
	var finalError error
	clientId := *keycloakConfig.Client.ID
	for _, policy := range keycloakConfig.Changes.Policies {
		if policy.Op == "ADD" {
			err := s.CreatePolicy(clientId, &policy.PolicySpec)
			if err != nil {
				finalError = err
			}
		} else if policy.Op == "UPD" {
			err := s.updatePolicy(clientId, &policy.PolicySpec)
			if err != nil {
				finalError = err
			}
		} else if policy.Op == "DEL" {
			err := s.deletePolicy(clientId, &policy.PolicySpec)
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

func InitializeService(keycloak *access.KeycloakContext, applyHandlers map[string]modules.ConfigurationHandler, diffHandlers map[string]modules.DiffHandler) {
	service := &policyService{
		client: keycloak.Client,
		ctx:    keycloak.Ctx,
		token:  keycloak.Token.AccessToken,
		realm:  keycloak.Realm,
	}
	applyHandlers["policies"] = service
	diffHandlers["policies"] = service
}

func (s *policyService) CreatePolicy(clientId string, policy *gocloak.PolicyRepresentation) error {
	_, err := s.client.CreatePolicy(s.ctx, s.token, s.realm, clientId, *policy)
	if err != nil {
		log.Err(err).Str("name", *policy.Name).Msg("Cannot create policy")
		return err
	} else {
		log.Info().Str("name", *policy.Name).Msg("Policy created")
	}
	return nil
}

func (s *policyService) deletePolicy(clientId string, policy *gocloak.PolicyRepresentation) error {
	err := s.client.DeletePolicy(s.ctx, s.token, s.realm, clientId, *policy.ID)
	if err != nil {
		log.Err(err).Str("name", *policy.Name).Msg("Cannot remove policy")
		return err
	} else {
		log.Info().Str("name", *policy.Name).Msg("Policy removed")
	}
	return nil
}
func (s *policyService) updatePolicy(clientId string, policy *gocloak.PolicyRepresentation) error {
	policyTemp := policy
	//since library does not create valid urls to update policies based on her type we need to do a hack
	switch *policy.Type {
	case "group":
		id := fmt.Sprintf("group/%s", *policy.ID)
		policyTemp.ID = &id
	case "role":
		id := fmt.Sprintf("role/%s", *policy.ID)
		policyTemp.ID = &id
	case "js":
		id := fmt.Sprintf("js/%s", *policy.ID)
		policyTemp.ID = &id
	}
	err := s.client.UpdatePolicy(s.ctx, s.token, s.realm, clientId, *policyTemp)
	if err != nil {
		log.Err(err).Str("name", *policy.Name).Msg("Cannot update policy")
		return err
	} else {
		log.Info().Str("name", *policy.Name).Msg("Policy updated")
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
		for _, existingPolicy := range existingPolicies {
			s.loadFullConfiguration(existingPolicy, *keycloakConfig.ClientOp.ClientSpec.ClientID)
		}
	}
	expectedPolicies := keycloakConfig.Declaration.Policies
	var expectedPoliciesMap map[string]gocloak.PolicyRepresentation = make(map[string]gocloak.PolicyRepresentation)
	for _, expectedPolicy := range expectedPolicies {
		expectedPoliciesMap[*expectedPolicy.Name] = expectedPolicy
	}
	for _, existingPolicy := range existingPolicies {
		name := *existingPolicy.Name
		expectedPolicy, found := expectedPoliciesMap[name]
		if found {
			applyDefaultsIfMissing(&expectedPolicy)
			if !s.policesEquals(&expectedPolicy, existingPolicy) {
				log.Info().Str("name", name).Msg("Policy update detected, update op required")
				expectedPolicy.ID = existingPolicy.ID
				ops = append(ops, modules.PoliciesOp{
					Op:         "UPD",
					PolicySpec: expectedPolicy,
				})
			}
			delete(expectedPoliciesMap, name)
		} else {
			log.Info().Str("name", name).Msg("Deprecated/Old Policy detected, delete op required")
			ops = append(ops, modules.PoliciesOp{
				Op:         "DEL",
				PolicySpec: *existingPolicy,
			})
		}
	}
	for key := range expectedPoliciesMap {
		policy := expectedPoliciesMap[key]
		log.Info().Str("name", *policy.Name).Str("key", key).Msg("New policy detected, add op required")
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
	policies, err := s.client.GetPolicies(s.ctx, s.token, s.realm, clientName, params)
	if err != nil {
		return nil, err
	}
	return policies, nil
}

// helper methods
func (s *policyService) policesEquals(first *gocloak.PolicyRepresentation, second *gocloak.PolicyRepresentation) bool {
	if tools.ObjectComparableInDepth(first, second) {
		if !tools.StringEquals(first.Description, second.Description) {
			return false
		}
		if !tools.StringEquals(first.Type, second.Type) {
			return false
		}
		if !tools.StringEquals((*string)(first.DecisionStrategy), (*string)(second.DecisionStrategy)) {
			return false
		}
		if !tools.StringEquals((*string)(first.Logic), (*string)(second.Logic)) {
			return false
		}
		policyType := *first.Type
		switch policyType {
		case "group":
			var firstGroups []string
			for _, group := range *first.Groups {
				firstGroups = append(firstGroups, *group.Path)
			}
			var secondGroups []string
			for _, group := range *second.Groups {
				secondGroups = append(secondGroups, *group.Path)
			}
			if !sliceEquals(firstGroups, secondGroups) {
				return false
			}
		case "role":
			var firstRoles []string
			for _, role := range *first.Roles {
				firstRoles = append(firstRoles, *role.ID)
			}
			var secondRoles []string
			for _, role := range *second.Roles {
				secondRoles = append(secondRoles, *role.ID)
			}
			if !sliceEquals(firstRoles, secondRoles) {
				return false
			}

		case "js":
			if bytes.Compare([]byte(*first.Code), []byte(*second.Code)) != 0 {
				return false
			}
		}
		return true
	}
	return false
}

func sliceEquals(first []string, second []string) bool {
	firstMap := make(map[string]struct{})
	for _, el1 := range first {
		firstMap[el1] = struct{}{}
	}
	secondMap := make(map[string]struct{})
	for _, el2 := range second {
		secondMap[el2] = struct{}{}
	}
	for key, _ := range firstMap {
		_, found := secondMap[key]
		if found {
			delete(firstMap, key)
			delete(secondMap, key)
		}
	}
	if len(firstMap) > 0 || len(secondMap) > 0 {
		return false
	}

	return true
}

func applyDefaultsIfMissing(policy *gocloak.PolicyRepresentation) {

	if policy.DecisionStrategy == nil {
		policy.DecisionStrategy = gocloak.UNANIMOUS
	}
	if policy.Logic == nil {
		policy.Logic = gocloak.POSITIVE
	}
}

func (s *policyService) loadFullConfiguration(policy *gocloak.PolicyRepresentation, clientName string) {
	switch *policy.Type {
	case "group":
		config := *policy.Config
		groupsConfig, _ := config["groups"]
		var secondGroups2 []gocloak.GroupDefinition
		json.Unmarshal([]byte(groupsConfig), &secondGroups2)
		policy.Groups = &secondGroups2
		for i, group := range secondGroups2 {
			groupId := group.ID
			fullGroup, _ := s.client.GetGroup(s.ctx, s.token, s.realm, *groupId)
			secondGroups2[i].Path = fullGroup.Path
		}
	case "role":
		config := *policy.Config
		rolesConfig, _ := config["roles"]
		var roles []gocloak.RoleDefinition
		json.Unmarshal([]byte(rolesConfig), &roles)
		policy.Roles = &roles
		for i, role := range roles {
			fullRole, _ := s.client.GetClientRoleByID(s.ctx, s.token, s.realm, *role.ID)
			//we should have in our configuration role declared by name not by keykloak id in format {clientName}/{clientRole}
			newRoleName := fmt.Sprintf("%s/%s", clientName, *fullRole.Name)
			roles[i].ID = &newRoleName
		}
	case "js":
		config := *policy.Config
		codeConfig, _ := config["code"]
		policy.Code = &codeConfig
	}

}
