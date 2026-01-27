package graphBetaDeleteUserFromSharedAppleDevice_test

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

// TestDeleteUserFromSharedAppleDeviceAction_Basic tests basic user deletion
func TestUnitActionDeleteUserFromSharedAppleDevice_01_Basic(t *testing.T) {
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

// TestDeleteUserFromSharedAppleDeviceAction_ConfigValidation tests configuration validation
func TestUnitActionDeleteUserFromSharedAppleDevice_02_ConfigValidation(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		Steps: []resource.TestStep{
			{
				Config: `
action "microsoft365_graph_beta_device_management_managed_device_delete_user_from_shared_apple_device" "test" {
  config {
    managed_devices = [
      {
        device_id           = "12345678-1234-1234-1234-123456789abc"
        user_principal_name = "user@example.com"
      }
    ]
  }
}
`,
			},
		},
	})
}

// TestDeleteUserFromSharedAppleDeviceAction_Maximal tests action with all features enabled
func TestUnitActionDeleteUserFromSharedAppleDevice_03_Maximal(t *testing.T) {
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

// TestDeleteUserFromSharedAppleDeviceAction_ComanagedOnly tests co-managed devices only
func TestUnitActionDeleteUserFromSharedAppleDevice_04_ComanagedOnly(t *testing.T) {
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
action "microsoft365_graph_beta_device_management_managed_device_delete_user_from_shared_apple_device" "test" {
  config {
    comanaged_devices = [
      {
        device_id           = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
        user_principal_name = "user@example.com"
      }
    ]
  }
}
`,
			},
		},
	})
}

// TestDeleteUserFromSharedAppleDeviceAction_IgnorePartialFailures tests ignore_partial_failures flag
func TestUnitActionDeleteUserFromSharedAppleDevice_05_IgnorePartialFailures(t *testing.T) {
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
action "microsoft365_graph_beta_device_management_managed_device_delete_user_from_shared_apple_device" "test" {
  config {
    managed_devices = [
      {
        device_id           = "12345678-1234-1234-1234-123456789abc"
        user_principal_name = "user@example.com"
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

// TestDeleteUserFromSharedAppleDeviceAction_DisableValidation tests validate_device_exists = false
func TestUnitActionDeleteUserFromSharedAppleDevice_06_DisableValidation(t *testing.T) {
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
action "microsoft365_graph_beta_device_management_managed_device_delete_user_from_shared_apple_device" "test" {
  config {
    managed_devices = [
      {
        device_id           = "12345678-1234-1234-1234-123456789abc"
        user_principal_name = "user@example.com"
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

// TestDeleteUserFromSharedAppleDeviceAction_CustomTimeout tests custom timeout configuration
func TestUnitActionDeleteUserFromSharedAppleDevice_07_CustomTimeout(t *testing.T) {
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
action "microsoft365_graph_beta_device_management_managed_device_delete_user_from_shared_apple_device" "test" {
  config {
    managed_devices = [
      {
        device_id           = "12345678-1234-1234-1234-123456789abc"
        user_principal_name = "user@example.com"
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

// TestDeleteUserFromSharedAppleDeviceAction_InvalidGUID tests validation for invalid GUID format
func TestUnitActionDeleteUserFromSharedAppleDevice_08_InvalidGUID(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		Steps: []resource.TestStep{
			{
				Config: `
action "microsoft365_graph_beta_device_management_managed_device_delete_user_from_shared_apple_device" "test" {
  config {
    managed_devices = [
      {
        device_id           = "invalid-guid"
        user_principal_name = "user@example.com"
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

// TestDeleteUserFromSharedAppleDeviceAction_NoDevices tests validation requiring at least one device
func TestUnitActionDeleteUserFromSharedAppleDevice_09_NoDevices(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		Steps: []resource.TestStep{
			{
				Config: `
action "microsoft365_graph_beta_device_management_managed_device_delete_user_from_shared_apple_device" "test" {
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

// TestDeleteUserFromSharedAppleDeviceAction_EmptyDeviceLists tests validation for empty device lists
func TestUnitActionDeleteUserFromSharedAppleDevice_10_EmptyDeviceLists(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		Steps: []resource.TestStep{
			{
				Config: `
action "microsoft365_graph_beta_device_management_managed_device_delete_user_from_shared_apple_device" "test" {
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
