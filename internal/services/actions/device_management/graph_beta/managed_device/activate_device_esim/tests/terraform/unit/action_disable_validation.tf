action "microsoft365_graph_beta_device_management_managed_device_activate_device_esim" "disable_validation" {
  config {
    managed_devices = [
      {
        device_id   = "12345678-1234-1234-1234-123456789abc"
        carrier_url = "https://carrier.example.com/esim/activate?token=test123"
      }
    ]

    validate_device_exists = false

    timeouts = {
      invoke = "5m"
    }
  }
}

