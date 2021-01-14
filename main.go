package main

import (
	"encoding/json"
	"io/ioutil"
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
	mode := "execute"

	if mode == "diff" {
		configFile := "product-service-sec-conf.json"
		var config modules.KeycloakConfig
		loadConfig(configFile, &config)
		ctx := createDiffCtx(config)
		diffConfig := modules.KeycloakOpsConfig{
			ClientConfig: modules.ClientConfigOpSpec{
				Name: config.ClientConfig.Name,
			},
		}
		handlers := make([]modules.DiffHandler, len(modules.DiffModules)+1)
		for _, handler := range modules.DiffModules {
			handlers[handler.Order()] = handler
		}
		for _, handler := range handlers {
			if handler != nil {
				handler.Diff(&ctx, &diffConfig)
			}
		}
		opsConfig, _ := json.MarshalIndent(diffConfig, "", "   ")
		log.Info().Msg(string(opsConfig))
	}
	if mode == "execute" {
		// configFileName := "product-service-sec.json"
		configFileName := "change-log-ops.json"
		var keycloakConfig modules.KeycloakOpsConfig
		loadConfig(configFileName, &keycloakConfig)
		context := createOpConfigCtx(keycloakConfig)
		handlers := make([]modules.ConfigurationHandler, len(modules.Modules))
		for _, handler := range modules.Modules {
			handlers[handler.Order()] = handler
		}
		for _, handler := range handlers {
			handler.Apply(&context)
		}
	}
}

func loadConfig(fileName string, target interface{}) {
	jsonFile, err := os.Open(fileName)
	if err != nil {
		panic("Something wrong with config file. " + err.Error())
	}
	defer jsonFile.Close()
	fileContent, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(fileContent, target)
}

func createOpConfigCtx(keycloakConfig modules.KeycloakOpsConfig) modules.ConfigurationContext {
	clientService := clients.New(access.KeycloakConnection())
	client, err := clientService.FindClientByName(keycloakConfig.ClientConfig.Name)
	if err != nil {
		log.Err(err).Str("client", keycloakConfig.ClientConfig.Name).Msg("Problem with locating client ")
		os.Exit(1)
	}
	context := modules.ConfigurationContext{
		Config: &keycloakConfig,
		Client: client,
	}
	return context
}

func createDiffCtx(keycloakConfig modules.KeycloakConfig) modules.DiffGenCtx {
	clientService := clients.New(access.KeycloakConnection())
	client, err := clientService.FindClientByName(keycloakConfig.ClientConfig.Name)
	if err != nil {
		log.Err(err).Str("client", keycloakConfig.ClientConfig.Name).Msg("Problem with locating client ")
		os.Exit(1)
	}
	context := modules.DiffGenCtx{
		Config: &keycloakConfig,
		Client: client,
	}
	return context
}
