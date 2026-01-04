action "microsoft365_graph_beta_device_management_managed_device_activate_device_esim" "maximal" {
  config {
    managed_devices = [
      {
        device_id   = "12345678-1234-1234-1234-123456789001"
        carrier_url = "https://carrier.example.com/esim/activate?token=managed1"
      },
      {
        device_id   = "12345678-1234-1234-1234-123456789002"
        carrier_url = "https://carrier.example.com/esim/activate?token=managed2"
      }
    ]

    comanaged_devices = [
      {
        device_id   = "87654321-4321-4321-4321-987654321001"
        carrier_url = "https://carrier.example.com/esim/activate?token=comanaged1"
      },
      {
        device_id   = "87654321-4321-4321-4321-987654321002"
        carrier_url = "https://carrier.example.com/esim/activate?token=comanaged2"
      }
    ]

    ignore_partial_failures = true
    validate_device_exists  = true

    timeouts = {
      invoke = "10m"
    }
  }
}

