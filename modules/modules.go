package modules

import (
	"keycloak-tools/access"

	"github.com/Nerzal/gocloak/v7"
)

//TODO replace with proper config from file or env vars
const REALM_NAME = "products"

type KeycloakClientDiffGenCtx struct {
	ClientOp *ClientOp
	// Config *KeycloakConfig
	Config *KeycloakClientConfig
}

// type KeycloakConfig struct {
// 	Groups       []gocloak.Group `json:"groups"`
// 	ClientConfig ClientConfig    `json:"clientConfig"`
// }

type KeycloakClientConfig struct {
	Definition  gocloak.Client                     `json:"clientDefinition"`
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
	//TODO move this to realm-dedicated part of code
	Groups       []GroupsOp         `json:"groups,omitempty"`
	ClientConfig ClientConfigOpSpec `json:"clientConfig,omitempty"`
	Scopes       []ScopesOp         `json:"scopes,omitempty"`
	Resources    []ResourcesOp      `json:"resources,omitempty"`
	Policies     []PoliciesOp       `json:"policies,omitempty"`
	Permissions  []PermissionsOp    `json:"permissions,omitempty"`
}

type GroupsOp struct {
	Op        string        `json:"op"`
	GroupSpec gocloak.Group `json:"groupSpec"`
}

type ClientConfigOpSpec struct {
	Declaration gocloak.Client  `json:"declaration,omitempty"`
	Op          string          `json:"op,omitempty"`
	Scopes      []ScopesOp      `json:"scopes,omitempty"`
	Resources   []ResourcesOp   `json:"resources,omitempty"`
	Policies    []PoliciesOp    `json:"policies,omitempty"`
	Permissions []PermissionsOp `json:"permissions,omitempty"`
}

type ClientOp struct {
	Op     string         `json:"op"`
	Client gocloak.Client `json:"client"`
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
	Diff(keycloakConfig *KeycloakClientDiffGenCtx, opsConfig *KeycloakOpsConfig) error
	Order() int
}
