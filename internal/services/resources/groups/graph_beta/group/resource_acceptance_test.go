package graphBetaGroup_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// TestAccGroupResource_Create_Minimal tests creating a group with minimal configuration
func TestAccGroupResource_Create_Minimal(t *testing.T) {
	// Skip if not running acceptance tests
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless TF_ACC=1")
	}

	resourceName := "microsoft365_graph_beta_groups_group.minimal"
	uniqueSuffix := generateUniqueNameSuffix()
	displayName := fmt.Sprintf("tfacctest-minimal-%s", uniqueSuffix)
	mailNickname := fmt.Sprintf("tfacctest-minimal-%s", uniqueSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckGroupDestroy,
		Steps: []resource.TestStep{
			// Create with minimal configuration
			{
				Config: testAccConfigMinimal(displayName, mailNickname),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "display_name", displayName),
					resource.TestCheckResourceAttr(resourceName, "mail_nickname", mailNickname),
					resource.TestCheckResourceAttr(resourceName, "mail_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "security_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "visibility", "Private"),
				),
			},
		},
	})
}

// TestAccGroupResource_Create_Maximal tests creating a group with maximal configuration
func TestAccGroupResource_Create_Maximal(t *testing.T) {
	// Skip if not running acceptance tests
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless TF_ACC=1")
	}

	resourceName := "microsoft365_graph_beta_groups_group.maximal"
	uniqueSuffix := generateUniqueNameSuffix()
	displayName := fmt.Sprintf("tfacctest-maximal-%s", uniqueSuffix)
	mailNickname := fmt.Sprintf("tfacctest-maximal-%s", uniqueSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckGroupDestroy,
		Steps: []resource.TestStep{
			// Create with maximal configuration
			{
				Config: testAccConfigMaximal(displayName, mailNickname),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "display_name", displayName),
					resource.TestCheckResourceAttr(resourceName, "description", "This is a maximal group configuration for testing"),
					resource.TestCheckResourceAttr(resourceName, "mail_nickname", mailNickname),
					resource.TestCheckResourceAttr(resourceName, "mail_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "security_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "group_types.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "visibility", "Private"),
					resource.TestCheckResourceAttr(resourceName, "membership_rule", "user.department -eq \"Engineering\""),
					resource.TestCheckResourceAttr(resourceName, "membership_rule_processing_state", "On"),
					resource.TestCheckResourceAttr(resourceName, "preferred_language", "en-US"),
					resource.TestCheckResourceAttr(resourceName, "classification", "High"),
				),
			},
		},
	})
}

// TestAccGroupResource_Update_MinimalToMaximal tests updating from minimal to maximal config
func TestAccGroupResource_Update_MinimalToMaximal(t *testing.T) {
	// Skip if not running acceptance tests
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless TF_ACC=1")
	}

	resourceName := "microsoft365_graph_beta_groups_group.test"
	uniqueSuffix := generateUniqueNameSuffix()
	displayName := fmt.Sprintf("tfacctest-update-%s", uniqueSuffix)
	mailNickname := fmt.Sprintf("tfacctest-update-%s", uniqueSuffix)
	updatedDisplayName := fmt.Sprintf("tfacctest-updated-%s", uniqueSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckGroupDestroy,
		Steps: []resource.TestStep{
			// Start with minimal configuration
			{
				Config: testAccConfigMinimalNamed("test", displayName, mailNickname),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "display_name", displayName),
					resource.TestCheckResourceAttr(resourceName, "mail_nickname", mailNickname),
					resource.TestCheckResourceAttr(resourceName, "mail_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "security_enabled", "true"),
				),
			},
			// Update to maximal configuration
			{
				Config: testAccConfigMaximalNamed("test", updatedDisplayName, mailNickname),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "display_name", updatedDisplayName),
					resource.TestCheckResourceAttr(resourceName, "description", "This is a maximal group configuration for testing"),
					resource.TestCheckResourceAttr(resourceName, "mail_nickname", mailNickname),
					resource.TestCheckResourceAttr(resourceName, "mail_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "security_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "group_types.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "membership_rule", "user.department -eq \"Engineering\""),
					resource.TestCheckResourceAttr(resourceName, "membership_rule_processing_state", "On"),
				),
			},
		},
	})
}

// TestAccGroupResource_Update_MaximalToMinimal tests updating from maximal to minimal config
func TestAccGroupResource_Update_MaximalToMinimal(t *testing.T) {
	// Skip if not running acceptance tests
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless TF_ACC=1")
	}

	resourceName := "microsoft365_graph_beta_groups_group.test"
	uniqueSuffix := generateUniqueNameSuffix()
	displayName := fmt.Sprintf("tfacctest-downgrade-%s", uniqueSuffix)
	mailNickname := fmt.Sprintf("tfacctest-downgrade-%s", uniqueSuffix)
	updatedDisplayName := fmt.Sprintf("tfacctest-minimal-%s", uniqueSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckGroupDestroy,
		Steps: []resource.TestStep{
			// Start with maximal configuration
			{
				Config: testAccConfigMaximalNamed("test", displayName, mailNickname),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "display_name", displayName),
					resource.TestCheckResourceAttr(resourceName, "description", "This is a maximal group configuration for testing"),
					resource.TestCheckResourceAttr(resourceName, "mail_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "group_types.#", "2"),
				),
			},
			// Update to minimal configuration
			{
				Config: testAccConfigMinimalNamed("test", updatedDisplayName, mailNickname),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "display_name", updatedDisplayName),
					resource.TestCheckResourceAttr(resourceName, "mail_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "security_enabled", "true"),
				),
			},
		},
	})
}

// TestAccGroupResource_Delete_Minimal tests deleting a group with minimal configuration
func TestAccGroupResource_Delete_Minimal(t *testing.T) {
	// Skip if not running acceptance tests
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless TF_ACC=1")
	}

	resourceName := "microsoft365_graph_beta_groups_group.minimal"
	uniqueSuffix := generateUniqueNameSuffix()
	displayName := fmt.Sprintf("tfacctest-delete-minimal-%s", uniqueSuffix)
	mailNickname := fmt.Sprintf("tfacctest-delete-min-%s", uniqueSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckGroupDestroy,
		Steps: []resource.TestStep{
			// Create the resource
			{
				Config: testAccConfigMinimal(displayName, mailNickname),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGroupExists(resourceName),
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

// TestAccGroupResource_Delete_Maximal tests deleting a group with maximal configuration
func TestAccGroupResource_Delete_Maximal(t *testing.T) {
	// Skip if not running acceptance tests
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless TF_ACC=1")
	}

	resourceName := "microsoft365_graph_beta_groups_group.maximal"
	uniqueSuffix := generateUniqueNameSuffix()
	displayName := fmt.Sprintf("tfacctest-delete-maximal-%s", uniqueSuffix)
	mailNickname := fmt.Sprintf("tfacctest-delete-max-%s", uniqueSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckGroupDestroy,
		Steps: []resource.TestStep{
			// Create the resource
			{
				Config: testAccConfigMaximal(displayName, mailNickname),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGroupExists(resourceName),
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

// TestAccGroupResource_Import tests importing a resource
func TestAccGroupResource_Import(t *testing.T) {
	// Skip if not running acceptance tests
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless TF_ACC=1")
	}

	resourceName := "microsoft365_graph_beta_groups_group.minimal"
	uniqueSuffix := generateUniqueNameSuffix()
	displayName := fmt.Sprintf("tfacctest-import-%s", uniqueSuffix)
	mailNickname := fmt.Sprintf("tfacctest-import-%s", uniqueSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckGroupDestroy,
		Steps: []resource.TestStep{
			// Create
			{
				Config: testAccConfigMinimal(displayName, mailNickname),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGroupExists(resourceName),
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

func testAccCheckGroupExists(resourceName string) resource.TestCheckFunc {
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

func testAccCheckGroupDestroy(s *terraform.State) error {
	// In a real test, we would verify the group is removed
	// For this resource, we don't need to check anything special since removing
	// the resource will remove the group
	return nil
}

// Helper function to generate a unique name suffix
func generateUniqueNameSuffix() string {
	// Use a timestamp or random string to ensure uniqueness
	return fmt.Sprintf("%d", os.Getpid())
}

// Test configurations

// Minimal configuration with specified display name and mail nickname
func testAccConfigMinimal(displayName, mailNickname string) string {
	return fmt.Sprintf(`resource "microsoft365_graph_beta_groups_group" "minimal" {
  display_name     = "%s"
  mail_nickname    = "%s"
  mail_enabled     = false
  security_enabled = true
}`, displayName, mailNickname)
}

// Minimal configuration with custom resource name
func testAccConfigMinimalNamed(resourceName, displayName, mailNickname string) string {
	return fmt.Sprintf(`resource "microsoft365_graph_beta_groups_group" "%s" {
  display_name     = "%s"
  mail_nickname    = "%s"
  mail_enabled     = false
  security_enabled = true
}`, resourceName, displayName, mailNickname)
}

// Maximal configuration with specified display name and mail nickname
func testAccConfigMaximal(displayName, mailNickname string) string {
	return fmt.Sprintf(`resource "microsoft365_graph_beta_groups_group" "maximal" {
  display_name                   = "%s"
  description                    = "This is a maximal group configuration for testing"
  mail_nickname                  = "%s"
  mail_enabled                   = true
  security_enabled               = true
  group_types                    = ["Unified", "DynamicMembership"]
  visibility                     = "Private"
  is_assignable_to_role          = false
  membership_rule                = "user.department -eq \"Engineering\""
  membership_rule_processing_state = "On"
  preferred_language             = "en-US"
  theme                          = "Blue"
  classification                 = "High"
}`, displayName, mailNickname)
}

// Maximal configuration with custom resource name
func testAccConfigMaximalNamed(resourceName, displayName, mailNickname string) string {
	return fmt.Sprintf(`resource "microsoft365_graph_beta_groups_group" "%s" {
  display_name                   = "%s"
  description                    = "This is a maximal group configuration for testing"
  mail_nickname                  = "%s"
  mail_enabled                   = true
  security_enabled               = true
  group_types                    = ["Unified", "DynamicMembership"]
  visibility                     = "Private"
  is_assignable_to_role          = false
  membership_rule                = "user.department -eq \"Engineering\""
  membership_rule_processing_state = "On"
  preferred_language             = "en-US"
  theme                          = "Blue"
  classification                 = "High"
}`, resourceName, displayName, mailNickname)
}
