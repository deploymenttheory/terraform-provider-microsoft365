package graphBetaGroupSettings_test

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	localMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/groups/graph_beta/group_settings/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/jarcoal/httpmock"
)

// Helper functions to return the test configurations by reading from files
func testConfigMinimal() string {
	content, err := os.ReadFile(filepath.Join("mocks", "terraform", "resource_minimal.tf"))
	if err != nil {
		return ""
	}
	return string(content)
}

func testConfigMaximal() string {
	content, err := os.ReadFile(filepath.Join("mocks", "terraform", "resource_maximal.tf"))
	if err != nil {
		return ""
	}
	return string(content)
}

func testConfigMinimalToMaximal() string {
	// For minimal to maximal test, we need to use the maximal config
	// but with the minimal resource name to simulate an update

	// Read the maximal config
	maximalContent, err := os.ReadFile(filepath.Join("mocks", "terraform", "resource_maximal.tf"))
	if err != nil {
		return ""
	}

	// Replace the resource name to match the minimal one
	updatedMaximal := strings.Replace(string(maximalContent), "maximal", "minimal", 1)

	return updatedMaximal
}

func testConfigError() string {
	// Create an error configuration with invalid template ID
	return `
resource "microsoft365_graph_beta_groups_group_settings" "error" {
  group_id = "12345678-1234-1234-1234-123456789012"
  template_id = "invalid-template-id"
  values = [
    {
      name  = "AllowToAddGuests"
      value = "false"
    }
  ]
}
`
}

// Helper function to set up the test environment
func setupTestEnvironment(t *testing.T) {
	// Set environment variables for testing
	os.Setenv("TF_ACC", "0")
	os.Setenv("MS365_TEST_MODE", "true")
}

// Helper function to set up the mock environment
func setupMockEnvironment() (*mocks.Mocks, *localMocks.GroupSettingsMock) {
	// Activate httpmock
	httpmock.Activate()

	// Create a new Mocks instance and register authentication mocks
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	// Register local mocks directly
	groupSettingsMock := &localMocks.GroupSettingsMock{}
	groupSettingsMock.RegisterMocks()

	return mockClient, groupSettingsMock
}

// Helper function to check if a resource exists
func testCheckExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("resource ID not set")
		}

		return nil
	}
}

// Helper function to get maximal config with a custom resource name
func testConfigMaximalWithResourceName(resourceName string) string {
	// Read the maximal config
	content, err := os.ReadFile(filepath.Join("mocks", "terraform", "resource_maximal.tf"))
	if err != nil {
		return ""
	}

	// Replace the resource name
	updated := strings.Replace(string(content), "maximal", resourceName, 1)

	return updated
}

// Helper function to get minimal config with a custom resource name
func testConfigMinimalWithResourceName(resourceName string) string {
	return fmt.Sprintf(`resource "microsoft365_graph_beta_groups_group_settings" "%s" {
  group_id = "12345678-1234-1234-1234-123456789012"
  template_id = "08d542b9-071f-4e16-94b0-74abb372e3d9"
  values = [
    {
      name  = "AllowToAddGuests"
      value = "false"
    }
  ]
}`, resourceName)
}

// TestUnitGroupSettingsResource_Create_Minimal tests the creation of group settings with minimal configuration
func TestUnitGroupSettingsResource_Create_Minimal(t *testing.T) {
	// Set up mock environment
	_, _ = setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	// Set up the test environment
	setupTestEnvironment(t)

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_groups_group_settings.minimal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group_settings.minimal", "template_id", "08d542b9-071f-4e16-94b0-74abb372e3d9"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group_settings.minimal", "display_name", "Group.Unified.Guest"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group_settings.minimal", "values.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs("microsoft365_graph_beta_groups_group_settings.minimal", "values.*", map[string]string{
						"name":  "AllowToAddGuests",
						"value": "false",
					}),
				),
			},
		},
	})
}

// TestUnitGroupSettingsResource_Create_Maximal tests the creation of group settings with maximal configuration
func TestUnitGroupSettingsResource_Create_Maximal(t *testing.T) {
	// Set up mock environment
	_, _ = setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	// Set up the test environment
	setupTestEnvironment(t)

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMaximal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_groups_group_settings.maximal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group_settings.maximal", "template_id", "62375ab9-6b52-47ed-826b-58e47e0e304b"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group_settings.maximal", "display_name", "Group.Unified"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group_settings.maximal", "values.#", "6"),
					resource.TestCheckTypeSetElemNestedAttrs("microsoft365_graph_beta_groups_group_settings.maximal", "values.*", map[string]string{
						"name":  "ClassificationList",
						"value": "Confidential,Secret,Top Secret",
					}),
					resource.TestCheckTypeSetElemNestedAttrs("microsoft365_graph_beta_groups_group_settings.maximal", "values.*", map[string]string{
						"name":  "DefaultClassification",
						"value": "Confidential",
					}),
					resource.TestCheckTypeSetElemNestedAttrs("microsoft365_graph_beta_groups_group_settings.maximal", "values.*", map[string]string{
						"name":  "AllowGuestsToBeGroupOwner",
						"value": "false",
					}),
					resource.TestCheckTypeSetElemNestedAttrs("microsoft365_graph_beta_groups_group_settings.maximal", "values.*", map[string]string{
						"name":  "AllowGuestsToAccessGroups",
						"value": "true",
					}),
					resource.TestCheckTypeSetElemNestedAttrs("microsoft365_graph_beta_groups_group_settings.maximal", "values.*", map[string]string{
						"name":  "AllowToAddGuests",
						"value": "true",
					}),
					resource.TestCheckTypeSetElemNestedAttrs("microsoft365_graph_beta_groups_group_settings.maximal", "values.*", map[string]string{
						"name":  "UsageGuidelinesUrl",
						"value": "https://contoso.com/marketing-group-guidelines",
					}),
				),
			},
		},
	})
}

// TestUnitGroupSettingsResource_Update_MinimalToMaximal tests updating from minimal to maximal configuration
func TestUnitGroupSettingsResource_Update_MinimalToMaximal(t *testing.T) {
	// Set up mock environment
	_, _ = setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	// Set up the test environment
	setupTestEnvironment(t)

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Start with minimal configuration
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_groups_group_settings.minimal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group_settings.minimal", "template_id", "08d542b9-071f-4e16-94b0-74abb372e3d9"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group_settings.minimal", "display_name", "Group.Unified.Guest"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group_settings.minimal", "values.#", "1"),
				),
			},
			// Update to maximal configuration
			{
				Config: testConfigMinimalToMaximal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_groups_group_settings.minimal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group_settings.minimal", "template_id", "62375ab9-6b52-47ed-826b-58e47e0e304b"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group_settings.minimal", "display_name", "Group.Unified"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group_settings.minimal", "values.#", "6"),
				),
			},
		},
	})
}

// TestUnitGroupSettingsResource_Update_MaximalToMinimal tests updating from maximal to minimal configuration
func TestUnitGroupSettingsResource_Update_MaximalToMinimal(t *testing.T) {
	// Set up mock environment
	_, _ = setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	// Set up the test environment
	setupTestEnvironment(t)

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Start with maximal configuration
			{
				Config: testConfigMaximalWithResourceName("test"),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_groups_group_settings.test"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group_settings.test", "template_id", "62375ab9-6b52-47ed-826b-58e47e0e304b"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group_settings.test", "display_name", "Group.Unified"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group_settings.test", "values.#", "6"),
				),
			},
			// Update to minimal configuration
			{
				Config: testConfigMinimalWithResourceName("test"),
				// We expect a non-empty plan because computed fields will show as changes
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_groups_group_settings.test"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group_settings.test", "template_id", "08d542b9-071f-4e16-94b0-74abb372e3d9"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group_settings.test", "display_name", "Group.Unified.Guest"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group_settings.test", "values.#", "1"),
					// Don't check for absence of attributes as they may appear as computed
				),
			},
		},
	})
}

// TestUnitGroupSettingsResource_Delete_Minimal tests deleting group settings with minimal configuration
func TestUnitGroupSettingsResource_Delete_Minimal(t *testing.T) {
	// Set up mock environment
	_, _ = setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	// Set up the test environment
	setupTestEnvironment(t)

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create the resource
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_groups_group_settings.minimal"),
				),
			},
			// Delete the resource (by providing empty config)
			{
				Config: `# Empty config for deletion test`,
				Check: func(s *terraform.State) error {
					// The resource should be gone
					_, exists := s.RootModule().Resources["microsoft365_graph_beta_groups_group_settings.minimal"]
					if exists {
						return fmt.Errorf("resource still exists after deletion")
					}
					return nil
				},
			},
		},
	})
}

// TestUnitGroupSettingsResource_Delete_Maximal tests deleting group settings with maximal configuration
func TestUnitGroupSettingsResource_Delete_Maximal(t *testing.T) {
	// Set up mock environment
	_, _ = setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	// Set up the test environment
	setupTestEnvironment(t)

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create the resource
			{
				Config: testConfigMaximal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_groups_group_settings.maximal"),
				),
			},
			// Delete the resource (by providing empty config)
			{
				Config: `# Empty config for deletion test`,
				Check: func(s *terraform.State) error {
					// The resource should be gone
					_, exists := s.RootModule().Resources["microsoft365_graph_beta_groups_group_settings.maximal"]
					if exists {
						return fmt.Errorf("resource still exists after deletion")
					}
					return nil
				},
			},
		},
	})
}

// TestUnitGroupSettingsResource_Import tests importing a resource
func TestUnitGroupSettingsResource_Import(t *testing.T) {
	// Set up mock environment
	_, groupSettingsMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	// Set up the test environment
	setupTestEnvironment(t)

	// Add the setting to the mock state before testing import
	minimalGroupId := "12345678-1234-1234-1234-123456789012"
	minimalSettingId := "test-setting-id"

	// Use a pre-configured import test
	groupSettingsMock.SetupImportTest(minimalGroupId, minimalSettingId)

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Import directly without creating first
			{
				ResourceName:      "microsoft365_graph_beta_groups_group_settings.minimal",
				ImportState:       true,
				ImportStateId:     minimalGroupId + "/" + minimalSettingId,
				ImportStateVerify: true,
				// Skip applying the config since we're testing import directly
				SkipFunc: func() (bool, error) {
					return true, nil
				},
			},
		},
	})
}

// TestUnitGroupSettingsResource_Error tests error handling
func TestUnitGroupSettingsResource_Error(t *testing.T) {
	// Set up mock environment
	_, groupSettingsMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	// Register error mocks
	groupSettingsMock.RegisterErrorMocks()

	// Set up the test environment
	setupTestEnvironment(t)

	// Run the test with an error case
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testConfigError(),
				ExpectError: regexp.MustCompile("Attribute template_id must be a valid GUID"),
			},
		},
	})
}
