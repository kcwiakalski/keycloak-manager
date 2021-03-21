package main

import (
	"keycloak-tools/clients"
	_ "keycloak-tools/groups"
	_ "keycloak-tools/permissions"
	_ "keycloak-tools/policies"
	_ "keycloak-tools/resources"
	_ "keycloak-tools/scopes"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	command := "client"
	// mode := "diff"
	mode := "execute"
	configFile := "product-service-sec-conf.json"
	configFileChange := "product-service-sec-change.json"

	if command == "client" {
		if mode == "diff" {
			clients.HandleClientDiffCommand(configFile, configFileChange)
		}
		if mode == "execute" {
			clients.ApplyClientChanged(configFileChange)
		}
	}
}
