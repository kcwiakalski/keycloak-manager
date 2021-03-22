package groups

import (
	"context"
	"errors"
	"fmt"
	"keycloak-tools/access"
	"keycloak-tools/modules"
	"strings"

	"github.com/Nerzal/gocloak/v7"
	"github.com/rs/zerolog/log"
)

type GroupService struct {
	client gocloak.GoCloak
	ctx    context.Context
	token  string
}

func new(ctx *access.KeycloakContext) *GroupService {
	return &GroupService{
		client: ctx.Client,
		ctx:    ctx.Ctx,
		token:  ctx.Token.AccessToken,
	}
}

var groupsService *GroupService

// groupsService := New(keycloak)
func (s *GroupService) Apply(keycloakConfig *modules.ClientChangeContext) error {
	// var finalError error
	// for _, group := range keycloakConfig.Config.Groups {
	// 	err := groupsService.AddGroup(&group.GroupSpec)
	// 	if err != nil {
	// 		finalError = err
	// 	}
	// }
	// return finalError
	return errors.New("not implemented yet")
}

func (s *GroupService) Order() int {
	return 0
}

func init() {
	groupsService = new(modules.Keycloak)
	modules.Modules["groups"] = groupsService
	modules.DiffModules["groups"] = groupsService
}
func (s *GroupService) getGroupByName(groupName string) *gocloak.Group {
	first := 0
	max := 20
	fullSearch := true
	var pathElements []string
	nestedSearch := false
	if strings.ContainsAny(groupName, "/") {
		pathElements = strings.Split(groupName, "/")
		nestedSearch = true
	}
	searchTerm := groupName
	if nestedSearch {
		searchTerm = pathElements[0]
	}
	params := gocloak.GetGroupsParams{
		First:  &first,
		Max:    &max,
		Search: &searchTerm,
		Full:   &fullSearch,
	}
	existingGroups, _ := s.client.GetGroups(s.ctx, s.token, modules.REALM_NAME, params)
	for _, group := range existingGroups {
		if *group.Name == groupName {
			return group
		}
	}
	return nil
}
func (s *GroupService) getGroupByPath(groupPath *gocloak.Group) *gocloak.Group {
	first := 0
	max := 20
	fullSearch := true
	var pathElements []string
	searchTerm := groupPath.Path
	if strings.ContainsAny(*groupPath.Path, "/") {
		pathElements = strings.Split(*groupPath.Path, "/")
		searchTerm = &pathElements[1]
	}
	params := gocloak.GetGroupsParams{
		First:  &first,
		Max:    &max,
		Search: searchTerm,
		Full:   &fullSearch,
	}
	existingGroups, _ := s.client.GetGroups(s.ctx, s.token, modules.REALM_NAME, params)
	for _, group := range existingGroups {
		matchedGroup := groupsService.findDirectParentGroup(groupPath, group)
		if matchedGroup != nil {
			return matchedGroup
		}
	}
	return nil
}

func (s *GroupService) findDirectParentGroup(group *gocloak.Group, topLevelGroup *gocloak.Group) *gocloak.Group {
	pathElements := strings.Split(*group.Path, "/")
	nestedGroup := topLevelGroup
	if pathElements[0] != "" {
		return nil
	}
	if pathElements[1] != *topLevelGroup.Name {
		return nil
	}
	index := 2
	for nestedGroup != nil && index < len(pathElements)-1 {
		var matchedGroup *gocloak.Group
		for _, innerGroup := range *nestedGroup.SubGroups {
			if *innerGroup.Name == pathElements[index] {
				matchedGroup = &innerGroup
			}
		}
		nestedGroup = matchedGroup
		index++
	}
	return nestedGroup
}

func (s *GroupService) groupExists(groupName string) bool {
	group := s.getGroupByName(groupName)
	if group != nil {
		return true
	}
	return false
}

func (s *GroupService) AddGroup(group *gocloak.Group) error {
	pathParts := strings.Count(*group.Path, "/")
	if pathParts == 1 && strings.TrimPrefix(*group.Path, "/") == *group.Name {
		if !s.groupExists(*group.Name) {
			groupId, err := s.client.CreateGroup(s.ctx, s.token, "products", *group)
			if err != nil {
				log.Err(err).Str("name", *group.Name).Msg("Error creating group")
				return err
			} else {
				log.Info().Str("name", *group.Name).Str("id", groupId).Msg("Group created")
			}
		} else {
			log.Warn().Str("name", *group.Name).Msg("Group already exists")
		}
	} else if pathParts > 1 {
		mainGroup := s.getGroupByName(strings.Split(*group.Path, "/")[1])
		directParent := s.findDirectParentGroup(group, mainGroup)
		for _, childGroups := range *directParent.SubGroups {
			if childGroups.Name == group.Name {
				log.Warn().Str("group", *group.Name).Str("parent", *directParent.Name).Msg("Group with parent is already defined")
				return fmt.Errorf("Group %s with parent%s is already defined", *group.Name, *directParent.Name)
			}
		}
		groupId, err := s.client.CreateChildGroup(s.ctx, s.token, "products", *directParent.ID, *group)
		if err != nil {
			log.Err(err).Str("group", *group.Name).Str("parent", *directParent.Name).Msg("Cannot create child group in parent")
			return fmt.Errorf("Cannot create child group %s in parent %s. %s", *group.Name, *directParent.Name, err.Error())
		} else {
			log.Info().Str("name", *group.Name).Str("parent", *directParent.Name).Str("id", groupId).Msg("Child group created")
		}
	} else {
		log.Error().Str("group", *group.Name).Msg("Invalid group definition")
		return fmt.Errorf("Invalid group definition for name %s", *group.Name)
	}
	return nil
}

func (s *GroupService) Diff(declaration *modules.ClientDiffContext, changes *modules.ClientChanges) error {
	var ops []modules.GroupsOp = make([]modules.GroupsOp, 0)
	for _, expectedGroup := range declaration.Declaration.Groups {
		existingGroup := s.getGroupByPath(&expectedGroup)
		if existingGroup == nil {
			ops = append(ops, modules.GroupsOp{
				Op:        "ADD",
				GroupSpec: expectedGroup,
			})
		}
	}
	if len(ops) > 0 {
		changes.Groups = ops
	}
	return nil
}
