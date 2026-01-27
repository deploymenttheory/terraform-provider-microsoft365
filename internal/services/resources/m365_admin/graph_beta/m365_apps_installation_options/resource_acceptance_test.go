package graphM365AppsInstallationOptions_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// TestAccM365AppsInstallationOptionsResource_Create_Minimal tests creating M365 Apps Installation Options with minimal configuration
func TestAccResourceM365AppsInstallationOptions_01_Create_Minimal(t *testing.T) {
	// Skip if not running acceptance tests
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless TF_ACC=1")
	}

	resourceName := "microsoft365_graph_m365_admin_m365_apps_installation_options.minimal"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create with minimal configuration
			{
				Config: testAccConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckM365AppsInstallationOptionsExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "update_channel", "current"),
					resource.TestCheckResourceAttr(resourceName, "apps_for_windows.is_microsoft_365_apps_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "apps_for_windows.is_skype_for_business_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "apps_for_mac.is_microsoft_365_apps_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "apps_for_mac.is_skype_for_business_enabled", "true"),
				),
			},
		},
	})
}

// TestAccM365AppsInstallationOptionsResource_Create_Maximal tests creating M365 Apps Installation Options with maximal configuration
func TestAccResourceM365AppsInstallationOptions_02_Create_Maximal(t *testing.T) {
	// Skip if not running acceptance tests
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless TF_ACC=1")
	}

	resourceName := "microsoft365_graph_m365_admin_m365_apps_installation_options.maximal"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create with maximal configuration
			{
				Config: testAccConfigMaximal(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckM365AppsInstallationOptionsExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "update_channel", "semiAnnual"),
					resource.TestCheckResourceAttr(resourceName, "apps_for_windows.is_microsoft_365_apps_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "apps_for_windows.is_skype_for_business_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "apps_for_mac.is_microsoft_365_apps_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "apps_for_mac.is_skype_for_business_enabled", "false"),
				),
			},
		},
	})
}

// TestAccM365AppsInstallationOptionsResource_Update_MinimalToMaximal tests updating from minimal to maximal configuration
func TestAccResourceM365AppsInstallationOptions_03_Update_MinimalToMaximal(t *testing.T) {
	// Skip if not running acceptance tests
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless TF_ACC=1")
	}

	resourceName := "microsoft365_graph_m365_admin_m365_apps_installation_options.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Start with minimal configuration
			{
				Config: testAccConfigMinimalNamed("test"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckM365AppsInstallationOptionsExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "update_channel", "current"),
					resource.TestCheckResourceAttr(resourceName, "apps_for_windows.is_microsoft_365_apps_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "apps_for_mac.is_microsoft_365_apps_enabled", "true"),
				),
			},
			// Update to maximal configuration
			{
				Config: testAccConfigMaximalNamed("test"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckM365AppsInstallationOptionsExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "update_channel", "semiAnnual"),
					resource.TestCheckResourceAttr(resourceName, "apps_for_windows.is_microsoft_365_apps_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "apps_for_mac.is_microsoft_365_apps_enabled", "false"),
				),
			},
		},
	})
}

// TestAccM365AppsInstallationOptionsResource_Update_MaximalToMinimal tests updating from maximal to minimal configuration
func TestAccResourceM365AppsInstallationOptions_04_Update_MaximalToMinimal(t *testing.T) {
	// Skip if not running acceptance tests
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless TF_ACC=1")
	}

	resourceName := "microsoft365_graph_m365_admin_m365_apps_installation_options.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Start with maximal configuration
			{
				Config: testAccConfigMaximalNamed("test"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckM365AppsInstallationOptionsExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "update_channel", "semiAnnual"),
					resource.TestCheckResourceAttr(resourceName, "apps_for_windows.is_microsoft_365_apps_enabled", "false"),
				),
			},
			// Update to minimal configuration
			{
				Config: testAccConfigMinimalNamed("test"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckM365AppsInstallationOptionsExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "update_channel", "current"),
					resource.TestCheckResourceAttr(resourceName, "apps_for_windows.is_microsoft_365_apps_enabled", "true"),
				),
			},
		},
	})
}

// TestAccM365AppsInstallationOptionsResource_Delete_Minimal tests deleting M365 Apps Installation Options with minimal configuration
func TestAccResourceM365AppsInstallationOptions_05_Delete_Minimal(t *testing.T) {
	// Skip if not running acceptance tests
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless TF_ACC=1")
	}

	resourceName := "microsoft365_graph_m365_admin_m365_apps_installation_options.minimal"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create the resource
			{
				Config: testAccConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckM365AppsInstallationOptionsExists(resourceName),
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

// TestAccM365AppsInstallationOptionsResource_Delete_Maximal tests deleting M365 Apps Installation Options with maximal configuration
func TestAccResourceM365AppsInstallationOptions_06_Delete_Maximal(t *testing.T) {
	// Skip if not running acceptance tests
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless TF_ACC=1")
	}

	resourceName := "microsoft365_graph_m365_admin_m365_apps_installation_options.maximal"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create the resource
			{
				Config: testAccConfigMaximal(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckM365AppsInstallationOptionsExists(resourceName),
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

// TestAccM365AppsInstallationOptionsResource_Import tests importing a resource
func TestAccResourceM365AppsInstallationOptions_07_Import(t *testing.T) {
	// Skip if not running acceptance tests
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless TF_ACC=1")
	}

	resourceName := "microsoft365_graph_m365_admin_m365_apps_installation_options.minimal"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create
			{
				Config: testAccConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckM365AppsInstallationOptionsExists(resourceName),
				),
			},
			// Import
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       false,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

// Helper functions for acceptance tests

func testAccCheckM365AppsInstallationOptionsExists(resourceName string) resource.TestCheckFunc {
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

// Test configurations

// Minimal configuration with default resource name
func testAccConfigMinimal() string {
	return `
resource "microsoft365_graph_m365_admin_m365_apps_installation_options" "minimal" {
  update_channel = "current"
  
  apps_for_windows = {
    is_microsoft_365_apps_enabled = true
    is_skype_for_business_enabled = true
  }
  
  apps_for_mac = {
    is_microsoft_365_apps_enabled = true
    is_skype_for_business_enabled = true
  }
}
`
}

// Minimal configuration with custom resource name
func testAccConfigMinimalNamed(resourceName string) string {
	return fmt.Sprintf(`
resource "microsoft365_graph_m365_admin_m365_apps_installation_options" "%s" {
  update_channel = "current"
  
  apps_for_windows = {
    is_microsoft_365_apps_enabled = true
    is_skype_for_business_enabled = true
  }
  
  apps_for_mac = {
    is_microsoft_365_apps_enabled = true
    is_skype_for_business_enabled = true
  }
}
`, resourceName)
}

// Maximal configuration with default resource name
func testAccConfigMaximal() string {
	return `
resource "microsoft365_graph_m365_admin_m365_apps_installation_options" "maximal" {
  update_channel = "semiAnnual"
  
  apps_for_windows = {
    is_microsoft_365_apps_enabled = false
    is_skype_for_business_enabled = false
  }
  
  apps_for_mac = {
    is_microsoft_365_apps_enabled = false
    is_skype_for_business_enabled = false
  }
  
  timeouts = {
    create = "30m"
    read   = "10m"
    update = "30m"
    delete = "30m"
  }
}
`
}

// Maximal configuration with custom resource name
func testAccConfigMaximalNamed(resourceName string) string {
	return fmt.Sprintf(`
resource "microsoft365_graph_m365_admin_m365_apps_installation_options" "%s" {
  update_channel = "semiAnnual"
  
  apps_for_windows = {
    is_microsoft_365_apps_enabled = false
    is_skype_for_business_enabled = false
  }
  
  apps_for_mac = {
    is_microsoft_365_apps_enabled = false
    is_skype_for_business_enabled = false
  }
  
  timeouts = {
    create = "30m"
    read   = "10m"
    update = "30m"
    delete = "30m"
  }
}
`, resourceName)
}
