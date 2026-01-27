package graphBetaWindowsQualityUpdatePolicy_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaWindowsQualityUpdatePolicy "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/windows_quality_update_policy"
	qualityUpdateMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/windows_quality_update_policy/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

// setupMockEnvironment sets up the mock environment using centralized mocks
func setupMockEnvironment() (*mocks.Mocks, *qualityUpdateMocks.WindowsQualityUpdatePolicyMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	qualityUpdateMock := &qualityUpdateMocks.WindowsQualityUpdatePolicyMock{}
	qualityUpdateMock.RegisterMocks()
	return mockClient, qualityUpdateMock
}

// setupErrorMockEnvironment sets up the mock environment for error testing
func setupErrorMockEnvironment() (*mocks.Mocks, *qualityUpdateMocks.WindowsQualityUpdatePolicyMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	qualityUpdateMock := &qualityUpdateMocks.WindowsQualityUpdatePolicyMock{}
	qualityUpdateMock.RegisterErrorMocks()
	return mockClient, qualityUpdateMock
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
func TestUnitResourceWindowsQualityUpdatePolicy_01_Scenario_Minimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, qualityUpdateMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer qualityUpdateMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("001_scenario_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaWindowsQualityUpdatePolicy.ResourceName+".test_001").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaWindowsQualityUpdatePolicy.ResourceName+".test_001").Key("display_name").HasValue("unit-test-windows-quality-update-policy-001-minimal"),
					check.That(graphBetaWindowsQualityUpdatePolicy.ResourceName+".test_001").Key("role_scope_tag_ids.#").HasValue("1"),
					check.That(graphBetaWindowsQualityUpdatePolicy.ResourceName+".test_001").Key("role_scope_tag_ids.0").HasValue("0"),
				),
			},
			{
				ResourceName:      graphBetaWindowsQualityUpdatePolicy.ResourceName + ".test_001",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Test 002: Scenario 2 - Maximal configuration without assignments
func TestUnitResourceWindowsQualityUpdatePolicy_02_Scenario_Maximal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, qualityUpdateMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer qualityUpdateMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("002_scenario_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaWindowsQualityUpdatePolicy.ResourceName+".test_002").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaWindowsQualityUpdatePolicy.ResourceName+".test_002").Key("display_name").HasValue("unit-test-windows-quality-update-policy-002-maximal"),
					check.That(graphBetaWindowsQualityUpdatePolicy.ResourceName+".test_002").Key("description").HasValue("Scenario 2: Maximal configuration without assignments"),
					check.That(graphBetaWindowsQualityUpdatePolicy.ResourceName+".test_002").Key("hotpatch_enabled").HasValue("true"),
					check.That(graphBetaWindowsQualityUpdatePolicy.ResourceName+".test_002").Key("role_scope_tag_ids.#").HasValue("2"),
				),
			},
			{
				ResourceName:      graphBetaWindowsQualityUpdatePolicy.ResourceName + ".test_002",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Test 003: Scenario 3 - Lifecycle from minimal to maximal
func TestUnitResourceWindowsQualityUpdatePolicy_03_Lifecycle_MinimalToMaximal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, qualityUpdateMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer qualityUpdateMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("003_lifecycle_minimal_to_maximal_step_1.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaWindowsQualityUpdatePolicy.ResourceName+".test_003").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaWindowsQualityUpdatePolicy.ResourceName+".test_003").Key("display_name").HasValue("unit-test-windows-quality-update-policy-003-lifecycle"),
					check.That(graphBetaWindowsQualityUpdatePolicy.ResourceName+".test_003").Key("role_scope_tag_ids.#").HasValue("1"),
				),
			},
			{
				Config: loadUnitTestTerraform("003_lifecycle_minimal_to_maximal_step_2.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaWindowsQualityUpdatePolicy.ResourceName+".test_003").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaWindowsQualityUpdatePolicy.ResourceName+".test_003").Key("display_name").HasValue("unit-test-windows-quality-update-policy-003-lifecycle"),
					check.That(graphBetaWindowsQualityUpdatePolicy.ResourceName+".test_003").Key("description").HasValue("Lifecycle Step 2: Updated to maximal configuration"),
					check.That(graphBetaWindowsQualityUpdatePolicy.ResourceName+".test_003").Key("hotpatch_enabled").HasValue("true"),
					check.That(graphBetaWindowsQualityUpdatePolicy.ResourceName+".test_003").Key("role_scope_tag_ids.#").HasValue("2"),
				),
			},
			{
				ResourceName:      graphBetaWindowsQualityUpdatePolicy.ResourceName + ".test_003",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Test 004: Scenario 4 - Lifecycle from maximal to minimal
func TestUnitResourceWindowsQualityUpdatePolicy_04_Lifecycle_MaximalToMinimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, qualityUpdateMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer qualityUpdateMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("004_lifecycle_maximal_to_minimal_step_1.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaWindowsQualityUpdatePolicy.ResourceName+".test_004").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaWindowsQualityUpdatePolicy.ResourceName+".test_004").Key("display_name").HasValue("unit-test-windows-quality-update-policy-004-lifecycle"),
					check.That(graphBetaWindowsQualityUpdatePolicy.ResourceName+".test_004").Key("description").HasValue("Lifecycle Step 1: Starting with maximal configuration"),
					check.That(graphBetaWindowsQualityUpdatePolicy.ResourceName+".test_004").Key("hotpatch_enabled").HasValue("true"),
					check.That(graphBetaWindowsQualityUpdatePolicy.ResourceName+".test_004").Key("role_scope_tag_ids.#").HasValue("2"),
				),
			},
			{
				Config: loadUnitTestTerraform("004_lifecycle_maximal_to_minimal_step_2.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaWindowsQualityUpdatePolicy.ResourceName+".test_004").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaWindowsQualityUpdatePolicy.ResourceName+".test_004").Key("display_name").HasValue("unit-test-windows-quality-update-policy-004-lifecycle"),
					check.That(graphBetaWindowsQualityUpdatePolicy.ResourceName+".test_004").Key("role_scope_tag_ids.#").HasValue("1"),
				),
			},
			{
				ResourceName:      graphBetaWindowsQualityUpdatePolicy.ResourceName + ".test_004",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Test 005: Scenario 5 - Minimal assignments
func TestUnitResourceWindowsQualityUpdatePolicy_05_AssignmentsMinimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, qualityUpdateMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer qualityUpdateMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("005_assignments_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaWindowsQualityUpdatePolicy.ResourceName+".test_005").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaWindowsQualityUpdatePolicy.ResourceName+".test_005").Key("display_name").HasValue("unit-test-windows-quality-update-policy-005-assignments-minimal"),
					check.That(graphBetaWindowsQualityUpdatePolicy.ResourceName+".test_005").Key("assignments.#").HasValue("1"),
					resource.TestCheckTypeSetElemNestedAttrs(graphBetaWindowsQualityUpdatePolicy.ResourceName+".test_005", "assignments.*", map[string]string{
						"type":     "groupAssignmentTarget",
						"group_id": "11111111-1111-1111-1111-111111111111",
					}),
				),
			},
			{
				ResourceName:      graphBetaWindowsQualityUpdatePolicy.ResourceName + ".test_005",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Test 006: Scenario 6 - Maximal assignments
func TestUnitResourceWindowsQualityUpdatePolicy_06_AssignmentsMaximal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, qualityUpdateMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer qualityUpdateMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("006_assignments_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaWindowsQualityUpdatePolicy.ResourceName+".test_006").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaWindowsQualityUpdatePolicy.ResourceName+".test_006").Key("display_name").HasValue("unit-test-windows-quality-update-policy-006-assignments-maximal"),
					check.That(graphBetaWindowsQualityUpdatePolicy.ResourceName+".test_006").Key("assignments.#").HasValue("3"),
					resource.TestCheckTypeSetElemNestedAttrs(graphBetaWindowsQualityUpdatePolicy.ResourceName+".test_006", "assignments.*", map[string]string{
						"type":     "groupAssignmentTarget",
						"group_id": "44444444-4444-4444-4444-444444444444",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(graphBetaWindowsQualityUpdatePolicy.ResourceName+".test_006", "assignments.*", map[string]string{
						"type":     "groupAssignmentTarget",
						"group_id": "33333333-3333-3333-3333-333333333333",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(graphBetaWindowsQualityUpdatePolicy.ResourceName+".test_006", "assignments.*", map[string]string{
						"type":     "exclusionGroupAssignmentTarget",
						"group_id": "77777777-7777-7777-7777-777777777777",
					}),
				),
			},
			{
				ResourceName:      graphBetaWindowsQualityUpdatePolicy.ResourceName + ".test_006",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Test 007: Scenario 7 - Assignments lifecycle minimal to maximal
func TestUnitResourceWindowsQualityUpdatePolicy_07_AssignmentsLifecycle_MinimalToMaximal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, qualityUpdateMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer qualityUpdateMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("007_assignments_lifecycle_minimal_to_maximal_step_1.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaWindowsQualityUpdatePolicy.ResourceName+".test_007").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaWindowsQualityUpdatePolicy.ResourceName+".test_007").Key("display_name").HasValue("unit-test-windows-quality-update-policy-007-lifecycle"),
					check.That(graphBetaWindowsQualityUpdatePolicy.ResourceName+".test_007").Key("assignments.#").HasValue("1"),
					resource.TestCheckTypeSetElemNestedAttrs(graphBetaWindowsQualityUpdatePolicy.ResourceName+".test_007", "assignments.*", map[string]string{
						"type":     "groupAssignmentTarget",
						"group_id": "11111111-1111-1111-1111-111111111111",
					}),
				),
			},
			{
				Config: loadUnitTestTerraform("007_assignments_lifecycle_minimal_to_maximal_step_2.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaWindowsQualityUpdatePolicy.ResourceName+".test_007").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaWindowsQualityUpdatePolicy.ResourceName+".test_007").Key("display_name").HasValue("unit-test-windows-quality-update-policy-007-lifecycle"),
					check.That(graphBetaWindowsQualityUpdatePolicy.ResourceName+".test_007").Key("assignments.#").HasValue("3"),
					resource.TestCheckTypeSetElemNestedAttrs(graphBetaWindowsQualityUpdatePolicy.ResourceName+".test_007", "assignments.*", map[string]string{
						"type":     "groupAssignmentTarget",
						"group_id": "11111111-1111-1111-1111-111111111111",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(graphBetaWindowsQualityUpdatePolicy.ResourceName+".test_007", "assignments.*", map[string]string{
						"type":     "groupAssignmentTarget",
						"group_id": "22222222-2222-2222-2222-222222222222",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(graphBetaWindowsQualityUpdatePolicy.ResourceName+".test_007", "assignments.*", map[string]string{
						"type":     "exclusionGroupAssignmentTarget",
						"group_id": "33333333-3333-3333-3333-333333333333",
					}),
				),
			},
			{
				ResourceName:      graphBetaWindowsQualityUpdatePolicy.ResourceName + ".test_007",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Test 008: Scenario 8 - Assignments lifecycle maximal to minimal
func TestUnitResourceWindowsQualityUpdatePolicy_08_AssignmentsLifecycle_MaximalToMinimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, qualityUpdateMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer qualityUpdateMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("008_assignments_lifecycle_maximal_to_minimal_step_1.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaWindowsQualityUpdatePolicy.ResourceName+".test_008").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaWindowsQualityUpdatePolicy.ResourceName+".test_008").Key("display_name").HasValue("unit-test-windows-quality-update-policy-008-lifecycle"),
					check.That(graphBetaWindowsQualityUpdatePolicy.ResourceName+".test_008").Key("assignments.#").HasValue("3"),
					resource.TestCheckTypeSetElemNestedAttrs(graphBetaWindowsQualityUpdatePolicy.ResourceName+".test_008", "assignments.*", map[string]string{
						"type":     "groupAssignmentTarget",
						"group_id": "11111111-1111-1111-1111-111111111111",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(graphBetaWindowsQualityUpdatePolicy.ResourceName+".test_008", "assignments.*", map[string]string{
						"type":     "groupAssignmentTarget",
						"group_id": "22222222-2222-2222-2222-222222222222",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(graphBetaWindowsQualityUpdatePolicy.ResourceName+".test_008", "assignments.*", map[string]string{
						"type":     "exclusionGroupAssignmentTarget",
						"group_id": "33333333-3333-3333-3333-333333333333",
					}),
				),
			},
			{
				Config: loadUnitTestTerraform("008_assignments_lifecycle_maximal_to_minimal_step_2.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaWindowsQualityUpdatePolicy.ResourceName+".test_008").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaWindowsQualityUpdatePolicy.ResourceName+".test_008").Key("display_name").HasValue("unit-test-windows-quality-update-policy-008-lifecycle"),
					check.That(graphBetaWindowsQualityUpdatePolicy.ResourceName+".test_008").Key("assignments.#").HasValue("1"),
					resource.TestCheckTypeSetElemNestedAttrs(graphBetaWindowsQualityUpdatePolicy.ResourceName+".test_008", "assignments.*", map[string]string{
						"type":     "groupAssignmentTarget",
						"group_id": "11111111-1111-1111-1111-111111111111",
					}),
				),
			},
			{
				ResourceName:      graphBetaWindowsQualityUpdatePolicy.ResourceName + ".test_008",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
