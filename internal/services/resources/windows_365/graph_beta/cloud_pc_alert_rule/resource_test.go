package graphBetaCloudPcAlertRule_test

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	cloudPcAlertRuleMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/windows_365/graph_beta/cloud_pc_alert_rule/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

func setupMockEnvironment() (*cloudPcAlertRuleMocks.CloudPcAlertRuleMock, *cloudPcAlertRuleMocks.CloudPcAlertRuleMock) {
	httpmock.Activate()
	mock := &cloudPcAlertRuleMocks.CloudPcAlertRuleMock{}
	errorMock := &cloudPcAlertRuleMocks.CloudPcAlertRuleMock{}
	return mock, errorMock
}

func setupTestEnvironment(t *testing.T) {
	// Set up any test-specific environment variables or configurations here if needed
}

// testCheckExists is a basic check to ensure the resource exists in the state
func testCheckExists(resourceName string) resource.TestCheckFunc {
	return resource.TestCheckResourceAttrSet(resourceName, "id")
}

// testConfigMinimal returns the minimal configuration for testing
func testConfigMinimal() string {
	content, err := os.ReadFile(filepath.Join("mocks", "terraform", "resource_minimal.tf"))
	if err != nil {
		return ""
	}
	return string(content)
}

// testConfigMaximal returns the maximal configuration for testing
func testConfigMaximal() string {
	content, err := os.ReadFile(filepath.Join("mocks", "terraform", "resource_maximal.tf"))
	if err != nil {
		return ""
	}
	return string(content)
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

	// Fix the display name to match test expectations
	updated = strings.Replace(updated, "Test Maximal Cloud PC Alert Rule - Unique", "Test Maximal Cloud PC Alert Rule", 1)

	return updated
}

// Helper function to get minimal config with a custom resource name
func testConfigMinimalWithResourceName(resourceName string) string {
	return fmt.Sprintf(`resource "microsoft365_graph_beta_windows_365_cloud_pc_alert_rule" "%s" {
  alert_rule_template = "cloudPcProvisionScenario"
  display_name   = "Test Minimal Cloud PC Alert Rule"
  severity       = "warning"
  enabled        = true
  is_system_rule = false

  notification_channels = [
    {
      notification_channel_type = "portal"
      notification_receivers = [
        {
          contact_information = "admin@test.com"
          locale             = "en-US"
        }
      ]
    }
  ]

  threshold = {
    aggregation = "count"
    operator    = "greaterOrEqual"
    target      = 1
  }

  conditions = [
    {
      relationship_type   = "and"
      condition_category  = "provisionFailures"
      aggregation        = "count"
      operator           = "greaterOrEqual"
      threshold_value    = "1"
    }
  ]
  
  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}`, resourceName)
}

// TestUnitResourceCloudPcAlertRule_01_CreateMinimal tests the creation of a Cloud PC alert rule with minimal configuration
func TestUnitResourceCloudPcAlertRule_01_CreateMinimal(t *testing.T) {
	// Set up mock environment
	_, _ = setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	// Set up the test environment
	setupTestEnvironment(t)

	// Register the mocks
	mock := &cloudPcAlertRuleMocks.CloudPcAlertRuleMock{}
	mock.RegisterMocks()

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_windows_365_cloud_pc_alert_rule.minimal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_alert_rule.minimal", "alert_rule_template", "cloudPcProvisionScenario"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_alert_rule.minimal", "display_name", "Test Minimal Cloud PC Alert Rule - Unique"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_alert_rule.minimal", "severity", "warning"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_alert_rule.minimal", "enabled", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_alert_rule.minimal", "notification_channels.#", "1"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_alert_rule.minimal", "notification_channels.0.notification_channel_type", "portal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_alert_rule.minimal", "notification_channels.0.notification_receivers.#", "1"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_alert_rule.minimal", "notification_channels.0.notification_receivers.0.contact_information", "admin@test.com"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_alert_rule.minimal", "threshold.aggregation", "count"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_alert_rule.minimal", "threshold.operator", "greaterOrEqual"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_alert_rule.minimal", "conditions.#", "1"),
				),
			},
		},
	})
}

// TestUnitResourceCloudPcAlertRule_02_CreateMaximal tests the creation of a Cloud PC alert rule with maximal configuration
func TestUnitResourceCloudPcAlertRule_02_CreateMaximal(t *testing.T) {
	// Set up mock environment
	_, _ = setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	// Set up the test environment
	setupTestEnvironment(t)

	// Register the mocks
	mock := &cloudPcAlertRuleMocks.CloudPcAlertRuleMock{}
	mock.RegisterMocks()

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMaximal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_windows_365_cloud_pc_alert_rule.maximal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_alert_rule.maximal", "alert_rule_template", "cloudPcProvisionScenario"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_alert_rule.maximal", "display_name", "Test Maximal Cloud PC Alert Rule - Unique"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_alert_rule.maximal", "description", "Comprehensive alert rule for testing Cloud PC provisioning failures with all features"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_alert_rule.maximal", "severity", "critical"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_alert_rule.maximal", "enabled", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_alert_rule.maximal", "notification_channels.#", "2"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_alert_rule.maximal", "threshold.aggregation", "count"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_alert_rule.maximal", "threshold.operator", "greaterOrEqual"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_alert_rule.maximal", "threshold.target", "5"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_alert_rule.maximal", "conditions.#", "2"),
					// Check notification channels
					resource.TestCheckTypeSetElemNestedAttrs("microsoft365_graph_beta_windows_365_cloud_pc_alert_rule.maximal", "notification_channels.*", map[string]string{
						"notification_channel_type": "portal",
					}),
					resource.TestCheckTypeSetElemNestedAttrs("microsoft365_graph_beta_windows_365_cloud_pc_alert_rule.maximal", "notification_channels.*", map[string]string{
						"notification_channel_type": "email",
					}),
				),
			},
		},
	})
}

// TestUnitResourceCloudPcAlertRule_03_UpdateMinimalToMaximal tests updating from minimal to maximal configuration
func TestUnitResourceCloudPcAlertRule_03_UpdateMinimalToMaximal(t *testing.T) {
	// Set up mock environment
	_, _ = setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	// Set up the test environment
	setupTestEnvironment(t)

	// Register the mocks
	mock := &cloudPcAlertRuleMocks.CloudPcAlertRuleMock{}
	mock.RegisterMocks()

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Start with minimal configuration
			{
				Config: testConfigMinimalWithResourceName("test"),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_windows_365_cloud_pc_alert_rule.test"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_alert_rule.test", "display_name", "Test Minimal Cloud PC Alert Rule"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_alert_rule.test", "severity", "warning"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_alert_rule.test", "notification_channels.#", "1"),
				),
			},
			// Update to maximal configuration (with the same resource name)
			{
				Config: testConfigMaximalWithResourceName("test"),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_windows_365_cloud_pc_alert_rule.test"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_alert_rule.test", "display_name", "Test Maximal Cloud PC Alert Rule"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_alert_rule.test", "severity", "critical"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_alert_rule.test", "description", "Comprehensive alert rule for testing Cloud PC provisioning failures with all features"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_alert_rule.test", "notification_channels.#", "2"),
				),
			},
		},
	})
}

// TestUnitResourceCloudPcAlertRule_04_UpdateMaximalToMinimal tests updating from maximal to minimal configuration
func TestUnitResourceCloudPcAlertRule_04_UpdateMaximalToMinimal(t *testing.T) {
	// Set up mock environment
	_, _ = setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	// Set up the test environment
	setupTestEnvironment(t)

	// Register the mocks
	mock := &cloudPcAlertRuleMocks.CloudPcAlertRuleMock{}
	mock.RegisterMocks()

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Start with maximal configuration
			{
				Config: testConfigMaximalWithResourceName("test"),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_windows_365_cloud_pc_alert_rule.test"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_alert_rule.test", "display_name", "Test Maximal Cloud PC Alert Rule"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_alert_rule.test", "severity", "critical"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_alert_rule.test", "notification_channels.#", "2"),
				),
			},
			// Update to minimal configuration (with the same resource name)
			{
				Config: testConfigMinimalWithResourceName("test"),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_windows_365_cloud_pc_alert_rule.test"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_alert_rule.test", "display_name", "Test Minimal Cloud PC Alert Rule"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_alert_rule.test", "severity", "warning"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_alert_rule.test", "notification_channels.#", "1"),
					resource.TestCheckNoResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_alert_rule.test", "description"),
				),
			},
		},
	})
}

// TestUnitResourceCloudPcAlertRule_05_DeleteMinimal tests deleting a Cloud PC alert rule with minimal configuration
func TestUnitResourceCloudPcAlertRule_05_DeleteMinimal(t *testing.T) {
	// Set up mock environment
	_, _ = setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	// Set up the test environment
	setupTestEnvironment(t)

	// Register the mocks
	mock := &cloudPcAlertRuleMocks.CloudPcAlertRuleMock{}
	mock.RegisterMocks()

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_windows_365_cloud_pc_alert_rule.minimal"),
				),
			},
		},
	})
}

// TestUnitResourceCloudPcAlertRule_06_DeleteMaximal tests deleting a Cloud PC alert rule with maximal configuration
func TestUnitResourceCloudPcAlertRule_06_DeleteMaximal(t *testing.T) {
	// Set up mock environment
	_, _ = setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	// Set up the test environment
	setupTestEnvironment(t)

	// Register the mocks
	mock := &cloudPcAlertRuleMocks.CloudPcAlertRuleMock{}
	mock.RegisterMocks()

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMaximal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_windows_365_cloud_pc_alert_rule.maximal"),
				),
			},
		},
	})
}

// TestUnitResourceCloudPcAlertRule_07_Import tests importing a Cloud PC alert rule
func TestUnitResourceCloudPcAlertRule_07_Import(t *testing.T) {
	// Set up mock environment
	_, _ = setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	// Set up the test environment
	setupTestEnvironment(t)

	// Register the mocks
	mock := &cloudPcAlertRuleMocks.CloudPcAlertRuleMock{}
	mock.RegisterMocks()

	// Run the test
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_windows_365_cloud_pc_alert_rule.minimal"),
				),
			},
			{
				ResourceName:      "microsoft365_graph_beta_windows_365_cloud_pc_alert_rule.minimal",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// TestUnitResourceCloudPcAlertRule_08_ErrorHandling tests error handling scenarios
func TestUnitResourceCloudPcAlertRule_08_ErrorHandling(t *testing.T) {
	// Set up mock environment
	_, errorMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	// Set up the test environment
	setupTestEnvironment(t)

	// Register error mocks
	errorMock.RegisterErrorMocks()

	// Run the test - this should fail due to the error mocks
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testConfigMinimal(),
				ExpectError: regexp.MustCompile("(error|Error|ERROR)"),
			},
		},
	})
}
