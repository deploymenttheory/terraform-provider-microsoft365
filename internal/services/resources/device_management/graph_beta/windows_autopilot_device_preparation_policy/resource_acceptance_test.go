package graphBetaWindowsAutopilotDevicePreparationPolicy_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccWindowsAutopilotDevicePreparationPolicyResource_Minimal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckWindowsAutopilotDevicePreparationPolicyDestroy,
		ExternalProviders: map[string]resource.ExternalProvider{

			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy.minimal", "name", "acc-test-windows-autopilot-device-preparation-policy-minimal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy.minimal", "description", "acc-test-windows-autopilot-device-preparation-policy-minimal"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy.minimal", "device_security_group"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy.minimal", "deployment_settings.deployment_mode", "enrollment_autopilot_dpp_deploymentmode_0"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy.minimal", "deployment_settings.deployment_type", "enrollment_autopilot_dpp_deploymenttype_0"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy.minimal", "deployment_settings.join_type", "enrollment_autopilot_dpp_jointype_0"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy.minimal", "deployment_settings.account_type", "enrollment_autopilot_dpp_accountype_0"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy.minimal", "oobe_settings.timeout_in_minutes", "60"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy.minimal", "oobe_settings.allow_skip", "false"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy.minimal", "oobe_settings.allow_diagnostics", "false"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy.minimal", "assignments.include_group_ids.#", "2"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy.minimal", "id"),
				),
			},
			{
				ResourceName:                         "microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy.minimal",
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "id",
				ImportStateVerifyIgnore: []string{
					"timeouts",
				},
			},
		},
	})
}

func TestAccWindowsAutopilotDevicePreparationPolicyResource_Enhanced(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckWindowsAutopilotDevicePreparationPolicyDestroy,
		ExternalProviders: map[string]resource.ExternalProvider{

			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccConfigEnhanced(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy.enhanced", "name", "acc-test-windows-autopilot-device-preparation-policy-enhanced"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy.enhanced", "deployment_settings.deployment_mode", "enrollment_autopilot_dpp_deploymentmode_1"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy.enhanced", "deployment_settings.deployment_type", "enrollment_autopilot_dpp_deploymenttype_1"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy.enhanced", "deployment_settings.account_type", "enrollment_autopilot_dpp_accountype_1"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy.enhanced", "oobe_settings.timeout_in_minutes", "120"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy.enhanced", "oobe_settings.allow_skip", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy.enhanced", "oobe_settings.allow_diagnostics", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy.enhanced", "allowed_apps.#", "1"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy.enhanced", "allowed_scripts.#", "1"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy.enhanced", "assignments.include_group_ids.#", "3"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy.enhanced", "id"),
				),
			},
		},
	})
}

func TestAccWindowsAutopilotDevicePreparationPolicyResource_SelfDeploying(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckWindowsAutopilotDevicePreparationPolicyDestroy,
		ExternalProviders: map[string]resource.ExternalProvider{

			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccConfigSelfDeploying(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy.self_deploying", "name", "acc-test-windows-autopilot-device-preparation-policy-self-deploying"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy.self_deploying", "deployment_settings.deployment_type", "enrollment_autopilot_dpp_deploymenttype_1"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy.self_deploying", "oobe_settings.timeout_in_minutes", "90"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy.self_deploying", "assignments.include_group_ids.#", "1"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy.self_deploying", "id"),
				),
			},
		},
	})
}

func TestAccWindowsAutopilotDevicePreparationPolicyResource_HybridJoined(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckWindowsAutopilotDevicePreparationPolicyDestroy,
		ExternalProviders: map[string]resource.ExternalProvider{

			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccConfigHybridJoined(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy.hybrid_joined", "name", "acc-test-windows-autopilot-device-preparation-policy-hybrid-joined"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy.hybrid_joined", "deployment_settings.join_type", "enrollment_autopilot_dpp_jointype_1"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy.hybrid_joined", "oobe_settings.timeout_in_minutes", "75"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy.hybrid_joined", "assignments.include_group_ids.#", "2"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy.hybrid_joined", "id"),
				),
			},
		},
	})
}

func testAccConfigMinimal() string {
	// Load dependencies
	groupsConfig, err := helpers.ParseHCLFile("../../../../../acceptance/terraform_dependancies/device_management/groups.tf")
	if err != nil {
		panic(fmt.Sprintf("failed to load groups dependency: %s", err.Error()))
	}

	autopilotGroupsConfig, err := helpers.ParseHCLFile("../../../../../acceptance/terraform_dependancies/device_management/autopilot_security_groups.tf")
	if err != nil {
		panic(fmt.Sprintf("failed to load autopilot groups dependency: %s", err.Error()))
	}

	// Load test configuration
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/autopilot_device_preparation_minimal.tf")
	if err != nil {
		panic(fmt.Sprintf("failed to load acceptance test config: %s", err.Error()))
	}

	// Combine configurations
	combinedConfig := groupsConfig + "\n\n" + autopilotGroupsConfig + "\n\n" + accTestConfig
	return acceptance.ConfiguredM365ProviderBlock(combinedConfig)
}

func testAccConfigEnhanced() string {
	// Load dependencies
	groupsConfig, err := helpers.ParseHCLFile("../../../../../acceptance/terraform_dependancies/device_management/groups.tf")
	if err != nil {
		panic(fmt.Sprintf("failed to load groups dependency: %s", err.Error()))
	}

	autopilotGroupsConfig, err := helpers.ParseHCLFile("../../../../../acceptance/terraform_dependancies/device_management/autopilot_security_groups.tf")
	if err != nil {
		panic(fmt.Sprintf("failed to load autopilot groups dependency: %s", err.Error()))
	}

	appConfig, err := helpers.ParseHCLFile("../../../../../acceptance/terraform_dependancies/device_management/win32_lob_app.tf")
	if err != nil {
		panic(fmt.Sprintf("failed to load app dependency: %s", err.Error()))
	}

	scriptConfig, err := helpers.ParseHCLFile("../../../../../acceptance/terraform_dependancies/device_management/device_shell_script.tf")
	if err != nil {
		panic(fmt.Sprintf("failed to load script dependency: %s", err.Error()))
	}

	// Load test configuration
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/autopilot_device_preparation_enhanced.tf")
	if err != nil {
		panic(fmt.Sprintf("failed to load acceptance test config: %s", err.Error()))
	}

	// Combine configurations
	combinedConfig := groupsConfig + "\n\n" + autopilotGroupsConfig + "\n\n" + appConfig + "\n\n" + scriptConfig + "\n\n" + accTestConfig
	return acceptance.ConfiguredM365ProviderBlock(combinedConfig)
}

func testAccConfigSelfDeploying() string {
	// Load dependencies
	groupsConfig, err := helpers.ParseHCLFile("../../../../../acceptance/terraform_dependancies/device_management/groups.tf")
	if err != nil {
		panic(fmt.Sprintf("failed to load groups dependency: %s", err.Error()))
	}

	autopilotGroupsConfig, err := helpers.ParseHCLFile("../../../../../acceptance/terraform_dependancies/device_management/autopilot_security_groups.tf")
	if err != nil {
		panic(fmt.Sprintf("failed to load autopilot groups dependency: %s", err.Error()))
	}

	// Load test configuration
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/autopilot_device_preparation_self_deploying.tf")
	if err != nil {
		panic(fmt.Sprintf("failed to load acceptance test config: %s", err.Error()))
	}

	// Combine configurations
	combinedConfig := groupsConfig + "\n\n" + autopilotGroupsConfig + "\n\n" + accTestConfig
	return acceptance.ConfiguredM365ProviderBlock(combinedConfig)
}

func testAccConfigHybridJoined() string {
	// Load dependencies
	groupsConfig, err := helpers.ParseHCLFile("../../../../../acceptance/terraform_dependancies/device_management/groups.tf")
	if err != nil {
		panic(fmt.Sprintf("failed to load groups dependency: %s", err.Error()))
	}

	autopilotGroupsConfig, err := helpers.ParseHCLFile("../../../../../acceptance/terraform_dependancies/device_management/autopilot_security_groups.tf")
	if err != nil {
		panic(fmt.Sprintf("failed to load autopilot groups dependency: %s", err.Error()))
	}

	// Load test configuration
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/autopilot_device_preparation_hybrid_joined.tf")
	if err != nil {
		panic(fmt.Sprintf("failed to load acceptance test config: %s", err.Error()))
	}

	// Combine configurations
	combinedConfig := groupsConfig + "\n\n" + autopilotGroupsConfig + "\n\n" + accTestConfig
	return acceptance.ConfiguredM365ProviderBlock(combinedConfig)
}

func testAccCheckWindowsAutopilotDevicePreparationPolicyDestroy(s *terraform.State) error {
	graphClient, err := acceptance.TestGraphClient()
	if err != nil {
		return fmt.Errorf("error creating Graph client for CheckDestroy: %v", err)
	}
	ctx := context.Background()
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy" {
			continue
		}
		_, err := graphClient.
			DeviceManagement().
			ConfigurationPolicies().
			ByDeviceManagementConfigurationPolicyId(rs.Primary.ID).
			Get(ctx, nil)

		if err != nil {
			errorInfo := errors.GraphError(ctx, err)
			if errorInfo.StatusCode == 404 || errorInfo.ErrorCode == "ResourceNotFound" || errorInfo.ErrorCode == "ItemNotFound" {
				fmt.Printf("DEBUG: Resource %s successfully destroyed (404/NotFound)\n", rs.Primary.ID)
				continue
			}
			return fmt.Errorf("error checking if windows autopilot device preparation policy %s was destroyed: %v", rs.Primary.ID, err)
		}
		return fmt.Errorf("windows autopilot device preparation policy %s still exists", rs.Primary.ID)
	}
	return nil
}
