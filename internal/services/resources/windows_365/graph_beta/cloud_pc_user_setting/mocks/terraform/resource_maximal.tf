resource "microsoft365_graph_beta_windows_365_cloud_pc_user_setting" "maximal" {
  display_name        = "Test Maximal User Setting"
  local_admin_enabled = true
  reset_enabled       = true

  restore_point_setting = {
    user_restore_enabled = true
    frequency_in_hours   = 12
  }

  cross_region_disaster_recovery_setting = {
    cross_region_disaster_recovery_enabled      = true
    maintain_cross_region_restore_point_enabled = true
    user_initiated_disaster_recovery_allowed    = true
    disaster_recovery_type                      = "premium"

    disaster_recovery_network_setting = {
      region_name  = "automatic"
      region_group = "usCentral"
    }
  }

  notification_setting = {
    restart_prompts_disabled = false
  }

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}
