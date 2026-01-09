package graphBetaMacOSPlatformScript_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	platformScriptMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/macos_platform_script/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

// setupMockEnvironment initializes the mock environment for testing
func setupMockEnvironment() (*mocks.Mocks, *platformScriptMocks.MacOSPlatformScriptMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	macOSPlatformScriptMock := &platformScriptMocks.MacOSPlatformScriptMock{}
	macOSPlatformScriptMock.RegisterMocks()
	return mockClient, macOSPlatformScriptMock
}

// loadUnitTestTerraform loads a Terraform configuration file for unit testing
func loadUnitTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/unit/" + filename)
	if err != nil {
		panic("failed to load unit test config " + filename + ": " + err.Error())
	}
	return config
}

// ====================================================================================
// Scenario 01: Minimal Resource
// ====================================================================================

func TestMacOSPlatformScriptResource_Scenario01_Minimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, macOSPlatformScriptMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer macOSPlatformScriptMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_scenario_01_minimal_001.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".minimal").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".minimal").Key("display_name").HasValue("unit-test-minimal-macos-script"),
					check.That(resourceType+".minimal").Key("file_name").HasValue("minimal_test.sh"),
					check.That(resourceType+".minimal").Key("script_content").HasValue("#!/bin/bash\necho 'Min Test'\nexit 0"),
					check.That(resourceType+".minimal").Key("run_as_account").HasValue("system"),

					// Verify optional fields are not set
					check.That(resourceType+".minimal").Key("description").IsEmpty(),
					check.That(resourceType+".minimal").Key("role_scope_tag_ids.#").HasValue("1"),
					check.That(resourceType+".minimal").Key("role_scope_tag_ids.*").ContainsTypeSetElement("0"),

					// Verify no assignments
					check.That(resourceType+".minimal").Key("assignments.#").HasValue("0"),
				),
			},
		},
	})
}

// ====================================================================================
// Scenario 02: Maximal Resource
// ====================================================================================

func TestMacOSPlatformScriptResource_Scenario02_Maximal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, macOSPlatformScriptMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer macOSPlatformScriptMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_scenario_02_maximal_001.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".maximal").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".maximal").Key("display_name").HasValue("unit-test-maximal-macos-script"),
					check.That(resourceType+".maximal").Key("description").HasValue("Comprehensive macOS platform script with all features enabled for unit testing"),
					check.That(resourceType+".maximal").Key("file_name").HasValue("maximal_test.sh"),
					check.That(resourceType+".maximal").Key("run_as_account").HasValue("user"),
					check.That(resourceType+".maximal").Key("execution_frequency").HasValue("P1D"),
					check.That(resourceType+".maximal").Key("block_execution_notifications").HasValue("true"),
					check.That(resourceType+".maximal").Key("retry_count").HasValue("3"),

					// Verify role_scope_tag_ids
					check.That(resourceType+".maximal").Key("role_scope_tag_ids.#").HasValue("3"),
					check.That(resourceType+".maximal").Key("role_scope_tag_ids.*").ContainsTypeSetElement("0"),
					check.That(resourceType+".maximal").Key("role_scope_tag_ids.*").ContainsTypeSetElement("1"),
					check.That(resourceType+".maximal").Key("role_scope_tag_ids.*").ContainsTypeSetElement("2"),

					// Verify no assignments (base resource test)
					check.That(resourceType+".maximal").Key("assignments.#").HasValue("0"),
				),
			},
		},
	})
}

// ====================================================================================
// Scenario 03: Minimal to Maximal Update
// ====================================================================================

func TestMacOSPlatformScriptResource_Scenario03_MinimalToMaximal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, macOSPlatformScriptMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer macOSPlatformScriptMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Step 1: Create minimal configuration
			{
				Config: loadUnitTestTerraform("resource_scenario_03_minimal_to_maximal_step_01_001.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".update_test").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".update_test").Key("display_name").HasValue("unit-test-update-test-script"),
					check.That(resourceType+".update_test").Key("file_name").HasValue("update_test.sh"),
					check.That(resourceType+".update_test").Key("run_as_account").HasValue("system"),
					check.That(resourceType+".update_test").Key("description").IsEmpty(),
					check.That(resourceType+".update_test").Key("assignments.#").HasValue("0"),
				),
			},
			// Step 2: Update to maximal configuration
			{
				Config: loadUnitTestTerraform("resource_scenario_03_minimal_to_maximal_step_02_001.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".update_test").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".update_test").Key("display_name").HasValue("unit-test-update-test-script-updated"),
					check.That(resourceType+".update_test").Key("description").HasValue("Updated to maximal configuration"),
					check.That(resourceType+".update_test").Key("file_name").HasValue("update_test_maximal.sh"),
					check.That(resourceType+".update_test").Key("run_as_account").HasValue("user"),
					check.That(resourceType+".update_test").Key("execution_frequency").HasValue("PT12H"),
					check.That(resourceType+".update_test").Key("block_execution_notifications").HasValue("true"),
					check.That(resourceType+".update_test").Key("retry_count").HasValue("2"),
					check.That(resourceType+".update_test").Key("role_scope_tag_ids.#").HasValue("2"),
					check.That(resourceType+".update_test").Key("assignments.#").HasValue("0"),
				),
			},
		},
	})
}

// ====================================================================================
// Scenario 04: Maximal to Minimal Update
// ====================================================================================

func TestMacOSPlatformScriptResource_Scenario04_MaximalToMinimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, macOSPlatformScriptMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer macOSPlatformScriptMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Step 1: Create maximal configuration
			{
				Config: loadUnitTestTerraform("resource_scenario_04_maximal_to_minimal_step_01_001.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".downgrade_test").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".downgrade_test").Key("display_name").HasValue("unit-test-downgrade-test-script"),
					check.That(resourceType+".downgrade_test").Key("description").HasValue("Initial maximal configuration for downgrade testing"),
					check.That(resourceType+".downgrade_test").Key("run_as_account").HasValue("user"),
					check.That(resourceType+".downgrade_test").Key("execution_frequency").HasValue("P1D"),
					check.That(resourceType+".downgrade_test").Key("block_execution_notifications").HasValue("true"),
					check.That(resourceType+".downgrade_test").Key("retry_count").HasValue("5"),
					check.That(resourceType+".downgrade_test").Key("role_scope_tag_ids.#").HasValue("3"),
					check.That(resourceType+".downgrade_test").Key("assignments.#").HasValue("0"),
				),
			},
			// Step 2: Downgrade to minimal configuration
			{
				Config: loadUnitTestTerraform("resource_scenario_04_maximal_to_minimal_step_02_001.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".downgrade_test").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".downgrade_test").Key("display_name").HasValue("unit-test-downgrade-test-script-minimal"),
					check.That(resourceType+".downgrade_test").Key("file_name").HasValue("downgrade_test_minimal.sh"),
					check.That(resourceType+".downgrade_test").Key("run_as_account").HasValue("system"),
					check.That(resourceType+".downgrade_test").Key("script_content").HasValue("#!/bin/bash\necho 'Min Test'\nexit 0"),
				),
			},
		},
	})
}

// ====================================================================================
// Scenario 05: No Assignments to Minimal Assignment
// ====================================================================================

func TestMacOSPlatformScriptResource_Scenario05_NoAssignmentsToMinimalAssignment(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, macOSPlatformScriptMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer macOSPlatformScriptMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Step 1: Create resource with no assignments
			{
				Config: loadUnitTestTerraform("resource_scenario_05_no_assignments_to_minimal_assignments_step_01_001.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".add_minimal_assignment").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".add_minimal_assignment").Key("display_name").HasValue("unit-test-add-minimal-assignment"),
					check.That(resourceType+".add_minimal_assignment").Key("assignments.#").HasValue("0"),
				),
			},
			// Step 2: Add single assignment
			{
				Config: loadUnitTestTerraform("resource_scenario_05_no_assignments_to_minimal_assignments_step_02_001.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".add_minimal_assignment").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".add_minimal_assignment").Key("display_name").HasValue("unit-test-add-minimal-assignment"),
					check.That(resourceType+".add_minimal_assignment").Key("assignments.#").HasValue("1"),
					check.That(resourceType+".add_minimal_assignment").Key("assignments.0.type").HasValue("groupAssignmentTarget"),
					check.That(resourceType+".add_minimal_assignment").Key("assignments.0.group_id").HasValue("11111111-1111-1111-1111-111111111111"),
				),
			},
		},
	})
}

// ====================================================================================
// Scenario 06: No Assignments to Maximal Assignments
// ====================================================================================

func TestMacOSPlatformScriptResource_Scenario06_NoAssignmentsToMaximalAssignments(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, macOSPlatformScriptMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer macOSPlatformScriptMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Step 1: Create resource with no assignments
			{
				Config: loadUnitTestTerraform("resource_scenario_06_no_assignments_to_maximal_assignments_step_01_001.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".add_maximal_assignments").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".add_maximal_assignments").Key("display_name").HasValue("unit-test-add-maximal-assignments"),
					check.That(resourceType+".add_maximal_assignments").Key("assignments.#").HasValue("0"),
				),
			},
			// Step 2: Add all 4 assignment types
			{
				Config: loadUnitTestTerraform("resource_scenario_06_no_assignments_to_maximal_assignments_step_02_001.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".add_maximal_assignments").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".add_maximal_assignments").Key("display_name").HasValue("unit-test-add-maximal-assignments"),
					check.That(resourceType+".add_maximal_assignments").Key("assignments.#").HasValue("4"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceType+".add_maximal_assignments", "assignments.*", map[string]string{
						"type": "allDevicesAssignmentTarget",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceType+".add_maximal_assignments", "assignments.*", map[string]string{
						"type": "allLicensedUsersAssignmentTarget",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceType+".add_maximal_assignments", "assignments.*", map[string]string{
						"type":     "groupAssignmentTarget",
						"group_id": "22222222-2222-2222-2222-222222222222",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceType+".add_maximal_assignments", "assignments.*", map[string]string{
						"type":     "exclusionGroupAssignmentTarget",
						"group_id": "33333333-3333-3333-3333-333333333333",
					}),
				),
			},
		},
	})
}

// ====================================================================================
// Scenario 07: Minimal to Maximal Assignments
// ====================================================================================

func TestMacOSPlatformScriptResource_Scenario07_MinimalToMaximalAssignments(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, macOSPlatformScriptMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer macOSPlatformScriptMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Step 1: Create with single assignment
			{
				Config: loadUnitTestTerraform("resource_scenario_07_minimal_to_maximal_assignments_step_01_001.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".assignment_update").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".assignment_update").Key("display_name").HasValue("unit-test-assignment-update"),
					check.That(resourceType+".assignment_update").Key("assignments.#").HasValue("1"),
					check.That(resourceType+".assignment_update").Key("assignments.0.type").HasValue("groupAssignmentTarget"),
					check.That(resourceType+".assignment_update").Key("assignments.0.group_id").HasValue("55555555-5555-5555-5555-555555555555"),
				),
			},
			// Step 2: Update to all 4 assignment types
			{
				Config: loadUnitTestTerraform("resource_scenario_07_minimal_to_maximal_assignments_step_02_001.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".assignment_update").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".assignment_update").Key("display_name").HasValue("unit-test-assignment-update"),
					check.That(resourceType+".assignment_update").Key("assignments.#").HasValue("4"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceType+".assignment_update", "assignments.*", map[string]string{
						"type": "allDevicesAssignmentTarget",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceType+".assignment_update", "assignments.*", map[string]string{
						"type": "allLicensedUsersAssignmentTarget",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceType+".assignment_update", "assignments.*", map[string]string{
						"type":     "groupAssignmentTarget",
						"group_id": "66666666-6666-6666-6666-666666666666",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceType+".assignment_update", "assignments.*", map[string]string{
						"type":     "exclusionGroupAssignmentTarget",
						"group_id": "77777777-7777-7777-7777-777777777777",
					}),
				),
			},
		},
	})
}

// ====================================================================================
// Scenario 08: Maximal to Minimal Assignments
// ====================================================================================

func TestMacOSPlatformScriptResource_Scenario08_MaximalToMinimalAssignments(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, macOSPlatformScriptMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer macOSPlatformScriptMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Step 1: Create with all 4 assignment types
			{
				Config: loadUnitTestTerraform("resource_scenario_08_maximal_to_minimal_assignments_step_01_001.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".assignment_downgrade").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".assignment_downgrade").Key("display_name").HasValue("unit-test-assignment-downgrade"),
					check.That(resourceType+".assignment_downgrade").Key("assignments.#").HasValue("4"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceType+".assignment_downgrade", "assignments.*", map[string]string{
						"type": "allDevicesAssignmentTarget",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceType+".assignment_downgrade", "assignments.*", map[string]string{
						"type": "allLicensedUsersAssignmentTarget",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceType+".assignment_downgrade", "assignments.*", map[string]string{
						"type":     "groupAssignmentTarget",
						"group_id": "88888888-8888-8888-8888-888888888888",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceType+".assignment_downgrade", "assignments.*", map[string]string{
						"type":     "exclusionGroupAssignmentTarget",
						"group_id": "99999999-9999-9999-9999-999999999999",
					}),
				),
			},
			// Step 2: Downgrade to single assignment
			{
				Config: loadUnitTestTerraform("resource_scenario_08_maximal_to_minimal_assignments_step_02_001.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".assignment_downgrade").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".assignment_downgrade").Key("display_name").HasValue("unit-test-assignment-downgrade"),
					check.That(resourceType+".assignment_downgrade").Key("assignments.#").HasValue("1"),
					check.That(resourceType+".assignment_downgrade").Key("assignments.0.type").HasValue("groupAssignmentTarget"),
					check.That(resourceType+".assignment_downgrade").Key("assignments.0.group_id").HasValue("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"),
				),
			},
		},
	})
}

// ====================================================================================
// Scenario 09: Error Cases
// ====================================================================================

func TestMacOSPlatformScriptResource_Scenario09_Error_InvalidRunAs(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, macOSPlatformScriptMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer macOSPlatformScriptMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("resource_scenario_09_error_invalid_run_as_001.tf"),
				ExpectError: regexp.MustCompile(`Attribute run_as_account value must be one of`),
			},
		},
	})
}

func TestMacOSPlatformScriptResource_Scenario09_Error_InvalidDuration(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, macOSPlatformScriptMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer macOSPlatformScriptMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("resource_scenario_09_error_invalid_duration_002.tf"),
				ExpectError: regexp.MustCompile(`must be a valid ISO 8601 duration`),
			},
		},
	})
}

func TestMacOSPlatformScriptResource_Scenario09_Error_DescriptionLength(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, macOSPlatformScriptMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer macOSPlatformScriptMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("resource_scenario_09_error_description_length_003.tf"),
				ExpectError: regexp.MustCompile(`Attribute description string length must be at most 1500`),
			},
		},
	})
}
