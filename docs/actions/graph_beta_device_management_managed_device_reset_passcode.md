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
- [Device passcode reset](https://learn.microsoft.com/en-us/intune/intune-service/remote-actions/device-passcode-reset)

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
| **Android** | ✅ Full Support | Work profile or fully managed devices |
| **iOS** | ❌ Not Supported | Not available for iOS devices |
| **iPadOS** | ❌ Not Supported | Not available for iPadOS devices |
| **macOS** | ❌ Not Supported | Not available for macOS devices |
| **Windows** | ❌ Not Supported | Not available for Windows devices |
| **ChromeOS** | ❌ Not Supported | Not available for ChromeOS devices |

### Platform-Specific Requirements

#### Android
- Device must be fully managed or have work profile
- May only reset work profile passcode on BYOD devices
- Device needs to be online to receive the command
- Applies to screen lock PIN, password, or pattern
- User will be prompted to set new passcode after reset

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
# Example 1: Reset passcode for a single device - Minimal
action "microsoft365_graph_beta_device_management_managed_device_reset_passcode" "reset_single" {
  config {
    device_ids = [
      "12345678-1234-1234-1234-123456789abc"
    ]
  }
}

# Example 2: Reset passcodes for multiple devices
action "microsoft365_graph_beta_device_management_managed_device_reset_passcode" "reset_batch" {
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

# Example 3: Reset passcodes with validation and failure handling - Maximal
action "microsoft365_graph_beta_device_management_managed_device_reset_passcode" "reset_with_validation" {
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

# Example 4: Reset passcodes for locked devices from data source
data "microsoft365_graph_beta_device_management_managed_device" "locked_ios_devices" {
  filter_type  = "odata"
  odata_filter = "(operatingSystem eq 'iOS') and (isSupervised eq true)"
}

action "microsoft365_graph_beta_device_management_managed_device_reset_passcode" "reset_locked_devices" {
  config {
    device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.locked_ios_devices.items : device.id]

    validate_device_exists = true

    timeouts = {
      invoke = "15m"
    }
  }
}

# Example 5: Emergency passcode reset for specific user's devices
data "microsoft365_graph_beta_device_management_managed_device" "user_devices" {
  filter_type  = "odata"
  odata_filter = "userPrincipalName eq 'user@example.com'"
}

action "microsoft365_graph_beta_device_management_managed_device_reset_passcode" "reset_user_devices" {
  config {
    device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.user_devices.items : device.id]

    ignore_partial_failures = true

    timeouts = {
      invoke = "10m"
    }
  }
}

# Example 6: Reset passcodes for Android devices
data "microsoft365_graph_beta_device_management_managed_device" "android_devices" {
  filter_type  = "odata"
  odata_filter = "operatingSystem eq 'Android'"
}

action "microsoft365_graph_beta_device_management_managed_device_reset_passcode" "reset_android" {
  config {
    device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.android_devices.items : device.id]

    timeouts = {
      invoke = "15m"
    }
  }
}

# Output examples
output "reset_device_count" {
  value       = length(action.microsoft365_graph_beta_device_management_managed_device_reset_passcode.reset_batch.config.device_ids)
  description = "Number of devices for which passcodes were reset"
}

output "locked_devices_reset_count" {
  value       = length(action.microsoft365_graph_beta_device_management_managed_device_reset_passcode.reset_locked_devices.config.device_ids)
  description = "Number of locked devices for which passcodes were reset"
}
```

<!-- action schema generated by tfplugindocs -->
## Schema

### Required

- `device_ids` (List of String) List of managed device IDs to reset passcodes for. Each ID must be a valid GUID format. Multiple device passcodes can be reset in a single action. Example: `["12345678-1234-1234-1234-123456789abc", "87654321-4321-4321-4321-ba9876543210"]`

**Important:** The new temporary passcode for each device will be displayed in Intune and must be communicated to the device user. Users should change this temporary passcode immediately after regaining access.

### Optional

- `ignore_partial_failures` (Boolean) If set to `true`, the action will succeed even if some operations fail. Failed operations will be reported as warnings instead of errors. Default: `false` (action fails if any operation fails).
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))
- `validate_device_exists` (Boolean) Whether to validate that devices exist and are Android devices before attempting passcode reset. Disabling this can speed up planning but may result in runtime errors for non-existent or unsupported devices. Default: `true`.

<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `invoke` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).

