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

func TestWindowsAutopilotDeploymentProfileResource_SelfDeployingOsDefaultLocale(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, profileMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer profileMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfig01SelfDeployingOsDefaultLocale(),
				Check: resource.ComposeTestCheckFunc(
					// Basic attributes
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven", "display_name", "acc test user driven autopilot profile with os default locale"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven", "description", "user driven autopilot profile with os default locale"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven", "device_name_template", "thing-%RAND:5%"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven", "locale", "os-default"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven", "device_type", "windowsPc"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven", "device_join_type", "microsoft_entra_joined"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven", "preprovisioning_allowed", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven", "hardware_hash_extraction_enabled", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven", "hybrid_azure_ad_join_skip_connectivity_check", "false"),

					// OOBE Settings
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven", "out_of_box_experience_setting.device_usage_type", "singleUser"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven", "out_of_box_experience_setting.privacy_settings_hidden", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven", "out_of_box_experience_setting.eula_hidden", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven", "out_of_box_experience_setting.user_type", "standard"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven", "out_of_box_experience_setting.keyboard_selection_page_skipped", "true"),

					// Role Scope Tags
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven", "role_scope_tag_ids.#", "2"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven", "role_scope_tag_ids.0", "0"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven", "role_scope_tag_ids.1", "1"),

					// Assignments
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven", "assignments.#", "3"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven", "assignments.0.type", "groupAssignmentTarget"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven", "assignments.0.group_id", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven", "assignments.1.type", "groupAssignmentTarget"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven", "assignments.1.group_id", "00000000-0000-0000-0000-000000000002"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven", "assignments.2.type", "exclusionGroupAssignmentTarget"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven", "assignments.2.group_id", "00000000-0000-0000-0000-000000000003"),

					// Computed attributes
					testCheckExists("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven", "created_date_time"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven", "last_modified_date_time"),
				),
			},
		},
	})
}

func TestWindowsAutopilotDeploymentProfileResource_UserDrivenHybridDomainJoin(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, profileMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer profileMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfig02UserDrivenHybridDomainJoin(),
				Check: resource.ComposeTestCheckFunc(
					// Basic attributes
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven_japanese_preprovisioned_with_assignments", "display_name", "unit_test_user_driven_japanese_preprovisioned"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven_japanese_preprovisioned_with_assignments", "description", "user driven autopilot profile with japanese locale and allow pre provisioned deployment"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven_japanese_preprovisioned_with_assignments", "locale", "ja-JP"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven_japanese_preprovisioned_with_assignments", "preprovisioning_allowed", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven_japanese_preprovisioned_with_assignments", "hardware_hash_extraction_enabled", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven_japanese_preprovisioned_with_assignments", "device_type", "windowsPc"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven_japanese_preprovisioned_with_assignments", "device_join_type", "microsoft_entra_hybrid_joined"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven_japanese_preprovisioned_with_assignments", "hybrid_azure_ad_join_skip_connectivity_check", "true"),

					// OOBE Settings
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven_japanese_preprovisioned_with_assignments", "out_of_box_experience_setting.device_usage_type", "singleUser"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven_japanese_preprovisioned_with_assignments", "out_of_box_experience_setting.privacy_settings_hidden", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven_japanese_preprovisioned_with_assignments", "out_of_box_experience_setting.eula_hidden", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven_japanese_preprovisioned_with_assignments", "out_of_box_experience_setting.user_type", "standard"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven_japanese_preprovisioned_with_assignments", "out_of_box_experience_setting.keyboard_selection_page_skipped", "true"),

					// Role Scope Tags
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven_japanese_preprovisioned_with_assignments", "role_scope_tag_ids.#", "1"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven_japanese_preprovisioned_with_assignments", "role_scope_tag_ids.0", "0"),

					// Assignments
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven_japanese_preprovisioned_with_assignments", "assignments.#", "1"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven_japanese_preprovisioned_with_assignments", "assignments.0.type", "allDevicesAssignmentTarget"),

					testCheckExists("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven_japanese_preprovisioned_with_assignments"),
				),
			},
		},
	})
}

func TestWindowsAutopilotDeploymentProfileResource_UserDrivenWithGroupAssignments(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, profileMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer profileMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfig03UserDrivenWithGroupAssignments(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven", "display_name", "user driven autopilot with group assignments"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven", "description", "user driven autopilot profile with os default locale"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven", "device_name_template", "thing-%RAND:5%"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven", "locale", "os-default"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven", "preprovisioning_allowed", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven", "device_type", "windowsPc"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven", "hardware_hash_extraction_enabled", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven", "device_join_type", "microsoft_entra_joined"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven", "hybrid_azure_ad_join_skip_connectivity_check", "false"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven", "out_of_box_experience_setting.device_usage_type", "singleUser"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven", "out_of_box_experience_setting.privacy_settings_hidden", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven", "out_of_box_experience_setting.eula_hidden", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven", "out_of_box_experience_setting.user_type", "standard"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven", "out_of_box_experience_setting.keyboard_selection_page_skipped", "true"),

					// Role Scope Tags
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven", "role_scope_tag_ids.#", "2"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven", "role_scope_tag_ids.0", "0"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven", "role_scope_tag_ids.1", "1"),

					// Assignments
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven", "assignments.#", "3"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven", "assignments.0.type", "groupAssignmentTarget"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven", "assignments.0.group_id", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven", "assignments.1.type", "groupAssignmentTarget"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven", "assignments.1.group_id", "00000000-0000-0000-0000-000000000002"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven", "assignments.2.type", "exclusionGroupAssignmentTarget"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven", "assignments.2.group_id", "00000000-0000-0000-0000-000000000003"),

					testCheckExists("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven"),
				),
			},
		},
	})
}

func TestWindowsAutopilotDeploymentProfileResource_HoloLensWithAllDeviceAssignment(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, profileMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer profileMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfig04HoloLensWithAllDeviceAssignment(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.hololens_with_all_device_assignment", "display_name", "unit_test_hololens_with_all_device_assignment"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.hololens_with_all_device_assignment", "description", "hololens autopilot profile with hk locale and all device assignment"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.hololens_with_all_device_assignment", "device_name_template", "thing-%RAND:2%"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.hololens_with_all_device_assignment", "device_type", "holoLens"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.hololens_with_all_device_assignment", "locale", "zh-HK"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.hololens_with_all_device_assignment", "hardware_hash_extraction_enabled", "false"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.hololens_with_all_device_assignment", "preprovisioning_allowed", "false"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.hololens_with_all_device_assignment", "device_join_type", "microsoft_entra_joined"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.hololens_with_all_device_assignment", "hybrid_azure_ad_join_skip_connectivity_check", "false"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.hololens_with_all_device_assignment", "out_of_box_experience_setting.device_usage_type", "shared"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.hololens_with_all_device_assignment", "out_of_box_experience_setting.privacy_settings_hidden", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.hololens_with_all_device_assignment", "out_of_box_experience_setting.eula_hidden", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.hololens_with_all_device_assignment", "out_of_box_experience_setting.user_type", "standard"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.hololens_with_all_device_assignment", "out_of_box_experience_setting.keyboard_selection_page_skipped", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.hololens_with_all_device_assignment", "role_scope_tag_ids.#", "1"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.hololens_with_all_device_assignment", "role_scope_tag_ids.0", "0"),

					// Assignments
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.hololens_with_all_device_assignment", "assignments.#", "1"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.hololens_with_all_device_assignment", "assignments.0.type", "allDevicesAssignmentTarget"),

					testCheckExists("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.hololens_with_all_device_assignment"),
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
				Config:      testConfig01SelfDeployingOsDefaultLocale(),
				ExpectError: regexp.MustCompile("Bad Request - 400"),
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
				Config: testConfig01SelfDeployingOsDefaultLocale(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven"),
				),
			},
			{
				ResourceName:            "microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"hybrid_azure_ad_join_skip_connectivity_check"},
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
				Config: testConfig01SelfDeployingOsDefaultLocale(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven", "description", "user driven autopilot profile with os default locale"),
					testCheckExists("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven"),
				),
			},
			{
				Config: testConfigUserDrivenUpdated(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven", "display_name", "acc test user driven autopilot profile updated"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven", "description", "Updated unit test description"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven", "locale", "ja-JP"),
					testCheckExists("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven"),
				),
			},
		},
	})
}

// Configuration functions
func testConfig01SelfDeployingOsDefaultLocale() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/01_self_deploying_os_default_locale.tf")
	if err != nil {
		panic("failed to load 01_self_deploying_os_default_locale config: " + err.Error())
	}
	return unitTestConfig
}

func testConfig02UserDrivenHybridDomainJoin() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/02_user_driven_hybrid_domain_join.tf")
	if err != nil {
		panic("failed to load 02_user_driven_hybrid_domain_join config: " + err.Error())
	}
	return unitTestConfig
}

func testConfig03UserDrivenWithGroupAssignments() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/03_user_driven_with_group_assignments.tf")
	if err != nil {
		panic("failed to load 03_user_driven_with_group_assignments config: " + err.Error())
	}
	return unitTestConfig
}

func testConfig04HoloLensWithAllDeviceAssignment() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/04_hololens_with_all_device_assignment.tf")
	if err != nil {
		panic("failed to load 04_hololens_with_all_device_assignment config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigUserDrivenUpdated() string {
	return `
resource "microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile" "user_driven" {
  display_name                                 = "acc test user driven autopilot profile updated"
  description                                  = "Updated unit test description"
  device_name_template                         = "thing-%RAND:3%"
  locale                                       = "ja-JP"
  preprovisioning_allowed                      = true
  device_type                                  = "windowsPc"
  hardware_hash_extraction_enabled             = true
  role_scope_tag_ids                           = ["0", "1"]
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
      type     = "groupAssignmentTarget"
      group_id = "00000000-0000-0000-0000-000000000001"
    },
    {
      type = "allDevicesAssignmentTarget"
    }
  ]
}
`
}