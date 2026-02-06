package graphBetaMacOSSoftwareUpdateConfiguration_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	softwareUpdateConfigurationMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/macos_software_update_configuration/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

func setupMockEnvironment() (*mocks.Mocks, *softwareUpdateConfigurationMocks.MacOSSoftwareUpdateConfigurationMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	softwareUpdateConfigurationMock := &softwareUpdateConfigurationMocks.MacOSSoftwareUpdateConfigurationMock{}
	softwareUpdateConfigurationMock.RegisterMocks()
	return mockClient, softwareUpdateConfigurationMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *softwareUpdateConfigurationMocks.MacOSSoftwareUpdateConfigurationMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	softwareUpdateConfigurationMock := &softwareUpdateConfigurationMocks.MacOSSoftwareUpdateConfigurationMock{}
	softwareUpdateConfigurationMock.RegisterErrorMocks()
	return mockClient, softwareUpdateConfigurationMock
}

func loadUnitTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/unit/" + filename)
	if err != nil {
		panic("failed to load unit test config " + filename + ": " + err.Error())
	}
	return config
}

func TestUnitResourceMacOSSoftwareUpdateConfiguration_01_CreateMinimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, softwareUpdateConfigurationMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer softwareUpdateConfigurationMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_01_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_01_minimal").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_01_minimal").Key("display_name").HasValue("Test 01: Minimal macOS Software Update Configuration"),
					check.That(resourceType+".test_01_minimal").Key("update_schedule_type").HasValue("alwaysUpdate"),
					check.That(resourceType+".test_01_minimal").Key("critical_update_behavior").HasValue("installASAP"),
					check.That(resourceType+".test_01_minimal").Key("config_data_update_behavior").HasValue("installASAP"),
					check.That(resourceType+".test_01_minimal").Key("firmware_update_behavior").HasValue("installASAP"),
					check.That(resourceType+".test_01_minimal").Key("all_other_update_behavior").HasValue("installASAP"),
					check.That(resourceType+".test_01_minimal").Key("role_scope_tag_ids.#").HasValue("1"),
					check.That(resourceType+".test_01_minimal").Key("role_scope_tag_ids.0").HasValue("0"),
				),
			},
		},
	})
}

func TestUnitResourceMacOSSoftwareUpdateConfiguration_02_CreateMaximal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, softwareUpdateConfigurationMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer softwareUpdateConfigurationMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_02_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_02_maximal").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_02_maximal").Key("display_name").HasValue("Test 02: Maximal macOS Software Update Configuration"),
					check.That(resourceType+".test_02_maximal").Key("description").HasValue("Maximal software update configuration for testing with all features"),
					check.That(resourceType+".test_02_maximal").Key("update_schedule_type").HasValue("updateDuringTimeWindows"),
					check.That(resourceType+".test_02_maximal").Key("critical_update_behavior").HasValue("installASAP"),
					check.That(resourceType+".test_02_maximal").Key("max_user_deferrals_count").HasValue("5"),
					check.That(resourceType+".test_02_maximal").Key("priority").HasValue("high"),
					check.That(resourceType+".test_02_maximal").Key("role_scope_tag_ids.#").HasValue("2"),
					check.That(resourceType+".test_02_maximal").Key("role_scope_tag_ids.0").HasValue("0"),
					check.That(resourceType+".test_02_maximal").Key("role_scope_tag_ids.1").HasValue("1"),
					check.That(resourceType+".test_02_maximal").Key("custom_update_time_windows.#").HasValue("2"),
				),
			},
		},
	})
}

func TestUnitResourceMacOSSoftwareUpdateConfiguration_03_UpdateMinimalToMaximal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, softwareUpdateConfigurationMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer softwareUpdateConfigurationMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_03_minimal_step.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_03_progression").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_03_progression").Key("display_name").HasValue("Test 03: Progression macOS Software Update Configuration"),
					check.That(resourceType+".test_03_progression").Key("update_schedule_type").HasValue("alwaysUpdate"),
					check.That(resourceType+".test_03_progression").Key("role_scope_tag_ids.#").HasValue("1"),
				),
			},
			{
				Config: loadUnitTestTerraform("resource_03_intermediate_step.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_03_progression").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_03_progression").Key("display_name").HasValue("Test 03: Progression macOS Software Update Configuration"),
					check.That(resourceType+".test_03_progression").Key("update_schedule_type").HasValue("updateDuringTimeWindows"),
					check.That(resourceType+".test_03_progression").Key("priority").HasValue("low"),
				),
			},
			{
				Config: loadUnitTestTerraform("resource_03_maximal_step.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_03_progression").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_03_progression").Key("display_name").HasValue("Test 03: Progression macOS Software Update Configuration"),
					check.That(resourceType+".test_03_progression").Key("update_schedule_type").HasValue("updateDuringTimeWindows"),
					check.That(resourceType+".test_03_progression").Key("description").HasValue("Maximal software update configuration with all features"),
					check.That(resourceType+".test_03_progression").Key("priority").HasValue("high"),
					check.That(resourceType+".test_03_progression").Key("role_scope_tag_ids.#").HasValue("2"),
				),
			},
		},
	})
}

func TestUnitResourceMacOSSoftwareUpdateConfiguration_04_UpdateMaximalToMinimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, softwareUpdateConfigurationMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer softwareUpdateConfigurationMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_04_maximal_step.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_04_regression").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_04_regression").Key("display_name").HasValue("Test 04: Regression macOS Software Update Configuration"),
					check.That(resourceType+".test_04_regression").Key("update_schedule_type").HasValue("updateDuringTimeWindows"),
					check.That(resourceType+".test_04_regression").Key("priority").HasValue("high"),
				),
			},
			{
				Config: loadUnitTestTerraform("resource_04_intermediate_step.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_04_regression").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_04_regression").Key("display_name").HasValue("Test 04: Regression macOS Software Update Configuration"),
					check.That(resourceType+".test_04_regression").Key("update_schedule_type").HasValue("updateDuringTimeWindows"),
				),
			},
			{
				Config: loadUnitTestTerraform("resource_04_minimal_step.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_04_regression").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_04_regression").Key("display_name").HasValue("Test 04: Regression macOS Software Update Configuration"),
					check.That(resourceType+".test_04_regression").Key("update_schedule_type").HasValue("alwaysUpdate"),
					check.That(resourceType+".test_04_regression").Key("role_scope_tag_ids.#").HasValue("1"),
				),
			},
		},
	})
}

func TestUnitResourceMacOSSoftwareUpdateConfiguration_05_MinimalAssignments(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, softwareUpdateConfigurationMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer softwareUpdateConfigurationMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_05_minimal_assignments.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_05_min_assignments").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_05_min_assignments").Key("display_name").HasValue("Test 05: Minimal Assignments macOS Software Update Configuration"),
					check.That(resourceType+".test_05_min_assignments").Key("assignments.#").HasValue("1"),
				),
			},
		},
	})
}

func TestUnitResourceMacOSSoftwareUpdateConfiguration_06_MaximalAssignments(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, softwareUpdateConfigurationMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer softwareUpdateConfigurationMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_06_maximal_assignments.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_06_max_assignments").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_06_max_assignments").Key("display_name").HasValue("Test 06: Maximal Assignments macOS Software Update Configuration"),
					check.That(resourceType+".test_06_max_assignments").Key("assignments.#").HasValue("4"),
				),
			},
		},
	})
}

func TestUnitResourceMacOSSoftwareUpdateConfiguration_07_MinimalToMaximalAssignments(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, softwareUpdateConfigurationMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer softwareUpdateConfigurationMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_07_minimal_assignments_step.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_07_assignments_progression").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_07_assignments_progression").Key("display_name").HasValue("Test 07: Assignments Progression macOS Software Update Configuration"),
					check.That(resourceType+".test_07_assignments_progression").Key("assignments.#").HasValue("1"),
				),
			},
			{
				Config: loadUnitTestTerraform("resource_07_maximal_assignments_step.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_07_assignments_progression").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_07_assignments_progression").Key("display_name").HasValue("Test 07: Assignments Progression macOS Software Update Configuration"),
					check.That(resourceType+".test_07_assignments_progression").Key("assignments.#").HasValue("4"),
				),
			},
		},
	})
}

func TestUnitResourceMacOSSoftwareUpdateConfiguration_08_MaximalToMinimalAssignments(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, softwareUpdateConfigurationMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer softwareUpdateConfigurationMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_08_maximal_assignments_step.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_08_assignments_regression").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_08_assignments_regression").Key("display_name").HasValue("Test 08: Assignments Regression macOS Software Update Configuration"),
					check.That(resourceType+".test_08_assignments_regression").Key("assignments.#").HasValue("4"),
				),
			},
			{
				Config: loadUnitTestTerraform("resource_08_minimal_assignments_step.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_08_assignments_regression").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_08_assignments_regression").Key("display_name").HasValue("Test 08: Assignments Regression macOS Software Update Configuration"),
					check.That(resourceType+".test_08_assignments_regression").Key("assignments.#").HasValue("1"),
				),
			},
		},
	})
}

func TestUnitResourceMacOSSoftwareUpdateConfiguration_09_Import(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, softwareUpdateConfigurationMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer softwareUpdateConfigurationMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_01_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_01_minimal").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
				),
			},
			{
				ResourceName:      resourceType + ".test_01_minimal",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestUnitResourceMacOSSoftwareUpdateConfiguration_10_Error(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, softwareUpdateConfigurationMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer softwareUpdateConfigurationMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("resource_01_minimal.tf"),
				ExpectError: regexp.MustCompile("Validation error: Invalid display name"),
			},
		},
	})
}
