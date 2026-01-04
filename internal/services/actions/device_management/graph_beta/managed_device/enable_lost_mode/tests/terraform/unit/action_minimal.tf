action "microsoft365_graph_beta_device_management_managed_device_enable_lost_mode" "minimal" {
  config {
    managed_devices = [
      {
        device_id    = "12345678-1234-1234-1234-123456789abc"
        message      = "This device has been lost"
        phone_number = "+1234567890"
      }
    ]
  }
}
