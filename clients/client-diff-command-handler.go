package clients

import (
	"encoding/json"
	"io/ioutil"
	"keycloak-manager/access"
	"keycloak-manager/modules"
	"keycloak-manager/tools"

	"github.com/Nerzal/gocloak/v7"
	"github.com/rs/zerolog/log"
)

func HandleClientDiffCommand(cfgFileName string, changesFileName string, keycloak *access.KeycloakContext) {
	var config modules.ClientDeclaration
	tools.LoadConfigFile(cfgFileName, &config)
	ctx := createClientDiffCtx(config, keycloak)
	diffConfig := modules.ClientChanges{}
	if ctx.ClientOp.Op == "NONE" {
		diffConfig.Client = modules.ClientOp{
			ClientSpec: gocloak.Client{
				ClientID: ctx.ClientOp.ClientSpec.ClientID,
			},
			Op: "NONE",
		}
	} else {
		diffConfig.Client = *ctx.ClientOp
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
	opsConfig, err := json.MarshalIndent(diffConfig, "", "   ")
	if err != nil {
		log.Err(err).Msg("Cannot serialize config changes to json")
	}
	err = ioutil.WriteFile(changesFileName, opsConfig, 0644)
	if err != nil {
		log.Err(err).Str("fileName", changesFileName).Msg("Cannot write config changes to file")
	}
}

func createClientDiffCtx(config modules.ClientDeclaration, keycloak *access.KeycloakContext) modules.ClientDiffContext {
	clientService := New(keycloak)
	client, err := clientService.FindClientByName(*config.Client.ClientID)
	var clientOp modules.ClientOp
	if err != nil {
		log.Info().Str("client", *config.Client.ClientID).Msg("Client does not exists. Creating new")
		clientOp = modules.ClientOp{
			Op:         "ADD",
			ClientSpec: config.Client,
		}
	} else {
		clientOp = modules.ClientOp{
			Op:         "NONE",
			ClientSpec: *client,
		}
	}
	context := modules.ClientDiffContext{
		Declaration: &config,
		ClientOp:    &clientOp,
	}
	return context
}
