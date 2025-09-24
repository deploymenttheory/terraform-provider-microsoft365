resource "microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile" "test" {
  display_name     = "Test Windows Autopilot Deployment Profile"
  description      = "Test description for Windows Autopilot deployment profile"
  device_join_type = "microsoft_entra_joined"
  # This should trigger a validation error - too long (18 characters)
  device_name_template = "VERYLONGDEVICENAME"

  out_of_box_experience_setting = {
    privacy_settings_hidden          = false
    eula_hidden                     = false
    user_type                       = "administrator"
    device_usage_type               = "singleUser"
    keyboard_selection_page_skipped = false
    escape_link_hidden              = false
  }

  enrollment_status_screen_settings = {
    hide_installation_progress                                 = false
    allow_device_use_before_profile_and_app_install_complete  = false
    block_device_setup_retry_by_user                         = false
    allow_log_collection_on_install_failure                  = false
    allow_device_use_on_install_failure                      = false
  }
}