package graphBetaCleanWindowsManagedDevice_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/jarcoal/httpmock"
)

// Helper function to load test configs from unit directory
func loadUnitTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/unit/" + filename)
	if err != nil {
		panic("failed to load unit test config " + filename + ": " + err.Error())
	}
	return config
}

// TestCleanWindowsDeviceAction_Basic tests basic Windows device clean
func TestUnitActionCleanWindowsManagedDevice_01_Basic(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Register auth mocks
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		CheckDestroy: func(s *terraform.State) error {
			// Actions don't create persistent state, so nothing to destroy
			return nil
		},
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("action_minimal.tf"),
				// Actions don't persist in state, so we just verify the configuration is valid
				// The test passes if Terraform accepts the config and doesn't error during plan
			},
		},
	})
}

// TestCleanWindowsDeviceAction_ConfigValidation tests configuration validation
func TestUnitActionCleanWindowsManagedDevice_02_ConfigValidation(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		Steps: []resource.TestStep{
			{
				Config: `
action "microsoft365_graph_beta_device_management_managed_device_clean_windows_device" "test" {
  config {
    managed_devices = [
      {
        device_id      = "12345678-1234-1234-1234-123456789abc"
        keep_user_data = false
      }
    ]
  }
}
`,
				// Verify the action configuration is accepted
			},
		},
	})
}

// TestCleanWindowsDeviceAction_Maximal tests action with all features enabled
func TestUnitActionCleanWindowsManagedDevice_03_Maximal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		CheckDestroy: func(s *terraform.State) error {
			return nil
		},
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("action_maximal.tf"),
			},
		},
	})
}

// TestCleanWindowsDeviceAction_ComanagedOnly tests co-managed devices only
func TestUnitActionCleanWindowsManagedDevice_04_ComanagedOnly(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		CheckDestroy: func(s *terraform.State) error {
			return nil
		},
		Steps: []resource.TestStep{
			{
				Config: `
action "microsoft365_graph_beta_device_management_managed_device_clean_windows_device" "test" {
  config {
    comanaged_devices = [
      {
        device_id      = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
        keep_user_data = false
      }
    ]
  }
}
`,
			},
		},
	})
}

// TestCleanWindowsDeviceAction_KeepUserData tests keep_user_data flag
func TestUnitActionCleanWindowsManagedDevice_05_KeepUserData(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		CheckDestroy: func(s *terraform.State) error {
			return nil
		},
		Steps: []resource.TestStep{
			{
				Config: `
action "microsoft365_graph_beta_device_management_managed_device_clean_windows_device" "test" {
  config {
    managed_devices = [
      {
        device_id      = "12345678-1234-1234-1234-123456789abc"
        keep_user_data = true
      }
    ]
  }
}
`,
			},
		},
	})
}

// TestCleanWindowsDeviceAction_IgnorePartialFailures tests ignore_partial_failures flag
func TestUnitActionCleanWindowsManagedDevice_06_IgnorePartialFailures(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		CheckDestroy: func(s *terraform.State) error {
			return nil
		},
		Steps: []resource.TestStep{
			{
				Config: `
action "microsoft365_graph_beta_device_management_managed_device_clean_windows_device" "test" {
  config {
    managed_devices = [
      {
        device_id      = "12345678-1234-1234-1234-123456789abc"
        keep_user_data = false
      }
    ]
    ignore_partial_failures = true
  }
}
`,
			},
		},
	})
}

// TestCleanWindowsDeviceAction_DisableValidation tests validate_device_exists = false
func TestUnitActionCleanWindowsManagedDevice_07_DisableValidation(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		CheckDestroy: func(s *terraform.State) error {
			return nil
		},
		Steps: []resource.TestStep{
			{
				Config: `
action "microsoft365_graph_beta_device_management_managed_device_clean_windows_device" "test" {
  config {
    managed_devices = [
      {
        device_id      = "12345678-1234-1234-1234-123456789abc"
        keep_user_data = false
      }
    ]
    validate_device_exists = false
  }
}
`,
			},
		},
	})
}

// TestCleanWindowsDeviceAction_CustomTimeout tests custom timeout configuration
func TestUnitActionCleanWindowsManagedDevice_08_CustomTimeout(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		CheckDestroy: func(s *terraform.State) error {
			return nil
		},
		Steps: []resource.TestStep{
			{
				Config: `
action "microsoft365_graph_beta_device_management_managed_device_clean_windows_device" "test" {
  config {
    managed_devices = [
      {
        device_id      = "12345678-1234-1234-1234-123456789abc"
        keep_user_data = false
      }
    ]
    timeouts = {
      invoke = "10m"
    }
  }
}
`,
			},
		},
	})
}

// TestCleanWindowsDeviceAction_InvalidGUID tests validation for invalid GUID format
func TestUnitActionCleanWindowsManagedDevice_09_InvalidGUID(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		Steps: []resource.TestStep{
			{
				Config: `
action "microsoft365_graph_beta_device_management_managed_device_clean_windows_device" "test" {
  config {
    managed_devices = [
      {
        device_id      = "invalid-guid"
        keep_user_data = false
      }
    ]
  }
}
`,
				ExpectError: regexp.MustCompile(`device_id must be a valid GUID format`),
			},
		},
	})
}

// TestCleanWindowsDeviceAction_NoDevices tests validation requiring at least one device
func TestUnitActionCleanWindowsManagedDevice_10_NoDevices(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		Steps: []resource.TestStep{
			{
				Config: `
action "microsoft365_graph_beta_device_management_managed_device_clean_windows_device" "test" {
  config {
    timeouts = {
      invoke = "5m"
    }
  }
}
`,
				ExpectError: regexp.MustCompile(`No Devices Specified`),
			},
		},
	})
}

// TestCleanWindowsDeviceAction_EmptyDeviceLists tests validation for empty device lists
func TestUnitActionCleanWindowsManagedDevice_11_EmptyDeviceLists(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		Steps: []resource.TestStep{
			{
				Config: `
action "microsoft365_graph_beta_device_management_managed_device_clean_windows_device" "test" {
  config {
    managed_devices   = []
    comanaged_devices = []
  }
}
`,
				ExpectError: regexp.MustCompile(`No Devices Specified`),
			},
		},
	})
}

// TestCleanWindowsDeviceAction_MissingKeepUserData tests validation requiring keep_user_data
func TestUnitActionCleanWindowsManagedDevice_12_MissingKeepUserData(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		Steps: []resource.TestStep{
			{
				Config: `
action "microsoft365_graph_beta_device_management_managed_device_clean_windows_device" "test" {
  config {
    managed_devices = [
      {
        device_id = "12345678-1234-1234-1234-123456789abc"
      }
    ]
  }
}
`,
				ExpectError: regexp.MustCompile(`element 0: attribute\s+"keep_user_data" is required`),
			},
		},
	})
}
