resource "microsoft365_graph_beta_device_management_macos_software_update_configuration" "minimal" {
  display_name                             = "Test Minimal macOS Software Update Configuration - Unique"
  update_schedule_type                     = "alwaysUpdate"
  critical_update_behavior                 = "installASAP"
  config_data_update_behavior              = "installASAP"
  firmware_update_behavior                 = "installASAP"
  all_other_update_behavior                = "installASAP"
  update_time_window_utc_offset_in_minutes = 0

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}