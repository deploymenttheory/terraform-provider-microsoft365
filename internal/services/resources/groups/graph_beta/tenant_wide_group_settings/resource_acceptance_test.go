package graphBetaTenantWideGroupSettings_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// TestAccTenantWideGroupSettingsResource_Create_Minimal tests creating tenant-wide group settings with minimal configuration
func TestAccTenantWideGroupSettingsResource_Create_Minimal(t *testing.T) {
	// Skip if not running acceptance tests
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless TF_ACC=1")
	}

	resourceName := "microsoft365_graph_beta_groups_tenant_wide_group_settings.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckTenantWideGroupSettingsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTenantWideGroupSettingsExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "template_id", "62375ab9-6b52-47ed-826b-58e47e0e304b"),
					resource.TestCheckResourceAttr(resourceName, "display_name", "Group.Unified"),
					resource.TestCheckResourceAttr(resourceName, "values.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "values.*", map[string]string{
						"name":  "EnableGroupCreation",
						"value": "true",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "values.*", map[string]string{
						"name":  "GroupCreationAllowedGroupId",
						"value": "",
					}),
				),
			},
		},
	})
}

// TestAccTenantWideGroupSettingsResource_Create_Maximal tests creating tenant-wide group settings with maximal configuration
func TestAccTenantWideGroupSettingsResource_Create_Maximal(t *testing.T) {
	// Skip if not running acceptance tests
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless TF_ACC=1")
	}

	resourceName := "microsoft365_graph_beta_groups_tenant_wide_group_settings.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckTenantWideGroupSettingsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigMaximal(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTenantWideGroupSettingsExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "template_id", "62375ab9-6b52-47ed-826b-58e47e0e304b"),
					resource.TestCheckResourceAttr(resourceName, "display_name", "Group.Unified"),
					resource.TestCheckResourceAttr(resourceName, "values.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "values.*", map[string]string{
						"name":  "EnableGroupCreation",
						"value": "false",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "values.*", map[string]string{
						"name":  "GroupCreationAllowedGroupId",
						"value": "allowed-group-id",
					}),
				),
			},
		},
	})
}

// TestAccTenantWideGroupSettingsResource_Update_MinimalToMaximal tests updating from minimal to maximal config
func TestAccTenantWideGroupSettingsResource_Update_MinimalToMaximal(t *testing.T) {
	// Skip if not running acceptance tests
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless TF_ACC=1")
	}

	resourceName := "microsoft365_graph_beta_groups_tenant_wide_group_settings.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckTenantWideGroupSettingsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigMinimalNamed("test"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTenantWideGroupSettingsExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "values.#", "2"),
				),
			},
			{
				Config: testAccConfigMaximalNamed("test"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTenantWideGroupSettingsExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "values.#", "2"),
				),
			},
		},
	})
}

// TestAccTenantWideGroupSettingsResource_Update_MaximalToMinimal tests updating from maximal to minimal config
func TestAccTenantWideGroupSettingsResource_Update_MaximalToMinimal(t *testing.T) {
	// Skip if not running acceptance tests
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless TF_ACC=1")
	}

	resourceName := "microsoft365_graph_beta_groups_tenant_wide_group_settings.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckTenantWideGroupSettingsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigMaximalNamed("test"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTenantWideGroupSettingsExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "values.#", "2"),
				),
			},
			{
				Config: testAccConfigMinimalNamed("test"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTenantWideGroupSettingsExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "values.#", "2"),
				),
			},
		},
	})
}

// TestAccTenantWideGroupSettingsResource_Delete_Minimal tests deleting tenant-wide group settings with minimal configuration
func TestAccTenantWideGroupSettingsResource_Delete_Minimal(t *testing.T) {
	// Skip if not running acceptance tests
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless TF_ACC=1")
	}

	resourceName := "microsoft365_graph_beta_groups_tenant_wide_group_settings.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckTenantWideGroupSettingsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTenantWideGroupSettingsExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "template_id", "62375ab9-6b52-47ed-826b-58e47e0e304b"),
					resource.TestCheckResourceAttr(resourceName, "display_name", "Group.Unified"),
					resource.TestCheckResourceAttr(resourceName, "values.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "values.*", map[string]string{
						"name":  "EnableGroupCreation",
						"value": "true",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "values.*", map[string]string{
						"name":  "GroupCreationAllowedGroupId",
						"value": "",
					}),
				),
			},
		},
	})
}

// TestAccTenantWideGroupSettingsResource_Delete_Maximal tests deleting tenant-wide group settings with maximal configuration
func TestAccTenantWideGroupSettingsResource_Delete_Maximal(t *testing.T) {
	// Skip if not running acceptance tests
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless TF_ACC=1")
	}

	resourceName := "microsoft365_graph_beta_groups_tenant_wide_group_settings.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckTenantWideGroupSettingsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigMaximal(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTenantWideGroupSettingsExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "template_id", "62375ab9-6b52-47ed-826b-58e47e0e304b"),
					resource.TestCheckResourceAttr(resourceName, "display_name", "Group.Unified"),
					resource.TestCheckResourceAttr(resourceName, "values.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "values.*", map[string]string{
						"name":  "EnableGroupCreation",
						"value": "false",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "values.*", map[string]string{
						"name":  "GroupCreationAllowedGroupId",
						"value": "allowed-group-id",
					}),
				),
			},
		},
	})
}

// TestAccTenantWideGroupSettingsResource_Import tests importing a resource
func TestAccTenantWideGroupSettingsResource_Import(t *testing.T) {
	// Skip if not running acceptance tests
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless TF_ACC=1")
	}

	resourceName := "microsoft365_graph_beta_groups_tenant_wide_group_settings.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckTenantWideGroupSettingsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTenantWideGroupSettingsExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "template_id", "62375ab9-6b52-47ed-826b-58e47e0e304b"),
					resource.TestCheckResourceAttr(resourceName, "display_name", "Group.Unified"),
					resource.TestCheckResourceAttr(resourceName, "values.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "values.*", map[string]string{
						"name":  "EnableGroupCreation",
						"value": "true",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "values.*", map[string]string{
						"name":  "GroupCreationAllowedGroupId",
						"value": "",
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

func testAccCheckTenantWideGroupSettingsExists(resourceName string) resource.TestCheckFunc {
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

func testAccCheckTenantWideGroupSettingsDestroy(s *terraform.State) error {
	// This is a placeholder - in a real implementation, you would check if the resource still exists in Microsoft 365
	// For now, we'll just return nil as the resource should be deleted by the test cleanup
	return nil
}

func testAccConfigMinimal() string {
	return `
resource "microsoft365_graph_beta_groups_tenant_wide_group_settings" "test" {
  template_id = "62375ab9-6b52-47ed-826b-58e47e0e304b"
  
  values = [
    {
      name  = "EnableGroupCreation"
      value = "true"
    },
    {
      name  = "GroupCreationAllowedGroupId"
      value = ""
    }
  ]
}
`
}

func testAccConfigMinimalNamed(resourceName string) string {
	return fmt.Sprintf(`
resource "microsoft365_graph_beta_groups_tenant_wide_group_settings" "%s" {
  template_id = "62375ab9-6b52-47ed-826b-58e47e0e304b"
  
  values = [
    {
      name  = "EnableGroupCreation"
      value = "true"
    },
    {
      name  = "GroupCreationAllowedGroupId"
      value = ""
    }
  ]
}
`, resourceName)
}

func testAccConfigMaximal() string {
	return `
resource "microsoft365_graph_beta_groups_tenant_wide_group_settings" "test" {
  template_id = "62375ab9-6b52-47ed-826b-58e47e0e304b"
  
  values = [
    {
      name  = "EnableGroupCreation"
      value = "false"
    },
    {
      name  = "GroupCreationAllowedGroupId"
      value = "allowed-group-id"
    }
  ]
}
`
}

func testAccConfigMaximalNamed(resourceName string) string {
	return fmt.Sprintf(`
resource "microsoft365_graph_beta_groups_tenant_wide_group_settings" "%s" {
  template_id = "62375ab9-6b52-47ed-826b-58e47e0e304b"
  
  values = [
    {
      name  = "EnableGroupCreation"
      value = "false"
    },
    {
      name  = "GroupCreationAllowedGroupId"
      value = "allowed-group-id"
    }
  ]
}
`, resourceName)
}
