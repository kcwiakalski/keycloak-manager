package main

import (
	"keycloak-tools/clients"
	_ "keycloak-tools/groups"
	"keycloak-tools/model"
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
