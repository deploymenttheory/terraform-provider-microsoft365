package graphBetaWindowsBiosConfigurationsAndOtherSettingsTemplate_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaWindowsBiosConfigurationsAndOtherSettingsTemplate "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/windows_bios_configurations_and_other_settings_template"
	biosMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/windows_bios_configurations_and_other_settings_template/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

const (
	minimalConfigurationFileContent = "W2NjdGtdCjtEZWxsIENvbW1hbmQgfCBDb25maWd1cmUgLSBtaW5pbWFsIHVuaXQgdGVzdCBwYWNrYWdlClNlY3VyZUJvb3Q9RW5hYmxlZApUcG09T24K"
	maximalConfigurationFileContent = "W2NjdGtdCjtEZWxsIENvbW1hbmQgfCBDb25maWd1cmUgLSBtYXhpbWFsIHVuaXQgdGVzdCBwYWNrYWdlCk51bUxvY2s9RW5hYmxlZApTZWN1cmVCb290PUVuYWJsZWQKVHBtPU9uClRwbUFjdGl2YXRpb249QWN0aXZhdGUKVXNiQm9vdD1EaXNhYmxlZApWaXJ0dWFsaXphdGlvbj1FbmFibGVkCldha2VPbkxhbj1MYW5Pbmx5Cg=="
)

func resourceName(label string) string {
	return graphBetaWindowsBiosConfigurationsAndOtherSettingsTemplate.ResourceName + "." + label
}

func setupMockEnvironment() (*mocks.Mocks, *biosMocks.WindowsBiosConfigurationsAndOtherSettingsTemplateMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	biosMock := &biosMocks.WindowsBiosConfigurationsAndOtherSettingsTemplateMock{}
	biosMock.RegisterMocks()
	return mockClient, biosMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *biosMocks.WindowsBiosConfigurationsAndOtherSettingsTemplateMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	biosMock := &biosMocks.WindowsBiosConfigurationsAndOtherSettingsTemplateMock{}
	biosMock.RegisterErrorMocks()
	return mockClient, biosMock
}

func loadUnitTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/unit/" + filename)
	if err != nil {
		panic("failed to load unit test config " + filename + ": " + err.Error())
	}
	return config
}

func TestUnitResourceWindowsBiosConfigurationsAndOtherSettingsTemplate_01_Scenario_Minimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, biosMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer biosMock.CleanupMockState()

	name := resourceName("test_001")

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("001_scenario_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(name).Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(name).Key("display_name").HasValue("unit-test-windows-bios-configuration-001-minimal"),
					check.That(name).Key("file_name").HasValue("test-bios-001.cctk"),
					check.That(name).Key("configuration_file_content").HasValue(minimalConfigurationFileContent),
					check.That(name).Key("hardware_configuration_format").HasValue("dell"),
					check.That(name).Key("per_device_password_disabled").HasValue("false"),
					check.That(name).Key("role_scope_tag_ids.#").HasValue("1"),
					check.That(name).Key("role_scope_tag_ids.0").HasValue("0"),
					check.That(name).Key("version").Exists(),
					check.That(name).Key("created_date_time").Exists(),
					check.That(name).Key("last_modified_date_time").Exists(),
					check.That(name).Key("assignments").DoesNotExist(),
				),
			},
			{
				ResourceName:      name,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestUnitResourceWindowsBiosConfigurationsAndOtherSettingsTemplate_02_Scenario_Maximal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, biosMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer biosMock.CleanupMockState()

	name := resourceName("test_002")

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("002_scenario_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(name).Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(name).Key("display_name").HasValue("unit-test-windows-bios-configuration-002-maximal"),
					check.That(name).Key("description").HasValue("Maximal BIOS configuration template"),
					check.That(name).Key("file_name").HasValue("test-bios-002.cctk"),
					check.That(name).Key("configuration_file_content").HasValue(maximalConfigurationFileContent),
					check.That(name).Key("hardware_configuration_format").HasValue("dell"),
					check.That(name).Key("per_device_password_disabled").HasValue("true"),
					check.That(name).Key("role_scope_tag_ids.#").HasValue("2"),
				),
			},
			{
				ResourceName:      name,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestUnitResourceWindowsBiosConfigurationsAndOtherSettingsTemplate_03_Lifecycle_MinimalToMaximal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, biosMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer biosMock.CleanupMockState()

	name := resourceName("test_003")

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("003_lifecycle_minimal_to_maximal_step_1.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(name).Key("display_name").HasValue("unit-test-windows-bios-configuration-003-lifecycle"),
					check.That(name).Key("configuration_file_content").HasValue(minimalConfigurationFileContent),
					check.That(name).Key("per_device_password_disabled").HasValue("false"),
					check.That(name).Key("role_scope_tag_ids.#").HasValue("1"),
				),
			},
			{
				Config: loadUnitTestTerraform("003_lifecycle_minimal_to_maximal_step_2.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(name).Key("description").HasValue("Promoted to a maximal BIOS configuration template"),
					check.That(name).Key("configuration_file_content").HasValue(maximalConfigurationFileContent),
					check.That(name).Key("hardware_configuration_format").HasValue("dell"),
					check.That(name).Key("per_device_password_disabled").HasValue("true"),
					check.That(name).Key("role_scope_tag_ids.#").HasValue("2"),
				),
			},
		},
	})
}

func TestUnitResourceWindowsBiosConfigurationsAndOtherSettingsTemplate_04_Lifecycle_MaximalToMinimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, biosMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer biosMock.CleanupMockState()

	name := resourceName("test_004")

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("004_lifecycle_maximal_to_minimal_step_1.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(name).Key("description").HasValue("Starts as a maximal BIOS configuration template"),
					check.That(name).Key("configuration_file_content").HasValue(maximalConfigurationFileContent),
					check.That(name).Key("per_device_password_disabled").HasValue("true"),
					check.That(name).Key("role_scope_tag_ids.#").HasValue("2"),
				),
			},
			{
				Config: loadUnitTestTerraform("004_lifecycle_maximal_to_minimal_step_2.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(name).Key("description").HasValue("Reduced to a minimal BIOS configuration template"),
					check.That(name).Key("configuration_file_content").HasValue(minimalConfigurationFileContent),
					check.That(name).Key("per_device_password_disabled").HasValue("false"),
					check.That(name).Key("role_scope_tag_ids.#").HasValue("1"),
					check.That(name).Key("role_scope_tag_ids.0").HasValue("0"),
				),
			},
		},
	})
}

func TestUnitResourceWindowsBiosConfigurationsAndOtherSettingsTemplate_05_AssignmentsMinimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, biosMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer biosMock.CleanupMockState()

	name := resourceName("test_005")

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("005_assignments_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(name).Key("assignments.#").HasValue("1"),
					check.That(name).Key("assignments.0.type").HasValue("groupAssignmentTarget"),
					check.That(name).Key("assignments.0.group_id").HasValue("11111111-1111-1111-1111-111111111111"),
					check.That(name).Key("assignments.0.filter_type").HasValue("none"),
					check.That(name).Key("assignments.0.filter_id").HasValue("00000000-0000-0000-0000-000000000000"),
				),
			},
			{
				ResourceName:      name,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestUnitResourceWindowsBiosConfigurationsAndOtherSettingsTemplate_06_AssignmentsMaximal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, biosMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer biosMock.CleanupMockState()

	name := resourceName("test_006")

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("006_assignments_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(name).Key("assignments.#").HasValue("5"),
					resource.TestCheckTypeSetElemNestedAttrs(name, "assignments.*", map[string]string{
						"type":        "groupAssignmentTarget",
						"group_id":    "11111111-1111-1111-1111-111111111111",
						"filter_id":   "44444444-4444-4444-4444-444444444444",
						"filter_type": "include",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(name, "assignments.*", map[string]string{
						"type":        "groupAssignmentTarget",
						"group_id":    "22222222-2222-2222-2222-222222222222",
						"filter_id":   "55555555-5555-5555-5555-555555555555",
						"filter_type": "exclude",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(name, "assignments.*", map[string]string{
						"type":        "groupAssignmentTarget",
						"group_id":    "33333333-3333-3333-3333-333333333333",
						"filter_id":   "00000000-0000-0000-0000-000000000000",
						"filter_type": "none",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(name, "assignments.*", map[string]string{
						"type":     "exclusionGroupAssignmentTarget",
						"group_id": "66666666-6666-6666-6666-666666666666",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(name, "assignments.*", map[string]string{
						"type":     "exclusionGroupAssignmentTarget",
						"group_id": "77777777-7777-7777-7777-777777777777",
					}),
				),
			},
			{
				ResourceName:      name,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestUnitResourceWindowsBiosConfigurationsAndOtherSettingsTemplate_07_AssignmentsLifecycle_MinimalToMaximal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, biosMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer biosMock.CleanupMockState()

	name := resourceName("test_007")

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("007_assignments_lifecycle_minimal_to_maximal_step_1.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(name).Key("assignments.#").HasValue("1"),
				),
			},
			{
				Config: loadUnitTestTerraform("007_assignments_lifecycle_minimal_to_maximal_step_2.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(name).Key("assignments.#").HasValue("5"),
					resource.TestCheckTypeSetElemNestedAttrs(name, "assignments.*", map[string]string{
						"type":        "groupAssignmentTarget",
						"group_id":    "11111111-1111-1111-1111-111111111111",
						"filter_id":   "44444444-4444-4444-4444-444444444444",
						"filter_type": "include",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(name, "assignments.*", map[string]string{
						"type":     "exclusionGroupAssignmentTarget",
						"group_id": "77777777-7777-7777-7777-777777777777",
					}),
				),
			},
		},
	})
}

func TestUnitResourceWindowsBiosConfigurationsAndOtherSettingsTemplate_08_AssignmentsLifecycle_MaximalToMinimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, biosMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer biosMock.CleanupMockState()

	name := resourceName("test_008")

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("008_assignments_lifecycle_maximal_to_minimal_step_1.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(name).Key("assignments.#").HasValue("5"),
					resource.TestCheckTypeSetElemNestedAttrs(name, "assignments.*", map[string]string{
						"type":        "groupAssignmentTarget",
						"group_id":    "11111111-1111-1111-1111-111111111111",
						"filter_id":   "44444444-4444-4444-4444-444444444444",
						"filter_type": "include",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(name, "assignments.*", map[string]string{
						"type":     "exclusionGroupAssignmentTarget",
						"group_id": "66666666-6666-6666-6666-666666666666",
					}),
				),
			},
			{
				Config: loadUnitTestTerraform("008_assignments_lifecycle_maximal_to_minimal_step_2.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(name).Key("assignments.#").HasValue("1"),
					check.That(name).Key("assignments.0.type").HasValue("groupAssignmentTarget"),
					check.That(name).Key("assignments.0.group_id").HasValue("11111111-1111-1111-1111-111111111111"),
				),
			},
		},
	})
}

func TestUnitResourceWindowsBiosConfigurationsAndOtherSettingsTemplate_09_ErrorHandling(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, biosMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer biosMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("009_error_scenario.tf"),
				ExpectError: regexp.MustCompile("Invalid Windows BIOS configuration template data"),
			},
		},
	})
}

func TestUnitResourceWindowsBiosConfigurationsAndOtherSettingsTemplate_10_Validation_InvalidBase64(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, biosMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer biosMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("010_validation_invalid_base64.tf"),
				ExpectError: regexp.MustCompile("must be a base64 encoded string"),
			},
		},
	})
}

func TestUnitResourceWindowsBiosConfigurationsAndOtherSettingsTemplate_11_Validation_InvalidFormat(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, biosMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer biosMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("011_validation_invalid_format.tf"),
				ExpectError: regexp.MustCompile("Invalid Attribute Value Match"),
			},
		},
	})
}

func TestUnitResourceWindowsBiosConfigurationsAndOtherSettingsTemplate_12_Validation_InvalidAssignmentTarget(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, biosMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer biosMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("012_validation_invalid_assignment_target.tf"),
				ExpectError: regexp.MustCompile("Invalid Attribute Value Match"),
			},
		},
	})
}
