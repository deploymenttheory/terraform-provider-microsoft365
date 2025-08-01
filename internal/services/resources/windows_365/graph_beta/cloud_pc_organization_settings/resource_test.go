package graphBetaCloudPcOrganizationSettings_test

import (
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	cloudPcOrganizationSettingsMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/windows_365/graph_beta/cloud_pc_organization_settings/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

func TestMain(m *testing.M) {
	exitCode := m.Run()
	os.Exit(exitCode)
}

func setupMockEnvironment() (*cloudPcOrganizationSettingsMocks.CloudPcOrganizationSettingsMock, *cloudPcOrganizationSettingsMocks.CloudPcOrganizationSettingsMock) {
	httpmock.Activate()
	mock := &cloudPcOrganizationSettingsMocks.CloudPcOrganizationSettingsMock{}
	errorMock := &cloudPcOrganizationSettingsMocks.CloudPcOrganizationSettingsMock{}
	return mock, errorMock
}

func setupTestEnvironment(t *testing.T) {
	// Set up any test-specific environment variables or configurations here if needed
}

// testCheckExists is a basic check to ensure the resource exists in the state
func testCheckExists(resourceName string) resource.TestCheckFunc {
	return resource.TestCheckResourceAttrSet(resourceName, "id")
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
	return `
resource "microsoft365_graph_beta_windows_365_cloud_pc_organization_settings" "` + resourceName + `" {
  enable_mem_auto_enroll = true
  enable_single_sign_on  = true
  os_version             = "windows11"
  user_account_type      = "standardUser"
  windows_settings = {
    language = "en-US"
  }
  
  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}
`
}

// TestUnitCloudPcOrganizationSettingsResource_Create_Maximal tests the creation of Cloud PC organization settings with maximal configuration
func TestUnitCloudPcOrganizationSettingsResource_Create_Maximal(t *testing.T) {
	// Set up mock environment
	_, _ = setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	// Set up the test environment
	setupTestEnvironment(t)

	// Register the mocks
	mock := &cloudPcOrganizationSettingsMocks.CloudPcOrganizationSettingsMock{}
	mock.RegisterMocks()

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMaximal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_windows_365_cloud_pc_organization_settings.maximal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_organization_settings.maximal", "enable_mem_auto_enroll", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_organization_settings.maximal", "enable_single_sign_on", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_organization_settings.maximal", "os_version", "windows11"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_organization_settings.maximal", "user_account_type", "standardUser"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_organization_settings.maximal", "windows_settings.language", "en-US"),
				),
			},
		},
	})
}

// TestUnitCloudPcOrganizationSettingsResource_Update tests updating the Cloud PC organization settings
func TestUnitCloudPcOrganizationSettingsResource_Update(t *testing.T) {
	// Set up mock environment
	_, _ = setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	// Set up the test environment
	setupTestEnvironment(t)

	// Register the mocks
	mock := &cloudPcOrganizationSettingsMocks.CloudPcOrganizationSettingsMock{}
	mock.RegisterMocks()

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Start with initial configuration
			{
				Config: testConfigMaximalWithResourceName("test"),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_windows_365_cloud_pc_organization_settings.test"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_organization_settings.test", "enable_mem_auto_enroll", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_organization_settings.test", "os_version", "windows11"),
				),
			},
			// Update configuration
			{
				Config: `
resource "microsoft365_graph_beta_windows_365_cloud_pc_organization_settings" "test" {
  enable_mem_auto_enroll = false
  enable_single_sign_on  = false
  os_version             = "windows10"
  user_account_type      = "administrator"
  windows_settings = {
    language = "fr-FR"
  }
  
  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}
`,
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_windows_365_cloud_pc_organization_settings.test"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_organization_settings.test", "enable_mem_auto_enroll", "false"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_organization_settings.test", "enable_single_sign_on", "false"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_organization_settings.test", "os_version", "windows10"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_organization_settings.test", "user_account_type", "administrator"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_organization_settings.test", "windows_settings.language", "fr-FR"),
				),
			},
		},
	})
}

// TestUnitCloudPcOrganizationSettingsResource_Delete tests deleting the Cloud PC organization settings
func TestUnitCloudPcOrganizationSettingsResource_Delete(t *testing.T) {
	// Set up mock environment
	_, _ = setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	// Set up the test environment
	setupTestEnvironment(t)

	// Register the mocks
	mock := &cloudPcOrganizationSettingsMocks.CloudPcOrganizationSettingsMock{}
	mock.RegisterMocks()

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMaximal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_windows_365_cloud_pc_organization_settings.maximal"),
				),
			},
		},
	})
}

// TestUnitCloudPcOrganizationSettingsResource_Import tests importing the Cloud PC organization settings
func TestUnitCloudPcOrganizationSettingsResource_Import(t *testing.T) {
	// Set up mock environment
	_, _ = setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	// Set up the test environment
	setupTestEnvironment(t)

	// Register the mocks
	mock := &cloudPcOrganizationSettingsMocks.CloudPcOrganizationSettingsMock{}
	mock.RegisterMocks()

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMaximal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_windows_365_cloud_pc_organization_settings.maximal"),
				),
			},
			{
				ResourceName:      "microsoft365_graph_beta_windows_365_cloud_pc_organization_settings.maximal",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// TestUnitCloudPcOrganizationSettingsResource_Error tests error handling
func TestUnitCloudPcOrganizationSettingsResource_Error(t *testing.T) {
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
				Config:      testConfigMaximal(),
				ExpectError: regexp.MustCompile("Validation error: Invalid OS version"),
			},
		},
	})
}