resource "microsoft365_graph_beta_device_management_macos_software_update_configuration" "minimal" {
  display_name                             = "Minimal macOS Software Update Configuration"
  critical_update_behavior                 = "default"
  config_data_update_behavior              = "default"
  firmware_update_behavior                 = "default"
  all_other_update_behavior                = "default"
  update_schedule_type                     = "alwaysUpdate"
  update_time_window_utc_offset_in_minutes = 0

  assignments = {
    all_devices = false
    all_users   = false
  }
} 