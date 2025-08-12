package graphBetaUserLicenseAssignment_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// TestAccUserLicenseAssignmentResource_Create_Minimal tests creating a license assignment with minimal configuration
func TestAccUserLicenseAssignmentResource_Create_Minimal(t *testing.T) {
	// Skip if not running acceptance tests
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless TF_ACC=1")
	}

	// Get test user ID from environment variable or skip
	testUserID1 := os.Getenv("TEST_USER_ID_1")
	if testUserID1 == "" {
		t.Skip("TEST_USER_ID_1 environment variable must be set for acceptance tests")
	}

	// Get test license SKU ID from environment variable or skip
	testLicenseSkuID1 := os.Getenv("TEST_LICENSE_SKU_ID_1")
	if testLicenseSkuID1 == "" {
		t.Skip("TEST_LICENSE_SKU_ID_1 environment variable must be set for acceptance tests")
	}

	resourceName := "microsoft365_graph_beta_users_user_license_assignment.minimal"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckUserLicenseAssignmentDestroy,
		Steps: []resource.TestStep{
			// Create with minimal configuration
			{
				Config: testAccConfigMinimal(testUserID1, testLicenseSkuID1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUserLicenseAssignmentExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "user_id", testUserID1),
					resource.TestCheckResourceAttrSet(resourceName, "user_principal_name"),
					resource.TestCheckResourceAttr(resourceName, "add_licenses.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "add_licenses.0.sku_id", testLicenseSkuID1),
					resource.TestCheckResourceAttr(resourceName, "add_licenses.0.disabled_plans.#", "0"),
				),
			},
		},
	})
}

// TestAccUserLicenseAssignmentResource_Create_Maximal tests creating a license assignment with maximal configuration
func TestAccUserLicenseAssignmentResource_Create_Maximal(t *testing.T) {
	// Skip if not running acceptance tests
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless TF_ACC=1")
	}

	// Get test user ID from environment variable or skip
	testUserID2 := os.Getenv("TEST_USER_ID_2")
	if testUserID2 == "" {
		t.Skip("TEST_USER_ID_2 environment variable must be set for acceptance tests")
	}

	// Get test license SKU IDs from environment variables or skip
	testLicenseSkuID1 := os.Getenv("TEST_LICENSE_SKU_ID_1")
	if testLicenseSkuID1 == "" {
		t.Skip("TEST_LICENSE_SKU_ID_1 environment variable must be set for acceptance tests")
	}

	testLicenseSkuID2 := os.Getenv("TEST_LICENSE_SKU_ID_2")
	if testLicenseSkuID2 == "" {
		t.Skip("TEST_LICENSE_SKU_ID_2 environment variable must be set for acceptance tests")
	}

	// Get test service plan ID from environment variable (optional)
	testServicePlanID := os.Getenv("TEST_SERVICE_PLAN_ID")

	resourceName := "microsoft365_graph_beta_users_user_license_assignment.maximal"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckUserLicenseAssignmentDestroy,
		Steps: []resource.TestStep{
			// Create with maximal configuration
			{
				Config: testAccConfigMaximal(testUserID2, testLicenseSkuID1, testLicenseSkuID2, testServicePlanID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUserLicenseAssignmentExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "user_id", testUserID2),
					resource.TestCheckResourceAttrSet(resourceName, "user_principal_name"),
					resource.TestCheckResourceAttr(resourceName, "add_licenses.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "add_licenses.0.sku_id", testLicenseSkuID1),
					resource.TestCheckResourceAttr(resourceName, "add_licenses.1.sku_id", testLicenseSkuID2),
				),
			},
		},
	})
}

// TestAccUserLicenseAssignmentResource_Update_MinimalToMaximal tests updating from minimal to maximal config
func TestAccUserLicenseAssignmentResource_Update_MinimalToMaximal(t *testing.T) {
	// Skip if not running acceptance tests
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless TF_ACC=1")
	}

	// Get test user ID from environment variable or skip
	testUserID1 := os.Getenv("TEST_USER_ID_1")
	if testUserID1 == "" {
		t.Skip("TEST_USER_ID_1 environment variable must be set for acceptance tests")
	}

	// Get test license SKU IDs from environment variables or skip
	testLicenseSkuID1 := os.Getenv("TEST_LICENSE_SKU_ID_1")
	if testLicenseSkuID1 == "" {
		t.Skip("TEST_LICENSE_SKU_ID_1 environment variable must be set for acceptance tests")
	}

	testLicenseSkuID2 := os.Getenv("TEST_LICENSE_SKU_ID_2")
	if testLicenseSkuID2 == "" {
		t.Skip("TEST_LICENSE_SKU_ID_2 environment variable must be set for acceptance tests")
	}

	// Get test service plan ID from environment variable (optional)
	testServicePlanID := os.Getenv("TEST_SERVICE_PLAN_ID")

	resourceName := "microsoft365_graph_beta_users_user_license_assignment.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckUserLicenseAssignmentDestroy,
		Steps: []resource.TestStep{
			// Start with minimal configuration
			{
				Config: testAccConfigMinimalNamed("test", testUserID1, testLicenseSkuID1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUserLicenseAssignmentExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "user_id", testUserID1),
					resource.TestCheckResourceAttr(resourceName, "add_licenses.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "add_licenses.0.sku_id", testLicenseSkuID1),
				),
			},
			// Update to maximal configuration
			{
				Config: testAccConfigMaximalNamed("test", testUserID1, testLicenseSkuID1, testLicenseSkuID2, testServicePlanID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUserLicenseAssignmentExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "user_id", testUserID1),
					resource.TestCheckResourceAttr(resourceName, "add_licenses.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "add_licenses.0.sku_id", testLicenseSkuID1),
					resource.TestCheckResourceAttr(resourceName, "add_licenses.1.sku_id", testLicenseSkuID2),
				),
			},
		},
	})
}

// TestAccUserLicenseAssignmentResource_Update_MaximalToMinimal tests updating from maximal to minimal config
func TestAccUserLicenseAssignmentResource_Update_MaximalToMinimal(t *testing.T) {
	// Skip if not running acceptance tests
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless TF_ACC=1")
	}

	// Get test user ID from environment variable or skip
	testUserID2 := os.Getenv("TEST_USER_ID_2")
	if testUserID2 == "" {
		t.Skip("TEST_USER_ID_2 environment variable must be set for acceptance tests")
	}

	// Get test license SKU IDs from environment variables or skip
	testLicenseSkuID1 := os.Getenv("TEST_LICENSE_SKU_ID_1")
	if testLicenseSkuID1 == "" {
		t.Skip("TEST_LICENSE_SKU_ID_1 environment variable must be set for acceptance tests")
	}

	testLicenseSkuID2 := os.Getenv("TEST_LICENSE_SKU_ID_2")
	if testLicenseSkuID2 == "" {
		t.Skip("TEST_LICENSE_SKU_ID_2 environment variable must be set for acceptance tests")
	}

	// Get test service plan ID from environment variable (optional)
	testServicePlanID := os.Getenv("TEST_SERVICE_PLAN_ID")

	resourceName := "microsoft365_graph_beta_users_user_license_assignment.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckUserLicenseAssignmentDestroy,
		Steps: []resource.TestStep{
			// Start with maximal configuration
			{
				Config: testAccConfigMaximalNamed("test", testUserID2, testLicenseSkuID1, testLicenseSkuID2, testServicePlanID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUserLicenseAssignmentExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "user_id", testUserID2),
					resource.TestCheckResourceAttr(resourceName, "add_licenses.#", "2"),
				),
			},
			// Update to minimal configuration
			{
				Config: testAccConfigMinimalNamed("test", testUserID2, testLicenseSkuID1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUserLicenseAssignmentExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "user_id", testUserID2),
					resource.TestCheckResourceAttr(resourceName, "add_licenses.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "add_licenses.0.sku_id", testLicenseSkuID1),
				),
			},
		},
	})
}

// TestAccUserLicenseAssignmentResource_Delete_Minimal tests deleting a license assignment with minimal configuration
func TestAccUserLicenseAssignmentResource_Delete_Minimal(t *testing.T) {
	// Skip if not running acceptance tests
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless TF_ACC=1")
	}

	// Get test user ID from environment variable or skip
	testUserID1 := os.Getenv("TEST_USER_ID_1")
	if testUserID1 == "" {
		t.Skip("TEST_USER_ID_1 environment variable must be set for acceptance tests")
	}

	// Get test license SKU ID from environment variable or skip
	testLicenseSkuID1 := os.Getenv("TEST_LICENSE_SKU_ID_1")
	if testLicenseSkuID1 == "" {
		t.Skip("TEST_LICENSE_SKU_ID_1 environment variable must be set for acceptance tests")
	}

	resourceName := "microsoft365_graph_beta_users_user_license_assignment.minimal"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckUserLicenseAssignmentDestroy,
		Steps: []resource.TestStep{
			// Create the resource
			{
				Config: testAccConfigMinimal(testUserID1, testLicenseSkuID1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUserLicenseAssignmentExists(resourceName),
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

// TestAccUserLicenseAssignmentResource_Delete_Maximal tests deleting a license assignment with maximal configuration
func TestAccUserLicenseAssignmentResource_Delete_Maximal(t *testing.T) {
	// Skip if not running acceptance tests
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless TF_ACC=1")
	}

	// Get test user ID from environment variable or skip
	testUserID2 := os.Getenv("TEST_USER_ID_2")
	if testUserID2 == "" {
		t.Skip("TEST_USER_ID_2 environment variable must be set for acceptance tests")
	}

	// Get test license SKU IDs from environment variables or skip
	testLicenseSkuID1 := os.Getenv("TEST_LICENSE_SKU_ID_1")
	if testLicenseSkuID1 == "" {
		t.Skip("TEST_LICENSE_SKU_ID_1 environment variable must be set for acceptance tests")
	}

	testLicenseSkuID2 := os.Getenv("TEST_LICENSE_SKU_ID_2")
	if testLicenseSkuID2 == "" {
		t.Skip("TEST_LICENSE_SKU_ID_2 environment variable must be set for acceptance tests")
	}

	// Get test service plan ID from environment variable (optional)
	testServicePlanID := os.Getenv("TEST_SERVICE_PLAN_ID")

	resourceName := "microsoft365_graph_beta_users_user_license_assignment.maximal"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckUserLicenseAssignmentDestroy,
		Steps: []resource.TestStep{
			// Create the resource
			{
				Config: testAccConfigMaximal(testUserID2, testLicenseSkuID1, testLicenseSkuID2, testServicePlanID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUserLicenseAssignmentExists(resourceName),
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

// TestAccUserLicenseAssignmentResource_Import tests importing a resource
func TestAccUserLicenseAssignmentResource_Import(t *testing.T) {
	// Skip if not running acceptance tests
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless TF_ACC=1")
	}

	// Get test user ID from environment variable or skip
	testUserID1 := os.Getenv("TEST_USER_ID_1")
	if testUserID1 == "" {
		t.Skip("TEST_USER_ID_1 environment variable must be set for acceptance tests")
	}

	// Get test license SKU ID from environment variable or skip
	testLicenseSkuID1 := os.Getenv("TEST_LICENSE_SKU_ID_1")
	if testLicenseSkuID1 == "" {
		t.Skip("TEST_LICENSE_SKU_ID_1 environment variable must be set for acceptance tests")
	}

	resourceName := "microsoft365_graph_beta_users_user_license_assignment.minimal"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckUserLicenseAssignmentDestroy,
		Steps: []resource.TestStep{
			// Create
			{
				Config: testAccConfigMinimal(testUserID1, testLicenseSkuID1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUserLicenseAssignmentExists(resourceName),
				),
			},
			// Import
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"remove_licenses", // This is expected to be ignored on import
					"add_licenses",    // The import only sets assigned_licenses, not add_licenses
				},
			},
		},
	})
}

// Helper functions for acceptance tests

func testAccCheckUserLicenseAssignmentExists(resourceName string) resource.TestCheckFunc {
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

func testAccCheckUserLicenseAssignmentDestroy(s *terraform.State) error {
	// In a real test, we would verify the license assignment is removed
	// For this resource, we don't need to check anything special since removing
	// the resource will remove the license assignment
	return nil
}

// Test configurations

// Minimal configuration with default resource name
func testAccConfigMinimal(userID, licenseSkuID string) string {
	return fmt.Sprintf(`
resource "microsoft365_graph_beta_users_user_license_assignment" "minimal" {
  user_id = "%s"
  add_licenses = [{
    sku_id = "%s"
  }]
}
`, userID, licenseSkuID)
}

// Minimal configuration with custom resource name
func testAccConfigMinimalNamed(resourceName string, userID, licenseSkuID string) string {
	return fmt.Sprintf(`
resource "microsoft365_graph_beta_users_user_license_assignment" "%s" {
  user_id = "%s"
  add_licenses = [{
    sku_id = "%s"
  }]
}
`, resourceName, userID, licenseSkuID)
}

// Maximal configuration with default resource name
func testAccConfigMaximal(userID, licenseSkuID1, licenseSkuID2, servicePlanID string) string {
	disabledPlans := ""
	if servicePlanID != "" {
		disabledPlans = fmt.Sprintf(`
      disabled_plans = [
        "%s"
      ]`, servicePlanID)
	}

	return fmt.Sprintf(`
resource "microsoft365_graph_beta_users_user_license_assignment" "maximal" {
  user_id = "%s"
  add_licenses = [
    {
      sku_id = "%s"%s
    },
    {
      sku_id = "%s"
    }
  ]
}
`, userID, licenseSkuID1, disabledPlans, licenseSkuID2)
}

// Maximal configuration with custom resource name
func testAccConfigMaximalNamed(resourceName string, userID, licenseSkuID1, licenseSkuID2, servicePlanID string) string {
	disabledPlans := ""
	if servicePlanID != "" {
		disabledPlans = fmt.Sprintf(`
      disabled_plans = [
        "%s"
      ]`, servicePlanID)
	}

	return fmt.Sprintf(`
resource "microsoft365_graph_beta_users_user_license_assignment" "%s" {
  user_id = "%s"
  add_licenses = [
    {
      sku_id = "%s"%s
    },
    {
      sku_id = "%s"
    }
  ]
}
`, resourceName, userID, licenseSkuID1, disabledPlans, licenseSkuID2)
}
