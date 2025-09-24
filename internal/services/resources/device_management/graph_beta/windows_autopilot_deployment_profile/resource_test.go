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

func TestWindowsAutopilotDeploymentProfileResource_UserDriven(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, profileMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer profileMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigUserDriven(),
				Check: resource.ComposeTestCheckFunc(
					// Basic attributes
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven", "display_name", "unit-test-user-driven"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven", "description", "user driven autopilot profile with os default locale"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven", "device_name_template", "thing-%RAND:5%"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven", "locale", "os-default"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven", "device_type", "windowsPc"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven", "device_join_type", "microsoft_entra_joined"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven", "preprovisioning_allowed", "false"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven", "hardware_hash_extraction_enabled", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven", "hybrid_azure_ad_join_skip_connectivity_check", "false"),

					// OOBE Settings
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven", "out_of_box_experience_setting.device_usage_type", "singleUser"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven", "out_of_box_experience_setting.privacy_settings_hidden", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven", "out_of_box_experience_setting.eula_hidden", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven", "out_of_box_experience_setting.user_type", "standard"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven", "out_of_box_experience_setting.keyboard_selection_page_skipped", "true"),

					// Role Scope Tags
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven", "role_scope_tag_ids.#", "1"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven", "role_scope_tag_ids.0", "0"),

					// Computed attributes
					testCheckExists("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven", "created_date_time"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven", "last_modified_date_time"),
				),
			},
		},
	})
}

func TestWindowsAutopilotDeploymentProfileResource_UserDrivenJapanesePreprovisioned(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, profileMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer profileMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigUserDrivenJapanesePreprovisioned(),
				Check: resource.ComposeTestCheckFunc(
					// Basic attributes
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven_japanese_preprovisioned", "display_name", "unit-test-user-driven-japanese-preprovisioned"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven_japanese_preprovisioned", "description", "user driven autopilot profile with japanese locale and allow pre provisioned deployment"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven_japanese_preprovisioned", "device_name_template", "thing-%RAND:3%"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven_japanese_preprovisioned", "locale", "ja-JP"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven_japanese_preprovisioned", "preprovisioning_allowed", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven_japanese_preprovisioned", "hardware_hash_extraction_enabled", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven_japanese_preprovisioned", "device_type", "windowsPc"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven_japanese_preprovisioned", "device_join_type", "microsoft_entra_joined"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven_japanese_preprovisioned", "hybrid_azure_ad_join_skip_connectivity_check", "false"),

					// OOBE Settings
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven_japanese_preprovisioned", "out_of_box_experience_setting.device_usage_type", "singleUser"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven_japanese_preprovisioned", "out_of_box_experience_setting.privacy_settings_hidden", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven_japanese_preprovisioned", "out_of_box_experience_setting.eula_hidden", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven_japanese_preprovisioned", "out_of_box_experience_setting.user_type", "standard"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven_japanese_preprovisioned", "out_of_box_experience_setting.keyboard_selection_page_skipped", "true"),

					// Role Scope Tags
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven_japanese_preprovisioned", "role_scope_tag_ids.#", "1"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven_japanese_preprovisioned", "role_scope_tag_ids.0", "0"),

					testCheckExists("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven_japanese_preprovisioned"),
				),
			},
		},
	})
}

func TestWindowsAutopilotDeploymentProfileResource_SelfDeploying(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, profileMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer profileMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigSelfDeploying(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.self_deploying", "display_name", "unit-test-self-deploying"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.self_deploying", "description", "self deploying autopilot profile with os default locale"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.self_deploying", "device_name_template", "thing-%RAND:2%"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.self_deploying", "locale", "os-default"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.self_deploying", "preprovisioning_allowed", "false"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.self_deploying", "device_type", "windowsPc"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.self_deploying", "hardware_hash_extraction_enabled", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.self_deploying", "device_join_type", "microsoft_entra_joined"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.self_deploying", "hybrid_azure_ad_join_skip_connectivity_check", "false"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.self_deploying", "out_of_box_experience_setting.device_usage_type", "shared"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.self_deploying", "out_of_box_experience_setting.privacy_settings_hidden", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.self_deploying", "out_of_box_experience_setting.eula_hidden", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.self_deploying", "out_of_box_experience_setting.user_type", "standard"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.self_deploying", "out_of_box_experience_setting.keyboard_selection_page_skipped", "true"),
					testCheckExists("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.self_deploying"),
				),
			},
		},
	})
}

func TestWindowsAutopilotDeploymentProfileResource_HoloLens(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, profileMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer profileMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigHoloLens(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.hololens", "display_name", "unit-test-hololens"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.hololens", "description", "hololens autopilot profile with os default locale"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.hololens", "device_name_template", "thing-%RAND:2%"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.hololens", "device_type", "holoLens"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.hololens", "locale", "zh-HK"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.hololens", "hardware_hash_extraction_enabled", "false"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.hololens", "preprovisioning_allowed", "false"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.hololens", "device_join_type", "microsoft_entra_joined"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.hololens", "hybrid_azure_ad_join_skip_connectivity_check", "false"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.hololens", "out_of_box_experience_setting.device_usage_type", "shared"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.hololens", "out_of_box_experience_setting.privacy_settings_hidden", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.hololens", "out_of_box_experience_setting.eula_hidden", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.hololens", "out_of_box_experience_setting.user_type", "standard"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.hololens", "out_of_box_experience_setting.keyboard_selection_page_skipped", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.hololens", "role_scope_tag_ids.#", "1"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.hololens", "role_scope_tag_ids.0", "0"),
					testCheckExists("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.hololens"),
				),
			},
		},
	})
}

func TestWindowsAutopilotDeploymentProfileResource_ValidationError(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, profileMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer profileMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testConfigUserDriven(),
				ExpectError: regexp.MustCompile("BadRequest"),
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
				Config: testConfigUserDriven(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven", "description", "user driven autopilot profile with os default locale"),
					testCheckExists("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven"),
				),
			},
			{
				Config: testConfigUserDrivenUpdated(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven", "description", "Updated unit test description"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven", "locale", "en-US"),
					testCheckExists("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven"),
				),
			},
		},
	})
}

// Configuration functions
func testConfigUserDriven() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/user_driven_minimal.tf")
	if err != nil {
		panic("failed to load user driven config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigUserDrivenJapanesePreprovisioned() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/user_driven_japanese_preprovisioned.tf")
	if err != nil {
		panic("failed to load user driven japanese preprovisioned config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigSelfDeploying() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/self_deploying.tf")
	if err != nil {
		panic("failed to load self deploying config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigHoloLens() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/hololens.tf")
	if err != nil {
		panic("failed to load hololens config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigUserDrivenUpdated() string {
	return `
resource "microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile" "user_driven" {
  display_name                                 = "acc_test_user_driven_japanese_preprovisioned"
  description                                  = "user driven autopilot profile with japanese locale and allow pre provisioned deployment"
  device_name_template                         = "thing-%RAND:3%"
  locale                                       = "ja-JP"
  preprovisioning_allowed                      = true
  device_type                                  = "windowsPc"
  hardware_hash_extraction_enabled             = true
  role_scope_tag_ids                           = ["0"]
  device_join_type                             = "microsoft_entra_joined"
  hybrid_azure_ad_join_skip_connectivity_check = false

  out_of_box_experience_setting = {
    device_usage_type               = "singleUser"
    privacy_settings_hidden         = true
    eula_hidden                     = true
    user_type                       = "standard"
    keyboard_selection_page_skipped = true
  }

  assignments = [
    {
      type = "allDevicesAssignmentTarget"
    }
  ]
}
`
}