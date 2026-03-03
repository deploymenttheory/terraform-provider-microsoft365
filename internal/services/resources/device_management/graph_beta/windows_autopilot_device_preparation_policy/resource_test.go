package graphBetaWindowsAutopilotDevicePreparationPolicy_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	policyMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/windows_autopilot_device_preparation_policy/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

func setupMockEnvironment() (*mocks.Mocks, *policyMocks.WindowsAutopilotDevicePreparationPolicyMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	policyMock := &policyMocks.WindowsAutopilotDevicePreparationPolicyMock{}
	policyMock.RegisterMocks()
	return mockClient, policyMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *policyMocks.WindowsAutopilotDevicePreparationPolicyMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	policyMock := &policyMocks.WindowsAutopilotDevicePreparationPolicyMock{}
	policyMock.RegisterErrorMocks()
	return mockClient, policyMock
}

func loadUnitTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/unit/" + filename)
	if err != nil {
		panic("failed to load unit test config " + filename + ": " + err.Error())
	}
	return config
}

// Test 01: Automatic mode minimal configuration
func TestUnitResourceWindowsAutopilotDevicePreparationPolicy_01_AutomaticMinimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, policyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer policyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("001_scenario_automatic_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".auto_minimal").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".auto_minimal").Key("name").HasValue("unit-test-autopilot-dpp-auto-minimal"),
					check.That(resourceType+".auto_minimal").Key("deployment_settings.deployment_type").HasValue("enrollment_autopilot_dpp_deploymenttype_1"),
					check.That(resourceType+".auto_minimal").Key("allowed_apps.#").HasValue("1"),
					check.That(resourceType+".auto_minimal").Key("allowed_apps.0.app_id").HasValue("00000000-0000-0000-0000-000000000001"),
					check.That(resourceType+".auto_minimal").Key("allowed_apps.0.app_type").HasValue("winGetApp"),
				),
			},
		},
	})
}

// Test 02: Automatic mode maximal configuration
func TestUnitResourceWindowsAutopilotDevicePreparationPolicy_02_AutomaticMaximal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, policyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer policyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("002_scenario_automatic_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".auto_maximal").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".auto_maximal").Key("name").HasValue("unit-test-autopilot-dpp-auto-maximal"),
					check.That(resourceType+".auto_maximal").Key("deployment_settings.deployment_type").HasValue("enrollment_autopilot_dpp_deploymenttype_1"),
					check.That(resourceType+".auto_maximal").Key("allowed_apps.#").HasValue("3"),
					check.That(resourceType+".auto_maximal").Key("allowed_scripts.#").HasValue("2"),
				),
			},
		},
	})
}

// Test 03: User-driven mode minimal configuration
func TestUnitResourceWindowsAutopilotDevicePreparationPolicy_03_UserDrivenMinimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, policyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer policyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("003_scenario_user_driven_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".ud_minimal").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".ud_minimal").Key("name").HasValue("unit-test-autopilot-dpp-ud-minimal"),
					check.That(resourceType+".ud_minimal").Key("deployment_settings.deployment_type").HasValue("enrollment_autopilot_dpp_deploymenttype_0"),
					check.That(resourceType+".ud_minimal").Key("device_security_group").HasValue("00000000-0000-0000-0000-000000000001"),
					check.That(resourceType+".ud_minimal").Key("deployment_settings.deployment_mode").HasValue("enrollment_autopilot_dpp_deploymentmode_0"),
					check.That(resourceType+".ud_minimal").Key("oobe_settings.timeout_in_minutes").HasValue("60"),
					check.That(resourceType+".ud_minimal").Key("allowed_apps.#").HasValue("1"),
				),
			},
		},
	})
}

// Test 04: User-driven mode maximal configuration
func TestUnitResourceWindowsAutopilotDevicePreparationPolicy_04_UserDrivenMaximal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, policyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer policyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("004_scenario_user_driven_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".ud_maximal").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".ud_maximal").Key("name").HasValue("unit-test-autopilot-dpp-ud-maximal"),
					check.That(resourceType+".ud_maximal").Key("deployment_settings.deployment_type").HasValue("enrollment_autopilot_dpp_deploymenttype_0"),
					check.That(resourceType+".ud_maximal").Key("deployment_settings.deployment_mode").HasValue("enrollment_autopilot_dpp_deploymentmode_1"),
					check.That(resourceType+".ud_maximal").Key("oobe_settings.allow_skip").HasValue("true"),
					check.That(resourceType+".ud_maximal").Key("allowed_apps.#").HasValue("3"),
					check.That(resourceType+".ud_maximal").Key("allowed_scripts.#").HasValue("2"),
				),
			},
		},
	})
}

// Test 05: User-driven mode with minimal assignments
func TestUnitResourceWindowsAutopilotDevicePreparationPolicy_05_UserDrivenMinimalWithMinimalAssignments(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, policyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer policyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("005_scenario_user_driven_minimal_with_minimal_assignments.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".ud_min_assign").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".ud_min_assign").Key("name").HasValue("unit-test-autopilot-dpp-ud-min-assign"),
					check.That(resourceType+".ud_min_assign").Key("deployment_settings.deployment_type").HasValue("enrollment_autopilot_dpp_deploymenttype_0"),
					check.That(resourceType+".ud_min_assign").Key("assignments.#").HasValue("1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceType+".ud_min_assign", "assignments.*", map[string]string{
						"type":     "groupAssignmentTarget",
						"group_id": "00000000-0000-0000-0000-000000000003",
					}),
				),
			},
		},
	})
}

// Test 06: User-driven mode with maximal assignments
func TestUnitResourceWindowsAutopilotDevicePreparationPolicy_06_UserDrivenMaximalWithMaximalAssignments(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, policyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer policyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("006_scenario_user_driven_minimal_with_maximal_assignments.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".ud_max_assign").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".ud_max_assign").Key("name").HasValue("unit-test-autopilot-dpp-ud-max-assign"),
					check.That(resourceType+".ud_max_assign").Key("deployment_settings.deployment_type").HasValue("enrollment_autopilot_dpp_deploymenttype_0"),
					check.That(resourceType+".ud_max_assign").Key("assignments.#").HasValue("4"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceType+".ud_max_assign", "assignments.*", map[string]string{
						"type": "allLicensedUsersAssignmentTarget",
					}),
				),
			},
		},
	})
}

// Test 07: Error handling - API rejection
func TestUnitResourceWindowsAutopilotDevicePreparationPolicy_07_ErrorHandling(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, policyMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer policyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("003_scenario_user_driven_minimal.tf"),
				ExpectError: regexp.MustCompile("Invalid Group ID in include_group_ids"),
			},
		},
	})
}
