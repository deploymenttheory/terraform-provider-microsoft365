action "microsoft365_graph_beta_device_management_managed_device_set_device_name" "maximal" {
  config {
    managed_devices = [
      {
        device_id   = "00000000-0000-0000-0000-000000000001"
        device_name = "NYC-Marketing-Laptop-01"
      },
      {
        device_id   = "00000000-0000-0000-0000-000000000002"
        device_name = "NYC-IT-Desktop-05"
      }
    ]
    comanaged_devices = [
      {
        device_id   = "00000000-0000-0000-0000-000000000003"
        device_name = "NYC-HR-Laptop-03"
      }
    ]
    ignore_partial_failures = true
    validate_device_exists  = true

    timeouts = {
      invoke = "5m"
    }
  }
}

