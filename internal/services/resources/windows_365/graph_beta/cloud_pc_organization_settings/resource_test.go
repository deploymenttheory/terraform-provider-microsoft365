package graphBetaCloudPcOrganizationSettings_test

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	organizationSettingsMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/windows_365/graph_beta/cloud_pc_organization_settings/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

func TestMain(m *testing.M) {
	exitCode := m.Run()
	os.Exit(exitCode)
}

func setupMockEnvironment() (*organizationSettingsMocks.OrganizationSettingsMock, *organizationSettingsMocks.OrganizationSettingsMock) {
	httpmock.Activate()
	mock := &organizationSettingsMocks.OrganizationSettingsMock{}
	errorMock := &organizationSettingsMocks.OrganizationSettingsMock{}
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

	return updated
}

// Helper function to get minimal config with a custom resource name
func testConfigMinimalWithResourceName(resourceName string) string {
	return fmt.Sprintf(`resource "microsoft365_graph_beta_windows_365_cloud_pc_organization_settings" "%s" {
  enable_mem_auto_enroll = false
  enable_single_sign_on  = false
  
  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}`, resourceName)
}

// TestUnitOrganizationSettingsResource_Create_Minimal tests the creation of organization settings with minimal configuration
func TestUnitOrganizationSettingsResource_Create_Minimal(t *testing.T) {
	// Set up mock environment
	_, _ = setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	// Set up the test environment
	setupTestEnvironment(t)

	// Register the mocks
	mock := &organizationSettingsMocks.OrganizationSettingsMock{}
	mock.RegisterMocks()

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_windows_365_cloud_pc_organization_settings.minimal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_organization_settings.minimal", "enable_mem_auto_enroll", "false"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_organization_settings.minimal", "enable_single_sign_on", "false"),
				),
			},
		},
	})
}

// TestUnitOrganizationSettingsResource_Create_Maximal tests the creation of organization settings with maximal configuration
func TestUnitOrganizationSettingsResource_Create_Maximal(t *testing.T) {
	// Set up mock environment
	_, _ = setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	// Set up the test environment
	setupTestEnvironment(t)

	// Register the mocks
	mock := &organizationSettingsMocks.OrganizationSettingsMock{}
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
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_organization_settings.maximal", "user_account_type", "administrator"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_organization_settings.maximal", "windows_settings.language", "en-GB"),
				),
			},
		},
	})
}

// TestUnitOrganizationSettingsResource_Update_MinimalToMaximal tests updating from minimal to maximal configuration
func TestUnitOrganizationSettingsResource_Update_MinimalToMaximal(t *testing.T) {
	// Set up mock environment
	_, _ = setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	// Set up the test environment
	setupTestEnvironment(t)

	// Register the mocks
	mock := &organizationSettingsMocks.OrganizationSettingsMock{}
	mock.RegisterMocks()

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Start with minimal configuration
			{
				Config: testConfigMinimalWithResourceName("test"),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_windows_365_cloud_pc_organization_settings.test"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_organization_settings.test", "enable_mem_auto_enroll", "false"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_organization_settings.test", "enable_single_sign_on", "false"),
				),
			},
			// Update to maximal configuration (with the same resource name)
			{
				Config: testConfigMaximalWithResourceName("test"),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_windows_365_cloud_pc_organization_settings.test"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_organization_settings.test", "enable_mem_auto_enroll", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_organization_settings.test", "enable_single_sign_on", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_organization_settings.test", "os_version", "windows11"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_organization_settings.test", "user_account_type", "administrator"),
				),
			},
		},
	})
}

// TestUnitOrganizationSettingsResource_Update_MaximalToMinimal tests updating from maximal to minimal configuration
func TestUnitOrganizationSettingsResource_Update_MaximalToMinimal(t *testing.T) {
	// Set up mock environment
	_, _ = setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	// Set up the test environment
	setupTestEnvironment(t)

	// Register the mocks
	mock := &organizationSettingsMocks.OrganizationSettingsMock{}
	mock.RegisterMocks()

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Start with maximal configuration
			{
				Config: testConfigMaximalWithResourceName("test"),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_windows_365_cloud_pc_organization_settings.test"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_organization_settings.test", "enable_mem_auto_enroll", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_organization_settings.test", "os_version", "windows11"),
				),
			},
			// Update to minimal configuration (with the same resource name)
			{
				Config: testConfigMinimalWithResourceName("test"),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_windows_365_cloud_pc_organization_settings.test"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_organization_settings.test", "enable_mem_auto_enroll", "false"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_organization_settings.test", "enable_single_sign_on", "false"),
				),
			},
		},
	})
}

// TestUnitOrganizationSettingsResource_Read tests reading the organization settings
func TestUnitOrganizationSettingsResource_Read(t *testing.T) {
	// Set up mock environment
	_, _ = setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	// Set up the test environment
	setupTestEnvironment(t)

	// Register the mocks
	mock := &organizationSettingsMocks.OrganizationSettingsMock{}
	mock.RegisterMocks()

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_windows_365_cloud_pc_organization_settings.minimal"),
				),
			},
		},
	})
}

// TestUnitOrganizationSettingsResource_Import tests importing organization settings
func TestUnitOrganizationSettingsResource_Import(t *testing.T) {
	// Set up mock environment
	_, _ = setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	// Set up the test environment
	setupTestEnvironment(t)

	// Register the mocks
	mock := &organizationSettingsMocks.OrganizationSettingsMock{}
	mock.RegisterMocks()

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_windows_365_cloud_pc_organization_settings.minimal"),
				),
			},
			{
				ResourceName:      "microsoft365_graph_beta_windows_365_cloud_pc_organization_settings.minimal",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// TestUnitOrganizationSettingsResource_Error tests error handling
func TestUnitOrganizationSettingsResource_Error(t *testing.T) {
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
				ExpectError: regexp.MustCompile("Validation error: Invalid settings"),
			},
		},
	})
}