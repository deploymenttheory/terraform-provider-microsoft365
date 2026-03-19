package graphBetaWindowsUpdatesAutopatchUpdatePolicy_test

import (
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/testlog"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func loadAcceptanceTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/acceptance/" + filename)
	if err != nil {
		panic("failed to load acceptance test config " + filename + ": " + err.Error())
	}
	return config
}

// TestAccResourceWindowsUpdatesUpdatePolicy_01_CreateUpdatePolicy creates an update policy
// and verifies import.
func TestAccResourceWindowsUpdatesUpdatePolicy_01_CreateUpdatePolicy(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Create update policy")
				},
				Config: loadAcceptanceTestTerraform("01_create_update_policy.tf"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceType+".test", "id"),
					resource.TestCheckResourceAttrSet(resourceType+".test", "created_date_time"),
					resource.TestCheckResourceAttrSet(resourceType+".test", "audience_id"),
					resource.TestCheckResourceAttr(resourceType+".test", "compliance_change_rules.#", "1"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Import update policy")
				},
				ResourceName:      resourceType + ".test",
				ImportState:       true,
				ImportStateVerify: true,
				// Compliance changes is write-only and not returned by the API
				ImportStateVerifyIgnore: []string{"timeouts", "compliance_changes"},
			},
		},
	})
}

// TestAccResourceWindowsUpdatesUpdatePolicy_02_LifecycleUpdate tests updating an update policy
func TestAccResourceWindowsUpdatesUpdatePolicy_02_LifecycleUpdate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Create initial update policy")
				},
				Config: loadAcceptanceTestTerraform("01_create_update_policy.tf"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceType+".test", "id"),
					resource.TestCheckResourceAttrSet(resourceType+".test", "audience_id"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Update policy settings")
				},
				Config: loadAcceptanceTestTerraform("02_lifecycle_update.tf"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceType+".test", "id"),
					resource.TestCheckResourceAttr(resourceType+".test", "compliance_change_rules.#", "1"),
				),
			},
		},
	})
}

// TestAccResourceWindowsUpdatesUpdatePolicy_03_MinimalPolicy tests creating a minimal update policy
func TestAccResourceWindowsUpdatesUpdatePolicy_03_MinimalPolicy(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Create minimal update policy")
				},
				Config: loadAcceptanceTestTerraform("03_minimal_policy.tf"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceType+".test", "id"),
					resource.TestCheckResourceAttrSet(resourceType+".test", "created_date_time"),
					resource.TestCheckResourceAttrSet(resourceType+".test", "audience_id"),
				),
			},
		},
	})
}
