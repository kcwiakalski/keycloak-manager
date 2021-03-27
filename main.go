package main

import (
	"keycloak-manager/clients"
	_ "keycloak-manager/groups"
	"keycloak-manager/model"
	_ "keycloak-manager/permissions"
	_ "keycloak-manager/policies"
	_ "keycloak-manager/resources"
	_ "keycloak-manager/scopes"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})

	switch model.Ctx.Command() {
	case "client":
		file := model.CLI.Client.File
		mode := model.CLI.Client.Mode
		if mode == "diff" {
			output := model.CLI.Client.Output
			clients.HandleClientDiffCommand(file, output)
		}
		if mode == "apply" {
			clients.ApplyClientChanged(file)
		}
	}
}
