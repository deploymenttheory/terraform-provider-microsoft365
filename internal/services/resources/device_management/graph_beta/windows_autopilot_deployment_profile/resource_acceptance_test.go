package graphBetaWindowsAutopilotDeploymentProfile_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccResourceWindowsAutopilotDeploymentProfile_01_SelfDeployingOSDefaultLocale(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckWindowsAutopilotDeploymentProfileDestroy,
		ExternalProviders: map[string]resource.ExternalProvider{

			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccConfigSelfDeployingOSDefaultLocale(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven", "display_name", "acc test user driven autopilot profile with os default locale"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven", "description", "user driven autopilot profile with os default locale"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven", "locale", "os-default"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven", "preprovisioning_allowed", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven", "device_type", "windowsPc"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven", "hardware_hash_extraction_enabled", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven", "device_join_type", "microsoft_entra_joined"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven", "hybrid_azure_ad_join_skip_connectivity_check", "false"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven", "out_of_box_experience_setting.device_usage_type", "singleUser"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven", "out_of_box_experience_setting.privacy_settings_hidden", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven", "out_of_box_experience_setting.eula_hidden", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven", "out_of_box_experience_setting.user_type", "standard"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven", "out_of_box_experience_setting.keyboard_selection_page_skipped", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven", "assignments.#", "3"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven", "id"),
				),
			},
			{
				ResourceName:                         "microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven",
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "id",
				ImportStateVerifyIgnore: []string{
					"timeouts",
					"hybrid_azure_ad_join_skip_connectivity_check",
				},
			},
		},
	})
}

func TestAccResourceWindowsAutopilotDeploymentProfile_02_UserDrivenHybridDomainJoin(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckWindowsAutopilotDeploymentProfileDestroy,
		ExternalProviders: map[string]resource.ExternalProvider{

			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccConfigUserDrivenHybridDomainJoin(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven_japanese_preprovisioned_with_assignments", "display_name", "acc_test_user_driven_japanese_preprovisioned"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven_japanese_preprovisioned_with_assignments", "description", "user driven autopilot profile with japanese locale and allow pre provisioned deployment"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven_japanese_preprovisioned_with_assignments", "locale", "ja-JP"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven_japanese_preprovisioned_with_assignments", "preprovisioning_allowed", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven_japanese_preprovisioned_with_assignments", "device_type", "windowsPc"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven_japanese_preprovisioned_with_assignments", "device_join_type", "microsoft_entra_hybrid_joined"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven_japanese_preprovisioned_with_assignments", "hybrid_azure_ad_join_skip_connectivity_check", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven_japanese_preprovisioned_with_assignments", "assignments.#", "1"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven_japanese_preprovisioned_with_assignments", "id"),
				),
			},
		},
	})
}

func TestAccResourceWindowsAutopilotDeploymentProfile_02_UserDrivenWithGroupAssignments(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckWindowsAutopilotDeploymentProfileDestroy,
		ExternalProviders: map[string]resource.ExternalProvider{

			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccConfigUserDrivenWithGroupAssignments(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven", "display_name", "acc test user driven autopilot with group assignments"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven", "description", "user driven autopilot profile with os default locale"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven", "locale", "os-default"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven", "preprovisioning_allowed", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven", "device_type", "windowsPc"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven", "device_join_type", "microsoft_entra_joined"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven", "assignments.#", "3"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.user_driven", "id"),
				),
			},
		},
	})
}

func TestAccResourceWindowsAutopilotDeploymentProfile_04_HololensWithAllDeviceAssignment(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckWindowsAutopilotDeploymentProfileDestroy,
		ExternalProviders: map[string]resource.ExternalProvider{

			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccConfigHololensWithAllDeviceAssignment(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.hololens_with_all_device_assignment", "display_name", "acc_test_hololens_with_all_device_assignment"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.hololens_with_all_device_assignment", "description", "hololens autopilot profile with hk locale and group assignment"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.hololens_with_all_device_assignment", "locale", "zh-HK"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.hololens_with_all_device_assignment", "preprovisioning_allowed", "false"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.hololens_with_all_device_assignment", "device_type", "holoLens"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.hololens_with_all_device_assignment", "hardware_hash_extraction_enabled", "false"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.hololens_with_all_device_assignment", "device_join_type", "microsoft_entra_joined"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.hololens_with_all_device_assignment", "out_of_box_experience_setting.device_usage_type", "shared"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.hololens_with_all_device_assignment", "out_of_box_experience_setting.privacy_settings_hidden", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.hololens_with_all_device_assignment", "out_of_box_experience_setting.eula_hidden", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.hololens_with_all_device_assignment", "out_of_box_experience_setting.user_type", "standard"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.hololens_with_all_device_assignment", "out_of_box_experience_setting.keyboard_selection_page_skipped", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.hololens_with_all_device_assignment", "assignments.#", "1"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.hololens_with_all_device_assignment", "id"),
				),
			},
		},
	})
}

func testAccConfigSelfDeployingOSDefaultLocale() string {
	// Load dependencies
	groupsConfig, err := helpers.ParseHCLFile("../../../../../acceptance/terraform_dependancies/device_management/groups.tf")
	if err != nil {
		panic(fmt.Sprintf("failed to load groups dependency: %s", err.Error()))
	}

	// Load test configuration
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/01_self_deploying_os_default_locale.tf")
	if err != nil {
		panic(fmt.Sprintf("failed to load acceptance test config: %s", err.Error()))
	}

	// Combine configurations
	combinedConfig := groupsConfig + "\n\n" + accTestConfig
	return acceptance.ConfiguredM365ProviderBlock(combinedConfig)
}

func testAccConfigUserDrivenHybridDomainJoin() string {
	// Load dependencies
	groupsConfig, err := helpers.ParseHCLFile("../../../../../acceptance/terraform_dependancies/device_management/groups.tf")
	if err != nil {
		panic(fmt.Sprintf("failed to load groups dependency: %s", err.Error()))
	}

	// Load test configuration
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/02_user_driven_hybrid_domain_join.tf")
	if err != nil {
		panic(fmt.Sprintf("failed to load acceptance test config: %s", err.Error()))
	}

	// Combine configurations
	combinedConfig := groupsConfig + "\n\n" + accTestConfig
	return acceptance.ConfiguredM365ProviderBlock(combinedConfig)
}

func testAccConfigUserDrivenWithGroupAssignments() string {
	// Load dependencies
	groupsConfig, err := helpers.ParseHCLFile("../../../../../acceptance/terraform_dependancies/device_management/groups.tf")
	if err != nil {
		panic(fmt.Sprintf("failed to load groups dependency: %s", err.Error()))
	}

	// Load test configuration
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/03_user_driven_with_group_assignments.tf")
	if err != nil {
		panic(fmt.Sprintf("failed to load acceptance test config: %s", err.Error()))
	}

	// Combine configurations
	combinedConfig := groupsConfig + "\n\n" + accTestConfig
	return acceptance.ConfiguredM365ProviderBlock(combinedConfig)
}

func testAccConfigHololensWithAllDeviceAssignment() string {
	// Load dependencies
	groupsConfig, err := helpers.ParseHCLFile("../../../../../acceptance/terraform_dependancies/device_management/groups.tf")
	if err != nil {
		panic(fmt.Sprintf("failed to load groups dependency: %s", err.Error()))
	}

	// Load test configuration
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/04_hololens_with_all_device_assignment.tf")
	if err != nil {
		panic(fmt.Sprintf("failed to load acceptance test config: %s", err.Error()))
	}

	// Combine configurations
	combinedConfig := groupsConfig + "\n\n" + accTestConfig
	return acceptance.ConfiguredM365ProviderBlock(combinedConfig)
}

func testAccCheckWindowsAutopilotDeploymentProfileDestroy(s *terraform.State) error {
	graphClient, err := acceptance.TestGraphClient()
	if err != nil {
		return fmt.Errorf("error creating Graph client for CheckDestroy: %v", err)
	}
	ctx := context.Background()
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile" {
			continue
		}
		_, err := graphClient.
			DeviceManagement().
			WindowsAutopilotDeploymentProfiles().
			ByWindowsAutopilotDeploymentProfileId(rs.Primary.ID).
			Get(ctx, nil)

		if err != nil {
			errorInfo := errors.GraphError(ctx, err)
			if errorInfo.StatusCode == 404 || errorInfo.ErrorCode == "ResourceNotFound" || errorInfo.ErrorCode == "ItemNotFound" {
				fmt.Printf("DEBUG: Resource %s successfully destroyed (404/NotFound)\n", rs.Primary.ID)
				continue
			}
			return fmt.Errorf("error checking if windows autopilot deployment profile %s was destroyed: %v", rs.Primary.ID, err)
		}
		return fmt.Errorf("windows autopilot deployment profile %s still exists", rs.Primary.ID)
	}
	return nil
}
