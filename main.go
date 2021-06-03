package main

import (
	"fmt"
	"keycloak-manager/access"
	"keycloak-manager/clients"
	"keycloak-manager/groups"
	"keycloak-manager/model"
	"keycloak-manager/modules"
	"keycloak-manager/permissions"
	"keycloak-manager/policies"
	"keycloak-manager/resources"
	"keycloak-manager/roles"
	"keycloak-manager/scopes"
	"os"
	"syscall"

	"github.com/alecthomas/kong"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/ssh/terminal"
)

const version string = "1.1"

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})

	if os.Args[0] == "--version" {
		log.Printf("Keycloak Manager Version %s", version)
		return
	}
	var params model.CLI
	ctx := kong.Parse(&params)
	switch ctx.Command() {
	case "version":
		log.Info().Msg("Keycloak Manager Version " + version)
		break
	case "client":
		if params.Pass == "" {
			fmt.Printf("Enter password for user %s: ", params.User)
			password, _ := terminal.ReadPassword(int(syscall.Stdin))
			pass := string(password)
			if len(pass) < 1 {
				panic("Password for user missing")
			}
			params.Pass = pass
		}
		kecloak := bootstrapServices(&params)

		file := params.Client.File
		mode := params.Client.Mode
		if mode == "diff" {
			output := params.Client.Output
			clients.HandleClientDiffCommand(file, output, kecloak)
			log.Info().Msg("Config processing and diff generation finished")
		}
		if mode == "apply" {
			clients.ApplyClientChanged(file, kecloak)
			log.Info().Msg("Configuration diff application finished")
		}
	}
}

func bootstrapServices(params *model.CLI) *access.KeycloakContext {
	kecloak := access.NewConnection(params)
	scopes.InitializeService(kecloak, modules.Modules, modules.DiffModules)
	roles.InitializeService(kecloak, modules.Modules, modules.DiffModules)
	resources.InitializeService(kecloak, modules.Modules, modules.DiffModules)
	policies.InitializeService(kecloak, modules.Modules, modules.DiffModules)
	permissions.InitializeService(kecloak, modules.Modules, modules.DiffModules)
	groups.InitializeService(kecloak, modules.Modules, modules.DiffModules)
	return kecloak
}
