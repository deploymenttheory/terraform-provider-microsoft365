package graphBetaUpdateWindowsDeviceAccount_test

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

func TestUpdateWindowsDeviceAccountAction_Basic(t *testing.T) {
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

func TestUpdateWindowsDeviceAccountAction_Maximal(t *testing.T) {
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

func TestUpdateWindowsDeviceAccountAction_ComanagedOnly(t *testing.T) {
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
action "microsoft365_graph_beta_device_management_managed_device_update_windows_device_account" "comanaged_only" {
  config {
    comanaged_devices = [
      {
        device_id                 = "00000000-0000-0000-0000-000000000001"
        device_account_email      = "teams-room@company.com"
        password                  = "SecurePassword123!"
        password_rotation_enabled = true
        calendar_sync_enabled     = true
      }
    ]
  }
}
`,
			},
		},
	})
}

func TestUpdateWindowsDeviceAccountAction_PartialFailures(t *testing.T) {
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
action "microsoft365_graph_beta_device_management_managed_device_update_windows_device_account" "partial_failures" {
  config {
    managed_devices = [
      {
        device_id                 = "00000000-0000-0000-0000-000000000001"
        device_account_email      = "conference-room@company.com"
        password                  = "SecurePassword123!"
        password_rotation_enabled = true
        calendar_sync_enabled     = true
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

func TestUpdateWindowsDeviceAccountAction_ValidationDisabled(t *testing.T) {
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
action "microsoft365_graph_beta_device_management_managed_device_update_windows_device_account" "validation_disabled" {
  config {
    managed_devices = [
      {
        device_id                 = "00000000-0000-0000-0000-000000000001"
        device_account_email      = "conference-room@company.com"
        password                  = "SecurePassword123!"
        password_rotation_enabled = true
        calendar_sync_enabled     = true
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

func TestUpdateWindowsDeviceAccountAction_ValidationEnabled(t *testing.T) {
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
action "microsoft365_graph_beta_device_management_managed_device_update_windows_device_account" "validation_enabled" {
  config {
    managed_devices = [
      {
        device_id                 = "00000000-0000-0000-0000-000000000001"
        device_account_email      = "conference-room@company.com"
        password                  = "SecurePassword123!"
        password_rotation_enabled = true
        calendar_sync_enabled     = true
      }
    ]
    validate_device_exists = true
  }
}
`,
			},
		},
	})
}

func TestUpdateWindowsDeviceAccountAction_InvalidDeviceID(t *testing.T) {
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
action "microsoft365_graph_beta_device_management_managed_device_update_windows_device_account" "invalid_device_id" {
  config {
    managed_devices = [
      {
        device_id                 = "not-a-valid-guid"
        device_account_email      = "conference-room@company.com"
        password                  = "SecurePassword123!"
        password_rotation_enabled = true
        calendar_sync_enabled     = true
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

func TestUpdateWindowsDeviceAccountAction_NoDevices(t *testing.T) {
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
action "microsoft365_graph_beta_device_management_managed_device_update_windows_device_account" "no_devices" {
  config {}
}
`,
				ExpectError: regexp.MustCompile(`No Devices Specified`),
			},
		},
	})
}

func TestUpdateWindowsDeviceAccountAction_MissingPassword(t *testing.T) {
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
action "microsoft365_graph_beta_device_management_managed_device_update_windows_device_account" "missing_password" {
  config {
    managed_devices = [
      {
        device_id                 = "00000000-0000-0000-0000-000000000001"
        device_account_email      = "conference-room@company.com"
        password_rotation_enabled = true
        calendar_sync_enabled     = true
      }
    ]
  }
}
`,
				ExpectError: regexp.MustCompile(`Inappropriate value for attribute "managed_devices": element 0: attribute\s+"password" is required`),
			},
		},
	})
}

func TestUpdateWindowsDeviceAccountAction_DuplicateDeviceIDs(t *testing.T) {
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
action "microsoft365_graph_beta_device_management_managed_device_update_windows_device_account" "duplicate_device_ids" {
  config {
    managed_devices = [
      {
        device_id                 = "00000000-0000-0000-0000-000000000001"
        device_account_email      = "conference-room@company.com"
        password                  = "SecurePassword123!"
        password_rotation_enabled = true
        calendar_sync_enabled     = true
      },
      {
        device_id                 = "00000000-0000-0000-0000-000000000001"
        device_account_email      = "another-room@company.com"
        password                  = "AnotherPassword456!"
        password_rotation_enabled = true
        calendar_sync_enabled     = true
      }
    ]
  }
}
`,
			},
		},
	})
}

