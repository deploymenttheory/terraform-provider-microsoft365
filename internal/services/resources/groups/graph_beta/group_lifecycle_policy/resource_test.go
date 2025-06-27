package graphBetaGroupLifecyclePolicy_test

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	localMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/groups/graph_beta/group_lifecycle_policy/mocks"
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
	return `resource "microsoft365_graph_beta_groups_group_lifecycle_policy" "minimal" {
  group_lifetime_in_days         = 365
  managed_group_types            = "Selected"
  alternate_notification_emails  = "admin@example.com;notifications@example.com"
}`
}

func testConfigMaximalToMinimal() string {
	return `resource "microsoft365_graph_beta_groups_group_lifecycle_policy" "maximal" {
  group_lifetime_in_days = 180
  managed_group_types    = "All"
}`
}

func testConfigError() string {
	// Read the minimal config and modify for error scenario
	content, err := os.ReadFile(filepath.Join("mocks", "terraform", "resource_minimal.tf"))
	if err != nil {
		return ""
	}

	// Replace resource name to create an error scenario
	updated := strings.Replace(string(content), "minimal", "error", 1)

	return updated
}

func testConfigImport() string {
	return `resource "microsoft365_graph_beta_groups_group_lifecycle_policy" "imported" {
  group_lifetime_in_days = 180
  managed_group_types    = "All"
}`
}

// Helper function to set up the test environment
func setupTestEnvironment(t *testing.T) {
	// Set environment variables for testing
	os.Setenv("TF_ACC", "0")
	os.Setenv("MS365_TEST_MODE", "true")
}

// Helper function to set up the mock environment
func setupMockEnvironment() (*mocks.Mocks, *localMocks.GroupLifecyclePolicyMock) {
	// Activate httpmock
	httpmock.Activate()

	// Create a new Mocks instance and register authentication mocks
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	// Register local mocks directly
	policyMock := &localMocks.GroupLifecyclePolicyMock{}
	policyMock.RegisterMocks()

	return mockClient, policyMock
}

// Helper function to check if a resource exists
func testCheckExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		fmt.Printf("[DEBUG] testCheckExists: resourceName=%s, ID=%s\n", resourceName, rs.Primary.ID)

		if rs.Primary.ID == "" {
			return fmt.Errorf("resource ID not set")
		}

		return nil
	}
}

// TestUnitGroupLifecyclePolicyResource_Create_Minimal tests the creation of a group lifecycle policy with minimal configuration
func TestUnitGroupLifecyclePolicyResource_Create_Minimal(t *testing.T) {
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
					testCheckExists("microsoft365_graph_beta_groups_group_lifecycle_policy.minimal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group_lifecycle_policy.minimal", "group_lifetime_in_days", "180"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group_lifecycle_policy.minimal", "managed_group_types", "All"),
				),
			},
		},
	})
}

// TestUnitGroupLifecyclePolicyResource_Create_Maximal tests the creation of a group lifecycle policy with maximal configuration
func TestUnitGroupLifecyclePolicyResource_Create_Maximal(t *testing.T) {
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
					testCheckExists("microsoft365_graph_beta_groups_group_lifecycle_policy.maximal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group_lifecycle_policy.maximal", "group_lifetime_in_days", "365"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group_lifecycle_policy.maximal", "managed_group_types", "Selected"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group_lifecycle_policy.maximal", "alternate_notification_emails", "admin@example.com;notifications@example.com"),
				),
			},
		},
	})
}

// TestUnitGroupLifecyclePolicyResource_Update_MinimalToMaximal tests updating from minimal to maximal configuration
func TestUnitGroupLifecyclePolicyResource_Update_MinimalToMaximal(t *testing.T) {
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
					testCheckExists("microsoft365_graph_beta_groups_group_lifecycle_policy.minimal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group_lifecycle_policy.minimal", "group_lifetime_in_days", "180"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group_lifecycle_policy.minimal", "managed_group_types", "All"),
					// Verify minimal config doesn't have these attributes
					resource.TestCheckNoResourceAttr("microsoft365_graph_beta_groups_group_lifecycle_policy.minimal", "alternate_notification_emails"),
				),
			},
			// Update to maximal configuration (with the same resource name)
			{
				Config: testConfigMinimalToMaximal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_groups_group_lifecycle_policy.minimal"),
					// Now check that it has maximal attributes
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group_lifecycle_policy.minimal", "group_lifetime_in_days", "365"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group_lifecycle_policy.minimal", "managed_group_types", "Selected"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group_lifecycle_policy.minimal", "alternate_notification_emails", "admin@example.com;notifications@example.com"),
				),
			},
		},
	})
}

// TestUnitGroupLifecyclePolicyResource_Update_MaximalToMinimal tests updating from maximal to minimal configuration
func TestUnitGroupLifecyclePolicyResource_Update_MaximalToMinimal(t *testing.T) {
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
				Config: testConfigMaximal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_groups_group_lifecycle_policy.maximal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group_lifecycle_policy.maximal", "group_lifetime_in_days", "365"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group_lifecycle_policy.maximal", "managed_group_types", "Selected"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group_lifecycle_policy.maximal", "alternate_notification_emails", "admin@example.com;notifications@example.com"),
				),
			},
			// Update to minimal configuration (with the same resource name)
			{
				Config: testConfigMaximalToMinimal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_groups_group_lifecycle_policy.maximal"),
					// Now check that it has minimal attributes
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group_lifecycle_policy.maximal", "group_lifetime_in_days", "180"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group_lifecycle_policy.maximal", "managed_group_types", "All"),
					// Verify minimal config doesn't have these attributes
					resource.TestCheckNoResourceAttr("microsoft365_graph_beta_groups_group_lifecycle_policy.maximal", "alternate_notification_emails"),
				),
			},
		},
	})
}

// TestUnitGroupLifecyclePolicyResource_Delete_Minimal tests deleting a group lifecycle policy with minimal configuration
func TestUnitGroupLifecyclePolicyResource_Delete_Minimal(t *testing.T) {
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
					testCheckExists("microsoft365_graph_beta_groups_group_lifecycle_policy.minimal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group_lifecycle_policy.minimal", "group_lifetime_in_days", "180"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group_lifecycle_policy.minimal", "managed_group_types", "All"),
				),
			},
			// Delete the resource (by providing empty config)
			{
				Config: `# Empty config for deletion test`,
				Check: resource.ComposeTestCheckFunc(
					// Verify the resource is gone
					func(s *terraform.State) error {
						if _, ok := s.RootModule().Resources["microsoft365_graph_beta_groups_group_lifecycle_policy.minimal"]; ok {
							fmt.Printf("[DEBUG] State after deletion: %+v\n", s.RootModule().Resources)
							return fmt.Errorf("resource 'microsoft365_graph_beta_groups_group_lifecycle_policy.minimal' still exists in state after deletion")
						}
						return nil
					},
				),
			},
		},
	})
}

// TestUnitGroupLifecyclePolicyResource_Delete_Maximal tests deleting a group lifecycle policy with maximal configuration
func TestUnitGroupLifecyclePolicyResource_Delete_Maximal(t *testing.T) {
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
					testCheckExists("microsoft365_graph_beta_groups_group_lifecycle_policy.maximal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group_lifecycle_policy.maximal", "group_lifetime_in_days", "365"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group_lifecycle_policy.maximal", "managed_group_types", "Selected"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group_lifecycle_policy.maximal", "alternate_notification_emails", "admin@example.com;notifications@example.com"),
				),
			},
			// Delete the resource (by providing empty config)
			{
				Config: `# Empty config for deletion test`,
				Check: resource.ComposeTestCheckFunc(
					// Verify the resource is gone
					func(s *terraform.State) error {
						if _, ok := s.RootModule().Resources["microsoft365_graph_beta_groups_group_lifecycle_policy.maximal"]; ok {
							fmt.Printf("[DEBUG] State after deletion: %+v\n", s.RootModule().Resources)
							return fmt.Errorf("resource 'microsoft365_graph_beta_groups_group_lifecycle_policy.maximal' still exists in state after deletion")
						}
						return nil
					},
				),
			},
		},
	})
}

// TestUnitGroupLifecyclePolicyResource_Import tests importing an existing group lifecycle policy
func TestUnitGroupLifecyclePolicyResource_Import(t *testing.T) {
	// Set up mock environment
	_, _ = setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	// Set up the test environment
	setupTestEnvironment(t)

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Import the resource
			{
				ResourceName:  "microsoft365_graph_beta_groups_group_lifecycle_policy.imported",
				ImportState:   true,
				ImportStateId: "test-policy-id",
				Config:        testConfigImport(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_groups_group_lifecycle_policy.imported"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group_lifecycle_policy.imported", "id", "test-policy-id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group_lifecycle_policy.imported", "group_lifetime_in_days", "180"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_groups_group_lifecycle_policy.imported", "managed_group_types", "All"),
				),
			},
		},
	})
}

// TestUnitGroupLifecyclePolicyResource_Error tests error handling
func TestUnitGroupLifecyclePolicyResource_Error(t *testing.T) {
	// Set up mock environment with error responses
	_, policyMock := setupMockEnvironment()
	policyMock.RegisterErrorMocks()
	defer httpmock.DeactivateAndReset()

	// Set up the test environment
	setupTestEnvironment(t)

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testConfigError(),
				ExpectError: regexp.MustCompile(`.*`), // Expect any error
			},
		},
	})
}
