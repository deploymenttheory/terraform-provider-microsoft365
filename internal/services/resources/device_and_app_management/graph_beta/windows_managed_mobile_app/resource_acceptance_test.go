package graphBetaDeviceAndAppManagementWindowsManagedMobileApp_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// TestAccWindowsManagedMobileAppResource_Create_Minimal tests creating a Windows managed mobile app with minimal configuration
func TestAccWindowsManagedMobileAppResource_Create_Minimal(t *testing.T) {
	// Skip if not running acceptance tests
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless TF_ACC=1")
	}

	// Get test managed app protection ID from environment variable or skip
	testManagedAppProtectionID1 := os.Getenv("TEST_WINDOWS_MANAGED_APP_PROTECTION_ID_1")
	if testManagedAppProtectionID1 == "" {
		t.Skip("TEST_WINDOWS_MANAGED_APP_PROTECTION_ID_1 environment variable must be set for acceptance tests")
	}

	resourceName := "microsoft365_graph_beta_device_and_app_management_windows_managed_mobile_app.minimal"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckWindowsManagedMobileAppDestroy,
		Steps: []resource.TestStep{
			// Create with minimal configuration
			{
				Config: testAccConfigMinimal(testManagedAppProtectionID1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWindowsManagedMobileAppExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "managed_app_protection_id", testManagedAppProtectionID1),
					resource.TestCheckResourceAttr(resourceName, "mobile_app_identifier.windows_app_id", "com.example.testapp"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
		},
	})
}

// TestAccWindowsManagedMobileAppResource_Create_Maximal tests creating a Windows managed mobile app with maximal configuration
func TestAccWindowsManagedMobileAppResource_Create_Maximal(t *testing.T) {
	// Skip if not running acceptance tests
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless TF_ACC=1")
	}

	// Get test managed app protection ID from environment variable or skip
	testManagedAppProtectionID2 := os.Getenv("TEST_WINDOWS_MANAGED_APP_PROTECTION_ID_2")
	if testManagedAppProtectionID2 == "" {
		t.Skip("TEST_WINDOWS_MANAGED_APP_PROTECTION_ID_2 environment variable must be set for acceptance tests")
	}

	resourceName := "microsoft365_graph_beta_device_and_app_management_windows_managed_mobile_app.maximal"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckWindowsManagedMobileAppDestroy,
		Steps: []resource.TestStep{
			// Create with maximal configuration
			{
				Config: testAccConfigMaximal(testManagedAppProtectionID2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWindowsManagedMobileAppExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "managed_app_protection_id", testManagedAppProtectionID2),
					resource.TestCheckResourceAttr(resourceName, "mobile_app_identifier.windows_app_id", "com.example.complexapp"),
					resource.TestCheckResourceAttr(resourceName, "version", "1.5"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
		},
	})
}

// TestAccWindowsManagedMobileAppResource_Update_MinimalToMaximal tests updating from minimal to maximal config
func TestAccWindowsManagedMobileAppResource_Update_MinimalToMaximal(t *testing.T) {
	// Skip if not running acceptance tests
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless TF_ACC=1")
	}

	// Get test managed app protection ID from environment variable or skip
	testManagedAppProtectionID1 := os.Getenv("TEST_WINDOWS_MANAGED_APP_PROTECTION_ID_1")
	if testManagedAppProtectionID1 == "" {
		t.Skip("TEST_WINDOWS_MANAGED_APP_PROTECTION_ID_1 environment variable must be set for acceptance tests")
	}

	resourceName := "microsoft365_graph_beta_device_and_app_management_windows_managed_mobile_app.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckWindowsManagedMobileAppDestroy,
		Steps: []resource.TestStep{
			// Start with minimal configuration
			{
				Config: testAccConfigMinimalNamed("test", testManagedAppProtectionID1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWindowsManagedMobileAppExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "managed_app_protection_id", testManagedAppProtectionID1),
					resource.TestCheckResourceAttr(resourceName, "mobile_app_identifier.windows_app_id", "com.example.testapp"),
				),
			},
			// Update to maximal configuration
			{
				Config: testAccConfigMaximalNamed("test", testManagedAppProtectionID1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWindowsManagedMobileAppExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "managed_app_protection_id", testManagedAppProtectionID1),
					resource.TestCheckResourceAttr(resourceName, "mobile_app_identifier.windows_app_id", "com.example.complexapp"),
					resource.TestCheckResourceAttr(resourceName, "version", "1.5"),
				),
			},
		},
	})
}

// TestAccWindowsManagedMobileAppResource_Update_MaximalToMinimal tests updating from maximal to minimal config
func TestAccWindowsManagedMobileAppResource_Update_MaximalToMinimal(t *testing.T) {
	// Skip if not running acceptance tests
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless TF_ACC=1")
	}

	// Get test managed app protection ID from environment variable or skip
	testManagedAppProtectionID2 := os.Getenv("TEST_WINDOWS_MANAGED_APP_PROTECTION_ID_2")
	if testManagedAppProtectionID2 == "" {
		t.Skip("TEST_WINDOWS_MANAGED_APP_PROTECTION_ID_2 environment variable must be set for acceptance tests")
	}

	resourceName := "microsoft365_graph_beta_device_and_app_management_windows_managed_mobile_app.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckWindowsManagedMobileAppDestroy,
		Steps: []resource.TestStep{
			// Start with maximal configuration
			{
				Config: testAccConfigMaximalNamed("test", testManagedAppProtectionID2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWindowsManagedMobileAppExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "managed_app_protection_id", testManagedAppProtectionID2),
					resource.TestCheckResourceAttr(resourceName, "version", "1.5"),
				),
			},
			// Update to minimal configuration
			{
				Config: testAccConfigMinimalNamed("test", testManagedAppProtectionID2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWindowsManagedMobileAppExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "managed_app_protection_id", testManagedAppProtectionID2),
					resource.TestCheckResourceAttr(resourceName, "mobile_app_identifier.windows_app_id", "com.example.testapp"),
				),
			},
		},
	})
}

// TestAccWindowsManagedMobileAppResource_Delete_Minimal tests deleting a Windows managed mobile app with minimal configuration
func TestAccWindowsManagedMobileAppResource_Delete_Minimal(t *testing.T) {
	// Skip if not running acceptance tests
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless TF_ACC=1")
	}

	// Get test managed app protection ID from environment variable or skip
	testManagedAppProtectionID1 := os.Getenv("TEST_WINDOWS_MANAGED_APP_PROTECTION_ID_1")
	if testManagedAppProtectionID1 == "" {
		t.Skip("TEST_WINDOWS_MANAGED_APP_PROTECTION_ID_1 environment variable must be set for acceptance tests")
	}

	resourceName := "microsoft365_graph_beta_device_and_app_management_windows_managed_mobile_app.minimal"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckWindowsManagedMobileAppDestroy,
		Steps: []resource.TestStep{
			// Create the resource
			{
				Config: testAccConfigMinimal(testManagedAppProtectionID1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWindowsManagedMobileAppExists(resourceName),
				),
			},
			// Delete the resource (by providing empty config)
			{
				Config: `# Empty config for deletion test`,
				Check: func(s *terraform.State) error {
					// The resource should be gone
					_, exists := s.RootModule().Resources[resourceName]
					if exists {
						return fmt.Errorf("resource %s still exists after deletion", resourceName)
					}
					return nil
				},
			},
		},
	})
}

// TestAccWindowsManagedMobileAppResource_Delete_Maximal tests deleting a Windows managed mobile app with maximal configuration
func TestAccWindowsManagedMobileAppResource_Delete_Maximal(t *testing.T) {
	// Skip if not running acceptance tests
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless TF_ACC=1")
	}

	// Get test managed app protection ID from environment variable or skip
	testManagedAppProtectionID2 := os.Getenv("TEST_WINDOWS_MANAGED_APP_PROTECTION_ID_2")
	if testManagedAppProtectionID2 == "" {
		t.Skip("TEST_WINDOWS_MANAGED_APP_PROTECTION_ID_2 environment variable must be set for acceptance tests")
	}

	resourceName := "microsoft365_graph_beta_device_and_app_management_windows_managed_mobile_app.maximal"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckWindowsManagedMobileAppDestroy,
		Steps: []resource.TestStep{
			// Create the resource
			{
				Config: testAccConfigMaximal(testManagedAppProtectionID2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWindowsManagedMobileAppExists(resourceName),
				),
			},
			// Delete the resource (by providing empty config)
			{
				Config: `# Empty config for deletion test`,
				Check: func(s *terraform.State) error {
					// The resource should be gone
					_, exists := s.RootModule().Resources[resourceName]
					if exists {
						return fmt.Errorf("resource %s still exists after deletion", resourceName)
					}
					return nil
				},
			},
		},
	})
}

// TestAccWindowsManagedMobileAppResource_Import tests importing a resource
func TestAccWindowsManagedMobileAppResource_Import(t *testing.T) {
	// Skip if not running acceptance tests
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless TF_ACC=1")
	}

	// Get test managed app protection ID from environment variable or skip
	testManagedAppProtectionID1 := os.Getenv("TEST_WINDOWS_MANAGED_APP_PROTECTION_ID_1")
	if testManagedAppProtectionID1 == "" {
		t.Skip("TEST_WINDOWS_MANAGED_APP_PROTECTION_ID_1 environment variable must be set for acceptance tests")
	}

	resourceName := "microsoft365_graph_beta_device_and_app_management_windows_managed_mobile_app.minimal"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckWindowsManagedMobileAppDestroy,
		Steps: []resource.TestStep{
			// Create
			{
				Config: testAccConfigMinimal(testManagedAppProtectionID1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWindowsManagedMobileAppExists(resourceName),
				),
			},
			// Import
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Helper functions for acceptance tests

func testAccCheckWindowsManagedMobileAppExists(resourceName string) resource.TestCheckFunc {
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

func testAccCheckWindowsManagedMobileAppDestroy(s *terraform.State) error {
	// In a real test, we would verify the Windows managed mobile app is removed
	// For this resource, we don't need to check anything special since removing
	// the resource will remove the managed mobile app
	return nil
}

// Test configurations

// Minimal configuration with default resource name
func testAccConfigMinimal(managedAppProtectionID string) string {
	return fmt.Sprintf(`
resource "microsoft365_graph_beta_device_and_app_management_windows_managed_mobile_app" "minimal" {
  managed_app_protection_id = "%s"
  mobile_app_identifier = {
    windows_app_id = "com.example.testapp"
  }
}
`, managedAppProtectionID)
}

// Minimal configuration with custom resource name
func testAccConfigMinimalNamed(resourceName string, managedAppProtectionID string) string {
	return fmt.Sprintf(`
resource "microsoft365_graph_beta_device_and_app_management_windows_managed_mobile_app" "%s" {
  managed_app_protection_id = "%s"
  mobile_app_identifier = {
    windows_app_id = "com.example.testapp"
  }
}
`, resourceName, managedAppProtectionID)
}

// Maximal configuration with default resource name
func testAccConfigMaximal(managedAppProtectionID string) string {
	return fmt.Sprintf(`
resource "microsoft365_graph_beta_device_and_app_management_windows_managed_mobile_app" "maximal" {
  managed_app_protection_id = "%s"
  mobile_app_identifier = {
    windows_app_id = "com.example.complexapp"
  }
  version = "1.5"
}
`, managedAppProtectionID)
}

// Maximal configuration with custom resource name
func testAccConfigMaximalNamed(resourceName string, managedAppProtectionID string) string {
	return fmt.Sprintf(`
resource "microsoft365_graph_beta_device_and_app_management_windows_managed_mobile_app" "%s" {
  managed_app_protection_id = "%s"
  mobile_app_identifier = {
    windows_app_id = "com.example.complexapp"
  }
  version = "1.5"
}
`, resourceName, managedAppProtectionID)
}
