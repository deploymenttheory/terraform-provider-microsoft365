package graphBetaUserLicenseAssignment_test

import (
	"testing"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/destroy"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/testlog"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaUserLicenseAssignment "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/users/graph_beta/license_assignment"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

var (
	// testResource is the test resource implementation for user license assignments
	testResource = graphBetaUserLicenseAssignment.UserLicenseAssignmentTestResource{}
)

// loadAcceptanceTestTerraform loads an acceptance test terraform configuration file
func loadAcceptanceTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/acceptance/" + filename)
	if err != nil {
		panic("failed to load acceptance config " + filename + ": " + err.Error())
	}
	return config
}

func TestAccResourceUserLicenseAssignment_01_Lifecycle(t *testing.T) {
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
					testlog.StepAction(resourceType, "Creating")
				},
				Config: loadAcceptanceTestTerraform("resource_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						testlog.WaitForConsistency("user license assignment", 90*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".minimal").ExistsInGraph(testResource),
					check.That(resourceType+".minimal").Key("id").Exists(),
					check.That(resourceType+".minimal").Key("sku_id").HasValue("a403ebcc-fae0-4ca2-8c8c-7a907fd6c235"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing")
				},
				ResourceName:      resourceType + ".minimal",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceUserLicenseAssignment_02_Maximal(t *testing.T) {
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
					testlog.StepAction(resourceType, "Creating")
				},
				Config: loadAcceptanceTestTerraform("resource_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("user license assignment", 90*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".maximal").ExistsInGraph(testResource),
					check.That(resourceType+".maximal").Key("id").Exists(),
					check.That(resourceType+".maximal").Key("sku_id").HasValue("a403ebcc-fae0-4ca2-8c8c-7a907fd6c235"),
					check.That(resourceType+".maximal").Key("disabled_plans.#").HasValue("1"),
					check.That(resourceType+".maximal").Key("disabled_plans.*").ContainsTypeSetElement("c948ea65-2053-4a5a-8a62-9eaaaf11b522"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing")
				},
				ResourceName:      resourceType + ".maximal",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// TestAccResourceUserLicenseAssignment_03_DisabledPlansLifecycle is the acceptance
// equivalent of the unit test for the same bug: proves that removing disabled_plans
// from config sends an explicit clear to the Graph API rather than silently retaining
// the previous values.
//
// Step 1 — Create with one disabled plan: verify disabled_plans.# = 1
// Step 2 — Remove disabled_plans from config: verify disabled_plans.# = 0 and that a
//
//	subsequent plan shows no diff (no spurious drift from stale API state)
func TestAccResourceUserLicenseAssignment_03_DisabledPlansLifecycle(t *testing.T) {
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
					testlog.StepAction(resourceType, "Creating with disabled plan")
				},
				Config: loadAcceptanceTestTerraform("resource_lifecycle_step1_with_disabled_plans.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("user license assignment with disabled plan", 90*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".lifecycle").ExistsInGraph(testResource),
					check.That(resourceType+".lifecycle").Key("sku_id").HasValue("a403ebcc-fae0-4ca2-8c8c-7a907fd6c235"),
					check.That(resourceType+".lifecycle").Key("disabled_plans.#").HasValue("1"),
					check.That(resourceType+".lifecycle").Key("disabled_plans.*").ContainsTypeSetElement("c948ea65-2053-4a5a-8a62-9eaaaf11b522"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Removing disabled_plans — verifying plans are cleared in API")
				},
				Config: loadAcceptanceTestTerraform("resource_lifecycle_step2_without_disabled_plans.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("user license assignment after clearing disabled plans", 90*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".lifecycle").ExistsInGraph(testResource),
					check.That(resourceType+".lifecycle").Key("disabled_plans.#").HasValue("0"),
				),
			},
			// Re-plan to prove no drift: if the API still had the old disabled plan,
			// this would produce a non-empty plan and fail with ExpectNonEmptyPlan implicitly.
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Verifying no drift after disabled plans cleared")
				},
				Config:   loadAcceptanceTestTerraform("resource_lifecycle_step2_without_disabled_plans.tf"),
				PlanOnly: true,
			},
		},
	})
}
