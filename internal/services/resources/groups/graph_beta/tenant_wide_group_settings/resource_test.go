package graphBetaTenantWideGroupSettings_test

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	localMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/groups/graph_beta/tenant_wide_group_settings/mocks"
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
	// Replace "test" with "minimal" in the resource name
	return strings.Replace(string(content), `"test"`, `"minimal"`, 1)
}

func testConfigMaximal() string {
	content, err := os.ReadFile(filepath.Join("mocks", "terraform", "resource_maximal.tf"))
	if err != nil {
		return ""
	}
	// Replace "test" with "maximal" in the resource name
	return strings.Replace(string(content), `"test"`, `"maximal"`, 1)
}

func testConfigMinimalToMaximal() string {
	// Read the maximal config
	maximalContent, err := os.ReadFile(filepath.Join("mocks", "terraform", "resource_maximal.tf"))
	if err != nil {
		return ""
	}

	// Replace "test" with "minimal" in the resource name
	updatedMaximal := strings.Replace(string(maximalContent), `"test"`, `"minimal"`, 1)

	return updatedMaximal
}

func testConfigError() string {
	// Read the error config
	content, err := os.ReadFile(filepath.Join("mocks", "terraform", "resource_error.tf"))
	if err != nil {
		return ""
	}
	return string(content)
}

// Helper function to set up the test environment
func setupTestEnvironment(t *testing.T) {
	// Set environment variables for testing
	os.Setenv("TF_ACC", "0")
	os.Setenv("MS365_TEST_MODE", "true")
}

// Helper function to set up the mock environment
func setupMockEnvironment() {
	// Activate httpmock
	httpmock.Activate()

	// Register local mocks directly
	localMocks.SetupTenantWideGroupSettingsMocks()
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

// TestUnitTenantWideGroupSettingsResource_Create_Minimal tests the creation of tenant-wide group settings with minimal configuration
func TestUnitTenantWideGroupSettingsResource_Create_Minimal(t *testing.T) {
	// Set up mock environment
	setupMockEnvironment()
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
					testCheckExists("microsoft365_graph_beta_groups_tenant_wide_group_settings.minimal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_tenant_wide_group_settings.minimal", "template_id", "62375ab9-6b52-47ed-826b-58e47e0e304b"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_tenant_wide_group_settings.minimal", "display_name", "Group.Unified"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_tenant_wide_group_settings.minimal", "values.#", "3"),
					resource.TestCheckTypeSetElemNestedAttrs("microsoft365_graph_beta_groups_tenant_wide_group_settings.minimal", "values.*", map[string]string{
						"name":  "EnableGroupCreation",
						"value": "true",
					}),
					resource.TestCheckTypeSetElemNestedAttrs("microsoft365_graph_beta_groups_tenant_wide_group_settings.minimal", "values.*", map[string]string{
						"name":  "AllowGuestsToAccessGroups",
						"value": "true",
					}),
					resource.TestCheckTypeSetElemNestedAttrs("microsoft365_graph_beta_groups_tenant_wide_group_settings.minimal", "values.*", map[string]string{
						"name":  "AllowToAddGuests",
						"value": "true",
					}),
				),
			},
		},
	})
}

// TestUnitTenantWideGroupSettingsResource_Create_Maximal tests the creation of tenant-wide group settings with maximal configuration
func TestUnitTenantWideGroupSettingsResource_Create_Maximal(t *testing.T) {
	// Set up mock environment
	setupMockEnvironment()
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
					testCheckExists("microsoft365_graph_beta_groups_tenant_wide_group_settings.maximal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_tenant_wide_group_settings.maximal", "template_id", "62375ab9-6b52-47ed-826b-58e47e0e304b"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_tenant_wide_group_settings.maximal", "display_name", "Group.Unified"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_tenant_wide_group_settings.maximal", "values.#", "15"),
					// Check a few key values
					resource.TestCheckTypeSetElemNestedAttrs("microsoft365_graph_beta_groups_tenant_wide_group_settings.maximal", "values.*", map[string]string{
						"name":  "EnableGroupCreation",
						"value": "false",
					}),
					resource.TestCheckTypeSetElemNestedAttrs("microsoft365_graph_beta_groups_tenant_wide_group_settings.maximal", "values.*", map[string]string{
						"name":  "GroupCreationAllowedGroupId",
						"value": "12345678-1234-1234-1234-123456789012",
					}),
				),
			},
		},
	})
}

// TestUnitTenantWideGroupSettingsResource_Update_MinimalToMaximal tests updating from minimal to maximal configuration
func TestUnitTenantWideGroupSettingsResource_Update_MinimalToMaximal(t *testing.T) {
	// Set up mock environment
	setupMockEnvironment()
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
					testCheckExists("microsoft365_graph_beta_groups_tenant_wide_group_settings.minimal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_tenant_wide_group_settings.minimal", "values.#", "3"),
				),
			},
			// Update to maximal configuration
			{
				Config: testConfigMinimalToMaximal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_groups_tenant_wide_group_settings.minimal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_tenant_wide_group_settings.minimal", "values.#", "15"),
				),
			},
		},
	})
}

// TestUnitTenantWideGroupSettingsResource_Delete_Minimal tests deleting tenant-wide group settings with minimal configuration
func TestUnitTenantWideGroupSettingsResource_Delete_Minimal(t *testing.T) {
	// Set up mock environment
	setupMockEnvironment()
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
					testCheckExists("microsoft365_graph_beta_groups_tenant_wide_group_settings.minimal"),
				),
			},
			// Delete the resource (by providing empty config)
			{
				Config: `# Empty config for deletion test`,
				Check: func(s *terraform.State) error {
					// The resource should be gone
					_, exists := s.RootModule().Resources["microsoft365_graph_beta_groups_tenant_wide_group_settings.minimal"]
					if exists {
						return fmt.Errorf("resource still exists after deletion")
					}
					return nil
				},
			},
		},
	})
}

// TestUnitTenantWideGroupSettingsResource_Import tests importing a resource
func TestUnitTenantWideGroupSettingsResource_Import(t *testing.T) {
	// Set up mock environment
	setupMockEnvironment()
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
					testCheckExists("microsoft365_graph_beta_groups_tenant_wide_group_settings.minimal"),
				),
			},
			{
				ResourceName:      "microsoft365_graph_beta_groups_tenant_wide_group_settings.minimal",
				ImportState:       true,
				ImportStateId:     "test-tenant-setting-id",
				ImportStateVerify: true,
			},
		},
	})
}

// TestUnitTenantWideGroupSettingsResource_Error tests error handling
func TestUnitTenantWideGroupSettingsResource_Error(t *testing.T) {
	// Set up mock environment
	setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	// Register error mocks
	localMocks.RegisterErrorMocks()

	// Set up the test environment
	setupTestEnvironment(t)

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testConfigError(),
				ExpectError: regexp.MustCompile("Internal server error occurred"),
			},
		},
	})
}
