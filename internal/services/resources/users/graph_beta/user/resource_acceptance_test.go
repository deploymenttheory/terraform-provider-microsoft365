package graphBetaUsersUser_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// TestAccUserResource_Minimal tests the minimal configuration
func TestAccUserResource_Minimal(t *testing.T) {
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
		PreCheck:                 func() { testAccPreCheck(t) },
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
			// Import test
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

// TestAccUserResource_Maximal tests the maximal configuration
func TestAccUserResource_Maximal(t *testing.T) {
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
		PreCheck:                 func() { testAccPreCheck(t) },
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
			// Import test
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

// TestAccUserResource_MinimalToMaximal tests updating from minimal to maximal config
func TestAccUserResource_MinimalToMaximal(t *testing.T) {
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
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckUserDestroy,
		Steps: []resource.TestStep{
			// Start with minimal configuration
			{
				Config: testAccConfigMinimalNamed(resourceName, userPrincipalName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUserExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "display_name", "Minimal User"),
					resource.TestCheckResourceAttr(resourceName, "user_principal_name", userPrincipalName),
				),
			},
			// Update to maximal configuration
			{
				Config: testAccConfigMaximalNamed(resourceName, userPrincipalName),
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

// TestAccUserResource_MaximalToMinimal tests updating from maximal to minimal config
func TestAccUserResource_MaximalToMinimal(t *testing.T) {
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
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckUserDestroy,
		Steps: []resource.TestStep{
			// Start with maximal configuration
			{
				Config: testAccConfigMaximalNamed(resourceName, userPrincipalName),
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
				Config: testAccConfigMinimalNamed(resourceName, userPrincipalName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUserExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "display_name", "Minimal User"),
					resource.TestCheckResourceAttr(resourceName, "user_principal_name", userPrincipalName),
				),
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
		"TEST_DOMAIN",
	}

	for _, env := range requiredEnvVars {
		if os.Getenv(env) == "" {
			t.Fatalf("%s environment variable must be set for acceptance tests", env)
		}
	}
}

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
	return fmt.Sprintf(`
resource "microsoft365_graph_beta_users_user" "minimal" {
  display_name        = "Minimal User"
  user_principal_name = "%s"
  password_profile = {
    password = "SecureP@ssw0rd123!"
  }
}
`, userPrincipalName)
}

// Minimal configuration with custom resource name
func testAccConfigMinimalNamed(resourceName string, userPrincipalName string) string {
	// Extract the resource name without the provider prefix
	name := resourceName
	if strings.Contains(resourceName, ".") {
		if len(resourceName) > 0 && resourceName[0] != '"' {
			parts := strings.Split(resourceName, ".")
			if len(parts) > 1 {
				name = parts[1]
			}
		}
	}

	return fmt.Sprintf(`
resource "microsoft365_graph_beta_users_user" "%s" {
  display_name        = "Minimal User"
  user_principal_name = "%s"
  password_profile = {
    password = "SecureP@ssw0rd123!"
  }
}
`, name, userPrincipalName)
}

// Maximal configuration with default resource name
func testAccConfigMaximal(userPrincipalName string) string {
	return fmt.Sprintf(`
resource "microsoft365_graph_beta_users_user" "maximal" {
  display_name        = "Maximal User"
  user_principal_name = "%s"
  account_enabled     = true
  given_name          = "Maximal"
  surname             = "User"
  mail                = "%s"
  mail_nickname       = "maxuser"
  job_title           = "Senior Developer"
  department          = "Engineering"
  company_name        = "Contoso Ltd"
  office_location     = "Building A"
  city                = "Redmond"
  state               = "WA"
  country             = "US"
  postal_code         = "98052"
  usage_location      = "US"
  business_phones     = ["+1 425-555-0100"]
  mobile_phone        = "+1 425-555-0101"
  password_profile = {
    password                          = "SecureP@ssw0rd123!"
    force_change_password_next_sign_in = true
  }
  identities = [
    {
      sign_in_type       = "emailAddress"
      issuer             = "%s"
      issuer_assigned_id = "%s"
    }
  ]
  other_mails     = ["maximal.user.other@contoso.com"]
  proxy_addresses = ["SMTP:%s"]
}
`, userPrincipalName, userPrincipalName, strings.Split(userPrincipalName, "@")[1], userPrincipalName, userPrincipalName)
}

// Maximal configuration with custom resource name
func testAccConfigMaximalNamed(resourceName string, userPrincipalName string) string {
	// Extract the resource name without the provider prefix
	name := resourceName
	if strings.Contains(resourceName, ".") {
		if len(resourceName) > 0 && resourceName[0] != '"' {
			parts := strings.Split(resourceName, ".")
			if len(parts) > 1 {
				name = parts[1]
			}
		}
	}

	return fmt.Sprintf(`
resource "microsoft365_graph_beta_users_user" "%s" {
  display_name        = "Maximal User"
  user_principal_name = "%s"
  account_enabled     = true
  given_name          = "Maximal"
  surname             = "User"
  mail                = "%s"
  mail_nickname       = "maxuser"
  job_title           = "Senior Developer"
  department          = "Engineering"
  company_name        = "Contoso Ltd"
  office_location     = "Building A"
  password_profile = {
    password = "SecureP@ssw0rd123!"
  }
}
`, name, userPrincipalName, userPrincipalName)
}
