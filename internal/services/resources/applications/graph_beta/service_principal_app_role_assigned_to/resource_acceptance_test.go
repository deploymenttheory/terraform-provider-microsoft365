package graphBetaServicePrincipalAppRoleAssignedTo_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/destroy"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/testlog"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaServicePrincipalAppRoleAssignedTo "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/applications/graph_beta/service_principal_app_role_assigned_to"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

var (
	// Resource type name from the resource package
	accResourceType = graphBetaServicePrincipalAppRoleAssignedTo.ResourceName

	// testResource is the test resource implementation for app role assignments
	testResource = graphBetaServicePrincipalAppRoleAssignedTo.ServicePrincipalAppRoleAssignedToTestResource{}
)

// TestAccServicePrincipalAppRoleAssignedToResource_ToServicePrincipal tests assigning an app role
// to a regular service principal created via the azuread provider (fallback when this provider
// doesn't have the required resource type)
func TestAccServicePrincipalAppRoleAssignedToResource_ToServicePrincipal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			accResourceType,
			30*time.Second,
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
					testlog.StepAction(accResourceType, "Creating app role assignment to service principal")
				},
				Config: testAccConfigToServicePrincipal(),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("app role assignment", 10*time.Second)
						time.Sleep(10 * time.Second)
						return nil
					},
					check.That(accResourceType+".user_read_all").ExistsInGraph(testResource),
					check.That(accResourceType+".user_read_all").Key("id").Exists(),
					check.That(accResourceType+".user_read_all").Key("resource_object_id").Exists(),
					check.That(accResourceType+".user_read_all").Key("app_role_id").HasValue("df021288-bdef-4463-88db-98f22de89214"),
					check.That(accResourceType+".user_read_all").Key("target_service_principal_object_id").Exists(),
					check.That(accResourceType+".user_read_all").Key("principal_type").HasValue("ServicePrincipal"),
					check.That(accResourceType+".user_read_all").Key("principal_display_name").Exists(),
					check.That(accResourceType+".user_read_all").Key("resource_display_name").Exists(),
					check.That(accResourceType+".user_read_all").Key("created_date_time").Exists(),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(accResourceType, "Importing app role assignment to service principal")
				},
				ResourceName: accResourceType + ".user_read_all",
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[accResourceType+".user_read_all"]
					if !ok {
						return "", fmt.Errorf("resource not found: %s", accResourceType+".user_read_all")
					}
					resourceObjectID := rs.Primary.Attributes["resource_object_id"]
					id := rs.Primary.Attributes["id"]
					return fmt.Sprintf("%s/%s", resourceObjectID, id), nil
				},
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"timeouts",
				},
			},
		},
	})
}

// Test configuration functions using mocks.LoadTerraformConfigFile and acceptance.ConfiguredM365ProviderBlock
func testAccConfigToServicePrincipal() string {
	config := mocks.LoadTerraformConfigFile("app_role_assignment_to_service_principle.tf")
	return acceptance.ConfiguredM365ProviderBlock(config)
}
