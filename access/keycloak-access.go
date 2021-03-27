package access

import (
	"context"
	"fmt"
	"keycloak-manager/model"

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
	context := &KeycloakContext{
		Client: client,
		Token:  token,
		Ctx:    ctx,
	}
	user, err := client.GetUserInfo(ctx, token.AccessToken, cli.Realm)
	if err != nil {
		noPermissions()
	}
	if !context.validateUserPermissions(user) {
		noPermissions()
	}
	return context
}

func noPermissions() {
	panic("User does not have permissions to manage realm. Grant user 'realm-admin' role on client 'realm-management'")
}

func (s *KeycloakContext) validateUserPermissions(userInfo *gocloak.UserInfo) bool {
	params := gocloak.GetUsersParams{
		Email: userInfo.Email,
	}
	users, err := s.Client.GetUsers(s.Ctx, s.Token.AccessToken, model.CLI.Realm, params)
	if err != nil || users == nil || len(users) < 1 {
		return false
	}
	for _, user := range users {
		roles, _ := s.Client.GetRoleMappingByUserID(s.Ctx, s.Token.AccessToken, model.CLI.Realm, *user.ID)
		for key, mapping := range roles.ClientMappings {
			if key == "realm-management" {
				for _, role := range *mapping.Mappings {
					if *role.Name == "realm-admin" {
						return true
					}
				}
			}
		}
	}
	return false
}

//"realm-admin"
