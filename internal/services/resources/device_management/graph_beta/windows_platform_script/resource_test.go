package graphBetaWindowsPlatformScript_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaWindowsPlatformScript "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/windows_platform_script"
	scriptMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/windows_platform_script/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

func setupMockEnvironment() (*mocks.Mocks, *scriptMocks.WindowsPlatformScriptMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	scriptMock := &scriptMocks.WindowsPlatformScriptMock{}
	scriptMock.RegisterMocks()
	return mockClient, scriptMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *scriptMocks.WindowsPlatformScriptMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	scriptMock := &scriptMocks.WindowsPlatformScriptMock{}
	scriptMock.RegisterErrorMocks()
	return mockClient, scriptMock
}

func loadUnitTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/unit/" + filename)
	if err != nil {
		panic("failed to load unit test config " + filename + ": " + err.Error())
	}
	return config
}

func TestWindowsPlatformScriptResource_001_Scenario_Minimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, scriptMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer scriptMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("001_scenario_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaWindowsPlatformScript.ResourceName+".test_001").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaWindowsPlatformScript.ResourceName+".test_001").Key("display_name").HasValue("unit-test-windows-platform-script-001-minimal"),
					check.That(graphBetaWindowsPlatformScript.ResourceName+".test_001").Key("run_as_account").HasValue("system"),
					check.That(graphBetaWindowsPlatformScript.ResourceName+".test_001").Key("role_scope_tag_ids.#").HasValue("1"),
					check.That(graphBetaWindowsPlatformScript.ResourceName+".test_001").Key("role_scope_tag_ids.0").HasValue("0"),
				),
			},
			{
				ResourceName:      graphBetaWindowsPlatformScript.ResourceName + ".test_001",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestWindowsPlatformScriptResource_002_Scenario_Maximal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, scriptMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer scriptMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("002_scenario_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaWindowsPlatformScript.ResourceName+".test_002").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaWindowsPlatformScript.ResourceName+".test_002").Key("display_name").HasValue("unit-test-windows-platform-script-002-maximal"),
					check.That(graphBetaWindowsPlatformScript.ResourceName+".test_002").Key("description").HasValue("Maximal test configuration for Windows platform script"),
					check.That(graphBetaWindowsPlatformScript.ResourceName+".test_002").Key("run_as_account").HasValue("user"),
					check.That(graphBetaWindowsPlatformScript.ResourceName+".test_002").Key("enforce_signature_check").HasValue("true"),
					check.That(graphBetaWindowsPlatformScript.ResourceName+".test_002").Key("run_as_32_bit").HasValue("true"),
					check.That(graphBetaWindowsPlatformScript.ResourceName+".test_002").Key("role_scope_tag_ids.#").HasValue("2"),
				),
			},
			{
				ResourceName:      graphBetaWindowsPlatformScript.ResourceName + ".test_002",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestWindowsPlatformScriptResource_003_Lifecycle_MinimalToMaximal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, scriptMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer scriptMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("003_lifecycle_minimal_to_maximal_step_1.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaWindowsPlatformScript.ResourceName+".test_003").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaWindowsPlatformScript.ResourceName+".test_003").Key("display_name").HasValue("unit-test-windows-platform-script-003-lifecycle"),
					check.That(graphBetaWindowsPlatformScript.ResourceName+".test_003").Key("run_as_account").HasValue("system"),
				),
			},
			{
				Config: loadUnitTestTerraform("003_lifecycle_minimal_to_maximal_step_2.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaWindowsPlatformScript.ResourceName+".test_003").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaWindowsPlatformScript.ResourceName+".test_003").Key("display_name").HasValue("unit-test-windows-platform-script-003-lifecycle"),
					check.That(graphBetaWindowsPlatformScript.ResourceName+".test_003").Key("description").HasValue("Maximal lifecycle test configuration"),
					check.That(graphBetaWindowsPlatformScript.ResourceName+".test_003").Key("run_as_account").HasValue("user"),
					check.That(graphBetaWindowsPlatformScript.ResourceName+".test_003").Key("enforce_signature_check").HasValue("true"),
					check.That(graphBetaWindowsPlatformScript.ResourceName+".test_003").Key("run_as_32_bit").HasValue("true"),
					check.That(graphBetaWindowsPlatformScript.ResourceName+".test_003").Key("role_scope_tag_ids.#").HasValue("2"),
				),
			},
		},
	})
}

func TestWindowsPlatformScriptResource_004_Lifecycle_MaximalToMinimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, scriptMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer scriptMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("004_lifecycle_maximal_to_minimal_step_1.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaWindowsPlatformScript.ResourceName+".test_004").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaWindowsPlatformScript.ResourceName+".test_004").Key("display_name").HasValue("unit-test-windows-platform-script-004-lifecycle"),
					check.That(graphBetaWindowsPlatformScript.ResourceName+".test_004").Key("description").HasValue("Maximal lifecycle test configuration"),
					check.That(graphBetaWindowsPlatformScript.ResourceName+".test_004").Key("run_as_account").HasValue("user"),
					check.That(graphBetaWindowsPlatformScript.ResourceName+".test_004").Key("enforce_signature_check").HasValue("true"),
					check.That(graphBetaWindowsPlatformScript.ResourceName+".test_004").Key("run_as_32_bit").HasValue("true"),
				),
			},
			{
				Config: loadUnitTestTerraform("004_lifecycle_maximal_to_minimal_step_2.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaWindowsPlatformScript.ResourceName+".test_004").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaWindowsPlatformScript.ResourceName+".test_004").Key("display_name").HasValue("unit-test-windows-platform-script-004-lifecycle"),
					check.That(graphBetaWindowsPlatformScript.ResourceName+".test_004").Key("run_as_account").HasValue("system"),
					check.That(graphBetaWindowsPlatformScript.ResourceName+".test_004").Key("role_scope_tag_ids.#").HasValue("1"),
				),
			},
		},
	})
}

func TestWindowsPlatformScriptResource_005_AssignmentsMinimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, scriptMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer scriptMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("005_assignments_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaWindowsPlatformScript.ResourceName+".test_005").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaWindowsPlatformScript.ResourceName+".test_005").Key("display_name").HasValue("unit-test-windows-platform-script-005-assignments-minimal"),
					check.That(graphBetaWindowsPlatformScript.ResourceName+".test_005").Key("assignments.#").HasValue("1"),
				),
			},
			{
				ResourceName:      graphBetaWindowsPlatformScript.ResourceName + ".test_005",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestWindowsPlatformScriptResource_006_AssignmentsMaximal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, scriptMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer scriptMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("006_assignments_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaWindowsPlatformScript.ResourceName+".test_006").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaWindowsPlatformScript.ResourceName+".test_006").Key("display_name").HasValue("unit-test-windows-platform-script-006-assignments-maximal"),
					check.That(graphBetaWindowsPlatformScript.ResourceName+".test_006").Key("description").HasValue("Maximal test with multiple assignments"),
					check.That(graphBetaWindowsPlatformScript.ResourceName+".test_006").Key("assignments.#").HasValue("5"),
				),
			},
			{
				ResourceName:      graphBetaWindowsPlatformScript.ResourceName + ".test_006",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestWindowsPlatformScriptResource_007_AssignmentsLifecycle_MinimalToMaximal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, scriptMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer scriptMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("007_assignments_lifecycle_minimal_to_maximal_step_1.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaWindowsPlatformScript.ResourceName+".test_007").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaWindowsPlatformScript.ResourceName+".test_007").Key("display_name").HasValue("unit-test-windows-platform-script-007-assignments-lifecycle"),
					check.That(graphBetaWindowsPlatformScript.ResourceName+".test_007").Key("assignments.#").HasValue("1"),
				),
			},
			{
				Config: loadUnitTestTerraform("007_assignments_lifecycle_minimal_to_maximal_step_2.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaWindowsPlatformScript.ResourceName+".test_007").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaWindowsPlatformScript.ResourceName+".test_007").Key("display_name").HasValue("unit-test-windows-platform-script-007-assignments-lifecycle"),
					check.That(graphBetaWindowsPlatformScript.ResourceName+".test_007").Key("assignments.#").HasValue("5"),
				),
			},
		},
	})
}

func TestWindowsPlatformScriptResource_008_AssignmentsLifecycle_MaximalToMinimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, scriptMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer scriptMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("008_assignments_lifecycle_maximal_to_minimal_step_1.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaWindowsPlatformScript.ResourceName+".test_008").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaWindowsPlatformScript.ResourceName+".test_008").Key("display_name").HasValue("unit-test-windows-platform-script-008-assignments-lifecycle"),
					check.That(graphBetaWindowsPlatformScript.ResourceName+".test_008").Key("description").HasValue("Maximal assignments lifecycle test"),
					check.That(graphBetaWindowsPlatformScript.ResourceName+".test_008").Key("assignments.#").HasValue("5"),
				),
			},
			{
				Config: loadUnitTestTerraform("008_assignments_lifecycle_maximal_to_minimal_step_2.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaWindowsPlatformScript.ResourceName+".test_008").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaWindowsPlatformScript.ResourceName+".test_008").Key("display_name").HasValue("unit-test-windows-platform-script-008-assignments-lifecycle"),
					check.That(graphBetaWindowsPlatformScript.ResourceName+".test_008").Key("description").HasValue("Maximal assignments lifecycle test"),
					check.That(graphBetaWindowsPlatformScript.ResourceName+".test_008").Key("assignments.#").HasValue("1"),
				),
			},
		},
	})
}

func TestWindowsPlatformScriptResource_009_ErrorHandling(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, scriptMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer scriptMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("009_error_scenario.tf"),
				ExpectError: regexp.MustCompile("Invalid Windows Platform Script data"),
			},
		},
	})
}
