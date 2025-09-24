package graphBetaWindowsAutopilotDeploymentProfile_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	profileMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/windows_autopilot_deployment_profile/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

func setupMockEnvironment() (*mocks.Mocks, *profileMocks.WindowsAutopilotDeploymentProfileMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	profileMock := &profileMocks.WindowsAutopilotDeploymentProfileMock{}
	profileMock.RegisterMocks()
	return mockClient, profileMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *profileMocks.WindowsAutopilotDeploymentProfileMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	profileMock := &profileMocks.WindowsAutopilotDeploymentProfileMock{}
	profileMock.RegisterErrorMocks()
	return mockClient, profileMock
}

func testCheckExists(resourceName string) resource.TestCheckFunc {
	return resource.TestCheckResourceAttrSet(resourceName, "id")
}

func TestWindowsAutopilotDeploymentProfileResource_Schema(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, profileMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer profileMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					// Basic attributes
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.test", "display_name", "Test Windows Autopilot Deployment Profile"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.test", "description", "Test description for Windows Autopilot deployment profile"),
					resource.TestMatchResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.test", "id", regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.test", "device_join_type", "microsoft_entra_joined"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.test", "locale", "os-default"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.test", "hardware_hash_extraction_enabled", "false"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.test", "preprovisioning_allowed", "false"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.test", "device_name_template", ""),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.test", "role_scope_tag_ids.#", "1"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.test", "role_scope_tag_ids.*", "0"),

					// Out of box experience settings
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.test", "out_of_box_experience_setting.privacy_settings_hidden", "false"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.test", "out_of_box_experience_setting.eula_hidden", "false"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.test", "out_of_box_experience_setting.user_type", "administrator"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.test", "out_of_box_experience_setting.device_usage_type", "singleUser"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.test", "out_of_box_experience_setting.keyboard_selection_page_skipped", "false"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.test", "out_of_box_experience_setting.escape_link_hidden", "false"),

					// Enrollment status screen settings
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.test", "enrollment_status_screen_settings.hide_installation_progress", "false"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.test", "enrollment_status_screen_settings.allow_device_use_before_profile_and_app_install_complete", "false"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.test", "enrollment_status_screen_settings.block_device_setup_retry_by_user", "false"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.test", "enrollment_status_screen_settings.allow_log_collection_on_install_failure", "false"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.test", "enrollment_status_screen_settings.allow_device_use_on_install_failure", "false"),

					// Computed attributes
					testCheckExists("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.test"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.test", "created_date_time"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.test", "last_modified_date_time"),
				),
			},
		},
	})
}

func TestWindowsAutopilotDeploymentProfileResource_MaximalSettings(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, profileMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer profileMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMaximal(),
				Check: resource.ComposeTestCheckFunc(
					// Basic attributes
					testCheckExists("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.test"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.test", "display_name", "Test Windows Autopilot Deployment Profile - maximal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.test", "description", "Test description for Windows Autopilot deployment profile with maximal configuration"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.test", "device_join_type", "microsoft_entra_joined"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.test", "locale", "en-US"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.test", "hardware_hash_extraction_enabled", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.test", "preprovisioning_allowed", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.test", "device_name_template", "AUTO-%SERIAL%"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.test", "device_type", "windowsPc"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.test", "role_scope_tag_ids.#", "3"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.test", "role_scope_tag_ids.*", "0"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.test", "role_scope_tag_ids.*", "1"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.test", "role_scope_tag_ids.*", "2"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.test", "management_service_app_id", "12345678-1234-1234-1234-123456789abc"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.test", "hybrid_azure_ad_join_skip_connectivity_check", "true"),

					// Out of box experience settings
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.test", "out_of_box_experience_setting.privacy_settings_hidden", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.test", "out_of_box_experience_setting.eula_hidden", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.test", "out_of_box_experience_setting.user_type", "standard"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.test", "out_of_box_experience_setting.device_usage_type", "shared"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.test", "out_of_box_experience_setting.keyboard_selection_page_skipped", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.test", "out_of_box_experience_setting.escape_link_hidden", "true"),

					// Enrollment status screen settings
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.test", "enrollment_status_screen_settings.hide_installation_progress", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.test", "enrollment_status_screen_settings.allow_device_use_before_profile_and_app_install_complete", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.test", "enrollment_status_screen_settings.block_device_setup_retry_by_user", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.test", "enrollment_status_screen_settings.allow_log_collection_on_install_failure", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.test", "enrollment_status_screen_settings.custom_error_message", "Custom error message for installation failure"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.test", "enrollment_status_screen_settings.install_progress_timeout_in_minutes", "60"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.test", "enrollment_status_screen_settings.allow_device_use_on_install_failure", "true"),

					// Computed attributes
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.test", "created_date_time"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.test", "last_modified_date_time"),
				),
			},
		},
	})
}

func TestWindowsAutopilotDeploymentProfileResource_Update(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, profileMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer profileMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.test", "display_name", "Test Windows Autopilot Deployment Profile"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.test", "locale", "os-default"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.test", "hardware_hash_extraction_enabled", "false"),
				),
			},
			{
				Config: testConfigUpdate(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.test", "display_name", "Updated Windows Autopilot Deployment Profile"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.test", "description", "Updated description for Windows Autopilot deployment profile"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.test", "locale", "fr-FR"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.test", "hardware_hash_extraction_enabled", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.test", "device_name_template", "UPD-%SERIAL%"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.test", "preprovisioning_allowed", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.test", "out_of_box_experience_setting.privacy_settings_hidden", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.test", "out_of_box_experience_setting.user_type", "standard"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.test", "enrollment_status_screen_settings.custom_error_message", "Updated error message"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.test", "enrollment_status_screen_settings.install_progress_timeout_in_minutes", "90"),
				),
			},
		},
	})
}

func TestWindowsAutopilotDeploymentProfileResource_Import(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, profileMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer profileMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.test"),
				),
			},
			{
				ResourceName:      "microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"timeouts",
				},
			},
		},
	})
}

func TestWindowsAutopilotDeploymentProfileResource_ErrorHandling(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, profileMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer profileMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testConfigMinimal(),
				ExpectError: regexp.MustCompile("Invalid request body"),
			},
		},
	})
}

func testConfigMinimal() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/windows_autopilot_deployment_profile_minimal.tf")
	if err != nil {
		panic("failed to load minimal config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigMaximal() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/windows_autopilot_deployment_profile_maximal.tf")
	if err != nil {
		panic("failed to load maximal config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigUpdate() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/windows_autopilot_deployment_profile_update.tf")
	if err != nil {
		panic("failed to load update config: " + err.Error())
	}
	return unitTestConfig
}

func TestWindowsAutopilotDeploymentProfileResource_DeviceNameValidation(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, profileMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer profileMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testConfigInvalidDeviceName(),
				ExpectError: regexp.MustCompile("Device name template length.*exceeds maximum allowed length of 15"),
			},
		},
	})
}

func testConfigInvalidDeviceName() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/windows_autopilot_deployment_profile_invalid_device_name.tf")
	if err != nil {
		panic("failed to load invalid device name config: " + err.Error())
	}
	return unitTestConfig
}