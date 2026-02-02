package graphBetaTenantAppManagementPolicy_test

import (
	"testing"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/destroy"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/testlog"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// Helper function to load test configs from acceptance directory
func loadAcceptanceTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/acceptance/" + filename)
	if err != nil {
		panic("failed to load acceptance config " + filename + ": " + err.Error())
	}
	return acceptance.ConfiguredM365ProviderBlock(config)
}

// TestAccResourceTenantAppManagementPolicy_01_Minimal tests creating tenant app management policy with minimal configuration
func TestAccResourceTenantAppManagementPolicy_01_Minimal(t *testing.T) {
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
			"time": {
				Source:            "hashicorp/time",
				VersionConstraint: constants.ExternalProviderTimeVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating tenant app management policy with minimal configuration")
				},
				Config: loadAcceptanceTestTerraform("resource_01_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("tenant app management policy", 10*time.Second)
						time.Sleep(10 * time.Second)
						return nil
					},
					check.That(resourceType+".minimal").ExistsInGraph(testResource),
					check.That(resourceType+".minimal").Key("id").HasValue("00000000-0000-0000-0000-000000000000"),
					check.That(resourceType+".minimal").Key("is_enabled").HasValue("true"),
					check.That(resourceType+".minimal").Key("application_restrictions.password_credentials.0.restriction_type").HasValue("passwordLifetime"),
					check.That(resourceType+".minimal").Key("application_restrictions.password_credentials.0.state").HasValue("enabled"),
					check.That(resourceType+".minimal").Key("application_restrictions.password_credentials.0.max_lifetime").HasValue("P90D"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing tenant app management policy")
				},
				ResourceName:      resourceType + ".minimal",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"timeouts",
				},
			},
		},
	})
}

// TestAccResourceTenantAppManagementPolicy_02_Maximal tests creating tenant app management policy with maximal configuration
func TestAccResourceTenantAppManagementPolicy_02_Maximal(t *testing.T) {
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
			"time": {
				Source:            "hashicorp/time",
				VersionConstraint: constants.ExternalProviderTimeVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating tenant app management policy with maximal configuration")
				},
				Config: loadAcceptanceTestTerraform("resource_02_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("tenant app management policy", 10*time.Second)
						time.Sleep(10 * time.Second)
						return nil
					},
					check.That(resourceType+".maximal").ExistsInGraph(testResource),
					check.That(resourceType+".maximal").Key("id").HasValue("00000000-0000-0000-0000-000000000000"),
					check.That(resourceType+".maximal").Key("is_enabled").HasValue("true"),
					check.That(resourceType+".maximal").Key("display_name").HasValue("Custom Tenant App Management Policy"),
					check.That(resourceType+".maximal").Key("description").HasValue("Enforces comprehensive app management restrictions"),
					check.That(resourceType+".maximal").Key("restore_to_default_upon_delete").HasValue("false"),
					check.That(resourceType+".maximal").Key("application_restrictions.password_credentials.0.restriction_type").HasValue("passwordLifetime"),
					check.That(resourceType+".maximal").Key("application_restrictions.password_credentials.1.restriction_type").HasValue("symmetricKeyLifetime"),
					check.That(resourceType+".maximal").Key("application_restrictions.key_credentials.0.restriction_type").HasValue("asymmetricKeyLifetime"),
					check.That(resourceType+".maximal").Key("service_principal_restrictions.password_credentials.0.restriction_type").HasValue("passwordLifetime"),
					check.That(resourceType+".maximal").Key("service_principal_restrictions.key_credentials.0.restriction_type").HasValue("asymmetricKeyLifetime"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing tenant app management policy")
				},
				ResourceName:      resourceType + ".maximal",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"timeouts",
				},
			},
		},
	})
}

// TestAccResourceTenantAppManagementPolicy_03_Update tests updating tenant app management policy from minimal to maximal configuration
func TestAccResourceTenantAppManagementPolicy_03_Update(t *testing.T) {
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
			"time": {
				Source:            "hashicorp/time",
				VersionConstraint: constants.ExternalProviderTimeVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating tenant app management policy with minimal configuration")
				},
				Config: loadAcceptanceTestTerraform("resource_01_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("tenant app management policy", 10*time.Second)
						time.Sleep(10 * time.Second)
						return nil
					},
					check.That(resourceType+".minimal").ExistsInGraph(testResource),
					check.That(resourceType+".minimal").Key("is_enabled").HasValue("true"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Updating tenant app management policy to maximal configuration")
				},
				Config: loadAcceptanceTestTerraform("resource_02_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("tenant app management policy", 10*time.Second)
						time.Sleep(10 * time.Second)
						return nil
					},
					check.That(resourceType+".maximal").ExistsInGraph(testResource),
					check.That(resourceType+".maximal").Key("display_name").HasValue("Custom Tenant App Management Policy"),
					check.That(resourceType+".maximal").Key("description").HasValue("Enforces comprehensive app management restrictions"),
				),
			},
		},
	})
}
