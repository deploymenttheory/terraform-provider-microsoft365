package graphBetaTriggerConfigurationManagerActionManagedDevice_test

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

func TestTriggerConfigurationManagerAction_Basic(t *testing.T) {
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

func TestTriggerConfigurationManagerAction_Maximal(t *testing.T) {
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

func TestTriggerConfigurationManagerAction_ComanagedOnly(t *testing.T) {
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
action "microsoft365_graph_beta_device_management_managed_device_trigger_configuration_manager_action" "comanaged_only" {
  config {
    comanaged_devices = [
      {
        device_id = "00000000-0000-0000-0000-000000000001"
        action    = "wakeUpClient"
      }
    ]
  }
}
`,
			},
		},
	})
}

func TestTriggerConfigurationManagerAction_AllActions(t *testing.T) {
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
action "microsoft365_graph_beta_device_management_managed_device_trigger_configuration_manager_action" "all_actions" {
  config {
    managed_devices = [
      {
        device_id = "00000000-0000-0000-0000-000000000001"
        action    = "refreshMachinePolicy"
      },
      {
        device_id = "00000000-0000-0000-0000-000000000002"
        action    = "refreshUserPolicy"
      },
      {
        device_id = "00000000-0000-0000-0000-000000000003"
        action    = "wakeUpClient"
      },
      {
        device_id = "00000000-0000-0000-0000-000000000004"
        action    = "appEvaluation"
      },
      {
        device_id = "00000000-0000-0000-0000-000000000005"
        action    = "quickScan"
      },
      {
        device_id = "00000000-0000-0000-0000-000000000006"
        action    = "fullScan"
      },
      {
        device_id = "00000000-0000-0000-0000-000000000007"
        action    = "windowsDefenderUpdateSignatures"
      }
    ]
  }
}
`,
			},
		},
	})
}

func TestTriggerConfigurationManagerAction_PartialFailures(t *testing.T) {
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
action "microsoft365_graph_beta_device_management_managed_device_trigger_configuration_manager_action" "partial_failures" {
  config {
    managed_devices = [
      {
        device_id = "00000000-0000-0000-0000-000000000001"
        action    = "refreshMachinePolicy"
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

func TestTriggerConfigurationManagerAction_ValidationDisabled(t *testing.T) {
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
action "microsoft365_graph_beta_device_management_managed_device_trigger_configuration_manager_action" "validation_disabled" {
  config {
    managed_devices = [
      {
        device_id = "00000000-0000-0000-0000-000000000001"
        action    = "refreshMachinePolicy"
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

func TestTriggerConfigurationManagerAction_ValidationEnabled(t *testing.T) {
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
action "microsoft365_graph_beta_device_management_managed_device_trigger_configuration_manager_action" "validation_enabled" {
  config {
    managed_devices = [
      {
        device_id = "00000000-0000-0000-0000-000000000001"
        action    = "refreshMachinePolicy"
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

func TestTriggerConfigurationManagerAction_InvalidDeviceID(t *testing.T) {
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
action "microsoft365_graph_beta_device_management_managed_device_trigger_configuration_manager_action" "invalid_device_id" {
  config {
    managed_devices = [
      {
        device_id = "not-a-valid-guid"
        action    = "refreshMachinePolicy"
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

func TestTriggerConfigurationManagerAction_InvalidAction(t *testing.T) {
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
action "microsoft365_graph_beta_device_management_managed_device_trigger_configuration_manager_action" "invalid_action" {
  config {
    managed_devices = [
      {
        device_id = "00000000-0000-0000-0000-000000000001"
        action    = "invalidAction"
      }
    ]
  }
}
`,
				ExpectError: regexp.MustCompile(`Attribute managed_devices\[0\]\.action value must be one of`),
			},
		},
	})
}

func TestTriggerConfigurationManagerAction_NoDevices(t *testing.T) {
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
action "microsoft365_graph_beta_device_management_managed_device_trigger_configuration_manager_action" "no_devices" {
  config {}
}
`,
				ExpectError: regexp.MustCompile(`No Devices Specified`),
			},
		},
	})
}

func TestTriggerConfigurationManagerAction_DuplicateDeviceIDs(t *testing.T) {
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
action "microsoft365_graph_beta_device_management_managed_device_trigger_configuration_manager_action" "duplicate_device_ids" {
  config {
    managed_devices = [
      {
        device_id = "00000000-0000-0000-0000-000000000001"
        action    = "refreshMachinePolicy"
      },
      {
        device_id = "00000000-0000-0000-0000-000000000001"
        action    = "appEvaluation"
      }
    ]
  }
}
`,
			},
		},
	})
}
