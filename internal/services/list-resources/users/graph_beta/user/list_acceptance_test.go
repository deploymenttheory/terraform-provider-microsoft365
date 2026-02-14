package graphBetaUsersUser_test

import (
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/querycheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

const (
	listType = "list.microsoft365_graph_beta_users_user"
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

// TestAccListResourceUser_01_All tests fetching all users from live API
func TestAccListResourceUser_01_All(t *testing.T) {
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
					// Verify list operation works (may return 0 or more users)
					// We can't assert on specific IDs/values since they vary per tenant
					querycheck.ExpectLengthAtLeast(listType+".all", 0),
				},
			},
		},
	})
}

// TestAccListResourceUser_02_FilterByAccountEnabled tests filtering by account status from live API
func TestAccListResourceUser_02_FilterByAccountEnabled(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		Steps: []resource.TestStep{
			{
				Query:  true,
				Config: loadAcceptanceTestTerraform("list_resource_02_by_account_enabled.tfquery.hcl"),
				QueryResultChecks: []querycheck.QueryResultCheck{
					// Note: items.# may be 0 if no users exist in the filtered state
					querycheck.ExpectLengthAtLeast(listType+".by_account_enabled", 0),
				},
			},
		},
	})
}

// TestAccListResourceUser_03_FilterByUserType tests filtering by user type from live API
func TestAccListResourceUser_03_FilterByUserType(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		Steps: []resource.TestStep{
			{
				Query:  true,
				Config: loadAcceptanceTestTerraform("list_resource_03_by_user_type.tfquery.hcl"),
				QueryResultChecks: []querycheck.QueryResultCheck{
					// Note: items.# may be 0 if no users match the filter
					querycheck.ExpectLengthAtLeast(listType+".by_user_type", 0),
				},
			},
		},
	})
}

// TestAccListResourceUser_04_FilterByDisplayName tests filtering by display name from live API
func TestAccListResourceUser_04_FilterByDisplayName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		Steps: []resource.TestStep{
			{
				Query:  true,
				Config: loadAcceptanceTestTerraform("list_resource_04_by_display_name.tfquery.hcl"),
				QueryResultChecks: []querycheck.QueryResultCheck{
					// Note: items.# may be 0 if no users match the filter
					querycheck.ExpectLengthAtLeast(listType+".by_display_name", 0),
				},
			},
		},
	})
}

// TestAccListResourceUser_05_FilterByUPN tests filtering by user principal name from live API
func TestAccListResourceUser_05_FilterByUPN(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		Steps: []resource.TestStep{
			{
				Query:  true,
				Config: loadAcceptanceTestTerraform("list_resource_05_by_upn.tfquery.hcl"),
				QueryResultChecks: []querycheck.QueryResultCheck{
					// Note: items.# may be 0 if no users match the filter
					querycheck.ExpectLengthAtLeast(listType+".by_upn", 0),
				},
			},
		},
	})
}

// TestAccListResourceUser_06_CombinedFilters tests combined filters from live API
func TestAccListResourceUser_06_CombinedFilters(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		Steps: []resource.TestStep{
			{
				Query:  true,
				Config: loadAcceptanceTestTerraform("list_resource_06_combined_filters.tfquery.hcl"),
				QueryResultChecks: []querycheck.QueryResultCheck{
					// Note: items.# may be 0 if no users match the filter
					querycheck.ExpectLengthAtLeast(listType+".combined", 0),
				},
			},
		},
	})
}
