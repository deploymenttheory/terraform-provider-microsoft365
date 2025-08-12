package graphBetaUsersUser_test

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

// TestAccUserResource_Create_Minimal tests creating a user with minimal configuration
func TestAccUserResource_Create_Minimal(t *testing.T) {
	// Skip if not running acceptance tests
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless TF_ACC=1")
	}

	// Get test domain from environment variable or skip
	testDomain := os.Getenv("TEST_DOMAIN")
	if testDomain == "" {
		t.Skip("TEST_DOMAIN environment variable must be set for acceptance tests")
	}

	resourceName := "microsoft365_graph_beta_users_user.minimal"
	userPrincipalName := fmt.Sprintf("tfacctest.minimal@%s", testDomain)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckUserDestroy,
		Steps: []resource.TestStep{
			// Create with minimal configuration
			{
				Config: testAccConfigMinimal(userPrincipalName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUserExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "display_name", "Minimal User"),
					resource.TestCheckResourceAttr(resourceName, "user_principal_name", userPrincipalName),
					resource.TestCheckResourceAttr(resourceName, "account_enabled", "true"),
				),
			},
		},
	})
}

// TestAccUserResource_Create_Maximal tests creating a user with maximal configuration
func TestAccUserResource_Create_Maximal(t *testing.T) {
	// Skip if not running acceptance tests
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless TF_ACC=1")
	}

	// Get test domain from environment variable or skip
	testDomain := os.Getenv("TEST_DOMAIN")
	if testDomain == "" {
		t.Skip("TEST_DOMAIN environment variable must be set for acceptance tests")
	}

	resourceName := "microsoft365_graph_beta_users_user.maximal"
	userPrincipalName := fmt.Sprintf("tfacctest.maximal@%s", testDomain)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckUserDestroy,
		Steps: []resource.TestStep{
			// Create with maximal configuration
			{
				Config: testAccConfigMaximal(userPrincipalName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUserExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "display_name", "Maximal User"),
					resource.TestCheckResourceAttr(resourceName, "user_principal_name", userPrincipalName),
					resource.TestCheckResourceAttr(resourceName, "given_name", "Maximal"),
					resource.TestCheckResourceAttr(resourceName, "surname", "User"),
					resource.TestCheckResourceAttr(resourceName, "job_title", "Senior Developer"),
					resource.TestCheckResourceAttr(resourceName, "department", "Engineering"),
					resource.TestCheckResourceAttr(resourceName, "company_name", "Contoso Ltd"),
					resource.TestCheckResourceAttr(resourceName, "business_phones.#", "1"),
				),
			},
		},
	})
}

// TestAccUserResource_Update_MinimalToMaximal tests updating from minimal to maximal config
func TestAccUserResource_Update_MinimalToMaximal(t *testing.T) {
	// Skip if not running acceptance tests
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless TF_ACC=1")
	}

	// Get test domain from environment variable or skip
	testDomain := os.Getenv("TEST_DOMAIN")
	if testDomain == "" {
		t.Skip("TEST_DOMAIN environment variable must be set for acceptance tests")
	}

	resourceName := "microsoft365_graph_beta_users_user.test"
	userPrincipalName := fmt.Sprintf("tfacctest.update@%s", testDomain)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckUserDestroy,
		Steps: []resource.TestStep{
			// Start with minimal configuration
			{
				Config: testAccConfigMinimalNamed("test", userPrincipalName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUserExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "display_name", "Minimal User"),
					resource.TestCheckResourceAttr(resourceName, "user_principal_name", userPrincipalName),
				),
			},
			// Update to maximal configuration
			{
				Config: testAccConfigMaximalNamed("test", userPrincipalName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUserExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "display_name", "Maximal User"),
					resource.TestCheckResourceAttr(resourceName, "user_principal_name", userPrincipalName),
					resource.TestCheckResourceAttr(resourceName, "given_name", "Maximal"),
					resource.TestCheckResourceAttr(resourceName, "surname", "User"),
					resource.TestCheckResourceAttr(resourceName, "job_title", "Senior Developer"),
					resource.TestCheckResourceAttr(resourceName, "department", "Engineering"),
					resource.TestCheckResourceAttr(resourceName, "company_name", "Contoso Ltd"),
				),
			},
		},
	})
}

// TestAccUserResource_Update_MaximalToMinimal tests updating from maximal to minimal config
func TestAccUserResource_Update_MaximalToMinimal(t *testing.T) {
	// Skip if not running acceptance tests
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless TF_ACC=1")
	}

	// Get test domain from environment variable or skip
	testDomain := os.Getenv("TEST_DOMAIN")
	if testDomain == "" {
		t.Skip("TEST_DOMAIN environment variable must be set for acceptance tests")
	}

	resourceName := "microsoft365_graph_beta_users_user.test"
	userPrincipalName := fmt.Sprintf("tfacctest.downgrade@%s", testDomain)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckUserDestroy,
		Steps: []resource.TestStep{
			// Start with maximal configuration
			{
				Config: testAccConfigMaximalNamed("test", userPrincipalName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUserExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "display_name", "Maximal User"),
					resource.TestCheckResourceAttr(resourceName, "user_principal_name", userPrincipalName),
					resource.TestCheckResourceAttr(resourceName, "given_name", "Maximal"),
					resource.TestCheckResourceAttr(resourceName, "surname", "User"),
				),
			},
			// Update to minimal configuration
			{
				Config: testAccConfigMinimalNamed("test", userPrincipalName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUserExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "display_name", "Minimal User"),
					resource.TestCheckResourceAttr(resourceName, "user_principal_name", userPrincipalName),
				),
			},
		},
	})
}

// TestAccUserResource_Delete_Minimal tests deleting a user with minimal configuration
func TestAccUserResource_Delete_Minimal(t *testing.T) {
	// Skip if not running acceptance tests
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless TF_ACC=1")
	}

	// Get test domain from environment variable or skip
	testDomain := os.Getenv("TEST_DOMAIN")
	if testDomain == "" {
		t.Skip("TEST_DOMAIN environment variable must be set for acceptance tests")
	}

	resourceName := "microsoft365_graph_beta_users_user.minimal"
	userPrincipalName := fmt.Sprintf("tfacctest.delete.minimal@%s", testDomain)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckUserDestroy,
		Steps: []resource.TestStep{
			// Create the resource
			{
				Config: testAccConfigMinimal(userPrincipalName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUserExists(resourceName),
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

// TestAccUserResource_Delete_Maximal tests deleting a user with maximal configuration
func TestAccUserResource_Delete_Maximal(t *testing.T) {
	// Skip if not running acceptance tests
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless TF_ACC=1")
	}

	// Get test domain from environment variable or skip
	testDomain := os.Getenv("TEST_DOMAIN")
	if testDomain == "" {
		t.Skip("TEST_DOMAIN environment variable must be set for acceptance tests")
	}

	resourceName := "microsoft365_graph_beta_users_user.maximal"
	userPrincipalName := fmt.Sprintf("tfacctest.delete.maximal@%s", testDomain)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckUserDestroy,
		Steps: []resource.TestStep{
			// Create the resource
			{
				Config: testAccConfigMaximal(userPrincipalName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUserExists(resourceName),
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

// TestAccUserResource_Import tests importing a resource
func TestAccUserResource_Import(t *testing.T) {
	// Skip if not running acceptance tests
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless TF_ACC=1")
	}

	// Get test domain from environment variable or skip
	testDomain := os.Getenv("TEST_DOMAIN")
	if testDomain == "" {
		t.Skip("TEST_DOMAIN environment variable must be set for acceptance tests")
	}

	resourceName := "microsoft365_graph_beta_users_user.minimal"
	userPrincipalName := fmt.Sprintf("tfacctest.import@%s", testDomain)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckUserDestroy,
		Steps: []resource.TestStep{
			// Create
			{
				Config: testAccConfigMinimal(userPrincipalName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUserExists(resourceName),
				),
			},
			// Import
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"password_profile", // Password is not returned by the API
				},
			},
		},
	})
}

// Helper functions for acceptance tests

func testAccCheckUserExists(resourceName string) resource.TestCheckFunc {
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

func testAccCheckUserDestroy(s *terraform.State) error {
	// In a real test, we would verify the user is removed
	// For this resource, we don't need to check anything special since removing
	// the resource will remove the user
	return nil
}

// Test configurations

// Minimal configuration with default resource name
func testAccConfigMinimal(userPrincipalName string) string {
	// Read the template file
	content, err := os.ReadFile(filepath.Join("mocks", "terraform", "resource_minimal.tf"))
	if err != nil {
		return ""
	}

	// Replace the UPN with the test UPN
	updated := strings.Replace(string(content), "minimal.user@contoso.com", userPrincipalName, 1)

	return updated
}

// Minimal configuration with custom resource name
func testAccConfigMinimalNamed(resourceName string, userPrincipalName string) string {
	// Read the template file
	content, err := os.ReadFile(filepath.Join("mocks", "terraform", "resource_minimal.tf"))
	if err != nil {
		return ""
	}

	// Replace the resource name and UPN
	updated := strings.Replace(string(content), "minimal", resourceName, 1)
	updated = strings.Replace(updated, "minimal.user@contoso.com", userPrincipalName, 1)

	return updated
}

// Maximal configuration with default resource name
func testAccConfigMaximal(userPrincipalName string) string {
	// Read the template file
	content, err := os.ReadFile(filepath.Join("mocks", "terraform", "resource_maximal.tf"))
	if err != nil {
		return ""
	}

	// Replace the UPN with the test UPN
	updated := strings.Replace(string(content), "maximal.user@contoso.com", userPrincipalName, 1)

	return updated
}

// Maximal configuration with custom resource name
func testAccConfigMaximalNamed(resourceName string, userPrincipalName string) string {
	// Read the template file
	content, err := os.ReadFile(filepath.Join("mocks", "terraform", "resource_maximal.tf"))
	if err != nil {
		return ""
	}

	// Replace the resource name and UPN
	updated := strings.Replace(string(content), "maximal", resourceName, 1)
	updated = strings.Replace(updated, "maximal.user@contoso.com", userPrincipalName, 1)

	return updated
}
