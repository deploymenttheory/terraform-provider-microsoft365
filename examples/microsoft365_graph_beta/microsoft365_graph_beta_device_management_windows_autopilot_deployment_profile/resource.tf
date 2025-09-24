# User-Driven Deployment Profile Example
resource "microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile" "user_driven" {
  display_name                                 = "acc_test_user_driven"
  description                                  = "user driven autopilot profile with os default locale"
  device_name_template                         = "thing-%RAND:5%" // Apply device name template
  locale                                       = "os-default"
  preprovisioning_allowed                      = false // Allow pre-provisioned deployment
  device_type                                  = "windowsPc"
  hardware_hash_extraction_enabled             = true
  role_scope_tag_ids                           = ["0"]
  device_join_type                             = "microsoft_entra_joined"
  hybrid_azure_ad_join_skip_connectivity_check = false

  out_of_box_experience_setting = {
    device_usage_type               = "singleUser"
    privacy_settings_hidden         = true // Privacy settings
    eula_hidden                     = true // Microsoft Software License Terms
    user_type                       = "standard"
    keyboard_selection_page_skipped = true // Automatically configure keyboard
  }

 // Optional assignments
    assignments = [
    {
      type     = "groupAssignmentTarget"
      group_id = microsoft365_graph_beta_groups_group.acc_test_group_1.id
    },
    {
      type     = "groupAssignmentTarget"
      group_id = microsoft365_graph_beta_groups_group.acc_test_group_2.id
    },
    {
      type     = "exclusionGroupAssignmentTarget"
      group_id = microsoft365_graph_beta_groups_group.acc_test_group_3.id
    }
  ]
}

# User-Driven with Japanese Language and Allow Pre-provisioned Deployment
resource "microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile" "user_driven_japanese_preprovisioned_with_assignments" {
  display_name                                 = "acc_test_user_driven_japanese_preprovisioned"
  description                                  = "user driven autopilot profile with japanese locale and allow pre provisioned deployment"
  device_name_template                         = "thing-%RAND:3%"
  locale                                       = "ja-JP"
  preprovisioning_allowed                      = true
  device_type                                  = "windowsPc"
  hardware_hash_extraction_enabled             = true
  role_scope_tag_ids                           = ["0"]
  device_join_type                             = "microsoft_entra_joined"
  hybrid_azure_ad_join_skip_connectivity_check = false

  out_of_box_experience_setting = {
    device_usage_type               = "singleUser"
    privacy_settings_hidden         = true
    eula_hidden                     = true
    user_type                       = "standard"
    keyboard_selection_page_skipped = true
  }

    assignments = [
    {
      type     = "allDevicesAssignmentTarget"
    }
  ]
}

# Self-Deploying Deployment Profile Example
resource "microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile" "self_deploying" {
  display_name                                 = "acc_test_self_deploying"
  description                                  = "self deploying autopilot profile with os default locale"
  device_name_template                         = "thing-%RAND:2%"
  locale                                       = "os-default"
  preprovisioning_allowed                      = false
  device_type                                  = "windowsPc"
  hardware_hash_extraction_enabled             = true
  role_scope_tag_ids                           = ["0"]
  device_join_type                             = "microsoft_entra_joined"
  hybrid_azure_ad_join_skip_connectivity_check = false

  out_of_box_experience_setting = {
    device_usage_type               = "shared"
    privacy_settings_hidden         = true
    eula_hidden                     = true
    user_type                       = "standard"
    keyboard_selection_page_skipped = true
  }
}

# HoloLens Deployment Profile Example
resource "microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile" "hololens" {
  display_name                                 = "hololens"
  description                                  = "hololens autopilot profile with os default locale"
  device_name_template                         = "thing-%RAND:2%"
  locale                                       = "zh-HK"
  preprovisioning_allowed                      = false
  device_type                                  = "holoLens"
  hardware_hash_extraction_enabled             = false
  role_scope_tag_ids                           = ["0"]
  device_join_type                             = "microsoft_entra_joined"
  hybrid_azure_ad_join_skip_connectivity_check = false

  out_of_box_experience_setting = {
    device_usage_type               = "shared"
    privacy_settings_hidden         = true
    eula_hidden                     = true
    user_type                       = "standard"
    keyboard_selection_page_skipped = true
  }
}