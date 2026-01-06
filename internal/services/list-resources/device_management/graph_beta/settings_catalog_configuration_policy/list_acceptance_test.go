package graphBetaSettingsCatalogConfigurationPolicy_test

import (
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/querycheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
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

// TestAccSettingsCatalogList_All tests fetching all policies from live API
func TestAccSettingsCatalogList_All(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		Steps: []resource.TestStep{
			{
				Query:  true,
				Config: loadAcceptanceTestTerraform("list_all.tfquery.hcl"),
				QueryResultChecks: []querycheck.QueryResultCheck{
					// Should return at least 1 policy from live API
					// We can't assert on specific IDs/values since they vary per tenant
					querycheck.ExpectLengthAtLeast(listType+".all", 1),
				},
			},
		},
	})
}

// TestAccSettingsCatalogList_ByPlatform tests filtering by platform from live API
func TestAccSettingsCatalogList_ByPlatform(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		Steps: []resource.TestStep{
			{
				Query:  true,
				Config: loadAcceptanceTestTerraform("list_by_platform.tfquery.hcl"),
				QueryResultChecks: []querycheck.QueryResultCheck{
					// Note: items.# may be 0 if no Windows policies exist
					querycheck.ExpectLengthAtLeast(listType+".by_platform", 0),
				},
			},
		},
	})
}

// TestAccSettingsCatalogList_AssignedOnly tests filtering by assignment status from live API
func TestAccSettingsCatalogList_AssignedOnly(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		Steps: []resource.TestStep{
			{
				Query:  true,
				Config: loadAcceptanceTestTerraform("list_assigned_only.tfquery.hcl"),
				QueryResultChecks: []querycheck.QueryResultCheck{
					// Note: items.# may be 0 if no assigned policies exist
					querycheck.ExpectLengthAtLeast(listType+".assigned_only", 0),
				},
			},
		},
	})
}
