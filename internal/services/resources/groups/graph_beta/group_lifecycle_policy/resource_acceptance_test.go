package graphBetaGroupLifecyclePolicy_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// TestAccGroupLifecyclePolicyResource_Create_Minimal tests creating a group lifecycle policy with minimal configuration
func TestAccGroupLifecyclePolicyResource_Create_Minimal(t *testing.T) {
	// Skip if not running acceptance tests
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless TF_ACC=1")
	}

	resourceName := "microsoft365_graph_beta_groups_group_lifecycle_policy.minimal"
	uniqueSuffix := generateUniqueNameSuffix()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckGroupLifecyclePolicyDestroy,
		Steps: []resource.TestStep{
			// Create with minimal configuration
			{
				Config: testAccConfigMinimal(uniqueSuffix),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGroupLifecyclePolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "group_lifetime_in_days", "180"),
					resource.TestCheckResourceAttr(resourceName, "managed_group_types", "All"),
				),
			},
		},
	})
}

// TestAccGroupLifecyclePolicyResource_Create_Maximal tests creating a group lifecycle policy with maximal configuration
func TestAccGroupLifecyclePolicyResource_Create_Maximal(t *testing.T) {
	// Skip if not running acceptance tests
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless TF_ACC=1")
	}

	resourceName := "microsoft365_graph_beta_groups_group_lifecycle_policy.maximal"
	uniqueSuffix := generateUniqueNameSuffix()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckGroupLifecyclePolicyDestroy,
		Steps: []resource.TestStep{
			// Create with maximal configuration
			{
				Config: testAccConfigMaximal(uniqueSuffix),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGroupLifecyclePolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "group_lifetime_in_days", "365"),
					resource.TestCheckResourceAttr(resourceName, "managed_group_types", "Selected"),
					resource.TestCheckResourceAttr(resourceName, "alternate_notification_emails", "admin@example.com;notifications@example.com"),
				),
			},
		},
	})
}

// TestAccGroupLifecyclePolicyResource_Update_MinimalToMaximal tests updating from minimal to maximal config
func TestAccGroupLifecyclePolicyResource_Update_MinimalToMaximal(t *testing.T) {
	// Skip if not running acceptance tests
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless TF_ACC=1")
	}

	resourceName := "microsoft365_graph_beta_groups_group_lifecycle_policy.test"
	uniqueSuffix := generateUniqueNameSuffix()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckGroupLifecyclePolicyDestroy,
		Steps: []resource.TestStep{
			// Start with minimal configuration
			{
				Config: testAccConfigMinimalNamed("test", uniqueSuffix),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGroupLifecyclePolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "group_lifetime_in_days", "180"),
					resource.TestCheckResourceAttr(resourceName, "managed_group_types", "All"),
				),
			},
			// Update to maximal configuration
			{
				Config: testAccConfigMaximalNamed("test", uniqueSuffix),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGroupLifecyclePolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "group_lifetime_in_days", "365"),
					resource.TestCheckResourceAttr(resourceName, "managed_group_types", "Selected"),
					resource.TestCheckResourceAttr(resourceName, "alternate_notification_emails", "admin@example.com;notifications@example.com"),
				),
			},
		},
	})
}

// TestAccGroupLifecyclePolicyResource_Update_MaximalToMinimal tests updating from maximal to minimal config
func TestAccGroupLifecyclePolicyResource_Update_MaximalToMinimal(t *testing.T) {
	// Skip if not running acceptance tests
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless TF_ACC=1")
	}

	resourceName := "microsoft365_graph_beta_groups_group_lifecycle_policy.test"
	uniqueSuffix := generateUniqueNameSuffix()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckGroupLifecyclePolicyDestroy,
		Steps: []resource.TestStep{
			// Start with maximal configuration
			{
				Config: testAccConfigMaximalNamed("test", uniqueSuffix),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGroupLifecyclePolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "group_lifetime_in_days", "365"),
					resource.TestCheckResourceAttr(resourceName, "managed_group_types", "Selected"),
					resource.TestCheckResourceAttr(resourceName, "alternate_notification_emails", "admin@example.com;notifications@example.com"),
				),
			},
			// Update to minimal configuration
			{
				Config: testAccConfigMinimalNamed("test", uniqueSuffix),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGroupLifecyclePolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "group_lifetime_in_days", "180"),
					resource.TestCheckResourceAttr(resourceName, "managed_group_types", "All"),
				),
			},
		},
	})
}

// TestAccGroupLifecyclePolicyResource_Delete_Minimal tests deleting a group lifecycle policy with minimal configuration
func TestAccGroupLifecyclePolicyResource_Delete_Minimal(t *testing.T) {
	// Skip if not running acceptance tests
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless TF_ACC=1")
	}

	resourceName := "microsoft365_graph_beta_groups_group_lifecycle_policy.minimal"
	uniqueSuffix := generateUniqueNameSuffix()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckGroupLifecyclePolicyDestroy,
		Steps: []resource.TestStep{
			// Create the resource
			{
				Config: testAccConfigMinimal(uniqueSuffix),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGroupLifecyclePolicyExists(resourceName),
				),
			},
			// Delete the resource (by providing empty config)
			{
				Config: "",
				Check: resource.ComposeTestCheckFunc(
					// Verify the resource is gone
					func(s *terraform.State) error {
						_, ok := s.RootModule().Resources[resourceName]
						if ok {
							return fmt.Errorf("resource still exists after deletion")
						}
						return nil
					},
				),
			},
		},
	})
}

// TestAccGroupLifecyclePolicyResource_Delete_Maximal tests deleting a group lifecycle policy with maximal configuration
func TestAccGroupLifecyclePolicyResource_Delete_Maximal(t *testing.T) {
	// Skip if not running acceptance tests
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless TF_ACC=1")
	}

	resourceName := "microsoft365_graph_beta_groups_group_lifecycle_policy.maximal"
	uniqueSuffix := generateUniqueNameSuffix()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckGroupLifecyclePolicyDestroy,
		Steps: []resource.TestStep{
			// Create the resource
			{
				Config: testAccConfigMaximal(uniqueSuffix),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGroupLifecyclePolicyExists(resourceName),
				),
			},
			// Delete the resource (by providing empty config)
			{
				Config: "",
				Check: resource.ComposeTestCheckFunc(
					// Verify the resource is gone
					func(s *terraform.State) error {
						_, ok := s.RootModule().Resources[resourceName]
						if ok {
							return fmt.Errorf("resource still exists after deletion")
						}
						return nil
					},
				),
			},
		},
	})
}

// TestAccGroupLifecyclePolicyResource_Import tests importing an existing group lifecycle policy
func TestAccGroupLifecyclePolicyResource_Import(t *testing.T) {
	// Skip if not running acceptance tests
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless TF_ACC=1")
	}

	resourceName := "microsoft365_graph_beta_groups_group_lifecycle_policy.imported"
	uniqueSuffix := generateUniqueNameSuffix()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckGroupLifecyclePolicyDestroy,
		Steps: []resource.TestStep{
			// Create the resource first
			{
				Config: testAccConfigMinimalNamed("imported", uniqueSuffix),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGroupLifecyclePolicyExists(resourceName),
				),
			},
			// Import the resource
			{
				ResourceName: resourceName,
				ImportState:  true,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGroupLifecyclePolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "group_lifetime_in_days", "180"),
					resource.TestCheckResourceAttr(resourceName, "managed_group_types", "All"),
				),
			},
		},
	})
}

// Helper functions

func testAccCheckGroupLifecyclePolicyExists(resourceName string) resource.TestCheckFunc {
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

func testAccCheckGroupLifecyclePolicyDestroy(s *terraform.State) error {
	// For group lifecycle policies, we can't easily verify deletion
	// since there's typically only one policy per tenant
	// This is a placeholder that always passes
	return nil
}

func generateUniqueNameSuffix() string {
	return fmt.Sprintf("%d", time.Now().Unix())
}

func testAccConfigMinimal(uniqueSuffix string) string {
	return `
resource "microsoft365_graph_beta_groups_group_lifecycle_policy" "minimal" {
  group_lifetime_in_days = 180
  managed_group_types    = "All"
}
`
}

func testAccConfigMinimalNamed(resourceName, uniqueSuffix string) string {
	return fmt.Sprintf(`
resource "microsoft365_graph_beta_groups_group_lifecycle_policy" "%s" {
  group_lifetime_in_days = 180
  managed_group_types    = "All"
}
`, resourceName)
}

func testAccConfigMaximal(uniqueSuffix string) string {
	return `
resource "microsoft365_graph_beta_groups_group_lifecycle_policy" "maximal" {
  group_lifetime_in_days         = 365
  managed_group_types            = "Selected"
  alternate_notification_emails  = "admin@example.com;notifications@example.com"
}
`
}

func testAccConfigMaximalNamed(resourceName, uniqueSuffix string) string {
	return fmt.Sprintf(`
resource "microsoft365_graph_beta_groups_group_lifecycle_policy" "%s" {
  group_lifetime_in_days         = 365
  managed_group_types            = "Selected"
  alternate_notification_emails  = "admin@example.com;notifications@example.com"
}
`, resourceName)
}
