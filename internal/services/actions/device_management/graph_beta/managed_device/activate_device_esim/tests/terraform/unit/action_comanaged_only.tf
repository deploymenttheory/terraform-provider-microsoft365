action "microsoft365_graph_beta_device_management_managed_device_activate_device_esim" "comanaged_only" {
  config {
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

    timeouts = {
      invoke = "5m"
    }
  }
}

