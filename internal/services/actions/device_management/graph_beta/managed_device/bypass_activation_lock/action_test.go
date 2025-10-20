package graphBetaBypassActivationLockManagedDevice_test

import (
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	bypassActivationLockMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/actions/device_management/graph_beta/managed_device/bypass_activation_lock/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

func setupUnitTestEnvironment(t *testing.T) {
	// Set environment variables for testing
	os.Setenv("TF_ACC", "0")
	os.Setenv("MS365_TEST_MODE", "true")
}

// setupMockEnvironment sets up the mock environment using centralized mocks
func setupMockEnvironment() (*mocks.Mocks, *bypassActivationLockMocks.BypassActivationLockMock) {
	// Activate httpmock
	httpmock.Activate()

	// Create a new Mocks instance and register authentication mocks
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	// Register local mocks directly
	bypassActivationLockMock := &bypassActivationLockMocks.BypassActivationLockMock{}
	bypassActivationLockMock.RegisterMocks()

	return mockClient, bypassActivationLockMock
}

// setupErrorMockEnvironment sets up the mock environment for error testing
func setupErrorMockEnvironment() (*mocks.Mocks, *bypassActivationLockMocks.BypassActivationLockMock) {
	// Activate httpmock
	httpmock.Activate()

	// Create a new Mocks instance and register authentication mocks
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	// Register error mocks
	bypassActivationLockMock := &bypassActivationLockMocks.BypassActivationLockMock{}
	bypassActivationLockMock.RegisterErrorMocks()

	return mockClient, bypassActivationLockMock
}

// testConfigMinimal returns the minimal configuration for testing
func testConfigMinimal() string {
	content, err := os.ReadFile(filepath.Join("tests", "terraform", "unit", "action_minimal.tf"))
	if err != nil {
		return ""
	}
	return string(content)
}

// testConfigMaximal returns the maximal configuration for testing
func testConfigMaximal() string {
	content, err := os.ReadFile(filepath.Join("tests", "terraform", "unit", "action_maximal.tf"))
	if err != nil {
		return ""
	}
	return string(content)
}

// TestBypassActivationLockAction_Schema validates the action schema
func TestBypassActivationLockAction_Schema(t *testing.T) {
	setupUnitTestEnvironment(t)
	_, bypassActivationLockMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer bypassActivationLockMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					// Check device IDs configuration
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_managed_device_bypass_activation_lock.minimal", "device_ids.#", "1"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_managed_device_bypass_activation_lock.minimal", "device_ids.0", "12345678-1234-1234-1234-123456789abc"),

					// Check default values for options
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_managed_device_bypass_activation_lock.minimal", "ignore_partial_failures", "false"),
				),
			},
		},
	})
}

// TestBypassActivationLockAction_SingleDevice tests single device bypass
func TestBypassActivationLockAction_SingleDevice(t *testing.T) {
	setupUnitTestEnvironment(t)
	_, bypassActivationLockMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer bypassActivationLockMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_managed_device_bypass_activation_lock.minimal", "device_ids.#", "1"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_managed_device_bypass_activation_lock.minimal", "device_ids.0", "12345678-1234-1234-1234-123456789abc"),
				),
			},
		},
	})
}

// TestBypassActivationLockAction_MultipleDevices tests multiple device bypass
func TestBypassActivationLockAction_MultipleDevices(t *testing.T) {
	setupUnitTestEnvironment(t)
	_, bypassActivationLockMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer bypassActivationLockMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMaximal(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_managed_device_bypass_activation_lock.maximal", "device_ids.#", "3"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_managed_device_bypass_activation_lock.maximal", "device_ids.0", "12345678-1234-1234-1234-123456789abc"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_managed_device_bypass_activation_lock.maximal", "device_ids.1", "87654321-4321-4321-4321-987654321cba"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_managed_device_bypass_activation_lock.maximal", "device_ids.2", "11111111-2222-3333-4444-555555555555"),
				),
			},
		},
	})
}

// TestBypassActivationLockAction_ConfigurationOptions tests the configuration options
func TestBypassActivationLockAction_ConfigurationOptions(t *testing.T) {
	setupUnitTestEnvironment(t)
	_, bypassActivationLockMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer bypassActivationLockMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
action "microsoft365_graph_beta_device_management_managed_device_bypass_activation_lock" "config_test" {
  device_ids = [
    "12345678-1234-1234-1234-123456789abc"
  ]
  
  ignore_partial_failures = true
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_managed_device_bypass_activation_lock.config_test", "ignore_partial_failures", "true"),
				),
			},
		},
	})
}

// TestBypassActivationLockAction_RequiredFields tests required field validation
func TestBypassActivationLockAction_RequiredFields(t *testing.T) {
	setupUnitTestEnvironment(t)
	_, bypassActivationLockMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer bypassActivationLockMock.CleanupMockState()

	testCases := []struct {
		name          string
		config        string
		expectedError string
	}{
		{
			name: "missing device_ids",
			config: `
action "microsoft365_graph_beta_device_management_managed_device_bypass_activation_lock" "test" {
  timeouts = {
    create = "5m"
  }
}
`,
			expectedError: `The argument "device_ids" is required`,
		},
		{
			name: "empty device_ids list",
			config: `
action "microsoft365_graph_beta_device_management_managed_device_bypass_activation_lock" "test" {
  device_ids = []
}
`,
			expectedError: `Attribute device_ids list must contain at least 1 elements`,
		},
		{
			name: "invalid device_id format",
			config: `
action "microsoft365_graph_beta_device_management_managed_device_bypass_activation_lock" "test" {
  device_ids = ["invalid-guid"]
}
`,
			expectedError: `each device ID must be a valid GUID format`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config:      tc.config,
						ExpectError: regexp.MustCompile(tc.expectedError),
					},
				},
			})
		})
	}
}

// TestBypassActivationLockAction_ErrorHandling tests error scenarios
func TestBypassActivationLockAction_ErrorHandling(t *testing.T) {
	setupUnitTestEnvironment(t)
	_, bypassActivationLockMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer bypassActivationLockMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
action "microsoft365_graph_beta_device_management_managed_device_bypass_activation_lock" "test" {
  device_ids = ["error-id"]
}
`,
				ExpectError: regexp.MustCompile(`BadRequest|Invalid request`),
			},
		},
	})
}

// TestBypassActivationLockAction_DeviceNotFound tests device not found scenarios
func TestBypassActivationLockAction_DeviceNotFound(t *testing.T) {
	setupUnitTestEnvironment(t)
	_, bypassActivationLockMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer bypassActivationLockMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
action "microsoft365_graph_beta_device_management_managed_device_bypass_activation_lock" "test" {
  device_ids = ["not-found-id"]
}
`,
				ExpectError: regexp.MustCompile(`Device not found|NotFound`),
			},
		},
	})
}

// TestBypassActivationLockAction_PartialFailureHandling tests partial failure scenarios
func TestBypassActivationLockAction_PartialFailureHandling(t *testing.T) {
	setupUnitTestEnvironment(t)
	_, bypassActivationLockMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer bypassActivationLockMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
action "microsoft365_graph_beta_device_management_managed_device_bypass_activation_lock" "partial_failure" {
  device_ids = [
    "12345678-1234-1234-1234-123456789abc",
    "error-id"
  ]
  ignore_partial_failures = true
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_managed_device_bypass_activation_lock.partial_failure", "ignore_partial_failures", "true"),
				),
			},
		},
	})
}
