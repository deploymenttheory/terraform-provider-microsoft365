package graphBetaCloudPcProvisioningPolicy_test

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	provisioningPolicyMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/windows_365/graph_beta/cloud_pc_provisioning_policy/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

func TestMain(m *testing.M) {
	exitCode := m.Run()
	os.Exit(exitCode)
}

func setupMockEnvironment() (*provisioningPolicyMocks.CloudPcProvisioningPolicyMock, *provisioningPolicyMocks.CloudPcProvisioningPolicyMock) {
	httpmock.Activate()
	mock := &provisioningPolicyMocks.CloudPcProvisioningPolicyMock{}
	errorMock := &provisioningPolicyMocks.CloudPcProvisioningPolicyMock{}
	return mock, errorMock
}

func setupTestEnvironment(t *testing.T) {
	// Set up any test-specific environment variables or configurations here if needed
}

// testCheckExists is a basic check to ensure the resource exists in the state
func testCheckExists(resourceName string) resource.TestCheckFunc {
	return resource.TestCheckResourceAttrSet(resourceName, "id")
}

// testConfigMinimal returns the minimal configuration for testing
func testConfigMinimal() string {
	content, err := os.ReadFile(filepath.Join("mocks", "terraform", "resource_minimal.tf"))
	if err != nil {
		return ""
	}
	return string(content)
}

// testConfigMaximal returns the maximal configuration for testing
func testConfigMaximal() string {
	content, err := os.ReadFile(filepath.Join("mocks", "terraform", "resource_maximal.tf"))
	if err != nil {
		return ""
	}
	return string(content)
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
	
	// Fix the display name to match test expectations
	updated = strings.Replace(updated, "Test Maximal Provisioning Policy - Unique", "Test Maximal Provisioning Policy", 1)

	return updated
}

// Helper function to get minimal config with a custom resource name
func testConfigMinimalWithResourceName(resourceName string) string {
	return fmt.Sprintf(`resource "microsoft365_graph_beta_windows_365_cloud_pc_provisioning_policy" "%s" {
  display_name = "Test Minimal Provisioning Policy"
  image_id     = "microsoftwindowsdesktop_windows-ent-cpc_win11-23h2-ent-cpc"
  
  microsoft_managed_desktop = {
    # Uses default values: managed_type = "notManaged", profile = "4aa9b805-9494-4eed-a04b-ed51ec9e631e"
  }
  
  windows_setting = {
    locale = "en-US"
  }
  
  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}`, resourceName)
}

// TestUnitCloudPcProvisioningPolicyResource_Create_Minimal tests the creation of a provisioning policy with minimal configuration
func TestUnitCloudPcProvisioningPolicyResource_Create_Minimal(t *testing.T) {
	// Set up mock environment
	_, _ = setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	// Set up the test environment
	setupTestEnvironment(t)

	// Register the mocks
	mock := &provisioningPolicyMocks.CloudPcProvisioningPolicyMock{}
	mock.RegisterMocks()

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_windows_365_cloud_pc_provisioning_policy.minimal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_provisioning_policy.minimal", "display_name", "Test Minimal Provisioning Policy - Unique"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_provisioning_policy.minimal", "image_id", "microsoftwindowsdesktop_windows-ent-cpc_win11-23h2-ent-cpc"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_provisioning_policy.minimal", "windows_setting.locale", "en-US"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_provisioning_policy.minimal", "microsoft_managed_desktop.managed_type", "notManaged"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_provisioning_policy.minimal", "provisioning_type", "dedicated"), // Default value
				),
			},
		},
	})
}

// TestUnitCloudPcProvisioningPolicyResource_Create_Maximal tests the creation of a provisioning policy with maximal configuration
func TestUnitCloudPcProvisioningPolicyResource_Create_Maximal(t *testing.T) {
	// Set up mock environment
	_, _ = setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	// Set up the test environment
	setupTestEnvironment(t)

	// Register the mocks
	mock := &provisioningPolicyMocks.CloudPcProvisioningPolicyMock{}
	mock.RegisterMocks()

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMaximal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_windows_365_cloud_pc_provisioning_policy.maximal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_provisioning_policy.maximal", "display_name", "Test Maximal Provisioning Policy - Unique"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_provisioning_policy.maximal", "description", "Maximal policy for testing with all features"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_provisioning_policy.maximal", "cloud_pc_naming_template", "CPC-MAX-%USERNAME:5%-%RAND:5%"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_provisioning_policy.maximal", "provisioning_type", "dedicated"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_provisioning_policy.maximal", "image_id", "microsoftwindowsdesktop_windows-ent-cpc_win11-24H2-ent-cpc-m365"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_provisioning_policy.maximal", "image_type", "gallery"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_provisioning_policy.maximal", "enable_single_sign_on", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_provisioning_policy.maximal", "local_admin_enabled", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_provisioning_policy.maximal", "managed_by", "windows365"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_provisioning_policy.maximal", "windows_setting.locale", "en-US"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_provisioning_policy.maximal", "microsoft_managed_desktop.managed_type", "notManaged"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_provisioning_policy.maximal", "microsoft_managed_desktop.profile", "4aa9b805-9494-4eed-a04b-ed51ec9e631e"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_provisioning_policy.maximal", "domain_join_configurations.0.domain_join_type", "hybridAzureADJoin"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_provisioning_policy.maximal", "domain_join_configurations.0.region_group", "usWest"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_provisioning_policy.maximal", "autopatch.autopatch_group_id", "4aa9b805-9494-4eed-a04b-ed51ec9e631e"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_provisioning_policy.maximal", "autopilot_configuration.device_preparation_profile_id", "12345678-1234-1234-1234-123456789012"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_provisioning_policy.maximal", "autopilot_configuration.application_timeout_in_minutes", "60"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_provisioning_policy.maximal", "autopilot_configuration.on_failure_device_access_denied", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_provisioning_policy.maximal", "apply_to_existing_cloud_pcs.microsoft_entra_single_sign_on_for_all_devices", "false"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_provisioning_policy.maximal", "apply_to_existing_cloud_pcs.region_or_azure_network_connection_for_all_devices", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_provisioning_policy.maximal", "apply_to_existing_cloud_pcs.region_or_azure_network_connection_for_select_devices", "false"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_provisioning_policy.maximal", "assignments.0.type", "groupAssignmentTarget"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_provisioning_policy.maximal", "assignments.0.group_id", "44444444-4444-4444-4444-444444444444"),
				),
			},
		},
	})
}

// TestUnitCloudPcProvisioningPolicyResource_Update_MinimalToMaximal tests updating from minimal to maximal configuration
func TestUnitCloudPcProvisioningPolicyResource_Update_MinimalToMaximal(t *testing.T) {
	// Set up mock environment
	_, _ = setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	// Set up the test environment
	setupTestEnvironment(t)

	// Register the mocks
	mock := &provisioningPolicyMocks.CloudPcProvisioningPolicyMock{}
	mock.RegisterMocks()

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Start with minimal configuration
			{
				Config: testConfigMinimalWithResourceName("test"),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_windows_365_cloud_pc_provisioning_policy.test"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_provisioning_policy.test", "display_name", "Test Minimal Provisioning Policy"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_provisioning_policy.test", "enable_single_sign_on", "false"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_provisioning_policy.test", "local_admin_enabled", "false"),
				),
			},
			// Update to maximal configuration (with the same resource name)
			{
				Config: testConfigMaximalWithResourceName("test"),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_windows_365_cloud_pc_provisioning_policy.test"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_provisioning_policy.test", "display_name", "Test Maximal Provisioning Policy"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_provisioning_policy.test", "enable_single_sign_on", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_provisioning_policy.test", "local_admin_enabled", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_provisioning_policy.test", "grace_period_in_hours", "8"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_provisioning_policy.test", "microsoft_managed_desktop.managed_type", "notManaged"),
				),
			},
		},
	})
}

// TestUnitCloudPcProvisioningPolicyResource_Update_MaximalToMinimal tests updating from maximal to minimal configuration
func TestUnitCloudPcProvisioningPolicyResource_Update_MaximalToMinimal(t *testing.T) {
	// Set up mock environment
	_, _ = setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	// Set up the test environment
	setupTestEnvironment(t)

	// Register the mocks
	mock := &provisioningPolicyMocks.CloudPcProvisioningPolicyMock{}
	mock.RegisterMocks()

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Start with maximal configuration
			{
				Config: testConfigMaximalWithResourceName("test"),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_windows_365_cloud_pc_provisioning_policy.test"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_provisioning_policy.test", "display_name", "Test Maximal Provisioning Policy"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_provisioning_policy.test", "microsoft_managed_desktop.managed_type", "notManaged"),
				),
			},
			// Update to minimal configuration (with the same resource name)
			{
				Config: testConfigMinimalWithResourceName("test"),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_windows_365_cloud_pc_provisioning_policy.test"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_provisioning_policy.test", "display_name", "Test Minimal Provisioning Policy"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_provisioning_policy.test", "enable_single_sign_on", "false"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_provisioning_policy.test", "local_admin_enabled", "false"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_provisioning_policy.test", "grace_period_in_hours", "4"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_provisioning_policy.test", "microsoft_managed_desktop.managed_type", "notManaged"),
				),
			},
		},
	})
}

// TestUnitCloudPcProvisioningPolicyResource_Delete_Minimal tests deleting a provisioning policy with minimal configuration
func TestUnitCloudPcProvisioningPolicyResource_Delete_Minimal(t *testing.T) {
	// Set up mock environment
	_, _ = setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	// Set up the test environment
	setupTestEnvironment(t)

	// Register the mocks
	mock := &provisioningPolicyMocks.CloudPcProvisioningPolicyMock{}
	mock.RegisterMocks()

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_windows_365_cloud_pc_provisioning_policy.minimal"),
				),
			},
		},
	})
}

// TestUnitCloudPcProvisioningPolicyResource_Delete_Maximal tests deleting a provisioning policy with maximal configuration
func TestUnitCloudPcProvisioningPolicyResource_Delete_Maximal(t *testing.T) {
	// Set up mock environment
	_, _ = setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	// Set up the test environment
	setupTestEnvironment(t)

	// Register the mocks
	mock := &provisioningPolicyMocks.CloudPcProvisioningPolicyMock{}
	mock.RegisterMocks()

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMaximal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_windows_365_cloud_pc_provisioning_policy.maximal"),
				),
			},
		},
	})
}

// TestUnitCloudPcProvisioningPolicyResource_Import tests importing a provisioning policy
func TestUnitCloudPcProvisioningPolicyResource_Import(t *testing.T) {
	// Set up mock environment
	_, _ = setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	// Set up the test environment
	setupTestEnvironment(t)

	// Register the mocks
	mock := &provisioningPolicyMocks.CloudPcProvisioningPolicyMock{}
	mock.RegisterMocks()

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_windows_365_cloud_pc_provisioning_policy.minimal"),
				),
			},
			{
				ResourceName:      "microsoft365_graph_beta_windows_365_cloud_pc_provisioning_policy.minimal",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// TestUnitCloudPcProvisioningPolicyResource_Error tests error handling
func TestUnitCloudPcProvisioningPolicyResource_Error(t *testing.T) {
	// Set up mock environment
	_, errorMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	// Set up the test environment
	setupTestEnvironment(t)

	// Register the error mocks
	errorMock.RegisterErrorMocks()

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testConfigMinimal(),
				ExpectError: regexp.MustCompile("Validation error: Invalid display name"),
			},
		},
	})
}
