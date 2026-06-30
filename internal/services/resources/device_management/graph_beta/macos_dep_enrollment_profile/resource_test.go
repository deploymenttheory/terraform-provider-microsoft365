package graphBetaMacOSDepEnrollmentProfile_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	macosMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/macos_dep_enrollment_profile/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

func setupMockEnvironment() (*mocks.Mocks, *macosMocks.MacOSDepEnrollmentProfileMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	macosMock := &macosMocks.MacOSDepEnrollmentProfileMock{}
	macosMock.RegisterMocks()
	return mockClient, macosMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *macosMocks.MacOSDepEnrollmentProfileMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	macosMock := &macosMocks.MacOSDepEnrollmentProfileMock{}
	macosMock.RegisterErrorMocks()
	return mockClient, macosMock
}

func testCheckExists(resourceName string) resource.TestCheckFunc {
	return resource.TestCheckResourceAttrSet(resourceName, "id")
}

func TestUnitResourceMacOSDepEnrollmentProfile_01_Schema(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, macosMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer macosMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_dep_enrollment_profile.minimal", "display_name", "Test Minimal macOS DEP Enrollment Profile - Unique"),
					resource.TestMatchResourceAttr("microsoft365_graph_beta_device_management_macos_dep_enrollment_profile.minimal", "id", regexp.MustCompile(`^[0-9a-fA-F-]+_[0-9a-fA-F-]+$`)),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_dep_enrollment_profile.minimal", "requires_user_authentication", "false"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_macos_dep_enrollment_profile.minimal", "configuration_endpoint_url"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_macos_dep_enrollment_profile.minimal", "dep_onboarding_settings_id"),
				),
			},
		},
	})
}

func TestUnitResourceMacOSDepEnrollmentProfile_02_SkipSetupAndAdminAccount(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, macosMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer macosMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigSkipSetup(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_macos_dep_enrollment_profile.skip_setup"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_dep_enrollment_profile.skip_setup", "display_name", "Test Skip Setup macOS DEP Enrollment Profile - Unique"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_dep_enrollment_profile.skip_setup", "await_device_configured", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_dep_enrollment_profile.skip_setup", "supervised_mode_enabled", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_dep_enrollment_profile.skip_setup", "siri_disabled", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_dep_enrollment_profile.skip_setup", "file_vault_disabled", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_dep_enrollment_profile.skip_setup", "admin_account_user_name", "localadmin"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_dep_enrollment_profile.skip_setup", "hide_admin_account", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_dep_enrollment_profile.skip_setup", "enabled_skip_keys.#", "6"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_dep_enrollment_profile.skip_setup", "admin_account_password_rotation.auto_rotation_period_in_days", "30"),
				),
			},
		},
	})
}

func TestUnitResourceMacOSDepEnrollmentProfile_03_ValidationErrors(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, macosMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer macosMock.CleanupMockState()

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

func TestUnitResourceMacOSDepEnrollmentProfile_04_ErrorHandling(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, macosMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer macosMock.CleanupMockState()

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

func testConfigSkipSetup() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_skip_setup.tf")
	if err != nil {
		panic("failed to load skip setup config: " + err.Error())
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
