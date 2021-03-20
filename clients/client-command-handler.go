package clients

import (
	"encoding/json"
	"keycloak-tools/access"
	"keycloak-tools/modules"
	"keycloak-tools/tools"

	"github.com/Nerzal/gocloak/v7"
	"github.com/rs/zerolog/log"
)

func HandleClientDiffCommand(cfgFileName string) {
	var config modules.KeycloakClientConfig
	tools.LoadConfigFile(cfgFileName, &config)
	ctx := createClientDiffCtx(config)
	diffConfig := modules.KeycloakOpsConfig{
		ClientConfig: modules.ClientConfigOpSpec{},
	}
	if ctx.ClientOp.Op != "NONE" {
		diffConfig.ClientConfig = modules.ClientConfigOpSpec{
			Declaration: ctx.ClientOp.Client,
			Op:          ctx.ClientOp.Op,
		}
	} else {
		diffConfig.ClientConfig = modules.ClientConfigOpSpec{
			Declaration: gocloak.Client{
				ClientID: ctx.ClientOp.Client.ClientID,
			},
			Op: "NONE",
		}
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

func createClientDiffCtx(config modules.KeycloakClientConfig) modules.KeycloakClientDiffGenCtx {
	clientService := New(access.KeycloakConnection())
	client, err := clientService.FindClientByName(*config.Definition.ClientID)
	var clientOp modules.ClientOp
	if err != nil {
		log.Info().Str("client", *config.Definition.ClientID).Msg("Client does not exists. Creating new")
		clientOp = modules.ClientOp{
			Op:     "ADD",
			Client: config.Definition,
		}
	} else {
		clientOp = modules.ClientOp{
			Op:     "NONE",
			Client: *client,
		}
	}
	context := modules.KeycloakClientDiffGenCtx{
		Config:   &config,
		ClientOp: &clientOp,
	}
	return context
}
