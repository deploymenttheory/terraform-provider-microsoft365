action "microsoft365_graph_beta_device_management_managed_device_activate_device_esim" "ignore_failures" {
  config {
    managed_devices = [
      {
        device_id   = "12345678-1234-1234-1234-123456789001"
        carrier_url = "https://carrier.example.com/esim/activate?token=device1"
      },
      {
        device_id   = "12345678-1234-1234-1234-123456789002"
        carrier_url = "https://carrier.example.com/esim/activate?token=device2"
      },
      {
        device_id   = "12345678-1234-1234-1234-123456789003"
        carrier_url = "https://carrier.example.com/esim/activate?token=device3"
      }
    ]

    ignore_partial_failures = true

    timeouts = {
      invoke = "5m"
    }
  }
}

