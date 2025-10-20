# Example 1: Disable lost mode for a single recovered device
action "microsoft365_graph_beta_device_management_managed_device_disable_lost_mode" "disable_single_recovered" {

  device_ids = [
    "12345678-1234-1234-1234-123456789abc"
  ]

  timeouts = {
    invoke = "5m"
  }
}

# Example 2: Disable lost mode for multiple recovered devices
action "microsoft365_graph_beta_device_management_managed_device_disable_lost_mode" "disable_batch_recovered" {

  device_ids = [
    "12345678-1234-1234-1234-123456789abc",
    "87654321-4321-4321-4321-ba9876543210",
    "abcdef12-3456-7890-abcd-ef1234567890"
  ]

  timeouts = {
    invoke = "10m"
  }
}

# Example 3: Disable lost mode for iOS devices in lost mode state
data "microsoft365_graph_beta_device_management_managed_device" "ios_in_lost_mode" {
  filter_type  = "odata"
  odata_filter = "operatingSystem eq 'iOS' and lostModeState eq 'enabled'"
}

action "microsoft365_graph_beta_device_management_managed_device_disable_lost_mode" "disable_recovered_ios" {

  device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.ios_in_lost_mode.items : device.id]

  timeouts = {
    invoke = "10m"
  }
}

# Example 4: Disable lost mode for iPadOS devices that were recovered
data "microsoft365_graph_beta_device_management_managed_device" "ipados_in_lost_mode" {
  filter_type  = "odata"
  odata_filter = "operatingSystem eq 'iPadOS' and lostModeState eq 'enabled'"
}

action "microsoft365_graph_beta_device_management_managed_device_disable_lost_mode" "disable_recovered_ipados" {

  device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.ipados_in_lost_mode.items : device.id]

  timeouts = {
    invoke = "10m"
  }
}

# Example 5: Disable lost mode for a specific user's recovered devices
data "microsoft365_graph_beta_device_management_managed_device" "user_devices_in_lost_mode" {
  filter_type  = "odata"
  odata_filter = "userId eq 'user@example.com' and lostModeState eq 'enabled'"
}

action "microsoft365_graph_beta_device_management_managed_device_disable_lost_mode" "disable_user_recovered" {

  device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.user_devices_in_lost_mode.items : device.id]

  timeouts = {
    invoke = "5m"
  }
}

# Example 6: Disable lost mode for supervised iOS devices (confirmed recovered)
data "microsoft365_graph_beta_device_management_managed_device" "supervised_ios_lost_mode" {
  filter_type  = "odata"
  odata_filter = "operatingSystem eq 'iOS' and isSupervised eq true and lostModeState eq 'enabled'"
}

action "microsoft365_graph_beta_device_management_managed_device_disable_lost_mode" "disable_supervised_recovered" {

  device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.supervised_ios_lost_mode.items : device.id]

  timeouts = {
    invoke = "10m"
  }
}

# Output examples
output "disabled_lost_mode_count" {
  value       = length(action.disable_batch_recovered.device_ids)
  description = "Number of devices that had lost mode disabled"
}

output "recovered_ios_count" {
  value       = length(action.disable_recovered_ios.device_ids)
  description = "Number of iOS devices returned to normal operation"
}

# Important Notes:
# Lost Mode Features:
# - Only available for iOS and iPadOS devices (iOS 9.3+)
# - Devices must be supervised to use lost mode
# - Lost mode locks device and displays custom message with contact info
# - Lost mode enables device location tracking
# - Disabling lost mode returns device to normal operation
#
# When to Disable Lost Mode:
# - Device has been physically recovered
# - Device location has been confirmed and it's safe
# - User has regained possession of their device
# - Lost mode was enabled in error
# - Device is being returned to service
#
# What Happens When Lost Mode is Disabled:
# - Device returns to normal operation
# - Custom lock screen message is removed
# - Device can be unlocked with regular passcode
# - Location tracking associated with lost mode stops
# - Device becomes fully functional again
#
# Platform Requirements:
# - iOS/iPadOS: Fully supported (iOS 9.3+, supervised devices)
# - macOS: Not supported (lost mode is iOS/iPadOS only)
# - Windows: Not supported
# - Android: Not supported
#
# Best Practices:
# - Verify device has been physically recovered before disabling
# - Document recovery details for audit/compliance
# - Confirm user identity before returning device
# - Check device hasn't been compromised during loss
# - Consider security policy before re-enabling full access
# - Update device compliance status if needed
#
# Security Considerations:
# - Ensure device recovery is legitimate
# - Check for signs of tampering
# - Verify no unauthorized access occurred
# - Consider resetting passcode as additional security measure
# - Review device logs for suspicious activity
# - Confirm device certificate/profile integrity
#
# Related Actions:
# - Enable Lost Mode: Use Intune portal (not yet available as provider action)
# - Remote Lock: Lock device immediately without lost mode
# - Locate Device: Use Intune portal to track device location
# - Wipe Device: Factory reset if device cannot be recovered
# - Reset Passcode: Change device passcode remotely
#
# Reference:
# https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-disablelostmode?view=graph-rest-beta

