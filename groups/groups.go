package groups

import (
	"context"
	"fmt"
	"keycloak-tools/access"
	"log"
	"strings"

	"github.com/Nerzal/gocloak/v7"
)

type GroupService struct {
	client gocloak.GoCloak
	ctx    context.Context
	token  string
}

func New(ctx *access.KeycloakContext) *GroupService {
	return &GroupService{
		client: ctx.Client,
		ctx:    ctx.Ctx,
		token:  ctx.Token.AccessToken,
	}
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
	existingGroups, _ := s.client.GetGroups(s.ctx, s.token, "products", params)
	for _, group := range existingGroups {
		if *group.Name == groupName {
			return group
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
				log.Printf("Error creating group %s. %s", *group.Name, err.Error())
				return err
			} else {
				log.Printf("Group %s created with id %s", *group.Name, groupId)
			}
		} else {
			log.Printf("Group with name %s already exists", *group.Name)
		}
	} else if pathParts > 1 {
		mainGroup := s.getGroupByName(strings.Split(*group.Path, "/")[1])
		directParent := s.findDirectParentGroup(group, mainGroup)
		for _, childGroups := range *directParent.SubGroups {
			if childGroups.Name == group.Name {
				return fmt.Errorf("Group %s with parent%s is already defined", *group.Name, *directParent.Name)
			}
		}
		groupId, err := s.client.CreateChildGroup(s.ctx, s.token, "products", *directParent.ID, *group)
		if err != nil {
			return fmt.Errorf("Cannor create child group %s in parent %s. %s", *group.Name, *directParent.Name, err.Error())
		} else {
			log.Printf("Child group %s created in parent %s with id %s", *group.Name, *directParent.Name, groupId)
		}
	} else {
		return fmt.Errorf("Invalid group definition for name %s", *group.Name)
	}
	return nil
}
