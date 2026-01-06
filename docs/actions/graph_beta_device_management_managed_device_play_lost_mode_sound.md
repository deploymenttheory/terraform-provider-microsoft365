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
| v0.40.0-alpha | Experimental | Example fixes and refactored sync progress logic |


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
# Example 1: Play lost mode sound for a single device - Minimal
action "microsoft365_graph_beta_device_management_managed_device_play_lost_mode_sound" "play_sound_single" {
  config {
    managed_devices = [
      {
        device_id = "12345678-1234-1234-1234-123456789abc"
      }
    ]
  }
}

# Example 2: Play lost mode sound for multiple devices
action "microsoft365_graph_beta_device_management_managed_device_play_lost_mode_sound" "play_sound_multiple" {
  config {
    managed_devices = [
      {
        device_id = "12345678-1234-1234-1234-123456789abc"
      },
      {
        device_id = "87654321-4321-4321-4321-ba9876543210"
      },
      {
        device_id = "abcdef12-3456-7890-abcd-ef1234567890"
      }
    ]

    timeouts = {
      invoke = "10m"
    }
  }
}

# Example 3: Maximal configuration with validation
action "microsoft365_graph_beta_device_management_managed_device_play_lost_mode_sound" "play_sound_maximal" {
  config {
    managed_devices = [
      {
        device_id = "12345678-1234-1234-1234-123456789abc"
      },
      {
        device_id = "87654321-4321-4321-4321-ba9876543210"
      }
    ]

    comanaged_devices = [
      {
        device_id = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
      }
    ]

    ignore_partial_failures = true
    validate_device_exists  = true

    timeouts = {
      invoke = "5m"
    }
  }
}

# Example 4: Play sound for all devices in lost mode
data "microsoft365_graph_beta_device_management_managed_device" "devices_in_lost_mode" {
  filter_type  = "odata"
  odata_filter = "lostModeState eq 'enabled'"
}

action "microsoft365_graph_beta_device_management_managed_device_play_lost_mode_sound" "play_sound_all_lost_mode" {
  config {
    managed_devices = [
      for device in data.microsoft365_graph_beta_device_management_managed_device.devices_in_lost_mode.items : {
        device_id = device.id
      }
    ]

    validate_device_exists = true

    timeouts = {
      invoke = "15m"
    }
  }
}

# Example 5: Play sound for user's devices in lost mode
data "microsoft365_graph_beta_device_management_managed_device" "user_lost_devices" {
  filter_type  = "odata"
  odata_filter = "(userPrincipalName eq 'user@example.com') and (lostModeState eq 'enabled')"
}

action "microsoft365_graph_beta_device_management_managed_device_play_lost_mode_sound" "play_sound_user_devices" {
  config {
    managed_devices = [
      for device in data.microsoft365_graph_beta_device_management_managed_device.user_lost_devices.items : {
        device_id = device.id
      }
    ]

    timeouts = {
      invoke = "5m"
    }
  }
}

# Example 6: Play sound for co-managed device
action "microsoft365_graph_beta_device_management_managed_device_play_lost_mode_sound" "play_sound_comanaged" {
  config {
    comanaged_devices = [
      {
        device_id = "abcdef12-3456-7890-abcd-ef1234567890"
      }
    ]

    timeouts = {
      invoke = "5m"
    }
  }
}
```

<!-- action schema generated by tfplugindocs -->
## Schema

### Optional

- `comanaged_devices` (Attributes List) List of co-managed devices to play lost mode sound on. These are iOS/iPadOS devices managed by both Intune and Configuration Manager (SCCM). Each entry specifies a device ID and the duration to play the sound.

**Note:** At least one of `managed_devices` or `comanaged_devices` must be provided. (see [below for nested schema](#nestedatt--comanaged_devices))
- `ignore_partial_failures` (Boolean) If set to `true`, the action will succeed even if some operations fail. Failed operations will be reported as warnings instead of errors. Default: `false` (action fails if any operation fails).
- `managed_devices` (Attributes List) List of managed devices to play lost mode sound on. These are iOS/iPadOS devices fully managed by Intune only. Each entry specifies a device ID and the duration to play the sound.

**Note:** At least one of `managed_devices` or `comanaged_devices` must be provided. You can provide both to play sounds on different types of devices in one action. (see [below for nested schema](#nestedatt--managed_devices))
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))
- `validate_device_exists` (Boolean) Whether to validate that devices exist, are iOS/iPadOS devices, are supervised, and are in lost mode before attempting to play the sound. Disabling this can speed up planning but may result in runtime errors for non-existent, unsupported, or devices not in lost mode. Default: `true`.

<a id="nestedatt--comanaged_devices"></a>
### Nested Schema for `comanaged_devices`

Required:

- `device_id` (String) The unique identifier (GUID) of the co-managed device to play sound on. Example: `"12345678-1234-1234-1234-123456789abc"`

Optional:

- `duration_in_minutes` (String) The duration in minutes to play the lost mode sound. Example: `"5"`


<a id="nestedatt--managed_devices"></a>
### Nested Schema for `managed_devices`

Required:

- `device_id` (String) The unique identifier (GUID) of the managed device to play sound on. Example: `"12345678-1234-1234-1234-123456789abc"`

Optional:

- `duration_in_minutes` (String) The duration in minutes to play the lost mode sound. If not specified, the sound will play for the default duration. Example: `"5"` for 5 minutes


<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `invoke` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).


