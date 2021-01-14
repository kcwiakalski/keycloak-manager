package modules

import (
	"keycloak-tools/access"

	"github.com/Nerzal/gocloak/v7"
)

type DiffGenCtx struct {
	Client *gocloak.Client
	Config *KeycloakConfig
}

type KeycloakConfig struct {
	Groups       []gocloak.Group `json:"groups"`
	ClientConfig ClientConfig    `json:"clientConfig"`
}
type ClientConfig struct {
	Name        string                             `json:"name"`
	Scopes      []gocloak.ScopeRepresentation      `json:"scopes,omitempty"`
	Resources   []gocloak.ResourceRepresentation   `json:"resources,omitempty"`
	Policies    []gocloak.PolicyRepresentation     `json:"policies,omitempty"`
	Permissions []gocloak.PermissionRepresentation `json:"permissions,omitempty"`
}

type ConfigurationContext struct {
	Config *KeycloakOpsConfig
	Client *gocloak.Client
}

type KeycloakOpsConfig struct {
	Groups       []GroupsOp         `json:"groups,omitempty"`
	ClientConfig ClientConfigOpSpec `json:"clientConfig,omitempty"`
}

type GroupsOp struct {
	Op        string        `json:"op"`
	GroupSpec gocloak.Group `json:"groupSpec"`
}

type ClientConfigOpSpec struct {
	Name        string          `json:"name"`
	Scopes      []ScopesOp      `json:"scopes,omitempty"`
	Resources   []ResourcesOp   `json:"resources,omitempty"`
	Policies    []PoliciesOp    `json:"policies,omitempty"`
	Permissions []PermissionsOp `json:"permissions,omitempty"`
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
var DiffModules map[string]DiffHandler = make(map[string]DiffHandler)

func init() {
	Keycloak = access.KeycloakConnection()
}

type ConfigurationHandler interface {
	Apply(keycloakConfig *ConfigurationContext) error
	Order() int
}

type DiffHandler interface {
	// method generating operations required to perform, so server match with config declaration
	Diff(keycloakConfig *DiffGenCtx, opsConfig *KeycloakOpsConfig) error
	Order() int
}
