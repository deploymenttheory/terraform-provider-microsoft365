package graphBetaServicePrincipalAppRoleAssignedTo_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/destroy"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/testlog"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaServicePrincipalAppRoleAssignedTo "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/applications/graph_beta/service_principal_app_role_assigned_to"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

var (
	// Resource type name from the resource package
	resourceType = graphBetaServicePrincipalAppRoleAssignedTo.ResourceName

	// testResource is the test resource implementation for app role assignments
	testResource = graphBetaServicePrincipalAppRoleAssignedTo.ServicePrincipalAppRoleAssignedToTestResource{}
)

// Helper function to load test configs from acceptance directory
func loadAcceptanceTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/acceptance/" + filename)
	if err != nil {
		panic("failed to load acceptance config " + filename + ": " + err.Error())
	}
	return acceptance.ConfiguredM365ProviderBlock(config)
}

// TestAccServicePrincipalAppRoleAssignedToResource_ToServicePrincipal tests assigning an app role
// to a regular service principal created via the azuread provider (fallback when this provider
// doesn't have the required resource type)
func TestAccServicePrincipalAppRoleAssignedToResource_ToServicePrincipal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			resourceType,
			30*time.Second,
		),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating app role assignment to service principal")
				},
				Config: loadAcceptanceTestTerraform("app_role_assignment_to_service_principal.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("app role assignment", 10*time.Second)
						time.Sleep(10 * time.Second)
						return nil
					},
					check.That(resourceType+".user_read_all").ExistsInGraph(testResource),
					check.That(resourceType+".user_read_all").Key("id").Exists(),
					check.That(resourceType+".user_read_all").Key("resource_object_id").Exists(),
					check.That(resourceType+".user_read_all").Key("app_role_id").HasValue("df021288-bdef-4463-88db-98f22de89214"),
					check.That(resourceType+".user_read_all").Key("target_service_principal_object_id").Exists(),
					check.That(resourceType+".user_read_all").Key("principal_type").HasValue("ServicePrincipal"),
					check.That(resourceType+".user_read_all").Key("principal_display_name").Exists(),
					check.That(resourceType+".user_read_all").Key("resource_display_name").Exists(),
					check.That(resourceType+".user_read_all").Key("created_date_time").Exists(),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing app role assignment to service principal")
				},
				ResourceName: resourceType + ".user_read_all",
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceType+".user_read_all"]
					if !ok {
						return "", fmt.Errorf("resource not found: %s", resourceType+".user_read_all")
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
