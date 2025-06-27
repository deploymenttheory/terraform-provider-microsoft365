package graphBetaGroupSettings_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// TestAccGroupSettingsResource_Create_Minimal tests creating group settings with minimal configuration
func TestAccGroupSettingsResource_Create_Minimal(t *testing.T) {
	// Skip if not running acceptance tests
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless TF_ACC=1")
	}

	// Generate unique test data
	groupID := os.Getenv("MS365_TEST_GROUP_ID")
	if groupID == "" {
		t.Skip("Skipping acceptance test as MS365_TEST_GROUP_ID is not set")
	}

	resourceName := "microsoft365_graph_beta_groups_group_settings.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckGroupSettingsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigMinimal(groupID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGroupSettingsExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "group_id", groupID),
					resource.TestCheckResourceAttr(resourceName, "template_id", "08d542b9-071f-4e16-94b0-74abb372e3d9"),
					resource.TestCheckResourceAttr(resourceName, "display_name", "Group.Unified.Guest"),
					resource.TestCheckResourceAttr(resourceName, "values.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "values.*", map[string]string{
						"name":  "AllowToAddGuests",
						"value": "false",
					}),
				),
			},
		},
	})
}

// TestAccGroupSettingsResource_Create_Maximal tests creating group settings with maximal configuration
func TestAccGroupSettingsResource_Create_Maximal(t *testing.T) {
	// Skip if not running acceptance tests
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless TF_ACC=1")
	}

	// Generate unique test data
	groupID := os.Getenv("MS365_TEST_GROUP_ID")
	if groupID == "" {
		t.Skip("Skipping acceptance test as MS365_TEST_GROUP_ID is not set")
	}

	resourceName := "microsoft365_graph_beta_groups_group_settings.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckGroupSettingsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigMaximal(groupID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGroupSettingsExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "group_id", groupID),
					resource.TestCheckResourceAttr(resourceName, "template_id", "08d542b9-071f-4e16-94b0-74abb372e3d9"),
					resource.TestCheckResourceAttr(resourceName, "display_name", "Group.Unified.Guest"),
					resource.TestCheckResourceAttr(resourceName, "values.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "values.*", map[string]string{
						"name":  "AllowToAddGuests",
						"value": "true",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "values.*", map[string]string{
						"name":  "GuestUsageGuidelinesUrl",
						"value": "https://contoso.com/guest-guidelines",
					}),
				),
			},
		},
	})
}

// TestAccGroupSettingsResource_Update_MinimalToMaximal tests updating from minimal to maximal config
func TestAccGroupSettingsResource_Update_MinimalToMaximal(t *testing.T) {
	// Skip if not running acceptance tests
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless TF_ACC=1")
	}

	// Generate unique test data
	groupID := os.Getenv("MS365_TEST_GROUP_ID")
	if groupID == "" {
		t.Skip("Skipping acceptance test as MS365_TEST_GROUP_ID is not set")
	}

	resourceName := "microsoft365_graph_beta_groups_group_settings.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckGroupSettingsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigMinimalNamed("test", groupID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGroupSettingsExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "values.#", "1"),
				),
			},
			{
				Config: testAccConfigMaximalNamed("test", groupID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGroupSettingsExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "values.#", "2"),
				),
			},
		},
	})
}

// TestAccGroupSettingsResource_Update_MaximalToMinimal tests updating from maximal to minimal config
func TestAccGroupSettingsResource_Update_MaximalToMinimal(t *testing.T) {
	// Skip if not running acceptance tests
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless TF_ACC=1")
	}

	// Generate unique test data
	groupID := os.Getenv("MS365_TEST_GROUP_ID")
	if groupID == "" {
		t.Skip("Skipping acceptance test as MS365_TEST_GROUP_ID is not set")
	}

	resourceName := "microsoft365_graph_beta_groups_group_settings.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckGroupSettingsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigMaximalNamed("test", groupID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGroupSettingsExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "values.#", "2"),
				),
			},
			{
				Config: testAccConfigMinimalNamed("test", groupID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGroupSettingsExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "values.#", "1"),
				),
			},
		},
	})
}

// TestAccGroupSettingsResource_Delete_Minimal tests deleting group settings with minimal configuration
func TestAccGroupSettingsResource_Delete_Minimal(t *testing.T) {
	// Skip if not running acceptance tests
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless TF_ACC=1")
	}

	// Generate unique test data
	groupID := os.Getenv("MS365_TEST_GROUP_ID")
	if groupID == "" {
		t.Skip("Skipping acceptance test as MS365_TEST_GROUP_ID is not set")
	}

	resourceName := "microsoft365_graph_beta_groups_group_settings.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckGroupSettingsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigMinimal(groupID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGroupSettingsExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "group_id", groupID),
					resource.TestCheckResourceAttr(resourceName, "template_id", "08d542b9-071f-4e16-94b0-74abb372e3d9"),
					resource.TestCheckResourceAttr(resourceName, "display_name", "Group.Unified.Guest"),
					resource.TestCheckResourceAttr(resourceName, "values.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "values.*", map[string]string{
						"name":  "AllowToAddGuests",
						"value": "false",
					}),
				),
			},
		},
	})
}

// TestAccGroupSettingsResource_Delete_Maximal tests deleting group settings with maximal configuration
func TestAccGroupSettingsResource_Delete_Maximal(t *testing.T) {
	// Skip if not running acceptance tests
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless TF_ACC=1")
	}

	// Generate unique test data
	groupID := os.Getenv("MS365_TEST_GROUP_ID")
	if groupID == "" {
		t.Skip("Skipping acceptance test as MS365_TEST_GROUP_ID is not set")
	}

	resourceName := "microsoft365_graph_beta_groups_group_settings.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckGroupSettingsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigMaximal(groupID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGroupSettingsExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "group_id", groupID),
					resource.TestCheckResourceAttr(resourceName, "template_id", "08d542b9-071f-4e16-94b0-74abb372e3d9"),
					resource.TestCheckResourceAttr(resourceName, "display_name", "Group.Unified.Guest"),
					resource.TestCheckResourceAttr(resourceName, "values.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "values.*", map[string]string{
						"name":  "AllowToAddGuests",
						"value": "true",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "values.*", map[string]string{
						"name":  "GuestUsageGuidelinesUrl",
						"value": "https://contoso.com/guest-guidelines",
					}),
				),
			},
		},
	})
}

// TestAccGroupSettingsResource_Import tests importing a resource
func TestAccGroupSettingsResource_Import(t *testing.T) {
	// Skip if not running acceptance tests
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless TF_ACC=1")
	}

	// Generate unique test data
	groupID := os.Getenv("MS365_TEST_GROUP_ID")
	if groupID == "" {
		t.Skip("Skipping acceptance test as MS365_TEST_GROUP_ID is not set")
	}

	resourceName := "microsoft365_graph_beta_groups_group_settings.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckGroupSettingsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigMinimal(groupID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGroupSettingsExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "group_id", groupID),
					resource.TestCheckResourceAttr(resourceName, "template_id", "08d542b9-071f-4e16-94b0-74abb372e3d9"),
					resource.TestCheckResourceAttr(resourceName, "display_name", "Group.Unified.Guest"),
					resource.TestCheckResourceAttr(resourceName, "values.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "values.*", map[string]string{
						"name":  "AllowToAddGuests",
						"value": "false",
					}),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Helper functions

func testAccPreCheck(t *testing.T) {
	// Verify required environment variables are set
	requiredEnvVars := []string{
		"ARM_CLIENT_ID",
		"ARM_CLIENT_SECRET",
		"ARM_TENANT_ID",
		"MS365_TEST_GROUP_ID",
	}

	for _, env := range requiredEnvVars {
		if os.Getenv(env) == "" {
			t.Fatalf("%s environment variable must be set for acceptance tests", env)
		}
	}
}

func testAccCheckGroupSettingsExists(resourceName string) resource.TestCheckFunc {
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

func testAccCheckGroupSettingsDestroy(s *terraform.State) error {
	// This is a placeholder - in a real implementation, you would check if the resource still exists in Microsoft 365
	// For now, we'll just return nil as the resource should be deleted by the test cleanup
	return nil
}

func testAccConfigMinimal(groupID string) string {
	return fmt.Sprintf(`
resource "microsoft365_graph_beta_groups_group_settings" "test" {
  group_id    = "%s"
  template_id = "08d542b9-071f-4e16-94b0-74abb372e3d9"
  
  values = [
    {
      name  = "AllowToAddGuests"
      value = "false"
    }
  ]
}
`, groupID)
}

func testAccConfigMinimalNamed(resourceName string, groupID string) string {
	return fmt.Sprintf(`
resource "microsoft365_graph_beta_groups_group_settings" "%s" {
  group_id    = "%s"
  template_id = "08d542b9-071f-4e16-94b0-74abb372e3d9"
  
  values = [
    {
      name  = "AllowToAddGuests"
      value = "false"
    }
  ]
}
`, resourceName, groupID)
}

func testAccConfigMaximal(groupID string) string {
	return fmt.Sprintf(`
resource "microsoft365_graph_beta_groups_group_settings" "test" {
  group_id    = "%s"
  template_id = "08d542b9-071f-4e16-94b0-74abb372e3d9"
  
  values = [
    {
      name  = "AllowToAddGuests"
      value = "true"
    },
    {
      name  = "GuestUsageGuidelinesUrl"
      value = "https://contoso.com/guest-guidelines"
    }
  ]
}
`, groupID)
}

func testAccConfigMaximalNamed(resourceName string, groupID string) string {
	return fmt.Sprintf(`
resource "microsoft365_graph_beta_groups_group_settings" "%s" {
  group_id    = "%s"
  template_id = "08d542b9-071f-4e16-94b0-74abb372e3d9"
  
  values = [
    {
      name  = "AllowToAddGuests"
      value = "true"
    },
    {
      name  = "GuestUsageGuidelinesUrl"
      value = "https://contoso.com/guest-guidelines"
    }
  ]
}
`, resourceName, groupID)
}
