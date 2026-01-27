package graphBetaAndroidManagedDeviceAppConfigurationPolicy_test

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

// Helper function to load test configs from acceptance directory
func loadAcceptanceTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/acceptance/" + filename)
	if err != nil {
		panic("failed to load acceptance config " + filename + ": " + err.Error())
	}
	return acceptance.ConfiguredM365ProviderBlock(config)
}

// TestAccAndroidManagedDeviceAppConfigurationPolicyResource_Lifecycle tests full lifecycle of the resource
func TestAccResourceAndroidManagedDeviceAppConfigurationPolicy_01_Lifecycle(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckAndroidManagedDeviceAppConfigurationPolicyDestroy,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			// Create minimal configuration
			{
				Config: loadAcceptanceTestTerraform("resource_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.minimal", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.minimal", "display_name", "acc-test-android-managed-device-app-configuration-policy-minimal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.minimal", "description", "Acceptance test Android managed store app configuration"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.minimal", "package_id", "app:com.microsoft.office.officehubrow"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.minimal", "profile_applicability", "androidDeviceOwner"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.minimal", "connected_apps_enabled", "true"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.minimal", "payload_json"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.minimal", "role_scope_tag_ids.#", "1"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.minimal", "role_scope_tag_ids.*", "0"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.minimal", "version"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.minimal", "targeted_mobile_apps.#", "1"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.minimal", "app_supports_oem_config", "false"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.minimal",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// TestAccAndroidManagedDeviceAppConfigurationPolicyResource_MicrosoftAuthenticator tests Microsoft Authenticator configuration
func TestAccResourceAndroidManagedDeviceAppConfigurationPolicy_02_MicrosoftAuthenticator(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckAndroidManagedDeviceAppConfigurationPolicyDestroy,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform("resource_microsoft_authenticator_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_authenticator_maximal", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_authenticator_maximal", "package_id", "app:com.azure.authenticator"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_authenticator_maximal", "profile_applicability", "androidDeviceOwner"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_authenticator_maximal", "connected_apps_enabled", "true"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_authenticator_maximal", "payload_json"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_authenticator_maximal", "permission_actions.#", "33"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_authenticator_maximal", "version"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_authenticator_maximal", "app_supports_oem_config", "false"),
				),
			},
		},
	})
}

// TestAccAndroidManagedDeviceAppConfigurationPolicyResource_Microsoft365Copilot tests Microsoft 365 Copilot configuration
func TestAccResourceAndroidManagedDeviceAppConfigurationPolicy_03_Microsoft365Copilot(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckAndroidManagedDeviceAppConfigurationPolicyDestroy,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform("resource_microsoft_365_copilot_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_365_copilot_maximal", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_365_copilot_maximal", "package_id", "app:com.microsoft.office.officehubrow"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_365_copilot_maximal", "profile_applicability", "androidWorkProfile"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_365_copilot_maximal", "connected_apps_enabled", "true"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_365_copilot_maximal", "payload_json"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_365_copilot_maximal", "permission_actions.#", "33"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_365_copilot_maximal", "version"),
				),
			},
		},
	})
}

// TestAccAndroidManagedDeviceAppConfigurationPolicyResource_ManagedHomeScreen tests Managed Home Screen kiosk configuration
func TestAccResourceAndroidManagedDeviceAppConfigurationPolicy_04_ManagedHomeScreen(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckAndroidManagedDeviceAppConfigurationPolicyDestroy,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform("resource_managed_home_screen_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.managed_home_screen_maximal", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.managed_home_screen_maximal", "package_id", "app:com.microsoft.launcher.enterprise"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.managed_home_screen_maximal", "profile_applicability", "default"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.managed_home_screen_maximal", "connected_apps_enabled", "true"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.managed_home_screen_maximal", "payload_json"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.managed_home_screen_maximal", "permission_actions.#", "33"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.managed_home_screen_maximal", "version"),
				),
			},
		},
	})
}

// TestAccAndroidManagedDeviceAppConfigurationPolicyResource_MicrosoftDefender tests Microsoft Defender configuration
func TestAccResourceAndroidManagedDeviceAppConfigurationPolicy_05_MicrosoftDefender(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckAndroidManagedDeviceAppConfigurationPolicyDestroy,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform("resource_microsoft_defender_antivirus_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_defender_antivirus_maximal", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_defender_antivirus_maximal", "package_id", "app:com.microsoft.scmx"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_defender_antivirus_maximal", "profile_applicability", "default"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_defender_antivirus_maximal", "connected_apps_enabled", "true"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_defender_antivirus_maximal", "payload_json"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_defender_antivirus_maximal", "permission_actions.#", "33"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_defender_antivirus_maximal", "version"),
				),
			},
		},
	})
}

// TestAccAndroidManagedDeviceAppConfigurationPolicyResource_MicrosoftEdge tests Microsoft Edge browser configuration
func TestAccResourceAndroidManagedDeviceAppConfigurationPolicy_06_MicrosoftEdge(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckAndroidManagedDeviceAppConfigurationPolicyDestroy,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform("resource_microsoft_edge_browser_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_edge_browser_maximal", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_edge_browser_maximal", "package_id", "app:com.microsoft.emmx"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_edge_browser_maximal", "profile_applicability", "androidWorkProfile"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_edge_browser_maximal", "connected_apps_enabled", "true"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_edge_browser_maximal", "payload_json"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_edge_browser_maximal", "permission_actions.#", "33"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_edge_browser_maximal", "version"),
				),
			},
		},
	})
}

// TestAccAndroidManagedDeviceAppConfigurationPolicyResource_MicrosoftExcel tests Microsoft Excel configuration
func TestAccResourceAndroidManagedDeviceAppConfigurationPolicy_07_MicrosoftExcel(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckAndroidManagedDeviceAppConfigurationPolicyDestroy,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform("resource_microsoft_excel_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_excel_maximal", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_excel_maximal", "package_id", "app:com.microsoft.office.excel"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_excel_maximal", "profile_applicability", "androidDeviceOwner"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_excel_maximal", "connected_apps_enabled", "true"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_excel_maximal", "payload_json"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_excel_maximal", "permission_actions.#", "33"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_excel_maximal", "version"),
				),
			},
		},
	})
}

// TestAccAndroidManagedDeviceAppConfigurationPolicyResource_MicrosoftPowerPoint tests Microsoft PowerPoint configuration
func TestAccResourceAndroidManagedDeviceAppConfigurationPolicy_08_MicrosoftPowerPoint(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckAndroidManagedDeviceAppConfigurationPolicyDestroy,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform("resource_microsoft_powerpoint_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_powerpoint_maximal", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_powerpoint_maximal", "package_id", "app:com.microsoft.office.powerpoint"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_powerpoint_maximal", "profile_applicability", "androidWorkProfile"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_powerpoint_maximal", "connected_apps_enabled", "true"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_powerpoint_maximal", "payload_json"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_powerpoint_maximal", "permission_actions.#", "33"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_powerpoint_maximal", "version"),
				),
			},
		},
	})
}

// TestAccAndroidManagedDeviceAppConfigurationPolicyResource_MicrosoftWord tests Microsoft Word configuration
func TestAccResourceAndroidManagedDeviceAppConfigurationPolicy_09_MicrosoftWord(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckAndroidManagedDeviceAppConfigurationPolicyDestroy,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform("resource_microsoft_word_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_word_maximal", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_word_maximal", "package_id", "app:com.microsoft.office.word"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_word_maximal", "profile_applicability", "androidDeviceOwner"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_word_maximal", "connected_apps_enabled", "true"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_word_maximal", "payload_json"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_word_maximal", "permission_actions.#", "33"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_word_maximal", "version"),
				),
			},
		},
	})
}

// TestAccAndroidManagedDeviceAppConfigurationPolicyResource_MicrosoftOneNote tests Microsoft OneNote configuration
func TestAccResourceAndroidManagedDeviceAppConfigurationPolicy_10_MicrosoftOneNote(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckAndroidManagedDeviceAppConfigurationPolicyDestroy,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform("resource_microsoft_onenote_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_onenote_maximal", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_onenote_maximal", "package_id", "app:com.microsoft.office.onenote"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_onenote_maximal", "profile_applicability", "androidWorkProfile"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_onenote_maximal", "connected_apps_enabled", "true"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_onenote_maximal", "payload_json"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_onenote_maximal", "permission_actions.#", "33"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_onenote_maximal", "version"),
				),
			},
		},
	})
}

// TestAccAndroidManagedDeviceAppConfigurationPolicyResource_MicrosoftOneDrive tests Microsoft OneDrive configuration
func TestAccResourceAndroidManagedDeviceAppConfigurationPolicy_11_MicrosoftOneDrive(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckAndroidManagedDeviceAppConfigurationPolicyDestroy,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform("resource_microsoft_onedrive_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_onedrive_maximal", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_onedrive_maximal", "package_id", "app:com.microsoft.skydrive"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_onedrive_maximal", "profile_applicability", "androidDeviceOwner"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_onedrive_maximal", "connected_apps_enabled", "true"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_onedrive_maximal", "payload_json"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_onedrive_maximal", "permission_actions.#", "33"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_onedrive_maximal", "version"),
				),
			},
		},
	})
}

// TestAccAndroidManagedDeviceAppConfigurationPolicyResource_MicrosoftOutlook tests Outlook email configuration
func TestAccResourceAndroidManagedDeviceAppConfigurationPolicy_12_MicrosoftOutlook(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckAndroidManagedDeviceAppConfigurationPolicyDestroy,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform("resource_microsoft_outlook_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_outlook_maximal", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_outlook_maximal", "package_id", "app:com.microsoft.office.outlook"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_outlook_maximal", "profile_applicability", "androidDeviceOwner"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_outlook_maximal", "connected_apps_enabled", "true"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_outlook_maximal", "payload_json"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_outlook_maximal", "permission_actions.#", "33"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_outlook_maximal", "version"),
				),
			},
		},
	})
}

// TestAccAndroidManagedDeviceAppConfigurationPolicyResource_MicrosoftTeams tests Teams collaboration configuration
func TestAccResourceAndroidManagedDeviceAppConfigurationPolicy_13_MicrosoftTeams(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckAndroidManagedDeviceAppConfigurationPolicyDestroy,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform("resource_microsoft_teams_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_teams_maximal", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_teams_maximal", "package_id", "app:com.microsoft.teams"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_teams_maximal", "profile_applicability", "androidWorkProfile"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_teams_maximal", "connected_apps_enabled", "true"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_teams_maximal", "payload_json"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_teams_maximal", "permission_actions.#", "33"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_teams_maximal", "version"),
				),
			},
		},
	})
}

// testAccCheckAndroidManagedDeviceAppConfigurationPolicyDestroy verifies the resource has been destroyed
func testAccCheckAndroidManagedDeviceAppConfigurationPolicyDestroy(s *terraform.State) error {
	graphClient, err := acceptance.TestGraphClient()
	if err != nil {
		return fmt.Errorf("error creating Graph client for CheckDestroy: %v", err)
	}

	ctx := context.Background()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy" {
			continue
		}

		_, err := graphClient.
			DeviceAppManagement().
			MobileAppConfigurations().
			ByManagedDeviceMobileAppConfigurationId(rs.Primary.ID).
			Get(ctx, nil)

		if err != nil {
			errorInfo := errors.GraphError(ctx, err)

			if errorInfo.StatusCode == 404 ||
				errorInfo.ErrorCode == "ResourceNotFound" ||
				errorInfo.ErrorCode == "ItemNotFound" {
				continue // Resource successfully destroyed
			}
			return fmt.Errorf("error checking if Android managed device app configuration policy %s was destroyed: %v", rs.Primary.ID, err)
		}

		return fmt.Errorf("Android managed device app configuration policy %s still exists", rs.Primary.ID)
	}

	return nil
}
