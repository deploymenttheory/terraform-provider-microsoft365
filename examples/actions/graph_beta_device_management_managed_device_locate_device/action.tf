# Example 1: Locate a single lost device - Minimal
action "microsoft365_graph_beta_device_management_managed_device_locate_device" "locate_single" {
  config {
    device_ids = [
      "12345678-1234-1234-1234-123456789abc"
    ]
  }
}

# Example 2: Locate multiple devices
action "microsoft365_graph_beta_device_management_managed_device_locate_device" "locate_batch" {
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

# Example 3: Locate with validation and failure handling - Maximal
action "microsoft365_graph_beta_device_management_managed_device_locate_device" "locate_with_validation" {
  config {
    device_ids = [
      "12345678-1234-1234-1234-123456789abc",
      "87654321-4321-4321-4321-ba9876543210",
      "abcdef12-3456-7890-abcd-ef1234567890"
    ]

    ignore_partial_failures = true
    validate_device_exists  = true

    timeouts = {
      invoke = "5m"
    }
  }
}

# Example 4: Locate devices in lost mode
data "microsoft365_graph_beta_device_management_managed_device" "devices_in_lost_mode" {
  filter_type  = "odata"
  odata_filter = "lostModeState eq 'enabled'"
}

action "microsoft365_graph_beta_device_management_managed_device_locate_device" "locate_lost_mode_devices" {
  config {
    device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.devices_in_lost_mode.items : device.id]

    validate_device_exists = true

    timeouts = {
      invoke = "10m"
    }
  }
}

# Example 5: Locate all iOS devices for a specific user
data "microsoft365_graph_beta_device_management_managed_device" "user_ios_devices" {
  filter_type  = "odata"
  odata_filter = "(userPrincipalName eq 'user@example.com') and (operatingSystem eq 'iOS')"
}

action "microsoft365_graph_beta_device_management_managed_device_locate_device" "locate_user_ios" {
  config {
    device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.user_ios_devices.items : device.id]

    timeouts = {
      invoke = "5m"
    }
  }
}

# Example 6: Locate supervised iOS/iPadOS devices
data "microsoft365_graph_beta_device_management_managed_device" "supervised_apple_devices" {
  filter_type  = "odata"
  odata_filter = "((operatingSystem eq 'iOS') or (operatingSystem eq 'iPadOS')) and (isSupervised eq true)"
}

action "microsoft365_graph_beta_device_management_managed_device_locate_device" "locate_supervised_apple" {
  config {
    device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.supervised_apple_devices.items : device.id]

    ignore_partial_failures = true

    timeouts = {
      invoke = "10m"
    }
  }
}

# Example 7: Locate corporate-owned devices
data "microsoft365_graph_beta_device_management_managed_device" "corporate_devices" {
  filter_type  = "odata"
  odata_filter = "managedDeviceOwnerType eq 'company'"
}

action "microsoft365_graph_beta_device_management_managed_device_locate_device" "locate_corporate" {
  config {
    device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.corporate_devices.items : device.id]

    timeouts = {
      invoke = "15m"
    }
  }
}

# Output examples
output "located_device_count" {
  value       = length(action.microsoft365_graph_beta_device_management_managed_device_locate_device.locate_batch.config.device_ids)
  description = "Number of devices that received locate command"
}

output "lost_mode_device_count" {
  value       = length(action.microsoft365_graph_beta_device_management_managed_device_locate_device.locate_lost_mode_devices.config.device_ids)
  description = "Number of devices in lost mode that received locate command"
}
