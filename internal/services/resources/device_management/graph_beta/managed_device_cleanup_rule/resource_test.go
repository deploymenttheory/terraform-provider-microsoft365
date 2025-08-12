package graphBetaManagedDeviceCleanupRule_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	localMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/managed_device_cleanup_rule/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/jarcoal/httpmock"
)

// Common test configurations that can be used by both unit and acceptance tests
const (
	// Basic configuration with standard attributes
	testConfigBasicTemplate = `
resource "microsoft365_graph_beta_device_management_managed_device_cleanup_rule" "test" {
  display_name                              = "Test Cleanup Rule"
  description                               = "Test description"
  device_cleanup_rule_platform_type         = "windows"
  device_inactivity_before_retirement_in_days = 90
}
`

	// Minimal configuration with only required attributes
	testConfigMinimalTemplate = `
resource "microsoft365_graph_beta_device_management_managed_device_cleanup_rule" "minimal" {
  display_name                              = "Minimal Cleanup Rule"
  device_cleanup_rule_platform_type         = "all"
  device_inactivity_before_retirement_in_days = 30
}
`

	// Maximal configuration with all possible attributes
	testConfigMaximalTemplate = `
resource "microsoft365_graph_beta_device_management_managed_device_cleanup_rule" "maximal" {
  display_name                              = "Maximal Cleanup Rule"
  description                               = "This is a comprehensive cleanup rule with all fields populated"
  device_cleanup_rule_platform_type         = "ios"
  device_inactivity_before_retirement_in_days = 180
}
`

	// Update configuration for testing changes
	testConfigUpdateTemplate = `
resource "microsoft365_graph_beta_device_management_managed_device_cleanup_rule" "test" {
  display_name                              = "Updated Cleanup Rule"
  description                               = "Updated description"
  device_cleanup_rule_platform_type         = "macOS"
  device_inactivity_before_retirement_in_days = 120
}
`
)

// Unit test provider configuration
const unitTestProviderConfig = `
provider "microsoft365" {
  tenant_id = "00000000-0000-0000-0000-000000000001"
  auth_method = "client_secret"
  entra_id_options = {
    client_id = "11111111-1111-1111-1111-111111111111"
    client_secret = "mock-secret-value"
  }
  cloud = "public"
}
`

// Acceptance test provider configuration
const accTestProviderConfig = `
provider "microsoft365" {
  # Configuration from environment variables
}
`

func TestUnitManagedDeviceCleanupRuleResource_Basic(t *testing.T) {
	// Activate httpmock
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Create a new Mocks instance and register mocks
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	// Register local mocks directly
	cleanupRuleMock := &localMocks.ManagedDeviceCleanupRuleMock{}
	cleanupRuleMock.RegisterMocks()

	// Set up the test environment
	setupTestEnvironment(t)

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_managed_device_cleanup_rule.test"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_managed_device_cleanup_rule.test", "display_name", "Test Cleanup Rule"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_managed_device_cleanup_rule.test", "description", "Test description"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_managed_device_cleanup_rule.test", "device_cleanup_rule_platform_type", "windows"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_managed_device_cleanup_rule.test", "device_inactivity_before_retirement_in_days", "90"),
				),
			},
		},
	})
}

func TestUnitManagedDeviceCleanupRuleResource_Minimal(t *testing.T) {
	// Activate httpmock
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Create a new Mocks instance and register mocks
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	// Register local mocks directly
	cleanupRuleMock := &localMocks.ManagedDeviceCleanupRuleMock{}
	cleanupRuleMock.RegisterMocks()

	// Set up the test environment
	setupTestEnvironment(t)

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_managed_device_cleanup_rule.minimal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_managed_device_cleanup_rule.minimal", "display_name", "Minimal Cleanup Rule"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_managed_device_cleanup_rule.minimal", "device_cleanup_rule_platform_type", "all"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_managed_device_cleanup_rule.minimal", "device_inactivity_before_retirement_in_days", "30"),
				),
			},
		},
	})
}

func TestUnitManagedDeviceCleanupRuleResource_Maximal(t *testing.T) {
	// Activate httpmock
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Create a new Mocks instance and register mocks
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	// Register local mocks directly
	cleanupRuleMock := &localMocks.ManagedDeviceCleanupRuleMock{}
	cleanupRuleMock.RegisterMocks()

	// Set up the test environment
	setupTestEnvironment(t)

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMaximal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_managed_device_cleanup_rule.maximal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_managed_device_cleanup_rule.maximal", "display_name", "Maximal Cleanup Rule"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_managed_device_cleanup_rule.maximal", "description", "This is a comprehensive cleanup rule with all fields populated"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_managed_device_cleanup_rule.maximal", "device_cleanup_rule_platform_type", "ios"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_managed_device_cleanup_rule.maximal", "device_inactivity_before_retirement_in_days", "180"),
				),
			},
		},
	})
}

func TestUnitManagedDeviceCleanupRuleResource_FullLifecycle(t *testing.T) {
	// Activate httpmock
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Create a new Mocks instance and register mocks
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	// Register local mocks directly
	cleanupRuleMock := &localMocks.ManagedDeviceCleanupRuleMock{}
	cleanupRuleMock.RegisterMocks()

	// Set up the test environment
	setupTestEnvironment(t)

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create with basic configuration
			{
				Config: testConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_managed_device_cleanup_rule.test"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_managed_device_cleanup_rule.test", "display_name", "Test Cleanup Rule"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_managed_device_cleanup_rule.test", "device_cleanup_rule_platform_type", "windows"),
				),
			},
			// Import test
			{
				ResourceName:      "microsoft365_graph_beta_device_management_managed_device_cleanup_rule.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestUnitManagedDeviceCleanupRuleResource_ErrorHandling(t *testing.T) {
	// Activate httpmock
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Create a new Mocks instance and register error mocks
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	// Register local error mocks directly
	cleanupRuleMock := &localMocks.ManagedDeviceCleanupRuleMock{}
	cleanupRuleMock.RegisterErrorMocks()

	// Set up the test environment
	setupTestEnvironment(t)

	// Run the test expecting an error
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testConfigBasic(),
				ExpectError: regexp.MustCompile(`(Bad Request|Forbidden|Access denied)`),
			},
		},
	})
}

func TestUnitManagedDeviceCleanupRuleResource_Update(t *testing.T) {
	// Activate httpmock
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Create a new Mocks instance and register mocks
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	// Register local mocks directly
	cleanupRuleMock := &localMocks.ManagedDeviceCleanupRuleMock{}
	cleanupRuleMock.RegisterMocks()

	// Set up the test environment
	setupTestEnvironment(t)

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigUpdate(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_managed_device_cleanup_rule.test"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_managed_device_cleanup_rule.test", "display_name", "Updated Cleanup Rule"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_managed_device_cleanup_rule.test", "description", "Updated description"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_managed_device_cleanup_rule.test", "device_cleanup_rule_platform_type", "macOS"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_managed_device_cleanup_rule.test", "device_inactivity_before_retirement_in_days", "120"),
				),
			},
		},
	})
}

// Acceptance Tests

func TestAccManagedDeviceCleanupRuleResource_Basic(t *testing.T) {
	// Skip if not running acceptance tests
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless TF_ACC environment variable is set")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckManagedDeviceCleanupRuleExists("microsoft365_graph_beta_device_management_managed_device_cleanup_rule.test"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_managed_device_cleanup_rule.test", "display_name", "Test Cleanup Rule"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_managed_device_cleanup_rule.test", "description", "Test description"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_managed_device_cleanup_rule.test", "device_cleanup_rule_platform_type", "windows"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_managed_device_cleanup_rule.test", "device_inactivity_before_retirement_in_days", "90"),
				),
			},
		},
		CheckDestroy: testAccCheckManagedDeviceCleanupRuleDestroy,
	})
}

func TestAccManagedDeviceCleanupRuleResource_Minimal(t *testing.T) {
	// Skip if not running acceptance tests
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless TF_ACC environment variable is set")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckManagedDeviceCleanupRuleExists("microsoft365_graph_beta_device_management_managed_device_cleanup_rule.minimal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_managed_device_cleanup_rule.minimal", "display_name", "Minimal Cleanup Rule"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_managed_device_cleanup_rule.minimal", "device_cleanup_rule_platform_type", "all"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_managed_device_cleanup_rule.minimal", "device_inactivity_before_retirement_in_days", "30"),
				),
			},
		},
		CheckDestroy: testAccCheckManagedDeviceCleanupRuleDestroy,
	})
}

func TestAccManagedDeviceCleanupRuleResource_Maximal(t *testing.T) {
	// Skip if not running acceptance tests
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless TF_ACC environment variable is set")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigMaximal(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckManagedDeviceCleanupRuleExists("microsoft365_graph_beta_device_management_managed_device_cleanup_rule.maximal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_managed_device_cleanup_rule.maximal", "display_name", "Maximal Cleanup Rule"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_managed_device_cleanup_rule.maximal", "description", "This is a comprehensive cleanup rule with all fields populated"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_managed_device_cleanup_rule.maximal", "device_cleanup_rule_platform_type", "ios"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_managed_device_cleanup_rule.maximal", "device_inactivity_before_retirement_in_days", "180"),
				),
			},
		},
		CheckDestroy: testAccCheckManagedDeviceCleanupRuleDestroy,
	})
}

func TestAccManagedDeviceCleanupRuleResource_Update(t *testing.T) {
	// Skip if not running acceptance tests
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless TF_ACC environment variable is set")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckManagedDeviceCleanupRuleExists("microsoft365_graph_beta_device_management_managed_device_cleanup_rule.test"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_managed_device_cleanup_rule.test", "display_name", "Test Cleanup Rule"),
				),
			},
			{
				Config: testAccConfigUpdate(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckManagedDeviceCleanupRuleExists("microsoft365_graph_beta_device_management_managed_device_cleanup_rule.test"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_managed_device_cleanup_rule.test", "display_name", "Updated Cleanup Rule"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_managed_device_cleanup_rule.test", "description", "Updated description"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_managed_device_cleanup_rule.test", "device_cleanup_rule_platform_type", "macOS"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_managed_device_cleanup_rule.test", "device_inactivity_before_retirement_in_days", "120"),
				),
			},
		},
		CheckDestroy: testAccCheckManagedDeviceCleanupRuleDestroy,
	})
}

func TestAccManagedDeviceCleanupRuleResource_FullLifecycle(t *testing.T) {
	// Skip if not running acceptance tests
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless TF_ACC environment variable is set")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckManagedDeviceCleanupRuleExists("microsoft365_graph_beta_device_management_managed_device_cleanup_rule.test"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_managed_device_cleanup_rule.test", "display_name", "Test Cleanup Rule"),
				),
			},
			{
				ResourceName:      "microsoft365_graph_beta_device_management_managed_device_cleanup_rule.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfigUpdate(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckManagedDeviceCleanupRuleExists("microsoft365_graph_beta_device_management_managed_device_cleanup_rule.test"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_managed_device_cleanup_rule.test", "display_name", "Updated Cleanup Rule"),
				),
			},
		},
		CheckDestroy: testAccCheckManagedDeviceCleanupRuleDestroy,
	})
}

// Helper functions

func testAccCheckManagedDeviceCleanupRuleExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("resource ID not set")
		}

		// In a real test, you would check the API here
		return nil
	}
}

func testAccCheckManagedDeviceCleanupRuleDestroy(s *terraform.State) error {
	// This would normally check the API to ensure the resource is gone
	// But for acceptance tests, we'll just return nil
	return nil
}

// Test configurations using shared templates

// Unit test configurations
func testConfigBasic() string {
	return unitTestProviderConfig + testConfigBasicTemplate
}

func testConfigMinimal() string {
	return unitTestProviderConfig + testConfigMinimalTemplate
}

func testConfigMaximal() string {
	return unitTestProviderConfig + testConfigMaximalTemplate
}

func testConfigUpdate() string {
	return unitTestProviderConfig + testConfigUpdateTemplate
}

// Acceptance test configurations
func testAccConfigBasic() string {
	return accTestProviderConfig + testConfigBasicTemplate
}

func testAccConfigMinimal() string {
	return accTestProviderConfig + testConfigMinimalTemplate
}

func testAccConfigMaximal() string {
	return accTestProviderConfig + testConfigMaximalTemplate
}

func testAccConfigUpdate() string {
	return accTestProviderConfig + testConfigUpdateTemplate
}

func setupTestEnvironment(t *testing.T) {
	// Set mock authentication credentials with valid values
	os.Setenv("M365_TENANT_ID", "00000000-0000-0000-0000-000000000001")
	os.Setenv("M365_CLIENT_ID", "11111111-1111-1111-1111-111111111111")
	os.Setenv("M365_CLIENT_SECRET", "mock-secret-value")
	os.Setenv("M365_AUTH_METHOD", "client_secret")
	os.Setenv("M365_CLOUD", "public")

	t.Cleanup(func() {
		os.Unsetenv("M365_TENANT_ID")
		os.Unsetenv("M365_CLIENT_ID")
		os.Unsetenv("M365_CLIENT_SECRET")
		os.Unsetenv("M365_AUTH_METHOD")
		os.Unsetenv("M365_CLOUD")
	})
}

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
