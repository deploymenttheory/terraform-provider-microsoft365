package graphBetaAuthorizationPolicy_test

import (
	"fmt"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/testlog"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaAuthorizationPolicy "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/authorization_policy"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var (
	// Resource type name from the resource package
	resourceType = graphBetaAuthorizationPolicy.ResourceName

	// testResource is the test resource implementation for authorization policies
	testResource = graphBetaAuthorizationPolicy.AuthorizationPolicyTestResource{}
)

// loadAcceptanceTestTerraform loads test configurations from the acceptance directory.
func loadAcceptanceTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/acceptance/" + filename)
	if err != nil {
		panic(fmt.Sprintf("failed to load acceptance config %s: %s", filename, err))
	}
	return acceptance.ConfiguredM365ProviderBlock(config)
}

// TestAccResourceAuthorizationPolicy_01_Lifecycle tests adoption, update, import, and stable planning.
// CheckDestroy is nil because this singleton remains in the tenant after Terraform destroy.
// The pre-check registers cleanup to restore and verify the original policy settings.
// These tests configure tenant-wide authorization settings and should use a dedicated test tenant.
func TestAccResourceAuthorizationPolicy_01_Lifecycle(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { authorizationPolicyPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { testlog.StepAction(resourceType, "Adopting the singleton authorization policy") },
				Config:    loadAcceptanceTestTerraform("resource_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").ExistsInGraph(testResource),
					check.That(resourceType+".test").Key("id").HasValue("authorizationPolicy"),
					check.That(resourceType+".test").Key("allow_invites_from").HasValue("adminsAndGuestInviters"),
					check.That(resourceType+".test").Key("allowed_to_use_sspr").HasValue("true"),
				),
			},
			{
				PreConfig: func() { testlog.StepAction(resourceType, "Updating authorization policy settings") },
				Config:    loadAcceptanceTestTerraform("resource_complete.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").ExistsInGraph(testResource),
					check.That(resourceType+".test").Key("allow_invites_from").HasValue("none"),
					check.That(resourceType+".test").Key("guest_user_role_id").HasValue("2af84b1e-32c8-42b7-82bc-daa82404023b"),
					check.That(resourceType+".test").Key("default_user_role_permissions.allowed_to_create_apps").HasValue("false"),
				),
			},
			{
				PreConfig:               func() { testlog.StepAction(resourceType, "Importing the singleton authorization policy") },
				ResourceName:            resourceType + ".test",
				ImportState:             true,
				ImportStateId:           "authorizationPolicy",
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
			{Config: loadAcceptanceTestTerraform("resource_complete.tf"), PlanOnly: true},
		},
	})
}

// TestAccResourceAuthorizationPolicy_02_EmptySets tests clearing collection settings explicitly.
func TestAccResourceAuthorizationPolicy_02_EmptySets(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { authorizationPolicyPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { testlog.StepAction(resourceType, "Clearing authorization policy collections") },
				Config:    loadAcceptanceTestTerraform("resource_empty_sets.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").ExistsInGraph(testResource),
					check.That(resourceType+".test").Key("enabled_preview_features.#").HasValue("0"),
					check.That(resourceType+".test").Key("permission_grant_policy_ids_assigned_to_default_user_role.#").HasValue("0"),
				),
			},
			{Config: loadAcceptanceTestTerraform("resource_empty_sets.tf"), PlanOnly: true},
		},
	})
}
