package clients

import (
	"context"
	"fmt"
	"keycloak-tools/access"

	"github.com/Nerzal/gocloak/v7"
)

type ClientService struct {
	client gocloak.GoCloak
	ctx    context.Context
	token  string
}

func New(ctx *access.KeycloakContext) *ClientService {
	return &ClientService{
		client: ctx.Client,
		ctx:    ctx.Ctx,
		token:  ctx.Token.AccessToken,
	}
}
func (s *ClientService) FindClientByName(name string) (*gocloak.Client, error) {
	clients, err := s.client.GetClients(s.ctx, s.token, "products", gocloak.GetClientsParams{})
	if err != nil {
		return nil, err
	}
	for _, client := range clients {
		if *client.ClientID == name {
			return client, nil
		}
	}
	return nil, fmt.Errorf("Cannot find client with name %s", name)
}
