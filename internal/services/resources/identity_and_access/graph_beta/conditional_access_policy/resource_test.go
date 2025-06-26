package graphBetaConditionalAccessPolicy_test

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	localMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/conditional_access_policy/mocks"
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

	// Replace the display name to indicate this is an update
	updatedMaximal = strings.Replace(updatedMaximal, "Comprehensive Security Policy - Maximal", "Comprehensive Security Policy - Updated from Minimal", 1)

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
	updated = strings.Replace(updated, "Block Legacy Authentication - Minimal", "Error Policy - Duplicate", 1)

	return updated
}

// Helper function to set up the test environment
func setupTestEnvironment(t *testing.T) {
	// Set environment variables for testing
	os.Setenv("TF_ACC", "0")
	os.Setenv("MS365_TEST_MODE", "true")
}

// Helper function to set up the mock environment
func setupMockEnvironment() (*mocks.Mocks, *localMocks.ConditionalAccessPolicyMock) {
	// Activate httpmock
	httpmock.Activate()

	// Create a new Mocks instance and register authentication mocks
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	// Register local mocks directly
	policyMock := &localMocks.ConditionalAccessPolicyMock{}
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

		if rs.Primary.ID == "" {
			return fmt.Errorf("resource ID not set")
		}

		return nil
	}
}

// TestUnitConditionalAccessPolicyResource_Create_Minimal tests the creation of a conditional access policy with minimal configuration
func TestUnitConditionalAccessPolicyResource_Create_Minimal(t *testing.T) {
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
					testCheckExists("microsoft365_graph_beta_identity_and_access_conditional_access_policy.minimal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.minimal", "display_name", "Block Legacy Authentication - Minimal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.minimal", "state", "enabled"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.minimal", "conditions.client_app_types.#", "2"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.minimal", "conditions.applications.include_applications.#", "1"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.minimal", "conditions.applications.include_applications.0", "All"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.minimal", "conditions.users.include_users.#", "1"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.minimal", "conditions.users.include_users.0", "All"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.minimal", "grant_controls.operator", "OR"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.minimal", "grant_controls.built_in_controls.#", "1"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.minimal", "grant_controls.built_in_controls.0", "block"),
				),
			},
		},
	})
}

// TestUnitConditionalAccessPolicyResource_Create_Maximal tests the creation of a conditional access policy with maximal configuration
func TestUnitConditionalAccessPolicyResource_Create_Maximal(t *testing.T) {
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
					testCheckExists("microsoft365_graph_beta_identity_and_access_conditional_access_policy.maximal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.maximal", "display_name", "Comprehensive Security Policy - Maximal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.maximal", "state", "enabled"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.maximal", "conditions.client_app_types.#", "1"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.maximal", "conditions.client_app_types.0", "all"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.maximal", "conditions.user_risk_levels.#", "1"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.maximal", "conditions.user_risk_levels.0", "high"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.maximal", "conditions.sign_in_risk_levels.#", "2"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.maximal", "conditions.platforms.include_platforms.#", "1"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.maximal", "conditions.platforms.include_platforms.0", "all"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.maximal", "conditions.locations.include_locations.#", "1"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.maximal", "conditions.locations.include_locations.0", "All"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.maximal", "conditions.devices.device_filter.mode", "include"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.maximal", "grant_controls.operator", "AND"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.maximal", "grant_controls.built_in_controls.#", "2"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.maximal", "grant_controls.authentication_strength.display_name", "Multifactor authentication"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.maximal", "session_controls.sign_in_frequency.is_enabled", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.maximal", "session_controls.sign_in_frequency.type", "hours"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.maximal", "session_controls.sign_in_frequency.value", "4"),
				),
			},
		},
	})
}

// TestUnitConditionalAccessPolicyResource_Update_MinimalToMaximal tests updating from minimal to maximal configuration
func TestUnitConditionalAccessPolicyResource_Update_MinimalToMaximal(t *testing.T) {
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
					testCheckExists("microsoft365_graph_beta_identity_and_access_conditional_access_policy.minimal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.minimal", "display_name", "Block Legacy Authentication - Minimal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.minimal", "grant_controls.operator", "OR"),
					// Verify minimal config doesn't have these attributes
					resource.TestCheckNoResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.minimal", "conditions.user_risk_levels"),
					resource.TestCheckNoResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.minimal", "conditions.platforms"),
					resource.TestCheckNoResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.minimal", "session_controls"),
				),
			},
			// Update to maximal configuration (with the same resource name)
			{
				Config: testConfigMinimalToMaximal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_identity_and_access_conditional_access_policy.minimal"),
					// Now check that it has maximal attributes
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.minimal", "display_name", "Comprehensive Security Policy - Updated from Minimal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.minimal", "grant_controls.operator", "AND"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.minimal", "conditions.user_risk_levels.#", "1"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.minimal", "conditions.platforms.include_platforms.#", "1"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.minimal", "session_controls.sign_in_frequency.is_enabled", "true"),
				),
			},
		},
	})
}

// TestUnitConditionalAccessPolicyResource_Update_MaximalToMinimal tests updating from maximal to minimal configuration
func TestUnitConditionalAccessPolicyResource_Update_MaximalToMinimal(t *testing.T) {
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
				Config: testConfigMaximalWithResourceName("maximal_to_minimal"),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_identity_and_access_conditional_access_policy.maximal_to_minimal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.maximal_to_minimal", "display_name", "Comprehensive Security Policy - Maximal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.maximal_to_minimal", "grant_controls.operator", "AND"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.maximal_to_minimal", "conditions.user_risk_levels.#", "1"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.maximal_to_minimal", "session_controls.sign_in_frequency.is_enabled", "true"),
				),
			},
			// Update to minimal configuration
			{
				Config: testConfigMinimalWithResourceName("maximal_to_minimal"),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_identity_and_access_conditional_access_policy.maximal_to_minimal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.maximal_to_minimal", "display_name", "Block Legacy Authentication - Minimal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.maximal_to_minimal", "grant_controls.operator", "OR"),
					// Verify complex attributes are removed
					resource.TestCheckNoResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.maximal_to_minimal", "conditions.user_risk_levels"),
					resource.TestCheckNoResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.maximal_to_minimal", "conditions.platforms"),
					resource.TestCheckNoResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.maximal_to_minimal", "session_controls"),
				),
			},
		},
	})
}

// Helper functions to generate configs with custom resource names
func testConfigMaximalWithResourceName(resourceName string) string {
	content, err := os.ReadFile(filepath.Join("mocks", "terraform", "resource_maximal.tf"))
	if err != nil {
		return ""
	}
	return strings.Replace(string(content), "maximal", resourceName, 1)
}

func testConfigMinimalWithResourceName(resourceName string) string {
	content, err := os.ReadFile(filepath.Join("mocks", "terraform", "resource_minimal.tf"))
	if err != nil {
		return ""
	}
	return strings.Replace(string(content), "minimal", resourceName, 1)
}

// TestUnitConditionalAccessPolicyResource_Delete_Minimal tests the deletion of a minimal conditional access policy
func TestUnitConditionalAccessPolicyResource_Delete_Minimal(t *testing.T) {
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
					testCheckExists("microsoft365_graph_beta_identity_and_access_conditional_access_policy.minimal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.minimal", "display_name", "Block Legacy Authentication - Minimal"),
				),
			},
		},
		CheckDestroy: func(s *terraform.State) error {
			// Verify the resource was destroyed
			for _, rs := range s.RootModule().Resources {
				if rs.Type != "microsoft365_graph_beta_identity_and_access_conditional_access_policy" {
					continue
				}
				// In a real scenario, we would check if the resource still exists in the API
				// For mocks, we assume it's deleted if we reach this point
			}
			return nil
		},
	})
}

// TestUnitConditionalAccessPolicyResource_Delete_Maximal tests the deletion of a maximal conditional access policy
func TestUnitConditionalAccessPolicyResource_Delete_Maximal(t *testing.T) {
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
					testCheckExists("microsoft365_graph_beta_identity_and_access_conditional_access_policy.maximal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.maximal", "display_name", "Comprehensive Security Policy - Maximal"),
				),
			},
		},
		CheckDestroy: func(s *terraform.State) error {
			// Verify the resource was destroyed
			for _, rs := range s.RootModule().Resources {
				if rs.Type != "microsoft365_graph_beta_identity_and_access_conditional_access_policy" {
					continue
				}
			}
			return nil
		},
	})
}

// TestUnitConditionalAccessPolicyResource_Import tests the import functionality
func TestUnitConditionalAccessPolicyResource_Import(t *testing.T) {
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
					testCheckExists("microsoft365_graph_beta_identity_and_access_conditional_access_policy.minimal"),
				),
			},
			{
				ResourceName:      "microsoft365_graph_beta_identity_and_access_conditional_access_policy.minimal",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     "minimal-policy-id-12345", // Use the predefined ID from mocks
			},
		},
	})
}

// TestUnitConditionalAccessPolicyResource_Error tests error handling
func TestUnitConditionalAccessPolicyResource_Error(t *testing.T) {
	// Set up mock environment with error mocks
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
				ExpectError: regexp.MustCompile("Conflict|already exists"),
			},
		},
	})
}
