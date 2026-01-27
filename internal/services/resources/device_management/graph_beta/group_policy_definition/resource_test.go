package graphBetaGroupPolicyDefinition_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaGroupPolicyDefinition "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/group_policy_definition"
	groupPolicyDefinitionMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/group_policy_definition/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

// setupMockEnvironment sets up the mock environment using centralized mocks
func setupMockEnvironment() (*mocks.Mocks, *groupPolicyDefinitionMocks.GroupPolicyDefinitionMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	groupPolicyDefinitionMock := &groupPolicyDefinitionMocks.GroupPolicyDefinitionMock{}
	groupPolicyDefinitionMock.RegisterMocks()
	return mockClient, groupPolicyDefinitionMock
}

// setupErrorMockEnvironment sets up the mock environment for error testing
func setupErrorMockEnvironment() (*mocks.Mocks, *groupPolicyDefinitionMocks.GroupPolicyDefinitionMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	groupPolicyDefinitionMock := &groupPolicyDefinitionMocks.GroupPolicyDefinitionMock{}
	groupPolicyDefinitionMock.RegisterErrorMocks()
	return mockClient, groupPolicyDefinitionMock
}

// Helper function to load test configs from unit directory
func loadUnitTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/unit/" + filename)
	if err != nil {
		panic("failed to load unit test config " + filename + ": " + err.Error())
	}
	return config
}

// Test 001: Boolean - Minimal configuration (2 checkboxes)
func TestUnitResourceGroupPolicyDefinition_01_Boolean_Minimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, groupPolicyDefinitionMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer groupPolicyDefinitionMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("001_boolean_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaGroupPolicyDefinition.ResourceName+".test_001").Key("id").Exists(),
					check.That(graphBetaGroupPolicyDefinition.ResourceName+".test_001").Key("group_policy_configuration_id").HasValue("config-001"),
					check.That(graphBetaGroupPolicyDefinition.ResourceName+".test_001").Key("policy_name").HasValue("Test Policy Boolean Minimal"),
					check.That(graphBetaGroupPolicyDefinition.ResourceName+".test_001").Key("class_type").HasValue("machine"),
					check.That(graphBetaGroupPolicyDefinition.ResourceName+".test_001").Key("category_path").HasValue("\\Test\\Boolean\\Minimal"),
					check.That(graphBetaGroupPolicyDefinition.ResourceName+".test_001").Key("enabled").HasValue("true"),
					check.That(graphBetaGroupPolicyDefinition.ResourceName+".test_001").Key("values.#").HasValue("2"),
					// Verify IDs are computed
					check.That(graphBetaGroupPolicyDefinition.ResourceName+".test_001").Key("values.0.id").Exists(),
					check.That(graphBetaGroupPolicyDefinition.ResourceName+".test_001").Key("values.1.id").Exists(),
					check.That(graphBetaGroupPolicyDefinition.ResourceName+".test_001").Key("created_date_time").Exists(),
					check.That(graphBetaGroupPolicyDefinition.ResourceName+".test_001").Key("last_modified_date_time").Exists(),
				),
			},
			{
				ResourceName:      graphBetaGroupPolicyDefinition.ResourceName + ".test_001",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Test 002: Boolean - Maximal configuration (25 checkboxes)
func TestUnitResourceGroupPolicyDefinition_02_Boolean_Maximal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, groupPolicyDefinitionMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer groupPolicyDefinitionMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("002_boolean_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaGroupPolicyDefinition.ResourceName+".test_002").Key("id").Exists(),
					check.That(graphBetaGroupPolicyDefinition.ResourceName+".test_002").Key("policy_name").HasValue("Test Policy Boolean Maximal"),
					check.That(graphBetaGroupPolicyDefinition.ResourceName+".test_002").Key("enabled").HasValue("true"),
					check.That(graphBetaGroupPolicyDefinition.ResourceName+".test_002").Key("values.#").HasValue("25"),
					// Verify all IDs are computed
					check.That(graphBetaGroupPolicyDefinition.ResourceName+".test_002").Key("values.0.id").Exists(),
					check.That(graphBetaGroupPolicyDefinition.ResourceName+".test_002").Key("values.24.id").Exists(),
				),
			},
			{
				ResourceName:      graphBetaGroupPolicyDefinition.ResourceName + ".test_002",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Test 003: TextBox - Single string value
func TestUnitResourceGroupPolicyDefinition_03_TextBox(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, groupPolicyDefinitionMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer groupPolicyDefinitionMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("003_textbox.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaGroupPolicyDefinition.ResourceName+".test_003").Key("id").Exists(),
					check.That(graphBetaGroupPolicyDefinition.ResourceName+".test_003").Key("policy_name").HasValue("Test Policy TextBox"),
					check.That(graphBetaGroupPolicyDefinition.ResourceName+".test_003").Key("enabled").HasValue("true"),
					check.That(graphBetaGroupPolicyDefinition.ResourceName+".test_003").Key("values.#").HasValue("1"),
					check.That(graphBetaGroupPolicyDefinition.ResourceName+".test_003").Key("values.0.id").Exists(),
					check.That(graphBetaGroupPolicyDefinition.ResourceName+".test_003").Key("values.0.label").HasValue("Text Setting"),
					check.That(graphBetaGroupPolicyDefinition.ResourceName+".test_003").Key("values.0.value").Exists(),
				),
			},
			{
				ResourceName:      graphBetaGroupPolicyDefinition.ResourceName + ".test_003",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Test 004: Decimal - Numeric value
func TestUnitResourceGroupPolicyDefinition_04_Decimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, groupPolicyDefinitionMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer groupPolicyDefinitionMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("004_decimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaGroupPolicyDefinition.ResourceName+".test_004").Key("id").Exists(),
					check.That(graphBetaGroupPolicyDefinition.ResourceName+".test_004").Key("policy_name").HasValue("Test Policy Decimal"),
					check.That(graphBetaGroupPolicyDefinition.ResourceName+".test_004").Key("enabled").HasValue("true"),
					check.That(graphBetaGroupPolicyDefinition.ResourceName+".test_004").Key("values.#").HasValue("1"),
					check.That(graphBetaGroupPolicyDefinition.ResourceName+".test_004").Key("values.0.id").Exists(),
					check.That(graphBetaGroupPolicyDefinition.ResourceName+".test_004").Key("values.0.value").HasValue("7200"),
				),
			},
			{
				ResourceName:      graphBetaGroupPolicyDefinition.ResourceName + ".test_004",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Test 005: MultiText - Multi-line value
func TestUnitResourceGroupPolicyDefinition_05_MultiText(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, groupPolicyDefinitionMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer groupPolicyDefinitionMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("005_multitext.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaGroupPolicyDefinition.ResourceName+".test_005").Key("id").Exists(),
					check.That(graphBetaGroupPolicyDefinition.ResourceName+".test_005").Key("policy_name").HasValue("Test Policy MultiText"),
					check.That(graphBetaGroupPolicyDefinition.ResourceName+".test_005").Key("enabled").HasValue("true"),
					check.That(graphBetaGroupPolicyDefinition.ResourceName+".test_005").Key("values.#").HasValue("1"),
					check.That(graphBetaGroupPolicyDefinition.ResourceName+".test_005").Key("values.0.id").Exists(),
				),
			},
			{
				ResourceName:      graphBetaGroupPolicyDefinition.ResourceName + ".test_005",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Test 006: Dropdown - Select option value
func TestUnitResourceGroupPolicyDefinition_06_Dropdown(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, groupPolicyDefinitionMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer groupPolicyDefinitionMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("006_dropdown.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaGroupPolicyDefinition.ResourceName+".test_006").Key("id").Exists(),
					check.That(graphBetaGroupPolicyDefinition.ResourceName+".test_006").Key("policy_name").HasValue("Test Policy Dropdown"),
					check.That(graphBetaGroupPolicyDefinition.ResourceName+".test_006").Key("enabled").HasValue("true"),
					check.That(graphBetaGroupPolicyDefinition.ResourceName+".test_006").Key("values.#").HasValue("1"),
					check.That(graphBetaGroupPolicyDefinition.ResourceName+".test_006").Key("values.0.id").Exists(),
					check.That(graphBetaGroupPolicyDefinition.ResourceName+".test_006").Key("values.0.value").HasValue("1"),
				),
			},
			{
				ResourceName:      graphBetaGroupPolicyDefinition.ResourceName + ".test_006",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Test 007: Lifecycle - Type transitions through all presentation types
func TestUnitResourceGroupPolicyDefinition_07_Lifecycle_TypeTransitions(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, groupPolicyDefinitionMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer groupPolicyDefinitionMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Step 1: Boolean
			{
				Config: loadUnitTestTerraform("007_lifecycle_step_1_boolean.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaGroupPolicyDefinition.ResourceName+".test_007").Key("id").Exists(),
					check.That(graphBetaGroupPolicyDefinition.ResourceName+".test_007").Key("policy_name").HasValue("Test Policy Boolean Minimal"),
					check.That(graphBetaGroupPolicyDefinition.ResourceName+".test_007").Key("values.#").HasValue("2"),
				),
			},
			// Step 2: TextBox
			{
				Config: loadUnitTestTerraform("007_lifecycle_step_2_textbox.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaGroupPolicyDefinition.ResourceName+".test_007").Key("id").Exists(),
					check.That(graphBetaGroupPolicyDefinition.ResourceName+".test_007").Key("policy_name").HasValue("Test Policy TextBox"),
					check.That(graphBetaGroupPolicyDefinition.ResourceName+".test_007").Key("values.#").HasValue("1"),
				),
			},
			// Step 3: Decimal
			{
				Config: loadUnitTestTerraform("007_lifecycle_step_3_decimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaGroupPolicyDefinition.ResourceName+".test_007").Key("id").Exists(),
					check.That(graphBetaGroupPolicyDefinition.ResourceName+".test_007").Key("policy_name").HasValue("Test Policy Decimal"),
					check.That(graphBetaGroupPolicyDefinition.ResourceName+".test_007").Key("values.0.value").HasValue("7200"),
				),
			},
			// Step 4: MultiText
			{
				Config: loadUnitTestTerraform("007_lifecycle_step_4_multitext.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaGroupPolicyDefinition.ResourceName+".test_007").Key("id").Exists(),
					check.That(graphBetaGroupPolicyDefinition.ResourceName+".test_007").Key("policy_name").HasValue("Test Policy MultiText"),
				),
			},
			// Step 5: Dropdown
			{
				Config: loadUnitTestTerraform("007_lifecycle_step_5_dropdown.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaGroupPolicyDefinition.ResourceName+".test_007").Key("id").Exists(),
					check.That(graphBetaGroupPolicyDefinition.ResourceName+".test_007").Key("policy_name").HasValue("Test Policy Dropdown"),
					check.That(graphBetaGroupPolicyDefinition.ResourceName+".test_007").Key("values.0.value").HasValue("1"),
				),
			},
			// Note: Import verification removed as lifecycle test creates/destroys multiple policies,
			// making import verification complex in the mock environment
		},
	})
}

// Test 008: Validation - Invalid label
func TestUnitResourceGroupPolicyDefinition_08_Validation_InvalidLabel(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, groupPolicyDefinitionMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer groupPolicyDefinitionMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("008_validation_invalid_label.tf"),
				ExpectError: regexp.MustCompile(`(?s)label.*not found in policy`),
			},
		},
	})
}

// Test 009: Validation - Invalid boolean value
func TestUnitResourceGroupPolicyDefinition_09_Validation_InvalidBooleanValue(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, groupPolicyDefinitionMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer groupPolicyDefinitionMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("009_validation_invalid_boolean.tf"),
				ExpectError: regexp.MustCompile(`(?s)requires a boolean value`),
			},
		},
	})
}

// Test 010: Validation - Invalid decimal value
func TestUnitResourceGroupPolicyDefinition_10_Validation_InvalidDecimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, groupPolicyDefinitionMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer groupPolicyDefinitionMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("010_validation_invalid_decimal.tf"),
				ExpectError: regexp.MustCompile(`requires a numeric value`),
			},
		},
	})
}

// Test 011: Validation - Read-only presentation (Text/Label)
func TestUnitResourceGroupPolicyDefinition_11_Validation_ReadOnlyPresentation(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, groupPolicyDefinitionMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer groupPolicyDefinitionMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("011_validation_readonly_presentation.tf"),
				ExpectError: regexp.MustCompile(`(?s)read-only.*cannot have a value set`),
			},
		},
	})
}

// Test 012: Validation - Missing required fields
func TestUnitResourceGroupPolicyDefinition_12_Validation_MissingFields(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, groupPolicyDefinitionMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer groupPolicyDefinitionMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("012_validation_missing_policy_name.tf"),
				ExpectError: regexp.MustCompile(`The argument "policy_name" is required`),
			},
			{
				Config:      loadUnitTestTerraform("012_validation_missing_class_type.tf"),
				ExpectError: regexp.MustCompile(`The argument "class_type" is required`),
			},
			{
				Config:      loadUnitTestTerraform("012_validation_missing_category_path.tf"),
				ExpectError: regexp.MustCompile(`The argument "category_path" is required`),
			},
		},
	})
}

// Test 013: Validation - Invalid class_type
func TestUnitResourceGroupPolicyDefinition_13_Validation_InvalidClassType(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, groupPolicyDefinitionMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer groupPolicyDefinitionMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("013_validation_invalid_class_type.tf"),
				ExpectError: regexp.MustCompile(`Attribute class_type value must be one of`),
			},
		},
	})
}
