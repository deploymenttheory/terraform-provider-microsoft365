package graphBetaUsersUser_test

import (
	"testing"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/destroy"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/testlog"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaUsersUser "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/users/graph_beta/user"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var (
	// Resource type name from the resource package
	resourceType = graphBetaUsersUser.ResourceName

	// testResource is the test resource implementation for users
	testResource = graphBetaUsersUser.UserTestResource{}
)

func TestAccUserResource_Lifecycle(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			resourceType,
			5*time.Second,
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
					testlog.StepAction(resourceType, "Creating")
				},
				Config: testAccUserConfig_minimal(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".minimal").ExistsInGraph(testResource),
					check.That(resourceType+".minimal").Key("id").Exists(),
					check.That(resourceType+".minimal").Key("display_name").IsNotEmpty(),
					check.That(resourceType+".minimal").Key("user_principal_name").IsNotEmpty(),
					check.That(resourceType+".minimal").Key("account_enabled").HasValue("true"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing")
				},
				ResourceName:      resourceType + ".minimal",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"password_profile",
					"password_profile.%",
					"password_profile.password",
					"password_profile.force_change_password_next_sign_in",
					"password_profile.force_change_password_next_sign_in_with_mfa",
				},
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Updating")
				},
				Config: testAccUserConfig_updated(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".maximal").ExistsInGraph(testResource),
					check.That(resourceType+".maximal").Key("id").Exists(),
					check.That(resourceType+".maximal").Key("display_name").IsNotEmpty(),
					check.That(resourceType+".maximal").Key("given_name").IsNotEmpty(),
					check.That(resourceType+".maximal").Key("surname").IsNotEmpty(),
					check.That(resourceType+".maximal").Key("job_title").IsNotEmpty(),
				),
			},
			{
				ResourceName:      resourceType + ".maximal",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Test configuration functions
func testAccUserConfig_minimal() string {
	config := mocks.LoadTerraformConfigFile("resource_minimal.tf")
	return acceptance.ConfiguredM365ProviderBlock(config)
}

func testAccUserConfig_updated() string {
	config := mocks.LoadTerraformConfigFile("resource_maximal.tf")
	return acceptance.ConfiguredM365ProviderBlock(config)
}
