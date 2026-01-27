package graphBetaWindowsDeviceComplianceScript_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	scriptMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/windows_device_compliance_script/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

func setupMockEnvironment() (*mocks.Mocks, *scriptMocks.WindowsDeviceComplianceScriptMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	scriptMock := &scriptMocks.WindowsDeviceComplianceScriptMock{}
	scriptMock.RegisterMocks()
	return mockClient, scriptMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *scriptMocks.WindowsDeviceComplianceScriptMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	scriptMock := &scriptMocks.WindowsDeviceComplianceScriptMock{}
	scriptMock.RegisterErrorMocks()
	return mockClient, scriptMock
}

func testCheckExists(resourceName string) resource.TestCheckFunc {
	return resource.TestCheckResourceAttrSet(resourceName, "id")
}

func TestUnitResourceWindowsDeviceComplianceScript_01_Schema(t *testing.T) {
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
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_device_compliance_script.minimal", "display_name", "Test Minimal Windows Device Compliance Script - Unique"),
					resource.TestMatchResourceAttr("microsoft365_graph_beta_device_management_windows_device_compliance_script.minimal", "id", regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_device_compliance_script.minimal", "role_scope_tag_ids.#", "1"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_windows_device_compliance_script.minimal", "role_scope_tag_ids.*", "0"),
				),
			},
		},
	})
}

func TestUnitResourceWindowsDeviceComplianceScript_02_RunAsAccount(t *testing.T) {
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
					testCheckExists("microsoft365_graph_beta_device_management_windows_device_compliance_script.system_account"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_device_compliance_script.system_account", "display_name", "Test System Account Windows Device Compliance Script - Unique"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_device_compliance_script.system_account", "run_as_account", "system"),
				),
			},
			{
				Config: testConfigUserAccount(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_windows_device_compliance_script.user_account"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_device_compliance_script.user_account", "display_name", "Test User Account Windows Device Compliance Script - Unique"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_device_compliance_script.user_account", "run_as_account", "user"),
				),
			},
		},
	})
}

func TestUnitResourceWindowsDeviceComplianceScript_03_ErrorHandling(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, scriptMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer scriptMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testConfigMinimal(),
				ExpectError: regexp.MustCompile("Invalid Windows Device Compliance Script data"),
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
