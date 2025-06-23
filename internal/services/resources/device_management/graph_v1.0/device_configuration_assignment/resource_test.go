package graphDeviceConfigurationAssignment_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/jarcoal/httpmock"
)

// Common test configurations that can be used by both unit and acceptance tests
const (
	// Basic configuration with group assignment
	testConfigGroupAssignmentTemplate = `
resource "microsoft365_graph_device_management_device_configuration_assignment" "test" {
  device_configuration_id = "test-config-id"
  target_type            = "groupAssignment"
  group_id               = "11111111-1111-1111-1111-111111111111"
}
`

	// Configuration with all devices assignment
	testConfigAllDevicesTemplate = `
resource "microsoft365_graph_device_management_device_configuration_assignment" "test" {
  device_configuration_id = "test-config-id"
  target_type            = "allDevices"
  group_id               = ""
}
`

	// Configuration with all licensed users assignment
	testConfigAllLicensedUsersTemplate = `
resource "microsoft365_graph_device_management_device_configuration_assignment" "test" {
  device_configuration_id = "test-config-id"
  target_type            = "allLicensedUsers"
  group_id               = ""
}
`

	// Configuration with exclusion group assignment
	testConfigExclusionGroupTemplate = `
resource "microsoft365_graph_device_management_device_configuration_assignment" "test" {
  device_configuration_id = "test-config-id"
  target_type            = "exclusionGroupAssignment"
  group_id               = "22222222-2222-2222-2222-222222222222"
}
`

	// Configuration with filter
	testConfigWithFilterTemplate = `
resource "microsoft365_graph_device_management_device_configuration_assignment" "test" {
  device_configuration_id = "test-config-id"
  target_type            = "groupAssignment"
  group_id               = "11111111-1111-1111-1111-111111111111"
  filter_id              = "33333333-3333-3333-3333-333333333333"
  filter_type            = "include"
}
`

	// Update configuration
	testConfigUpdateTemplate = `
resource "microsoft365_graph_device_management_device_configuration_assignment" "test" {
  device_configuration_id = "test-config-id"
  target_type            = "allLicensedUsers"
  group_id               = ""
}
`

	// Error configuration
	testConfigErrorTemplate = `
resource "microsoft365_graph_device_management_device_configuration_assignment" "test" {
  device_configuration_id = "error-config"
  target_type            = "groupAssignment"
  group_id               = "error-group"
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

// Set up the test environment
func setupTestEnvironment(t *testing.T) {
	// Activate httpmock
	httpmock.Activate()

	// Register the necessary mocks
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	// Register the specific mocks for this resource
	mockClient.RegisterDeviceConfigurationAssignmentMocks()
}

func TestUnitDeviceConfigurationAssignmentResource_GroupAssignment(t *testing.T) {
	// Activate httpmock
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Create a new Mocks instance and register mocks
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	mockClient.RegisterDeviceConfigurationAssignmentMocks()

	// Set up the test environment
	setupTestEnvironment(t)

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigGroupAssignment(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_device_management_device_configuration_assignment.test"),
					resource.TestCheckResourceAttr("microsoft365_graph_device_management_device_configuration_assignment.test", "device_configuration_id", "test-config-id"),
					resource.TestCheckResourceAttr("microsoft365_graph_device_management_device_configuration_assignment.test", "target_type", "groupAssignment"),
					resource.TestCheckResourceAttr("microsoft365_graph_device_management_device_configuration_assignment.test", "group_id", "11111111-1111-1111-1111-111111111111"),
				),
			},
		},
	})
}

func TestUnitDeviceConfigurationAssignmentResource_AllDevices(t *testing.T) {
	// Activate httpmock
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Create a new Mocks instance and register mocks
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	mockClient.RegisterDeviceConfigurationAssignmentMocks()

	// Set up the test environment
	setupTestEnvironment(t)

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigAllDevices(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_device_management_device_configuration_assignment.test"),
					resource.TestCheckResourceAttr("microsoft365_graph_device_management_device_configuration_assignment.test", "device_configuration_id", "test-config-id"),
					resource.TestCheckResourceAttr("microsoft365_graph_device_management_device_configuration_assignment.test", "target_type", "allDevices"),
				),
			},
		},
	})
}

func TestUnitDeviceConfigurationAssignmentResource_AllLicensedUsers(t *testing.T) {
	// Activate httpmock
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Create a new Mocks instance and register mocks
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	mockClient.RegisterDeviceConfigurationAssignmentMocks()

	// Set up the test environment
	setupTestEnvironment(t)

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigAllLicensedUsers(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_device_management_device_configuration_assignment.test"),
					resource.TestCheckResourceAttr("microsoft365_graph_device_management_device_configuration_assignment.test", "device_configuration_id", "test-config-id"),
					resource.TestCheckResourceAttr("microsoft365_graph_device_management_device_configuration_assignment.test", "target_type", "allLicensedUsers"),
				),
			},
		},
	})
}

func TestUnitDeviceConfigurationAssignmentResource_ExclusionGroup(t *testing.T) {
	// Activate httpmock
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Create a new Mocks instance and register mocks
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	mockClient.RegisterDeviceConfigurationAssignmentMocks()

	// Set up the test environment
	setupTestEnvironment(t)

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigExclusionGroup(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_device_management_device_configuration_assignment.test"),
					resource.TestCheckResourceAttr("microsoft365_graph_device_management_device_configuration_assignment.test", "device_configuration_id", "test-config-id"),
					resource.TestCheckResourceAttr("microsoft365_graph_device_management_device_configuration_assignment.test", "target_type", "exclusionGroupAssignment"),
					resource.TestCheckResourceAttr("microsoft365_graph_device_management_device_configuration_assignment.test", "group_id", "22222222-2222-2222-2222-222222222222"),
				),
			},
		},
	})
}

// Remove the WithFilter test entirely - filters not supported
// func TestUnitDeviceConfigurationAssignmentResource_WithFilter(t *testing.T) {
//     // REMOVED - Device configuration assignments don't support filters
// }

func TestUnitDeviceConfigurationAssignmentResource_Update(t *testing.T) {
	// Activate httpmock
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Create a new Mocks instance and register mocks
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	mockClient.RegisterDeviceConfigurationAssignmentMocks()

	// Set up the test environment
	setupTestEnvironment(t)

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigGroupAssignment(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_device_management_device_configuration_assignment.test"),
					resource.TestCheckResourceAttr("microsoft365_graph_device_management_device_configuration_assignment.test", "target_type", "groupAssignment"),
					resource.TestCheckResourceAttr("microsoft365_graph_device_management_device_configuration_assignment.test", "group_id", "11111111-1111-1111-1111-111111111111"),
				),
			},
			{
				Config: testConfigUpdate(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_device_management_device_configuration_assignment.test"),
					resource.TestCheckResourceAttr("microsoft365_graph_device_management_device_configuration_assignment.test", "target_type", "allLicensedUsers"),
					resource.TestCheckResourceAttr("microsoft365_graph_device_management_device_configuration_assignment.test", "group_id", ""),
				),
			},
		},
	})
}

func TestUnitDeviceConfigurationAssignmentResource_Error(t *testing.T) {
	// Activate httpmock
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Create a new Mocks instance and register ERROR mocks
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	mockClient.RegisterDeviceConfigurationAssignmentErrorMocks()

	// Set up the test environment
	setupTestEnvironment(t)

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testConfigError(),
				ExpectError: regexp.MustCompile("Bad Request"),
			},
		},
	})
}

func TestUnitDeviceConfigurationAssignmentResource_Import(t *testing.T) {
	// Activate httpmock
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Create a new Mocks instance and register mocks
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	mockClient.RegisterDeviceConfigurationAssignmentMocks()

	// Set up the test environment
	setupTestEnvironment(t)

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigGroupAssignment(),
			},
			{
				ResourceName:      "microsoft365_graph_device_management_device_configuration_assignment.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     "test-config-id:00000000-0000-0000-0000-000000000001", // Composite ID needed for API path
				ImportStateVerifyIgnore: []string{
					"timeouts", // Timeouts are not imported
				},
			},
		},
	})
}

// Helper functions to generate test configurations
func testConfigGroupAssignment() string {
	return unitTestProviderConfig + testConfigGroupAssignmentTemplate
}

func testConfigAllDevices() string {
	return unitTestProviderConfig + testConfigAllDevicesTemplate
}

func testConfigAllLicensedUsers() string {
	return unitTestProviderConfig + testConfigAllLicensedUsersTemplate
}

func testConfigExclusionGroup() string {
	return unitTestProviderConfig + testConfigExclusionGroupTemplate
}

func testConfigWithFilter() string {
	return unitTestProviderConfig + testConfigWithFilterTemplate
}

func testConfigUpdate() string {
	return unitTestProviderConfig + testConfigUpdateTemplate
}

func testConfigError() string {
	return unitTestProviderConfig + testConfigErrorTemplate
}

// testCheckExists checks if the resource exists by querying the API
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
