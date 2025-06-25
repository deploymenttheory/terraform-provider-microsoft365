package graphBetaGroup_test

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	localMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/groups/graph_beta/group/mocks"
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
	// Read the minimal config and modify for error scenario
	content, err := os.ReadFile(filepath.Join("mocks", "terraform", "resource_minimal.tf"))
	if err != nil {
		return ""
	}

	// Replace resource name and display name to create an error scenario
	updated := strings.Replace(string(content), "minimal", "error", 1)
	updated = strings.Replace(updated, "Minimal Group", "Error Group", 1)

	return updated
}

// Helper function to set up the test environment
func setupTestEnvironment(t *testing.T) {
	// Set environment variables for testing
	os.Setenv("TF_ACC", "0")
	os.Setenv("MS365_TEST_MODE", "true")
}

// Helper function to set up the mock environment
func setupMockEnvironment() (*mocks.Mocks, *localMocks.GroupMock) {
	// Activate httpmock
	httpmock.Activate()

	// Create a new Mocks instance and register authentication mocks
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	// Register local mocks directly
	groupMock := &localMocks.GroupMock{}
	groupMock.RegisterMocks()

	return mockClient, groupMock
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

// TestUnitGroupResource_Create_Minimal tests the creation of a group with minimal configuration
func TestUnitGroupResource_Create_Minimal(t *testing.T) {
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
					testCheckExists("microsoft365_graph_beta_groups_group.minimal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group.minimal", "display_name", "Minimal Group"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group.minimal", "mail_nickname", "minimal.group"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group.minimal", "mail_enabled", "false"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group.minimal", "security_enabled", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group.minimal", "visibility", "Private"),
				),
			},
		},
	})
}

// TestUnitGroupResource_Create_Maximal tests the creation of a group with maximal configuration
func TestUnitGroupResource_Create_Maximal(t *testing.T) {
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
					testCheckExists("microsoft365_graph_beta_groups_group.maximal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group.maximal", "display_name", "Maximal Group"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group.maximal", "description", "This is a maximal group configuration for testing"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group.maximal", "mail_nickname", "maximal.group"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group.maximal", "mail_enabled", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group.maximal", "security_enabled", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group.maximal", "group_types.#", "2"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group.maximal", "visibility", "Private"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group.maximal", "membership_rule", "user.department -eq \"Engineering\""),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group.maximal", "membership_rule_processing_state", "On"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group.maximal", "preferred_data_location", "NAM"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group.maximal", "preferred_language", "en-US"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group.maximal", "theme", "Blue"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group.maximal", "classification", "High"),
				),
			},
		},
	})
}

// TestUnitGroupResource_Update_MinimalToMaximal tests updating from minimal to maximal configuration
func TestUnitGroupResource_Update_MinimalToMaximal(t *testing.T) {
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
					testCheckExists("microsoft365_graph_beta_groups_group.minimal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group.minimal", "display_name", "Minimal Group"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group.minimal", "mail_nickname", "minimal.group"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group.minimal", "mail_enabled", "false"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group.minimal", "security_enabled", "true"),
					// Verify minimal config doesn't have these attributes
					resource.TestCheckNoResourceAttr("microsoft365_graph_beta_groups_group.minimal", "description"),
					resource.TestCheckNoResourceAttr("microsoft365_graph_beta_groups_group.minimal", "membership_rule"),
				),
			},
			// Update to maximal configuration (with the same resource name)
			{
				Config: testConfigMinimalToMaximal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_groups_group.minimal"),
					// Now check that it has maximal attributes
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group.minimal", "display_name", "Maximal Group"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group.minimal", "description", "This is a maximal group configuration for testing"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group.minimal", "mail_nickname", "maximal.group"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group.minimal", "mail_enabled", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group.minimal", "security_enabled", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group.minimal", "group_types.#", "2"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group.minimal", "membership_rule", "user.department -eq \"Engineering\""),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group.minimal", "membership_rule_processing_state", "On"),
				),
			},
		},
	})
}

// TestUnitGroupResource_Update_MaximalToMinimal tests updating from maximal to minimal configuration
func TestUnitGroupResource_Update_MaximalToMinimal(t *testing.T) {
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
					testCheckExists("microsoft365_graph_beta_groups_group.test"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group.test", "display_name", "Maximal Group"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group.test", "description", "This is a maximal group configuration for testing"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group.test", "group_types.#", "2"),
				),
			},
			// Update to minimal configuration (with the same resource name)
			{
				Config: testConfigMinimalWithResourceName("test"),
				// We expect a non-empty plan because computed fields will show as changes
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_groups_group.test"),
					// Verify it now has only minimal attributes
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group.test", "display_name", "Minimal Group"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group.test", "mail_enabled", "false"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group.test", "security_enabled", "true"),
					// Don't check for absence of attributes as they may appear as computed
				),
			},
		},
	})
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
	return fmt.Sprintf(`resource "microsoft365_graph_beta_groups_group" "%s" {
  display_name     = "Minimal Group"
  mail_nickname    = "test.group"
  mail_enabled     = false
  security_enabled = true
}`, resourceName)
}

// TestUnitGroupResource_Delete_Minimal tests deleting a group with minimal configuration
func TestUnitGroupResource_Delete_Minimal(t *testing.T) {
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
					testCheckExists("microsoft365_graph_beta_groups_group.minimal"),
				),
			},
			// Delete the resource (by providing empty config)
			{
				Config: `# Empty config for deletion test`,
				Check: func(s *terraform.State) error {
					// The resource should be gone
					_, exists := s.RootModule().Resources["microsoft365_graph_beta_groups_group.minimal"]
					if exists {
						return fmt.Errorf("resource still exists after deletion")
					}
					return nil
				},
			},
		},
	})
}

// TestUnitGroupResource_Delete_Maximal tests deleting a group with maximal configuration
func TestUnitGroupResource_Delete_Maximal(t *testing.T) {
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
					testCheckExists("microsoft365_graph_beta_groups_group.maximal"),
				),
			},
			// Delete the resource (by providing empty config)
			{
				Config: `# Empty config for deletion test`,
				Check: func(s *terraform.State) error {
					// The resource should be gone
					_, exists := s.RootModule().Resources["microsoft365_graph_beta_groups_group.maximal"]
					if exists {
						return fmt.Errorf("resource still exists after deletion")
					}
					return nil
				},
			},
		},
	})
}

// TestUnitGroupResource_Import tests importing a resource
func TestUnitGroupResource_Import(t *testing.T) {
	// Set up mock environment
	_, _ = setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	// Set up the test environment
	setupTestEnvironment(t)

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_groups_group.minimal"),
				),
			},
			// Import
			{
				ResourceName:      "microsoft365_graph_beta_groups_group.minimal",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestUnitGroupResource_Error(t *testing.T) {
	// Set up mock environment
	_, groupMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	// Register error mocks
	groupMock.RegisterErrorMocks()

	// Set up the test environment
	setupTestEnvironment(t)

	// Run the test with an error case
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testConfigError(),
				ExpectError: regexp.MustCompile("Group with this displayName already exists"),
			},
		},
	})
}
