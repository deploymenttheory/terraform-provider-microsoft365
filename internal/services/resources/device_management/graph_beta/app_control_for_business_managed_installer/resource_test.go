package graphBetaAppControlForBusinessManagedInstaller_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	managedInstallerMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/app_control_for_business_managed_installer/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

func setupMockEnvironment() (*mocks.Mocks, *managedInstallerMocks.AppControlForBusinessManagedInstallerMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	managedInstallerMock := &managedInstallerMocks.AppControlForBusinessManagedInstallerMock{}
	managedInstallerMock.RegisterMocks()
	return mockClient, managedInstallerMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *managedInstallerMocks.AppControlForBusinessManagedInstallerMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	managedInstallerMock := &managedInstallerMocks.AppControlForBusinessManagedInstallerMock{}
	managedInstallerMock.RegisterErrorMocks()
	return mockClient, managedInstallerMock
}

func testCheckExists(resourceName string) resource.TestCheckFunc {
	return resource.TestCheckResourceAttrSet(resourceName, "id")
}

func testConfigMinimal() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_minimal.tf")
	if err != nil {
		panic("failed to load minimal config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigMaximal() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_maximal.tf")
	if err != nil {
		panic("failed to load maximal config: " + err.Error())
	}
	return unitTestConfig
}

// TestAppControlForBusinessManagedInstallerResource_Schema validates the resource schema
func TestAppControlForBusinessManagedInstallerResource_Schema(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, managedInstallerMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer managedInstallerMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					// Check required attributes
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_app_control_for_business_managed_installer.minimal", "intune_management_extension_as_managed_installer", "Disabled"),

					// Check computed attributes are set
					resource.TestMatchResourceAttr("microsoft365_graph_beta_device_management_app_control_for_business_managed_installer.minimal", "id", regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_app_control_for_business_managed_installer.minimal", "available_version", "1.93.102.0"),
				),
			},
		},
	})
}

// TestAppControlForBusinessManagedInstallerResource_MaximalSettings tests with enabled configuration
func TestAppControlForBusinessManagedInstallerResource_MaximalSettings(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, managedInstallerMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer managedInstallerMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMaximal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_app_control_for_business_managed_installer.maximal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_app_control_for_business_managed_installer.maximal", "intune_management_extension_as_managed_installer", "Enabled"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_app_control_for_business_managed_installer.maximal", "available_version", "1.93.102.0"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_app_control_for_business_managed_installer.maximal", "managed_installer_configured_date_time"),
				),
			},
		},
	})
}

// TestAppControlForBusinessManagedInstallerResource_ErrorHandling tests error scenarios
func TestAppControlForBusinessManagedInstallerResource_ErrorHandling(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, managedInstallerMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer managedInstallerMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testConfigMinimal(),
				ExpectError: regexp.MustCompile("Windows Management App not found|ResourceNotFound"),
			},
		},
	})
}