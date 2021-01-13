package modules

import (
	"keycloak-tools/access"

	"github.com/Nerzal/gocloak/v7"
)

type ConfigurationContext struct {
	Config *KeycloakConfig
	Client *gocloak.Client
}

type KeycloakConfig struct {
	Groups       []GroupsOp         `json:"groups"`
	ClientConfig ClientConfigOpSpec `json:"clientConfig"`
}

type GroupsOp struct {
	Op        string        `json:"op"`
	GroupSpec gocloak.Group `json:"groupSpec"`
}

type ClientConfigOpSpec struct {
	Name        string          `json:"name"`
	Scopes      []ScopesOp      `json:"scopes"`
	Resources   []ResourcesOp   `json:"resources"`
	Policies    []PoliciesOp    `json:"policies"`
	Permissions []PermissionsOp `json:"permissions"`
}

type ScopesOp struct {
	Op        string                      `json:"op"`
	ScopeSpec gocloak.ScopeRepresentation `json:"scopeSpec"`
}

type ResourcesOp struct {
	Op           string                         `json:"op"`
	ResourceSpec gocloak.ResourceRepresentation `json:"resourceSpec"`
}
type PoliciesOp struct {
	Op         string                       `json:"op"`
	PolicySpec gocloak.PolicyRepresentation `json:"policySpec"`
}
type PermissionsOp struct {
	Op       string                           `json:"op"`
	PermSpec gocloak.PermissionRepresentation `json:"permSpec"`
}

var Keycloak *access.KeycloakContext

var Modules map[string]ConfigurationHandler = make(map[string]ConfigurationHandler)

func init() {
	Keycloak = access.KeycloakConnection()
}

type ConfigurationHandler interface {
	Apply(keycloakConfig *ConfigurationContext) error
	Order() int
}
