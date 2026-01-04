package graphBetaCreateDeviceLogCollectionRequestManagedDevice_test

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

// TestCreateDeviceLogCollectionRequestAction_Basic tests basic log collection
func TestCreateDeviceLogCollectionRequestAction_Basic(t *testing.T) {
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

// TestCreateDeviceLogCollectionRequestAction_ConfigValidation tests configuration validation
func TestCreateDeviceLogCollectionRequestAction_ConfigValidation(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		Steps: []resource.TestStep{
			{
				Config: `
action "microsoft365_graph_beta_device_management_managed_device_create_device_log_collection_request" "test" {
  config {
    managed_devices = [
      {
        device_id     = "12345678-1234-1234-1234-123456789abc"
        template_type = "predefined"
      }
    ]
  }
}
`,
			},
		},
	})
}

// TestCreateDeviceLogCollectionRequestAction_Maximal tests action with all features enabled
func TestCreateDeviceLogCollectionRequestAction_Maximal(t *testing.T) {
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

// TestCreateDeviceLogCollectionRequestAction_ComanagedOnly tests co-managed devices only
func TestCreateDeviceLogCollectionRequestAction_ComanagedOnly(t *testing.T) {
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
action "microsoft365_graph_beta_device_management_managed_device_create_device_log_collection_request" "test" {
  config {
    comanaged_devices = [
      {
        device_id     = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
        template_type = "predefined"
      }
    ]
  }
}
`,
			},
		},
	})
}

// TestCreateDeviceLogCollectionRequestAction_IgnorePartialFailures tests ignore_partial_failures flag
func TestCreateDeviceLogCollectionRequestAction_IgnorePartialFailures(t *testing.T) {
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
action "microsoft365_graph_beta_device_management_managed_device_create_device_log_collection_request" "test" {
  config {
    managed_devices = [
      {
        device_id     = "12345678-1234-1234-1234-123456789abc"
        template_type = "predefined"
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

// TestCreateDeviceLogCollectionRequestAction_DisableValidation tests validate_device_exists = false
func TestCreateDeviceLogCollectionRequestAction_DisableValidation(t *testing.T) {
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
action "microsoft365_graph_beta_device_management_managed_device_create_device_log_collection_request" "test" {
  config {
    managed_devices = [
      {
        device_id     = "12345678-1234-1234-1234-123456789abc"
        template_type = "predefined"
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

// TestCreateDeviceLogCollectionRequestAction_CustomTimeout tests custom timeout configuration
func TestCreateDeviceLogCollectionRequestAction_CustomTimeout(t *testing.T) {
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
action "microsoft365_graph_beta_device_management_managed_device_create_device_log_collection_request" "test" {
  config {
    managed_devices = [
      {
        device_id     = "12345678-1234-1234-1234-123456789abc"
        template_type = "predefined"
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

// TestCreateDeviceLogCollectionRequestAction_InvalidGUID tests validation for invalid GUID format
func TestCreateDeviceLogCollectionRequestAction_InvalidGUID(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		Steps: []resource.TestStep{
			{
				Config: `
action "microsoft365_graph_beta_device_management_managed_device_create_device_log_collection_request" "test" {
  config {
    managed_devices = [
      {
        device_id     = "invalid-guid"
        template_type = "predefined"
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

// TestCreateDeviceLogCollectionRequestAction_InvalidTemplateType tests validation for invalid template type
func TestCreateDeviceLogCollectionRequestAction_InvalidTemplateType(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		Steps: []resource.TestStep{
			{
				Config: `
action "microsoft365_graph_beta_device_management_managed_device_create_device_log_collection_request" "test" {
  config {
    managed_devices = [
      {
        device_id     = "12345678-1234-1234-1234-123456789abc"
        template_type = "invalid"
      }
    ]
  }
}
`,
				ExpectError: regexp.MustCompile(`value must be one of`),
			},
		},
	})
}

// TestCreateDeviceLogCollectionRequestAction_NoDevices tests validation requiring at least one device
func TestCreateDeviceLogCollectionRequestAction_NoDevices(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		Steps: []resource.TestStep{
			{
				Config: `
action "microsoft365_graph_beta_device_management_managed_device_create_device_log_collection_request" "test" {
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

// TestCreateDeviceLogCollectionRequestAction_EmptyDeviceLists tests validation for empty device lists
func TestCreateDeviceLogCollectionRequestAction_EmptyDeviceLists(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		Steps: []resource.TestStep{
			{
				Config: `
action "microsoft365_graph_beta_device_management_managed_device_create_device_log_collection_request" "test" {
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
