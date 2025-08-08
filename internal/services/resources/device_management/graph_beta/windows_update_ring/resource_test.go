package graphBetaWindowsUpdateRing_test

import (
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	windowsUpdateRingMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/windows_update_ring/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

func TestMain(m *testing.M) {
	exitCode := m.Run()
	os.Exit(exitCode)
}

// setupMockEnvironment sets up the mock environment using centralized mocks
func setupMockEnvironment() (*mocks.Mocks, *windowsUpdateRingMocks.WindowsUpdateRingMock) {
	// Activate httpmock
	httpmock.Activate()

	// Create a new Mocks instance and register authentication mocks
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	// Register local mocks directly
	windowsUpdateRingMock := &windowsUpdateRingMocks.WindowsUpdateRingMock{}
	windowsUpdateRingMock.RegisterMocks()

	return mockClient, windowsUpdateRingMock
}

// setupErrorMockEnvironment sets up the mock environment for error testing
func setupErrorMockEnvironment() (*mocks.Mocks, *windowsUpdateRingMocks.WindowsUpdateRingMock) {
	// Activate httpmock
	httpmock.Activate()

	// Create a new Mocks instance and register authentication mocks
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	// Register error mocks
	windowsUpdateRingMock := &windowsUpdateRingMocks.WindowsUpdateRingMock{}
	windowsUpdateRingMock.RegisterErrorMocks()

	return mockClient, windowsUpdateRingMock
}

// testCheckExists is a basic check to ensure the resource exists in the state
func testCheckExists(resourceName string) resource.TestCheckFunc {
	return resource.TestCheckResourceAttrSet(resourceName, "id")
}

// resourceMinimalUnitTestData returns the minimal tf configuration for testing
func resourceMinimalUnitTestData() string {
	content, err := os.ReadFile(filepath.Join("tests", "terraform", "unit", "resource_minimal.tf"))
	if err != nil {
		return ""
	}
	return string(content)
}

// resourceMaximalUnitTestData returns the maximal tf configuration for testing
func resourceMaximalUnitTestData() string {
	content, err := os.ReadFile(filepath.Join("tests", "terraform", "unit", "resource_maximal.tf"))
	if err != nil {
		return ""
	}
	return string(content)
}

// TestWindowsUpdateRingResource_Schema validates the resource schema
func TestWindowsUpdateRingResource_Schema(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, windowsUpdateRingMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer windowsUpdateRingMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: resourceMinimalUnitTestData(),
				Check: resource.ComposeTestCheckFunc(
					// Check required attributes
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_update_ring.minimal", "display_name", "Test Minimal Windows Update Ring - Unique"),

					// Check computed attributes are set
					resource.TestMatchResourceAttr("microsoft365_graph_beta_device_management_windows_update_ring.minimal", "id", regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_update_ring.minimal", "role_scope_tag_ids.#", "1"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_windows_update_ring.minimal", "role_scope_tag_ids.*", "0"),
				),
			},
		},
	})
}

// TestWindowsUpdateRingResource_Minimal tests basic CRUD operations
func TestWindowsUpdateRingResource_Minimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, windowsUpdateRingMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer windowsUpdateRingMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: resourceMinimalUnitTestData(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_windows_update_ring.minimal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_update_ring.minimal", "display_name", "Test Minimal Windows Update Ring - Unique"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "microsoft365_graph_beta_device_management_windows_update_ring.minimal",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: resourceMaximalUnitTestData(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_windows_update_ring.maximal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_update_ring.maximal", "display_name", "Test Maximal Windows Update Ring - Unique"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_update_ring.maximal", "description", "Maximal Windows update ring for testing with all features"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_update_ring.maximal", "role_scope_tag_ids.#", "2"),
				),
			},
		},
	})
}

// TestWindowsUpdateRingResource_UpdateInPlace tests in-place updates
func TestWindowsUpdateRingResource_UpdateInPlace(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, windowsUpdateRingMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer windowsUpdateRingMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: resourceMinimalUnitTestData(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_windows_update_ring.minimal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_update_ring.minimal", "display_name", "Test Minimal Windows Update Ring - Unique"),
				),
			},
			{
				Config: resourceMaximalUnitTestData(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_windows_update_ring.maximal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_update_ring.maximal", "display_name", "Test Maximal Windows Update Ring - Unique"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_update_ring.maximal", "description", "Maximal Windows update ring for testing with all features"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_update_ring.maximal", "role_scope_tag_ids.#", "2"),
				),
			},
		},
	})
}

// TestWindowsUpdateRingResource_RequiredFields tests required field validation
func TestWindowsUpdateRingResource_RequiredFields(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, windowsUpdateRingMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer windowsUpdateRingMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "microsoft365_graph_beta_device_management_windows_update_ring" "test" {
  # Missing display_name
  microsoft_update_service_allowed = true
  drivers_excluded = false
}
`,
				ExpectError: regexp.MustCompile(`The argument "display_name" is required`),
			},
		},
	})
}

// TestWindowsUpdateRingResource_ErrorHandling tests error scenarios
func TestWindowsUpdateRingResource_ErrorHandling(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, windowsUpdateRingMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer windowsUpdateRingMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "microsoft365_graph_beta_device_management_windows_update_ring" "test" {
  display_name = "Test Windows Update Ring"
  microsoft_update_service_allowed = true
  drivers_excluded = false
  quality_updates_deferral_period_in_days = 0
  feature_updates_deferral_period_in_days = 0
  allow_windows11_upgrade = true
  skip_checks_before_restart = false
  automatic_update_mode = "userDefined"
  feature_updates_rollback_window_in_days = 10
}
`,
				ExpectError: regexp.MustCompile(`Invalid Windows update ring data|BadRequest`),
			},
		},
	})
}

// TestWindowsUpdateRingResource_RoleScopeTags tests role scope tags handling
func TestWindowsUpdateRingResource_RoleScopeTags(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, windowsUpdateRingMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer windowsUpdateRingMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "microsoft365_graph_beta_device_management_windows_update_ring" "test" {
  display_name = "Test Windows Update Ring"
  microsoft_update_service_allowed = true
  drivers_excluded = false
  quality_updates_deferral_period_in_days = 0
  feature_updates_deferral_period_in_days = 0
  allow_windows11_upgrade = true
  skip_checks_before_restart = false
  automatic_update_mode = "userDefined"
  feature_updates_rollback_window_in_days = 10
  role_scope_tag_ids = ["0", "1", "2"]
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_update_ring.test", "role_scope_tag_ids.#", "3"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_windows_update_ring.test", "role_scope_tag_ids.*", "0"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_windows_update_ring.test", "role_scope_tag_ids.*", "1"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_windows_update_ring.test", "role_scope_tag_ids.*", "2"),
				),
			},
		},
	})
}

// TestWindowsUpdateRingResource_GroupAssignments tests group assignment functionality
func TestWindowsUpdateRingResource_GroupAssignments(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, windowsUpdateRingMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer windowsUpdateRingMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigGroupAssignments(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_windows_update_ring.group_assignments"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_update_ring.group_assignments", "display_name", "Test Group Assignments Windows Update Ring - Unique"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_update_ring.group_assignments", "assignments.#", "2"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_update_ring.group_assignments", "assignments.0.type", "groupAssignmentTarget"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_update_ring.group_assignments", "assignments.1.type", "groupAssignmentTarget"),
				),
			},
		},
	})
}

// TestWindowsUpdateRingResource_AllUsersAssignment tests all licensed users assignment functionality
func TestWindowsUpdateRingResource_AllUsersAssignment(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, windowsUpdateRingMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer windowsUpdateRingMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigAllUsersAssignment(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_windows_update_ring.all_users_assignment"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_update_ring.all_users_assignment", "display_name", "Test All Users Assignment Windows Update Ring - Unique"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_update_ring.all_users_assignment", "assignments.#", "1"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_update_ring.all_users_assignment", "assignments.0.type", "allLicensedUsersAssignmentTarget"),
				),
			},
		},
	})
}

// TestWindowsUpdateRingResource_AllDevicesAssignment tests all devices assignment functionality
func TestWindowsUpdateRingResource_AllDevicesAssignment(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, windowsUpdateRingMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer windowsUpdateRingMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigAllDevicesAssignment(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_windows_update_ring.all_devices_assignment"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_update_ring.all_devices_assignment", "display_name", "Test All Devices Assignment Windows Update Ring - Unique"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_update_ring.all_devices_assignment", "assignments.#", "1"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_update_ring.all_devices_assignment", "assignments.0.type", "allDevicesAssignmentTarget"),
				),
			},
		},
	})
}

// TestWindowsUpdateRingResource_ExclusionAssignment tests exclusion group assignment functionality
func TestWindowsUpdateRingResource_ExclusionAssignment(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, windowsUpdateRingMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer windowsUpdateRingMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigExclusionAssignment(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_windows_update_ring.exclusion_assignment"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_update_ring.exclusion_assignment", "display_name", "Test Exclusion Assignment Windows Update Ring - Unique"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_update_ring.exclusion_assignment", "assignments.#", "1"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_update_ring.exclusion_assignment", "assignments.0.type", "exclusionGroupAssignmentTarget"),
				),
			},
		},
	})
}

// TestWindowsUpdateRingResource_AllAssignmentTypes tests all assignment types together
func TestWindowsUpdateRingResource_AllAssignmentTypes(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, windowsUpdateRingMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer windowsUpdateRingMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigAllAssignmentTypes(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_windows_update_ring.all_assignment_types"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_update_ring.all_assignment_types", "display_name", "Test All Assignment Types Windows Update Ring - Unique"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_update_ring.all_assignment_types", "assignments.#", "5"),
					// Verify all assignment types are present
					resource.TestCheckTypeSetElemNestedAttrs("microsoft365_graph_beta_device_management_windows_update_ring.all_assignment_types", "assignments.*", map[string]string{"type": "groupAssignmentTarget"}),
					resource.TestCheckTypeSetElemNestedAttrs("microsoft365_graph_beta_device_management_windows_update_ring.all_assignment_types", "assignments.*", map[string]string{"type": "allLicensedUsersAssignmentTarget"}),
					resource.TestCheckTypeSetElemNestedAttrs("microsoft365_graph_beta_device_management_windows_update_ring.all_assignment_types", "assignments.*", map[string]string{"type": "allDevicesAssignmentTarget"}),
					resource.TestCheckTypeSetElemNestedAttrs("microsoft365_graph_beta_device_management_windows_update_ring.all_assignment_types", "assignments.*", map[string]string{"type": "exclusionGroupAssignmentTarget"}),
				),
			},
		},
	})
}

// testConfigGroupAssignments returns the group assignments configuration for testing
func testConfigGroupAssignments() string {
	content, err := os.ReadFile(filepath.Join("tests", "terraform", "unit", "resource_with_group_assignments.tf"))
	if err != nil {
		return ""
	}
	return string(content)
}

// testConfigAllUsersAssignment returns the all users assignment configuration for testing
func testConfigAllUsersAssignment() string {
	content, err := os.ReadFile(filepath.Join("tests", "terraform", "unit", "resource_with_all_users_assignment.tf"))
	if err != nil {
		return ""
	}
	return string(content)
}

// testConfigAllDevicesAssignment returns the all devices assignment configuration for testing
func testConfigAllDevicesAssignment() string {
	content, err := os.ReadFile(filepath.Join("tests", "terraform", "unit", "resource_with_all_devices_assignment.tf"))
	if err != nil {
		return ""
	}
	return string(content)
}

// testConfigExclusionAssignment returns the exclusion assignment configuration for testing
func testConfigExclusionAssignment() string {
	content, err := os.ReadFile(filepath.Join("tests", "terraform", "unit", "resource_with_exclusion_assignment.tf"))
	if err != nil {
		return ""
	}
	return string(content)
}

// testConfigAllAssignmentTypes returns the all assignment types configuration for testing
func testConfigAllAssignmentTypes() string {
	content, err := os.ReadFile(filepath.Join("tests", "terraform", "unit", "resource_with_all_assignment_types.tf"))
	if err != nil {
		return ""
	}
	return string(content)
}
