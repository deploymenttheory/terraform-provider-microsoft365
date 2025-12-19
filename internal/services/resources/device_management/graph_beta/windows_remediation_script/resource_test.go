package graphBetaWindowsRemediationScript_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaWindowsRemediationScript "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/windows_remediation_script"
	windowsRemediationScriptMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/windows_remediation_script/mocks"
	groupMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/groups/graph_beta/group/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

// setupMockEnvironment sets up the mock environment using centralized mocks
func setupMockEnvironment() (*mocks.Mocks, *windowsRemediationScriptMocks.WindowsRemediationScriptMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	// Register group mocks for tests that create groups
	groupMock := &groupMocks.GroupMock{}
	groupMock.RegisterMocks()

	windowsRemediationScriptMock := &windowsRemediationScriptMocks.WindowsRemediationScriptMock{}
	windowsRemediationScriptMock.RegisterMocks()
	return mockClient, windowsRemediationScriptMock
}

// setupErrorMockEnvironment sets up the mock environment for error testing
func setupErrorMockEnvironment() (*mocks.Mocks, *windowsRemediationScriptMocks.WindowsRemediationScriptMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	// Register group mocks for tests that create groups
	groupMock := &groupMocks.GroupMock{}
	groupMock.RegisterMocks()

	windowsRemediationScriptMock := &windowsRemediationScriptMocks.WindowsRemediationScriptMock{}
	windowsRemediationScriptMock.RegisterErrorMocks()
	return mockClient, windowsRemediationScriptMock
}

// Helper function to load test configs from unit directory
func loadUnitTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/unit/" + filename)
	if err != nil {
		panic("failed to load unit test config " + filename + ": " + err.Error())
	}
	return config
}

// Test 001: Scenario 1 - Minimal configuration without assignments
func TestWindowsRemediationScriptResource_001_Scenario_Minimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, windowsRemediationScriptMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer windowsRemediationScriptMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("scenario_001_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaWindowsRemediationScript.ResourceName+".test_001").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaWindowsRemediationScript.ResourceName+".test_001").Key("display_name").HasValue("unit-test-windows-remediation-script-001-minimal"),
					check.That(graphBetaWindowsRemediationScript.ResourceName+".test_001").Key("description").HasValue("Scenario 1: Minimal configuration without assignments"),
					check.That(graphBetaWindowsRemediationScript.ResourceName+".test_001").Key("publisher").HasValue("Terraform Provider Test"),
					check.That(graphBetaWindowsRemediationScript.ResourceName+".test_001").Key("run_as_account").HasValue("system"),
					check.That(graphBetaWindowsRemediationScript.ResourceName+".test_001").Key("run_as_32_bit").HasValue("false"),
					check.That(graphBetaWindowsRemediationScript.ResourceName+".test_001").Key("enforce_signature_check").HasValue("false"),
					check.That(graphBetaWindowsRemediationScript.ResourceName+".test_001").Key("detection_script_content").Exists(),
					check.That(graphBetaWindowsRemediationScript.ResourceName+".test_001").Key("remediation_script_content").Exists(),
					check.That(graphBetaWindowsRemediationScript.ResourceName+".test_001").Key("role_scope_tag_ids.#").HasValue("1"),
					check.That(graphBetaWindowsRemediationScript.ResourceName+".test_001").Key("version").Exists(),
					check.That(graphBetaWindowsRemediationScript.ResourceName+".test_001").Key("is_global_script").HasValue("false"),
					check.That(graphBetaWindowsRemediationScript.ResourceName+".test_001").Key("device_health_script_type").HasValue("deviceHealthScript"),
				),
			},
			{
				ResourceName:      graphBetaWindowsRemediationScript.ResourceName + ".test_001",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Test 002: Scenario 2 - Maximal configuration without assignments
func TestWindowsRemediationScriptResource_002_Scenario_Maximal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, windowsRemediationScriptMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer windowsRemediationScriptMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("scenario_002_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaWindowsRemediationScript.ResourceName+".test_002").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaWindowsRemediationScript.ResourceName+".test_002").Key("display_name").HasValue("unit-test-windows-remediation-script-002-maximal"),
					check.That(graphBetaWindowsRemediationScript.ResourceName+".test_002").Key("description").HasValue("Scenario 2: Maximal configuration without assignments"),
					check.That(graphBetaWindowsRemediationScript.ResourceName+".test_002").Key("publisher").HasValue("Terraform Provider Test Suite"),
					check.That(graphBetaWindowsRemediationScript.ResourceName+".test_002").Key("run_as_account").HasValue("user"),
					check.That(graphBetaWindowsRemediationScript.ResourceName+".test_002").Key("run_as_32_bit").HasValue("true"),
					check.That(graphBetaWindowsRemediationScript.ResourceName+".test_002").Key("enforce_signature_check").HasValue("true"),
					check.That(graphBetaWindowsRemediationScript.ResourceName+".test_002").Key("detection_script_content").Exists(),
					check.That(graphBetaWindowsRemediationScript.ResourceName+".test_002").Key("remediation_script_content").Exists(),
					check.That(graphBetaWindowsRemediationScript.ResourceName+".test_002").Key("role_scope_tag_ids.#").HasValue("2"),
					check.That(graphBetaWindowsRemediationScript.ResourceName+".test_002").Key("version").Exists(),
					check.That(graphBetaWindowsRemediationScript.ResourceName+".test_002").Key("is_global_script").HasValue("false"),
					check.That(graphBetaWindowsRemediationScript.ResourceName+".test_002").Key("device_health_script_type").HasValue("deviceHealthScript"),
				),
			},
			{
				ResourceName:      graphBetaWindowsRemediationScript.ResourceName + ".test_002",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Test 003: Scenario 3 - Lifecycle from minimal to maximal
func TestWindowsRemediationScriptResource_003_Lifecycle_MinimalToMaximal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, windowsRemediationScriptMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer windowsRemediationScriptMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("lifecycle_step_1_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaWindowsRemediationScript.ResourceName+".test_003").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaWindowsRemediationScript.ResourceName+".test_003").Key("display_name").HasValue("unit-test-windows-remediation-script-003-lifecycle"),
					check.That(graphBetaWindowsRemediationScript.ResourceName+".test_003").Key("run_as_account").HasValue("system"),
					check.That(graphBetaWindowsRemediationScript.ResourceName+".test_003").Key("run_as_32_bit").HasValue("false"),
					check.That(graphBetaWindowsRemediationScript.ResourceName+".test_003").Key("enforce_signature_check").HasValue("false"),
					check.That(graphBetaWindowsRemediationScript.ResourceName+".test_003").Key("detection_script_content").Exists(),
					check.That(graphBetaWindowsRemediationScript.ResourceName+".test_003").Key("remediation_script_content").Exists(),
				),
			},
			{
				Config: loadUnitTestTerraform("lifecycle_step_2_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaWindowsRemediationScript.ResourceName+".test_003").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaWindowsRemediationScript.ResourceName+".test_003").Key("display_name").HasValue("unit-test-windows-remediation-script-003-lifecycle"),
					check.That(graphBetaWindowsRemediationScript.ResourceName+".test_003").Key("run_as_account").HasValue("user"),
					check.That(graphBetaWindowsRemediationScript.ResourceName+".test_003").Key("run_as_32_bit").HasValue("true"),
					check.That(graphBetaWindowsRemediationScript.ResourceName+".test_003").Key("enforce_signature_check").HasValue("true"),
					check.That(graphBetaWindowsRemediationScript.ResourceName+".test_003").Key("role_scope_tag_ids.#").HasValue("2"),
					check.That(graphBetaWindowsRemediationScript.ResourceName+".test_003").Key("detection_script_content").Exists(),
					check.That(graphBetaWindowsRemediationScript.ResourceName+".test_003").Key("remediation_script_content").Exists(),
				),
			},
			{
				ResourceName:      graphBetaWindowsRemediationScript.ResourceName + ".test_003",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Test 004: Scenario 4 - Lifecycle from maximal to minimal
func TestWindowsRemediationScriptResource_004_Lifecycle_MaximalToMinimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, windowsRemediationScriptMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer windowsRemediationScriptMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("lifecycle_step_1_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaWindowsRemediationScript.ResourceName+".test_004").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaWindowsRemediationScript.ResourceName+".test_004").Key("display_name").HasValue("unit-test-windows-remediation-script-004-downgrade"),
					check.That(graphBetaWindowsRemediationScript.ResourceName+".test_004").Key("run_as_account").HasValue("user"),
					check.That(graphBetaWindowsRemediationScript.ResourceName+".test_004").Key("run_as_32_bit").HasValue("true"),
					check.That(graphBetaWindowsRemediationScript.ResourceName+".test_004").Key("enforce_signature_check").HasValue("true"),
					check.That(graphBetaWindowsRemediationScript.ResourceName+".test_004").Key("role_scope_tag_ids.#").HasValue("2"),
				),
			},
			{
				Config: loadUnitTestTerraform("lifecycle_step_2_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaWindowsRemediationScript.ResourceName+".test_004").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaWindowsRemediationScript.ResourceName+".test_004").Key("display_name").HasValue("unit-test-windows-remediation-script-004-downgrade"),
					check.That(graphBetaWindowsRemediationScript.ResourceName+".test_004").Key("run_as_account").HasValue("system"),
					check.That(graphBetaWindowsRemediationScript.ResourceName+".test_004").Key("run_as_32_bit").HasValue("false"),
					check.That(graphBetaWindowsRemediationScript.ResourceName+".test_004").Key("enforce_signature_check").HasValue("false"),
					check.That(graphBetaWindowsRemediationScript.ResourceName+".test_004").Key("role_scope_tag_ids.#").HasValue("1"),
				),
			},
			{
				ResourceName:      graphBetaWindowsRemediationScript.ResourceName + ".test_004",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Test 005: Scenario 5 - Minimal assignments
func TestWindowsRemediationScriptResource_005_AssignmentsMinimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, windowsRemediationScriptMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer windowsRemediationScriptMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("assignments_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaWindowsRemediationScript.ResourceName+".test_005").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaWindowsRemediationScript.ResourceName+".test_005").Key("display_name").HasValue("unit-test-windows-remediation-script-005-assignments-minimal"),
					check.That(graphBetaWindowsRemediationScript.ResourceName+".test_005").Key("description").HasValue("Scenario 5: Minimal assignments"),
					check.That(graphBetaWindowsRemediationScript.ResourceName+".test_005").Key("assignments.#").HasValue("1"),
					check.That(graphBetaWindowsRemediationScript.ResourceName+".test_005").Key("detection_script_content").Exists(),
					check.That(graphBetaWindowsRemediationScript.ResourceName+".test_005").Key("remediation_script_content").Exists(),
				),
			},
			{
				ResourceName:      graphBetaWindowsRemediationScript.ResourceName + ".test_005",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Test 006: Scenario 6 - Maximal assignments
func TestWindowsRemediationScriptResource_006_AssignmentsMaximal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, windowsRemediationScriptMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer windowsRemediationScriptMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("assignments_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaWindowsRemediationScript.ResourceName+".test_006").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaWindowsRemediationScript.ResourceName+".test_006").Key("display_name").HasValue("unit-test-windows-remediation-script-006-assignments-maximal"),
					check.That(graphBetaWindowsRemediationScript.ResourceName+".test_006").Key("description").HasValue("Scenario 6: Maximal assignments"),
					check.That(graphBetaWindowsRemediationScript.ResourceName+".test_006").Key("assignments.#").HasValue("5"),
					check.That(graphBetaWindowsRemediationScript.ResourceName+".test_006").Key("detection_script_content").Exists(),
					check.That(graphBetaWindowsRemediationScript.ResourceName+".test_006").Key("remediation_script_content").Exists(),
					// Verify assignment types are present
					resource.TestCheckTypeSetElemNestedAttrs(graphBetaWindowsRemediationScript.ResourceName+".test_006", "assignments.*", map[string]string{
						"type":        "groupAssignmentTarget",
						"group_id":    "44444444-4444-4444-4444-444444444444",
						"filter_type": "include",
						"filter_id":   "55555555-5555-5555-5555-555555555555",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(graphBetaWindowsRemediationScript.ResourceName+".test_006", "assignments.*", map[string]string{
						"type":        "groupAssignmentTarget",
						"group_id":    "33333333-3333-3333-3333-333333333333",
						"filter_type": "exclude",
						"filter_id":   "66666666-6666-6666-6666-666666666666",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(graphBetaWindowsRemediationScript.ResourceName+".test_006", "assignments.*", map[string]string{
						"type": "allLicensedUsersAssignmentTarget",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(graphBetaWindowsRemediationScript.ResourceName+".test_006", "assignments.*", map[string]string{
						"type": "allDevicesAssignmentTarget",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(graphBetaWindowsRemediationScript.ResourceName+".test_006", "assignments.*", map[string]string{
						"type":     "exclusionGroupAssignmentTarget",
						"group_id": "77777777-7777-7777-7777-777777777777",
					}),
				),
			},
			{
				ResourceName:      graphBetaWindowsRemediationScript.ResourceName + ".test_006",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Test 007: Scenario 7 - Assignments lifecycle minimal to maximal
func TestWindowsRemediationScriptResource_007_AssignmentsLifecycle_MinimalToMaximal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, windowsRemediationScriptMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer windowsRemediationScriptMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("assignments_lifecycle_step_1_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaWindowsRemediationScript.ResourceName+".test_007").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaWindowsRemediationScript.ResourceName+".test_007").Key("display_name").HasValue("unit-test-windows-remediation-script-007-assignments-lifecycle"),
					check.That(graphBetaWindowsRemediationScript.ResourceName+".test_007").Key("assignments.#").HasValue("1"),
					resource.TestCheckTypeSetElemNestedAttrs(graphBetaWindowsRemediationScript.ResourceName+".test_007", "assignments.*", map[string]string{
						"type": "allDevicesAssignmentTarget",
					}),
				),
			},
			{
				Config: loadUnitTestTerraform("assignments_lifecycle_step_2_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaWindowsRemediationScript.ResourceName+".test_007").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaWindowsRemediationScript.ResourceName+".test_007").Key("display_name").HasValue("unit-test-windows-remediation-script-007-assignments-lifecycle"),
					check.That(graphBetaWindowsRemediationScript.ResourceName+".test_007").Key("assignments.#").HasValue("5"),
					resource.TestCheckTypeSetElemNestedAttrs(graphBetaWindowsRemediationScript.ResourceName+".test_007", "assignments.*", map[string]string{
						"type": "allLicensedUsersAssignmentTarget",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(graphBetaWindowsRemediationScript.ResourceName+".test_007", "assignments.*", map[string]string{
						"type": "allDevicesAssignmentTarget",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(graphBetaWindowsRemediationScript.ResourceName+".test_007", "assignments.*", map[string]string{
						"type":     "exclusionGroupAssignmentTarget",
						"group_id": "77777777-7777-7777-7777-777777777777",
					}),
				),
			},
			{
				ResourceName:      graphBetaWindowsRemediationScript.ResourceName + ".test_007",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Test 008: Scenario 8 - Assignments lifecycle maximal to minimal
func TestWindowsRemediationScriptResource_008_AssignmentsLifecycle_MaximalToMinimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, windowsRemediationScriptMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer windowsRemediationScriptMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("assignments_downgrade_step_1_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaWindowsRemediationScript.ResourceName+".test_008").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaWindowsRemediationScript.ResourceName+".test_008").Key("display_name").HasValue("unit-test-windows-remediation-script-008-assignments-downgrade"),
					check.That(graphBetaWindowsRemediationScript.ResourceName+".test_008").Key("assignments.#").HasValue("5"),
					resource.TestCheckTypeSetElemNestedAttrs(graphBetaWindowsRemediationScript.ResourceName+".test_008", "assignments.*", map[string]string{
						"type": "allLicensedUsersAssignmentTarget",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(graphBetaWindowsRemediationScript.ResourceName+".test_008", "assignments.*", map[string]string{
						"type": "allDevicesAssignmentTarget",
					}),
				),
			},
			{
				Config: loadUnitTestTerraform("assignments_downgrade_step_2_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaWindowsRemediationScript.ResourceName+".test_008").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaWindowsRemediationScript.ResourceName+".test_008").Key("display_name").HasValue("unit-test-windows-remediation-script-008-assignments-downgrade"),
					check.That(graphBetaWindowsRemediationScript.ResourceName+".test_008").Key("assignments.#").HasValue("1"),
					resource.TestCheckTypeSetElemNestedAttrs(graphBetaWindowsRemediationScript.ResourceName+".test_008", "assignments.*", map[string]string{
						"type": "allDevicesAssignmentTarget",
					}),
				),
			},
			{
				ResourceName:      graphBetaWindowsRemediationScript.ResourceName + ".test_008",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Test 009: Scenario 9 - Validation errors
func TestWindowsRemediationScriptResource_009_ValidationErrors(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, windowsRemediationScriptMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer windowsRemediationScriptMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("validation_missing_display_name.tf"),
				ExpectError: regexp.MustCompile(`The argument "display_name" is required`),
			},
			{
				Config:      loadUnitTestTerraform("validation_invalid_run_as_account.tf"),
				ExpectError: regexp.MustCompile(`Attribute run_as_account value must be one of`),
			},
		},
	})
}
