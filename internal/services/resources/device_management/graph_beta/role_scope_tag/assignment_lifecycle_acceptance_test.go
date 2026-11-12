package graphBetaRoleScopeTag_test

import (
	"testing"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/destroy"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceRoleScopeTag_04_AssignmentLifecycle(t *testing.T) {
	address := resourceType + ".lifecycle"
	var tagID string
	var steps []resource.TestStep
	for i, phase := range []struct{ name, count string }{
		{"unassigned", ""}, {"maximal", "2"}, {"minimal", "1"},
		{"empty", "0"}, {"maximal", "2"}, {"unassigned", ""},
	} {
		config := loadAcceptanceTestTerraform("resource_lifecycle_" + phase.name + ".tf")
		var checks []resource.TestCheckFunc
		if i == 0 {
			checks = append(checks, resource.TestCheckResourceAttrWith(address, "id", func(value string) error { tagID = value; return nil }))
		} else {
			checks = append(checks, resource.TestCheckResourceAttrPtr(address, "id", &tagID))
		}
		if phase.count == "" {
			checks = append(checks, resource.TestCheckNoResourceAttr(address, "assignments.#"))
		} else {
			checks = append(checks, resource.TestCheckResourceAttr(address, "assignments.#", phase.count))
		}
		steps = append(steps, resource.TestStep{Config: config, Check: resource.ComposeTestCheckFunc(checks...)})
		// Graph does not retain whether an empty assignment set was explicit or omitted.
		if phase.name != "empty" {
			steps = append(steps, resource.TestStep{ResourceName: address, ImportState: true, ImportStateVerify: true})
		}
		steps = append(steps, resource.TestStep{Config: config, PlanOnly: true})
	}
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		ExternalProviders:        map[string]resource.ExternalProvider{"random": {Source: "hashicorp/random", VersionConstraint: constants.ExternalProviderRandomVersion}},
		CheckDestroy: destroy.CheckDestroyedTypesFunc(30*time.Second,
			destroy.ResourceTypeMapping{ResourceType: resourceType, TestResource: testResource},
			destroy.ResourceTypeMapping{ResourceType: groupResourceType, TestResource: groupTestResource},
		),
		Steps: steps,
	})
}
