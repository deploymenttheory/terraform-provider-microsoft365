package graphBetaEnableLostModeManagedDevice_test

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

func loadUnitTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/unit/" + filename)
	if err != nil {
		panic("failed to load unit test config " + filename + ": " + err.Error())
	}
	return config
}

func TestUnitActionEnableLostModeManagedDevice_01_Basic(t *testing.T) {
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

func TestUnitActionEnableLostModeManagedDevice_02_Maximal(t *testing.T) {
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

func TestUnitActionEnableLostModeManagedDevice_03_PartialFailures(t *testing.T) {
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
action "microsoft365_graph_beta_device_management_managed_device_enable_lost_mode" "test" {
  config {
    managed_devices = [
      {
        device_id    = "12345678-1234-1234-1234-123456789abc"
        message      = "Device is lost"
        phone_number = "+1234567890"
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

func TestUnitActionEnableLostModeManagedDevice_04_ValidationDisabled(t *testing.T) {
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
action "microsoft365_graph_beta_device_management_managed_device_enable_lost_mode" "test" {
  config {
    managed_devices = [
      {
        device_id    = "12345678-1234-1234-1234-123456789abc"
        message      = "Device is lost"
        phone_number = "+1234567890"
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

func TestUnitActionEnableLostModeManagedDevice_05_CustomTimeout(t *testing.T) {
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
action "microsoft365_graph_beta_device_management_managed_device_enable_lost_mode" "test" {
  config {
    managed_devices = [
      {
        device_id    = "12345678-1234-1234-1234-123456789abc"
        message      = "Device is lost"
        phone_number = "+1234567890"
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

func TestEnableLostModeManagedDeviceAction_NoDevices(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		Steps: []resource.TestStep{
			{
				Config: `
action "microsoft365_graph_beta_device_management_managed_device_enable_lost_mode" "test" {
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

func TestUnitActionEnableLostModeManagedDevice_07_MissingMessage(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		Steps: []resource.TestStep{
			{
				Config: `
action "microsoft365_graph_beta_device_management_managed_device_enable_lost_mode" "test" {
  config {
    managed_devices = [
      {
        device_id    = "12345678-1234-1234-1234-123456789abc"
        phone_number = "+1234567890"
      }
    ]
  }
}
`,
				ExpectError: regexp.MustCompile(`attribute\s+"message" is required`),
			},
		},
	})
}

func TestUnitActionEnableLostModeManagedDevice_08_MissingPhoneNumber(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		Steps: []resource.TestStep{
			{
				Config: `
action "microsoft365_graph_beta_device_management_managed_device_enable_lost_mode" "test" {
  config {
    managed_devices = [
      {
        device_id = "12345678-1234-1234-1234-123456789abc"
        message   = "Device is lost"
      }
    ]
  }
}
`,
				ExpectError: regexp.MustCompile(`attribute\s+"phone_number" is required`),
			},
		},
	})
}

func TestUnitActionEnableLostModeManagedDevice_09_InvalidDeviceID(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		Steps: []resource.TestStep{
			{
				Config: `
action "microsoft365_graph_beta_device_management_managed_device_enable_lost_mode" "test" {
  config {
    managed_devices = [
      {
        device_id    = "not-a-valid-guid"
        message      = "Device is lost"
        phone_number = "+1234567890"
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

func TestUnitActionEnableLostModeManagedDevice_10_BothManagedAndComanaged(t *testing.T) {
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
action "microsoft365_graph_beta_device_management_managed_device_enable_lost_mode" "test" {
  config {
    managed_devices = [
      {
        device_id    = "12345678-1234-1234-1234-123456789abc"
        message      = "Device is lost"
        phone_number = "+1234567890"
      }
    ]
    comanaged_devices = [
      {
        device_id    = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
        message      = "Co-managed device is lost"
        phone_number = "+0987654321"
      }
    ]
  }
}
`,
			},
		},
	})
}
