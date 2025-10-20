# Example 1: Reset passcode for a single device
action "microsoft365_graph_beta_device_management_managed_device_reset_passcode" "reset_single" {

  device_ids = [
    "12345678-1234-1234-1234-123456789abc"
  ]

  timeouts = {
    invoke = "5m"
  }
}

# Example 2: Reset passcodes for multiple devices
action "microsoft365_graph_beta_device_management_managed_device_reset_passcode" "reset_batch" {

  device_ids = [
    "12345678-1234-1234-1234-123456789abc",
    "87654321-4321-4321-4321-ba9876543210",
    "abcdef12-3456-7890-abcd-ef1234567890"
  ]

  timeouts = {
    invoke = "10m"
  }
}

# Example 3: Reset passcodes for locked devices from data source
data "microsoft365_graph_beta_device_management_managed_device" "locked_ios_devices" {
  filter_type  = "odata"
  odata_filter = "operatingSystem eq 'iOS'"
}

action "microsoft365_graph_beta_device_management_managed_device_reset_passcode" "reset_locked_devices" {

  # Filter for supervised iOS devices (required for passcode reset)
  device_ids = [
    for device in data.microsoft365_graph_beta_device_management_managed_device.locked_ios_devices.items :
    device.id
    # Note: In production, you would filter for supervised devices only
  ]

  timeouts = {
    invoke = "15m"
  }
}

# Example 4: Emergency passcode reset for specific user's devices
data "microsoft365_graph_beta_device_management_managed_device" "user_devices" {
  filter_type  = "user_id"
  filter_value = "user@example.com"
}

action "microsoft365_graph_beta_device_management_managed_device_reset_passcode" "reset_user_devices" {

  device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.user_devices.items : device.id]

  timeouts = {
    invoke = "10m"
  }
}

# Example 5: Reset passcodes for Android devices
data "microsoft365_graph_beta_device_management_managed_device" "android_devices" {
  filter_type  = "odata"
  odata_filter = "operatingSystem eq 'Android'"
}

action "microsoft365_graph_beta_device_management_managed_device_reset_passcode" "reset_android" {

  device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.android_devices.items : device.id]

  timeouts = {
    invoke = "15m"
  }
}

# Example 6: Reset passcode for Windows 10 devices
data "microsoft365_graph_beta_device_management_managed_device" "windows_devices" {
  filter_type  = "odata"
  odata_filter = "operatingSystem eq 'Windows' and osVersion startsWith '10.0'"
}

action "microsoft365_graph_beta_device_management_managed_device_reset_passcode" "reset_windows" {

  device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.windows_devices.items : device.id]

  timeouts = {
    invoke = "20m"
  }
}

# Output examples
output "reset_device_count" {
  value       = length(action.reset_batch.device_ids)
  description = "Number of devices for which passcodes were reset"
}

output "locked_devices_reset_count" {
  value       = length(action.reset_locked_devices.device_ids)
  description = "Number of locked devices for which passcodes were reset"
}