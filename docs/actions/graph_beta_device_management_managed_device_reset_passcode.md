---
page_title: "microsoft365_graph_beta_device_management_managed_device_reset_passcode Action - terraform-provider-microsoft365"
subcategory: "Device Management"

description: |-
  Resets the passcode on managed devices using the /deviceManagement/managedDevices/{managedDeviceId}/resetPasscode endpoint. This action removes the current device passcode/password and generates a new temporary passcode. The new passcode is displayed to the administrator and must be communicated to the device user. This action supports resetting passcodes on multiple devices in a single operation.
  Important Notes:
  The device must be online and able to receive the commandOn iOS/iPadOS devices, the device must be supervisedOn Android devices, this removes the passcode requirement temporarilyOn Windows devices, the functionality varies by Windows versionThe new passcode is a temporary system-generated code that should be changed by the userThis action requires the device to be enrolled and actively managed by Intune
  Use Cases:
  User forgot device passcode and cannot unlock deviceDevice locked after too many failed passcode attemptsAdministrative access needed for troubleshootingSecurity incident requiring immediate access to device
  Reference: Microsoft Graph API - Reset Passcode https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-resetpasscode?view=graph-rest-beta
---

# microsoft365_graph_beta_device_management_managed_device_reset_passcode (Action)

Resets the passcode on managed devices using the `/deviceManagement/managedDevices/{managedDeviceId}/resetPasscode` endpoint. This action removes the current device passcode/password and generates a new temporary passcode. The new passcode is displayed to the administrator and must be communicated to the device user. This action supports resetting passcodes on multiple devices in a single operation.

**Important Notes:**
- The device must be online and able to receive the command
- On iOS/iPadOS devices, the device must be supervised
- On Android devices, this removes the passcode requirement temporarily
- On Windows devices, the functionality varies by Windows version
- The new passcode is a temporary system-generated code that should be changed by the user
- This action requires the device to be enrolled and actively managed by Intune

**Use Cases:**
- User forgot device passcode and cannot unlock device
- Device locked after too many failed passcode attempts
- Administrative access needed for troubleshooting
- Security incident requiring immediate access to device

**Reference:** [Microsoft Graph API - Reset Passcode](https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-resetpasscode?view=graph-rest-beta)

## Microsoft Documentation

### Graph API References
- [resetPasscode action](https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-resetpasscode?view=graph-rest-beta)
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
| **iOS** | ✅ Full Support | Supervised devices only |
| **iPadOS** | ✅ Full Support | Supervised devices only |
| **macOS** | ⚠️ Limited | Supervised or DEP enrolled devices |
| **Windows** | ⚠️ Limited | Azure AD joined devices |
| **Android** | ⚠️ Limited | Work profile or fully managed devices |

### Platform-Specific Requirements

#### iOS/iPadOS
- Device must be supervised
- Enrollment via DEP/ABM or Apple Configurator
- User will be prompted to set new passcode after unlock

#### Android
- Device must be fully managed or have work profile
- May only reset work profile passcode on BYOD
- Device may need to be online

#### Windows
- Device must be Azure AD joined
- May require specific Windows version
- BitLocker recovery key access may be needed

#### macOS
- Best support with supervised or DEP enrolled devices
- May require user to be present
- FileVault considerations apply

### How Passcode Reset Works

1. Reset command issued via this action
2. Intune generates temporary passcode
3. Passcode displayed in Intune admin portal
4. User enters temporary passcode to unlock device
5. User prompted to create new permanent passcode
6. Temporary passcode expires after first use

### Retrieving New Passcode

- Navigate to Intune admin center
- Go to Devices > All devices
- Select the device
- New passcode displayed in device details
- Passcode typically 6-8 digits/characters
- Securely communicate passcode to user

### Common Use Cases

- User forgot device passcode
- Locked out device recovery
- Security incident response
- Departing employee device recovery
- Compliance enforcement
- Lost device preparation for return
- Device provisioning for new user
- Emergency access requirements

### Best Practices

- Verify user identity before providing new passcode
- Communicate passcode securely (not via email)
- Document reason for passcode reset
- Follow up to ensure user sets permanent passcode
- Consider privacy and compliance requirements
- Use in conjunction with other security measures
- Monitor for repeated reset requests
- Train help desk on passcode retrieval

## Example Usage

```terraform
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
```

<!-- action schema generated by tfplugindocs -->
## Schema

### Required

- `device_ids` (List of String) List of managed device IDs to reset passcodes for. Each ID must be a valid GUID format. Multiple device passcodes can be reset in a single action. Example: `["12345678-1234-1234-1234-123456789abc", "87654321-4321-4321-4321-ba9876543210"]`

**Important:** The new temporary passcode for each device will be displayed in Intune and must be communicated to the device user. Users should change this temporary passcode immediately after regaining access.

### Optional

- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).

