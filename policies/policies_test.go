package policies

import (
	"encoding/json"
	"log"
	"testing"

	"github.com/Nerzal/gocloak/v7"
)

func TestPolicyCreation(t *testing.T) {
	service := PolicyService{}
	log.Println(service)
}
func TestSome(t *testing.T) {
	name := "some-name"
	path := "/sample2"
	extend := true
	groups := []gocloak.GroupDefinition{{Path: &path, ExtendChildren: &extend}}
	var policy gocloak.PolicyRepresentation = gocloak.PolicyRepresentation{
		Name: &name,
		GroupPolicyRepresentation: gocloak.GroupPolicyRepresentation{
			Groups: &groups,
		},
	}
	x0, _ := json.MarshalIndent(policy, "", "\t")
	t.Fatalf("%s", string(x0))
}
