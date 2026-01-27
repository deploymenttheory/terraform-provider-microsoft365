package graphBetaAppleConfiguratorEnrollmentPolicy_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	appleMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/apple_configurator_enrollment_policy/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

func setupMockEnvironment() (*mocks.Mocks, *appleMocks.AppleConfiguratorEnrollmentPolicyMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	appleMock := &appleMocks.AppleConfiguratorEnrollmentPolicyMock{}
	appleMock.RegisterMocks()
	return mockClient, appleMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *appleMocks.AppleConfiguratorEnrollmentPolicyMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	appleMock := &appleMocks.AppleConfiguratorEnrollmentPolicyMock{}
	appleMock.RegisterErrorMocks()
	return mockClient, appleMock
}

func testCheckExists(resourceName string) resource.TestCheckFunc {
	return resource.TestCheckResourceAttrSet(resourceName, "id")
}

func TestUnitResourceAppleConfiguratorEnrollmentPolicy_01_Schema(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, appleMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer appleMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_apple_configurator_enrollment_policy.minimal", "display_name", "Test Minimal Apple Configurator Enrollment Policy - Unique"),
					resource.TestMatchResourceAttr("microsoft365_graph_beta_device_management_apple_configurator_enrollment_policy.minimal", "id", regexp.MustCompile(`^[0-9a-fA-F-]+_[0-9a-fA-F-]+$`)),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_apple_configurator_enrollment_policy.minimal", "requires_user_authentication", "false"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_apple_configurator_enrollment_policy.minimal", "enable_authentication_via_company_portal", "false"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_apple_configurator_enrollment_policy.minimal", "require_company_portal_on_setup_assistant_enrolled_devices", "false"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_apple_configurator_enrollment_policy.minimal", "configuration_endpoint_url"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_apple_configurator_enrollment_policy.minimal", "dep_onboarding_settings_id"),
				),
			},
		},
	})
}

func TestUnitResourceAppleConfiguratorEnrollmentPolicy_02_AuthenticationScenarios(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, appleMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer appleMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigWithCompanyPortal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_apple_configurator_enrollment_policy.with_company_portal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_apple_configurator_enrollment_policy.with_company_portal", "display_name", "Test Company Portal Apple Configurator Enrollment Policy - Unique"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_apple_configurator_enrollment_policy.with_company_portal", "enable_authentication_via_company_portal", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_apple_configurator_enrollment_policy.with_company_portal", "require_company_portal_on_setup_assistant_enrolled_devices", "false"),
				),
			},
			{
				Config: testConfigWithSetupAssistant(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_apple_configurator_enrollment_policy.with_setup_assistant"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_apple_configurator_enrollment_policy.with_setup_assistant", "display_name", "Test Setup Assistant Apple Configurator Enrollment Policy - Unique"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_apple_configurator_enrollment_policy.with_setup_assistant", "enable_authentication_via_company_portal", "false"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_apple_configurator_enrollment_policy.with_setup_assistant", "require_company_portal_on_setup_assistant_enrolled_devices", "true"),
				),
			},
		},
	})
}

func TestUnitResourceAppleConfiguratorEnrollmentPolicy_03_ValidationErrors(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, appleMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer appleMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testConfigValidationError(),
				ExpectError: regexp.MustCompile("Mutually Exclusive Fields"),
			},
		},
	})
}

func TestUnitResourceAppleConfiguratorEnrollmentPolicy_04_ErrorHandling(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, appleMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer appleMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testConfigMinimal(),
				ExpectError: regexp.MustCompile("Internal server error|InternalServerError|500"),
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

func testConfigWithCompanyPortal() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_with_company_portal.tf")
	if err != nil {
		panic("failed to load company portal config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigWithSetupAssistant() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_with_setup_assistant.tf")
	if err != nil {
		panic("failed to load setup assistant config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigValidationError() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_validation_error.tf")
	if err != nil {
		panic("failed to load validation error config: " + err.Error())
	}
	return unitTestConfig
}
