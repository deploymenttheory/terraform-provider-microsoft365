---
page_title: "microsoft365_graph_beta_device_management_managed_device_play_lost_mode_sound Action - terraform-provider-microsoft365"
subcategory: "Device Management"

description: |-
  Plays a sound on iOS/iPadOS managed devices in lost mode using the /deviceManagement/managedDevices/{managedDeviceId}/playLostModeSound and /deviceManagement/comanagedDevices/{managedDeviceId}/playLostModeSound endpoints. This action helps locate lost devices by triggering an audible alert that plays even if the device is in silent mode. The sound plays for a specified duration to assist in physically locating the device. This action supports playing sounds on multiple devices in a single operation with per-device duration settings.
  Important Notes:
  Only applicable to iOS and iPadOS devices in lost modeDevice must be supervisedDevice must currently be in lost modeSound plays even if device is in silent modeRequires device to be online to receive commandEach device can have its own sound duration
  Use Cases:
  Device is nearby but cannot be visually locatedDevice is in lost mode and needs audible alertAssisting user in finding device in office or homeConfirming device location before recovery
  Platform Support:
  iOS/iPadOS: Fully supported (supervised devices in lost mode only)Other Platforms: Not applicable - lost mode is iOS/iPadOS only
  Reference: Microsoft Graph API - Play Lost Mode Sound https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-playlostmodesound?view=graph-rest-beta
---

# microsoft365_graph_beta_device_management_managed_device_play_lost_mode_sound (Action)

Plays a sound on iOS/iPadOS managed devices in lost mode using the `/deviceManagement/managedDevices/{managedDeviceId}/playLostModeSound` and `/deviceManagement/comanagedDevices/{managedDeviceId}/playLostModeSound` endpoints. This action helps locate lost devices by triggering an audible alert that plays even if the device is in silent mode. The sound plays for a specified duration to assist in physically locating the device. This action supports playing sounds on multiple devices in a single operation with per-device duration settings.

**Important Notes:**
- Only applicable to iOS and iPadOS devices in lost mode
- Device must be supervised
- Device must currently be in lost mode
- Sound plays even if device is in silent mode
- Requires device to be online to receive command
- Each device can have its own sound duration

**Use Cases:**
- Device is nearby but cannot be visually located
- Device is in lost mode and needs audible alert
- Assisting user in finding device in office or home
- Confirming device location before recovery

**Platform Support:**
- **iOS/iPadOS**: Fully supported (supervised devices in lost mode only)
- **Other Platforms**: Not applicable - lost mode is iOS/iPadOS only

**Reference:** [Microsoft Graph API - Play Lost Mode Sound](https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-playlostmodesound?view=graph-rest-beta)

## Microsoft Documentation

### Graph API References
- [playLostModeSound action](https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-playlostmodesound?view=graph-rest-beta)
- [managedDevice resource type](https://learn.microsoft.com/en-us/graph/api/resources/intune-devices-manageddevice?view=graph-rest-beta)

### Intune Remote Actions Guides
- [Play lost mode sound - iOS](https://learn.microsoft.com/en-us/intune/intune-service/remote-actions/device-play-lost-mode-sound?pivots=ios)
- [Play lost mode sound - Android](https://learn.microsoft.com/en-us/intune/intune-service/remote-actions/device-play-lost-mode-sound?pivots=android)

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
| **iOS** | ✅ Full Support | Supervised devices in lost mode only |
| **iPadOS** | ✅ Full Support | Supervised devices in lost mode only |
| **macOS** | ❌ Not Supported | Lost mode is iOS/iPadOS only |
| **Windows** | ❌ Not Supported | Lost mode is iOS/iPadOS only |
| **Android** | ❌ Not Supported | Lost mode is iOS/iPadOS only |

### What is Play Lost Mode Sound?

Play Lost Mode Sound is a feature that:
- Plays an audible alert on iOS/iPadOS devices
- Works only when device is in lost mode
- Sound plays even if device is in silent mode
- Helps physically locate devices that are nearby
- Sound duration can be customized per device
- Useful for finding devices in known general locations

### When to Play Lost Mode Sound

- Device is believed to be nearby but cannot be visually located
- Confirming device location before physical recovery
- Assisting user in finding device in office, home, or vehicle
- Device is in lost mode and approximate location is known
- Need audible confirmation of device presence in area
- Supplementing GPS location tracking with audio cues

### What Happens When Sound is Played

- Device immediately begins playing an audible alert
- Sound plays regardless of device silent mode setting
- Alert continues for the specified duration
- Helps pinpoint exact physical location of device
- Lost mode status remains unchanged
- Action can be repeated as needed
- No visual indication on device screen beyond lost mode message

## Example Usage

```terraform
# Example 1: Play lost mode sound on a single device with default duration
action "microsoft365_graph_beta_device_management_managed_device_play_lost_mode_sound" "play_sound_single" {

  managed_devices {
    device_id = "12345678-1234-1234-1234-123456789abc"
  }

  timeouts = {
    invoke = "5m"
  }
}

# Example 2: Play lost mode sound with specific duration
action "microsoft365_graph_beta_device_management_managed_device_play_lost_mode_sound" "play_sound_with_duration" {

  managed_devices {
    device_id          = "12345678-1234-1234-1234-123456789abc"
    duration_in_minutes = "5"
  }

  timeouts = {
    invoke = "5m"
  }
}

# Example 3: Play sound on multiple devices with different durations
action "microsoft365_graph_beta_device_management_managed_device_play_lost_mode_sound" "play_sound_multiple" {

  managed_devices {
    device_id          = "12345678-1234-1234-1234-123456789abc"
    duration_in_minutes = "3"
  }

  managed_devices {
    device_id          = "87654321-4321-4321-4321-ba9876543210"
    duration_in_minutes = "10"
  }

  timeouts = {
    invoke = "10m"
  }
}

# Example 4: Play sound on all devices currently in lost mode
data "microsoft365_graph_beta_device_management_managed_device" "devices_in_lost_mode" {
  filter_type  = "odata"
  odata_filter = "lostModeState ne 'disabled' and operatingSystem eq 'iOS'"
}

action "microsoft365_graph_beta_device_management_managed_device_play_lost_mode_sound" "play_sound_all_lost_mode" {

  dynamic "managed_devices" {
    for_each = data.microsoft365_graph_beta_device_management_managed_device.devices_in_lost_mode.items
    content {
      device_id          = managed_devices.value.id
      duration_in_minutes = "5"
    }
  }

  timeouts = {
    invoke = "15m"
  }
}

# Example 5: Play sound for specific user's devices in lost mode
data "microsoft365_graph_beta_device_management_managed_device" "user_lost_devices" {
  filter_type  = "odata"
  odata_filter = "userId eq 'user@example.com' and lostModeState ne 'disabled'"
}

action "microsoft365_graph_beta_device_management_managed_device_play_lost_mode_sound" "play_sound_user_devices" {

  dynamic "managed_devices" {
    for_each = data.microsoft365_graph_beta_device_management_managed_device.user_lost_devices.items
    content {
      device_id          = managed_devices.value.id
      duration_in_minutes = "3"
    }
  }

  timeouts = {
    invoke = "10m"
  }
}

# Example 6: Play sound on co-managed device
action "microsoft365_graph_beta_device_management_managed_device_play_lost_mode_sound" "play_sound_comanaged" {

  comanaged_devices {
    device_id          = "abcdef12-3456-7890-abcd-ef1234567890"
    duration_in_minutes = "5"
  }

  timeouts = {
    invoke = "5m"
  }
}

# Example 7: Play sound to locate device nearby
action "microsoft365_graph_beta_device_management_managed_device_play_lost_mode_sound" "locate_nearby_device" {

  managed_devices {
    device_id          = "12345678-1234-1234-1234-123456789abc"
    duration_in_minutes = "2"
  }

  timeouts = {
    invoke = "5m"
  }
}

# Output examples
output "devices_with_sound" {
  value       = length(action.play_sound_multiple.managed_devices)
  description = "Number of devices that had lost mode sound played"
}

output "lost_mode_devices_count" {
  value       = length(action.play_sound_all_lost_mode.managed_devices)
  description = "Number of devices in lost mode that received sound command"
}
```

<!-- action schema generated by tfplugindocs -->
## Schema

### Optional

- `comanaged_devices` (Block List) List of co-managed devices to play lost mode sound on. These are iOS/iPadOS devices managed by both Intune and Configuration Manager (SCCM). Each entry specifies a device ID and the duration to play the sound.

**Note:** At least one of `managed_devices` or `comanaged_devices` must be provided. (see [below for nested schema](#nestedblock--comanaged_devices))
- `managed_devices` (Block List) List of managed devices to play lost mode sound on. These are iOS/iPadOS devices fully managed by Intune only. Each entry specifies a device ID and the duration to play the sound.

**Note:** At least one of `managed_devices` or `comanaged_devices` must be provided. You can provide both to play sounds on different types of devices in one action. (see [below for nested schema](#nestedblock--managed_devices))
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

<a id="nestedblock--comanaged_devices"></a>
### Nested Schema for `comanaged_devices`

Required:

- `device_id` (String) The unique identifier (GUID) of the co-managed device to play sound on. Example: `"12345678-1234-1234-1234-123456789abc"`

Optional:

- `duration_in_minutes` (String) The duration in minutes to play the lost mode sound. Example: `"5"`


<a id="nestedblock--managed_devices"></a>
### Nested Schema for `managed_devices`

Required:

- `device_id` (String) The unique identifier (GUID) of the managed device to play sound on. Example: `"12345678-1234-1234-1234-123456789abc"`

Optional:

- `duration_in_minutes` (String) The duration in minutes to play the lost mode sound. If not specified, the sound will play for the default duration. Example: `"5"` for 5 minutes


<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).


