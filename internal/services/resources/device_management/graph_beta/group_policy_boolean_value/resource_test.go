package graphBetaGroupPolicyBooleanValue_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaGroupPolicyBooleanValue "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/group_policy_boolean_value"
	groupPolicyBooleanValueMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/group_policy_boolean_value/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

// setupMockEnvironment sets up the mock environment using centralized mocks
func setupMockEnvironment() (*mocks.Mocks, *groupPolicyBooleanValueMocks.GroupPolicyBooleanValueMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	groupPolicyBooleanValueMock := &groupPolicyBooleanValueMocks.GroupPolicyBooleanValueMock{}
	groupPolicyBooleanValueMock.RegisterMocks()

	// Setup mock configurations and definitions for testing
	setupTestData(groupPolicyBooleanValueMock)

	return mockClient, groupPolicyBooleanValueMock
}

// setupErrorMockEnvironment sets up the mock environment for error testing
func setupErrorMockEnvironment() (*mocks.Mocks, *groupPolicyBooleanValueMocks.GroupPolicyBooleanValueMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	groupPolicyBooleanValueMock := &groupPolicyBooleanValueMocks.GroupPolicyBooleanValueMock{}
	groupPolicyBooleanValueMock.RegisterErrorMocks()

	return mockClient, groupPolicyBooleanValueMock
}

// setupTestData creates mock configurations and definitions for testing
func setupTestData(mock *groupPolicyBooleanValueMocks.GroupPolicyBooleanValueMock) {
	// Setup mock configurations
	mock.SetupMockConfiguration("00000000-0000-0000-0000-000000000001", "Test Configuration 001")
	mock.SetupMockConfiguration("00000000-0000-0000-0000-000000000002", "Test Configuration 002")
	mock.SetupMockConfiguration("00000000-0000-0000-0000-000000000003", "Test Configuration 003")
	mock.SetupMockConfiguration("00000000-0000-0000-0000-000000000004", "Test Configuration 004")

	// Setup mock definitions with presentations
	mock.SetupMockDefinition(
		"def-template-001",
		"Test Policy Minimal",
		"machine",
		"\\Test\\Category",
		[]map[string]any{
			{
				"@odata.type": "#microsoft.graph.groupPolicyPresentationCheckBox",
				"id":          "presentation-001",
				"label":       "Enable feature",
			},
		},
	)

	mock.SetupMockDefinition(
		"def-template-002",
		"Test Policy Maximal",
		"user",
		"\\Test\\Category\\Maximal",
		[]map[string]any{
			{
				"@odata.type": "#microsoft.graph.groupPolicyPresentationCheckBox",
				"id":          "presentation-002-1",
				"label":       "Enable feature 1",
			},
			{
				"@odata.type": "#microsoft.graph.groupPolicyPresentationCheckBox",
				"id":          "presentation-002-2",
				"label":       "Enable feature 2",
			},
			{
				"@odata.type": "#microsoft.graph.groupPolicyPresentationCheckBox",
				"id":          "presentation-002-3",
				"label":       "Enable feature 3",
			},
		},
	)

	mock.SetupMockDefinition(
		"def-template-003",
		"Test Policy Lifecycle",
		"machine",
		"\\Test\\Category\\Lifecycle",
		[]map[string]any{
			{
				"@odata.type": "#microsoft.graph.groupPolicyPresentationCheckBox",
				"id":          "presentation-003-1",
				"label":       "Feature 1",
			},
			{
				"@odata.type": "#microsoft.graph.groupPolicyPresentationCheckBox",
				"id":          "presentation-003-2",
				"label":       "Feature 2",
			},
			{
				"@odata.type": "#microsoft.graph.groupPolicyPresentationCheckBox",
				"id":          "presentation-003-3",
				"label":       "Feature 3",
			},
		},
	)

	mock.SetupMockDefinition(
		"def-template-004",
		"Test Policy Downgrade",
		"user",
		"\\Test\\Category\\Downgrade",
		[]map[string]any{
			{
				"@odata.type": "#microsoft.graph.groupPolicyPresentationCheckBox",
				"id":          "presentation-004-1",
				"label":       "Feature 1",
			},
			{
				"@odata.type": "#microsoft.graph.groupPolicyPresentationCheckBox",
				"id":          "presentation-004-2",
				"label":       "Feature 2",
			},
			{
				"@odata.type": "#microsoft.graph.groupPolicyPresentationCheckBox",
				"id":          "presentation-004-3",
				"label":       "Feature 3",
			},
		},
	)
}

// Helper function to load test configs from unit directory
func loadUnitTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/unit/" + filename)
	if err != nil {
		panic("failed to load unit test config " + filename + ": " + err.Error())
	}
	return config
}

// Test 001: Scenario 1 - Minimal configuration
func TestGroupPolicyBooleanValueResource_001_Scenario_Minimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, groupPolicyBooleanValueMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer groupPolicyBooleanValueMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("001_scenario_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaGroupPolicyBooleanValue.ResourceName+".test_001").Key("id").Exists(),
					check.That(graphBetaGroupPolicyBooleanValue.ResourceName+".test_001").Key("group_policy_configuration_id").HasValue("00000000-0000-0000-0000-000000000001"),
					check.That(graphBetaGroupPolicyBooleanValue.ResourceName+".test_001").Key("policy_name").HasValue("Test Policy Minimal"),
					check.That(graphBetaGroupPolicyBooleanValue.ResourceName+".test_001").Key("class_type").HasValue("machine"),
					check.That(graphBetaGroupPolicyBooleanValue.ResourceName+".test_001").Key("category_path").HasValue("\\Test\\Category"),
					check.That(graphBetaGroupPolicyBooleanValue.ResourceName+".test_001").Key("enabled").HasValue("true"),
					check.That(graphBetaGroupPolicyBooleanValue.ResourceName+".test_001").Key("values.#").HasValue("1"),
					check.That(graphBetaGroupPolicyBooleanValue.ResourceName+".test_001").Key("values.0.value").HasValue("true"),
					check.That(graphBetaGroupPolicyBooleanValue.ResourceName+".test_001").Key("group_policy_definition_value_id").Exists(),
					check.That(graphBetaGroupPolicyBooleanValue.ResourceName+".test_001").Key("created_date_time").Exists(),
					check.That(graphBetaGroupPolicyBooleanValue.ResourceName+".test_001").Key("last_modified_date_time").Exists(),
				),
			},
			{
				ResourceName:      graphBetaGroupPolicyBooleanValue.ResourceName + ".test_001",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Test 002: Scenario 2 - Maximal configuration
func TestGroupPolicyBooleanValueResource_002_Scenario_Maximal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, groupPolicyBooleanValueMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer groupPolicyBooleanValueMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("002_scenario_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaGroupPolicyBooleanValue.ResourceName+".test_002").Key("id").Exists(),
					check.That(graphBetaGroupPolicyBooleanValue.ResourceName+".test_002").Key("group_policy_configuration_id").HasValue("00000000-0000-0000-0000-000000000002"),
					check.That(graphBetaGroupPolicyBooleanValue.ResourceName+".test_002").Key("policy_name").HasValue("Test Policy Maximal"),
					check.That(graphBetaGroupPolicyBooleanValue.ResourceName+".test_002").Key("class_type").HasValue("user"),
					check.That(graphBetaGroupPolicyBooleanValue.ResourceName+".test_002").Key("category_path").HasValue("\\Test\\Category\\Maximal"),
					check.That(graphBetaGroupPolicyBooleanValue.ResourceName+".test_002").Key("enabled").HasValue("true"),
					check.That(graphBetaGroupPolicyBooleanValue.ResourceName+".test_002").Key("values.#").HasValue("3"),
					check.That(graphBetaGroupPolicyBooleanValue.ResourceName+".test_002").Key("values.0.value").HasValue("true"),
					check.That(graphBetaGroupPolicyBooleanValue.ResourceName+".test_002").Key("values.1.value").HasValue("false"),
					check.That(graphBetaGroupPolicyBooleanValue.ResourceName+".test_002").Key("values.2.value").HasValue("true"),
				),
			},
			{
				ResourceName:      graphBetaGroupPolicyBooleanValue.ResourceName + ".test_002",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Test 003: Scenario 3 - Lifecycle from minimal to maximal
func TestGroupPolicyBooleanValueResource_003_Lifecycle_MinimalToMaximal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, groupPolicyBooleanValueMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer groupPolicyBooleanValueMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("003_lifecycle_minimal_to_maximal_step_1.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaGroupPolicyBooleanValue.ResourceName+".test_003").Key("id").Exists(),
					check.That(graphBetaGroupPolicyBooleanValue.ResourceName+".test_003").Key("enabled").HasValue("false"),
					check.That(graphBetaGroupPolicyBooleanValue.ResourceName+".test_003").Key("values.#").HasValue("1"),
					check.That(graphBetaGroupPolicyBooleanValue.ResourceName+".test_003").Key("values.0.value").HasValue("false"),
				),
			},
			{
				Config: loadUnitTestTerraform("003_lifecycle_minimal_to_maximal_step_2.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaGroupPolicyBooleanValue.ResourceName+".test_003").Key("id").Exists(),
					check.That(graphBetaGroupPolicyBooleanValue.ResourceName+".test_003").Key("enabled").HasValue("true"),
					check.That(graphBetaGroupPolicyBooleanValue.ResourceName+".test_003").Key("values.#").HasValue("3"),
					check.That(graphBetaGroupPolicyBooleanValue.ResourceName+".test_003").Key("values.0.value").HasValue("true"),
					check.That(graphBetaGroupPolicyBooleanValue.ResourceName+".test_003").Key("values.1.value").HasValue("false"),
					check.That(graphBetaGroupPolicyBooleanValue.ResourceName+".test_003").Key("values.2.value").HasValue("true"),
				),
			},
			{
				ResourceName:      graphBetaGroupPolicyBooleanValue.ResourceName + ".test_003",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Test 004: Scenario 4 - Lifecycle from maximal to minimal
func TestGroupPolicyBooleanValueResource_004_Lifecycle_MaximalToMinimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, groupPolicyBooleanValueMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer groupPolicyBooleanValueMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("004_lifecycle_maximal_to_minimal_step_1.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaGroupPolicyBooleanValue.ResourceName+".test_004").Key("id").Exists(),
					check.That(graphBetaGroupPolicyBooleanValue.ResourceName+".test_004").Key("enabled").HasValue("true"),
					check.That(graphBetaGroupPolicyBooleanValue.ResourceName+".test_004").Key("values.#").HasValue("3"),
				),
			},
			{
				Config: loadUnitTestTerraform("004_lifecycle_maximal_to_minimal_step_2.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaGroupPolicyBooleanValue.ResourceName+".test_004").Key("id").Exists(),
					check.That(graphBetaGroupPolicyBooleanValue.ResourceName+".test_004").Key("enabled").HasValue("false"),
					check.That(graphBetaGroupPolicyBooleanValue.ResourceName+".test_004").Key("values.#").HasValue("1"),
					check.That(graphBetaGroupPolicyBooleanValue.ResourceName+".test_004").Key("values.0.value").HasValue("false"),
				),
			},
			{
				ResourceName:      graphBetaGroupPolicyBooleanValue.ResourceName + ".test_004",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Test 005: Scenario 5 - Validation errors
func TestGroupPolicyBooleanValueResource_005_ValidationErrors(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, groupPolicyBooleanValueMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer groupPolicyBooleanValueMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("005_validation_missing_policy_name.tf"),
				ExpectError: regexp.MustCompile(`The argument "policy_name" is required`),
			},
			{
				Config:      loadUnitTestTerraform("005_validation_invalid_class_type.tf"),
				ExpectError: regexp.MustCompile(`Attribute class_type value must be one of`),
			},
		},
	})
}

