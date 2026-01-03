action "microsoft365_graph_beta_device_management_managed_device_activate_device_esim" "maximal" {
  managed_devices = [
    {
      device_id   = "12345678-1234-1234-1234-123456789abc"
      carrier_url = "https://carrier.example.com/esim/activate?token=managed123"
    },
    {
      device_id   = "87654321-4321-4321-4321-987654321cba"
      carrier_url = "https://carrier.example.com/esim/activate?token=managed456"
    }
  ]

  comanaged_devices = [
    {
      device_id   = "11111111-2222-3333-4444-555555555555"
      carrier_url = "https://carrier.example.com/esim/activate?token=comanaged123"
    },
    {
      device_id   = "66666666-7777-8888-9999-000000000000"
      carrier_url = "https://carrier.example.com/esim/activate?token=comanaged456"
    }
  ]

  timeouts = {
    invoke = "10m"
  }
}