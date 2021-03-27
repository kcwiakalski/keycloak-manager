package modules

import (
	"keycloak-manager/access"
	"keycloak-manager/model"

	"github.com/Nerzal/gocloak/v7"
	"github.com/alecthomas/kong"
)

type ClientDeclaration struct {
	Client      gocloak.Client                     `json:"clientDefinition"`
	Scopes      []gocloak.ScopeRepresentation      `json:"scopes,omitempty"`
	Resources   []gocloak.ResourceRepresentation   `json:"resources,omitempty"`
	Policies    []gocloak.PolicyRepresentation     `json:"policies,omitempty"`
	Permissions []gocloak.PermissionRepresentation `json:"permissions,omitempty"`
	Groups      []gocloak.Group                    `json:"groups,omitempty"`
	Roles       []gocloak.Role                     `json:"roles,omitempty"`
}

type ClientChanges struct {
	Client      ClientOp        `json:"client"`
	Scopes      []ScopesOp      `json:"scopes,omitempty"`
	Resources   []ResourcesOp   `json:"resources,omitempty"`
	Policies    []PoliciesOp    `json:"policies,omitempty"`
	Permissions []PermissionsOp `json:"permissions,omitempty"`
	Groups      []GroupsOp      `json:"groups,omitempty"`
	Roles       []RolesOp       `json:"roles,omitempty"`
}

type ClientChangeContext struct {
	Changes     *ClientChanges
	Declaration *ClientDeclaration
	Client      *gocloak.Client
}
type ClientDiffContext struct {
	ClientOp    *ClientOp
	Declaration *ClientDeclaration
}

type GroupsOp struct {
	Op        string        `json:"op"`
	GroupSpec gocloak.Group `json:"groupSpec"`
}

type ClientOp struct {
	Op         string         `json:"op"`
	ClientSpec gocloak.Client `json:"clientSpec"`
}
type ScopesOp struct {
	Op        string                      `json:"op"`
	ScopeSpec gocloak.ScopeRepresentation `json:"scopeSpec"`
}
type RolesOp struct {
	Op       string       `json:"op"`
	RoleSpec gocloak.Role `json:"roleSpec"`
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
	model.Ctx = kong.Parse(&model.CLI)
	Keycloak = access.KeycloakConnection()
}

type ConfigurationHandler interface {
	Apply(changeCtx *ClientChangeContext) error
	Order() int
}

type DiffHandler interface {
	// method generating operations required to perform, so server match with config declaration
	Diff(declaration *ClientDiffContext, changes *ClientChanges) error
	Order() int
}
