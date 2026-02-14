package graphBetaUsersUser_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	listMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/list-resources/users/graph_beta/user/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/querycheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/jarcoal/httpmock"
)

// Helper function to load test configs from unit directory
func loadUnitTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/unit/" + filename)
	if err != nil {
		panic("failed to load unit test config " + filename + ": " + err.Error())
	}
	return config
}

func setupMockEnvironment() (*mocks.Mocks, *listMocks.UserListMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	listMock := &listMocks.UserListMock{}
	listMock.RegisterMocks()
	return mockClient, listMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *listMocks.UserListMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	listMock := &listMocks.UserListMock{}
	listMock.RegisterErrorMocks()
	return mockClient, listMock
}

func TestUnitListResourceUser_01_All(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, listMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer listMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		Steps: []resource.TestStep{
			{
				Query:  true,
				Config: loadUnitTestTerraform("list_resource_01_all.tfquery.hcl"),
				QueryResultChecks: []querycheck.QueryResultCheck{
					// Validate users are present
					querycheck.ExpectLengthAtLeast(listType+".all", 1),
				},
			},
		},
	})
}

func TestUnitListResourceUser_02_ByDisplayName(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, listMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer listMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		Steps: []resource.TestStep{
			{
				Query:  true,
				Config: loadUnitTestTerraform("list_resource_02_by_display_name.tfquery.hcl"),
				QueryResultChecks: []querycheck.QueryResultCheck{
					// Should return users with "John" in display name
					querycheck.ExpectLengthAtLeast(listType+".by_display_name", 1),
				},
			},
		},
	})
}

func TestUnitListResourceUser_03_ByUserPrincipalName(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, listMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer listMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		Steps: []resource.TestStep{
			{
				Query:  true,
				Config: loadUnitTestTerraform("list_resource_03_by_upn.tfquery.hcl"),
				QueryResultChecks: []querycheck.QueryResultCheck{
					// Should return users with UPN starting with "admin"
					querycheck.ExpectLengthAtLeast(listType+".by_upn", 1),
				},
			},
		},
	})
}

func TestUnitListResourceUser_04_ByAccountEnabled(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, listMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer listMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		Steps: []resource.TestStep{
			{
				Query:  true,
				Config: loadUnitTestTerraform("list_resource_04_by_account_enabled.tfquery.hcl"),
				QueryResultChecks: []querycheck.QueryResultCheck{
					// Should return enabled users
					querycheck.ExpectLengthAtLeast(listType+".by_account_enabled", 1),
				},
			},
		},
	})
}

func TestUnitListResourceUser_05_ByUserType(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, listMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer listMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		Steps: []resource.TestStep{
			{
				Query:  true,
				Config: loadUnitTestTerraform("list_resource_05_by_user_type.tfquery.hcl"),
				QueryResultChecks: []querycheck.QueryResultCheck{
					// Should return member users
					querycheck.ExpectLengthAtLeast(listType+".by_user_type", 1),
				},
			},
		},
	})
}

func TestUnitListResourceUser_06_CombinedFilters(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, listMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer listMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		Steps: []resource.TestStep{
			{
				Query:  true,
				Config: loadUnitTestTerraform("list_resource_06_combined_filters.tfquery.hcl"),
				QueryResultChecks: []querycheck.QueryResultCheck{
					// Should return enabled member users
					querycheck.ExpectLengthAtLeast(listType+".combined", 1),
				},
			},
		},
	})
}

func TestUnitListResourceUser_07_Error(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, listMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer listMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		Steps: []resource.TestStep{
			{
				Query:       true,
				Config:      loadUnitTestTerraform("list_resource_01_all.tfquery.hcl"),
				ExpectError: regexp.MustCompile("403|Forbidden"),
			},
		},
	})
}
