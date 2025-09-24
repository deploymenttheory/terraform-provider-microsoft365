resource "microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile" "self_deploying" {
  display_name                                 = "unit_test_self_deploying"
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