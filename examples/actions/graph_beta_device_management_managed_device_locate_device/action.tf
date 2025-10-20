# Example 1: Locate a single lost device
action "microsoft365_graph_beta_device_management_managed_device_locate_device" "locate_single" {

  device_ids = [
    "12345678-1234-1234-1234-123456789abc"
  ]

  timeouts = {
    invoke = "5m"
  }
}

# Example 2: Locate multiple devices
action "microsoft365_graph_beta_device_management_managed_device_locate_device" "locate_batch" {

  device_ids = [
    "12345678-1234-1234-1234-123456789abc",
    "87654321-4321-4321-4321-ba9876543210",
    "abcdef12-3456-7890-abcd-ef1234567890"
  ]

  timeouts = {
    invoke = "10m"
  }
}

# Example 3: Locate devices in lost mode
data "microsoft365_graph_beta_device_management_managed_device" "devices_in_lost_mode" {
  filter_type  = "odata"
  odata_filter = "lostModeState eq 'enabled'"
}

action "microsoft365_graph_beta_device_management_managed_device_locate_device" "locate_lost_mode_devices" {

  device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.devices_in_lost_mode.items : device.id]

  timeouts = {
    invoke = "10m"
  }
}

# Example 4: Locate all iOS devices for a specific user
data "microsoft365_graph_beta_device_management_managed_device" "user_ios_devices" {
  filter_type  = "odata"
  odata_filter = "userId eq 'user@example.com' and operatingSystem eq 'iOS'"
}

action "microsoft365_graph_beta_device_management_managed_device_locate_device" "locate_user_ios" {

  device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.user_ios_devices.items : device.id]

  timeouts = {
    invoke = "5m"
  }
}

# Example 5: Locate supervised iOS/iPadOS devices
data "microsoft365_graph_beta_device_management_managed_device" "supervised_apple_devices" {
  filter_type  = "odata"
  odata_filter = "(operatingSystem eq 'iOS' or operatingSystem eq 'iPadOS') and isSupervised eq true"
}

action "microsoft365_graph_beta_device_management_managed_device_locate_device" "locate_supervised_apple" {

  device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.supervised_apple_devices.items : device.id]

  timeouts = {
    invoke = "10m"
  }
}

# Example 6: Locate Android Enterprise devices
data "microsoft365_graph_beta_device_management_managed_device" "android_devices" {
  filter_type  = "odata"
  odata_filter = "operatingSystem eq 'Android'"
}

action "microsoft365_graph_beta_device_management_managed_device_locate_device" "locate_android" {

  device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.android_devices.items : device.id]

  timeouts = {
    invoke = "10m"
  }
}

# Example 7: Locate macOS devices
data "microsoft365_graph_beta_device_management_managed_device" "macos_devices" {
  filter_type  = "odata"
  odata_filter = "operatingSystem eq 'macOS'"
}

action "microsoft365_graph_beta_device_management_managed_device_locate_device" "locate_macos" {

  device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.macos_devices.items : device.id]

  timeouts = {
    invoke = "10m"
  }
}

# Example 8: Locate corporate-owned devices
data "microsoft365_graph_beta_device_management_managed_device" "corporate_devices" {
  filter_type  = "odata"
  odata_filter = "managedDeviceOwnerType eq 'company'"
}

action "microsoft365_graph_beta_device_management_managed_device_locate_device" "locate_corporate" {

  device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.corporate_devices.items : device.id]

  timeouts = {
    invoke = "15m"
  }
}

# Output examples
output "located_device_count" {
  value       = length(action.locate_batch.device_ids)
  description = "Number of devices that received locate command"
}

output "lost_mode_device_count" {
  value       = length(action.locate_lost_mode_devices.device_ids)
  description = "Number of devices in lost mode that received locate command"
}

# Important Notes:
#
# What is Device Location?
# - Requests device to report its current geographic position
# - Uses GPS, WiFi triangulation, or cellular tower data
# - Location data displayed in Intune admin center
# - Can be requested multiple times to track device movement
# - Essential feature for lost device recovery
#
# When to Use Locate Device:
# - Device reported lost or stolen
# - User cannot find their device
# - Security incident requiring device location verification
# - Before performing remote wipe on lost device
# - Tracking device movement over time
# - Asset recovery operations
# - Compliance verification for device location
#
# Platform Support:
# - iOS/iPadOS: Full support (iOS 9.3+, supervised devices)
# - Android: Supported (fully managed, dedicated, work profile)
# - macOS: Supported (macOS 10.13+, supervised or user-approved MDM)
# - Windows: Limited support (requires location services enabled)
#
# Requirements:
# - Device must be online to receive and process command
# - Location services must be enabled on device
# - Device must have GPS or location capability
# - User permissions may be required (varies by platform)
#
# How It Works:
# 1. Issue locate command via this action
# 2. Command is queued in Intune
# 3. Device checks in and receives command
# 4. Device queries its location hardware
# 5. Device reports location back to Intune
# 6. Location displayed in Intune portal with timestamp
# 7. Location includes latitude, longitude, accuracy radius
#
# Location Data:
# - Displayed in Intune admin center > Devices > Device name > Hardware
# - Shows latitude and longitude coordinates
# - Includes location accuracy radius (meters)
# - Timestamp of when location was captured
# - May show altitude (if available)
# - Location history may be available (platform dependent)
#
# Privacy Considerations:
# - Users may be notified that location was requested
# - Follow organizational privacy policies
# - Document business justification for location requests
# - Consider legal requirements in your jurisdiction
# - Some regions have strict privacy laws regarding employee tracking
# - Ensure device ownership agreements cover location tracking
#
# Common Use Cases:
# - Lost Device Recovery: User misplaced device, needs to find it
# - Theft Investigation: Device stolen, locate for law enforcement
# - Asset Tracking: Verify device location for inventory purposes
# - Security Incident: Locate device involved in security breach
# - Compliance: Verify devices are in authorized locations
# - Lost Mode Tracking: Track device in lost mode over time
#
# Limitations:
# - Device must be powered on and have network connectivity
# - Indoor locations may be less accurate (no GPS)
# - Accuracy depends on available location technologies
# - Location services must be enabled by user (some platforms)
# - Location data may be stale if device offline
# - Some enterprise networks block location services
#
# Best Practices:
# - Use in conjunction with lost mode for lost devices
# - Request location multiple times to track movement
# - Combine with remote lock if device in unauthorized location
# - Document reason for location request (audit trail)
# - Follow up with appropriate action based on location
# - Consider local laws before tracking personal devices
# - Train users that corporate devices may be located
#
# Workflow Example:
# 1. User reports device lost
# 2. IT enables lost mode (if supported)
# 3. IT issues locate device command
# 4. IT checks Intune portal for location data
# 5. If device found: Coordinate with user for recovery
# 6. If device not recovered: Consider remote wipe
# 7. Document actions taken for compliance
#
# Integration with Other Actions:
# - Often used with: Enable lost mode (iOS/iPadOS)
# - May be followed by: Remote lock
# - Last resort action: Remote wipe
# - Can be combined with: Disable lost mode (after recovery)
#
# Troubleshooting:
#
# Issue: Location not updating
# Solution: Verify device is online and location services enabled
#
# Issue: Location shows "Unknown" or error
# Solution: Device may be offline or location services disabled
#
# Issue: Location accuracy poor (large radius)
# Solution: Device may be indoors (no GPS) or using WiFi/cellular only
#
# Issue: Old timestamp on location
# Solution: Device hasn't checked in recently, wait or force sync
#
# Issue: No location data appears in portal
# Solution: Command may not have reached device yet, or platform unsupported
#
# Reference:
# https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-locatedevice?view=graph-rest-beta

