resource "microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile" "hybrid_joined" {
  display_name                                   = "acc-test-windows-autopilot-deployment-profile-hybrid-joined"
  description                                    = "acc-test-windows-autopilot-deployment-profile-hybrid-joined"
  device_join_type                               = "microsoft_entra_hybrid_joined"
  hybrid_azure_ad_join_skip_connectivity_check  = true

  out_of_box_experience_setting = {
    privacy_settings_hidden = false
    eula_hidden            = false
    user_type              = "administrator"
    device_usage_type      = "singleUser"
  }

  enrollment_status_screen_settings = {
    hide_installation_progress           = false
    install_progress_timeout_in_minutes = 90
  }
}