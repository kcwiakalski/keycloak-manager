package main

import (
	"encoding/json"
	"io/ioutil"
	"keycloak-tools/access"
	"keycloak-tools/clients"
	_ "keycloak-tools/groups"
	"keycloak-tools/modules"
	"keycloak-tools/permissions"
	"keycloak-tools/policies"
	_ "keycloak-tools/resources"
	_ "keycloak-tools/scopes"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// var groupsService *groups.GroupService

var keycloak *access.KeycloakContext

// var scopeService *scopes.ScopeService
var clientService *clients.ClientService

// var resourceService *resources.ResourceService
var policyService *policies.PolicyService
var permissionService *permissions.PermissionService

func init() {
	keycloak = access.KeycloakConnection()
	// groupsService = groups.New(keycloak)
	// scopeService = scopes.New(keycloak)
	clientService = clients.New(keycloak)
	// resourceService = resources.New(keycloak)
	policyService = policies.New(keycloak)
	permissionService = permissions.New(keycloak)
}

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	jsonFile, err := os.Open("product-service-sec.json")
	if err != nil {
		panic("Something wrong with config file. " + err.Error())
	}
	defer jsonFile.Close()
	fileContent, _ := ioutil.ReadAll(jsonFile)
	var keycloakConfig modules.KeycloakConfig
	json.Unmarshal(fileContent, &keycloakConfig)

	// for _, group := range keycloakConfig.Groups {
	// 	err := groupsService.AddGroup(&group.GroupSpec)
	// 	if err != nil {
	// 		log.Println(err)
	// 	}
	// }
	client, err := clientService.FindClientByName(keycloakConfig.ClientConfig.Name)
	context := modules.ConfigurationContext{
		Config: &keycloakConfig,
		Client: client,
	}

	handlers := make([]modules.ConfigurationHandler, len(modules.Modules))
	for _, handler := range modules.Modules {
		handlers[handler.Order()] = handler
	}
	for _, handler := range handlers {
		handler.Apply(&context)
	}

	if err != nil {
		log.Printf("Error locating client config, %s", err.Error())
	} else {
		// for _, scope := range keycloakConfig.ClientConfig.Scopes {
		// 	err = scopeService.AddScope(*client.ID, &scope.ScopeSpec)
		// 	if err != nil {
		// 		log.Printf("Problem with scope %s creation. %s", *scope.ScopeSpec.ID, err.Error())
		// 	}
		// }
		// for _, resource := range keycloakConfig.ClientConfig.Resources {
		// 	err = resourceService.AddResource(*client.ID, resource.ResourceSpec)
		// 	if err != nil {
		// 		log.Printf("Problem with resource %s creation. %s", *resource.ResourceSpec.Name, err.Error())
		// 	}
		// }
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
