package groups

import (
	"encoding/json"
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

const fullGroupJson = `  {
    "id": "1c447424-b540-4940-a723-a8d003eba228",
    "name": "employee-standard",
    "path": "/employee-standard",
    "subGroups": [
      {
        "id": "e693d5a6-745e-4392-9496-09d0103fea13",
        "name": "employee-data-admin",
        "path": "/employee-standard/employee-data-admin",
        "subGroups": [
          {
            "id": "62871dbc-4341-442e-8b2f-850348c60c08",
            "name": "employee-system-admin",
            "path": "/employee-standard/employee-data-admin/employee-system-admin",
            "subGroups": []
          }
        ]
      }
    ]
  }`

func TestShouldNotFindGroupInTopLevelGroup(t *testing.T) {
	//given
	var group gocloak.Group
	groupJson := `{
		"name": "employee-system-superadmin",
		"path": "/employee-standard/employee-data-admin/employee-system-admin/super-admin"
	 }`
	json.Unmarshal([]byte(groupJson), &group)

	var fullGroup gocloak.Group
	json.Unmarshal([]byte(fullGroupJson), &fullGroup)
	//when
	service := GroupService{}

	matchedGroup := service.findGroupInTopLevelGroup(&group, &fullGroup)

	//then
	if matchedGroup != nil {
		t.Errorf("Should not find group")
	}
}

func TestShouldFindGroupInTopLevelGroup(t *testing.T) {
	//given
	var group gocloak.Group
	groupJson := `{
		"name": "employee-system-admin",
		"path": "/employee-standard/employee-data-admin/employee-system-admin"
	 }`
	json.Unmarshal([]byte(groupJson), &group)

	var fullGroup gocloak.Group
	json.Unmarshal([]byte(fullGroupJson), &fullGroup)
	//when
	service := GroupService{}

	matchedGroup := service.findGroupInTopLevelGroup(&group, &fullGroup)

	//then
	if matchedGroup == nil {
		t.Errorf("Should not find group")
	}
}

func TestShouldNotFindGroupWhenTopLevelPathDontMatch(t *testing.T) {
	//given
	var group gocloak.Group
	groupJson := `{
		"name": "employee-system-admin",
		"path": "/some-other-group/employee-data-admin/employee-system-admin"
	 }`
	json.Unmarshal([]byte(groupJson), &group)

	var fullGroup gocloak.Group
	json.Unmarshal([]byte(fullGroupJson), &fullGroup)
	//when
	service := GroupService{}

	matchedGroup := service.findGroupInTopLevelGroup(&group, &fullGroup)

	//then
	if matchedGroup != nil {
		t.Errorf("Should not find group")
	}
}
