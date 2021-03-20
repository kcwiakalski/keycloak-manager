package main

import (
	"keycloak-tools/access"
	"keycloak-tools/clients"
	_ "keycloak-tools/groups"
	"keycloak-tools/modules"
	_ "keycloak-tools/permissions"
	_ "keycloak-tools/policies"
	_ "keycloak-tools/resources"
	_ "keycloak-tools/scopes"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// var groupsService *groups.GroupService

// var scopeService *scopes.ScopeService

// var resourceService *resources.ResourceService
// var policyService *policies.PolicyService
// var permissionService *permissions.PermissionService

func init() {
}

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	command := "client"
	mode := "diff"
	configFile := "product-service-sec-conf.json"

	if command == "client" {
		if mode == "diff" {
			clients.HandleClientDiffCommand(configFile)
		}
		if mode == "execute" {
			// configFileName := "product-service-sec.json"
			// configFileName := "change-log-ops.json"
			// var declaration modules.ClientConfig
			// loadConfig(configFileName, &declaration)
			// context := createOpConfigCtx(declaration)
			// handlers := make([]modules.ConfigurationHandler, len(modules.Modules))
			// for _, handler := range modules.Modules {
			// 	handlers[handler.Order()] = handler
			// }
			// for _, handler := range handlers {
			// 	handler.Apply(&context)
			// }
		}
	}
}

func createOpConfigCtx(config modules.KeycloakOpsConfig) modules.ConfigurationContext {
	clientService := clients.New(access.KeycloakConnection())
	client, err := clientService.FindClientByName(*config.ClientConfig.Declaration.Name)
	if err != nil {
		log.Info().Str("client", *config.ClientConfig.Declaration.Name).Msg("Client does not exists. Creating new")
		client = &config.ClientConfig.Declaration
	}
	context := modules.ConfigurationContext{
		Config: &config,
		Client: client,
	}
	return context
}
