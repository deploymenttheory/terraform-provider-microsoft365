# User-Driven Deployment Profile Example
resource "microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile" "user_driven" {
  display_name                                 = "acc test user driven autopilot profile with os default locale"
  description                                  = "user driven autopilot profile with os default locale"
  device_name_template                         = "thing-%RAND:5%" // Apply device name template max 15 characters
  locale                                       = "os-default"
  preprovisioning_allowed                      = true // Allow pre-provisioned deployment
  device_type                                  = "windowsPc"
  hardware_hash_extraction_enabled             = true
  role_scope_tag_ids                           = ["0", "1"]
  device_join_type                             = "microsoft_entra_joined"
  hybrid_azure_ad_join_skip_connectivity_check = false // always false when using microsoft_entra_joined

  out_of_box_experience_setting = {
    device_usage_type               = "singleUser"
    privacy_settings_hidden         = true       // Privacy settings
    eula_hidden                     = true       // Microsoft Software License Terms
    user_type                       = "standard" // standard or administrator
    keyboard_selection_page_skipped = true       // Automatically configure keyboard
  }

  // Optional assignments, can be either group based or all devices based
  assignments = [
    {
      type     = "groupAssignmentTarget"
      group_id = "00000000-0000-0000-0000-000000000001"
    },
    {
      type     = "groupAssignmentTarget"
      group_id = "00000000-0000-0000-0000-000000000002"
    },
    {
      type     = "exclusionGroupAssignmentTarget"
      group_id = "00000000-0000-0000-0000-000000000003"
    }
  ]
}