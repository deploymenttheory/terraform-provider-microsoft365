action "microsoft365_graph_beta_device_management_managed_device_delete_user_from_shared_apple_device" "minimal" {
  config {
    managed_devices = [
      {
        device_id           = "12345678-1234-1234-1234-123456789abc"
        user_principal_name = "user@example.com"
      }
    ]
  }
}

