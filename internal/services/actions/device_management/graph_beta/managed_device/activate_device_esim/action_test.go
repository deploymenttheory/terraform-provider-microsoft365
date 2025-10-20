package graphBetaActivateDeviceEsimManagedDevice_test

import (
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	activateDeviceEsimMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/actions/device_management/graph_beta/managed_device/activate_device_esim/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

func setupUnitTestEnvironment(t *testing.T) {
	// Set environment variables for testing
	os.Setenv("TF_ACC", "0")
	os.Setenv("MS365_TEST_MODE", "true")
}

// setupMockEnvironment sets up the mock environment using centralized mocks
func setupMockEnvironment() (*mocks.Mocks, *activateDeviceEsimMocks.ActivateDeviceEsimMock) {
	// Activate httpmock
	httpmock.Activate()

	// Create a new Mocks instance and register authentication mocks
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	// Register local mocks directly
	activateDeviceEsimMock := &activateDeviceEsimMocks.ActivateDeviceEsimMock{}
	activateDeviceEsimMock.RegisterMocks()

	return mockClient, activateDeviceEsimMock
}

// setupErrorMockEnvironment sets up the mock environment for error testing
func setupErrorMockEnvironment() (*mocks.Mocks, *activateDeviceEsimMocks.ActivateDeviceEsimMock) {
	// Activate httpmock
	httpmock.Activate()

	// Create a new Mocks instance and register authentication mocks
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	// Register error mocks
	activateDeviceEsimMock := &activateDeviceEsimMocks.ActivateDeviceEsimMock{}
	activateDeviceEsimMock.RegisterErrorMocks()

	return mockClient, activateDeviceEsimMock
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

// testConfigComanagedOnly returns the co-managed only configuration for testing
func testConfigComanagedOnly() string {
	content, err := os.ReadFile(filepath.Join("tests", "terraform", "unit", "action_comanaged_only.tf"))
	if err != nil {
		return ""
	}
	return string(content)
}

// TestActivateDeviceEsimAction_Schema validates the action schema
func TestActivateDeviceEsimAction_Schema(t *testing.T) {
	setupUnitTestEnvironment(t)
	_, activateDeviceEsimMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer activateDeviceEsimMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					// Check managed devices configuration
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_managed_device_activate_device_esim.minimal", "managed_devices.#", "1"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_managed_device_activate_device_esim.minimal", "managed_devices.0.device_id", "12345678-1234-1234-1234-123456789abc"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_managed_device_activate_device_esim.minimal", "managed_devices.0.carrier_url", "https://carrier.example.com/esim/activate?token=test123"),
					
					// Check default values for simplified options
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_managed_device_activate_device_esim.minimal", "ignore_partial_failures", "false"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_managed_device_activate_device_esim.minimal", "validate_device_exists", "true"),
				),
			},
		},
	})
}

// TestActivateDeviceEsimAction_ManagedDevices tests managed devices activation
func TestActivateDeviceEsimAction_ManagedDevices(t *testing.T) {
	setupUnitTestEnvironment(t)
	_, activateDeviceEsimMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer activateDeviceEsimMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_managed_device_activate_device_esim.minimal", "managed_devices.#", "1"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_managed_device_activate_device_esim.minimal", "managed_devices.0.device_id", "12345678-1234-1234-1234-123456789abc"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_managed_device_activate_device_esim.minimal", "managed_devices.0.carrier_url", "https://carrier.example.com/esim/activate?token=test123"),
				),
			},
		},
	})
}

// TestActivateDeviceEsimAction_ComanagedDevices tests co-managed devices activation
func TestActivateDeviceEsimAction_ComanagedDevices(t *testing.T) {
	setupUnitTestEnvironment(t)
	_, activateDeviceEsimMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer activateDeviceEsimMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigComanagedOnly(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_managed_device_activate_device_esim.comanaged_only", "comanaged_devices.#", "1"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_managed_device_activate_device_esim.comanaged_only", "comanaged_devices.0.device_id", "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_managed_device_activate_device_esim.comanaged_only", "comanaged_devices.0.carrier_url", "https://carrier.example.com/esim/activate?code=comanaged789"),
				),
			},
		},
	})
}

// TestActivateDeviceEsimAction_BothDeviceTypes tests both managed and co-managed devices
func TestActivateDeviceEsimAction_BothDeviceTypes(t *testing.T) {
	setupUnitTestEnvironment(t)
	_, activateDeviceEsimMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer activateDeviceEsimMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMaximal(),
				Check: resource.ComposeTestCheckFunc(
					// Check managed devices
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_managed_device_activate_device_esim.maximal", "managed_devices.#", "2"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_managed_device_activate_device_esim.maximal", "managed_devices.0.device_id", "12345678-1234-1234-1234-123456789abc"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_managed_device_activate_device_esim.maximal", "managed_devices.1.device_id", "87654321-4321-4321-4321-987654321cba"),
					
					// Check co-managed devices
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_managed_device_activate_device_esim.maximal", "comanaged_devices.#", "2"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_managed_device_activate_device_esim.maximal", "comanaged_devices.0.device_id", "11111111-2222-3333-4444-555555555555"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_managed_device_activate_device_esim.maximal", "comanaged_devices.1.device_id", "66666666-7777-8888-9999-000000000000"),
				),
			},
		},
	})
}

// TestActivateDeviceEsimAction_RequiredFields tests required field validation
func TestActivateDeviceEsimAction_RequiredFields(t *testing.T) {
	setupUnitTestEnvironment(t)
	_, activateDeviceEsimMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer activateDeviceEsimMock.CleanupMockState()

	testCases := []struct {
		name          string
		config        string
		expectedError string
	}{
		{
			name: "missing device_id in managed_devices",
			config: `
action "microsoft365_graph_beta_device_management_managed_device_activate_device_esim" "test" {
  managed_devices = [
    {
      carrier_url = "https://carrier.example.com/esim/activate?token=test"
    }
  ]
}
`,
			expectedError: `The argument "device_id" is required`,
		},
		{
			name: "missing carrier_url in managed_devices",
			config: `
action "microsoft365_graph_beta_device_management_managed_device_activate_device_esim" "test" {
  managed_devices = [
    {
      device_id = "12345678-1234-1234-1234-123456789abc"
    }
  ]
}
`,
			expectedError: `The argument "carrier_url" is required`,
		},
		{
			name: "invalid device_id format",
			config: `
action "microsoft365_graph_beta_device_management_managed_device_activate_device_esim" "test" {
  managed_devices = [
    {
      device_id   = "invalid-guid"
      carrier_url = "https://carrier.example.com/esim/activate?token=test"
    }
  ]
}
`,
			expectedError: `device_id must be a valid GUID format`,
		},
		{
			name: "empty carrier_url",
			config: `
action "microsoft365_graph_beta_device_management_managed_device_activate_device_esim" "test" {
  managed_devices = [
    {
      device_id   = "12345678-1234-1234-1234-123456789abc"
      carrier_url = ""
    }
  ]
}
`,
			expectedError: `string length must be at least 1`,
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

// TestActivateDeviceEsimAction_NoDevicesProvided tests validation when no devices are provided
func TestActivateDeviceEsimAction_NoDevicesProvided(t *testing.T) {
	setupUnitTestEnvironment(t)
	_, activateDeviceEsimMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer activateDeviceEsimMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
action "microsoft365_graph_beta_device_management_managed_device_activate_device_esim" "test" {
  timeouts = {
    create = "5m"
  }
}
`,
				ExpectError: regexp.MustCompile(`At least one of managed_devices or comanaged_devices must be provided`),
			},
		},
	})
}

// TestActivateDeviceEsimAction_ErrorHandling tests error scenarios
func TestActivateDeviceEsimAction_ErrorHandling(t *testing.T) {
	setupUnitTestEnvironment(t)
	_, activateDeviceEsimMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer activateDeviceEsimMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
action "microsoft365_graph_beta_device_management_managed_device_activate_device_esim" "test" {
  managed_devices = [
    {
      device_id   = "error-id"
      carrier_url = "https://carrier.example.com/esim/activate?token=error"
    }
  ]
}
`,
				ExpectError: regexp.MustCompile(`Invalid carrier URL|BadRequest`),
			},
		},
	})
}

// TestActivateDeviceEsimAction_DeviceNotFound tests device not found scenarios
func TestActivateDeviceEsimAction_DeviceNotFound(t *testing.T) {
	setupUnitTestEnvironment(t)
	_, activateDeviceEsimMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer activateDeviceEsimMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
action "microsoft365_graph_beta_device_management_managed_device_activate_device_esim" "test" {
  managed_devices = [
    {
      device_id   = "not-found-id"
      carrier_url = "https://carrier.example.com/esim/activate?token=test"
    }
  ]
}
`,
				ExpectError: regexp.MustCompile(`Device not found|NotFound`),
			},
		},
	})
}

// TestActivateDeviceEsimAction_ConfigurationOptions tests the simplified configuration options
func TestActivateDeviceEsimAction_ConfigurationOptions(t *testing.T) {
	setupUnitTestEnvironment(t)
	_, activateDeviceEsimMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer activateDeviceEsimMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
action "microsoft365_graph_beta_device_management_managed_device_activate_device_esim" "config_test" {
  managed_devices = [
    {
      device_id   = "12345678-1234-1234-1234-123456789abc"
      carrier_url = "https://carrier.example.com/esim/activate?token=test123"
    }
  ]
  
  ignore_partial_failures = true
  validate_device_exists   = false
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_managed_device_activate_device_esim.config_test", "ignore_partial_failures", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_managed_device_activate_device_esim.config_test", "validate_device_exists", "false"),
				),
			},
		},
	})
}

// TestActivateDeviceEsimAction_PartialFailureHandling tests partial failure scenarios
func TestActivateDeviceEsimAction_PartialFailureHandling(t *testing.T) {
	setupUnitTestEnvironment(t)
	_, activateDeviceEsimMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer activateDeviceEsimMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
action "microsoft365_graph_beta_device_management_managed_device_activate_device_esim" "partial_failure" {
  managed_devices = [
    {
      device_id   = "12345678-1234-1234-1234-123456789abc"
      carrier_url = "https://carrier.example.com/esim/activate?token=success"
    },
    {
      device_id   = "error-id"
      carrier_url = "https://carrier.example.com/esim/activate?token=error"
    }
  ]
  ignore_partial_failures = true
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_managed_device_activate_device_esim.partial_failure", "ignore_partial_failures", "true"),
				),
			},
		},
	})
}