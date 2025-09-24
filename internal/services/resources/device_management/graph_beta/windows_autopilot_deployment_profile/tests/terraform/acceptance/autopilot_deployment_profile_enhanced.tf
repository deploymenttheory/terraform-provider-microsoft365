resource "microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile" "enhanced" {
  display_name                      = "acc-test-windows-autopilot-deployment-profile-enhanced"
  description                       = "acc-test-windows-autopilot-deployment-profile-enhanced"
  locale                           = "en-US"
  device_join_type                 = "microsoft_entra_joined"
  hardware_hash_extraction_enabled = true
  device_name_template             = "TEST-%RAND:3%"
  device_type                      = "windowsPc"
  preprovisioning_allowed          = true

  out_of_box_experience_setting = {
    privacy_settings_hidden          = true
    eula_hidden                     = true
    user_type                       = "standard"
    device_usage_type               = "shared"
    keyboard_selection_page_skipped = true
    escape_link_hidden              = true
  }

  enrollment_status_screen_settings = {
    hide_installation_progress                                = true
    allow_device_use_before_profile_and_app_install_complete = true
    block_device_setup_retry_by_user                        = true
    allow_log_collection_on_install_failure                 = true
    custom_error_message                                     = "Please contact IT support for assistance"
    install_progress_timeout_in_minutes                     = 120
    allow_device_use_on_install_failure                     = true
  }
}