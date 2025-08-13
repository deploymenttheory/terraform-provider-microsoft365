package graphBetaWindowsQualityUpdatePolicy_test

import (
	"os"
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	qualityUpdateMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/windows_quality_update_policy/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

func TestMain(m *testing.M) {
	exitCode := m.Run()
	os.Exit(exitCode)
}

// setupMockEnvironment sets up the mock environment using centralized mocks
func setupMockEnvironment() (*mocks.Mocks, *qualityUpdateMocks.WindowsQualityUpdatePolicyMock) {
	// Activate httpmock
	httpmock.Activate()

	// Create a new Mocks instance and register authentication mocks
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	// Register local mocks directly
	qualityUpdateMock := &qualityUpdateMocks.WindowsQualityUpdatePolicyMock{}
	qualityUpdateMock.RegisterMocks()

	return mockClient, qualityUpdateMock
}

// setupErrorMockEnvironment sets up the mock environment for error testing
func setupErrorMockEnvironment() (*mocks.Mocks, *qualityUpdateMocks.WindowsQualityUpdatePolicyMock) {
	// Activate httpmock
	httpmock.Activate()

	// Create a new Mocks instance and register authentication mocks
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	// Register error mocks
	qualityUpdateMock := &qualityUpdateMocks.WindowsQualityUpdatePolicyMock{}
	qualityUpdateMock.RegisterErrorMocks()

	return mockClient, qualityUpdateMock
}

// testCheckExists is a basic check to ensure the resource exists in the state
func testCheckExists(resourceName string) resource.TestCheckFunc {
	return resource.TestCheckResourceAttrSet(resourceName, "id")
}

// TestWindowsQualityUpdatePolicyResource_Schema validates the resource schema
func TestWindowsQualityUpdatePolicyResource_Schema(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, qualityUpdateMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer qualityUpdateMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_quality_update_policy.minimal", "display_name", "Test Minimal Windows Quality Update Policy - Unique"),
					resource.TestMatchResourceAttr("microsoft365_graph_beta_device_management_windows_quality_update_policy.minimal", "id", regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_quality_update_policy.minimal", "role_scope_tag_ids.#", "1"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_windows_quality_update_policy.minimal", "role_scope_tag_ids.*", "0"),
				),
			},
		},
	})
}

// TestWindowsQualityUpdatePolicyResource_Minimal tests basic CRUD operations
func TestWindowsQualityUpdatePolicyResource_Minimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, qualityUpdateMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer qualityUpdateMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_windows_quality_update_policy.minimal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_quality_update_policy.minimal", "display_name", "Test Minimal Windows Quality Update Policy - Unique"),
				),
			},
			{
				ResourceName:      "microsoft365_graph_beta_device_management_windows_quality_update_policy.minimal",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testConfigMaximal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_windows_quality_update_policy.maximal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_quality_update_policy.maximal", "display_name", "Test Maximal Windows Quality Update Policy - Unique"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_quality_update_policy.maximal", "description", "Maximal Windows Quality Update Policy for testing with all features"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_quality_update_policy.maximal", "hotpatch_enabled", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_quality_update_policy.maximal", "role_scope_tag_ids.#", "2"),
				),
			},
		},
	})
}

// TestWindowsQualityUpdatePolicyResource_RequiredFields tests required field validation
func TestWindowsQualityUpdatePolicyResource_RequiredFields(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, qualityUpdateMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer qualityUpdateMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "microsoft365_graph_beta_device_management_windows_quality_update_policy" "test" {
  # Missing display_name
}
`,
				ExpectError: regexp.MustCompile(`The argument "display_name" is required`),
			},
		},
	})
}

// TestWindowsQualityUpdatePolicyResource_GroupAssignments tests group assignment functionality
func TestWindowsQualityUpdatePolicyResource_GroupAssignments(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, qualityUpdateMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer qualityUpdateMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigGroupAssignments(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_windows_quality_update_policy.group_assignments"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_quality_update_policy.group_assignments", "display_name", "Test Group Assignments Windows Quality Update Policy - Unique"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_quality_update_policy.group_assignments", "assignments.#", "2"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_quality_update_policy.group_assignments", "assignments.0.type", "groupAssignmentTarget"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_quality_update_policy.group_assignments", "assignments.1.type", "groupAssignmentTarget"),
				),
			},
		},
	})
}

// TestWindowsQualityUpdatePolicyResource_ExclusionAssignment tests exclusion group assignment functionality
func TestWindowsQualityUpdatePolicyResource_ExclusionAssignment(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, qualityUpdateMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer qualityUpdateMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigExclusionAssignment(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_windows_quality_update_policy.exclusion_assignment"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_quality_update_policy.exclusion_assignment", "display_name", "Test Exclusion Assignment Windows Quality Update Policy - Unique"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_quality_update_policy.exclusion_assignment", "assignments.#", "1"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_quality_update_policy.exclusion_assignment", "assignments.0.type", "exclusionGroupAssignmentTarget"),
				),
			},
		},
	})
}

// testConfigMinimal returns the minimal configuration for testing
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

// testConfigExclusionAssignment returns the exclusion assignment configuration for testing
func testConfigExclusionAssignment() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_with_exclusion_assignment.tf")
	if err != nil {
		panic("failed to load exclusion assignment config: " + err.Error())
	}
	return unitTestConfig
}
