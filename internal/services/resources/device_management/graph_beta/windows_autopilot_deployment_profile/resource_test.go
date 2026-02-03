package graphBetaWindowsAutopilotDeploymentProfile_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
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

func loadUnitTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/unit/" + filename)
	if err != nil {
		panic("failed to load unit test config " + filename + ": " + err.Error())
	}
	return config
}

func TestUnitResourceWindowsAutopilotDeploymentProfile_01_SelfDeployingOsDefaultLocale(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, profileMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer profileMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("01_self_deploying_os_default_locale.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".user_driven").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".user_driven").Key("display_name").HasValue("acc test user driven autopilot profile with os default locale"),
					check.That(resourceType+".user_driven").Key("description").HasValue("user driven autopilot profile with os default locale"),
					check.That(resourceType+".user_driven").Key("device_name_template").HasValue("thing-%RAND:5%"),
					check.That(resourceType+".user_driven").Key("locale").HasValue("os-default"),
					check.That(resourceType+".user_driven").Key("device_type").HasValue("windowsPc"),
					check.That(resourceType+".user_driven").Key("device_join_type").HasValue("microsoft_entra_joined"),
					check.That(resourceType+".user_driven").Key("preprovisioning_allowed").HasValue("true"),
					check.That(resourceType+".user_driven").Key("hardware_hash_extraction_enabled").HasValue("true"),
					check.That(resourceType+".user_driven").Key("hybrid_azure_ad_join_skip_connectivity_check").HasValue("false"),

					// OOBE Settings
					check.That(resourceType+".user_driven").Key("out_of_box_experience_setting.device_usage_type").HasValue("singleUser"),
					check.That(resourceType+".user_driven").Key("out_of_box_experience_setting.privacy_settings_hidden").HasValue("true"),
					check.That(resourceType+".user_driven").Key("out_of_box_experience_setting.eula_hidden").HasValue("true"),
					check.That(resourceType+".user_driven").Key("out_of_box_experience_setting.user_type").HasValue("standard"),
					check.That(resourceType+".user_driven").Key("out_of_box_experience_setting.keyboard_selection_page_skipped").HasValue("true"),

					// Role Scope Tags
					check.That(resourceType+".user_driven").Key("role_scope_tag_ids.#").HasValue("2"),
					check.That(resourceType+".user_driven").Key("role_scope_tag_ids.0").HasValue("0"),
					check.That(resourceType+".user_driven").Key("role_scope_tag_ids.1").HasValue("1"),

					// Assignments
					check.That(resourceType+".user_driven").Key("assignments.#").HasValue("3"),
					check.That(resourceType+".user_driven").Key("assignments.0.type").HasValue("groupAssignmentTarget"),
					check.That(resourceType+".user_driven").Key("assignments.0.group_id").HasValue("00000000-0000-0000-0000-000000000001"),
					check.That(resourceType+".user_driven").Key("assignments.1.type").HasValue("groupAssignmentTarget"),
					check.That(resourceType+".user_driven").Key("assignments.1.group_id").HasValue("00000000-0000-0000-0000-000000000002"),
					check.That(resourceType+".user_driven").Key("assignments.2.type").HasValue("exclusionGroupAssignmentTarget"),
					check.That(resourceType+".user_driven").Key("assignments.2.group_id").HasValue("00000000-0000-0000-0000-000000000003"),

					// Computed attributes
					check.That(resourceType+".user_driven").Key("created_date_time").Exists(),
					check.That(resourceType+".user_driven").Key("last_modified_date_time").Exists(),
				),
			},
		},
	})
}

func TestUnitResourceWindowsAutopilotDeploymentProfile_02_UserDrivenHybridDomainJoin(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, profileMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer profileMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("02_user_driven_hybrid_domain_join.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".user_driven_japanese_preprovisioned_with_assignments").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".user_driven_japanese_preprovisioned_with_assignments").Key("display_name").HasValue("unit_test_user_driven_japanese_preprovisioned"),
					check.That(resourceType+".user_driven_japanese_preprovisioned_with_assignments").Key("description").HasValue("user driven autopilot profile with japanese locale and allow pre provisioned deployment"),
					check.That(resourceType+".user_driven_japanese_preprovisioned_with_assignments").Key("locale").HasValue("ja-JP"),
					check.That(resourceType+".user_driven_japanese_preprovisioned_with_assignments").Key("preprovisioning_allowed").HasValue("true"),
					check.That(resourceType+".user_driven_japanese_preprovisioned_with_assignments").Key("hardware_hash_extraction_enabled").HasValue("true"),
					check.That(resourceType+".user_driven_japanese_preprovisioned_with_assignments").Key("device_type").HasValue("windowsPc"),
					check.That(resourceType+".user_driven_japanese_preprovisioned_with_assignments").Key("device_join_type").HasValue("microsoft_entra_hybrid_joined"),
					check.That(resourceType+".user_driven_japanese_preprovisioned_with_assignments").Key("hybrid_azure_ad_join_skip_connectivity_check").HasValue("true"),

					// OOBE Settings
					check.That(resourceType+".user_driven_japanese_preprovisioned_with_assignments").Key("out_of_box_experience_setting.device_usage_type").HasValue("singleUser"),
					check.That(resourceType+".user_driven_japanese_preprovisioned_with_assignments").Key("out_of_box_experience_setting.privacy_settings_hidden").HasValue("true"),
					check.That(resourceType+".user_driven_japanese_preprovisioned_with_assignments").Key("out_of_box_experience_setting.eula_hidden").HasValue("true"),
					check.That(resourceType+".user_driven_japanese_preprovisioned_with_assignments").Key("out_of_box_experience_setting.user_type").HasValue("standard"),
					check.That(resourceType+".user_driven_japanese_preprovisioned_with_assignments").Key("out_of_box_experience_setting.keyboard_selection_page_skipped").HasValue("true"),

					// Role Scope Tags
					check.That(resourceType+".user_driven_japanese_preprovisioned_with_assignments").Key("role_scope_tag_ids.#").HasValue("1"),
					check.That(resourceType+".user_driven_japanese_preprovisioned_with_assignments").Key("role_scope_tag_ids.0").HasValue("0"),

					// Assignments
					check.That(resourceType+".user_driven_japanese_preprovisioned_with_assignments").Key("assignments.#").HasValue("1"),
					check.That(resourceType+".user_driven_japanese_preprovisioned_with_assignments").Key("assignments.0.type").HasValue("allDevicesAssignmentTarget"),
				),
			},
		},
	})
}

func TestUnitResourceWindowsAutopilotDeploymentProfile_03_UserDrivenWithGroupAssignments(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, profileMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer profileMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("03_user_driven_with_group_assignments.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".user_driven").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".user_driven").Key("display_name").HasValue("user driven autopilot with group assignments"),
					check.That(resourceType+".user_driven").Key("description").HasValue("user driven autopilot profile with os default locale"),
					check.That(resourceType+".user_driven").Key("device_name_template").HasValue("thing-%RAND:5%"),
					check.That(resourceType+".user_driven").Key("locale").HasValue("os-default"),
					check.That(resourceType+".user_driven").Key("preprovisioning_allowed").HasValue("true"),
					check.That(resourceType+".user_driven").Key("device_type").HasValue("windowsPc"),
					check.That(resourceType+".user_driven").Key("hardware_hash_extraction_enabled").HasValue("true"),
					check.That(resourceType+".user_driven").Key("device_join_type").HasValue("microsoft_entra_joined"),
					check.That(resourceType+".user_driven").Key("hybrid_azure_ad_join_skip_connectivity_check").HasValue("false"),
					check.That(resourceType+".user_driven").Key("out_of_box_experience_setting.device_usage_type").HasValue("singleUser"),
					check.That(resourceType+".user_driven").Key("out_of_box_experience_setting.privacy_settings_hidden").HasValue("true"),
					check.That(resourceType+".user_driven").Key("out_of_box_experience_setting.eula_hidden").HasValue("true"),
					check.That(resourceType+".user_driven").Key("out_of_box_experience_setting.user_type").HasValue("standard"),
					check.That(resourceType+".user_driven").Key("out_of_box_experience_setting.keyboard_selection_page_skipped").HasValue("true"),

					// Role Scope Tags
					check.That(resourceType+".user_driven").Key("role_scope_tag_ids.#").HasValue("2"),
					check.That(resourceType+".user_driven").Key("role_scope_tag_ids.0").HasValue("0"),
					check.That(resourceType+".user_driven").Key("role_scope_tag_ids.1").HasValue("1"),

					// Assignments
					check.That(resourceType+".user_driven").Key("assignments.#").HasValue("3"),
					check.That(resourceType+".user_driven").Key("assignments.0.type").HasValue("groupAssignmentTarget"),
					check.That(resourceType+".user_driven").Key("assignments.0.group_id").HasValue("00000000-0000-0000-0000-000000000001"),
					check.That(resourceType+".user_driven").Key("assignments.1.type").HasValue("groupAssignmentTarget"),
					check.That(resourceType+".user_driven").Key("assignments.1.group_id").HasValue("00000000-0000-0000-0000-000000000002"),
					check.That(resourceType+".user_driven").Key("assignments.2.type").HasValue("exclusionGroupAssignmentTarget"),
					check.That(resourceType+".user_driven").Key("assignments.2.group_id").HasValue("00000000-0000-0000-0000-000000000003"),
				),
			},
		},
	})
}

func TestUnitResourceWindowsAutopilotDeploymentProfile_04_HoloLensWithAllDeviceAssignment(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, profileMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer profileMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("04_hololens_with_all_device_assignment.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".hololens_with_all_device_assignment").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".hololens_with_all_device_assignment").Key("display_name").HasValue("unit_test_hololens_with_all_device_assignment"),
					check.That(resourceType+".hololens_with_all_device_assignment").Key("description").HasValue("hololens autopilot profile with hk locale and all device assignment"),
					check.That(resourceType+".hololens_with_all_device_assignment").Key("device_name_template").HasValue("thing-%RAND:2%"),
					check.That(resourceType+".hololens_with_all_device_assignment").Key("device_type").HasValue("holoLens"),
					check.That(resourceType+".hololens_with_all_device_assignment").Key("locale").HasValue("zh-HK"),
					check.That(resourceType+".hololens_with_all_device_assignment").Key("hardware_hash_extraction_enabled").HasValue("false"),
					check.That(resourceType+".hololens_with_all_device_assignment").Key("preprovisioning_allowed").HasValue("false"),
					check.That(resourceType+".hololens_with_all_device_assignment").Key("device_join_type").HasValue("microsoft_entra_joined"),
					check.That(resourceType+".hololens_with_all_device_assignment").Key("hybrid_azure_ad_join_skip_connectivity_check").HasValue("false"),
					check.That(resourceType+".hololens_with_all_device_assignment").Key("out_of_box_experience_setting.device_usage_type").HasValue("shared"),
					check.That(resourceType+".hololens_with_all_device_assignment").Key("out_of_box_experience_setting.privacy_settings_hidden").HasValue("true"),
					check.That(resourceType+".hololens_with_all_device_assignment").Key("out_of_box_experience_setting.eula_hidden").HasValue("true"),
					check.That(resourceType+".hololens_with_all_device_assignment").Key("out_of_box_experience_setting.user_type").HasValue("standard"),
					check.That(resourceType+".hololens_with_all_device_assignment").Key("out_of_box_experience_setting.keyboard_selection_page_skipped").HasValue("true"),
					check.That(resourceType+".hololens_with_all_device_assignment").Key("role_scope_tag_ids.#").HasValue("1"),
					check.That(resourceType+".hololens_with_all_device_assignment").Key("role_scope_tag_ids.0").HasValue("0"),

					// Assignments
					check.That(resourceType+".hololens_with_all_device_assignment").Key("assignments.#").HasValue("1"),
					check.That(resourceType+".hololens_with_all_device_assignment").Key("assignments.0.type").HasValue("allDevicesAssignmentTarget"),
				),
			},
		},
	})
}

func TestUnitResourceWindowsAutopilotDeploymentProfile_05_ValidationError(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, profileMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer profileMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("01_self_deploying_os_default_locale.tf"),
				ExpectError: regexp.MustCompile("Bad Request - 400"),
			},
		},
	})
}

func TestUnitResourceWindowsAutopilotDeploymentProfile_06_Import(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, profileMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer profileMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("01_self_deploying_os_default_locale.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType + ".user_driven").Key("id").Exists(),
				),
			},
			{
				ResourceName:            resourceType + ".user_driven",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"hybrid_azure_ad_join_skip_connectivity_check"},
			},
		},
	})
}

func TestUnitResourceWindowsAutopilotDeploymentProfile_07_Update(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, profileMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer profileMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("01_self_deploying_os_default_locale.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".user_driven").Key("description").HasValue("user driven autopilot profile with os default locale"),
					check.That(resourceType+".user_driven").Key("id").Exists(),
				),
			},
			{
				Config: testConfigUserDrivenUpdated(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".user_driven").Key("display_name").HasValue("acc test user driven autopilot profile updated"),
					check.That(resourceType+".user_driven").Key("description").HasValue("Updated unit test description"),
					check.That(resourceType+".user_driven").Key("locale").HasValue("ja-JP"),
					check.That(resourceType+".user_driven").Key("id").Exists(),
				),
			},
		},
	})
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
