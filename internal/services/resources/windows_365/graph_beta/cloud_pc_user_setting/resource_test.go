package graphBetaCloudPcUserSetting_test

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	userSettingMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/windows_365/graph_beta/cloud_pc_user_setting/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

func setupMockEnvironment() (*userSettingMocks.UserSettingMock, *userSettingMocks.UserSettingMock) {
	httpmock.Activate()
	mock := &userSettingMocks.UserSettingMock{}
	errorMock := &userSettingMocks.UserSettingMock{}
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
	return fmt.Sprintf(`resource "microsoft365_graph_beta_windows_365_user_setting" "%s" {
  display_name         = "Test Minimal User Setting"
  local_admin_enabled  = false
  reset_enabled        = false
  self_service_enabled = false
  
  restore_point_setting = {
    frequency_in_hours   = 12
    frequency_type       = "default"
    user_restore_enabled = false
  }
  
  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}`, resourceName)
}

// TestUnitUserSettingResource_Create_Minimal tests the creation of a user setting with minimal configuration
func TestUnitUserSettingResource_Create_Minimal(t *testing.T) {
	// Set up mock environment
	_, _ = setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	// Set up the test environment
	setupTestEnvironment(t)

	// Register the mocks
	mock := &userSettingMocks.UserSettingMock{}
	mock.RegisterMocks()

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_windows_365_user_setting.minimal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_user_setting.minimal", "display_name", "Test Minimal User Setting"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_user_setting.minimal", "local_admin_enabled", "false"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_user_setting.minimal", "reset_enabled", "false"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_user_setting.minimal", "self_service_enabled", "false"),
				),
			},
		},
	})
}

// TestUnitUserSettingResource_Create_Maximal tests the creation of a user setting with maximal configuration
func TestUnitUserSettingResource_Create_Maximal(t *testing.T) {
	// Set up mock environment
	_, _ = setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	// Set up the test environment
	setupTestEnvironment(t)

	// Register the mocks
	mock := &userSettingMocks.UserSettingMock{}
	mock.RegisterMocks()

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMaximal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_windows_365_user_setting.maximal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_user_setting.maximal", "display_name", "Test Maximal User Setting"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_user_setting.maximal", "local_admin_enabled", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_user_setting.maximal", "reset_enabled", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_user_setting.maximal", "self_service_enabled", "false"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_user_setting.maximal", "restore_point_setting.frequency_in_hours", "12"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_user_setting.maximal", "restore_point_setting.frequency_type", "default"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_user_setting.maximal", "restore_point_setting.user_restore_enabled", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_user_setting.maximal", "notification_setting.restart_prompts_disabled", "false"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_user_setting.maximal", "cross_region_disaster_recovery_setting.maintain_cross_region_restore_point_enabled", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_user_setting.maximal", "cross_region_disaster_recovery_setting.user_initiated_disaster_recovery_allowed", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_user_setting.maximal", "cross_region_disaster_recovery_setting.disaster_recovery_type", "premium"),
				),
			},
		},
	})
}

// TestUnitUserSettingResource_Update_MinimalToMaximal tests updating from minimal to maximal configuration
func TestUnitUserSettingResource_Update_MinimalToMaximal(t *testing.T) {
	// Set up mock environment
	_, _ = setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	// Set up the test environment
	setupTestEnvironment(t)

	// Register the mocks
	mock := &userSettingMocks.UserSettingMock{}
	mock.RegisterMocks()

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Start with minimal configuration
			{
				Config: testConfigMinimalWithResourceName("test"),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_windows_365_user_setting.test"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_user_setting.test", "display_name", "Test Minimal User Setting"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_user_setting.test", "local_admin_enabled", "false"),
				),
			},
			// Update to maximal configuration (with the same resource name)
			{
				Config: testConfigMaximalWithResourceName("test"),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_windows_365_user_setting.test"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_user_setting.test", "display_name", "Test Maximal User Setting"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_user_setting.test", "local_admin_enabled", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_user_setting.test", "reset_enabled", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_user_setting.test", "restore_point_setting.frequency_in_hours", "12"),
				),
			},
		},
	})
}

// TestUnitUserSettingResource_Update_MaximalToMinimal tests updating from maximal to minimal configuration
func TestUnitUserSettingResource_Update_MaximalToMinimal(t *testing.T) {
	// Set up mock environment
	_, _ = setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	// Set up the test environment
	setupTestEnvironment(t)

	// Register the mocks
	mock := &userSettingMocks.UserSettingMock{}
	mock.RegisterMocks()

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Start with maximal configuration
			{
				Config: testConfigMaximalWithResourceName("test"),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_windows_365_user_setting.test"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_user_setting.test", "display_name", "Test Maximal User Setting"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_user_setting.test", "restore_point_setting.frequency_in_hours", "12"),
				),
			},
			// Update to minimal configuration (with the same resource name)
			{
				Config: testConfigMinimalWithResourceName("test"),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_windows_365_user_setting.test"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_user_setting.test", "display_name", "Test Minimal User Setting"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_user_setting.test", "local_admin_enabled", "false"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_user_setting.test", "reset_enabled", "false"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_user_setting.test", "self_service_enabled", "false"),
				),
			},
		},
	})
}

// TestUnitUserSettingResource_Delete_Minimal tests deleting a user setting with minimal configuration
func TestUnitUserSettingResource_Delete_Minimal(t *testing.T) {
	// Set up mock environment
	_, _ = setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	// Set up the test environment
	setupTestEnvironment(t)

	// Register the mocks
	mock := &userSettingMocks.UserSettingMock{}
	mock.RegisterMocks()

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_windows_365_user_setting.minimal"),
				),
			},
		},
	})
}

// TestUnitUserSettingResource_Delete_Maximal tests deleting a user setting with maximal configuration
func TestUnitUserSettingResource_Delete_Maximal(t *testing.T) {
	// Set up mock environment
	_, _ = setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	// Set up the test environment
	setupTestEnvironment(t)

	// Register the mocks
	mock := &userSettingMocks.UserSettingMock{}
	mock.RegisterMocks()

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMaximal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_windows_365_user_setting.maximal"),
				),
			},
		},
	})
}

// TestUnitUserSettingResource_Import tests importing a user setting
func TestUnitUserSettingResource_Import(t *testing.T) {
	// Set up mock environment
	_, _ = setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	// Set up the test environment
	setupTestEnvironment(t)

	// Register the mocks
	mock := &userSettingMocks.UserSettingMock{}
	mock.RegisterMocks()

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_windows_365_user_setting.minimal"),
				),
			},
			{
				ResourceName:      "microsoft365_graph_beta_windows_365_user_setting.minimal",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// TestUnitUserSettingResource_Error tests error handling
func TestUnitUserSettingResource_Error(t *testing.T) {
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
