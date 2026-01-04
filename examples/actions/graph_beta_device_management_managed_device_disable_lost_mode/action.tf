# Example 1: Disable lost mode for a single recovered device - Minimal
action "microsoft365_graph_beta_device_management_managed_device_disable_lost_mode" "disable_single_recovered" {
  config {
    managed_device_ids = [
      "12345678-1234-1234-1234-123456789abc"
    ]
  }
}

# Example 2: Disable lost mode for multiple recovered devices
action "microsoft365_graph_beta_device_management_managed_device_disable_lost_mode" "disable_batch_recovered" {
  config {
    managed_device_ids = [
      "12345678-1234-1234-1234-123456789abc",
      "87654321-4321-4321-4321-ba9876543210",
      "abcdef12-3456-7890-abcd-ef1234567890"
    ]

    timeouts = {
      invoke = "10m"
    }
  }
}

# Example 3: Disable lost mode with validation - Maximal
action "microsoft365_graph_beta_device_management_managed_device_disable_lost_mode" "disable_with_validation" {
  config {
    managed_device_ids = [
      "12345678-1234-1234-1234-123456789abc",
      "87654321-4321-4321-4321-ba9876543210"
    ]

    comanaged_device_ids = [
      "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
    ]

    ignore_partial_failures = true
    validate_device_exists  = true

    timeouts = {
      invoke = "5m"
    }
  }
}

# Example 4: Disable lost mode for recovered iOS devices
data "microsoft365_graph_beta_device_management_managed_device" "ios_in_lost_mode" {
  filter_type  = "odata"
  odata_filter = "(operatingSystem eq 'iOS') and (lostModeState eq 'enabled')"
}

action "microsoft365_graph_beta_device_management_managed_device_disable_lost_mode" "disable_recovered_ios" {
  config {
    managed_device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.ios_in_lost_mode.items : device.id]

    validate_device_exists = true

    timeouts = {
      invoke = "10m"
    }
  }
}

# Example 5: Disable lost mode for user's recovered devices
data "microsoft365_graph_beta_device_management_managed_device" "user_devices_in_lost_mode" {
  filter_type  = "odata"
  odata_filter = "(userPrincipalName eq 'user@example.com') and (lostModeState eq 'enabled')"
}

action "microsoft365_graph_beta_device_management_managed_device_disable_lost_mode" "disable_user_recovered" {
  config {
    managed_device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.user_devices_in_lost_mode.items : device.id]

    ignore_partial_failures = true

    timeouts = {
      invoke = "5m"
    }
  }
}

# Output examples
output "disabled_lost_mode_count" {
  value       = length(action.microsoft365_graph_beta_device_management_managed_device_disable_lost_mode.disable_batch_recovered.config.managed_device_ids)
  description = "Number of devices that had lost mode disabled"
}
