# Example 1: Enable lost mode for a single device
action "microsoft365_graph_beta_device_management_managed_device_enable_lost_mode" "enable_single_lost_device" {

  managed_devices {
    device_id    = "12345678-1234-1234-1234-123456789abc"
    message      = "This device has been lost. Please contact IT at 555-0123 to return."
    phone_number = "555-0123"
    footnote     = "Property of Contoso Corporation"
  }

  timeouts = {
    invoke = "5m"
  }
}

# Example 2: Enable lost mode for multiple devices with different messages
action "microsoft365_graph_beta_device_management_managed_device_enable_lost_mode" "enable_multiple_lost_devices" {

  managed_devices {
    device_id    = "12345678-1234-1234-1234-123456789abc"
    message      = "Lost iPhone - Please call John at IT to return"
    phone_number = "+1-555-0123"
    footnote     = "Reward available for return"
  }

  managed_devices {
    device_id    = "87654321-4321-4321-4321-ba9876543210"
    message      = "Lost iPad - Contact Mary in HR to return"
    phone_number = "+1-555-0456"
    footnote     = "Property of Contoso - Finance Department"
  }

  timeouts = {
    invoke = "10m"
  }
}

# Example 3: Enable lost mode for supervised iOS devices
data "microsoft365_graph_beta_device_management_managed_device" "supervised_ios" {
  filter_type  = "odata"
  odata_filter = "operatingSystem eq 'iOS' and isSupervised eq true and lostModeState eq 'disabled'"
}

action "microsoft365_graph_beta_device_management_managed_device_enable_lost_mode" "enable_for_supervised_ios" {

  dynamic "managed_devices" {
    for_each = data.microsoft365_graph_beta_device_management_managed_device.supervised_ios.items
    content {
      device_id    = managed_devices.value.id
      message      = "This iOS device has been lost. Please contact IT immediately."
      phone_number = "555-IT-HELP"
      footnote     = "Return to Contoso Corporation IT Department"
    }
  }

  timeouts = {
    invoke = "15m"
  }
}

# Example 4: Enable lost mode with user-specific messages
data "microsoft365_graph_beta_device_management_managed_device" "user_devices" {
  filter_type  = "odata"
  odata_filter = "userId eq 'user@example.com' and operatingSystem eq 'iOS'"
}

action "microsoft365_graph_beta_device_management_managed_device_enable_lost_mode" "enable_user_lost_devices" {

  dynamic "managed_devices" {
    for_each = data.microsoft365_graph_beta_device_management_managed_device.user_devices.items
    content {
      device_id    = managed_devices.value.id
      message      = format("Lost device belonging to %s. Please contact IT to return.", data.microsoft365_graph_beta_device_management_managed_device.user_devices.items[0].user_principal_name)
      phone_number = "555-0123"
      footnote     = "Corporate Device - Immediate Return Required"
    }
  }

  timeouts = {
    invoke = "10m"
  }
}

# Example 5: Enable lost mode for co-managed devices
action "microsoft365_graph_beta_device_management_managed_device_enable_lost_mode" "enable_comanaged_lost" {

  comanaged_devices {
    device_id    = "abcdef12-3456-7890-abcd-ef1234567890"
    message      = "This co-managed device has been lost. Contact IT."
    phone_number = "+1-555-9999"
  }

  timeouts = {
    invoke = "5m"
  }
}

# Example 6: Enable lost mode with emergency contact information
action "microsoft365_graph_beta_device_management_managed_device_enable_lost_mode" "enable_with_emergency_contact" {

  managed_devices {
    device_id    = "12345678-1234-1234-1234-123456789abc"
    message      = "LOST DEVICE - This device contains sensitive corporate data. Please return immediately!"
    phone_number = "+1-555-SECURITY"
    footnote     = "24/7 Security Hotline: +1-555-SEC-HELP | Reward: $100"
  }

  timeouts = {
    invoke = "5m"
  }
}

# Output examples
output "enabled_lost_mode_count" {
  value       = length(action.enable_multiple_lost_devices.managed_devices)
  description = "Number of devices that had lost mode enabled"
}

output "secured_ios_count" {
  value       = length(action.enable_for_supervised_ios.managed_devices)
  description = "Number of iOS devices now secured with lost mode"
}

# Important Notes:
# Lost Mode Features:
# - Only available for iOS and iPadOS devices (iOS 9.3+)
# - Devices must be supervised to use lost mode
# - Lost mode locks device and displays custom message with contact info
# - Lost mode enables device location tracking
# - Each device can have a unique message, phone number, and footnote
#
# When to Enable Lost Mode:
# - Device has been reported lost or stolen
# - Need to remotely lock and secure device immediately
# - Need to display recovery contact information
# - Need to track device location for recovery
# - Prevent unauthorized access to corporate data
# - User cannot be reached to secure device manually
#
# What Happens When Lost Mode is Enabled:
# - Device is immediately locked
# - Custom message displayed on lock screen with contact info
# - Location tracking is enabled
# - Device cannot be unlocked without proper credentials
# - Device data remains encrypted and protected
# - User must contact provided phone number to recover device
#
# Platform Requirements:
# - iOS/iPadOS: Fully supported (iOS 9.3+, supervised devices)
# - macOS: Not supported (lost mode is iOS/iPadOS only)
# - Windows: Not supported
# - Android: Not supported
#
# Message Best Practices:
# - Include clear instructions for returning the device
# - Provide a contact phone number that will be monitored
# - Keep message concise and professional
# - Include organization name or identification
# - Avoid including sensitive information
# - Consider including reward information if applicable
#
# Security Considerations:
# - Ensure contact number is monitored 24/7 if possible
# - Document when and why lost mode was enabled
# - Plan recovery process before enabling
# - Consider legal implications of location tracking
# - Have escalation procedure for device not returned
# - Verify identity of person returning device
#
# Related Actions:
# - Disable Lost Mode: Use to return device to normal operation after recovery
# - Remote Lock: Alternative for locking device without full lost mode
# - Locate Device: Use Intune portal to track device location
# - Wipe Device: Factory reset if device cannot be recovered
# - Reset Passcode: Change device passcode remotely
#
# Reference:
# https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-enablelostmode?view=graph-rest-beta


