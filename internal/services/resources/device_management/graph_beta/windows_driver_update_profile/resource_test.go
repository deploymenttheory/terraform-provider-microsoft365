package graphBetaWindowsDriverUpdateProfile_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	driverMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/windows_driver_update_profile/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

func setupMockEnvironment() (*mocks.Mocks, *driverMocks.WindowsDriverUpdateProfileMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	driverMock := &driverMocks.WindowsDriverUpdateProfileMock{}
	driverMock.RegisterMocks()
	return mockClient, driverMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *driverMocks.WindowsDriverUpdateProfileMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	driverMock := &driverMocks.WindowsDriverUpdateProfileMock{}
	driverMock.RegisterErrorMocks()
	return mockClient, driverMock
}

func testCheckExists(resourceName string) resource.TestCheckFunc {
	return resource.TestCheckResourceAttrSet(resourceName, "id")
}

func TestUnitResourceWindowsDriverUpdateProfile_01_Schema(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, driverMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer driverMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_driver_update_profile.minimal", "display_name", "Test Minimal Windows Driver Update Profile - Unique"),
					resource.TestMatchResourceAttr("microsoft365_graph_beta_device_management_windows_driver_update_profile.minimal", "id", regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_driver_update_profile.minimal", "role_scope_tag_ids.#", "1"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_windows_driver_update_profile.minimal", "role_scope_tag_ids.*", "0"),
				),
			},
		},
	})
}

func TestUnitResourceWindowsDriverUpdateProfile_02_ApprovalTypes(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, driverMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer driverMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigManualApproval(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_windows_driver_update_profile.manual_approval"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_driver_update_profile.manual_approval", "display_name", "Test Manual Approval Windows Driver Update Profile - Unique"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_driver_update_profile.manual_approval", "approval_type", "manual"),
				),
			},
			{
				Config: testConfigAutomaticApproval(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_windows_driver_update_profile.automatic_approval"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_driver_update_profile.automatic_approval", "display_name", "Test Automatic Approval Windows Driver Update Profile - Unique"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_driver_update_profile.automatic_approval", "approval_type", "automatic"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_driver_update_profile.automatic_approval", "deployment_deferral_in_days", "5"),
				),
			},
		},
	})
}

func TestUnitResourceWindowsDriverUpdateProfile_03_GroupAssignments(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, driverMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer driverMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigGroupAssignments(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_windows_driver_update_profile.group_assignments"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_driver_update_profile.group_assignments", "display_name", "Test Group Assignments Windows Driver Update Profile - Unique"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_driver_update_profile.group_assignments", "assignments.#", "2"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_driver_update_profile.group_assignments", "assignments.0.type", "groupAssignmentTarget"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_driver_update_profile.group_assignments", "assignments.1.type", "groupAssignmentTarget"),
				),
			},
		},
	})
}

func TestUnitResourceWindowsDriverUpdateProfile_04_ExclusionAssignment(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, driverMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer driverMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigExclusionAssignment(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_windows_driver_update_profile.exclusion_assignment"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_driver_update_profile.exclusion_assignment", "display_name", "Test Exclusion Assignment Windows Driver Update Profile - Unique"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_driver_update_profile.exclusion_assignment", "assignments.#", "1"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_driver_update_profile.exclusion_assignment", "assignments.0.type", "exclusionGroupAssignmentTarget"),
				),
			},
		},
	})
}

func TestUnitResourceWindowsDriverUpdateProfile_05_ErrorHandling(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, driverMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer driverMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testConfigMinimal(),
				ExpectError: regexp.MustCompile("Invalid Windows Driver Update Profile data"),
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

func testConfigManualApproval() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_manual_approval.tf")
	if err != nil {
		panic("failed to load manual approval config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigAutomaticApproval() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_automatic_approval.tf")
	if err != nil {
		panic("failed to load automatic approval config: " + err.Error())
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
