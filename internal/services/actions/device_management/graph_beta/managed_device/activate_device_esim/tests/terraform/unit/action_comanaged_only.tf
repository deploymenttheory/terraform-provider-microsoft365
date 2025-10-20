action "microsoft365_graph_beta_device_management_managed_device_activate_device_esim" "comanaged_only" {
  comanaged_devices = [
    {
      device_id   = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
      carrier_url = "https://carrier.example.com/esim/activate?code=comanaged789"
    }
  ]

  timeouts = {
    create = "5m"
  }
}