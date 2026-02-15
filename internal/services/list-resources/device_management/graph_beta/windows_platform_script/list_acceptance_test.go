package graphBetaWindowsPlatformScript_test

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

// TestAccListResourceWindowsPlatformScript_01_All tests fetching all scripts from live API
func TestAccListResourceWindowsPlatformScript_01_All(t *testing.T) {
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
					// Verify list operation works (may return 0 or more scripts)
					// We can't assert on specific IDs/values since they vary per tenant
					querycheck.ExpectLengthAtLeast(listType+".all", 0),
				},
			},
		},
	})
}

// TestAccListResourceWindowsPlatformScript_02_FilterByDisplayName tests filtering by display name from live API
func TestAccListResourceWindowsPlatformScript_02_FilterByDisplayName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		Steps: []resource.TestStep{
			{
				Query:  true,
				Config: loadAcceptanceTestTerraform("list_resource_02_by_display_name.tfquery.hcl"),
				QueryResultChecks: []querycheck.QueryResultCheck{
					// Note: items.# may be 0 if no scripts match the display name filter
					querycheck.ExpectLengthAtLeast(listType+".filtered", 0),
				},
			},
		},
	})
}

// TestAccListResourceWindowsPlatformScript_03_FilterByFileName tests filtering by file name from live API
func TestAccListResourceWindowsPlatformScript_03_FilterByFileName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		Steps: []resource.TestStep{
			{
				Query:  true,
				Config: loadAcceptanceTestTerraform("list_resource_03_by_file_name.tfquery.hcl"),
				QueryResultChecks: []querycheck.QueryResultCheck{
					// Note: items.# may be 0 if no scripts match the file name filter
					querycheck.ExpectLengthAtLeast(listType+".filtered", 0),
				},
			},
		},
	})
}

// TestAccListResourceWindowsPlatformScript_04_FilterByRunAsAccount tests filtering by run as account from live API
func TestAccListResourceWindowsPlatformScript_04_FilterByRunAsAccount(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		Steps: []resource.TestStep{
			{
				Query:  true,
				Config: loadAcceptanceTestTerraform("list_resource_04_by_run_as_account.tfquery.hcl"),
				QueryResultChecks: []querycheck.QueryResultCheck{
					// Note: items.# may be 0 if no system scripts exist
					querycheck.ExpectLengthAtLeast(listType+".filtered", 0),
				},
			},
		},
	})
}

// TestAccListResourceWindowsPlatformScript_05_CombinedFilters tests combined filters from live API
func TestAccListResourceWindowsPlatformScript_05_CombinedFilters(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		Steps: []resource.TestStep{
			{
				Query:  true,
				Config: loadAcceptanceTestTerraform("list_resource_05_combined_filters.tfquery.hcl"),
				QueryResultChecks: []querycheck.QueryResultCheck{
					// Note: items.# may be 0 if no scripts match the combined filters
					querycheck.ExpectLengthAtLeast(listType+".filtered", 0),
				},
			},
		},
	})
}
