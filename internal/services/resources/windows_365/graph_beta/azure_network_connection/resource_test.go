package graphBetaAzureNetworkConnection_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaAzureNetworkConnection "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/windows_365/graph_beta/azure_network_connection"
	localMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/windows_365/graph_beta/azure_network_connection/mocks"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/jarcoal/httpmock"
	msgraphbetamodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	"github.com/stretchr/testify/require"
)

// Helper functions to return the test configurations by reading from files
func testConfigMinimal() string {
	content, err := os.ReadFile(filepath.Join("mocks", "terraform", "resource_minimal.tf"))
	if err != nil {
		return ""
	}
	return string(content)
}

func testConfigMaximal() string {
	content, err := os.ReadFile(filepath.Join("mocks", "terraform", "resource_maximal.tf"))
	if err != nil {
		return ""
	}
	return string(content)
}

func testConfigMinimalToMaximal() string {
	// For minimal to maximal test, we need to use the maximal config
	// but with the minimal resource name to simulate an update

	// Read the maximal config
	maximalContent, err := os.ReadFile(filepath.Join("mocks", "terraform", "resource_maximal.tf"))
	if err != nil {
		return ""
	}

	// Replace the resource name to match the minimal one
	updatedMaximal := strings.Replace(string(maximalContent), "maximal", "minimal", 1)

	// Replace the subscription ID to match the minimal one
	updatedMaximal = strings.Replace(updatedMaximal, "22222222-2222-2222-2222-222222222222", "11111111-1111-1111-1111-111111111111", -1)

	return updatedMaximal
}

func testConfigError() string {
	// Create an error configuration with invalid resource group ID
	return `
resource "microsoft365_graph_beta_windows_365_azure_network_connection" "error" {
  display_name         = "Test Error Connection"
  connection_type      = "hybridAzureADJoin"
  ad_domain_name      = "example.local"
  ad_domain_username  = "testuser"
  ad_domain_password  = "TestPassword123!"
  resource_group_id   = "invalid-resource-group-id"
  subnet_id           = "/subscriptions/11111111-1111-1111-1111-111111111111/resourceGroups/test-rg/providers/Microsoft.Network/virtualNetworks/test-vnet/subnets/test-subnet"
  subscription_id     = "11111111-1111-1111-1111-111111111111"
  virtual_network_id  = "/subscriptions/11111111-1111-1111-1111-111111111111/resourceGroups/test-rg/providers/Microsoft.Network/virtualNetworks/test-vnet"
  
  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}
`
}

// Helper function to set up the test environment
func setupTestEnvironment(t *testing.T) {
	// Set environment variables for testing
	os.Setenv("TF_ACC", "0")
	os.Setenv("MS365_TEST_MODE", "true")
}

// Helper function to set up the mock environment
func setupMockEnvironment() (*mocks.Mocks, *localMocks.AzureNetworkConnectionMock) {
	// Activate httpmock
	httpmock.Activate()

	// Create a new Mocks instance and register authentication mocks
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	// Register local mocks directly
	connectionMock := &localMocks.AzureNetworkConnectionMock{}
	connectionMock.RegisterMocks()

	return mockClient, connectionMock
}

// Helper function to check if a resource exists
func testCheckExists(resourceName string) resource.TestCheckFunc {
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

// Helper function to get maximal config with a custom resource name
func testConfigMaximalWithResourceName(resourceName string) string {
	// Read the maximal config
	content, err := os.ReadFile(filepath.Join("mocks", "terraform", "resource_maximal.tf"))
	if err != nil {
		return ""
	}

	// Replace the resource name
	updated := strings.Replace(string(content), "maximal", resourceName, 1)

	return updated
}

// Helper function to get minimal config with a custom resource name
func testConfigMinimalWithResourceName(resourceName string) string {
	return fmt.Sprintf(`resource "microsoft365_graph_beta_windows_365_azure_network_connection" "%s" {
  display_name         = "Test Minimal Connection"
  connection_type      = "hybridAzureADJoin"
  ad_domain_name      = "example.local"
  ad_domain_username  = "testuser"
  ad_domain_password  = "TestPassword123!"
  resource_group_id   = "/subscriptions/11111111-1111-1111-1111-111111111111/resourceGroups/test-rg"
  subnet_id           = "/subscriptions/11111111-1111-1111-1111-111111111111/resourceGroups/test-rg/providers/Microsoft.Network/virtualNetworks/test-vnet/subnets/test-subnet"
  subscription_id     = "11111111-1111-1111-1111-111111111111"
  virtual_network_id  = "/subscriptions/11111111-1111-1111-1111-111111111111/resourceGroups/test-rg/providers/Microsoft.Network/virtualNetworks/test-vnet"
  
  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}`, resourceName)
}

// TestUnitResourceAzureNetworkConnection_01_CreateMinimal tests the creation of an azure network connection with minimal configuration
func TestUnitResourceAzureNetworkConnection_01_CreateMinimal(t *testing.T) {
	// Set up mock environment
	_, _ = setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	// Set up the test environment
	setupTestEnvironment(t)

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_windows_365_azure_network_connection.minimal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_azure_network_connection.minimal", "display_name", "Test Minimal Connection"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_azure_network_connection.minimal", "connection_type", "hybridAzureADJoin"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_azure_network_connection.minimal", "ad_domain_name", "example.local"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_azure_network_connection.minimal", "ad_domain_username", "testuser"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_azure_network_connection.minimal", "subscription_id", "11111111-1111-1111-1111-111111111111"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_azure_network_connection.minimal", "resource_group_id", "/subscriptions/11111111-1111-1111-1111-111111111111/resourceGroups/test-rg"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_azure_network_connection.minimal", "subnet_id", "/subscriptions/11111111-1111-1111-1111-111111111111/resourceGroups/test-rg/providers/Microsoft.Network/virtualNetworks/test-vnet/subnets/test-subnet"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_azure_network_connection.minimal", "virtual_network_id", "/subscriptions/11111111-1111-1111-1111-111111111111/resourceGroups/test-rg/providers/Microsoft.Network/virtualNetworks/test-vnet"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_windows_365_azure_network_connection.minimal", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_azure_network_connection.minimal", "health_check_status", "passed"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_azure_network_connection.minimal", "managed_by", "windows365"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_azure_network_connection.minimal", "in_use", "false"),
				),
			},
		},
	})
}

// TestUnitResourceAzureNetworkConnection_02_CreateMaximal tests the creation of an azure network connection with maximal configuration
func TestUnitResourceAzureNetworkConnection_02_CreateMaximal(t *testing.T) {
	// Set up mock environment
	_, _ = setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	// Set up the test environment
	setupTestEnvironment(t)

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMaximal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_windows_365_azure_network_connection.maximal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_azure_network_connection.maximal", "display_name", "Test Maximal Connection"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_azure_network_connection.maximal", "connection_type", "hybridAzureADJoin"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_azure_network_connection.maximal", "ad_domain_name", "example.local"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_azure_network_connection.maximal", "ad_domain_username", "testuser"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_azure_network_connection.maximal", "organizational_unit", "OU=CloudPCs,DC=example,DC=local"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_azure_network_connection.maximal", "subscription_id", "22222222-2222-2222-2222-222222222222"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_azure_network_connection.maximal", "resource_group_id", "/subscriptions/22222222-2222-2222-2222-222222222222/resourceGroups/test-rg-maximal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_azure_network_connection.maximal", "subnet_id", "/subscriptions/22222222-2222-2222-2222-222222222222/resourceGroups/test-rg-maximal/providers/Microsoft.Network/virtualNetworks/test-vnet-maximal/subnets/test-subnet-maximal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_azure_network_connection.maximal", "virtual_network_id", "/subscriptions/22222222-2222-2222-2222-222222222222/resourceGroups/test-rg-maximal/providers/Microsoft.Network/virtualNetworks/test-vnet-maximal"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_windows_365_azure_network_connection.maximal", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_azure_network_connection.maximal", "health_check_status", "passed"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_azure_network_connection.maximal", "managed_by", "windows365"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_azure_network_connection.maximal", "in_use", "false"),
				),
			},
		},
	})
}

// TestUnitResourceAzureNetworkConnection_03_UpdateMinimalToMaximal tests updating from minimal to maximal configuration
func TestUnitResourceAzureNetworkConnection_03_UpdateMinimalToMaximal(t *testing.T) {
	// Set up mock environment
	_, _ = setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	// Set up the test environment
	setupTestEnvironment(t)

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Start with minimal configuration
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_windows_365_azure_network_connection.minimal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_azure_network_connection.minimal", "display_name", "Test Minimal Connection"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_azure_network_connection.minimal", "connection_type", "hybridAzureADJoin"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_azure_network_connection.minimal", "subscription_id", "11111111-1111-1111-1111-111111111111"),
				),
			},
			// Update to maximal configuration (with the same resource name)
			{
				Config: testConfigMinimalToMaximal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_windows_365_azure_network_connection.minimal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_azure_network_connection.minimal", "display_name", "Test Maximal Connection"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_azure_network_connection.minimal", "connection_type", "hybridAzureADJoin"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_azure_network_connection.minimal", "organizational_unit", "OU=CloudPCs,DC=example,DC=local"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_azure_network_connection.minimal", "subscription_id", "11111111-1111-1111-1111-111111111111"),
				),
			},
		},
	})
}

// TestUnitResourceAzureNetworkConnection_04_UpdateMaximalToMinimal tests updating from maximal to minimal configuration
func TestUnitResourceAzureNetworkConnection_04_UpdateMaximalToMinimal(t *testing.T) {
	// Set up mock environment
	_, _ = setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	// Set up the test environment
	setupTestEnvironment(t)

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Start with maximal configuration
			{
				Config: testConfigMaximalWithResourceName("test"),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_windows_365_azure_network_connection.test"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_azure_network_connection.test", "display_name", "Test Maximal Connection"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_azure_network_connection.test", "organizational_unit", "OU=CloudPCs,DC=example,DC=local"),
				),
			},
			// Update to minimal configuration (with the same resource name)
			{
				Config: testConfigMinimalWithResourceName("test"),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_windows_365_azure_network_connection.test"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_azure_network_connection.test", "display_name", "Test Minimal Connection"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_azure_network_connection.test", "connection_type", "hybridAzureADJoin"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_azure_network_connection.test", "subscription_id", "11111111-1111-1111-1111-111111111111"),
				),
			},
		},
	})
}

// TestUnitResourceAzureNetworkConnection_05_DeleteMinimal tests deleting an azure network connection with minimal configuration
func TestUnitResourceAzureNetworkConnection_05_DeleteMinimal(t *testing.T) {
	// Set up mock environment
	_, _ = setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	// Set up the test environment
	setupTestEnvironment(t)

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create the resource
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_windows_365_azure_network_connection.minimal"),
				),
			},
			// Delete the resource (by providing empty config)
			{
				Config: `# Empty config for deletion test`,
				Check: func(s *terraform.State) error {
					// The resource should be gone
					_, exists := s.RootModule().Resources["microsoft365_graph_beta_windows_365_azure_network_connection.minimal"]
					if exists {
						return fmt.Errorf("resource still exists after deletion")
					}
					return nil
				},
			},
		},
	})
}

// TestUnitResourceAzureNetworkConnection_06_DeleteMaximal tests deleting an azure network connection with maximal configuration
func TestUnitResourceAzureNetworkConnection_06_DeleteMaximal(t *testing.T) {
	// Set up mock environment
	_, _ = setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	// Set up the test environment
	setupTestEnvironment(t)

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create the resource
			{
				Config: testConfigMaximal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_windows_365_azure_network_connection.maximal"),
				),
			},
			// Delete the resource (by providing empty config)
			{
				Config: `# Empty config for deletion test`,
				Check: func(s *terraform.State) error {
					// The resource should be gone
					_, exists := s.RootModule().Resources["microsoft365_graph_beta_windows_365_azure_network_connection.maximal"]
					if exists {
						return fmt.Errorf("resource still exists after deletion")
					}
					return nil
				},
			},
		},
	})
}

// TestUnitResourceAzureNetworkConnection_07_Import tests importing a resource
func TestUnitResourceAzureNetworkConnection_07_Import(t *testing.T) {
	// Set up mock environment
	_, _ = setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	// Set up the test environment
	setupTestEnvironment(t)

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_windows_365_azure_network_connection.minimal"),
				),
			},
			// Import
			{
				ResourceName:      "microsoft365_graph_beta_windows_365_azure_network_connection.minimal",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"ad_domain_password", // This is sensitive and not returned by the API
				},
			},
		},
	})
}

// TestUnitResourceAzureNetworkConnection_08_Error tests error handling
func TestUnitResourceAzureNetworkConnection_08_Error(t *testing.T) {
	// Set up mock environment
	_, connectionMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	// Register error mocks
	connectionMock.RegisterErrorMocks()

	// Set up the test environment
	setupTestEnvironment(t)

	// Run the test with an error case
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testConfigError(),
				ExpectError: regexp.MustCompile("Must be a valid Azure resource group ID"),
			},
		},
	})
}

func TestUnitResourceAzureNetworkConnection_09_AzureIDCasingLifecycle(t *testing.T) {
	_, connectionMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	mocks.SetupUnitTestEnvironment(t)
	connectionMock.AzureIDCasing = strings.ToLower

	const subscription = "ABCDEFAB-1234-5678-9ABC-ABCDEFABCDEF"
	const updatedSubscription = "ABCDEFAB-1234-5678-9ABC-123456ABCDEF"
	const address = graphBetaAzureNetworkConnection.ResourceName + ".minimal"
	config := strings.NewReplacer(
		"11111111-1111-1111-1111-111111111111", subscription,
		"test-rg", "ExampleGroup", "test-vnet", "ExampleNetwork", "test-subnet", "CloudPCs",
	).Replace(testConfigMinimal())
	updatedConfig := strings.NewReplacer(subscription, updatedSubscription, "ExampleGroup", "UpdatedGroup").Replace(config)

	checkIDs := func(subscriptionID, groupName string) resource.TestCheckFunc {
		group := "/subscriptions/" + subscriptionID + "/resourceGroups/" + groupName
		network := group + "/providers/Microsoft.Network/virtualNetworks/ExampleNetwork"
		return resource.ComposeTestCheckFunc(
			check.That(address).Key("resource_group_id").HasValue(group),
			check.That(address).Key("virtual_network_id").HasValue(network),
			check.That(address).Key("subnet_id").HasValue(network+"/subnets/CloudPCs"),
			check.That(address).Key("subscription_id").HasValue(subscriptionID),
		)
	}
	noChanges := resource.ConfigPlanChecks{PostApplyPostRefresh: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()}}
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check:  checkIDs(subscription, "ExampleGroup"),
			},
			{
				PreConfig:        func() { connectionMock.AzureIDCasing = strings.ToUpper },
				Config:           config,
				PlanOnly:         true,
				ConfigPlanChecks: noChanges,
			},
			{
				Config: updatedConfig,
				Check:  checkIDs(updatedSubscription, "UpdatedGroup"),
				ConfigPlanChecks: resource.ConfigPlanChecks{PreApply: []plancheck.PlanCheck{
					plancheck.ExpectResourceAction(address, plancheck.ResourceActionUpdate),
				}},
			},
			{
				PreConfig:        func() { connectionMock.AzureIDCasing = strings.ToLower },
				Config:           updatedConfig,
				PlanOnly:         true,
				ConfigPlanChecks: noChanges,
			},
		},
	})
}

func TestUnitResourceAzureNetworkConnection_10_AzureIDStateMapping(t *testing.T) {
	const subscription = "ABCDEFAB-1234-5678-9ABC-ABCDEFABCDEF"
	const group = "/subscriptions/" + subscription + "/resourceGroups/ExampleGroup"
	const network = group + "/providers/Microsoft.Network/virtualNetworks/ExampleNetwork"
	model := func(value func(string) types.String) graphBetaAzureNetworkConnection.CloudPcOnPremisesConnectionResourceModel {
		return graphBetaAzureNetworkConnection.CloudPcOnPremisesConnectionResourceModel{
			ResourceGroupId: value(group), VirtualNetworkId: value(network),
			SubnetId: value(network + "/subnets/CloudPCs"), SubscriptionId: value(subscription),
		}
	}
	lower := func(s string) types.String { return types.StringValue(strings.ToLower(s)) }
	upper := func(s string) types.String { return types.StringValue(strings.ToUpper(s)) }
	null := func(string) types.String { return types.StringNull() }
	unknown := func(string) types.String { return types.StringUnknown() }
	empty := func(string) types.String { return types.StringValue("") }
	changed := func(s string) types.String {
		return types.StringValue(strings.ReplaceAll(s, "ABCDEFAB", "123456AB"))
	}
	for _, tc := range []struct {
		name   string
		prior  func(string) types.String
		remote func(string) types.String
		want   func(string) types.String
	}{
		{"unchanged", types.StringValue, types.StringValue, types.StringValue},
		{"lowercase_response", types.StringValue, lower, types.StringValue},
		{"uppercase_response", lower, upper, lower},
		{"different_id", types.StringValue, changed, changed},
		{"import_without_prior_ids", null, types.StringValue, types.StringValue},
		{"unknown_prior_ids", unknown, upper, upper},
		{"missing_remote_ids", types.StringValue, null, null},
		{"empty_remote_ids", types.StringValue, empty, empty},
		{"null_prior_empty_remote", null, empty, empty},
		{"unknown_prior_empty_remote", unknown, empty, empty},
	} {
		t.Run(tc.name, func(t *testing.T) {
			data, remoteData, want := model(tc.prior), model(tc.remote), model(tc.want)
			remote := msgraphbetamodels.NewCloudPcOnPremisesConnection()
			remote.SetResourceGroupId(remoteData.ResourceGroupId.ValueStringPointer())
			remote.SetVirtualNetworkId(remoteData.VirtualNetworkId.ValueStringPointer())
			remote.SetSubnetId(remoteData.SubnetId.ValueStringPointer())
			remote.SetSubscriptionId(remoteData.SubscriptionId.ValueStringPointer())

			graphBetaAzureNetworkConnection.MapRemoteStateToTerraform(context.Background(), &data, remote, types.StringNull())

			require.Equal(t, want.ResourceGroupId, data.ResourceGroupId)
			require.Equal(t, want.VirtualNetworkId, data.VirtualNetworkId)
			require.Equal(t, want.SubnetId, data.SubnetId)
			require.Equal(t, want.SubscriptionId, data.SubscriptionId)
		})
	}
}
