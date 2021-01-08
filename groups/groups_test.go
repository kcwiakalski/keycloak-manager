package groups

import (
	"testing"

	"github.com/Nerzal/gocloak/v7"
)

func TestFindingParentGroup(t *testing.T) {
	groupService := GroupService{}
	parent1Name := "parent1"
	parent1Path := "/" + parent1Name
	parent2Name := "parent2"
	parent2Path := parent1Path + "/" + parent2Name
	parralelParent := "parrallel-parent"
	parralelParentPath := parent1Path + "/" + parralelParent
	parent2 := gocloak.Group{
		Name: &parent2Name,
		Path: &parent2Path,
	}
	parrallelParentGroup := gocloak.Group{
		Name: &parralelParent,
		Path: &parralelParentPath,
	}

	subGroups := []gocloak.Group{parrallelParentGroup, parent2}
	parent1 := gocloak.Group{
		Name:      &parent1Name,
		Path:      &parent1Path,
		SubGroups: &subGroups,
	}
	childGroupName := "childGroup"
	childGroupPath := parent2Path + "/" + childGroupName
	searchGroup := gocloak.Group{
		Name: &childGroupName,
		Path: &childGroupPath,
	}
	searchedParentGroup := groupService.findDirectParentGroup(&searchGroup, &parent1)

	if searchedParentGroup == nil {
		t.Fatalf("Cannot find parent group")
	}
	if searchedParentGroup.Name != &parent2Name {
		t.Errorf("Parent group should be %s but was %s", parent2Name, *searchGroup.Name)
	}
}
