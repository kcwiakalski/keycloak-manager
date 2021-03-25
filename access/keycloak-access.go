package access

import (
	"context"
	"fmt"
	"keycloak-tools/model"

	"github.com/Nerzal/gocloak/v7"
)

type KeycloakContext struct {
	Client gocloak.GoCloak
	Token  *gocloak.JWT
	Ctx    context.Context
}

func KeycloakConnection() *KeycloakContext {
	cli := model.CLI
	client := gocloak.NewClient(fmt.Sprintf("http://%s:%d", cli.Host, cli.Port))
	ctx := context.Background()
	token, err := client.LoginAdmin(ctx, cli.User, cli.Pass, cli.Realm)
	if err != nil {
		panic("Something wrong with creds or address. " + err.Error())
	}
	return &KeycloakContext{
		Client: client,
		Token:  token,
		Ctx:    ctx,
	}
}
