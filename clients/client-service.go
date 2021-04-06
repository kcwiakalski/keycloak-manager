package clients

import (
	"context"
	"fmt"
	"keycloak-manager/access"

	"github.com/Nerzal/gocloak/v7"
	"github.com/rs/zerolog/log"
)

type ClientService struct {
	client gocloak.GoCloak
	ctx    context.Context
	token  string
	realm  string
}

func New(ctx *access.KeycloakContext) *ClientService {
	return &ClientService{
		client: ctx.Client,
		ctx:    ctx.Ctx,
		token:  ctx.Token.AccessToken,
		realm:  ctx.Realm,
	}
}
func (s *ClientService) FindClientByName(name string) (*gocloak.Client, error) {
	clients, err := s.client.GetClients(s.ctx, s.token, s.realm, gocloak.GetClientsParams{})
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

func (s *ClientService) CreateClient(client gocloak.Client) (string, error) {
	id, err := s.client.CreateClient(s.ctx, s.token, s.realm, client)
	if err != nil {
		log.Fatal().Err(err).Str("clientName", *client.ClientID).Msg("Cannot create client")
	} else {
		client.Name = &id
	}
	return id, nil
}
