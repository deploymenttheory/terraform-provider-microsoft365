package graphBetaIdentityAndAccessMobilityManagementPolicy_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/testlog"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	graphBetaIdentityAndAccessMobilityManagementPolicy "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/mobility_management_policy"
)

var (
	// Resource type name from the resource package
	resourceType = graphBetaIdentityAndAccessMobilityManagementPolicy.ResourceName

	// testResource is the test resource implementation for mobility management policies
	testResource = graphBetaIdentityAndAccessMobilityManagementPolicy.MobilityManagementPolicyTestResource{}

	// externalProviders are the providers the acceptance configs use for unique names
	externalProviders = map[string]resource.ExternalProvider{
		"random": {
			Source:            "hashicorp/random",
			VersionConstraint: constants.ExternalProviderRandomVersion,
		},
	}
)

// Helper function to load test configs from acceptance directory
func loadAcceptanceTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/acceptance/" + filename)
	if err != nil {
		panic("failed to load acceptance config " + filename + ": " + err.Error())
	}
	return config
}

// waitForConsistency gives Microsoft Entra replicas time to converge before checks use a separate client
func waitForConsistency(_ *terraform.State) error {
	testlog.WaitForConsistency("mobility management policy", 15*time.Second)
	time.Sleep(15 * time.Second)
	return nil
}

// checkPolicyLeftInTenant confirms destroy only removed the resource from state; the service-owned policies remain
func checkPolicyLeftInTenant(_ *terraform.State) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	for _, policyType := range []string{"device", "app"} {
		if _, err := readPolicySnapshot(ctx, policyType); err != nil {
			return fmt.Errorf("the %s mobility management policy should remain after destroy: %w", policyType, err)
		}
	}
	return nil
}

// Minimal: adopt the Microsoft Intune MDM policy without writing anything, then import it.
// The MDM scope controls automatic Windows enrollment, so this test never changes it and fails if
// the policy differs from its snapshot afterwards.
func TestAccResourceMobilityManagementPolicy_01_MinimalDevice(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { snapshotPreCheck(t, "device") },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             checkPolicyLeftInTenant,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Adopting the Microsoft Intune MDM policy without changes")
				},
				Config: loadAcceptanceTestTerraform("resource_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".minimal").ExistsInGraph(testResource),
					check.That(resourceType+".minimal").Key("id").HasValue("device/0000000a-0000-0000-c000-000000000000"),
					check.That(resourceType+".minimal").Key("policy_type").HasValue("device"),
					check.That(resourceType+".minimal").Key("display_name").Exists(),
					check.That(resourceType+".minimal").Key("applies_to").Exists(),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing the Microsoft Intune MDM policy")
				},
				ResourceName:            resourceType + ".minimal",
				ImportState:             true,
				ImportStateId:           "device/0000000a-0000-0000-c000-000000000000",
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// App scope: scope the Microsoft Intune MAM policy to an acc-test group, import it, then restore the none scope.
// Runs only when the MAM scope is none, so no real user changes scope; the snapshot is checked afterwards.
func TestAccResourceMobilityManagementPolicy_02_AppSelectedGroup(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { appScopePreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             checkPolicyLeftInTenant,
		ExternalProviders:        externalProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Scoping the Microsoft Intune MAM policy to the test group")
				},
				Config: loadAcceptanceTestTerraform("resource_app_selected.tf"),
				Check: resource.ComposeTestCheckFunc(
					waitForConsistency,
					check.That(resourceType+".app").ExistsInGraph(testResource),
					check.That(resourceType+".app").Key("id").HasValue("app/0000000a-0000-0000-c000-000000000000"),
					check.That(resourceType+".app").Key("applies_to").HasValue("selected"),
					check.That(resourceType+".app").Key("included_groups.#").HasValue("1"),
					resource.TestCheckTypeSetElemAttrPair(resourceType+".app", "included_groups.*", "microsoft365_graph_beta_groups_group.scope", "id"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing the Microsoft Intune MAM policy")
				},
				ResourceName:            resourceType + ".app",
				ImportState:             true,
				ImportStateId:           "app/0000000a-0000-0000-c000-000000000000",
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Restoring the Microsoft Intune MAM policy scope to none")
				},
				Config: loadAcceptanceTestTerraform("resource_app_restore_none.tf"),
				Check: resource.ComposeTestCheckFunc(
					waitForConsistency,
					check.That(resourceType+".app").Key("applies_to").HasValue("none"),
					check.That(resourceType+".app").Key("included_groups.#").DoesNotExist(),
				),
			},
		},
	})
}
