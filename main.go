package main

import (
	"encoding/json"
	"io/ioutil"
	"keycloak-tools/access"
	"keycloak-tools/clients"
	"keycloak-tools/groups"
	"keycloak-tools/permissions"
	"keycloak-tools/policies"
	"keycloak-tools/resources"
	"keycloak-tools/scopes"
	"log"
	"os"

	"github.com/Nerzal/gocloak/v7"
)

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

var groupsService *groups.GroupService

var keycloak *access.KeycloakContext
var scopeService *scopes.ScopeService
var clientService *clients.ClientService
var resourceService *resources.ResourceService
var policyService *policies.PolicyService
var permissionService *permissions.PermissionService

func init() {
	keycloak = access.KeycloakConnection()
	groupsService = groups.New(keycloak)
	scopeService = scopes.New(keycloak)
	clientService = clients.New(keycloak)
	resourceService = resources.New(keycloak)
	policyService = policies.New(keycloak)
	permissionService = permissions.New(keycloak)
}

func main() {
	jsonFile, err := os.Open("product-service-sec.json")
	if err != nil {
		panic("Something wrong with config file. " + err.Error())
	}
	defer jsonFile.Close()
	fileContent, _ := ioutil.ReadAll(jsonFile)
	var keycloakConfig KeycloakConfig
	json.Unmarshal(fileContent, &keycloakConfig)

	for _, group := range keycloakConfig.Groups {
		err := groupsService.AddGroup(&group.GroupSpec)
		if err != nil {
			log.Println(err)
		}
	}
	client, err := clientService.FindClientByName(keycloakConfig.ClientConfig.Name)
	if err != nil {
		log.Printf("Error locating client config, %s", err.Error())
	} else {
		for _, scope := range keycloakConfig.ClientConfig.Scopes {
			err = scopeService.AddScope(*client.ID, &scope.ScopeSpec)
			if err != nil {
				log.Printf("Problem with scope %s creation. %s", *scope.ScopeSpec.ID, err.Error())
			}
		}
		for _, resource := range keycloakConfig.ClientConfig.Resources {
			err = resourceService.AddResource(*client.ID, resource.ResourceSpec)
			if err != nil {
				log.Printf("Problem with resource %s creation. %s", *resource.ResourceSpec.Name, err.Error())
			}
		}
		for _, policy := range keycloakConfig.ClientConfig.Policies {
			err = policyService.CreatePolicy(*client.ID, &policy.PolicySpec)
			if err != nil {
				log.Printf("Problem with policy %s creation. %s", *policy.PolicySpec.Name, err.Error())
			}
		}
		for _, perm := range keycloakConfig.ClientConfig.Permissions {
			err = permissionService.AddPermission(*client.ID, perm.PermSpec)
			if err != nil {
				log.Printf("Problem with permission %s creation. %s", *perm.PermSpec.Name, err.Error())
			}
		}
	}
}
