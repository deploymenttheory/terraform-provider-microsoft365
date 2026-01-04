action "microsoft365_graph_beta_device_management_managed_device_pause_configuration_refresh" "maximal" {
  config {
    managed_devices = [
      {
        device_id                    = "00000000-0000-0000-0000-000000000001"
        pause_time_period_in_minutes = 120
      },
      {
        device_id                    = "00000000-0000-0000-0000-000000000002"
        pause_time_period_in_minutes = 240
      }
    ]
    comanaged_devices = [
      {
        device_id                    = "00000000-0000-0000-0000-000000000003"
        pause_time_period_in_minutes = 480
      }
    ]
    ignore_partial_failures = true
    validate_device_exists  = true

    timeouts = {
      invoke = "5m"
    }
  }
}

