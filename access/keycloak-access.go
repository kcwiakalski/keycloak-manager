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
	Realm  string
}

func NewConnection(params *model.CLI) *KeycloakContext {

	client := gocloak.NewClient(fmt.Sprintf("http://%s:%d", params.Server, params.Port))
	ctx := context.Background()
	token, err := client.LoginAdmin(ctx, params.User, params.Pass, params.Realm)
	if err != nil {
		panic("Something wrong with creds or address. " + err.Error())
	}
	context := &KeycloakContext{
		Client: client,
		Token:  token,
		Ctx:    ctx,
		Realm:  params.Realm,
	}
	user, err := client.GetUserInfo(ctx, token.AccessToken, params.Realm)
	if err != nil {
		noPermissions(params.User)
	}
	if !context.validateUserPermissions(params, user) {
		noPermissions(params.User)
	}
	return context
}

func noPermissions(user string) {
	panic("User does not have permissions to manage realm. Grant user " + user + " role 'realm-admin' on client 'realm-management'")
}

func (s *KeycloakContext) validateUserPermissions(params *model.CLI, userInfo *gocloak.UserInfo) bool {
	query := gocloak.GetUsersParams{
		Email: userInfo.Email,
	}
	users, err := s.Client.GetUsers(s.Ctx, s.Token.AccessToken, params.Realm, query)
	if err != nil || users == nil || len(users) < 1 {
		return false
	}
	for _, user := range users {
		roles, _ := s.Client.GetRoleMappingByUserID(s.Ctx, s.Token.AccessToken, params.Realm, *user.ID)
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
