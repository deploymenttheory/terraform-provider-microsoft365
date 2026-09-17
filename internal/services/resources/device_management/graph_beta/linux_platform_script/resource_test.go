package graphBetaLinuxPlatformScript_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaLinuxPlatformScript "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/linux_platform_script"
	scriptMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/linux_platform_script/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/jarcoal/httpmock"
)

const minimalScriptContent = "#!/bin/bash\nprintf \"Linux platform script test\\n\"\n"
const maximalScriptContent = "#!/bin/bash\n# Updated script: café\nprintf \"Linux platform script updated\\n\"\n"

func setupMockEnvironment() (*mocks.Mocks, *scriptMocks.LinuxPlatformScriptMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	scriptMock := &scriptMocks.LinuxPlatformScriptMock{}
	scriptMock.RegisterMocks()
	return mockClient, scriptMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *scriptMocks.LinuxPlatformScriptMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	scriptMock := &scriptMocks.LinuxPlatformScriptMock{}
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

func TestUnitResourceLinuxPlatformScript_01_Scenario_Minimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, scriptMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer scriptMock.CleanupMockState()
	scriptMock.ExpectedWrites = []string{"validate_create/post_001_scenario_minimal.json"}

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		CheckDestroy:             func(_ *terraform.State) error { return scriptMock.CheckDestroyed() },
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("001_scenario_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaLinuxPlatformScript.ResourceName+".test_001").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaLinuxPlatformScript.ResourceName+".test_001").Key("name").HasValue("unit-test-linux-platform-script-001"),
					check.That(graphBetaLinuxPlatformScript.ResourceName+".test_001").Key("platforms").HasValue("linux"),
					check.That(graphBetaLinuxPlatformScript.ResourceName+".test_001").Key("technologies.0").HasValue("linuxMdm"),
					check.That(graphBetaLinuxPlatformScript.ResourceName+".test_001").Key("script_content").HasValue(minimalScriptContent),
					check.That(graphBetaLinuxPlatformScript.ResourceName+".test_001").Key("execution_context").HasValue("user"),
					check.That(graphBetaLinuxPlatformScript.ResourceName+".test_001").Key("execution_frequency").HasValue("15minutes"),
					check.That(graphBetaLinuxPlatformScript.ResourceName+".test_001").Key("execution_retries").HasValue("1"),
					check.That(graphBetaLinuxPlatformScript.ResourceName+".test_001").Key("role_scope_tag_ids.0").HasValue("0"),
					resource.TestCheckNoResourceAttr(graphBetaLinuxPlatformScript.ResourceName+".test_001", "assignments.#"),
				),
			},
			{
				ResourceName:            graphBetaLinuxPlatformScript.ResourceName + ".test_001",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

func TestUnitResourceLinuxPlatformScript_02_Scenario_Maximal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, scriptMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer scriptMock.CleanupMockState()
	scriptMock.ExpectedWrites = []string{"validate_create/post_002_scenario_maximal.json"}

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		CheckDestroy:             func(_ *terraform.State) error { return scriptMock.CheckDestroyed() },
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("002_scenario_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaLinuxPlatformScript.ResourceName+".test_002").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaLinuxPlatformScript.ResourceName+".test_002").Key("name").HasValue("unit-test-linux-platform-script-002"),
					check.That(graphBetaLinuxPlatformScript.ResourceName+".test_002").Key("platforms").HasValue("linux"),
					check.That(graphBetaLinuxPlatformScript.ResourceName+".test_002").Key("technologies.0").HasValue("linuxMdm"),
					check.That(graphBetaLinuxPlatformScript.ResourceName+".test_002").Key("description").HasValue("Maximal Linux platform script"),
					check.That(graphBetaLinuxPlatformScript.ResourceName+".test_002").Key("script_content").HasValue(maximalScriptContent),
					check.That(graphBetaLinuxPlatformScript.ResourceName+".test_002").Key("execution_context").HasValue("root"),
					check.That(graphBetaLinuxPlatformScript.ResourceName+".test_002").Key("execution_frequency").HasValue("1day"),
					check.That(graphBetaLinuxPlatformScript.ResourceName+".test_002").Key("execution_retries").HasValue("3"),
					check.That(graphBetaLinuxPlatformScript.ResourceName+".test_002").Key("role_scope_tag_ids.#").HasValue("2"),
					resource.TestCheckNoResourceAttr(graphBetaLinuxPlatformScript.ResourceName+".test_002", "assignments.#"),
				),
			},
			{
				ResourceName:            graphBetaLinuxPlatformScript.ResourceName + ".test_002",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

func TestUnitResourceLinuxPlatformScript_03_Lifecycle_MinimalToMaximal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, scriptMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer scriptMock.CleanupMockState()
	scriptMock.ExpectedWrites = []string{"validate_create/post_001_scenario_minimal.json", "validate_update/put_003_lifecycle_maximal.json"}

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		CheckDestroy:             func(_ *terraform.State) error { return scriptMock.CheckDestroyed() },
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("003_lifecycle_minimal_to_maximal_step_1.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaLinuxPlatformScript.ResourceName+".test_003").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaLinuxPlatformScript.ResourceName+".test_003").Key("name").HasValue("unit-test-linux-platform-script-003"),
					check.That(graphBetaLinuxPlatformScript.ResourceName+".test_003").Key("platforms").HasValue("linux"),
					check.That(graphBetaLinuxPlatformScript.ResourceName+".test_003").Key("technologies.0").HasValue("linuxMdm"),
					check.That(graphBetaLinuxPlatformScript.ResourceName+".test_003").Key("script_content").HasValue(minimalScriptContent),
					check.That(graphBetaLinuxPlatformScript.ResourceName+".test_003").Key("execution_context").HasValue("user"),
					check.That(graphBetaLinuxPlatformScript.ResourceName+".test_003").Key("execution_frequency").HasValue("15minutes"),
					check.That(graphBetaLinuxPlatformScript.ResourceName+".test_003").Key("execution_retries").HasValue("1"),
					check.That(graphBetaLinuxPlatformScript.ResourceName+".test_003").Key("role_scope_tag_ids.0").HasValue("0"),
					resource.TestCheckNoResourceAttr(graphBetaLinuxPlatformScript.ResourceName+".test_003", "assignments.#"),
				),
			},
			{
				Config: loadUnitTestTerraform("003_lifecycle_minimal_to_maximal_step_2.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaLinuxPlatformScript.ResourceName+".test_003").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaLinuxPlatformScript.ResourceName+".test_003").Key("name").HasValue("unit-test-linux-platform-script-003"),
					check.That(graphBetaLinuxPlatformScript.ResourceName+".test_003").Key("platforms").HasValue("linux"),
					check.That(graphBetaLinuxPlatformScript.ResourceName+".test_003").Key("technologies.0").HasValue("linuxMdm"),
					check.That(graphBetaLinuxPlatformScript.ResourceName+".test_003").Key("description").HasValue("Maximal Linux platform script"),
					check.That(graphBetaLinuxPlatformScript.ResourceName+".test_003").Key("script_content").HasValue(maximalScriptContent),
					check.That(graphBetaLinuxPlatformScript.ResourceName+".test_003").Key("execution_context").HasValue("root"),
					check.That(graphBetaLinuxPlatformScript.ResourceName+".test_003").Key("execution_frequency").HasValue("1day"),
					check.That(graphBetaLinuxPlatformScript.ResourceName+".test_003").Key("execution_retries").HasValue("3"),
					check.That(graphBetaLinuxPlatformScript.ResourceName+".test_003").Key("role_scope_tag_ids.#").HasValue("2"),
					resource.TestCheckNoResourceAttr(graphBetaLinuxPlatformScript.ResourceName+".test_003", "assignments.#"),
				),
			},
			{
				ResourceName:            graphBetaLinuxPlatformScript.ResourceName + ".test_003",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

func TestUnitResourceLinuxPlatformScript_04_Lifecycle_MaximalToMinimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, scriptMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer scriptMock.CleanupMockState()
	scriptMock.ExpectedWrites = []string{"validate_create/post_002_scenario_maximal.json", "validate_update/put_004_lifecycle_minimal.json"}

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		CheckDestroy:             func(_ *terraform.State) error { return scriptMock.CheckDestroyed() },
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("004_lifecycle_maximal_to_minimal_step_1.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaLinuxPlatformScript.ResourceName+".test_004").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaLinuxPlatformScript.ResourceName+".test_004").Key("name").HasValue("unit-test-linux-platform-script-004"),
					check.That(graphBetaLinuxPlatformScript.ResourceName+".test_004").Key("platforms").HasValue("linux"),
					check.That(graphBetaLinuxPlatformScript.ResourceName+".test_004").Key("technologies.0").HasValue("linuxMdm"),
					check.That(graphBetaLinuxPlatformScript.ResourceName+".test_004").Key("description").HasValue("Maximal Linux platform script"),
					check.That(graphBetaLinuxPlatformScript.ResourceName+".test_004").Key("script_content").HasValue(maximalScriptContent),
					check.That(graphBetaLinuxPlatformScript.ResourceName+".test_004").Key("execution_context").HasValue("root"),
					check.That(graphBetaLinuxPlatformScript.ResourceName+".test_004").Key("execution_frequency").HasValue("1day"),
					check.That(graphBetaLinuxPlatformScript.ResourceName+".test_004").Key("execution_retries").HasValue("3"),
					check.That(graphBetaLinuxPlatformScript.ResourceName+".test_004").Key("role_scope_tag_ids.#").HasValue("2"),
					resource.TestCheckNoResourceAttr(graphBetaLinuxPlatformScript.ResourceName+".test_004", "assignments.#"),
				),
			},
			{
				Config: loadUnitTestTerraform("004_lifecycle_maximal_to_minimal_step_2.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaLinuxPlatformScript.ResourceName+".test_004").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaLinuxPlatformScript.ResourceName+".test_004").Key("name").HasValue("unit-test-linux-platform-script-004"),
					check.That(graphBetaLinuxPlatformScript.ResourceName+".test_004").Key("platforms").HasValue("linux"),
					check.That(graphBetaLinuxPlatformScript.ResourceName+".test_004").Key("technologies.0").HasValue("linuxMdm"),
					check.That(graphBetaLinuxPlatformScript.ResourceName+".test_004").Key("script_content").HasValue(minimalScriptContent),
					check.That(graphBetaLinuxPlatformScript.ResourceName+".test_004").Key("execution_context").HasValue("user"),
					check.That(graphBetaLinuxPlatformScript.ResourceName+".test_004").Key("execution_frequency").HasValue("15minutes"),
					check.That(graphBetaLinuxPlatformScript.ResourceName+".test_004").Key("execution_retries").HasValue("1"),
					check.That(graphBetaLinuxPlatformScript.ResourceName+".test_004").Key("role_scope_tag_ids.0").HasValue("0"),
					resource.TestCheckNoResourceAttr(graphBetaLinuxPlatformScript.ResourceName+".test_004", "assignments.#"),
				),
			},
			{
				ResourceName:            graphBetaLinuxPlatformScript.ResourceName + ".test_004",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

func TestUnitResourceLinuxPlatformScript_05_AssignmentsMinimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, scriptMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer scriptMock.CleanupMockState()
	scriptMock.ExpectedWrites = []string{"validate_create/post_001_scenario_minimal.json"}

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		CheckDestroy:             func(_ *terraform.State) error { return scriptMock.CheckDestroyed() },
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("005_assignments_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaLinuxPlatformScript.ResourceName+".test_005").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaLinuxPlatformScript.ResourceName+".test_005").Key("name").HasValue("unit-test-linux-platform-script-005"),
					check.That(graphBetaLinuxPlatformScript.ResourceName+".test_005").Key("platforms").HasValue("linux"),
					check.That(graphBetaLinuxPlatformScript.ResourceName+".test_005").Key("technologies.0").HasValue("linuxMdm"),
					check.That(graphBetaLinuxPlatformScript.ResourceName+".test_005").Key("script_content").HasValue(minimalScriptContent),
					check.That(graphBetaLinuxPlatformScript.ResourceName+".test_005").Key("execution_context").HasValue("user"),
					check.That(graphBetaLinuxPlatformScript.ResourceName+".test_005").Key("execution_frequency").HasValue("15minutes"),
					check.That(graphBetaLinuxPlatformScript.ResourceName+".test_005").Key("execution_retries").HasValue("1"),
					check.That(graphBetaLinuxPlatformScript.ResourceName+".test_005").Key("role_scope_tag_ids.0").HasValue("0"),
					check.That(graphBetaLinuxPlatformScript.ResourceName+".test_005").Key("assignments.#").HasValue("1"),
					resource.TestCheckTypeSetElemNestedAttrs(graphBetaLinuxPlatformScript.ResourceName+".test_005", "assignments.*", map[string]string{"type": "groupAssignmentTarget", "group_id": "11111111-1111-1111-1111-111111111111", "filter_type": "none", "filter_id": "00000000-0000-0000-0000-000000000000"}),
				),
			},
			{
				ResourceName:            graphBetaLinuxPlatformScript.ResourceName + ".test_005",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

func TestUnitResourceLinuxPlatformScript_06_AssignmentsMaximal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, scriptMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer scriptMock.CleanupMockState()
	scriptMock.ExpectedWrites = []string{"validate_create/post_001_scenario_minimal.json"}

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		CheckDestroy:             func(_ *terraform.State) error { return scriptMock.CheckDestroyed() },
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("006_assignments_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaLinuxPlatformScript.ResourceName+".test_006").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaLinuxPlatformScript.ResourceName+".test_006").Key("name").HasValue("unit-test-linux-platform-script-006"),
					check.That(graphBetaLinuxPlatformScript.ResourceName+".test_006").Key("platforms").HasValue("linux"),
					check.That(graphBetaLinuxPlatformScript.ResourceName+".test_006").Key("technologies.0").HasValue("linuxMdm"),
					check.That(graphBetaLinuxPlatformScript.ResourceName+".test_006").Key("script_content").HasValue(minimalScriptContent),
					check.That(graphBetaLinuxPlatformScript.ResourceName+".test_006").Key("execution_context").HasValue("user"),
					check.That(graphBetaLinuxPlatformScript.ResourceName+".test_006").Key("execution_frequency").HasValue("15minutes"),
					check.That(graphBetaLinuxPlatformScript.ResourceName+".test_006").Key("execution_retries").HasValue("1"),
					check.That(graphBetaLinuxPlatformScript.ResourceName+".test_006").Key("role_scope_tag_ids.0").HasValue("0"),
					check.That(graphBetaLinuxPlatformScript.ResourceName+".test_006").Key("assignments.#").HasValue("5"),
					resource.TestCheckTypeSetElemNestedAttrs(graphBetaLinuxPlatformScript.ResourceName+".test_006", "assignments.*", map[string]string{"type": "groupAssignmentTarget", "group_id": "11111111-1111-1111-1111-111111111111", "filter_type": "none", "filter_id": "00000000-0000-0000-0000-000000000000"}),
					resource.TestCheckTypeSetElemNestedAttrs(graphBetaLinuxPlatformScript.ResourceName+".test_006", "assignments.*", map[string]string{"type": "groupAssignmentTarget", "group_id": "22222222-2222-2222-2222-222222222222", "filter_type": "include", "filter_id": "44444444-4444-4444-4444-444444444444"}),
					resource.TestCheckTypeSetElemNestedAttrs(graphBetaLinuxPlatformScript.ResourceName+".test_006", "assignments.*", map[string]string{"type": "exclusionGroupAssignmentTarget", "group_id": "33333333-3333-3333-3333-333333333333"}),
					resource.TestCheckTypeSetElemNestedAttrs(graphBetaLinuxPlatformScript.ResourceName+".test_006", "assignments.*", map[string]string{"type": "allDevicesAssignmentTarget"}),
					resource.TestCheckTypeSetElemNestedAttrs(graphBetaLinuxPlatformScript.ResourceName+".test_006", "assignments.*", map[string]string{"type": "allLicensedUsersAssignmentTarget"}),
				),
			},
			{
				ResourceName:            graphBetaLinuxPlatformScript.ResourceName + ".test_006",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

func TestUnitResourceLinuxPlatformScript_07_AssignmentsLifecycle_MinimalToMaximal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, scriptMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer scriptMock.CleanupMockState()
	scriptMock.ExpectedWrites = []string{"validate_create/post_001_scenario_minimal.json", "validate_create/post_001_scenario_minimal.json"}

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		CheckDestroy:             func(_ *terraform.State) error { return scriptMock.CheckDestroyed() },
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("007_assignments_lifecycle_minimal_to_maximal_step_1.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaLinuxPlatformScript.ResourceName+".test_007").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaLinuxPlatformScript.ResourceName+".test_007").Key("name").HasValue("unit-test-linux-platform-script-007"),
					check.That(graphBetaLinuxPlatformScript.ResourceName+".test_007").Key("platforms").HasValue("linux"),
					check.That(graphBetaLinuxPlatformScript.ResourceName+".test_007").Key("technologies.0").HasValue("linuxMdm"),
					check.That(graphBetaLinuxPlatformScript.ResourceName+".test_007").Key("script_content").HasValue(minimalScriptContent),
					check.That(graphBetaLinuxPlatformScript.ResourceName+".test_007").Key("execution_context").HasValue("user"),
					check.That(graphBetaLinuxPlatformScript.ResourceName+".test_007").Key("execution_frequency").HasValue("15minutes"),
					check.That(graphBetaLinuxPlatformScript.ResourceName+".test_007").Key("execution_retries").HasValue("1"),
					check.That(graphBetaLinuxPlatformScript.ResourceName+".test_007").Key("role_scope_tag_ids.0").HasValue("0"),
					check.That(graphBetaLinuxPlatformScript.ResourceName+".test_007").Key("assignments.#").HasValue("1"),
					resource.TestCheckTypeSetElemNestedAttrs(graphBetaLinuxPlatformScript.ResourceName+".test_007", "assignments.*", map[string]string{"type": "groupAssignmentTarget", "group_id": "11111111-1111-1111-1111-111111111111", "filter_type": "none", "filter_id": "00000000-0000-0000-0000-000000000000"}),
				),
			},
			{
				Config: loadUnitTestTerraform("007_assignments_lifecycle_minimal_to_maximal_step_2.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaLinuxPlatformScript.ResourceName+".test_007").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaLinuxPlatformScript.ResourceName+".test_007").Key("name").HasValue("unit-test-linux-platform-script-007"),
					check.That(graphBetaLinuxPlatformScript.ResourceName+".test_007").Key("platforms").HasValue("linux"),
					check.That(graphBetaLinuxPlatformScript.ResourceName+".test_007").Key("technologies.0").HasValue("linuxMdm"),
					check.That(graphBetaLinuxPlatformScript.ResourceName+".test_007").Key("script_content").HasValue(minimalScriptContent),
					check.That(graphBetaLinuxPlatformScript.ResourceName+".test_007").Key("execution_context").HasValue("user"),
					check.That(graphBetaLinuxPlatformScript.ResourceName+".test_007").Key("execution_frequency").HasValue("15minutes"),
					check.That(graphBetaLinuxPlatformScript.ResourceName+".test_007").Key("execution_retries").HasValue("1"),
					check.That(graphBetaLinuxPlatformScript.ResourceName+".test_007").Key("role_scope_tag_ids.0").HasValue("0"),
					check.That(graphBetaLinuxPlatformScript.ResourceName+".test_007").Key("assignments.#").HasValue("5"),
					resource.TestCheckTypeSetElemNestedAttrs(graphBetaLinuxPlatformScript.ResourceName+".test_007", "assignments.*", map[string]string{"type": "groupAssignmentTarget", "group_id": "11111111-1111-1111-1111-111111111111", "filter_type": "none", "filter_id": "00000000-0000-0000-0000-000000000000"}),
					resource.TestCheckTypeSetElemNestedAttrs(graphBetaLinuxPlatformScript.ResourceName+".test_007", "assignments.*", map[string]string{"type": "groupAssignmentTarget", "group_id": "22222222-2222-2222-2222-222222222222", "filter_type": "include", "filter_id": "44444444-4444-4444-4444-444444444444"}),
					resource.TestCheckTypeSetElemNestedAttrs(graphBetaLinuxPlatformScript.ResourceName+".test_007", "assignments.*", map[string]string{"type": "exclusionGroupAssignmentTarget", "group_id": "33333333-3333-3333-3333-333333333333"}),
					resource.TestCheckTypeSetElemNestedAttrs(graphBetaLinuxPlatformScript.ResourceName+".test_007", "assignments.*", map[string]string{"type": "allDevicesAssignmentTarget"}),
					resource.TestCheckTypeSetElemNestedAttrs(graphBetaLinuxPlatformScript.ResourceName+".test_007", "assignments.*", map[string]string{"type": "allLicensedUsersAssignmentTarget"}),
				),
			},
			{
				ResourceName:            graphBetaLinuxPlatformScript.ResourceName + ".test_007",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

func TestUnitResourceLinuxPlatformScript_08_AssignmentsLifecycle_MaximalToMinimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, scriptMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer scriptMock.CleanupMockState()
	scriptMock.ExpectedWrites = []string{"validate_create/post_001_scenario_minimal.json", "validate_create/post_001_scenario_minimal.json", "validate_create/post_001_scenario_minimal.json"}

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		CheckDestroy:             func(_ *terraform.State) error { return scriptMock.CheckDestroyed() },
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("008_assignments_lifecycle_maximal_to_minimal_step_1.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaLinuxPlatformScript.ResourceName+".test_008").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaLinuxPlatformScript.ResourceName+".test_008").Key("name").HasValue("unit-test-linux-platform-script-008"),
					check.That(graphBetaLinuxPlatformScript.ResourceName+".test_008").Key("platforms").HasValue("linux"),
					check.That(graphBetaLinuxPlatformScript.ResourceName+".test_008").Key("technologies.0").HasValue("linuxMdm"),
					check.That(graphBetaLinuxPlatformScript.ResourceName+".test_008").Key("script_content").HasValue(minimalScriptContent),
					check.That(graphBetaLinuxPlatformScript.ResourceName+".test_008").Key("execution_context").HasValue("user"),
					check.That(graphBetaLinuxPlatformScript.ResourceName+".test_008").Key("execution_frequency").HasValue("15minutes"),
					check.That(graphBetaLinuxPlatformScript.ResourceName+".test_008").Key("execution_retries").HasValue("1"),
					check.That(graphBetaLinuxPlatformScript.ResourceName+".test_008").Key("role_scope_tag_ids.0").HasValue("0"),
					check.That(graphBetaLinuxPlatformScript.ResourceName+".test_008").Key("assignments.#").HasValue("5"),
					resource.TestCheckTypeSetElemNestedAttrs(graphBetaLinuxPlatformScript.ResourceName+".test_008", "assignments.*", map[string]string{"type": "groupAssignmentTarget", "group_id": "11111111-1111-1111-1111-111111111111", "filter_type": "none", "filter_id": "00000000-0000-0000-0000-000000000000"}),
					resource.TestCheckTypeSetElemNestedAttrs(graphBetaLinuxPlatformScript.ResourceName+".test_008", "assignments.*", map[string]string{"type": "groupAssignmentTarget", "group_id": "22222222-2222-2222-2222-222222222222", "filter_type": "include", "filter_id": "44444444-4444-4444-4444-444444444444"}),
					resource.TestCheckTypeSetElemNestedAttrs(graphBetaLinuxPlatformScript.ResourceName+".test_008", "assignments.*", map[string]string{"type": "exclusionGroupAssignmentTarget", "group_id": "33333333-3333-3333-3333-333333333333"}),
					resource.TestCheckTypeSetElemNestedAttrs(graphBetaLinuxPlatformScript.ResourceName+".test_008", "assignments.*", map[string]string{"type": "allDevicesAssignmentTarget"}),
					resource.TestCheckTypeSetElemNestedAttrs(graphBetaLinuxPlatformScript.ResourceName+".test_008", "assignments.*", map[string]string{"type": "allLicensedUsersAssignmentTarget"}),
				),
			},
			{
				Config: loadUnitTestTerraform("008_assignments_lifecycle_maximal_to_minimal_step_2.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaLinuxPlatformScript.ResourceName+".test_008").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaLinuxPlatformScript.ResourceName+".test_008").Key("name").HasValue("unit-test-linux-platform-script-008"),
					check.That(graphBetaLinuxPlatformScript.ResourceName+".test_008").Key("platforms").HasValue("linux"),
					check.That(graphBetaLinuxPlatformScript.ResourceName+".test_008").Key("technologies.0").HasValue("linuxMdm"),
					check.That(graphBetaLinuxPlatformScript.ResourceName+".test_008").Key("script_content").HasValue(minimalScriptContent),
					check.That(graphBetaLinuxPlatformScript.ResourceName+".test_008").Key("execution_context").HasValue("user"),
					check.That(graphBetaLinuxPlatformScript.ResourceName+".test_008").Key("execution_frequency").HasValue("15minutes"),
					check.That(graphBetaLinuxPlatformScript.ResourceName+".test_008").Key("execution_retries").HasValue("1"),
					check.That(graphBetaLinuxPlatformScript.ResourceName+".test_008").Key("role_scope_tag_ids.0").HasValue("0"),
					check.That(graphBetaLinuxPlatformScript.ResourceName+".test_008").Key("assignments.#").HasValue("1"),
					resource.TestCheckTypeSetElemNestedAttrs(graphBetaLinuxPlatformScript.ResourceName+".test_008", "assignments.*", map[string]string{"type": "groupAssignmentTarget", "group_id": "11111111-1111-1111-1111-111111111111", "filter_type": "none", "filter_id": "00000000-0000-0000-0000-000000000000"}),
				),
			},
			{
				Config: loadUnitTestTerraform("008_assignments_lifecycle_maximal_to_minimal_step_3.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaLinuxPlatformScript.ResourceName+".test_008").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaLinuxPlatformScript.ResourceName+".test_008").Key("name").HasValue("unit-test-linux-platform-script-008"),
					check.That(graphBetaLinuxPlatformScript.ResourceName+".test_008").Key("platforms").HasValue("linux"),
					check.That(graphBetaLinuxPlatformScript.ResourceName+".test_008").Key("technologies.0").HasValue("linuxMdm"),
					check.That(graphBetaLinuxPlatformScript.ResourceName+".test_008").Key("script_content").HasValue(minimalScriptContent),
					check.That(graphBetaLinuxPlatformScript.ResourceName+".test_008").Key("execution_context").HasValue("user"),
					check.That(graphBetaLinuxPlatformScript.ResourceName+".test_008").Key("execution_frequency").HasValue("15minutes"),
					check.That(graphBetaLinuxPlatformScript.ResourceName+".test_008").Key("execution_retries").HasValue("1"),
					check.That(graphBetaLinuxPlatformScript.ResourceName+".test_008").Key("role_scope_tag_ids.0").HasValue("0"),
					resource.TestCheckNoResourceAttr(graphBetaLinuxPlatformScript.ResourceName+".test_008", "assignments.#"),
				),
			},
			{
				ResourceName:            graphBetaLinuxPlatformScript.ResourceName + ".test_008",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

func TestUnitResourceLinuxPlatformScript_09_ErrorHandling(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, scriptMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer scriptMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("009_error_scenario.tf"),
				ExpectError: regexp.MustCompile("Invalid Linux Platform Script data"),
			},
		},
	})
}
