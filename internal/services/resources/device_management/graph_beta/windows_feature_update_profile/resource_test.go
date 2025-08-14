package graphBetaWindowsFeatureUpdateProfile_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	featureMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/windows_feature_update_profile/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

func TestMain(m *testing.M) {
	exitCode := m.Run()
	os.Exit(exitCode)
}

func setupMockEnvironment() (*mocks.Mocks, *featureMocks.WindowsFeatureUpdateProfileMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	featureMock := &featureMocks.WindowsFeatureUpdateProfileMock{}
	featureMock.RegisterMocks()
	return mockClient, featureMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *featureMocks.WindowsFeatureUpdateProfileMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	featureMock := &featureMocks.WindowsFeatureUpdateProfileMock{}
	featureMock.RegisterErrorMocks()
	return mockClient, featureMock
}

func testCheckExists(resourceName string) resource.TestCheckFunc {
	return resource.TestCheckResourceAttrSet(resourceName, "id")
}

func TestWindowsFeatureUpdateProfileResource_Schema(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, featureMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer featureMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_feature_update_profile.minimal", "display_name", "Test Minimal Windows Feature Update Profile - Unique"),
					resource.TestMatchResourceAttr("microsoft365_graph_beta_device_management_windows_feature_update_profile.minimal", "id", regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_feature_update_profile.minimal", "role_scope_tag_ids.#", "1"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_windows_feature_update_profile.minimal", "role_scope_tag_ids.*", "0"),
				),
			},
		},
	})
}

func TestWindowsFeatureUpdateProfileResource_RolloutScenarios(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, featureMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer featureMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigScheduledRollout(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_windows_feature_update_profile.scheduled_rollout"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_feature_update_profile.scheduled_rollout", "display_name", "Test Scheduled Rollout Windows Feature Update Profile - Unique"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_feature_update_profile.scheduled_rollout", "feature_update_version", "Windows 11, version 22H2"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_feature_update_profile.scheduled_rollout", "rollout_settings.offer_start_date_time_in_utc", "2029-04-01T00:00:00Z"),
				),
			},
			{
				Config: testConfigSpecificDateRollout(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_windows_feature_update_profile.specific_date_rollout"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_feature_update_profile.specific_date_rollout", "display_name", "Test Specific Date Rollout Windows Feature Update Profile - Unique"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_feature_update_profile.specific_date_rollout", "feature_update_version", "Windows 11, version 24H2"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_feature_update_profile.specific_date_rollout", "rollout_settings.offer_start_date_time_in_utc", "2029-08-01T00:00:00Z"),
				),
			},
			{
				Config: testConfigImmediateRollout(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_windows_feature_update_profile.immediate_rollout"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_feature_update_profile.immediate_rollout", "display_name", "Test Immediate Rollout Windows Feature Update Profile - Unique"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_feature_update_profile.immediate_rollout", "feature_update_version", "Windows 10, version 22H2"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_feature_update_profile.immediate_rollout", "install_feature_updates_optional", "false"),
				),
			},
		},
	})
}

func TestWindowsFeatureUpdateProfileResource_GroupAssignments(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, featureMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer featureMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigGroupAssignments(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_windows_feature_update_profile.group_assignments"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_feature_update_profile.group_assignments", "display_name", "Test Group Assignments Windows Feature Update Profile - Unique"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_feature_update_profile.group_assignments", "assignments.#", "2"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_feature_update_profile.group_assignments", "assignments.0.type", "groupAssignmentTarget"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_feature_update_profile.group_assignments", "assignments.1.type", "groupAssignmentTarget"),
				),
			},
		},
	})
}

func TestWindowsFeatureUpdateProfileResource_ExclusionAssignment(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, featureMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer featureMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigExclusionAssignment(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_windows_feature_update_profile.exclusion_assignment"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_feature_update_profile.exclusion_assignment", "display_name", "Test Exclusion Assignment Windows Feature Update Profile - Unique"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_feature_update_profile.exclusion_assignment", "assignments.#", "1"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_feature_update_profile.exclusion_assignment", "assignments.0.type", "exclusionGroupAssignmentTarget"),
				),
			},
		},
	})
}

func TestWindowsFeatureUpdateProfileResource_AllFeatureUpdateVersions(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, featureMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer featureMock.CleanupMockState()

	featureVersions := []string{
		"Windows 11, version 24H2",
		"Windows 11, version 23H2", 
		"Windows 11, version 22H2",
		"Windows 10, version 22H2",
	}

	for i, version := range featureVersions {
		t.Run(fmt.Sprintf("FeatureVersion_%d_%s", i, version), func(t *testing.T) {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: testConfigWithFeatureVersion(version, i),
						Check: resource.ComposeTestCheckFunc(
							testCheckExists(fmt.Sprintf("microsoft365_graph_beta_device_management_windows_feature_update_profile.test_%d", i)),
							resource.TestCheckResourceAttr(fmt.Sprintf("microsoft365_graph_beta_device_management_windows_feature_update_profile.test_%d", i), "feature_update_version", version),
						),
					},
				},
			})
		})
	}
}

func TestWindowsFeatureUpdateProfileResource_ErrorHandling(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, featureMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer featureMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testConfigMinimal(),
				ExpectError: regexp.MustCompile("Invalid Windows Feature Update Profile data"),
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


func testConfigScheduledRollout() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_gradual_rollout.tf")
	if err != nil {
		panic("failed to load scheduled rollout config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigSpecificDateRollout() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_specific_date_rollout.tf")
	if err != nil {
		panic("failed to load specific date rollout config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigImmediateRollout() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_immediate_rollout.tf")
	if err != nil {
		panic("failed to load immediate rollout config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigWithFeatureVersion(version string, index int) string {
	return fmt.Sprintf(`
resource "microsoft365_graph_beta_device_management_windows_feature_update_profile" "test_%d" {
  display_name           = "Test Feature Update Profile %s - Unique"
  feature_update_version = "%s"

  install_feature_updates_optional                         = false
  install_latest_windows10_on_windows11_ineligible_device = false

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}`, index, version, version)
}

func testConfigGroupAssignments() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_with_group_assignments.tf")
	if err != nil {
		panic("failed to load group assignments config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigExclusionAssignment() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_with_exclusion_assignment.tf")
	if err != nil {
		panic("failed to load exclusion assignment config: " + err.Error())
	}
	return unitTestConfig
}
