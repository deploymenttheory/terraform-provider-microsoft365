package graphBetaSettingsCatalogConfigurationPolicy_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	listMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/list-resources/device_management/graph_beta/settings_catalog_configuration_policy/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/querycheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/jarcoal/httpmock"
)

const (
	listType = "list.microsoft365_graph_beta_device_management_settings_catalog_configuration_policy"
)

// Helper function to load test configs from unit directory
func loadUnitTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/unit/" + filename)
	if err != nil {
		panic("failed to load unit test config " + filename + ": " + err.Error())
	}
	return config
}

func setupMockEnvironment() (*mocks.Mocks, *listMocks.SettingsCatalogListMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	listMock := &listMocks.SettingsCatalogListMock{}
	listMock.RegisterMocks()
	return mockClient, listMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *listMocks.SettingsCatalogListMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	listMock := &listMocks.SettingsCatalogListMock{}
	listMock.RegisterErrorMocks()
	return mockClient, listMock
}

func TestUnitListResourceSettingsCatalogConfigurationPolicy_01_All(t *testing.T) {
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
				Config: loadUnitTestTerraform("list_all.tfquery.hcl"),
				QueryResultChecks: []querycheck.QueryResultCheck{
					// Should return all 6 policies from JSON
					querycheck.ExpectLength(listType+".all", 6),
				},
			},
		},
	})
}

func TestUnitListResourceSettingsCatalogConfigurationPolicy_02_ByName(t *testing.T) {
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
				Config: loadUnitTestTerraform("list_by_name.tfquery.hcl"),
				QueryResultChecks: []querycheck.QueryResultCheck{
					// Should filter policies with "Kerberos" in name (3 policies)
					querycheck.ExpectLength(listType+".by_name", 3),
				},
			},
		},
	})
}

func TestUnitListResourceSettingsCatalogConfigurationPolicy_03_ByPlatform(t *testing.T) {
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
				Config: loadUnitTestTerraform("list_by_platform.tfquery.hcl"),
				QueryResultChecks: []querycheck.QueryResultCheck{
					// Should return Windows10 policies (4 policies)
					querycheck.ExpectLength(listType+".by_platform", 4),
				},
			},
		},
	})
}

func TestUnitListResourceSettingsCatalogConfigurationPolicy_04_ByTemplateFamily(t *testing.T) {
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
				Config: loadUnitTestTerraform("list_by_template_family.tfquery.hcl"),
				QueryResultChecks: []querycheck.QueryResultCheck{
					// Should return endpointSecurityAntivirus policies (2 policies)
					querycheck.ExpectLength(listType+".by_template_family", 2),
				},
			},
		},
	})
}

func TestUnitListResourceSettingsCatalogConfigurationPolicy_05_AssignedOnly(t *testing.T) {
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
				Config: loadUnitTestTerraform("list_assigned_only.tfquery.hcl"),
				QueryResultChecks: []querycheck.QueryResultCheck{
					// Should return assigned policies (4 policies: pol-001, pol-003, pol-005, pol-006 is false)
					querycheck.ExpectLength(listType+".assigned_only", 3),
				},
			},
		},
	})
}

func TestUnitListResourceSettingsCatalogConfigurationPolicy_06_CombinedFilters(t *testing.T) {
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
				Config: loadUnitTestTerraform("list_combined_filters.tfquery.hcl"),
				QueryResultChecks: []querycheck.QueryResultCheck{
					// Should return Windows10 + Defender policies (2 policies)
					querycheck.ExpectLength(listType+".combined", 2),
				},
			},
		},
	})
}

func TestUnitListResourceSettingsCatalogConfigurationPolicy_07_Error(t *testing.T) {
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
				Config:      loadUnitTestTerraform("list_all.tfquery.hcl"),
				ExpectError: regexp.MustCompile("403|Forbidden"),
			},
		},
	})
}
