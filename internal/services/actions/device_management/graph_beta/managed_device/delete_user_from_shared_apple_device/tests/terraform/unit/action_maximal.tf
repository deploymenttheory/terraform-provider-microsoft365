action "microsoft365_graph_beta_device_management_managed_device_delete_user_from_shared_apple_device" "maximal" {
  config {
    managed_devices = [
      {
        device_id           = "12345678-1234-1234-1234-123456789abc"
        user_principal_name = "user1@example.com"
      },
      {
        device_id           = "87654321-4321-4321-4321-987654321cba"
        user_principal_name = "user2@example.com"
      }
    ]

    comanaged_devices = [
      {
        device_id           = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
        user_principal_name = "user3@example.com"
      }
    ]

    ignore_partial_failures = false
    validate_device_exists  = true

    timeouts = {
      invoke = "5m"
    }
  }
}

