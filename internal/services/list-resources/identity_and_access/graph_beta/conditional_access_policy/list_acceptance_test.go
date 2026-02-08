package graphBetaConditionalAccessPolicy_test

import (
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/querycheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

const (
	listType = "list.microsoft365_graph_beta_identity_and_access_conditional_access_policy"
)

// Helper function to load test configs from acceptance directory
func loadAcceptanceTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/acceptance/" + filename)
	if err != nil {
		panic("failed to load acceptance test config " + filename + ": " + err.Error())
	}
	// For Query mode tests, the provider block is already in the .HCL file
	return config
}

// TestAccListResourceConditionalAccessPolicy_01_All tests fetching all policies from live API
func TestAccListResourceConditionalAccessPolicy_01_All(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		Steps: []resource.TestStep{
			{
				Query:  true,
				Config: loadAcceptanceTestTerraform("list_resource_01_all.tfquery.hcl"),
				QueryResultChecks: []querycheck.QueryResultCheck{
					// Verify list operation works (may return 0 or more policies)
					// We can't assert on specific IDs/values since they vary per tenant
					querycheck.ExpectLengthAtLeast(listType+".all", 0),
				},
			},
		},
	})
}

// TestAccListResourceConditionalAccessPolicy_02_FilterByState tests filtering by policy state from live API
func TestAccListResourceConditionalAccessPolicy_02_FilterByState(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		Steps: []resource.TestStep{
			{
				Query:  true,
				Config: loadAcceptanceTestTerraform("list_resource_02_by_state.tfquery.hcl"),
				QueryResultChecks: []querycheck.QueryResultCheck{
					// Note: items.# may be 0 if no policies exist in the filtered state
					querycheck.ExpectLengthAtLeast(listType+".by_state", 0),
				},
			},
		},
	})
}

// TestAccListResourceConditionalAccessPolicy_03_FilterByDisplayName tests filtering by display name from live API
func TestAccListResourceConditionalAccessPolicy_03_FilterByDisplayName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		Steps: []resource.TestStep{
			{
				Query:  true,
				Config: loadAcceptanceTestTerraform("list_resource_03_by_display_name.tfquery.hcl"),
				QueryResultChecks: []querycheck.QueryResultCheck{
					// Note: items.# may be 0 if no policies match the display name filter
					querycheck.ExpectLengthAtLeast(listType+".by_display_name", 0),
				},
			},
		},
	})
}
