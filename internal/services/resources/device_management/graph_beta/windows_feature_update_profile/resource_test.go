package graphBetaWindowsFeatureUpdateProfile_test

import (
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

func TestWindowsFeatureUpdateProfileResource_MinimalToMax(t *testing.T) {
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
					testCheckExists("microsoft365_graph_beta_device_management_windows_feature_update_profile.minimal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_feature_update_profile.minimal", "display_name", "Test Minimal Windows Feature Update Profile - Unique"),
				),
			},
			{
				ResourceName:      "microsoft365_graph_beta_device_management_windows_feature_update_profile.minimal",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testConfigMaximal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_windows_feature_update_profile.maximal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_feature_update_profile.maximal", "display_name", "Test Maximal Windows Feature Update Profile - Unique"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_feature_update_profile.maximal", "description", "Maximal Windows Feature Update Profile for testing with all features"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_feature_update_profile.maximal", "install_feature_updates_optional", "true"),
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
