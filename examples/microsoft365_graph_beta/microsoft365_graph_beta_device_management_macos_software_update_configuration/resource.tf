resource "microsoft365_graph_beta_device_management_macos_software_update_configuration" "example" {
  display_name                        = "Example macOS Software Update Configuration"
  description                         = "Example configuration for macOS software updates"
  critical_update_behavior            = "default"
  config_data_update_behavior         = "default"
  firmware_update_behavior            = "default"
  all_other_update_behavior           = "default"
  update_schedule_type                = "alwaysUpdate"
  update_time_window_utc_offset_in_minutes = 0
  max_user_deferrals_count            = 0
  priority                            = "low"

  custom_update_time_windows = [
    {
      start_day  = "monday"
      end_day    = "monday"
      start_time = "01:00:00"
      end_time   = "02:00:00"
    },
    {
      start_day  = "friday"
      end_day    = "friday"
      start_time = "03:00:00"
      end_time   = "04:00:00"
    }
  ]

  assignments = {
    all_devices        = true
    all_users          = false
    include_group_ids  = ["00000000-0000-0000-0000-000000000001", "00000000-0000-0000-0000-000000000002"]
    exclude_group_ids  = ["00000000-0000-0000-0000-000000000003"]
  }
}