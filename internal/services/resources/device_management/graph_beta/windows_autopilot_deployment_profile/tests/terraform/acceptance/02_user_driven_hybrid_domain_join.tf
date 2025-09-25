# User-Driven with Japanese Language and Allow Pre-provisioned Deployment
resource "microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile" "user_driven_japanese_preprovisioned_with_assignments" {
  display_name                                 = "acc_test_user_driven_japanese_preprovisioned"
  description                                  = "user driven autopilot profile with japanese locale and allow pre provisioned deployment"
  locale                                       = "ja-JP"
  preprovisioning_allowed                      = true
  device_type                                  = "windowsPc"
  hardware_hash_extraction_enabled             = true
  role_scope_tag_ids                           = ["0"]
  device_join_type                             = "microsoft_entra_hybrid_joined"
  hybrid_azure_ad_join_skip_connectivity_check = true

  out_of_box_experience_setting = {
    device_usage_type               = "singleUser"
    privacy_settings_hidden         = true
    eula_hidden                     = true
    user_type                       = "standard"
    keyboard_selection_page_skipped = true
  }

  // Optional assignments, can be either group based or all devices based
  assignments = [
    {
      type     = "groupAssignmentTarget"
      group_id = microsoft365_graph_beta_groups_group.acc_test_group_1.id
    }
  ]
}