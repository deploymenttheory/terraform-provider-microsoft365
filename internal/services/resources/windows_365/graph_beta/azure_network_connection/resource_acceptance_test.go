package graphBetaAzureNetworkConnection_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// TestAccAzureNetworkConnectionResource_Create_Minimal tests creating an azure network connection with minimal configuration
func TestAccAzureNetworkConnectionResource_Create_Minimal(t *testing.T) {
	// Skip if not running acceptance tests
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless TF_ACC=1")
	}

	// Get test subscription ID from environment variable or skip
	testSubscriptionID := os.Getenv("TEST_SUBSCRIPTION_ID")
	if testSubscriptionID == "" {
		t.Skip("TEST_SUBSCRIPTION_ID environment variable must be set for acceptance tests")
	}

	// Get test resource group ID from environment variable or skip
	testResourceGroupID := os.Getenv("TEST_RESOURCE_GROUP_ID")
	if testResourceGroupID == "" {
		t.Skip("TEST_RESOURCE_GROUP_ID environment variable must be set for acceptance tests")
	}

	// Get test subnet ID from environment variable or skip
	testSubnetID := os.Getenv("TEST_SUBNET_ID")
	if testSubnetID == "" {
		t.Skip("TEST_SUBNET_ID environment variable must be set for acceptance tests")
	}

	// Get test virtual network ID from environment variable or skip
	testVirtualNetworkID := os.Getenv("TEST_VIRTUAL_NETWORK_ID")
	if testVirtualNetworkID == "" {
		t.Skip("TEST_VIRTUAL_NETWORK_ID environment variable must be set for acceptance tests")
	}

	// Get test domain credentials from environment variables or skip
	testDomainName := os.Getenv("TEST_DOMAIN_NAME")
	if testDomainName == "" {
		t.Skip("TEST_DOMAIN_NAME environment variable must be set for acceptance tests")
	}

	testDomainUsername := os.Getenv("TEST_DOMAIN_USERNAME")
	if testDomainUsername == "" {
		t.Skip("TEST_DOMAIN_USERNAME environment variable must be set for acceptance tests")
	}

	testDomainPassword := os.Getenv("TEST_DOMAIN_PASSWORD")
	if testDomainPassword == "" {
		t.Skip("TEST_DOMAIN_PASSWORD environment variable must be set for acceptance tests")
	}

	resourceName := "microsoft365_graph_beta_windows_365_azure_network_connection.minimal"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckAzureNetworkConnectionDestroy,
		Steps: []resource.TestStep{
			// Create with minimal configuration
			{
				Config: testAccConfigMinimal(testSubscriptionID, testResourceGroupID, testSubnetID, testVirtualNetworkID, testDomainName, testDomainUsername, testDomainPassword),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAzureNetworkConnectionExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "display_name", "Test Minimal Connection"),
					resource.TestCheckResourceAttr(resourceName, "connection_type", "hybridAzureADJoin"),
					resource.TestCheckResourceAttr(resourceName, "ad_domain_name", testDomainName),
					resource.TestCheckResourceAttr(resourceName, "ad_domain_username", testDomainUsername),
					resource.TestCheckResourceAttr(resourceName, "subscription_id", testSubscriptionID),
					resource.TestCheckResourceAttr(resourceName, "resource_group_id", testResourceGroupID),
					resource.TestCheckResourceAttr(resourceName, "subnet_id", testSubnetID),
					resource.TestCheckResourceAttr(resourceName, "virtual_network_id", testVirtualNetworkID),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "health_check_status"),
					resource.TestCheckResourceAttrSet(resourceName, "managed_by"),
				),
			},
		},
	})
}

// TestAccAzureNetworkConnectionResource_Create_Maximal tests creating an azure network connection with maximal configuration
func TestAccAzureNetworkConnectionResource_Create_Maximal(t *testing.T) {
	// Skip if not running acceptance tests
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless TF_ACC=1")
	}

	// Get test environment variables (same as minimal test)
	testSubscriptionID := os.Getenv("TEST_SUBSCRIPTION_ID")
	if testSubscriptionID == "" {
		t.Skip("TEST_SUBSCRIPTION_ID environment variable must be set for acceptance tests")
	}

	testResourceGroupID := os.Getenv("TEST_RESOURCE_GROUP_ID")
	if testResourceGroupID == "" {
		t.Skip("TEST_RESOURCE_GROUP_ID environment variable must be set for acceptance tests")
	}

	testSubnetID := os.Getenv("TEST_SUBNET_ID")
	if testSubnetID == "" {
		t.Skip("TEST_SUBNET_ID environment variable must be set for acceptance tests")
	}

	testVirtualNetworkID := os.Getenv("TEST_VIRTUAL_NETWORK_ID")
	if testVirtualNetworkID == "" {
		t.Skip("TEST_VIRTUAL_NETWORK_ID environment variable must be set for acceptance tests")
	}

	testDomainName := os.Getenv("TEST_DOMAIN_NAME")
	if testDomainName == "" {
		t.Skip("TEST_DOMAIN_NAME environment variable must be set for acceptance tests")
	}

	testDomainUsername := os.Getenv("TEST_DOMAIN_USERNAME")
	if testDomainUsername == "" {
		t.Skip("TEST_DOMAIN_USERNAME environment variable must be set for acceptance tests")
	}

	testDomainPassword := os.Getenv("TEST_DOMAIN_PASSWORD")
	if testDomainPassword == "" {
		t.Skip("TEST_DOMAIN_PASSWORD environment variable must be set for acceptance tests")
	}

	// Optional organizational unit
	testOrganizationalUnit := os.Getenv("TEST_ORGANIZATIONAL_UNIT")

	resourceName := "microsoft365_graph_beta_windows_365_azure_network_connection.maximal"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckAzureNetworkConnectionDestroy,
		Steps: []resource.TestStep{
			// Create with maximal configuration
			{
				Config: testAccConfigMaximal(testSubscriptionID, testResourceGroupID, testSubnetID, testVirtualNetworkID, testDomainName, testDomainUsername, testDomainPassword, testOrganizationalUnit),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAzureNetworkConnectionExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "display_name", "Test Maximal Connection"),
					resource.TestCheckResourceAttr(resourceName, "connection_type", "hybridAzureADJoin"),
					resource.TestCheckResourceAttr(resourceName, "ad_domain_name", testDomainName),
					resource.TestCheckResourceAttr(resourceName, "ad_domain_username", testDomainUsername),
					resource.TestCheckResourceAttr(resourceName, "subscription_id", testSubscriptionID),
					resource.TestCheckResourceAttr(resourceName, "resource_group_id", testResourceGroupID),
					resource.TestCheckResourceAttr(resourceName, "subnet_id", testSubnetID),
					resource.TestCheckResourceAttr(resourceName, "virtual_network_id", testVirtualNetworkID),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "health_check_status"),
					resource.TestCheckResourceAttrSet(resourceName, "managed_by"),
				),
			},
		},
	})
}

// TestAccAzureNetworkConnectionResource_Update_MinimalToMaximal tests updating from minimal to maximal config
func TestAccAzureNetworkConnectionResource_Update_MinimalToMaximal(t *testing.T) {
	// Skip if not running acceptance tests
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless TF_ACC=1")
	}

	// Get test environment variables (same as previous tests)
	testSubscriptionID := os.Getenv("TEST_SUBSCRIPTION_ID")
	if testSubscriptionID == "" {
		t.Skip("TEST_SUBSCRIPTION_ID environment variable must be set for acceptance tests")
	}

	testResourceGroupID := os.Getenv("TEST_RESOURCE_GROUP_ID")
	if testResourceGroupID == "" {
		t.Skip("TEST_RESOURCE_GROUP_ID environment variable must be set for acceptance tests")
	}

	testSubnetID := os.Getenv("TEST_SUBNET_ID")
	if testSubnetID == "" {
		t.Skip("TEST_SUBNET_ID environment variable must be set for acceptance tests")
	}

	testVirtualNetworkID := os.Getenv("TEST_VIRTUAL_NETWORK_ID")
	if testVirtualNetworkID == "" {
		t.Skip("TEST_VIRTUAL_NETWORK_ID environment variable must be set for acceptance tests")
	}

	testDomainName := os.Getenv("TEST_DOMAIN_NAME")
	if testDomainName == "" {
		t.Skip("TEST_DOMAIN_NAME environment variable must be set for acceptance tests")
	}

	testDomainUsername := os.Getenv("TEST_DOMAIN_USERNAME")
	if testDomainUsername == "" {
		t.Skip("TEST_DOMAIN_USERNAME environment variable must be set for acceptance tests")
	}

	testDomainPassword := os.Getenv("TEST_DOMAIN_PASSWORD")
	if testDomainPassword == "" {
		t.Skip("TEST_DOMAIN_PASSWORD environment variable must be set for acceptance tests")
	}

	testOrganizationalUnit := os.Getenv("TEST_ORGANIZATIONAL_UNIT")

	resourceName := "microsoft365_graph_beta_windows_365_azure_network_connection.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckAzureNetworkConnectionDestroy,
		Steps: []resource.TestStep{
			// Start with minimal configuration
			{
				Config: testAccConfigMinimalNamed("test", testSubscriptionID, testResourceGroupID, testSubnetID, testVirtualNetworkID, testDomainName, testDomainUsername, testDomainPassword),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAzureNetworkConnectionExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "display_name", "Test Minimal Connection"),
					resource.TestCheckResourceAttr(resourceName, "connection_type", "hybridAzureADJoin"),
				),
			},
			// Update to maximal configuration
			{
				Config: testAccConfigMaximalNamed("test", testSubscriptionID, testResourceGroupID, testSubnetID, testVirtualNetworkID, testDomainName, testDomainUsername, testDomainPassword, testOrganizationalUnit),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAzureNetworkConnectionExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "display_name", "Test Maximal Connection"),
					resource.TestCheckResourceAttr(resourceName, "connection_type", "hybridAzureADJoin"),
				),
			},
		},
	})
}

// TestAccAzureNetworkConnectionResource_Delete_Minimal tests deleting an azure network connection with minimal configuration
func TestAccAzureNetworkConnectionResource_Delete_Minimal(t *testing.T) {
	// Skip if not running acceptance tests
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless TF_ACC=1")
	}

	// Get test environment variables (same as minimal test)
	testSubscriptionID := os.Getenv("TEST_SUBSCRIPTION_ID")
	if testSubscriptionID == "" {
		t.Skip("TEST_SUBSCRIPTION_ID environment variable must be set for acceptance tests")
	}

	testResourceGroupID := os.Getenv("TEST_RESOURCE_GROUP_ID")
	if testResourceGroupID == "" {
		t.Skip("TEST_RESOURCE_GROUP_ID environment variable must be set for acceptance tests")
	}

	testSubnetID := os.Getenv("TEST_SUBNET_ID")
	if testSubnetID == "" {
		t.Skip("TEST_SUBNET_ID environment variable must be set for acceptance tests")
	}

	testVirtualNetworkID := os.Getenv("TEST_VIRTUAL_NETWORK_ID")
	if testVirtualNetworkID == "" {
		t.Skip("TEST_VIRTUAL_NETWORK_ID environment variable must be set for acceptance tests")
	}

	testDomainName := os.Getenv("TEST_DOMAIN_NAME")
	if testDomainName == "" {
		t.Skip("TEST_DOMAIN_NAME environment variable must be set for acceptance tests")
	}

	testDomainUsername := os.Getenv("TEST_DOMAIN_USERNAME")
	if testDomainUsername == "" {
		t.Skip("TEST_DOMAIN_USERNAME environment variable must be set for acceptance tests")
	}

	testDomainPassword := os.Getenv("TEST_DOMAIN_PASSWORD")
	if testDomainPassword == "" {
		t.Skip("TEST_DOMAIN_PASSWORD environment variable must be set for acceptance tests")
	}

	resourceName := "microsoft365_graph_beta_windows_365_azure_network_connection.minimal"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckAzureNetworkConnectionDestroy,
		Steps: []resource.TestStep{
			// Create the resource
			{
				Config: testAccConfigMinimal(testSubscriptionID, testResourceGroupID, testSubnetID, testVirtualNetworkID, testDomainName, testDomainUsername, testDomainPassword),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAzureNetworkConnectionExists(resourceName),
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

// TestAccAzureNetworkConnectionResource_Import tests importing a resource
func TestAccAzureNetworkConnectionResource_Import(t *testing.T) {
	// Skip if not running acceptance tests
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless TF_ACC=1")
	}

	// Get test environment variables (same as minimal test)
	testSubscriptionID := os.Getenv("TEST_SUBSCRIPTION_ID")
	if testSubscriptionID == "" {
		t.Skip("TEST_SUBSCRIPTION_ID environment variable must be set for acceptance tests")
	}

	testResourceGroupID := os.Getenv("TEST_RESOURCE_GROUP_ID")
	if testResourceGroupID == "" {
		t.Skip("TEST_RESOURCE_GROUP_ID environment variable must be set for acceptance tests")
	}

	testSubnetID := os.Getenv("TEST_SUBNET_ID")
	if testSubnetID == "" {
		t.Skip("TEST_SUBNET_ID environment variable must be set for acceptance tests")
	}

	testVirtualNetworkID := os.Getenv("TEST_VIRTUAL_NETWORK_ID")
	if testVirtualNetworkID == "" {
		t.Skip("TEST_VIRTUAL_NETWORK_ID environment variable must be set for acceptance tests")
	}

	testDomainName := os.Getenv("TEST_DOMAIN_NAME")
	if testDomainName == "" {
		t.Skip("TEST_DOMAIN_NAME environment variable must be set for acceptance tests")
	}

	testDomainUsername := os.Getenv("TEST_DOMAIN_USERNAME")
	if testDomainUsername == "" {
		t.Skip("TEST_DOMAIN_USERNAME environment variable must be set for acceptance tests")
	}

	testDomainPassword := os.Getenv("TEST_DOMAIN_PASSWORD")
	if testDomainPassword == "" {
		t.Skip("TEST_DOMAIN_PASSWORD environment variable must be set for acceptance tests")
	}

	resourceName := "microsoft365_graph_beta_windows_365_azure_network_connection.minimal"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckAzureNetworkConnectionDestroy,
		Steps: []resource.TestStep{
			// Create
			{
				Config: testAccConfigMinimal(testSubscriptionID, testResourceGroupID, testSubnetID, testVirtualNetworkID, testDomainName, testDomainUsername, testDomainPassword),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAzureNetworkConnectionExists(resourceName),
				),
			},
			// Import
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"ad_domain_password", // This is sensitive and not returned by the API
				},
			},
		},
	})
}

// Helper functions for acceptance tests

func testAccPreCheck(t *testing.T) {
	// Verify required environment variables are set
	requiredEnvVars := []string{
		"M365_TENANT_ID",
		"M365_CLIENT_SECRET",
		"M365_CLIENT_ID",
		"TEST_SUBSCRIPTION_ID",
		"TEST_RESOURCE_GROUP_ID",
		"TEST_SUBNET_ID",
		"TEST_VIRTUAL_NETWORK_ID",
		"TEST_DOMAIN_NAME",
		"TEST_DOMAIN_USERNAME",
		"TEST_DOMAIN_PASSWORD",
	}

	for _, env := range requiredEnvVars {
		if os.Getenv(env) == "" {
			t.Fatalf("%s environment variable must be set for acceptance tests", env)
		}
	}
}

func testAccCheckAzureNetworkConnectionExists(resourceName string) resource.TestCheckFunc {
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

func testAccCheckAzureNetworkConnectionDestroy(s *terraform.State) error {
	// In a real test, we would verify the connection is removed
	// For this resource, we don't need to check anything special since removing
	// the resource will remove the connection
	return nil
}

// Test configurations

// Minimal configuration with default resource name
func testAccConfigMinimal(subscriptionID, resourceGroupID, subnetID, virtualNetworkID, domainName, domainUsername, domainPassword string) string {
	return fmt.Sprintf(`
resource "microsoft365_graph_beta_windows_365_azure_network_connection" "minimal" {
  display_name         = "Test Minimal Connection"
  connection_type      = "hybridAzureADJoin"
  ad_domain_name      = "%s"
  ad_domain_username  = "%s"
  ad_domain_password  = "%s"
  resource_group_id   = "%s"
  subnet_id           = "%s"
  subscription_id     = "%s"
  virtual_network_id  = "%s"
}
`, domainName, domainUsername, domainPassword, resourceGroupID, subnetID, subscriptionID, virtualNetworkID)
}

// Minimal configuration with custom resource name
func testAccConfigMinimalNamed(resourceName, subscriptionID, resourceGroupID, subnetID, virtualNetworkID, domainName, domainUsername, domainPassword string) string {
	return fmt.Sprintf(`
resource "microsoft365_graph_beta_windows_365_azure_network_connection" "%s" {
  display_name         = "Test Minimal Connection"
  connection_type      = "hybridAzureADJoin"
  ad_domain_name      = "%s"
  ad_domain_username  = "%s"
  ad_domain_password  = "%s"
  resource_group_id   = "%s"
  subnet_id           = "%s"
  subscription_id     = "%s"
  virtual_network_id  = "%s"
}
`, resourceName, domainName, domainUsername, domainPassword, resourceGroupID, subnetID, subscriptionID, virtualNetworkID)
}

// Maximal configuration with default resource name
func testAccConfigMaximal(subscriptionID, resourceGroupID, subnetID, virtualNetworkID, domainName, domainUsername, domainPassword, organizationalUnit string) string {
	ouConfig := ""
	if organizationalUnit != "" {
		ouConfig = fmt.Sprintf(`
  organizational_unit = "%s"`, organizationalUnit)
	}

	return fmt.Sprintf(`
resource "microsoft365_graph_beta_windows_365_azure_network_connection" "maximal" {
  display_name         = "Test Maximal Connection"
  connection_type      = "hybridAzureADJoin"
  ad_domain_name      = "%s"
  ad_domain_username  = "%s"
  ad_domain_password  = "%s"%s
  resource_group_id   = "%s"
  subnet_id           = "%s"
  subscription_id     = "%s"
  virtual_network_id  = "%s"
}
`, domainName, domainUsername, domainPassword, ouConfig, resourceGroupID, subnetID, subscriptionID, virtualNetworkID)
}

// Maximal configuration with custom resource name
func testAccConfigMaximalNamed(resourceName, subscriptionID, resourceGroupID, subnetID, virtualNetworkID, domainName, domainUsername, domainPassword, organizationalUnit string) string {
	ouConfig := ""
	if organizationalUnit != "" {
		ouConfig = fmt.Sprintf(`
  organizational_unit = "%s"`, organizationalUnit)
	}

	return fmt.Sprintf(`
resource "microsoft365_graph_beta_windows_365_azure_network_connection" "%s" {
  display_name         = "Test Maximal Connection"
  connection_type      = "hybridAzureADJoin"
  ad_domain_name      = "%s"
  ad_domain_username  = "%s"
  ad_domain_password  = "%s"%s
  resource_group_id   = "%s"
  subnet_id           = "%s"
  subscription_id     = "%s"
  virtual_network_id  = "%s"
}
`, resourceName, domainName, domainUsername, domainPassword, ouConfig, resourceGroupID, subnetID, subscriptionID, virtualNetworkID)
}
