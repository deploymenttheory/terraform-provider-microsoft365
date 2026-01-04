# Example 1: Delete user from a single shared Apple device - Minimal
action "microsoft365_graph_beta_device_management_managed_device_delete_user_from_shared_apple_device" "delete_single" {
  config {
    managed_devices = [
      {
        device_id           = "12345678-1234-1234-1234-123456789abc"
        user_principal_name = "user@example.com"
      }
    ]
  }
}

# Example 2: Delete users from multiple shared Apple devices
action "microsoft365_graph_beta_device_management_managed_device_delete_user_from_shared_apple_device" "delete_multiple" {
  config {
    managed_devices = [
      {
        device_id           = "12345678-1234-1234-1234-123456789abc"
        user_principal_name = "user1@example.com"
      },
      {
        device_id           = "87654321-4321-4321-4321-ba9876543210"
        user_principal_name = "user2@example.com"
      }
    ]

    timeouts = {
      invoke = "10m"
    }
  }
}

# Example 3: Delete with validation - Maximal
action "microsoft365_graph_beta_device_management_managed_device_delete_user_from_shared_apple_device" "delete_maximal" {
  config {
    managed_devices = [
      {
        device_id           = "12345678-1234-1234-1234-123456789abc"
        user_principal_name = "user1@example.com"
      }
    ]

    comanaged_devices = [
      {
        device_id           = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
        user_principal_name = "user2@example.com"
      }
    ]

    ignore_partial_failures = true
    validate_device_exists  = true

    timeouts = {
      invoke = "5m"
    }
  }
}

# Example 4: Delete departing user from all shared iPads
data "microsoft365_graph_beta_device_management_managed_device" "shared_ipads" {
  filter_type  = "odata"
  odata_filter = "(operatingSystem eq 'iPadOS') and (managementMode eq 'shared')"
}

action "microsoft365_graph_beta_device_management_managed_device_delete_user_from_shared_apple_device" "delete_departing_user" {
  config {
    managed_devices = [
      for device in data.microsoft365_graph_beta_device_management_managed_device.shared_ipads.items : {
        device_id           = device.id
        user_principal_name = "departing.user@example.com"
      }
    ]

    validate_device_exists = true

    timeouts = {
      invoke = "15m"
    }
  }
}
