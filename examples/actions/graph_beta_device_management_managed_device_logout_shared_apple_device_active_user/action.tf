# Example 1: Logout active user from a single shared Apple device - Minimal
action "microsoft365_graph_beta_device_management_managed_device_logout_shared_apple_device_active_user" "logout_single" {
  config {
    device_ids = ["12345678-1234-1234-1234-123456789abc"]
  }
}

# Example 2: Logout active users from multiple shared Apple devices
action "microsoft365_graph_beta_device_management_managed_device_logout_shared_apple_device_active_user" "logout_multiple" {
  config {
    device_ids = [
      "12345678-1234-1234-1234-123456789abc",
      "87654321-4321-4321-4321-ba9876543210",
      "abcdef12-3456-7890-abcd-ef1234567890"
    ]

    timeouts = {
      invoke = "10m"
    }
  }
}

# Example 3: Logout with validation - Maximal
action "microsoft365_graph_beta_device_management_managed_device_logout_shared_apple_device_active_user" "logout_maximal" {
  config {
    device_ids = [
      "12345678-1234-1234-1234-123456789abc",
      "87654321-4321-4321-4321-ba9876543210"
    ]

    ignore_partial_failures = true
    validate_device_exists  = true

    timeouts = {
      invoke = "5m"
    }
  }
}

# Example 4: Logout users from all shared iPads
data "microsoft365_graph_beta_device_management_managed_device" "shared_ipads" {
  filter_type  = "odata"
  odata_filter = "(operatingSystem eq 'iPadOS') and (managementMode eq 'shared')"
}

action "microsoft365_graph_beta_device_management_managed_device_logout_shared_apple_device_active_user" "logout_all_shared_ipads" {
  config {
    device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.shared_ipads.items : device.id]

    validate_device_exists  = true
    ignore_partial_failures = true

    timeouts = {
      invoke = "15m"
    }
  }
}

# Example 5: Logout users from classroom iPads at end of day
data "microsoft365_graph_beta_device_management_managed_device" "classroom_ipads" {
  filter_type  = "odata"
  odata_filter = "(deviceCategoryDisplayName eq 'Classroom') and (operatingSystem eq 'iPadOS')"
}

action "microsoft365_graph_beta_device_management_managed_device_logout_shared_apple_device_active_user" "logout_classroom_ipads" {
  config {
    device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.classroom_ipads.items : device.id]

    timeouts = {
      invoke = "10m"
    }
  }
}
