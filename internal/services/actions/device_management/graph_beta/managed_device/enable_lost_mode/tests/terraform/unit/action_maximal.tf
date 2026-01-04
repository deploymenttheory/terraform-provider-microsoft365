action "microsoft365_graph_beta_device_management_managed_device_enable_lost_mode" "maximal" {
  config {
    managed_devices = [
      {
        device_id    = "12345678-1234-1234-1234-123456789abc"
        message      = "This device has been lost"
        phone_number = "+1234567890"
        footer       = "Please return to owner"
      }
    ]

    comanaged_devices = [
      {
        device_id    = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
        message      = "Lost device"
        phone_number = "+0987654321"
      }
    ]

    ignore_partial_failures = false
    validate_device_exists  = true

    timeouts = {
      invoke = "5m"
    }
  }
}
