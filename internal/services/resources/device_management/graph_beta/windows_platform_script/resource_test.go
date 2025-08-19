package graphBetaWindowsPlatformScript_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	scriptMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/windows_platform_script/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

func setupMockEnvironment() (*mocks.Mocks, *scriptMocks.WindowsPlatformScriptMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	scriptMock := &scriptMocks.WindowsPlatformScriptMock{}
	scriptMock.RegisterMocks()
	return mockClient, scriptMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *scriptMocks.WindowsPlatformScriptMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	scriptMock := &scriptMocks.WindowsPlatformScriptMock{}
	scriptMock.RegisterErrorMocks()
	return mockClient, scriptMock
}

func testCheckExists(resourceName string) resource.TestCheckFunc {
	return resource.TestCheckResourceAttrSet(resourceName, "id")
}

func TestWindowsPlatformScriptResource_Schema(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, scriptMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer scriptMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_platform_script.minimal", "display_name", "Test Minimal Windows Platform Script - Unique"),
					resource.TestMatchResourceAttr("microsoft365_graph_beta_device_management_windows_platform_script.minimal", "id", regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_platform_script.minimal", "role_scope_tag_ids.#", "1"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_windows_platform_script.minimal", "role_scope_tag_ids.*", "0"),
				),
			},
		},
	})
}

func TestWindowsPlatformScriptResource_RunAsAccount(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, scriptMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer scriptMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigSystemAccount(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_windows_platform_script.system_account"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_platform_script.system_account", "display_name", "Test System Account Windows Platform Script - Unique"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_platform_script.system_account", "run_as_account", "system"),
				),
			},
			{
				Config: testConfigUserAccount(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_windows_platform_script.user_account"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_platform_script.user_account", "display_name", "Test User Account Windows Platform Script - Unique"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_platform_script.user_account", "run_as_account", "user"),
				),
			},
		},
	})
}

func TestWindowsPlatformScriptResource_Assignments(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, scriptMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer scriptMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigWithAssignments(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_windows_platform_script.with_assignments"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_platform_script.with_assignments", "display_name", "Test Windows Platform Script with Assignments - Unique"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_platform_script.with_assignments", "assignments.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs("microsoft365_graph_beta_device_management_windows_platform_script.with_assignments", "assignments.*", map[string]string{
						"type": "groupAssignmentTarget",
					}),
					resource.TestCheckTypeSetElemNestedAttrs("microsoft365_graph_beta_device_management_windows_platform_script.with_assignments", "assignments.*", map[string]string{
						"type": "exclusionGroupAssignmentTarget",
					}),
				),
			},
		},
	})
}

func TestWindowsPlatformScriptResource_ErrorHandling(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, scriptMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer scriptMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testConfigMinimal(),
				ExpectError: regexp.MustCompile("Invalid Windows Platform Script data"),
			},
		},
	})
}

func testConfigMinimal() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_minimal.tf")
	if err != nil {
		panic("failed to load minimal config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigSystemAccount() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_system_account.tf")
	if err != nil {
		panic("failed to load system account config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigUserAccount() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_user_account.tf")
	if err != nil {
		panic("failed to load user account config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigWithAssignments() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_with_assignments.tf")
	if err != nil {
		panic("failed to load with assignments config: " + err.Error())
	}
	return unitTestConfig
}
