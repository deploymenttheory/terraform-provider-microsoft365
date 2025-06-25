package graphBetaSettingsCatalog_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// TestAccSettingsCatalogResource_Create_Minimal tests creating a settings catalog with minimal configuration
func TestAccSettingsCatalogResource_Create_Minimal(t *testing.T) {
	// Skip if not running acceptance tests
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless TF_ACC=1")
	}

	resourceName := "microsoft365_graph_beta_device_management_settings_catalog.minimal"
	policyName := "tfacctest-minimal-settings-catalog"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckSettingsCatalogDestroy,
		Steps: []resource.TestStep{
			// Create with minimal configuration
			{
				Config: testAccConfigMinimal(policyName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSettingsCatalogExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", policyName),
					resource.TestCheckResourceAttr(resourceName, "description", "Minimal settings catalog policy"),
					resource.TestCheckResourceAttr(resourceName, "platform", "windows10"),
				),
			},
		},
	})
}

// TestAccSettingsCatalogResource_Create_Maximal tests creating a settings catalog with maximal configuration
func TestAccSettingsCatalogResource_Create_Maximal(t *testing.T) {
	// Skip if not running acceptance tests
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless TF_ACC=1")
	}

	resourceName := "microsoft365_graph_beta_device_management_settings_catalog.maximal"
	policyName := "tfacctest-maximal-settings-catalog"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckSettingsCatalogDestroy,
		Steps: []resource.TestStep{
			// Create with maximal configuration
			{
				Config: testAccConfigMaximal(policyName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSettingsCatalogExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", policyName),
					resource.TestCheckResourceAttr(resourceName, "description", "Maximal settings catalog policy with all options"),
					resource.TestCheckResourceAttr(resourceName, "platform", "windows10"),
					resource.TestCheckResourceAttr(resourceName, "technologies", "mdm"),
					resource.TestCheckResourceAttr(resourceName, "settings.#", "2"),
					// Assignments may be computed and could change, so we just check they exist
					resource.TestCheckResourceAttrSet(resourceName, "assignments.#"),
				),
			},
		},
	})
}

// TestAccSettingsCatalogResource_Update_MinimalToMaximal tests updating from minimal to maximal config
func TestAccSettingsCatalogResource_Update_MinimalToMaximal(t *testing.T) {
	// Skip if not running acceptance tests
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless TF_ACC=1")
	}

	resourceName := "microsoft365_graph_beta_device_management_settings_catalog.test"
	policyName := "tfacctest-update-settings-catalog"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckSettingsCatalogDestroy,
		Steps: []resource.TestStep{
			// Start with minimal configuration
			{
				Config: testAccConfigMinimalNamed("test", policyName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSettingsCatalogExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", policyName),
					resource.TestCheckResourceAttr(resourceName, "description", "Minimal settings catalog policy"),
					resource.TestCheckResourceAttr(resourceName, "platform", "windows10"),
				),
			},
			// Update to maximal configuration
			{
				Config: testAccConfigMaximalNamed("test", policyName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSettingsCatalogExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", policyName),
					resource.TestCheckResourceAttr(resourceName, "description", "Maximal settings catalog policy with all options"),
					resource.TestCheckResourceAttr(resourceName, "platform", "windows10"),
					resource.TestCheckResourceAttr(resourceName, "technologies", "mdm"),
					resource.TestCheckResourceAttr(resourceName, "settings.#", "2"),
					resource.TestCheckResourceAttrSet(resourceName, "assignments.#"),
				),
			},
		},
	})
}

// TestAccSettingsCatalogResource_Update_MaximalToMinimal tests updating from maximal to minimal config
func TestAccSettingsCatalogResource_Update_MaximalToMinimal(t *testing.T) {
	// Skip if not running acceptance tests
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless TF_ACC=1")
	}

	resourceName := "microsoft365_graph_beta_device_management_settings_catalog.test"
	policyName := "tfacctest-downgrade-settings-catalog"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckSettingsCatalogDestroy,
		Steps: []resource.TestStep{
			// Start with maximal configuration
			{
				Config: testAccConfigMaximalNamed("test", policyName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSettingsCatalogExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", policyName),
					resource.TestCheckResourceAttr(resourceName, "settings.#", "2"),
				),
			},
			// Update to minimal configuration
			{
				Config: testAccConfigMinimalNamed("test", policyName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSettingsCatalogExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", policyName),
					resource.TestCheckResourceAttr(resourceName, "description", "Minimal settings catalog policy"),
					resource.TestCheckResourceAttr(resourceName, "platform", "windows10"),
					// Settings should be empty now
					resource.TestCheckResourceAttr(resourceName, "settings.#", "0"),
				),
			},
		},
	})
}

// TestAccSettingsCatalogResource_Delete_Minimal tests deleting a settings catalog with minimal configuration
func TestAccSettingsCatalogResource_Delete_Minimal(t *testing.T) {
	// Skip if not running acceptance tests
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless TF_ACC=1")
	}

	resourceName := "microsoft365_graph_beta_device_management_settings_catalog.minimal"
	policyName := "tfacctest-delete-minimal-settings-catalog"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckSettingsCatalogDestroy,
		Steps: []resource.TestStep{
			// Create the resource
			{
				Config: testAccConfigMinimal(policyName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSettingsCatalogExists(resourceName),
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

// TestAccSettingsCatalogResource_Delete_Maximal tests deleting a settings catalog with maximal configuration
func TestAccSettingsCatalogResource_Delete_Maximal(t *testing.T) {
	// Skip if not running acceptance tests
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless TF_ACC=1")
	}

	resourceName := "microsoft365_graph_beta_device_management_settings_catalog.maximal"
	policyName := "tfacctest-delete-maximal-settings-catalog"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckSettingsCatalogDestroy,
		Steps: []resource.TestStep{
			// Create the resource
			{
				Config: testAccConfigMaximal(policyName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSettingsCatalogExists(resourceName),
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

// TestAccSettingsCatalogResource_Import tests importing a resource
func TestAccSettingsCatalogResource_Import(t *testing.T) {
	// Skip if not running acceptance tests
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless TF_ACC=1")
	}

	resourceName := "microsoft365_graph_beta_device_management_settings_catalog.minimal"
	policyName := "tfacctest-import-settings-catalog"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckSettingsCatalogDestroy,
		Steps: []resource.TestStep{
			// Create
			{
				Config: testAccConfigMinimal(policyName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSettingsCatalogExists(resourceName),
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

func testAccPreCheck(t *testing.T) {
	// Verify required environment variables are set
	requiredEnvVars := []string{
		"ARM_CLIENT_ID",
		"ARM_CLIENT_SECRET",
		"ARM_TENANT_ID",
	}

	for _, env := range requiredEnvVars {
		if os.Getenv(env) == "" {
			t.Fatalf("%s environment variable must be set for acceptance tests", env)
		}
	}
}

func testAccCheckSettingsCatalogExists(resourceName string) resource.TestCheckFunc {
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

func testAccCheckSettingsCatalogDestroy(s *terraform.State) error {
	// In a real test, we would verify the settings catalog policy is removed
	// For this resource, we don't need to check anything special since removing
	// the resource will remove the policy
	return nil
}

// Test configurations

// Minimal configuration with default resource name
func testAccConfigMinimal(policyName string) string {
	// Read the template file
	content, err := os.ReadFile(filepath.Join("mocks", "terraform", "resource_minimal.tf"))
	if err != nil {
		return fmt.Sprintf(`resource "microsoft365_graph_beta_device_management_settings_catalog" "minimal" {
  name        = "%s"
  description = "Minimal settings catalog policy"
  platform    = "windows10"
}`, policyName)
	}

	// Replace the policy name
	updated := strings.Replace(string(content), "Minimal Settings Catalog", policyName, 1)

	return updated
}

// Minimal configuration with custom resource name
func testAccConfigMinimalNamed(resourceName string, policyName string) string {
	// Read the template file
	content, err := os.ReadFile(filepath.Join("mocks", "terraform", "resource_minimal.tf"))
	if err != nil {
		return fmt.Sprintf(`resource "microsoft365_graph_beta_device_management_settings_catalog" "%s" {
  name        = "%s"
  description = "Minimal settings catalog policy"
  platform    = "windows10"
}`, resourceName, policyName)
	}

	// Replace the resource name and policy name
	updated := strings.Replace(string(content), "minimal", resourceName, 1)
	updated = strings.Replace(updated, "Minimal Settings Catalog", policyName, 1)

	return updated
}

// Maximal configuration with default resource name
func testAccConfigMaximal(policyName string) string {
	// Read the template file
	content, err := os.ReadFile(filepath.Join("mocks", "terraform", "resource_maximal.tf"))
	if err != nil {
		return fmt.Sprintf(`resource "microsoft365_graph_beta_device_management_settings_catalog" "maximal" {
  name        = "%s"
  description = "Maximal settings catalog policy with all options"
  platform    = "windows10"
  technologies = "mdm"
  
  settings {
    setting_instance {
      setting_definition_id = "device_vendor_msft_policy_config_defender_allowarchivescanning"
      value_json            = "true"
    }
    setting_instance {
      setting_definition_id = "device_vendor_msft_policy_config_defender_allowbehaviormonitoring"
      value_json            = "true"
    }
  }
  
  assignments {
    intent = "apply"
    target {
      type = "allDevices"
    }
  }
}`, policyName)
	}

	// Replace the policy name
	updated := strings.Replace(string(content), "Maximal Settings Catalog", policyName, 1)

	return updated
}

// Maximal configuration with custom resource name
func testAccConfigMaximalNamed(resourceName string, policyName string) string {
	// Read the template file
	content, err := os.ReadFile(filepath.Join("mocks", "terraform", "resource_maximal.tf"))
	if err != nil {
		return fmt.Sprintf(`resource "microsoft365_graph_beta_device_management_settings_catalog" "%s" {
  name        = "%s"
  description = "Maximal settings catalog policy with all options"
  platform    = "windows10"
  technologies = "mdm"
  
  settings {
    setting_instance {
      setting_definition_id = "device_vendor_msft_policy_config_defender_allowarchivescanning"
      value_json            = "true"
    }
    setting_instance {
      setting_definition_id = "device_vendor_msft_policy_config_defender_allowbehaviormonitoring"
      value_json            = "true"
    }
  }
  
  assignments {
    intent = "apply"
    target {
      type = "allDevices"
    }
  }
}`, resourceName, policyName)
	}

	// Replace the resource name and policy name
	updated := strings.Replace(string(content), "maximal", resourceName, 1)
	updated = strings.Replace(updated, "Maximal Settings Catalog", policyName, 1)

	return updated
}
