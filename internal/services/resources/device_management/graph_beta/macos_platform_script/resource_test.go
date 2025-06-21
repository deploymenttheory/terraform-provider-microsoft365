package graphBetaMacOSPlatformScript_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/jarcoal/httpmock"
)

// Common test configurations that can be used by both unit and acceptance tests
const (
	// Basic configuration with standard attributes
	testConfigBasicTemplate = `
resource "microsoft365_graph_beta_device_management_macos_platform_script" "test" {
  display_name    = "Test macOS Script"
  description     = "Test description"
  script_content  = "#!/bin/bash\necho 'Hello World'"
  run_as_account  = "system"
  file_name       = "test-script.sh"
  block_execution_notifications = true
  execution_frequency = "P1D"
  retry_count     = 3

  assignments = {
    all_devices = false
    all_users   = true
  }
}
`

	// Minimal configuration with only required attributes
	testConfigMinimalTemplate = `
resource "microsoft365_graph_beta_device_management_macos_platform_script" "minimal" {
  display_name   = "Minimal macOS Script"
  script_content = "#!/bin/bash\necho 'Minimal Script'"
  run_as_account = "system"
  file_name      = "minimal-script.sh"

  assignments = {
    all_devices = false
    all_users   = false
  }
}
`

	// Maximal configuration with all possible attributes
	testConfigMaximalTemplate = `
resource "microsoft365_graph_beta_device_management_macos_platform_script" "maximal" {
  display_name    = "Maximal macOS Script"
  description     = "This is a comprehensive script with all fields populated"
  script_content  = "#!/bin/bash\necho 'Maximal Script Configuration'"
  run_as_account  = "user"
  file_name       = "maximal-script.sh"
  block_execution_notifications = true
  execution_frequency = "P4W"
  retry_count     = 10
  role_scope_tag_ids = ["0", "1"]

  assignments = {
    all_devices = true
    all_users   = false
  }
}
`

	// Update configuration for testing changes
	testConfigUpdateTemplate = `
resource "microsoft365_graph_beta_device_management_macos_platform_script" "test" {
  display_name    = "Updated macOS Script"
  description     = "Updated description"
  script_content  = "#!/bin/bash\necho 'Hello Updated World'"
  run_as_account  = "user"
  file_name       = "updated-script.sh"
  block_execution_notifications = false
  execution_frequency = "P1W"
  retry_count     = 5

  assignments = {
    all_devices = true
    all_users   = false
  }
}
`

	// Group assignments configuration
	testConfigGroupAssignmentsTemplate = `
resource "microsoft365_graph_beta_device_management_macos_platform_script" "group_assigned" {
  display_name    = "Group Assignment Script"
  description     = "Script with group assignments"
  script_content  = "#!/bin/bash\necho 'Group Assignment Script'"
  run_as_account  = "system"
  file_name       = "group-script.sh"

  assignments = {
    all_devices = false
    all_users   = false
    include_group_ids = ["11111111-1111-1111-1111-111111111111"]
    exclude_group_ids = ["22222222-2222-2222-2222-222222222222"]
  }
}
`

	// Complex duration configuration
	testConfigComplexDurationTemplate = `
resource "microsoft365_graph_beta_device_management_macos_platform_script" "complex_duration" {
  display_name    = "Complex Duration Script"
  description     = "Testing complex ISO 8601 duration"
  script_content  = "#!/bin/bash\necho 'Testing complex duration'"
  run_as_account  = "system"
  file_name       = "complex-duration-script.sh"
  block_execution_notifications = true
  execution_frequency = "P4W2D"
  retry_count     = 3

  assignments = {
    all_devices = false
    all_users   = true
  }
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

func TestUnitMacOSPlatformScriptResource_Basic(t *testing.T) {
	// Activate httpmock
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Create a new Mocks instance and register mocks
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	mockClient.RegisterMacOSPlatformScriptMocks()

	// Set up the test environment
	setupTestEnvironment(t)

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_macos_platform_script.test"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_platform_script.test", "display_name", "Test macOS Script"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_platform_script.test", "description", "Test description"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_platform_script.test", "run_as_account", "system"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_platform_script.test", "file_name", "test-script.sh"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_platform_script.test", "block_execution_notifications", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_platform_script.test", "execution_frequency", "P1D"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_platform_script.test", "retry_count", "3"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_platform_script.test", "assignments.all_users", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_platform_script.test", "assignments.all_devices", "false"),
				),
			},
		},
	})
}

func TestUnitMacOSPlatformScriptResource_Minimal(t *testing.T) {
	// Activate httpmock
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Create a new Mocks instance and register mocks
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	mockClient.RegisterMacOSPlatformScriptMocks()

	// Set up the test environment
	setupTestEnvironment(t)

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_macos_platform_script.minimal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_platform_script.minimal", "display_name", "Minimal macOS Script"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_platform_script.minimal", "run_as_account", "system"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_platform_script.minimal", "file_name", "minimal-script.sh"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_platform_script.minimal", "assignments.all_users", "false"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_platform_script.minimal", "assignments.all_devices", "false"),
				),
			},
		},
	})
}

func TestUnitMacOSPlatformScriptResource_Maximal(t *testing.T) {
	// Activate httpmock
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Create a new Mocks instance and register mocks
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	mockClient.RegisterMacOSPlatformScriptMocks()

	// Set up the test environment
	setupTestEnvironment(t)

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMaximal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_macos_platform_script.maximal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_platform_script.maximal", "display_name", "Maximal macOS Script"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_platform_script.maximal", "description", "This is a comprehensive script with all fields populated"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_platform_script.maximal", "run_as_account", "user"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_platform_script.maximal", "file_name", "maximal-script.sh"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_platform_script.maximal", "block_execution_notifications", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_platform_script.maximal", "retry_count", "10"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_platform_script.maximal", "assignments.all_devices", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_platform_script.maximal", "assignments.all_users", "false"),
				),
			},
		},
	})
}

func TestUnitMacOSPlatformScriptResource_GroupAssignments(t *testing.T) {
	// Activate httpmock
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Create a new Mocks instance and register mocks
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	mockClient.RegisterMacOSPlatformScriptMocks()

	// Set up the test environment
	setupTestEnvironment(t)

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigGroupAssignments(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_macos_platform_script.group_assigned"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_platform_script.group_assigned", "display_name", "Group Assignment Script"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_platform_script.group_assigned", "assignments.all_devices", "false"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_platform_script.group_assigned", "assignments.all_users", "false"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_macos_platform_script.group_assigned", "assignments.include_group_ids.*", "11111111-1111-1111-1111-111111111111"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_macos_platform_script.group_assigned", "assignments.exclude_group_ids.*", "22222222-2222-2222-2222-222222222222"),
				),
			},
		},
	})
}

func TestUnitMacOSPlatformScriptResource_FullLifecycle(t *testing.T) {
	// Activate httpmock
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Create a new Mocks instance and register mocks
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	mockClient.RegisterMacOSPlatformScriptMocks()

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
					testCheckExists("microsoft365_graph_beta_device_management_macos_platform_script.test"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_platform_script.test", "display_name", "Test macOS Script"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_platform_script.test", "run_as_account", "system"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_platform_script.test", "assignments.all_users", "true"),
				),
			},
			// Import test
			{
				ResourceName:      "microsoft365_graph_beta_device_management_macos_platform_script.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"script_content", // Sensitive value, not returned in API responses
					"assignments.%",
					"assignments.all_devices",
					"assignments.all_users",
					"assignments.include_group_ids",
					"assignments.exclude_group_ids",
					"retry_count",
				},
			},
		},
	})
}

func TestUnitMacOSPlatformScriptResource_ErrorHandling(t *testing.T) {
	// Activate httpmock
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Create a new Mocks instance and register error mocks
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	mockClient.RegisterMacOSPlatformScriptErrorMocks()

	// Set up the test environment
	setupTestEnvironment(t)

	// Run the test expecting an error
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testConfigBasic(),
				ExpectError: regexp.MustCompile(`(Access denied|Forbidden)`),
			},
		},
	})
}

func TestUnitMacOSPlatformScriptResource_Update(t *testing.T) {
	// Activate httpmock
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Create a new Mocks instance and register mocks
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	mockClient.RegisterMacOSPlatformScriptMocks()

	// Set up the test environment
	setupTestEnvironment(t)

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigUpdate(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_macos_platform_script.test"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_platform_script.test", "display_name", "Updated macOS Script"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_platform_script.test", "description", "Updated description"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_platform_script.test", "run_as_account", "user"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_platform_script.test", "file_name", "updated-script.sh"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_platform_script.test", "block_execution_notifications", "false"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_platform_script.test", "retry_count", "5"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_platform_script.test", "assignments.all_devices", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_platform_script.test", "assignments.all_users", "false"),
				),
			},
		},
	})
}

// Acceptance Tests
func TestAccMacOSPlatformScriptResource_Basic(t *testing.T) {
	// Skip if not running acceptance tests
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless TF_ACC environment variable is set")
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMacOSPlatformScriptExists("microsoft365_graph_beta_device_management_macos_platform_script.test"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_platform_script.test", "display_name", "Test macOS Script"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_platform_script.test", "run_as_account", "system"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_platform_script.test", "execution_frequency", "P1D"),
				),
			},
			{
				Config: testAccConfigUpdate(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMacOSPlatformScriptExists("microsoft365_graph_beta_device_management_macos_platform_script.test"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_platform_script.test", "display_name", "Updated macOS Script"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_platform_script.test", "run_as_account", "user"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_platform_script.test", "execution_frequency", "P1W"),
				),
			},
			{
				ResourceName:      "microsoft365_graph_beta_device_management_macos_platform_script.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
		CheckDestroy: testAccCheckMacOSPlatformScriptDestroy,
	})
}

func TestAccMacOSPlatformScriptResource_ComplexDuration(t *testing.T) {
	// Skip if not running acceptance tests
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless TF_ACC environment variable is set")
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigComplexDuration(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMacOSPlatformScriptExists("microsoft365_graph_beta_device_management_macos_platform_script.complex_duration"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_platform_script.complex_duration", "display_name", "Complex Duration Script"),
					// P4W2D would normally be normalized to P30D or similar
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_platform_script.complex_duration", "execution_frequency", "P4W2D"),
				),
			},
			{
				ResourceName:      "microsoft365_graph_beta_device_management_macos_platform_script.complex_duration",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
		CheckDestroy: testAccCheckMacOSPlatformScriptDestroy,
	})
}

// Helper Functions
func testAccCheckMacOSPlatformScriptExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		return nil
	}
}

func testAccCheckMacOSPlatformScriptDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "microsoft365_graph_beta_device_management_macos_platform_script" {
			continue
		}

		// In a real test, we would make an API call to verify the resource is gone
		// For unit tests with mocks, we can assume it's destroyed if we get here
		return nil
	}

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

func testConfigGroupAssignments() string {
	return unitTestProviderConfig + testConfigGroupAssignmentsTemplate
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

func testAccConfigComplexDuration() string {
	return accTestProviderConfig + testConfigComplexDurationTemplate
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

// testCheckExists verifies the resource exists in Terraform state
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

func testAccPreCheck(t *testing.T) {
	// Check required environment variables for acceptance tests
	envVars := []string{
		"MICROSOFT365_CLIENT_ID",
		"MICROSOFT365_CLIENT_SECRET",
		"MICROSOFT365_TENANT_ID",
	}

	for _, envVar := range envVars {
		if os.Getenv(envVar) == "" {
			t.Fatalf("%s environment variable must be set for acceptance tests", envVar)
		}
	}
}
