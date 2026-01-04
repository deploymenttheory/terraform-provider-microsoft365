action "microsoft365_graph_beta_device_management_managed_device_windows_defender_scan" "maximal" {
  config {
    managed_devices = [
      {
        device_id  = "00000000-0000-0000-0000-000000000001"
        quick_scan = true
      },
      {
        device_id  = "00000000-0000-0000-0000-000000000002"
        quick_scan = false
      }
    ]
    comanaged_devices = [
      {
        device_id  = "00000000-0000-0000-0000-000000000003"
        quick_scan = true
      }
    ]
    ignore_partial_failures = true
    validate_device_exists  = true

    timeouts = {
      invoke = "5m"
    }
  }
}

