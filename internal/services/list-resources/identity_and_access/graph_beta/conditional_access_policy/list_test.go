package graphBetaConditionalAccessPolicy_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	listMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/list-resources/identity_and_access/graph_beta/conditional_access_policy/mocks"
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

func setupMockEnvironment() (*mocks.Mocks, *listMocks.ConditionalAccessPolicyListMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	listMock := &listMocks.ConditionalAccessPolicyListMock{}
	listMock.RegisterMocks()
	return mockClient, listMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *listMocks.ConditionalAccessPolicyListMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	listMock := &listMocks.ConditionalAccessPolicyListMock{}
	listMock.RegisterErrorMocks()
	return mockClient, listMock
}

func TestUnitListResourceConditionalAccessPolicy_01_All(t *testing.T) {
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
					// Validate specific policies are present
					querycheck.ExpectLengthAtLeast(listType+".all", 1),
				},
			},
		},
	})
}

func TestUnitListResourceConditionalAccessPolicy_02_ByDisplayName(t *testing.T) {
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
					// Should return policies with "MFA" in display name
					querycheck.ExpectLengthAtLeast(listType+".by_display_name", 1),
				},
			},
		},
	})
}

func TestUnitListResourceConditionalAccessPolicy_03_ByState(t *testing.T) {
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
				Config: loadUnitTestTerraform("list_resource_03_by_state.tfquery.hcl"),
				QueryResultChecks: []querycheck.QueryResultCheck{
					// Should return enabled policies
					querycheck.ExpectLengthAtLeast(listType+".by_state", 1),
				},
			},
		},
	})
}

func TestUnitListResourceConditionalAccessPolicy_04_CombinedFilters(t *testing.T) {
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
				Config: loadUnitTestTerraform("list_resource_04_combined_filters.tfquery.hcl"),
				QueryResultChecks: []querycheck.QueryResultCheck{
					// Should return enabled policies with "MFA" in display name
					querycheck.ExpectLengthAtLeast(listType+".combined", 1),
				},
			},
		},
	})
}

func TestUnitListResourceConditionalAccessPolicy_05_Error(t *testing.T) {
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
