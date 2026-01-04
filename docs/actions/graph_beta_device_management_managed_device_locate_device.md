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
- [Device locate - Windows](https://learn.microsoft.com/en-us/intune/intune-service/remote-actions/device-locate?pivots=windows)
- [Device locate - iOS](https://learn.microsoft.com/en-us/intune/intune-service/remote-actions/device-locate?pivots=ios)
- [Device locate - Android](https://learn.microsoft.com/en-us/intune/intune-service/remote-actions/device-locate?pivots=android)

## API Permissions

The following API permissions are required in order to use this action.

### Microsoft Graph

- **Application**: `DeviceManagementManagedDevices.PrivilegedOperations.All`
- **Delegated**: `DeviceManagementManagedDevices.PrivilegedOperations.All`

## Version History

| Version | Status | Notes |
|---------|--------|-------|
| v0.33.0-alpha | Experimental | Initial release |
| v0.40.0-alpha | Experimental | Example fixes and refactored sync progress logic |


## Notes

### Platform Compatibility

| Platform | Support | Requirements |
|----------|---------|--------------|
| **Windows** | ⚠️ Limited | Location services must be enabled |
| **iOS** | ✅ Full Support | iOS 9.3+, supervised devices |
| **iPadOS** | ✅ Full Support | Supervised devices |
| **Android** | ✅ Supported | Fully managed, dedicated, or work profile |
| **macOS** | ❌ Not Supported | Not available on macOS devices |
| **ChromeOS** | ❌ Not Supported | Not available on ChromeOS devices |

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
```

<!-- action schema generated by tfplugindocs -->
## Schema

### Required

- `device_ids` (List of String) List of managed device IDs to locate. Each ID must be a valid GUID format. Multiple devices can be located in a single action. Example: `["12345678-1234-1234-1234-123456789abc", "87654321-4321-4321-4321-ba9876543210"]`

**Important:** Devices must be online and have location services enabled to respond to the locate request. Location data will be displayed in the Microsoft Intune admin center once the device reports its position.

### Optional

- `ignore_partial_failures` (Boolean) If set to `true`, the action will succeed even if some operations fail. Failed operations will be reported as warnings instead of errors. Default: `false` (action fails if any operation fails).
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))
- `validate_device_exists` (Boolean) Whether to validate that devices exist and support location services before attempting to locate them. Disabling this can speed up planning but may result in runtime errors for non-existent or unsupported devices. Default: `true`.

<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `invoke` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).

