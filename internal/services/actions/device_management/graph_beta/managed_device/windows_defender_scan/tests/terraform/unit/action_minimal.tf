action "microsoft365_graph_beta_device_management_managed_device_windows_defender_scan" "minimal" {
  config {
    managed_devices = [
      {
        device_id  = "00000000-0000-0000-0000-000000000001"
        quick_scan = true
      }
    ]
  }
}

