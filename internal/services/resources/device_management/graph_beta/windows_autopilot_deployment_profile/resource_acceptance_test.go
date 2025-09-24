package graphBetaWindowsAutopilotDeploymentProfile_test

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

func TestAccWindowsAutopilotDeploymentProfileResource_Minimal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckWindowsAutopilotDeploymentProfileDestroy,
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
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.minimal", "display_name", "acc-test-windows-autopilot-deployment-profile-minimal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.minimal", "description", "acc-test-windows-autopilot-deployment-profile-minimal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.minimal", "device_join_type", "microsoft_entra_joined"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.minimal", "locale", "os-default"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.minimal", "hardware_hash_extraction_enabled", "false"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.minimal", "preprovisioning_allowed", "false"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.minimal", "device_name_template", ""),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.minimal", "role_scope_tag_ids.#", "1"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.minimal", "role_scope_tag_ids.0", "0"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.minimal", "out_of_box_experience_setting.privacy_settings_hidden", "false"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.minimal", "out_of_box_experience_setting.eula_hidden", "false"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.minimal", "out_of_box_experience_setting.user_type", "administrator"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.minimal", "out_of_box_experience_setting.device_usage_type", "singleUser"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.minimal", "enrollment_status_screen_settings.hide_installation_progress", "false"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.minimal", "enrollment_status_screen_settings.allow_device_use_before_profile_and_app_install_complete", "false"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.minimal", "id"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.minimal", "created_date_time"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.minimal", "last_modified_date_time"),
				),
			},
			{
				ResourceName:                         "microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.minimal",
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

func TestAccWindowsAutopilotDeploymentProfileResource_Enhanced(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckWindowsAutopilotDeploymentProfileDestroy,
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
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.enhanced", "display_name", "acc-test-windows-autopilot-deployment-profile-enhanced"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.enhanced", "description", "acc-test-windows-autopilot-deployment-profile-enhanced"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.enhanced", "device_join_type", "microsoft_entra_joined"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.enhanced", "locale", "en-US"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.enhanced", "hardware_hash_extraction_enabled", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.enhanced", "preprovisioning_allowed", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.enhanced", "device_name_template", "TEST-%RAND:3%"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.enhanced", "device_type", "windowsPc"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.enhanced", "role_scope_tag_ids.#", "1"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.enhanced", "role_scope_tag_ids.0", "0"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.enhanced", "out_of_box_experience_setting.privacy_settings_hidden", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.enhanced", "out_of_box_experience_setting.eula_hidden", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.enhanced", "out_of_box_experience_setting.user_type", "standard"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.enhanced", "out_of_box_experience_setting.device_usage_type", "shared"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.enhanced", "enrollment_status_screen_settings.hide_installation_progress", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.enhanced", "enrollment_status_screen_settings.allow_device_use_before_profile_and_app_install_complete", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.enhanced", "enrollment_status_screen_settings.custom_error_message", "Please contact IT support for assistance"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.enhanced", "enrollment_status_screen_settings.install_progress_timeout_in_minutes", "120"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.enhanced", "id"),
				),
			},
		},
	})
}

func TestAccWindowsAutopilotDeploymentProfileResource_HybridJoined(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckWindowsAutopilotDeploymentProfileDestroy,
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
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.hybrid_joined", "display_name", "acc-test-windows-autopilot-deployment-profile-hybrid-joined"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.hybrid_joined", "device_join_type", "microsoft_entra_hybrid_joined"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.hybrid_joined", "hybrid_azure_ad_join_skip_connectivity_check", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.hybrid_joined", "out_of_box_experience_setting.user_type", "administrator"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.hybrid_joined", "enrollment_status_screen_settings.install_progress_timeout_in_minutes", "90"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.hybrid_joined", "id"),
				),
			},
		},
	})
}

func testAccConfigMinimal() string {
	// Load test configuration
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/autopilot_deployment_profile_minimal.tf")
	if err != nil {
		panic(fmt.Sprintf("failed to load acceptance test config: %s", err.Error()))
	}

	// Return configured provider block with test config
	return acceptance.ConfiguredM365ProviderBlock(accTestConfig)
}

func testAccConfigEnhanced() string {
	// Load test configuration
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/autopilot_deployment_profile_enhanced.tf")
	if err != nil {
		panic(fmt.Sprintf("failed to load acceptance test config: %s", err.Error()))
	}

	// Return configured provider block with test config
	return acceptance.ConfiguredM365ProviderBlock(accTestConfig)
}

func testAccConfigHybridJoined() string {
	// Load test configuration
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/autopilot_deployment_profile_hybrid_joined.tf")
	if err != nil {
		panic(fmt.Sprintf("failed to load acceptance test config: %s", err.Error()))
	}

	// Return configured provider block with test config
	return acceptance.ConfiguredM365ProviderBlock(accTestConfig)
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