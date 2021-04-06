package clients

import (
	"keycloak-manager/access"
	"keycloak-manager/modules"
	"keycloak-manager/tools"

	"github.com/rs/zerolog/log"
)

func ApplyClientChanged(configFileName string, keycloak *access.KeycloakContext) {
	var declaration modules.ClientChanges
	tools.LoadConfigFile(configFileName, &declaration)
	context := createOpConfigCtx(declaration, keycloak)
	handlers := make([]modules.ConfigurationHandler, len(modules.Modules))
	for _, handler := range modules.Modules {
		handlers[handler.Order()] = handler
	}
	for _, handler := range handlers {
		handler.Apply(&context)
	}
}
func createOpConfigCtx(config modules.ClientChanges, keycloak *access.KeycloakContext) modules.ClientChangeContext {
	clientService := New(keycloak)
	client, err := clientService.FindClientByName(*config.Client.ClientSpec.ClientID)
	if err != nil {
		log.Info().Str("client", *config.Client.ClientSpec.ClientID).Msg("Client does not exists. Creating new")
		clientId, _ := clientService.CreateClient(config.Client.ClientSpec)
		client = &config.Client.ClientSpec
		client.ID = &clientId
	}
	context := modules.ClientChangeContext{
		Changes: &config,
		Client:  client,
	}
	return context
}
