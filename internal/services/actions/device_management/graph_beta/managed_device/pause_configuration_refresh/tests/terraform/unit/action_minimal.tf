action "microsoft365_graph_beta_device_management_managed_device_pause_configuration_refresh" "minimal" {
  config {
    managed_devices = [
      {
        device_id                    = "00000000-0000-0000-0000-000000000001"
        pause_time_period_in_minutes = 60
      }
    ]
  }
}

