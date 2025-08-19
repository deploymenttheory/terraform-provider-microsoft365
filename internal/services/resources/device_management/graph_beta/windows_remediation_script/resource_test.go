package graphBetaWindowsRemediationScript_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	windowsRemediationScriptMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/windows_remediation_script/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

// setupMockEnvironment sets up the mock environment using centralized mocks
func setupMockEnvironment() (*mocks.Mocks, *windowsRemediationScriptMocks.WindowsRemediationScriptMock) {
	// Activate httpmock
	httpmock.Activate()

	// Create a new Mocks instance and register authentication mocks
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	// Register local mocks directly
	windowsRemediationScriptMock := &windowsRemediationScriptMocks.WindowsRemediationScriptMock{}
	windowsRemediationScriptMock.RegisterMocks()

	return mockClient, windowsRemediationScriptMock
}

// setupErrorMockEnvironment sets up the mock environment for error testing
func setupErrorMockEnvironment() (*mocks.Mocks, *windowsRemediationScriptMocks.WindowsRemediationScriptMock) {
	// Activate httpmock
	httpmock.Activate()

	// Create a new Mocks instance and register authentication mocks
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	// Register error mocks
	windowsRemediationScriptMock := &windowsRemediationScriptMocks.WindowsRemediationScriptMock{}
	windowsRemediationScriptMock.RegisterErrorMocks()

	return mockClient, windowsRemediationScriptMock
}

// testCheckExists is a basic check to ensure the resource exists in the state
func testCheckExists(resourceName string) resource.TestCheckFunc {
	return resource.TestCheckResourceAttrSet(resourceName, "id")
}

// TestWindowsRemediationScriptResource_Schema validates the resource schema
func TestWindowsRemediationScriptResource_Schema(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, windowsRemediationScriptMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer windowsRemediationScriptMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					// Check required attributes
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_remediation_script.minimal", "display_name", "Test Minimal Windows Remediation Script - Unique"),

					// Check computed attributes are set
					resource.TestMatchResourceAttr("microsoft365_graph_beta_device_management_windows_remediation_script.minimal", "id", regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_remediation_script.minimal", "role_scope_tag_ids.#", "1"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_windows_remediation_script.minimal", "role_scope_tag_ids.*", "0"),
				),
			},
		},
	})
}

// TestWindowsRemediationScriptResource_Minimal tests basic CRUD operations
func TestWindowsRemediationScriptResource_Minimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, windowsRemediationScriptMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer windowsRemediationScriptMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_windows_remediation_script.minimal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_remediation_script.minimal", "display_name", "Test Minimal Windows Remediation Script - Unique"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "microsoft365_graph_beta_device_management_windows_remediation_script.minimal",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testConfigMaximal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_windows_remediation_script.maximal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_remediation_script.maximal", "display_name", "Test Maximal Windows Remediation Script - Unique"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_remediation_script.maximal", "description", "Maximal Windows remediation script for testing with all features"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_remediation_script.maximal", "role_scope_tag_ids.#", "2"),
				),
			},
		},
	})
}

// TestWindowsRemediationScriptResource_UpdateInPlace tests in-place updates
func TestWindowsRemediationScriptResource_UpdateInPlace(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, windowsRemediationScriptMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer windowsRemediationScriptMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_windows_remediation_script.minimal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_remediation_script.minimal", "display_name", "Test Minimal Windows Remediation Script - Unique"),
				),
			},
			{
				Config: testConfigMaximal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_windows_remediation_script.maximal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_remediation_script.maximal", "display_name", "Test Maximal Windows Remediation Script - Unique"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_remediation_script.maximal", "description", "Maximal Windows remediation script for testing with all features"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_remediation_script.maximal", "role_scope_tag_ids.#", "2"),
				),
			},
		},
	})
}

// TestWindowsRemediationScriptResource_RequiredFields tests required field validation
func TestWindowsRemediationScriptResource_RequiredFields(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, windowsRemediationScriptMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer windowsRemediationScriptMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "microsoft365_graph_beta_device_management_windows_remediation_script" "test" {
  # Missing display_name
  publisher = "Test Publisher"
  run_as_account = "system"
}
`,
				ExpectError: regexp.MustCompile(`The argument "display_name" is required`),
			},
		},
	})
}

// TestWindowsRemediationScriptResource_ErrorHandling tests error scenarios
func TestWindowsRemediationScriptResource_ErrorHandling(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, windowsRemediationScriptMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer windowsRemediationScriptMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "microsoft365_graph_beta_device_management_windows_remediation_script" "test" {
  display_name = "Test Windows Remediation Script"
  publisher = "Test Publisher"
  run_as_account = "system"
  detection_script_content = "Write-Host 'Test'"
  remediation_script_content = "Write-Host 'Test'"
}
`,
				ExpectError: regexp.MustCompile(`Invalid Windows remediation script data|BadRequest`),
			},
		},
	})
}

// TestWindowsRemediationScriptResource_GroupAssignments tests group assignment functionality
func TestWindowsRemediationScriptResource_GroupAssignments(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, windowsRemediationScriptMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer windowsRemediationScriptMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigGroupAssignments(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_windows_remediation_script.group_assignments"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_remediation_script.group_assignments", "display_name", "Test Group Assignments Windows Remediation Script - Unique"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_remediation_script.group_assignments", "assignments.#", "2"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_remediation_script.group_assignments", "assignments.0.type", "groupAssignmentTarget"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_remediation_script.group_assignments", "assignments.1.type", "groupAssignmentTarget"),
				),
			},
		},
	})
}

// TestWindowsRemediationScriptResource_AllUsersAssignment tests all licensed users assignment functionality
func TestWindowsRemediationScriptResource_AllUsersAssignment(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, windowsRemediationScriptMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer windowsRemediationScriptMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigAllUsersAssignment(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_windows_remediation_script.all_users_assignment"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_remediation_script.all_users_assignment", "display_name", "Test All Users Assignment Windows Remediation Script - Unique"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_remediation_script.all_users_assignment", "assignments.#", "1"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_remediation_script.all_users_assignment", "assignments.0.type", "allLicensedUsersAssignmentTarget"),
				),
			},
		},
	})
}

// TestWindowsRemediationScriptResource_AllDevicesAssignment tests all devices assignment functionality
func TestWindowsRemediationScriptResource_AllDevicesAssignment(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, windowsRemediationScriptMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer windowsRemediationScriptMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigAllDevicesAssignment(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_windows_remediation_script.all_devices_assignment"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_remediation_script.all_devices_assignment", "display_name", "Test All Devices Assignment Windows Remediation Script - Unique"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_remediation_script.all_devices_assignment", "assignments.#", "1"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_remediation_script.all_devices_assignment", "assignments.0.type", "allDevicesAssignmentTarget"),
				),
			},
		},
	})
}

// TestWindowsRemediationScriptResource_ExclusionAssignment tests exclusion group assignment functionality
func TestWindowsRemediationScriptResource_ExclusionAssignment(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, windowsRemediationScriptMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer windowsRemediationScriptMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigExclusionAssignment(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_windows_remediation_script.exclusion_assignment"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_remediation_script.exclusion_assignment", "display_name", "Test Exclusion Assignment Windows Remediation Script - Unique"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_remediation_script.exclusion_assignment", "assignments.#", "1"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_remediation_script.exclusion_assignment", "assignments.0.type", "exclusionGroupAssignmentTarget"),
				),
			},
		},
	})
}

// TestWindowsRemediationScriptResource_AllAssignmentTypes tests all assignment types together
func TestWindowsRemediationScriptResource_AllAssignmentTypes(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, windowsRemediationScriptMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer windowsRemediationScriptMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigAllAssignmentTypes(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_windows_remediation_script.all_assignment_types"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_remediation_script.all_assignment_types", "display_name", "Test All Assignment Types Windows Remediation Script - Unique"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_remediation_script.all_assignment_types", "assignments.#", "5"),
					// Verify all assignment types are present
					resource.TestCheckTypeSetElemNestedAttrs("microsoft365_graph_beta_device_management_windows_remediation_script.all_assignment_types", "assignments.*", map[string]string{"type": "groupAssignmentTarget"}),
					resource.TestCheckTypeSetElemNestedAttrs("microsoft365_graph_beta_device_management_windows_remediation_script.all_assignment_types", "assignments.*", map[string]string{"type": "allLicensedUsersAssignmentTarget"}),
					resource.TestCheckTypeSetElemNestedAttrs("microsoft365_graph_beta_device_management_windows_remediation_script.all_assignment_types", "assignments.*", map[string]string{"type": "allDevicesAssignmentTarget"}),
					resource.TestCheckTypeSetElemNestedAttrs("microsoft365_graph_beta_device_management_windows_remediation_script.all_assignment_types", "assignments.*", map[string]string{"type": "exclusionGroupAssignmentTarget"}),
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

// testConfigMaximal returns the maximal configuration for testing
func testConfigMaximal() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_maximal.tf")
	if err != nil {
		panic("failed to load maximal config: " + err.Error())
	}
	return unitTestConfig
}

// testConfigGroupAssignments returns the group assignments configuration for testing
func testConfigGroupAssignments() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_with_group_assignments.tf")
	if err != nil {
		panic("failed to load group assignments config: " + err.Error())
	}
	return unitTestConfig
}

// testConfigAllUsersAssignment returns the all users assignment configuration for testing
func testConfigAllUsersAssignment() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_with_all_users_assignment.tf")
	if err != nil {
		panic("failed to load all users assignment config: " + err.Error())
	}
	return unitTestConfig
}

// testConfigAllDevicesAssignment returns the all devices assignment configuration for testing
func testConfigAllDevicesAssignment() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_with_all_devices_assignment.tf")
	if err != nil {
		panic("failed to load all devices assignment config: " + err.Error())
	}
	return unitTestConfig
}

// testConfigExclusionAssignment returns the exclusion assignment configuration for testing
func testConfigExclusionAssignment() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_with_exclusion_assignment.tf")
	if err != nil {
		panic("failed to load exclusion assignment config: " + err.Error())
	}
	return unitTestConfig
}

// testConfigAllAssignmentTypes returns the all assignment types configuration for testing
func testConfigAllAssignmentTypes() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_with_all_assignment_types.tf")
	if err != nil {
		panic("failed to load all assignment types config: " + err.Error())
	}
	return unitTestConfig
}
