package graphBetaCloudPcProvisioningPolicy_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceCloudPcProvisioningPolicy_01_Complete(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create with minimal configuration
			{
				Config: testAccCloudPcProvisioningPolicyConfig_minimal(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_windows_365_cloud_pc_provisioning_policy.test", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_provisioning_policy.test", "display_name", "Test Acceptance Provisioning Policy"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_provisioning_policy.test", "image_id", "microsoftwindowsdesktop_windows-ent-cpc_win11-23h2-ent-cpc"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_provisioning_policy.test", "enable_single_sign_on", "false"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_provisioning_policy.test", "local_admin_enabled", "false"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_provisioning_policy.test", "windows_setting.locale", "en-US"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_provisioning_policy.test", "microsoft_managed_desktop.managed_type", "notManaged"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "microsoft365_graph_beta_windows_365_cloud_pc_provisioning_policy.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update to maximal configuration
			{
				Config: testAccCloudPcProvisioningPolicyConfig_maximal(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_windows_365_cloud_pc_provisioning_policy.test", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_provisioning_policy.test", "display_name", "Test Acceptance Provisioning Policy - Updated"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_provisioning_policy.test", "description", "Updated description for acceptance testing"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_provisioning_policy.test", "enable_single_sign_on", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_provisioning_policy.test", "local_admin_enabled", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_provisioning_policy.test", "cloud_pc_naming_template", "CPC-ACC-%USERNAME:5%-%RAND:5%"),
				),
			},
			// Update back to minimal configuration
			{
				Config: testAccCloudPcProvisioningPolicyConfig_minimal(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_windows_365_cloud_pc_provisioning_policy.test", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_provisioning_policy.test", "display_name", "Test Acceptance Provisioning Policy"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_provisioning_policy.test", "enable_single_sign_on", "false"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_provisioning_policy.test", "local_admin_enabled", "false"),
				),
			},
		},
	})
}

func TestAccResourceCloudPcProvisioningPolicy_02_WithAssignments(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create with assignments
			{
				Config: testAccCloudPcProvisioningPolicyConfig_withAssignments(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_windows_365_cloud_pc_provisioning_policy.test_assignments", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_provisioning_policy.test_assignments", "display_name", "Test Provisioning Policy with Assignments"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_provisioning_policy.test_assignments", "assignments.#", "1"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_provisioning_policy.test_assignments", "assignments.0.type", "groupAssignmentTarget"),
				),
			},
		},
	})
}

func TestAccResourceCloudPcProvisioningPolicy_03_RequiredFields(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccCloudPcProvisioningPolicyConfig_missingDisplayName(),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
			{
				Config:      testAccCloudPcProvisioningPolicyConfig_missingImageId(),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
		},
	})
}

func testAccCloudPcProvisioningPolicyConfig_minimal() string {
	return `
resource "microsoft365_graph_beta_windows_365_cloud_pc_provisioning_policy" "test" {
  display_name = "Test Acceptance Provisioning Policy"
  image_id     = "microsoftwindowsdesktop_windows-ent-cpc_win11-23h2-ent-cpc"
  
  microsoft_managed_desktop = {
    # Uses default values: managed_type = "notManaged", profile = "4aa9b805-9494-4eed-a04b-ed51ec9e631e"
  }
  
  windows_setting = {
    locale = "en-US"
  }
}
`
}

func testAccCloudPcProvisioningPolicyConfig_maximal() string {
	return `
resource "microsoft365_graph_beta_windows_365_cloud_pc_provisioning_policy" "test" {
  display_name             = "Test Acceptance Provisioning Policy - Updated"
  description              = "Updated description for acceptance testing"
  cloud_pc_naming_template = "CPC-ACC-%USERNAME:5%-%RAND:5%"
  provisioning_type        = "dedicated"
  image_id                 = "microsoftwindowsdesktop_windows-ent-cpc_win11-24H2-ent-cpc-m365"
  image_type               = "gallery"
  enable_single_sign_on    = true
  local_admin_enabled      = true
  managed_by               = "windows365"

  windows_setting = {
    locale = "en-US"
  }

  microsoft_managed_desktop = {
    managed_type = "notManaged"
    profile      = "4aa9b805-9494-4eed-a04b-ed51ec9e631e"
  }

  domain_join_configurations = [
    {
      domain_join_type          = "hybridAzureADJoin"
      on_premises_connection_id = "33333333-3333-3333-3333-333333333333"
      region_name               = "automatic"
      region_group              = "usWest"
    }
  ]

  autopatch = {
    autopatch_group_id = "4aa9b805-9494-4eed-a04b-ed51ec9e631e"
  }

  autopilot_configuration = {
    device_preparation_profile_id   = "12345678-1234-1234-1234-123456789012"
    application_timeout_in_minutes  = 60
    on_failure_device_access_denied = true
  }

  apply_to_existing_cloud_pcs = {
    microsoft_entra_single_sign_on_for_all_devices        = false
    region_or_azure_network_connection_for_all_devices    = true
    region_or_azure_network_connection_for_select_devices = false
  }

  scope_ids = ["9", "8"]
}
`
}

func testAccCloudPcProvisioningPolicyConfig_withAssignments() string {
	return `
data "azuread_group" "test_group" {
  display_name = "Test Group"
}

resource "microsoft365_graph_beta_windows_365_cloud_pc_provisioning_policy" "test_assignments" {
  display_name = "Test Provisioning Policy with Assignments"
  image_id     = "microsoftwindowsdesktop_windows-ent-cpc_win11-23h2-ent-cpc"

  microsoft_managed_desktop = {
    # Uses default values
  }

  windows_setting = {
    locale = "en-US"
  }

  assignments = [
    {
      type     = "groupAssignmentTarget"
      group_id = data.azuread_group.test_group.object_id
    }
  ]
}
`
}

func testAccCloudPcProvisioningPolicyConfig_missingDisplayName() string {
	return `
resource "microsoft365_graph_beta_windows_365_cloud_pc_provisioning_policy" "test" {
  image_id = "microsoftwindowsdesktop_windows-ent-cpc_win11-23h2-ent-cpc"
  
  microsoft_managed_desktop = {}
  
  windows_setting = {
    locale = "en-US"
  }
}
`
}

func testAccCloudPcProvisioningPolicyConfig_missingImageId() string {
	return `
resource "microsoft365_graph_beta_windows_365_cloud_pc_provisioning_policy" "test" {
  display_name = "Test Policy"
  
  microsoft_managed_desktop = {}
  
  windows_setting = {
    locale = "en-US"
  }
}
`
}
