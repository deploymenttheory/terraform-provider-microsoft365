package graphBetaGroupAppRoleAssignment_test

import (
	"testing"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/destroy"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/testlog"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaGroupAppRoleAssignment "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/groups/graph_beta/group_app_role_assignment"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const (
	resourceTypeAcc = graphBetaGroupAppRoleAssignment.ResourceName
)

var (
	testResource = graphBetaGroupAppRoleAssignment.GroupAppRoleAssignmentTestResource{}
)

func testAccGroupAppRoleAssignmentConfig_Minimal() string {
	config, _ := helpers.ParseHCLFile("tests/terraform/acceptance/resource_minimal.tf")
	return config
}

func testAccGroupAppRoleAssignmentConfig_Maximal() string {
	config, _ := helpers.ParseHCLFile("tests/terraform/acceptance/resource_maximal.tf")
	return config
}

// TestAccGroupAppRoleAssignmentResource_Minimal tests minimal configuration
func TestAccGroupAppRoleAssignmentResource_Minimal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			resourceTypeAcc,
			20*time.Second,
		),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceTypeAcc, "Creating app role assignment with minimal config")
				},
				Config: testAccGroupAppRoleAssignmentConfig_Minimal(),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("group app role assignment", 20*time.Second)
						time.Sleep(20 * time.Second)
						return nil
					},
					check.That(resourceTypeAcc+".minimal").ExistsInGraph(testResource),
					check.That(resourceTypeAcc+".minimal").Key("id").Exists(),
					check.That(resourceTypeAcc+".minimal").Key("target_group_id").Exists(),
					check.That(resourceTypeAcc+".minimal").Key("resource_object_id").Exists(),
					check.That(resourceTypeAcc+".minimal").Key("app_role_id").HasValue("ea358ccf-c4a8-48ac-8b94-2558ae2f7a5c"),
					check.That(resourceTypeAcc+".minimal").Key("principal_display_name").Exists(),
					check.That(resourceTypeAcc+".minimal").Key("resource_display_name").HasValue("MileIQ Admin Center"),
					check.That(resourceTypeAcc+".minimal").Key("principal_type").HasValue("Group"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceTypeAcc, "Importing app role assignment with minimal config")
				},
				ResourceName:      resourceTypeAcc + ".minimal",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// TestAccGroupAppRoleAssignmentResource_Maximal tests maximal configuration
func TestAccGroupAppRoleAssignmentResource_Maximal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			resourceTypeAcc,
			20*time.Second,
		),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceTypeAcc, "Creating app role assignment with maximal config")
				},
				Config: testAccGroupAppRoleAssignmentConfig_Maximal(),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("group app role assignment", 20*time.Second)
						time.Sleep(20 * time.Second)
						return nil
					},
					check.That(resourceTypeAcc+".maximal").ExistsInGraph(testResource),
					check.That(resourceTypeAcc+".maximal").Key("id").Exists(),
					check.That(resourceTypeAcc+".maximal").Key("target_group_id").Exists(),
					check.That(resourceTypeAcc+".maximal").Key("resource_object_id").Exists(),
					check.That(resourceTypeAcc+".maximal").Key("app_role_id").HasValue("ea358ccf-c4a8-48ac-8b94-2558ae2f7a5c"),
					check.That(resourceTypeAcc+".maximal").Key("principal_display_name").Exists(),
					check.That(resourceTypeAcc+".maximal").Key("resource_display_name").HasValue("MileIQ Admin Center"),
					check.That(resourceTypeAcc+".maximal").Key("principal_type").HasValue("Group"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceTypeAcc, "Importing app role assignment with maximal config")
				},
				ResourceName:      resourceTypeAcc + ".maximal",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
