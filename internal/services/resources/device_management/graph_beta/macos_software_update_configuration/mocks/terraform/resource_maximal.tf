resource "microsoft365_graph_beta_device_management_macos_software_update_configuration" "maximal" {
  display_name                             = "Maximal macOS Software Update Configuration"
  description                              = "This is a comprehensive configuration with all fields populated"
  critical_update_behavior                 = "installASAP"
  config_data_update_behavior              = "notifyOnly"
  firmware_update_behavior                 = "downloadOnly"
  all_other_update_behavior                = "installLater"
  update_schedule_type                     = "updateDuringTimeWindows"
  update_time_window_utc_offset_in_minutes = 60
  max_user_deferrals_count                 = 3
  priority                                 = "high"
  role_scope_tag_ids                       = ["0", "1"]

  assignments = {
    all_devices = true
    all_users   = false
  }
} 