---
page_title: "microsoft365_graph_beta_device_management_managed_device_disable_lost_mode Action - terraform-provider-microsoft365"
subcategory: "Device Management"

description: |-
  Disables lost mode on iOS/iPadOS managed devices using the /deviceManagement/managedDevices/{managedDeviceId}/disableLostMode and /deviceManagement/comanagedDevices/{managedDeviceId}/disableLostMode endpoints. This action removes the device from lost mode, allowing normal device operation to resume. Lost mode is a feature that helps locate and secure lost iOS/iPadOS devices by locking them and displaying a custom message with contact information on the lock screen. This action supports disabling lost mode on multiple devices in a single operation.
  Important Notes:
  Only applicable to iOS and iPadOS devices (iOS 9.3+)Device must currently be in lost modeDevice must be supervisedRequires device to be online to receive commandOnce disabled, device returns to normal operationThe custom lock screen message is removed
  Use Cases:
  Device has been recovered and needs to be returned to serviceLost mode was enabled in errorDevice location has been confirmed and no longer needs trackingUser has regained possession of their device
  Platform Support:
  iOS/iPadOS: Fully supported (iOS 9.3+, supervised devices only)Other Platforms: Not applicable - lost mode is iOS/iPadOS only
  Reference: Microsoft Graph API - Disable Lost Mode https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-disablelostmode?view=graph-rest-beta
---

# microsoft365_graph_beta_device_management_managed_device_disable_lost_mode (Action)

Disables lost mode on iOS/iPadOS managed devices using the `/deviceManagement/managedDevices/{managedDeviceId}/disableLostMode` and `/deviceManagement/comanagedDevices/{managedDeviceId}/disableLostMode` endpoints. This action removes the device from lost mode, allowing normal device operation to resume. Lost mode is a feature that helps locate and secure lost iOS/iPadOS devices by locking them and displaying a custom message with contact information on the lock screen. This action supports disabling lost mode on multiple devices in a single operation.

**Important Notes:**
- Only applicable to iOS and iPadOS devices (iOS 9.3+)
- Device must currently be in lost mode
- Device must be supervised
- Requires device to be online to receive command
- Once disabled, device returns to normal operation
- The custom lock screen message is removed

**Use Cases:**
- Device has been recovered and needs to be returned to service
- Lost mode was enabled in error
- Device location has been confirmed and no longer needs tracking
- User has regained possession of their device

**Platform Support:**
- **iOS/iPadOS**: Fully supported (iOS 9.3+, supervised devices only)
- **Other Platforms**: Not applicable - lost mode is iOS/iPadOS only

**Reference:** [Microsoft Graph API - Disable Lost Mode](https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-disablelostmode?view=graph-rest-beta)

## Microsoft Documentation

### Graph API References
- [disableLostMode action](https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-disablelostmode?view=graph-rest-beta)
- [managedDevice resource type](https://learn.microsoft.com/en-us/graph/api/resources/intune-devices-manageddevice?view=graph-rest-beta)

### Intune Remote Actions Guides
- [Lost mode for iOS devices](https://learn.microsoft.com/en-us/intune/intune-service/remote-actions/device-lost-mode?pivots=ios)
- [Lost mode for ChromeOS devices](https://learn.microsoft.com/en-us/intune/intune-service/remote-actions/device-lost-mode?pivots=chromeos)

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
| **iOS** | ✅ Full Support | iOS 9.3+, supervised devices only |
| **iPadOS** | ✅ Full Support | Supervised devices only |
| **macOS** | ❌ Not Supported | Lost mode iOS/iPadOS only |
| **Windows** | ❌ Not Supported | Lost mode iOS/iPadOS only |
| **Android** | ❌ Not Supported | Lost mode iOS/iPadOS only |

### What is Lost Mode?

Lost mode is an Apple security feature that:
- Locks device and displays custom message with contact info
- Enables device location tracking
- Prevents unauthorized access to device data
- Helps recover lost or stolen devices
- Automatically enabled via Intune or Find My app

### When to Disable Lost Mode

- Device has been physically recovered
- Device location confirmed and device is safe
- User has regained possession of their device
- Lost mode was enabled in error
- Device is being returned to service
- Device recovery operation completed

### What Happens When Disabled

- Device returns to normal operation
- Custom lock screen message removed
- Device can be unlocked with regular passcode
- Location tracking associated with lost mode stops
- Device becomes fully functional again
- User can access all device features

## Example Usage

```terraform
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
```

<!-- action schema generated by tfplugindocs -->
## Schema

### Optional

- `comanaged_device_ids` (List of String) List of co-managed device IDs to disable lost mode for. These are iOS/iPadOS devices managed by both Intune and Configuration Manager (SCCM). Each ID must be a valid GUID format. Example: `["12345678-1234-1234-1234-123456789abc"]`

**Note:** At least one of `managed_device_ids` or `comanaged_device_ids` must be provided.
- `ignore_partial_failures` (Boolean) If set to `true`, the action will succeed even if some operations fail. Failed operations will be reported as warnings instead of errors. Default: `false` (action fails if any operation fails).
- `managed_device_ids` (List of String) List of managed device IDs to disable lost mode for. These are iOS/iPadOS devices fully managed by Intune only. Each ID must be a valid GUID format. Example: `["12345678-1234-1234-1234-123456789abc", "87654321-4321-4321-4321-ba9876543210"]`

**Note:** At least one of `managed_device_ids` or `comanaged_device_ids` must be provided. You can provide both to disable lost mode on different types of devices in one action.
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))
- `validate_device_exists` (Boolean) Whether to validate that devices exist and are iOS/iPadOS devices before attempting to disable lost mode. Disabling this can speed up planning but may result in runtime errors for non-existent or unsupported devices. Default: `true`.

<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `invoke` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).

