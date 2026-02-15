package graphBetaWindowsPlatformScript_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	listMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/list-resources/device_management/graph_beta/windows_platform_script/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/querycheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/jarcoal/httpmock"
)

const (
	listType = "list.microsoft365_graph_beta_device_management_windows_platform_script"
)

// Helper function to load test configs from unit directory
func loadUnitTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/unit/" + filename)
	if err != nil {
		panic("failed to load unit test config " + filename + ": " + err.Error())
	}
	return config
}

func setupMockEnvironment() (*mocks.Mocks, *listMocks.WindowsPlatformScriptListMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	listMock := &listMocks.WindowsPlatformScriptListMock{}
	listMock.RegisterMocks()
	return mockClient, listMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *listMocks.WindowsPlatformScriptListMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	listMock := &listMocks.WindowsPlatformScriptListMock{}
	listMock.RegisterErrorMocks()
	return mockClient, listMock
}

func TestUnitListResourceWindowsPlatformScript_01_All(t *testing.T) {
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
					// Should return all 6 scripts from JSON
					querycheck.ExpectLength(listType+".all", 6),
				},
			},
		},
	})
}

func TestUnitListResourceWindowsPlatformScript_02_ByDisplayName(t *testing.T) {
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
					// Should filter scripts with "User" in display name (2 scripts)
					querycheck.ExpectLength(listType+".filtered", 2),
				},
			},
		},
	})
}

func TestUnitListResourceWindowsPlatformScript_03_ByFileName(t *testing.T) {
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
				Config: loadUnitTestTerraform("list_resource_03_by_file_name.tfquery.hcl"),
				QueryResultChecks: []querycheck.QueryResultCheck{
					// Should filter scripts with "baseline_setup.ps1" in file name (1 script)
					querycheck.ExpectLength(listType+".filtered", 1),
				},
			},
		},
	})
}

func TestUnitListResourceWindowsPlatformScript_04_ByRunAsAccount(t *testing.T) {
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
				Config: loadUnitTestTerraform("list_resource_04_by_run_as_account.tfquery.hcl"),
				QueryResultChecks: []querycheck.QueryResultCheck{
					// Should return system scripts (4 scripts)
					querycheck.ExpectLength(listType+".filtered", 4),
				},
			},
		},
	})
}

func TestUnitListResourceWindowsPlatformScript_05_CombinedFilters(t *testing.T) {
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
				Config: loadUnitTestTerraform("list_resource_05_combined_filters.tfquery.hcl"),
				QueryResultChecks: []querycheck.QueryResultCheck{
					// Should return scripts matching combined filters (1 script)
					querycheck.ExpectLength(listType+".filtered", 1),
				},
			},
		},
	})
}

func TestUnitListResourceWindowsPlatformScript_06_Error(t *testing.T) {
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
