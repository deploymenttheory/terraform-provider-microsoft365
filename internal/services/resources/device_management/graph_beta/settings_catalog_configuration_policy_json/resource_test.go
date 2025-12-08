package graphBetaSettingsCatalogConfigurationPolicyJson_test

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	settingsCatalogJson "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/settings_catalog_configuration_policy_json"
	settingsCatalogJsonMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/settings_catalog_configuration_policy_json/mocks"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	terraformResource "github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

// setupMockEnvironment sets up the mock environment using centralized mocks
func setupMockEnvironment() (*mocks.Mocks, *settingsCatalogJsonMocks.SettingsCatalogConfigurationPolicyJsonMock) {
	// Activate httpmock
	httpmock.Activate()

	// Create a new Mocks instance and register authentication mocks
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	// Register local mocks directly
	settingsCatalogJsonMock := &settingsCatalogJsonMocks.SettingsCatalogConfigurationPolicyJsonMock{}
	settingsCatalogJsonMock.RegisterMocks()

	return mockClient, settingsCatalogJsonMock
}

// setupErrorMockEnvironment sets up the mock environment for error testing
func setupErrorMockEnvironment() (*mocks.Mocks, *settingsCatalogJsonMocks.SettingsCatalogConfigurationPolicyJsonMock) {
	// Activate httpmock
	httpmock.Activate()

	// Create a new Mocks instance and register authentication mocks
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	// Register error mocks
	settingsCatalogJsonMock := &settingsCatalogJsonMocks.SettingsCatalogConfigurationPolicyJsonMock{}
	settingsCatalogJsonMock.RegisterErrorMocks()

	return mockClient, settingsCatalogJsonMock
}

func TestUnitSettingsCatalogConfigurationPolicyJsonResource(t *testing.T) {
	t.Run("resource schema validation", func(t *testing.T) {
		// Test resource schema construction without full provider initialization
		// This avoids the deep recursion issue while still validating schema structure

		startTime := time.Now()

		// Create resource instance
		resourceInstance := settingsCatalogJson.NewSettingsCatalogJsonResource()

		// Create schema request/response
		req := resource.SchemaRequest{}
		resp := &resource.SchemaResponse{}

		// Test that schema construction completes within reasonable time
		resourceInstance.Schema(context.Background(), req, resp)

		elapsed := time.Since(startTime)
		if elapsed > time.Second*30 { // Allow reasonable time but avoid timeout
			t.Fatalf("Schema construction took too long: %v", elapsed)
		}

		// Validate that the schema was constructed successfully
		if resp.Schema.Attributes == nil {
			t.Fatal("Schema attributes should not be nil")
		}

		// Test that main resource attributes exist
		expectedAttrs := []string{"id", "name", "description", "platforms", "settings", "assignments"}
		for _, attr := range expectedAttrs {
			if _, exists := resp.Schema.Attributes[attr]; !exists {
				t.Fatalf("Resource attribute %s should exist", attr)
			}
		}

		// Test that settings attribute is correctly structured as StringAttribute
		settingsAttr, exists := resp.Schema.Attributes["settings"]
		if !exists {
			t.Fatal("settings attribute should exist")
		}

		stringAttr, ok := settingsAttr.(schema.StringAttribute)
		if !ok {
			t.Fatal("settings should be a StringAttribute")
		}

		if !stringAttr.Required {
			t.Fatal("settings attribute should be required")
		}

		t.Logf("Resource schema validation passed in %v - JSON settings properly configured", elapsed)
	})

	t.Run("assignment schema validation", func(t *testing.T) {
		// Test assignment-related attributes in schema

		resourceInstance := settingsCatalogJson.NewSettingsCatalogJsonResource()
		req := resource.SchemaRequest{}
		resp := &resource.SchemaResponse{}

		resourceInstance.Schema(context.Background(), req, resp)

		// Validate assignments attribute exists
		if assignmentAttr, exists := resp.Schema.Attributes["assignments"]; exists {
			// This is a complex nested attribute - just verify it exists and is structured properly
			if assignmentAttr == nil {
				t.Fatal("assignments attribute should not be nil")
			}
		} else {
			t.Fatal("assignments attribute should exist")
		}

		// Validate role_scope_tag_ids attribute
		if roleScopeAttr, exists := resp.Schema.Attributes["role_scope_tag_ids"]; exists {
			if roleScopeAttr == nil {
				t.Fatal("role_scope_tag_ids attribute should not be nil")
			}
		} else {
			t.Fatal("role_scope_tag_ids attribute should exist")
		}

		t.Log("Assignment schema validation passed")
	})

	t.Run("platform and technology validation", func(t *testing.T) {
		// Test platform and technology attributes validation

		resourceInstance := settingsCatalogJson.NewSettingsCatalogJsonResource()
		req := resource.SchemaRequest{}
		resp := &resource.SchemaResponse{}

		resourceInstance.Schema(context.Background(), req, resp)

		// Validate platforms attribute
		if platformAttr, exists := resp.Schema.Attributes["platforms"]; exists {
			stringAttr, ok := platformAttr.(schema.StringAttribute)
			if !ok {
				t.Fatal("platforms should be a StringAttribute")
			}
			if !stringAttr.Optional {
				t.Fatal("platforms should be optional")
			}
			if !stringAttr.Computed {
				t.Fatal("platforms should be computed")
			}
		} else {
			t.Fatal("platforms attribute should exist")
		}

		// Validate technologies attribute
		if techAttr, exists := resp.Schema.Attributes["technologies"]; exists {
			listAttr, ok := techAttr.(schema.ListAttribute)
			if !ok {
				t.Fatal("technologies should be a ListAttribute")
			}
			if !listAttr.Optional {
				t.Fatal("technologies should be optional")
			}
			if !listAttr.Computed {
				t.Fatal("technologies should be computed")
			}
		} else {
			t.Fatal("technologies attribute should exist")
		}

		t.Log("Platform and technology validation passed")
	})
}

// Unit tests for JSON settings validation
func TestUnitConstructSettingsCatalogJsonSettings(t *testing.T) {
	t.Run("resource schema validation", func(t *testing.T) {
		// Test resource schema construction without full provider initialization
		// This avoids the deep recursion issue while still validating schema structure

		startTime := time.Now()

		// Create resource instance
		resourceInstance := settingsCatalogJson.NewSettingsCatalogJsonResource()

		// Create schema request/response
		req := resource.SchemaRequest{}
		resp := &resource.SchemaResponse{}

		// Test that schema construction completes within reasonable time
		resourceInstance.Schema(context.Background(), req, resp)

		elapsed := time.Since(startTime)
		if elapsed > time.Second*30 { // Allow more time but still reasonable
			t.Fatalf("Schema construction took too long: %v", elapsed)
		}

		// Validate that the schema was constructed successfully
		if resp.Schema.Attributes == nil {
			t.Fatal("Schema attributes should not be nil")
		}

		// Test that main resource attributes exist
		expectedAttrs := []string{"id", "name", "description", "platforms", "settings", "assignments"}
		for _, attr := range expectedAttrs {
			if _, exists := resp.Schema.Attributes[attr]; !exists {
				t.Fatalf("Resource attribute %s should exist", attr)
			}
		}

		// Test that settings attribute is correctly structured as JSON string
		settingsAttr, exists := resp.Schema.Attributes["settings"]
		if !exists {
			t.Fatal("settings attribute should exist")
		}

		stringAttr, ok := settingsAttr.(schema.StringAttribute)
		if !ok {
			t.Fatal("settings should be a StringAttribute for JSON input")
		}

		if !stringAttr.Required {
			t.Fatal("settings attribute should be required")
		}

		t.Logf("Resource schema validation passed in %v - JSON settings properly configured", elapsed)
	})

	t.Run("schema performance validation", func(t *testing.T) {
		// Test that multiple schema constructions don't cause performance issues

		startTime := time.Now()

		for i := 0; i < 3; i++ { // Test multiple constructions
			resourceInstance := settingsCatalogJson.NewSettingsCatalogJsonResource()
			req := resource.SchemaRequest{}
			resp := &resource.SchemaResponse{}

			resourceInstance.Schema(context.Background(), req, resp)

			if resp.Schema.Attributes == nil {
				t.Fatalf("Schema attributes should not be nil on iteration %d", i)
			}
		}

		elapsed := time.Since(startTime)
		if elapsed > time.Minute*2 { // Allow reasonable time for multiple constructions
			t.Fatalf("Multiple schema constructions took too long: %v", elapsed)
		}

		t.Logf("Multiple schema constructions completed in %v", elapsed)
	})

	t.Run("basic attribute validation", func(t *testing.T) {
		// Test basic attribute structure without deep recursion

		resourceInstance := settingsCatalogJson.NewSettingsCatalogJsonResource()
		req := resource.SchemaRequest{}
		resp := &resource.SchemaResponse{}

		resourceInstance.Schema(context.Background(), req, resp)

		// Test basic required attributes
		if idAttr, exists := resp.Schema.Attributes["id"]; exists {
			stringAttr, ok := idAttr.(schema.StringAttribute)
			if !ok {
				t.Fatal("id should be a StringAttribute")
			}
			if !stringAttr.Computed {
				t.Fatal("id should be computed")
			}
		} else {
			t.Fatal("id attribute should exist")
		}

		if nameAttr, exists := resp.Schema.Attributes["name"]; exists {
			stringAttr, ok := nameAttr.(schema.StringAttribute)
			if !ok {
				t.Fatal("name should be a StringAttribute")
			}
			if !stringAttr.Required {
				t.Fatal("name should be required")
			}
		} else {
			t.Fatal("name attribute should exist")
		}

		t.Log("Basic attribute validation passed")
	})
}

// TestSettingsCatalogConfigurationPolicyJsonResource_ErrorHandling tests error scenarios
func TestSettingsCatalogConfigurationPolicyJsonResource_ErrorHandling(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, settingsCatalogJsonMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer settingsCatalogJsonMock.CleanupMockState()

	terraformResource.UnitTest(t, terraformResource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []terraformResource.TestStep{
			// Test invalid configuration - missing required name field
			{
				Config: `
resource "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy_json" "test" {
  platforms = "macOS"
  technologies = ["mdm"]
  settings = jsonencode({
    "settings": []
  })
}
`,
				ExpectError: regexp.MustCompile(`Missing required argument|name`),
			},
			// Test invalid platforms value
			{
				Config: `
resource "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy_json" "test" {
  name = "Test Policy"
  platforms = "invalid_platform"
  technologies = ["mdm"]
  settings = jsonencode({
    "settings": []
  })
}
`,
				ExpectError: regexp.MustCompile(`Attribute platforms value must be one of|invalid_platform`),
			},
			// Test invalid technologies value
			{
				Config: `
resource "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy_json" "test" {
  name = "Test Policy"
  platforms = "macOS"
  technologies = ["invalid_technology"]
  settings = jsonencode({
    "settings": []
  })
}
`,
				ExpectError: regexp.MustCompile(`invalid value for technologies|invalid_technology`),
			},
			// Test invalid JSON settings
			{
				Config: `
resource "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy_json" "test" {
  name = "Test Policy"
  platforms = "macOS"
  technologies = ["mdm"]
  settings = "invalid json string"
}
`,
				ExpectError: regexp.MustCompile(`Invalid JSON|invalid character`),
			},
		},
	})
}

// TestSettingsCatalogConfigurationPolicyJsonResource_JSONValidation tests JSON-specific scenarios
func TestSettingsCatalogConfigurationPolicyJsonResource_JSONValidation(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, settingsCatalogJsonMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer settingsCatalogJsonMock.CleanupMockState()

	terraformResource.UnitTest(t, terraformResource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []terraformResource.TestStep{
			// Test valid JSON structure - plan validation only
			{
				Config: `
resource "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy_json" "test" {
  name = "Test JSON Policy"
  platforms = "macOS"
  technologies = ["mdm"]
  settings = jsonencode({
    "settings": [
      {
        "id": "0",
        "settingInstance": {
          "@odata.type": "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance",
          "settingDefinitionId": "test.setting",
          "simpleSettingValue": {
            "@odata.type": "#microsoft.graph.deviceManagementConfigurationStringSettingValue",
            "value": "test_value"
          }
        }
      }
    ]
  })
}
`,
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			// Test secret setting with proper valueState
			{
				Config: `
resource "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy_json" "test" {
  name = "Test Secret JSON Policy"
  platforms = "macOS"
  technologies = ["mdm"]
  settings = jsonencode({
    "settings": [
      {
        "id": "0",
        "settingInstance": {
          "@odata.type": "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance",
          "settingDefinitionId": "test.secret.setting",
          "simpleSettingValue": {
            "@odata.type": "#microsoft.graph.deviceManagementConfigurationSecretSettingValue",
            "value": "secret_value",
            "valueState": "notEncrypted"
          }
        }
      }
    ]
  })
}
`,
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

// TestSettingsCatalogConfigurationPolicyJsonResource_Schema validates the resource schema
func TestSettingsCatalogConfigurationPolicyJsonResource_Schema(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, settingsCatalogJsonMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer settingsCatalogJsonMock.CleanupMockState()

	terraformResource.UnitTest(t, terraformResource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []terraformResource.TestStep{
			// Test JSON Simple String Setting Schema - Plan Only
			{
				Config:             testUnitSettingsCatalogJsonSimpleString(),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			// Test JSON Simple Secret Setting Schema - Plan Only
			{
				Config:             testUnitSettingsCatalogJsonSimpleSecret(),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			// Test JSON Choice Setting Schema - Plan Only
			{
				Config:             testUnitSettingsCatalogJsonChoice(),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			// Test JSON Simple Collection Setting Schema - Plan Only
			{
				Config:             testUnitSettingsCatalogJsonSimpleCollection(),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			// Test JSON Choice Collection Setting Schema - Plan Only
			{
				Config:             testUnitSettingsCatalogJsonChoiceCollection(),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			// Test JSON Group Collection Setting Schema - Plan Only
			{
				Config:             testUnitSettingsCatalogJsonGroupCollection(),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			// Test JSON Complex Group Collection Setting Schema - Plan Only
			{
				Config:             testUnitSettingsCatalogJsonComplexGroupCollection(),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

// Test configuration functions for different setting types (JSON variants)
func testUnitSettingsCatalogJsonSimpleString() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_simple_string.tf")
	if err != nil {
		panic("failed to load simple string config: " + err.Error())
	}
	return unitTestConfig
}

func testUnitSettingsCatalogJsonSimpleSecret() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_simple_secret.tf")
	if err != nil {
		panic("failed to load simple secret config: " + err.Error())
	}
	return unitTestConfig
}

func testUnitSettingsCatalogJsonChoice() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_choice.tf")
	if err != nil {
		panic("failed to load choice config: " + err.Error())
	}
	return unitTestConfig
}

func testUnitSettingsCatalogJsonSimpleCollection() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_simple_collection.tf")
	if err != nil {
		panic("failed to load simple collection config: " + err.Error())
	}
	return unitTestConfig
}

func testUnitSettingsCatalogJsonChoiceCollection() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_choice_collection.tf")
	if err != nil {
		panic("failed to load choice collection config: " + err.Error())
	}
	return unitTestConfig
}

func testUnitSettingsCatalogJsonGroupCollection() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_group_collection.tf")
	if err != nil {
		panic("failed to load group collection config: " + err.Error())
	}
	return unitTestConfig
}

func testUnitSettingsCatalogJsonComplexGroupCollection() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_complex_group_collection.tf")
	if err != nil {
		panic("failed to load complex group collection config: " + err.Error())
	}
	return unitTestConfig
}
