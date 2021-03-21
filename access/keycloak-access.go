package access

import (
	"context"

	"github.com/Nerzal/gocloak/v7"
)

type KeycloakContext struct {
	Client gocloak.GoCloak
	Token  *gocloak.JWT
	Ctx    context.Context
}

func KeycloakConnection() *KeycloakContext {
	client := gocloak.NewClient("http://swarm-local:9723")
	ctx := context.Background()
	token, err := client.LoginAdmin(ctx, "admin", "password", "master")
	if err != nil {
		panic("Something wrong with creds or address. " + err.Error())
	}
	return &KeycloakContext{
		Client: client,
		Token:  token,
		Ctx:    ctx,
	}
}
