resource "microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile" "test" {
  display_name                      = "Test Windows Autopilot Deployment Profile"
  description                       = "Test description for Windows Autopilot deployment profile"
  locale                           = "en-US"
  device_join_type                 = "microsoft_entra_joined"
  hardware_hash_extraction_enabled = false
  device_name_template             = "AP-%SERIAL%"
  device_type                      = "windowsPc"
  preprovisioning_allowed          = false
  role_scope_tag_ids               = ["0"]

  out_of_box_experience_setting = {
    privacy_settings_hidden          = true
    eula_hidden                     = true
    user_type                       = "administrator"
    device_usage_type               = "singleUser"
    keyboard_selection_page_skipped = true
    escape_link_hidden              = false
  }

  enrollment_status_screen_settings = {
    hide_installation_progress                                = false
    allow_device_use_before_profile_and_app_install_complete = true
    block_device_setup_retry_by_user                        = false
    allow_log_collection_on_install_failure                 = true
    custom_error_message                                     = "Please contact IT support for assistance"
    install_progress_timeout_in_minutes                     = 30
    allow_device_use_on_install_failure                     = false
  }
}