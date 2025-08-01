resource "microsoft365_graph_beta_device_management_macos_software_update_configuration" "maximal" {
  display_name                             = "Test Maximal macOS Software Update Configuration - Unique"
  description                              = "Maximal software update configuration for testing with all features"
  update_schedule_type                     = "updateDuringTimeWindows"
  critical_update_behavior                 = "installASAP"
  config_data_update_behavior              = "installASAP"
  firmware_update_behavior                 = "installASAP"
  all_other_update_behavior                = "installASAP"
  update_time_window_utc_offset_in_minutes = -480
  max_user_deferrals_count                 = 5
  priority                                 = "high"
  role_scope_tag_ids                       = ["0", "1"]

  custom_update_time_windows = [
    {
      start_day  = "monday"
      end_day    = "friday"
      start_time = "02:00:00"
      end_time   = "06:00:00"
    },
    {
      start_day  = "saturday"
      end_day    = "sunday"
      start_time = "01:00:00"
      end_time   = "05:00:00"
    }
  ]

  assignments = [
    {
      type     = "groupAssignmentTarget"
      group_id = "44444444-4444-4444-4444-444444444444"
    }
  ]

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}