package graphBetaWindowsQualityUpdateExpeditePolicy_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaWindowsQualityUpdateExpeditePolicy "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/windows_quality_update_expedite_policy"
	expediteMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/windows_quality_update_expedite_policy/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

func setupMockEnvironment() (*mocks.Mocks, *expediteMocks.WindowsQualityUpdateExpeditePolicyMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	expediteMock := &expediteMocks.WindowsQualityUpdateExpeditePolicyMock{}
	expediteMock.RegisterMocks()
	return mockClient, expediteMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *expediteMocks.WindowsQualityUpdateExpeditePolicyMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	expediteMock := &expediteMocks.WindowsQualityUpdateExpeditePolicyMock{}
	expediteMock.RegisterErrorMocks()
	return mockClient, expediteMock
}

func loadUnitTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/unit/" + filename)
	if err != nil {
		panic("failed to load unit test config " + filename + ": " + err.Error())
	}
	return config
}

func TestWindowsQualityUpdateExpeditePolicyResource_001_Scenario_Minimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, expediteMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer expediteMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("001_scenario_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaWindowsQualityUpdateExpeditePolicy.ResourceName+".test_001").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaWindowsQualityUpdateExpeditePolicy.ResourceName+".test_001").Key("display_name").HasValue("unit-test-windows-quality-update-expedite-policy-001-minimal"),
					check.That(graphBetaWindowsQualityUpdateExpeditePolicy.ResourceName+".test_001").Key("role_scope_tag_ids.#").HasValue("1"),
					check.That(graphBetaWindowsQualityUpdateExpeditePolicy.ResourceName+".test_001").Key("role_scope_tag_ids.0").HasValue("0"),
				),
			},
			{
				ResourceName:      graphBetaWindowsQualityUpdateExpeditePolicy.ResourceName + ".test_001",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestWindowsQualityUpdateExpeditePolicyResource_002_Scenario_Maximal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, expediteMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer expediteMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("002_scenario_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaWindowsQualityUpdateExpeditePolicy.ResourceName+".test_002").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaWindowsQualityUpdateExpeditePolicy.ResourceName+".test_002").Key("display_name").HasValue("unit-test-windows-quality-update-expedite-policy-002-maximal"),
					check.That(graphBetaWindowsQualityUpdateExpeditePolicy.ResourceName+".test_002").Key("description").HasValue("Scenario 2: Maximal configuration without assignments"),
					check.That(graphBetaWindowsQualityUpdateExpeditePolicy.ResourceName+".test_002").Key("role_scope_tag_ids.#").HasValue("2"),
					check.That(graphBetaWindowsQualityUpdateExpeditePolicy.ResourceName+".test_002").Key("expedited_update_settings.quality_update_release").HasValue("2025-11-20T00:00:00Z"),
					check.That(graphBetaWindowsQualityUpdateExpeditePolicy.ResourceName+".test_002").Key("expedited_update_settings.days_until_forced_reboot").HasValue("1"),
				),
			},
			{
				ResourceName:      graphBetaWindowsQualityUpdateExpeditePolicy.ResourceName + ".test_002",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestWindowsQualityUpdateExpeditePolicyResource_003_Lifecycle_MinimalToMaximal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, expediteMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer expediteMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("003_lifecycle_minimal_to_maximal_step_1.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaWindowsQualityUpdateExpeditePolicy.ResourceName+".test_003").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaWindowsQualityUpdateExpeditePolicy.ResourceName+".test_003").Key("display_name").HasValue("unit-test-windows-quality-update-expedite-policy-003-lifecycle"),
					check.That(graphBetaWindowsQualityUpdateExpeditePolicy.ResourceName+".test_003").Key("role_scope_tag_ids.#").HasValue("1"),
				),
			},
			{
				Config: loadUnitTestTerraform("003_lifecycle_minimal_to_maximal_step_2.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaWindowsQualityUpdateExpeditePolicy.ResourceName+".test_003").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaWindowsQualityUpdateExpeditePolicy.ResourceName+".test_003").Key("display_name").HasValue("unit-test-windows-quality-update-expedite-policy-003-lifecycle"),
					check.That(graphBetaWindowsQualityUpdateExpeditePolicy.ResourceName+".test_003").Key("description").HasValue("Lifecycle Step 2: Updated to maximal configuration"),
					check.That(graphBetaWindowsQualityUpdateExpeditePolicy.ResourceName+".test_003").Key("role_scope_tag_ids.#").HasValue("2"),
					check.That(graphBetaWindowsQualityUpdateExpeditePolicy.ResourceName+".test_003").Key("expedited_update_settings.quality_update_release").HasValue("2025-11-20T00:00:00Z"),
					check.That(graphBetaWindowsQualityUpdateExpeditePolicy.ResourceName+".test_003").Key("expedited_update_settings.days_until_forced_reboot").HasValue("1"),
				),
			},
			{
				ResourceName:      graphBetaWindowsQualityUpdateExpeditePolicy.ResourceName + ".test_003",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestWindowsQualityUpdateExpeditePolicyResource_004_Lifecycle_MaximalToMinimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, expediteMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer expediteMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("004_lifecycle_maximal_to_minimal_step_1.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaWindowsQualityUpdateExpeditePolicy.ResourceName+".test_004").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaWindowsQualityUpdateExpeditePolicy.ResourceName+".test_004").Key("display_name").HasValue("unit-test-windows-quality-update-expedite-policy-004-lifecycle"),
					check.That(graphBetaWindowsQualityUpdateExpeditePolicy.ResourceName+".test_004").Key("description").HasValue("Lifecycle Step 1: Starting with maximal configuration"),
					check.That(graphBetaWindowsQualityUpdateExpeditePolicy.ResourceName+".test_004").Key("role_scope_tag_ids.#").HasValue("2"),
					check.That(graphBetaWindowsQualityUpdateExpeditePolicy.ResourceName+".test_004").Key("expedited_update_settings.quality_update_release").HasValue("2025-11-20T00:00:00Z"),
					check.That(graphBetaWindowsQualityUpdateExpeditePolicy.ResourceName+".test_004").Key("expedited_update_settings.days_until_forced_reboot").HasValue("1"),
				),
			},
			{
				Config: loadUnitTestTerraform("004_lifecycle_maximal_to_minimal_step_2.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaWindowsQualityUpdateExpeditePolicy.ResourceName+".test_004").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaWindowsQualityUpdateExpeditePolicy.ResourceName+".test_004").Key("display_name").HasValue("unit-test-windows-quality-update-expedite-policy-004-lifecycle"),
					check.That(graphBetaWindowsQualityUpdateExpeditePolicy.ResourceName+".test_004").Key("role_scope_tag_ids.#").HasValue("1"),
				),
			},
			{
				ResourceName:      graphBetaWindowsQualityUpdateExpeditePolicy.ResourceName + ".test_004",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestWindowsQualityUpdateExpeditePolicyResource_005_AssignmentsMinimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, expediteMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer expediteMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("005_assignments_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaWindowsQualityUpdateExpeditePolicy.ResourceName+".test_005").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaWindowsQualityUpdateExpeditePolicy.ResourceName+".test_005").Key("display_name").HasValue("unit-test-windows-quality-update-expedite-policy-005-assignments-minimal"),
					check.That(graphBetaWindowsQualityUpdateExpeditePolicy.ResourceName+".test_005").Key("assignments.#").HasValue("1"),
					resource.TestCheckTypeSetElemNestedAttrs(graphBetaWindowsQualityUpdateExpeditePolicy.ResourceName+".test_005", "assignments.*", map[string]string{
						"type":     "groupAssignmentTarget",
						"group_id": "11111111-1111-1111-1111-111111111111",
					}),
				),
			},
			{
				ResourceName:      graphBetaWindowsQualityUpdateExpeditePolicy.ResourceName + ".test_005",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestWindowsQualityUpdateExpeditePolicyResource_006_AssignmentsMaximal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, expediteMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer expediteMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("006_assignments_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaWindowsQualityUpdateExpeditePolicy.ResourceName+".test_006").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaWindowsQualityUpdateExpeditePolicy.ResourceName+".test_006").Key("display_name").HasValue("unit-test-windows-quality-update-expedite-policy-006-assignments-maximal"),
					check.That(graphBetaWindowsQualityUpdateExpeditePolicy.ResourceName+".test_006").Key("assignments.#").HasValue("3"),
					resource.TestCheckTypeSetElemNestedAttrs(graphBetaWindowsQualityUpdateExpeditePolicy.ResourceName+".test_006", "assignments.*", map[string]string{
						"type":     "groupAssignmentTarget",
						"group_id": "44444444-4444-4444-4444-444444444444",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(graphBetaWindowsQualityUpdateExpeditePolicy.ResourceName+".test_006", "assignments.*", map[string]string{
						"type":     "groupAssignmentTarget",
						"group_id": "33333333-3333-3333-3333-333333333333",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(graphBetaWindowsQualityUpdateExpeditePolicy.ResourceName+".test_006", "assignments.*", map[string]string{
						"type":     "exclusionGroupAssignmentTarget",
						"group_id": "77777777-7777-7777-7777-777777777777",
					}),
				),
			},
			{
				ResourceName:      graphBetaWindowsQualityUpdateExpeditePolicy.ResourceName + ".test_006",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestWindowsQualityUpdateExpeditePolicyResource_007_AssignmentsLifecycle_MinimalToMaximal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, expediteMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer expediteMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("007_assignments_lifecycle_minimal_to_maximal_step_1.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaWindowsQualityUpdateExpeditePolicy.ResourceName+".test_007").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaWindowsQualityUpdateExpeditePolicy.ResourceName+".test_007").Key("display_name").HasValue("unit-test-windows-quality-update-expedite-policy-007-lifecycle"),
					check.That(graphBetaWindowsQualityUpdateExpeditePolicy.ResourceName+".test_007").Key("assignments.#").HasValue("1"),
					resource.TestCheckTypeSetElemNestedAttrs(graphBetaWindowsQualityUpdateExpeditePolicy.ResourceName+".test_007", "assignments.*", map[string]string{
						"type":     "groupAssignmentTarget",
						"group_id": "11111111-1111-1111-1111-111111111111",
					}),
				),
			},
			{
				Config: loadUnitTestTerraform("007_assignments_lifecycle_minimal_to_maximal_step_2.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaWindowsQualityUpdateExpeditePolicy.ResourceName+".test_007").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaWindowsQualityUpdateExpeditePolicy.ResourceName+".test_007").Key("display_name").HasValue("unit-test-windows-quality-update-expedite-policy-007-lifecycle"),
					check.That(graphBetaWindowsQualityUpdateExpeditePolicy.ResourceName+".test_007").Key("assignments.#").HasValue("3"),
					resource.TestCheckTypeSetElemNestedAttrs(graphBetaWindowsQualityUpdateExpeditePolicy.ResourceName+".test_007", "assignments.*", map[string]string{
						"type":     "groupAssignmentTarget",
						"group_id": "11111111-1111-1111-1111-111111111111",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(graphBetaWindowsQualityUpdateExpeditePolicy.ResourceName+".test_007", "assignments.*", map[string]string{
						"type":     "groupAssignmentTarget",
						"group_id": "22222222-2222-2222-2222-222222222222",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(graphBetaWindowsQualityUpdateExpeditePolicy.ResourceName+".test_007", "assignments.*", map[string]string{
						"type":     "exclusionGroupAssignmentTarget",
						"group_id": "33333333-3333-3333-3333-333333333333",
					}),
				),
			},
			{
				ResourceName:      graphBetaWindowsQualityUpdateExpeditePolicy.ResourceName + ".test_007",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestWindowsQualityUpdateExpeditePolicyResource_008_AssignmentsLifecycle_MaximalToMinimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, expediteMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer expediteMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("008_assignments_lifecycle_maximal_to_minimal_step_1.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaWindowsQualityUpdateExpeditePolicy.ResourceName+".test_008").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaWindowsQualityUpdateExpeditePolicy.ResourceName+".test_008").Key("display_name").HasValue("unit-test-windows-quality-update-expedite-policy-008-lifecycle"),
					check.That(graphBetaWindowsQualityUpdateExpeditePolicy.ResourceName+".test_008").Key("assignments.#").HasValue("3"),
					resource.TestCheckTypeSetElemNestedAttrs(graphBetaWindowsQualityUpdateExpeditePolicy.ResourceName+".test_008", "assignments.*", map[string]string{
						"type":     "groupAssignmentTarget",
						"group_id": "11111111-1111-1111-1111-111111111111",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(graphBetaWindowsQualityUpdateExpeditePolicy.ResourceName+".test_008", "assignments.*", map[string]string{
						"type":     "groupAssignmentTarget",
						"group_id": "22222222-2222-2222-2222-222222222222",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(graphBetaWindowsQualityUpdateExpeditePolicy.ResourceName+".test_008", "assignments.*", map[string]string{
						"type":     "exclusionGroupAssignmentTarget",
						"group_id": "33333333-3333-3333-3333-333333333333",
					}),
				),
			},
			{
				Config: loadUnitTestTerraform("008_assignments_lifecycle_maximal_to_minimal_step_2.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaWindowsQualityUpdateExpeditePolicy.ResourceName+".test_008").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaWindowsQualityUpdateExpeditePolicy.ResourceName+".test_008").Key("display_name").HasValue("unit-test-windows-quality-update-expedite-policy-008-lifecycle"),
					check.That(graphBetaWindowsQualityUpdateExpeditePolicy.ResourceName+".test_008").Key("assignments.#").HasValue("1"),
					resource.TestCheckTypeSetElemNestedAttrs(graphBetaWindowsQualityUpdateExpeditePolicy.ResourceName+".test_008", "assignments.*", map[string]string{
						"type":     "groupAssignmentTarget",
						"group_id": "11111111-1111-1111-1111-111111111111",
					}),
				),
			},
			{
				ResourceName:      graphBetaWindowsQualityUpdateExpeditePolicy.ResourceName + ".test_008",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestWindowsQualityUpdateExpeditePolicyResource_009_IntentionalErrors(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, expediteMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer expediteMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("001_scenario_minimal.tf"),
				ExpectError: regexp.MustCompile(`(?i)(error|failed|bad request)`),
			},
		},
	})
}
