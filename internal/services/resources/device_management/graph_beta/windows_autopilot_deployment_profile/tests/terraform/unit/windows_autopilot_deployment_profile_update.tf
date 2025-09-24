resource "microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile" "test" {
  display_name                      = "Updated Windows Autopilot Deployment Profile"
  description                       = "Updated description for Windows Autopilot deployment profile"
  locale                           = "fr-FR"
  device_join_type                 = "microsoft_entra_joined"
  hardware_hash_extraction_enabled = true
  device_name_template             = "UPD-%SERIAL%"
  preprovisioning_allowed          = true

  out_of_box_experience_setting = {
    privacy_settings_hidden          = true
    eula_hidden                     = false
    user_type                       = "standard"
    device_usage_type               = "singleUser"
    keyboard_selection_page_skipped = false
    escape_link_hidden              = true
  }

  enrollment_status_screen_settings = {
    hide_installation_progress                                = true
    allow_device_use_before_profile_and_app_install_complete = false
    block_device_setup_retry_by_user                        = false
    allow_log_collection_on_install_failure                 = true
    custom_error_message                                     = "Updated error message"
    install_progress_timeout_in_minutes                     = 90
    allow_device_use_on_install_failure                     = false
  }
}