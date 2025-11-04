---
page_title: "microsoft365_graph_beta_device_management_managed_device_enable_lost_mode Action - terraform-provider-microsoft365"
subcategory: "Device Management"

description: |-
  Enables lost mode on iOS/iPadOS managed devices using the /deviceManagement/managedDevices/{managedDeviceId}/enableLostMode and /deviceManagement/comanagedDevices/{managedDeviceId}/enableLostMode endpoints. This action locks the device and displays a custom message with contact information on the lock screen. Lost mode is a feature that helps locate and secure lost iOS/iPadOS devices by locking them and enabling device location tracking. This action supports enabling lost mode on multiple devices in a single operation with per-device messages.
  Important Notes:
  Only applicable to iOS and iPadOS devices (iOS 9.3+)Device must be supervisedRequires device to be online to receive commandLocks device and displays custom message with contact informationEnables device location trackingEach device can have its own custom message, phone number, and footnote
  Use Cases:
  Device has been reported lost or stolenNeed to lock device and display recovery contact informationNeed to track device location for recoveryPrevent unauthorized access to corporate data
  Platform Support:
  iOS/iPadOS: Fully supported (iOS 9.3+, supervised devices only)Other Platforms: Not applicable - lost mode is iOS/iPadOS only
  Reference: Microsoft Graph API - Enable Lost Mode https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-enablelostmode?view=graph-rest-beta
---

# microsoft365_graph_beta_device_management_managed_device_enable_lost_mode (Action)

Enables lost mode on iOS/iPadOS managed devices using the `/deviceManagement/managedDevices/{managedDeviceId}/enableLostMode` and `/deviceManagement/comanagedDevices/{managedDeviceId}/enableLostMode` endpoints. This action locks the device and displays a custom message with contact information on the lock screen. Lost mode is a feature that helps locate and secure lost iOS/iPadOS devices by locking them and enabling device location tracking. This action supports enabling lost mode on multiple devices in a single operation with per-device messages.

**Important Notes:**
- Only applicable to iOS and iPadOS devices (iOS 9.3+)
- Device must be supervised
- Requires device to be online to receive command
- Locks device and displays custom message with contact information
- Enables device location tracking
- Each device can have its own custom message, phone number, and footnote

**Use Cases:**
- Device has been reported lost or stolen
- Need to lock device and display recovery contact information
- Need to track device location for recovery
- Prevent unauthorized access to corporate data

**Platform Support:**
- **iOS/iPadOS**: Fully supported (iOS 9.3+, supervised devices only)
- **Other Platforms**: Not applicable - lost mode is iOS/iPadOS only

**Reference:** [Microsoft Graph API - Enable Lost Mode](https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-enablelostmode?view=graph-rest-beta)

## Microsoft Documentation

### Graph API References
- [enableLostMode action](https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-enablelostmode?view=graph-rest-beta)
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
- Can be enabled via Intune or Find My app

### When to Enable Lost Mode

- Device has been reported lost or stolen
- Need to remotely lock and secure device
- Need to display recovery contact information
- Need to track device location for recovery
- Prevent unauthorized access to corporate data
- User cannot be reached to secure device manually

### What Happens When Enabled

- Device is immediately locked
- Custom message displayed on lock screen with contact info
- Location tracking is enabled
- Device cannot be unlocked without proper credentials
- Device data remains encrypted and protected
- User must contact provided phone number to recover device

### Message Best Practices

- Include clear instructions for returning the device
- Provide a contact phone number that will be monitored
- Keep message concise and professional
- Include organization name or identification
- Avoid including sensitive information
- Consider including reward information if applicable

## Example Usage

```terraform
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
```

<!-- action schema generated by tfplugindocs -->
## Schema

### Optional

- `comanaged_devices` (Block List) List of co-managed devices to enable lost mode for. These are iOS/iPadOS devices managed by both Intune and Configuration Manager (SCCM). Each entry specifies a device ID and the custom lost mode configuration.

**Note:** At least one of `managed_devices` or `comanaged_devices` must be provided. (see [below for nested schema](#nestedblock--comanaged_devices))
- `managed_devices` (Block List) List of managed devices to enable lost mode for. These are iOS/iPadOS devices fully managed by Intune only. Each entry specifies a device ID and the custom lost mode configuration for that device.

**Note:** At least one of `managed_devices` or `comanaged_devices` must be provided. You can provide both to enable lost mode on different types of devices in one action. (see [below for nested schema](#nestedblock--managed_devices))
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

<a id="nestedblock--comanaged_devices"></a>
### Nested Schema for `comanaged_devices`

Required:

- `device_id` (String) The unique identifier (GUID) of the co-managed device to enable lost mode for. Example: `"12345678-1234-1234-1234-123456789abc"`
- `message` (String) The message to display on this device's lock screen. This message should provide information on how to return the device. Example: `"This device has been lost. Please contact IT at 555-0123 to return."`
- `phone_number` (String) The phone number to display on this device's lock screen. Example: `"555-0123"`

Optional:

- `footer` (String) An optional footer to display below the message on this device's lock screen.


<a id="nestedblock--managed_devices"></a>
### Nested Schema for `managed_devices`

Required:

- `device_id` (String) The unique identifier (GUID) of the managed device to enable lost mode for. Example: `"12345678-1234-1234-1234-123456789abc"`
- `message` (String) The message to display on this device's lock screen. This message should provide information on how to return the device. Example: `"This device has been lost. Please contact IT at 555-0123 to return."`

**Requirements:**
- Must not be empty
- Should include clear instructions for device return
- Recommended: Include contact information and identification details
- `phone_number` (String) The phone number to display on this device's lock screen. This should be a contact number for returning the device. Example: `"555-0123"` or `"+1-555-0123"`

**Requirements:**
- Must not be empty
- Should be a valid phone number format
- Can include international dialing codes

Optional:

- `footnote` (String) An optional footnote to display below the message on this device's lock screen. This can be used for additional instructions or legal information. Example: `"Property of Contoso Corporation"`


<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).


