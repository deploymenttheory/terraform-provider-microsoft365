package graphBetaWindowsUpdateRing_test

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	windowsUpdateRingMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/windows_update_ring/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

func TestMain(m *testing.M) {
	exitCode := m.Run()
	os.Exit(exitCode)
}

func setupMockEnvironment() (*windowsUpdateRingMocks.WindowsUpdateRingMock, *windowsUpdateRingMocks.WindowsUpdateRingMock) {
	httpmock.Activate()
	mock := &windowsUpdateRingMocks.WindowsUpdateRingMock{}
	errorMock := &windowsUpdateRingMocks.WindowsUpdateRingMock{}
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
	updated = strings.Replace(updated, "Test Maximal Windows Update Ring - Unique", "Test Maximal Windows Update Ring", 1)

	return updated
}

// Helper function to get minimal config with a custom resource name
func testConfigMinimalWithResourceName(resourceName string) string {
	return fmt.Sprintf(`resource "microsoft365_graph_beta_device_management_windows_update_ring" "%s" {
  display_name                             = "Test Minimal Windows Update Ring"
  microsoft_update_service_allowed         = true
  drivers_excluded                         = false
  quality_updates_deferral_period_in_days  = 0
  feature_updates_deferral_period_in_days  = 0
  allow_windows11_upgrade                  = true
  skip_checks_before_restart               = false
  automatic_update_mode                    = "userDefined"
  feature_updates_rollback_window_in_days  = 10
  
  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}`, resourceName)
}

// TestUnitWindowsUpdateRingResource_Create_Minimal tests the creation of a Windows update ring with minimal configuration
func TestUnitWindowsUpdateRingResource_Create_Minimal(t *testing.T) {
	// Set up mock environment
	_, _ = setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	// Set up the test environment
	setupTestEnvironment(t)

	// Register the mocks
	mock := &windowsUpdateRingMocks.WindowsUpdateRingMock{}
	mock.RegisterMocks()

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_windows_update_ring.minimal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_update_ring.minimal", "display_name", "Test Minimal Windows Update Ring - Unique"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_update_ring.minimal", "microsoft_update_service_allowed", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_update_ring.minimal", "drivers_excluded", "false"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_update_ring.minimal", "quality_updates_deferral_period_in_days", "0"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_update_ring.minimal", "feature_updates_deferral_period_in_days", "0"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_update_ring.minimal", "allow_windows11_upgrade", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_update_ring.minimal", "skip_checks_before_restart", "false"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_update_ring.minimal", "automatic_update_mode", "userDefined"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_update_ring.minimal", "feature_updates_rollback_window_in_days", "10"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_update_ring.minimal", "role_scope_tag_ids.#", "1"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_update_ring.minimal", "role_scope_tag_ids.0", "0"),
				),
			},
		},
	})
}

// TestUnitWindowsUpdateRingResource_Create_Maximal tests the creation of a Windows update ring with maximal configuration
func TestUnitWindowsUpdateRingResource_Create_Maximal(t *testing.T) {
	// Set up mock environment
	_, _ = setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	// Set up the test environment
	setupTestEnvironment(t)

	// Register the mocks
	mock := &windowsUpdateRingMocks.WindowsUpdateRingMock{}
	mock.RegisterMocks()

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMaximal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_windows_update_ring.maximal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_update_ring.maximal", "display_name", "Test Maximal Windows Update Ring - Unique"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_update_ring.maximal", "description", "Maximal Windows update ring for testing with all features"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_update_ring.maximal", "microsoft_update_service_allowed", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_update_ring.maximal", "drivers_excluded", "false"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_update_ring.maximal", "quality_updates_deferral_period_in_days", "7"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_update_ring.maximal", "feature_updates_deferral_period_in_days", "14"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_update_ring.maximal", "allow_windows11_upgrade", "false"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_update_ring.maximal", "skip_checks_before_restart", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_update_ring.maximal", "automatic_update_mode", "autoInstallAndRebootAtScheduledTime"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_update_ring.maximal", "business_ready_updates_only", "businessReadyOnly"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_update_ring.maximal", "delivery_optimization_mode", "httpWithPeeringNat"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_update_ring.maximal", "prerelease_features", "settingsOnly"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_update_ring.maximal", "update_weeks", "firstWeek"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_update_ring.maximal", "active_hours_start", "09:00:00"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_update_ring.maximal", "active_hours_end", "17:00:00"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_update_ring.maximal", "user_pause_access", "disabled"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_update_ring.maximal", "feature_updates_rollback_window_in_days", "10"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_update_ring.maximal", "engaged_restart_deadline_in_days", "3"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_update_ring.maximal", "role_scope_tag_ids.#", "2"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_update_ring.maximal", "role_scope_tag_ids.0", "0"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_update_ring.maximal", "role_scope_tag_ids.1", "1"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_update_ring.maximal", "deadline_settings.deadline_for_feature_updates_in_days", "7"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_update_ring.maximal", "assignments.0.type", "groupAssignmentTarget"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_update_ring.maximal", "assignments.0.group_id", "44444444-4444-4444-4444-444444444444"),
				),
			},
		},
	})
}

// TestUnitWindowsUpdateRingResource_Update_MinimalToMaximal tests updating from minimal to maximal configuration
func TestUnitWindowsUpdateRingResource_Update_MinimalToMaximal(t *testing.T) {
	// Set up mock environment
	_, _ = setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	// Set up the test environment
	setupTestEnvironment(t)

	// Register the mocks
	mock := &windowsUpdateRingMocks.WindowsUpdateRingMock{}
	mock.RegisterMocks()

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Start with minimal configuration
			{
				Config: testConfigMinimalWithResourceName("test"),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_windows_update_ring.test"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_update_ring.test", "display_name", "Test Minimal Windows Update Ring"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_update_ring.test", "automatic_update_mode", "userDefined"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_update_ring.test", "role_scope_tag_ids.#", "1"),
				),
			},
			// Update to maximal configuration (with the same resource name)
			{
				Config: testConfigMaximalWithResourceName("test"),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_windows_update_ring.test"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_update_ring.test", "display_name", "Test Maximal Windows Update Ring"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_update_ring.test", "automatic_update_mode", "autoInstallAndRebootAtScheduledTime"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_update_ring.test", "description", "Maximal Windows update ring for testing with all features"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_update_ring.test", "business_ready_updates_only", "businessReadyOnly"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_update_ring.test", "role_scope_tag_ids.#", "2"),
				),
			},
		},
	})
}

// TestUnitWindowsUpdateRingResource_Update_MaximalToMinimal tests updating from maximal to minimal configuration
func TestUnitWindowsUpdateRingResource_Update_MaximalToMinimal(t *testing.T) {
	// Set up mock environment
	_, _ = setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	// Set up the test environment
	setupTestEnvironment(t)

	// Register the mocks
	mock := &windowsUpdateRingMocks.WindowsUpdateRingMock{}
	mock.RegisterMocks()

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Start with maximal configuration
			{
				Config: testConfigMaximalWithResourceName("test"),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_windows_update_ring.test"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_update_ring.test", "display_name", "Test Maximal Windows Update Ring"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_update_ring.test", "automatic_update_mode", "autoInstallAndRebootAtScheduledTime"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_update_ring.test", "business_ready_updates_only", "businessReadyOnly"),
				),
			},
			// Update to minimal configuration (with the same resource name)
			{
				Config: testConfigMinimalWithResourceName("test"),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_windows_update_ring.test"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_update_ring.test", "display_name", "Test Minimal Windows Update Ring"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_update_ring.test", "automatic_update_mode", "userDefined"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_update_ring.test", "role_scope_tag_ids.#", "1"),
				),
			},
		},
	})
}

// TestUnitWindowsUpdateRingResource_Delete_Minimal tests deleting a Windows update ring with minimal configuration
func TestUnitWindowsUpdateRingResource_Delete_Minimal(t *testing.T) {
	// Set up mock environment
	_, _ = setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	// Set up the test environment
	setupTestEnvironment(t)

	// Register the mocks
	mock := &windowsUpdateRingMocks.WindowsUpdateRingMock{}
	mock.RegisterMocks()

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_windows_update_ring.minimal"),
				),
			},
		},
	})
}

// TestUnitWindowsUpdateRingResource_Delete_Maximal tests deleting a Windows update ring with maximal configuration
func TestUnitWindowsUpdateRingResource_Delete_Maximal(t *testing.T) {
	// Set up mock environment
	_, _ = setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	// Set up the test environment
	setupTestEnvironment(t)

	// Register the mocks
	mock := &windowsUpdateRingMocks.WindowsUpdateRingMock{}
	mock.RegisterMocks()

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMaximal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_windows_update_ring.maximal"),
				),
			},
		},
	})
}

// TestUnitWindowsUpdateRingResource_Import tests importing a Windows update ring
func TestUnitWindowsUpdateRingResource_Import(t *testing.T) {
	// Set up mock environment
	_, _ = setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	// Set up the test environment
	setupTestEnvironment(t)

	// Register the mocks
	mock := &windowsUpdateRingMocks.WindowsUpdateRingMock{}
	mock.RegisterMocks()

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_windows_update_ring.minimal"),
				),
			},
			{
				ResourceName:      "microsoft365_graph_beta_device_management_windows_update_ring.minimal",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// TestUnitWindowsUpdateRingResource_Error tests error handling
func TestUnitWindowsUpdateRingResource_Error(t *testing.T) {
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
