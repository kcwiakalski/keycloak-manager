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
	Groups       []GroupsConfig   `json:"groups"`
	ClientConfig ClientConfigSpec `json:"clientConfig"`
}

type GroupsConfig struct {
	Op        string        `json:"op"`
	GroupSpec gocloak.Group `json:"groupSpec"`
}

type ClientConfigSpec struct {
	Name        string             `json:"name"`
	Scopes      []ScopeConfig      `json:"scopes"`
	Resources   []ResourceConfig   `json:"resources"`
	Policies    []PolicyConfig     `json:"policies"`
	Permissions []PermissionConfig `json:"permissions"`
}

type ScopeConfig struct {
	Op        string                      `json:"op"`
	ScopeSpec gocloak.ScopeRepresentation `json:"scopeSpec"`
}

type ResourceConfig struct {
	Op           string                         `json:"op"`
	ResourceSpec gocloak.ResourceRepresentation `json:"resourceSpec"`
}
type PolicyConfig struct {
	Op         string                       `json:"op"`
	PolicySpec gocloak.PolicyRepresentation `json:"policySpec"`
}
type PermissionConfig struct {
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
