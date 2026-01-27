package graphBetaDeprovisionManagedDevice_test

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

// TestDeprovisionManagedDeviceAction_Basic tests basic device deprovision
func TestUnitActionDeprovisionManagedDevice_01_Basic(t *testing.T) {
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
				Config: loadUnitTestTerraform("action_minimal.tf"),
			},
		},
	})
}

// TestDeprovisionManagedDeviceAction_ConfigValidation tests configuration validation
func TestUnitActionDeprovisionManagedDevice_02_ConfigValidation(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		Steps: []resource.TestStep{
			{
				Config: `
action "microsoft365_graph_beta_device_management_managed_device_deprovision" "test" {
  config {
    managed_devices = [
      {
        device_id          = "12345678-1234-1234-1234-123456789abc"
        deprovision_reason = "Test"
      }
    ]
  }
}
`,
			},
		},
	})
}

// TestDeprovisionManagedDeviceAction_Maximal tests action with all features enabled
func TestUnitActionDeprovisionManagedDevice_03_Maximal(t *testing.T) {
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

// TestDeprovisionManagedDeviceAction_ComanagedOnly tests co-managed devices only
func TestUnitActionDeprovisionManagedDevice_04_ComanagedOnly(t *testing.T) {
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
action "microsoft365_graph_beta_device_management_managed_device_deprovision" "test" {
  config {
    comanaged_devices = [
      {
        device_id          = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
        deprovision_reason = "Removing co-management"
      }
    ]
  }
}
`,
			},
		},
	})
}

// TestDeprovisionManagedDeviceAction_IgnorePartialFailures tests ignore_partial_failures flag
func TestUnitActionDeprovisionManagedDevice_05_IgnorePartialFailures(t *testing.T) {
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
action "microsoft365_graph_beta_device_management_managed_device_deprovision" "test" {
  config {
    managed_devices = [
      {
        device_id          = "12345678-1234-1234-1234-123456789abc"
        deprovision_reason = "Test"
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

// TestDeprovisionManagedDeviceAction_DisableValidation tests validate_device_exists = false
func TestUnitActionDeprovisionManagedDevice_06_DisableValidation(t *testing.T) {
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
action "microsoft365_graph_beta_device_management_managed_device_deprovision" "test" {
  config {
    managed_devices = [
      {
        device_id          = "12345678-1234-1234-1234-123456789abc"
        deprovision_reason = "Test"
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

// TestDeprovisionManagedDeviceAction_CustomTimeout tests custom timeout configuration
func TestUnitActionDeprovisionManagedDevice_07_CustomTimeout(t *testing.T) {
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
action "microsoft365_graph_beta_device_management_managed_device_deprovision" "test" {
  config {
    managed_devices = [
      {
        device_id          = "12345678-1234-1234-1234-123456789abc"
        deprovision_reason = "Test"
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

// TestDeprovisionManagedDeviceAction_InvalidGUID tests validation for invalid GUID format
func TestUnitActionDeprovisionManagedDevice_08_InvalidGUID(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		Steps: []resource.TestStep{
			{
				Config: `
action "microsoft365_graph_beta_device_management_managed_device_deprovision" "test" {
  config {
    managed_devices = [
      {
        device_id          = "invalid-guid"
        deprovision_reason = "Test"
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

// TestDeprovisionManagedDeviceAction_NoDevices tests validation requiring at least one device
func TestUnitActionDeprovisionManagedDevice_09_NoDevices(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		Steps: []resource.TestStep{
			{
				Config: `
action "microsoft365_graph_beta_device_management_managed_device_deprovision" "test" {
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

// TestDeprovisionManagedDeviceAction_EmptyDeviceLists tests validation for empty device lists
func TestUnitActionDeprovisionManagedDevice_10_EmptyDeviceLists(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		Steps: []resource.TestStep{
			{
				Config: `
action "microsoft365_graph_beta_device_management_managed_device_deprovision" "test" {
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

// TestDeprovisionManagedDeviceAction_MissingDeprovisionReason tests validation for missing deprovision reason
func TestUnitActionDeprovisionManagedDevice_11_MissingDeprovisionReason(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		Steps: []resource.TestStep{
			{
				Config: `
action "microsoft365_graph_beta_device_management_managed_device_deprovision" "test" {
  config {
    managed_devices = [
      {
        device_id = "12345678-1234-1234-1234-123456789abc"
      }
    ]
  }
}
`,
				ExpectError: regexp.MustCompile(`deprovision_reason.*is required`),
			},
		},
	})
}
