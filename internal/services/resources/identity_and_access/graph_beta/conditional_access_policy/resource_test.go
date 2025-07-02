package graphBetaConditionalAccessPolicy_test

// import (
// 	"fmt"
// 	"os"
// 	"path/filepath"
// 	"regexp"
// 	"strings"
// 	"testing"

// 	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
// 	localMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/conditional_access_policy/mocks"
// 	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
// 	"github.com/hashicorp/terraform-plugin-testing/terraform"
// 	"github.com/jarcoal/httpmock"
// )

// // Helper functions to return the test configurations by reading from files
// func testConfigMinimal() string {
// 	content, err := os.ReadFile(filepath.Join("mocks", "terraform", "resource_minimal.tf"))
// 	if err != nil {
// 		return ""
// 	}
// 	return string(content)
// }

// func testConfigMaximal() string {
// 	content, err := os.ReadFile(filepath.Join("mocks", "terraform", "resource_maximal.tf"))
// 	if err != nil {
// 		return ""
// 	}
// 	return string(content)
// }

// func testConfigMinimalToMaximal() string {
// 	// For minimal to maximal test, we need to use the maximal config
// 	// but with the minimal resource name to simulate an update

// 	// Read the maximal config
// 	maximalContent, err := os.ReadFile(filepath.Join("mocks", "terraform", "resource_maximal.tf"))
// 	if err != nil {
// 		return ""
// 	}

// 	// Replace the resource name to match the minimal one
// 	updatedMaximal := strings.Replace(string(maximalContent), "maximal", "minimal", 1)

// 	// Replace the display name to indicate this is an update
// 	updatedMaximal = strings.Replace(updatedMaximal, "Comprehensive Security Policy - Maximal", "Comprehensive Security Policy - Updated from Minimal", 1)

// 	return updatedMaximal
// }

// func testConfigError() string {
// 	// Read the minimal config and modify for error scenario
// 	content, err := os.ReadFile(filepath.Join("mocks", "terraform", "resource_minimal.tf"))
// 	if err != nil {
// 		return ""
// 	}

// 	// Replace resource name and display name to create an error scenario
// 	updated := strings.Replace(string(content), "minimal", "error", 1)
// 	updated = strings.Replace(updated, "Block Legacy Authentication - Minimal", "Error Policy - Duplicate", 1)

// 	return updated
// }

// // Helper function to set up the test environment
// func setupTestEnvironment(t *testing.T) {
// 	// Set environment variables for testing
// 	os.Setenv("TF_ACC", "0")
// 	os.Setenv("MS365_TEST_MODE", "true")
// }

// // Helper function to set up the mock environment for unit tests
// func setupMockEnvironment() (*mocks.Mocks, *localMocks.ConditionalAccessPolicyMock) {
// 	// Activate httpmock
// 	httpmock.Activate()

// 	// Create a new Mocks instance and register authentication mocks
// 	mockClient := mocks.NewMocks()
// 	mockClient.AuthMocks.RegisterMocks()

// 	// Register local mocks directly
// 	policyMock := &localMocks.ConditionalAccessPolicyMock{}
// 	policyMock.RegisterMocks()

// 	return mockClient, policyMock
// }

// // Helper function to check if a resource exists
// func testCheckExists(resourceName string) resource.TestCheckFunc {
// 	return func(s *terraform.State) error {
// 		rs, ok := s.RootModule().Resources[resourceName]
// 		if !ok {
// 			return fmt.Errorf("resource not found: %s", resourceName)
// 		}

// 		if rs.Primary.ID == "" {
// 			return fmt.Errorf("resource ID not set")
// 		}

// 		return nil
// 	}
// }

// // Helper functions to generate configs with custom resource names
// func testConfigMaximalWithResourceName(resourceName string) string {
// 	content, err := os.ReadFile(filepath.Join("mocks", "terraform", "resource_maximal.tf"))
// 	if err != nil {
// 		return ""
// 	}
// 	return strings.Replace(string(content), "maximal", resourceName, 1)
// }

// func testConfigMinimalWithResourceName(resourceName string) string {
// 	content, err := os.ReadFile(filepath.Join("mocks", "terraform", "resource_minimal.tf"))
// 	if err != nil {
// 		return ""
// 	}
// 	return strings.Replace(string(content), "minimal", resourceName, 1)
// }

// // =============================================================================
// // UNIT TESTS
// // =============================================================================

// // TestUnitConditionalAccessPolicyResource_Create_Minimal tests the creation of a conditional access policy with minimal configuration
// func TestUnitConditionalAccessPolicyResource_Create_Minimal(t *testing.T) {
// 	// Set up mock environment
// 	_, _ = setupMockEnvironment()
// 	defer httpmock.DeactivateAndReset()

// 	// Set up the test environment
// 	setupTestEnvironment(t)

// 	// Run the test
// 	resource.UnitTest(t, resource.TestCase{
// 		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: testConfigMinimal(),
// 				Check: resource.ComposeTestCheckFunc(
// 					testCheckExists("microsoft365_graph_beta_identity_and_access_conditional_access_policy.minimal"),
// 					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.minimal", "display_name", "Block Legacy Authentication - Minimal"),
// 					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.minimal", "state", "enabled"),
// 					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.minimal", "conditions.client_app_types.#", "2"),
// 					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.minimal", "conditions.sign_in_risk_levels.#", "0"),
// 					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.minimal", "conditions.applications.include_applications.#", "1"),
// 					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.minimal", "conditions.applications.include_applications.0", "All"),
// 					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.minimal", "conditions.users.include_users.#", "1"),
// 					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.minimal", "conditions.users.include_users.0", "All"),
// 					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.minimal", "grant_controls.operator", "OR"),
// 					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.minimal", "grant_controls.built_in_controls.#", "1"),
// 					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.minimal", "grant_controls.built_in_controls.0", "block"),
// 				),
// 			},
// 		},
// 	})
// }

// // TestUnitConditionalAccessPolicyResource_Create_Maximal tests the creation of a conditional access policy with maximal configuration
// func TestUnitConditionalAccessPolicyResource_Create_Maximal(t *testing.T) {
// 	// Set up mock environment
// 	_, _ = setupMockEnvironment()
// 	defer httpmock.DeactivateAndReset()

// 	// Set up the test environment
// 	setupTestEnvironment(t)

// 	// Run the test
// 	resource.UnitTest(t, resource.TestCase{
// 		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: testConfigMaximal(),
// 				Check: resource.ComposeTestCheckFunc(
// 					testCheckExists("microsoft365_graph_beta_identity_and_access_conditional_access_policy.maximal"),
// 					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.maximal", "display_name", "Comprehensive Security Policy - Maximal"),
// 					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.maximal", "state", "enabled"),

// 					// Conditions
// 					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.maximal", "conditions.client_app_types.#", "1"),
// 					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.maximal", "conditions.client_app_types.0", "all"),
// 					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.maximal", "conditions.user_risk_levels.#", "1"),
// 					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.maximal", "conditions.user_risk_levels.0", "high"),
// 					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.maximal", "conditions.sign_in_risk_levels.#", "2"),
// 					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.maximal", "conditions.service_principal_risk_levels.#", "2"),

// 					// Applications
// 					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.maximal", "conditions.applications.include_applications.#", "1"),
// 					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.maximal", "conditions.applications.include_applications.0", "All"),
// 					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.maximal", "conditions.applications.exclude_applications.#", "1"),
// 					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.maximal", "conditions.applications.include_user_actions.#", "1"),
// 					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.maximal", "conditions.applications.application_filter.mode", "exclude"),

// 					// Users
// 					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.maximal", "conditions.users.include_users.#", "1"),
// 					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.maximal", "conditions.users.include_users.0", "All"),
// 					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.maximal", "conditions.users.exclude_users.#", "1"),
// 					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.maximal", "conditions.users.include_groups.#", "1"),
// 					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.maximal", "conditions.users.exclude_groups.#", "1"),
// 					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.maximal", "conditions.users.include_roles.#", "1"),
// 					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.maximal", "conditions.users.exclude_roles.#", "1"),

// 					// Platforms
// 					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.maximal", "conditions.platforms.include_platforms.#", "1"),
// 					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.maximal", "conditions.platforms.include_platforms.0", "all"),
// 					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.maximal", "conditions.platforms.exclude_platforms.#", "2"),

// 					// Locations
// 					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.maximal", "conditions.locations.include_locations.#", "1"),
// 					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.maximal", "conditions.locations.include_locations.0", "All"),
// 					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.maximal", "conditions.locations.exclude_locations.#", "2"),

// 					// Devices
// 					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.maximal", "conditions.devices.include_devices.#", "1"),
// 					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.maximal", "conditions.devices.exclude_devices.#", "1"),
// 					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.maximal", "conditions.devices.include_device_states.#", "1"),
// 					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.maximal", "conditions.devices.exclude_device_states.#", "1"),
// 					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.maximal", "conditions.devices.device_filter.mode", "include"),

// 					// Grant Controls
// 					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.maximal", "grant_controls.operator", "AND"),
// 					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.maximal", "grant_controls.built_in_controls.#", "3"),
// 					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.maximal", "grant_controls.custom_authentication_factors.#", "1"),
// 					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.maximal", "grant_controls.terms_of_use.#", "1"),
// 					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.maximal", "grant_controls.authentication_strength.display_name", "Multifactor authentication"),
// 					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.maximal", "grant_controls.authentication_strength.allowed_combinations.#", "7"),

// 					// Session Controls
// 					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.maximal", "session_controls.disable_resilience_defaults", "false"),
// 					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.maximal", "session_controls.application_enforced_restrictions.is_enabled", "true"),
// 					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.maximal", "session_controls.cloud_app_security.is_enabled", "true"),
// 					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.maximal", "session_controls.cloud_app_security.cloud_app_security_type", "mcasConfigured"),
// 					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.maximal", "session_controls.sign_in_frequency.is_enabled", "true"),
// 					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.maximal", "session_controls.sign_in_frequency.type", "hours"),
// 					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.maximal", "session_controls.sign_in_frequency.value", "4"),
// 					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.maximal", "session_controls.persistent_browser.is_enabled", "true"),
// 					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.maximal", "session_controls.persistent_browser.mode", "always"),
// 					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.maximal", "session_controls.continuous_access_evaluation.mode", "strict"),
// 					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.maximal", "session_controls.secure_sign_in_session.is_enabled", "true"),
// 				),
// 			},
// 		},
// 	})
// }

// // TestUnitConditionalAccessPolicyResource_Update_MinimalToMaximal tests updating from minimal to maximal configuration
// func TestUnitConditionalAccessPolicyResource_Update_MinimalToMaximal(t *testing.T) {
// 	// Set up mock environment
// 	_, _ = setupMockEnvironment()
// 	defer httpmock.DeactivateAndReset()

// 	// Set up the test environment
// 	setupTestEnvironment(t)

// 	// Run the test
// 	resource.UnitTest(t, resource.TestCase{
// 		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
// 		Steps: []resource.TestStep{
// 			// Start with minimal configuration
// 			{
// 				Config: testConfigMinimal(),
// 				Check: resource.ComposeTestCheckFunc(
// 					testCheckExists("microsoft365_graph_beta_identity_and_access_conditional_access_policy.minimal"),
// 					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.minimal", "display_name", "Block Legacy Authentication - Minimal"),
// 					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.minimal", "grant_controls.operator", "OR"),
// 					// Check that minimal config has empty risk levels
// 					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.minimal", "conditions.user_risk_levels.#", "0"),
// 					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.minimal", "conditions.sign_in_risk_levels.#", "0"),
// 				),
// 			},
// 			// Update to maximal configuration (with the same resource name)
// 			{
// 				Config: testConfigMinimalToMaximal(),
// 				Check: resource.ComposeTestCheckFunc(
// 					testCheckExists("microsoft365_graph_beta_identity_and_access_conditional_access_policy.minimal"),
// 					// Now check that it has maximal attributes
// 					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.minimal", "display_name", "Comprehensive Security Policy - Updated from Minimal"),
// 					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.minimal", "grant_controls.operator", "AND"),
// 					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.minimal", "conditions.user_risk_levels.#", "1"),
// 					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.minimal", "conditions.platforms.include_platforms.#", "1"),
// 					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.minimal", "session_controls.sign_in_frequency.is_enabled", "true"),
// 				),
// 			},
// 		},
// 	})
// }

// // TestUnitConditionalAccessPolicyResource_Update_MaximalToMinimal tests updating from maximal to minimal configuration
// func TestUnitConditionalAccessPolicyResource_Update_MaximalToMinimal(t *testing.T) {
// 	// Set up mock environment
// 	_, _ = setupMockEnvironment()
// 	defer httpmock.DeactivateAndReset()

// 	// Set up the test environment
// 	setupTestEnvironment(t)

// 	// Run the test
// 	resource.UnitTest(t, resource.TestCase{
// 		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
// 		Steps: []resource.TestStep{
// 			// Start with maximal configuration
// 			{
// 				Config: testConfigMaximalWithResourceName("maximal_to_minimal"),
// 				Check: resource.ComposeTestCheckFunc(
// 					testCheckExists("microsoft365_graph_beta_identity_and_access_conditional_access_policy.maximal_to_minimal"),
// 					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.maximal_to_minimal", "display_name", "Comprehensive Security Policy - Maximal"),
// 					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.maximal_to_minimal", "grant_controls.operator", "AND"),
// 					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.maximal_to_minimal", "conditions.user_risk_levels.#", "1"),
// 					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.maximal_to_minimal", "session_controls.sign_in_frequency.is_enabled", "true"),
// 				),
// 			},
// 			// Update to minimal configuration
// 			{
// 				Config: testConfigMinimalWithResourceName("maximal_to_minimal"),
// 				Check: resource.ComposeTestCheckFunc(
// 					testCheckExists("microsoft365_graph_beta_identity_and_access_conditional_access_policy.maximal_to_minimal"),
// 					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.maximal_to_minimal", "display_name", "Block Legacy Authentication - Minimal"),
// 					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.maximal_to_minimal", "grant_controls.operator", "OR"),
// 					// Verify complex attributes are removed/simplified
// 					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.maximal_to_minimal", "conditions.user_risk_levels.#", "0"),
// 					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.maximal_to_minimal", "conditions.sign_in_risk_levels.#", "0"),
// 				),
// 			},
// 		},
// 	})
// }

// // TestUnitConditionalAccessPolicyResource_Delete_Minimal tests the deletion of a minimal conditional access policy
// func TestUnitConditionalAccessPolicyResource_Delete_Minimal(t *testing.T) {
// 	// Set up mock environment
// 	_, _ = setupMockEnvironment()
// 	defer httpmock.DeactivateAndReset()

// 	// Set up the test environment
// 	setupTestEnvironment(t)

// 	// Run the test
// 	resource.UnitTest(t, resource.TestCase{
// 		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: testConfigMinimal(),
// 				Check: resource.ComposeTestCheckFunc(
// 					testCheckExists("microsoft365_graph_beta_identity_and_access_conditional_access_policy.minimal"),
// 					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.minimal", "display_name", "Block Legacy Authentication - Minimal"),
// 				),
// 			},
// 		},
// 		CheckDestroy: func(s *terraform.State) error {
// 			// Verify the resource was destroyed
// 			for _, rs := range s.RootModule().Resources {
// 				if rs.Type != "microsoft365_graph_beta_identity_and_access_conditional_access_policy" {
// 					continue
// 				}
// 				// In a real scenario, we would check if the resource still exists in the API
// 				// For mocks, we assume it's deleted if we reach this point
// 			}
// 			return nil
// 		},
// 	})
// }

// // TestUnitConditionalAccessPolicyResource_Delete_Maximal tests the deletion of a maximal conditional access policy
// func TestUnitConditionalAccessPolicyResource_Delete_Maximal(t *testing.T) {
// 	// Set up mock environment
// 	_, _ = setupMockEnvironment()
// 	defer httpmock.DeactivateAndReset()

// 	// Set up the test environment
// 	setupTestEnvironment(t)

// 	// Run the test
// 	resource.UnitTest(t, resource.TestCase{
// 		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: testConfigMaximal(),
// 				Check: resource.ComposeTestCheckFunc(
// 					testCheckExists("microsoft365_graph_beta_identity_and_access_conditional_access_policy.maximal"),
// 					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.maximal", "display_name", "Comprehensive Security Policy - Maximal"),
// 				),
// 			},
// 		},
// 		CheckDestroy: func(s *terraform.State) error {
// 			// Verify the resource was destroyed
// 			for _, rs := range s.RootModule().Resources {
// 				if rs.Type != "microsoft365_graph_beta_identity_and_access_conditional_access_policy" {
// 					continue
// 				}
// 			}
// 			return nil
// 		},
// 	})
// }

// // TestUnitConditionalAccessPolicyResource_Import tests the import functionality
// func TestUnitConditionalAccessPolicyResource_Import(t *testing.T) {
// 	// Set up mock environment
// 	_, _ = setupMockEnvironment()
// 	defer httpmock.DeactivateAndReset()

// 	// Set up the test environment
// 	setupTestEnvironment(t)

// 	// Run the test
// 	resource.UnitTest(t, resource.TestCase{
// 		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: testConfigMinimal(),
// 				Check: resource.ComposeTestCheckFunc(
// 					testCheckExists("microsoft365_graph_beta_identity_and_access_conditional_access_policy.minimal"),
// 				),
// 			},
// 			{
// 				ResourceName:      "microsoft365_graph_beta_identity_and_access_conditional_access_policy.minimal",
// 				ImportState:       true,
// 				ImportStateVerify: true,
// 				ImportStateId:     "minimal-policy-id-12345", // Use the predefined ID from mocks
// 				// Skip verification of computed fields that might be handled differently
// 				ImportStateVerifyIgnore: []string{"created_date_time", "modified_date_time"},
// 			},
// 		},
// 	})
// }

// // TestUnitConditionalAccessPolicyResource_Error tests error handling
// func TestUnitConditionalAccessPolicyResource_Error(t *testing.T) {
// 	// Set up mock environment with error mocks
// 	_, policyMock := setupMockEnvironment()
// 	policyMock.RegisterErrorMocks()
// 	defer httpmock.DeactivateAndReset()

// 	// Set up the test environment
// 	setupTestEnvironment(t)

// 	// Run the test
// 	resource.UnitTest(t, resource.TestCase{
// 		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
// 		Steps: []resource.TestStep{
// 			{
// 				Config:      testConfigError(),
// 				ExpectError: regexp.MustCompile("Conflict|already exists"),
// 			},
// 		},
// 	})
// }

// // =============================================================================
// // ACCEPTANCE TESTS
// // =============================================================================

// // TestAccConditionalAccessPolicyResource_Create_Minimal tests the creation of a conditional access policy with minimal configuration
// func TestAccConditionalAccessPolicyResource_Create_Minimal(t *testing.T) {
// 	// Skip if not running acceptance tests
// 	if os.Getenv("TF_ACC") == "" {
// 		t.Skip("Acceptance tests skipped unless TF_ACC=1")
// 	}

// 	// Get test domain from environment variable or skip
// 	testDomain := os.Getenv("TEST_DOMAIN")
// 	if testDomain == "" {
// 		t.Skip("TEST_DOMAIN environment variable must be set for acceptance tests")
// 	}

// 	resource.Test(t, resource.TestCase{
// 		PreCheck:                 func() { testAccPreCheck(t) },
// 		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: testConfigMinimal(),
// 				Check: resource.ComposeTestCheckFunc(
// 					testCheckExists("microsoft365_graph_beta_identity_and_access_conditional_access_policy.minimal"),
// 					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.minimal", "display_name", "Block Legacy Authentication - Minimal"),
// 					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.minimal", "state", "enabled"),
// 					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.minimal", "conditions.client_app_types.#", "2"),
// 					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.minimal", "conditions.applications.include_applications.#", "1"),
// 					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.minimal", "conditions.applications.include_applications.0", "All"),
// 					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.minimal", "conditions.users.include_users.#", "1"),
// 					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.minimal", "conditions.users.include_users.0", "All"),
// 					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.minimal", "grant_controls.operator", "OR"),
// 					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.minimal", "grant_controls.built_in_controls.#", "1"),
// 					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.minimal", "grant_controls.built_in_controls.0", "block"),
// 					// Verify computed fields are set
// 					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_identity_and_access_conditional_access_policy.minimal", "id"),
// 					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_identity_and_access_conditional_access_policy.minimal", "created_date_time"),
// 					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_identity_and_access_conditional_access_policy.minimal", "modified_date_time"),
// 				),
// 			},
// 		},
// 	})
// }

// // TestAccConditionalAccessPolicyResource_Create_Maximal tests the creation of a conditional access policy with maximal configuration
// func TestAccConditionalAccessPolicyResource_Create_Maximal(t *testing.T) {
// 	// Skip if not running acceptance tests
// 	if os.Getenv("TF_ACC") == "" {
// 		t.Skip("Acceptance tests skipped unless TF_ACC=1")
// 	}

// 	// Get test domain from environment variable or skip
// 	testDomain := os.Getenv("TEST_DOMAIN")
// 	if testDomain == "" {
// 		t.Skip("TEST_DOMAIN environment variable must be set for acceptance tests")
// 	}

// 	resource.Test(t, resource.TestCase{
// 		PreCheck:                 func() { testAccPreCheck(t) },
// 		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: testConfigMaximal(),
// 				Check: resource.ComposeTestCheckFunc(
// 					testCheckExists("microsoft365_graph_beta_identity_and_access_conditional_access_policy.maximal"),
// 					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.maximal", "display_name", "Comprehensive Security Policy - Maximal"),
// 					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.maximal", "state", "enabled"),
// 					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.maximal", "conditions.client_app_types.#", "1"),
// 					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.maximal", "conditions.user_risk_levels.#", "1"),
// 					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.maximal", "conditions.sign_in_risk_levels.#", "2"),
// 					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.maximal", "conditions.platforms.include_platforms.#", "1"),
// 					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.maximal", "conditions.locations.include_locations.#", "1"),
// 					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.maximal", "conditions.devices.device_filter.mode", "include"),
// 					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.maximal", "grant_controls.operator", "AND"),
// 					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.maximal", "grant_controls.built_in_controls.#", "2"),
// 					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.maximal", "grant_controls.authentication_strength.display_name", "Multifactor authentication"),
// 					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.maximal", "session_controls.sign_in_frequency.is_enabled", "true"),
// 					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.maximal", "session_controls.sign_in_frequency.type", "hours"),
// 					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.maximal", "session_controls.sign_in_frequency.value", "4"),
// 					// Verify computed fields are set
// 					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_identity_and_access_conditional_access_policy.maximal", "id"),
// 					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_identity_and_access_conditional_access_policy.maximal", "created_date_time"),
// 					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_identity_and_access_conditional_access_policy.maximal", "modified_date_time"),
// 				),
// 			},
// 		},
// 	})
// }

// // TestAccConditionalAccessPolicyResource_Update_MinimalToMaximal tests updating from minimal to maximal configuration
// func TestAccConditionalAccessPolicyResource_Update_MinimalToMaximal(t *testing.T) {
// 	// Skip if not running acceptance tests
// 	if os.Getenv("TF_ACC") == "" {
// 		t.Skip("Acceptance tests skipped unless TF_ACC=1")
// 	}

// 	// Get test domain from environment variable or skip
// 	testDomain := os.Getenv("TEST_DOMAIN")
// 	if testDomain == "" {
// 		t.Skip("TEST_DOMAIN environment variable must be set for acceptance tests")
// 	}

// 	resource.Test(t, resource.TestCase{
// 		PreCheck:                 func() { testAccPreCheck(t) },
// 		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
// 		Steps: []resource.TestStep{
// 			// Start with minimal configuration
// 			{
// 				Config: testConfigMinimal(),
// 				Check: resource.ComposeTestCheckFunc(
// 					testCheckExists("microsoft365_graph_beta_identity_and_access_conditional_access_policy.minimal"),
// 					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.minimal", "display_name", "Block Legacy Authentication - Minimal"),
// 					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.minimal", "grant_controls.operator", "OR"),
// 				),
// 			},
// 			// Update to maximal configuration (with the same resource name)
// 			{
// 				Config: testConfigMinimalToMaximal(),
// 				Check: resource.ComposeTestCheckFunc(
// 					testCheckExists("microsoft365_graph_beta_identity_and_access_conditional_access_policy.minimal"),
// 					// Now check that it has maximal attributes
// 					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.minimal", "display_name", "Comprehensive Security Policy - Updated from Minimal"),
// 					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.minimal", "grant_controls.operator", "AND"),
// 					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.minimal", "conditions.user_risk_levels.#", "1"),
// 					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.minimal", "conditions.platforms.include_platforms.#", "1"),
// 					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.minimal", "session_controls.sign_in_frequency.is_enabled", "true"),
// 				),
// 			},
// 		},
// 	})
// }

// // TestAccConditionalAccessPolicyResource_Update_MaximalToMinimal tests updating from maximal to minimal configuration
// func TestAccConditionalAccessPolicyResource_Update_MaximalToMinimal(t *testing.T) {

// 	// Skip if not running acceptance tests
// 	if os.Getenv("TF_ACC") == "" {
// 		t.Skip("Acceptance tests skipped unless TF_ACC=1")
// 	}

// 	// Get test domain from environment variable or skip
// 	testDomain := os.Getenv("TEST_DOMAIN")
// 	if testDomain == "" {
// 		t.Skip("TEST_DOMAIN environment variable must be set for acceptance tests")
// 	}
// 	resource.Test(t, resource.TestCase{
// 		PreCheck:                 func() { testAccPreCheck(t) },
// 		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
// 		Steps: []resource.TestStep{
// 			// Start with maximal configuration
// 			{
// 				Config: testConfigMaximalWithResourceName("maximal_to_minimal"),
// 				Check: resource.ComposeTestCheckFunc(
// 					testCheckExists("microsoft365_graph_beta_identity_and_access_conditional_access_policy.maximal_to_minimal"),
// 					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.maximal_to_minimal", "display_name", "Comprehensive Security Policy - Maximal"),
// 					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.maximal_to_minimal", "grant_controls.operator", "AND"),
// 					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.maximal_to_minimal", "conditions.user_risk_levels.#", "1"),
// 					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.maximal_to_minimal", "session_controls.sign_in_frequency.is_enabled", "true"),
// 				),
// 			},
// 			// Update to minimal configuration
// 			{
// 				Config: testConfigMinimalWithResourceName("maximal_to_minimal"),
// 				Check: resource.ComposeTestCheckFunc(
// 					testCheckExists("microsoft365_graph_beta_identity_and_access_conditional_access_policy.maximal_to_minimal"),
// 					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.maximal_to_minimal", "display_name", "Block Legacy Authentication - Minimal"),
// 					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.maximal_to_minimal", "grant_controls.operator", "OR"),
// 				),
// 			},
// 		},
// 	})
// }

// // TestAccConditionalAccessPolicyResource_Import tests the import functionality
// func TestAccConditionalAccessPolicyResource_Import(t *testing.T) {
// 	resource.Test(t, resource.TestCase{
// 		PreCheck:                 func() { testAccPreCheck(t) },
// 		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: testConfigMinimal(),
// 				Check: resource.ComposeTestCheckFunc(
// 					testCheckExists("microsoft365_graph_beta_identity_and_access_conditional_access_policy.minimal"),
// 				),
// 			},
// 			{
// 				ResourceName:      "microsoft365_graph_beta_identity_and_access_conditional_access_policy.minimal",
// 				ImportState:       true,
// 				ImportStateVerify: true,
// 			},
// 		},
// 	})
// }

// // testAccPreCheck validates the necessary test API credentials exist
// func testAccPreCheck(t *testing.T) {
// 	// Check for required environment variables for acceptance testing
// 	requiredEnvVars := []string{
// 		"MS365_TENANT_ID",
// 		"MS365_CLIENT_ID",
// 		"MS365_CLIENT_SECRET",
// 	}

// 	for _, envVar := range requiredEnvVars {
// 		if os.Getenv(envVar) == "" {
// 			t.Fatalf("%s must be set for acceptance tests", envVar)
// 		}
// 	}
// }
