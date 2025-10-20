action "microsoft365_graph_beta_device_management_managed_device_clean_windows_device" "minimal" {
  managed_devices = [
    {
      device_id      = "12345678-1234-1234-1234-123456789abc"
      keep_user_data = false
    }
  ]
}
