resource "microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile" "test" {
  display_name                                    = "Test Windows Autopilot Deployment Profile - maximal"
  description                                     = "Test description for Windows Autopilot deployment profile with maximal configuration"
  locale                                         = "en-US"
  device_join_type                               = "microsoft_entra_joined"
  hardware_hash_extraction_enabled               = true
  device_name_template                           = "AUTO-%SERIAL%"
  device_type                                    = "windowsPc"
  preprovisioning_allowed                        = true
  role_scope_tag_ids                             = ["0", "1", "2"]
  management_service_app_id                      = "12345678-1234-1234-1234-123456789abc"
  hybrid_azure_ad_join_skip_connectivity_check  = true

  out_of_box_experience_setting = {
    privacy_settings_hidden          = true
    eula_hidden                     = true
    user_type                       = "standard"
    device_usage_type               = "shared"
    keyboard_selection_page_skipped = true
    escape_link_hidden              = true
  }

  enrollment_status_screen_settings = {
    hide_installation_progress                                 = true
    allow_device_use_before_profile_and_app_install_complete  = true
    block_device_setup_retry_by_user                         = true
    allow_log_collection_on_install_failure                  = true
    custom_error_message                                      = "Custom error message for installation failure"
    install_progress_timeout_in_minutes                      = 60
    allow_device_use_on_install_failure                      = true
  }
}