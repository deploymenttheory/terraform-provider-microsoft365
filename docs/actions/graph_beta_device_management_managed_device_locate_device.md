---
page_title: "microsoft365_graph_beta_device_management_managed_device_locate_device Action - terraform-provider-microsoft365"
subcategory: "Device Management"

description: |-
  Triggers device location for one or more managed devices using the /deviceManagement/managedDevices/{managedDeviceId}/locateDevice endpoint. This action requests the device to report its current geographic location, which is then viewable in the Microsoft Intune admin center. The locate device feature is essential for finding lost or stolen devices and is commonly used in conjunction with lost mode.
  Important Notes:
  Device must be online to receive and respond to the locate commandLocation services must be enabled on the deviceDevice must have GPS/location hardware capabilityLocation data is displayed in the Intune portal, not returned via APIMultiple location requests can be sent over time to track device movementLocation accuracy depends on device capabilities (GPS, WiFi, cellular triangulation)
  Use Cases:
  Locating lost or stolen devicesTracking devices in lost modeVerifying device location for security/complianceFinding devices before performing remote wipeAssisting users who have misplaced their devicesAsset tracking and recovery operations
  Platform Support:
  iOS/iPadOS: Fully supported (iOS 9.3+, supervised devices)Android Enterprise: Supported (fully managed, dedicated, and work profile devices)macOS: Supported (macOS 10.13+, supervised devices or user-approved MDM)Windows: Limited support (Windows 10/11 with location services enabled)
  Location Data:
  Location data is displayed in the Intune admin center under device propertiesShows latitude, longitude, altitude (if available)Includes location accuracy radiusDisplays timestamp of when location was capturedLocation history may be available depending on platform
  Privacy Considerations:
  Users may receive notification that device location was requestedLocation tracking should comply with organizational privacy policiesDocument legitimate business reasons for location requestsConsider legal requirements in your jurisdiction
  Reference: Microsoft Graph API - Locate Device https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-locatedevice?view=graph-rest-beta
---

# microsoft365_graph_beta_device_management_managed_device_locate_device (Action)

Triggers device location for one or more managed devices using the `/deviceManagement/managedDevices/{managedDeviceId}/locateDevice` endpoint. This action requests the device to report its current geographic location, which is then viewable in the Microsoft Intune admin center. The locate device feature is essential for finding lost or stolen devices and is commonly used in conjunction with lost mode.

**Important Notes:**
- Device must be online to receive and respond to the locate command
- Location services must be enabled on the device
- Device must have GPS/location hardware capability
- Location data is displayed in the Intune portal, not returned via API
- Multiple location requests can be sent over time to track device movement
- Location accuracy depends on device capabilities (GPS, WiFi, cellular triangulation)

**Use Cases:**
- Locating lost or stolen devices
- Tracking devices in lost mode
- Verifying device location for security/compliance
- Finding devices before performing remote wipe
- Assisting users who have misplaced their devices
- Asset tracking and recovery operations

**Platform Support:**
- **iOS/iPadOS**: Fully supported (iOS 9.3+, supervised devices)
- **Android Enterprise**: Supported (fully managed, dedicated, and work profile devices)
- **macOS**: Supported (macOS 10.13+, supervised devices or user-approved MDM)
- **Windows**: Limited support (Windows 10/11 with location services enabled)

**Location Data:**
- Location data is displayed in the Intune admin center under device properties
- Shows latitude, longitude, altitude (if available)
- Includes location accuracy radius
- Displays timestamp of when location was captured
- Location history may be available depending on platform

**Privacy Considerations:**
- Users may receive notification that device location was requested
- Location tracking should comply with organizational privacy policies
- Document legitimate business reasons for location requests
- Consider legal requirements in your jurisdiction

**Reference:** [Microsoft Graph API - Locate Device](https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-locatedevice?view=graph-rest-beta)

## Microsoft Documentation

### Graph API References
- [locateDevice action](https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-locatedevice?view=graph-rest-beta)
- [managedDevice resource type](https://learn.microsoft.com/en-us/graph/api/resources/intune-devices-manageddevice?view=graph-rest-beta)

### Intune Remote Actions Guides
- [Windows Remote Actions](https://learn.microsoft.com/en-us/intune/intune-service/remote-actions/?tabs=windows)
- [iOS/iPadOS Remote Actions](https://learn.microsoft.com/en-us/intune/intune-service/remote-actions/?tabs=ios-ipados)
- [macOS Remote Actions](https://learn.microsoft.com/en-us/intune/intune-service/remote-actions/?tabs=macos)
- [Android Remote Actions](https://learn.microsoft.com/en-us/intune/intune-service/remote-actions/?tabs=android)
- [ChromeOS Remote Actions](https://learn.microsoft.com/en-us/intune/intune-service/remote-actions/?tabs=chromeos)

## API Permissions

The following API permissions are required in order to use this action.

### Microsoft Graph

- **Application**: `DeviceManagementManagedDevices.PrivilegedOperations.All`
- **Delegated**: `DeviceManagementManagedDevices.PrivilegedOperations.All`

## Version History

| Version | Status | Notes |
|---------|--------|-------|
| v0.33.0-alpha | Experimental | Initial release |

## Notes

### Platform Compatibility

| Platform | Support | Requirements |
|----------|---------|--------------|
| **iOS** | ✅ Full Support | iOS 9.3+, supervised devices |
| **iPadOS** | ✅ Full Support | Supervised devices |
| **macOS** | ✅ Supported | macOS 10.13+, supervised or user-approved MDM |
| **Windows** | ⚠️ Limited | Location services must be enabled |
| **Android** | ✅ Supported | Fully managed, dedicated, or work profile |

### How Location Works

- Device must be online to receive command
- Location services must be enabled
- Device queries GPS, WiFi, or cellular location
- Location data reported back to Intune
- Data displayed in Intune admin center
- Includes latitude, longitude, accuracy radius
- Timestamp of when location was captured
- Multiple requests can track device movement

### Location Data Access

Location information is displayed in:
- Microsoft Intune admin center
- Device properties > Hardware section
- Shows coordinates, accuracy, and timestamp
- May include altitude if available
- Location history (platform dependent)

### Privacy Considerations

- Users may receive notification of location request
- Follow organizational privacy policies
- Document business justification
- Consider legal requirements in your jurisdiction
- Ensure device ownership agreements cover tracking
- Comply with employee privacy laws

### Common Use Cases

- Lost device recovery
- Theft investigation
- Asset tracking and inventory
- Security incident response
- Compliance verification
- Lost mode tracking over time
- Coordinate device recovery with users

## Example Usage

```terraform
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
```

<!-- action schema generated by tfplugindocs -->
## Schema

### Required

- `device_ids` (List of String) List of managed device IDs to locate. Each ID must be a valid GUID format. Multiple devices can be located in a single action. Example: `["12345678-1234-1234-1234-123456789abc", "87654321-4321-4321-4321-ba9876543210"]`

**Important:** Devices must be online and have location services enabled to respond to the locate request. Location data will be displayed in the Microsoft Intune admin center once the device reports its position.

### Optional

- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).

