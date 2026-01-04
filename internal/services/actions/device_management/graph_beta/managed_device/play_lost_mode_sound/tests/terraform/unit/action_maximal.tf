action "microsoft365_graph_beta_device_management_managed_device_play_lost_mode_sound" "maximal" {
  config {
    managed_devices = [
      {
        device_id           = "00000000-0000-0000-0000-000000000001"
        duration_in_minutes = "5"
      },
      {
        device_id           = "00000000-0000-0000-0000-000000000002"
        duration_in_minutes = "10"
      }
    ]
    comanaged_devices = [
      {
        device_id           = "00000000-0000-0000-0000-000000000003"
        duration_in_minutes = "3"
      }
    ]
    ignore_partial_failures = true
    validate_device_exists  = true

    timeouts = {
      invoke = "5m"
    }
  }
}

