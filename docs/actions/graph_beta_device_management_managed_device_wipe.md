---
page_title: "microsoft365_graph_beta_device_management_managed_device_wipe Action - terraform-provider-microsoft365"
subcategory: "Device Management"

description: |-
  Wipes managed devices in Microsoft Intune using the /deviceManagement/managedDevices/{managedDeviceId}/wipe endpoint. This action is used to perform a factory reset, removing all data (company and personal) from the device. The device is returned to its out-of-box state and removed from Intune management. This action supports wiping multiple devices in a single operation.
  Important Notes:
  This action removes ALL data from the device unless keep_user_data is set to trueFor iOS/iPadOS devices, Activation Lock must be disabled or unlock code providedFor Windows devices, you can use protected wipe to maintain UEFI-embedded licensesFor Android devices, factory reset protection must be disabledThis action cannot be reversed - all data will be permanently deleted
  Reference: Microsoft Graph API - Wipe Managed Device https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-wipe?view=graph-rest-beta
---

# microsoft365_graph_beta_device_management_managed_device_wipe (Action)

Wipes managed devices in Microsoft Intune using the `/deviceManagement/managedDevices/{managedDeviceId}/wipe` endpoint. This action is used to perform a factory reset, removing all data (company and personal) from the device. The device is returned to its out-of-box state and removed from Intune management. This action supports wiping multiple devices in a single operation.

**Important Notes:**
- This action removes **ALL** data from the device unless `keep_user_data` is set to `true`
- For iOS/iPadOS devices, Activation Lock must be disabled or unlock code provided
- For Windows devices, you can use protected wipe to maintain UEFI-embedded licenses
- For Android devices, factory reset protection must be disabled
- This action cannot be reversed - all data will be permanently deleted

**Reference:** [Microsoft Graph API - Wipe Managed Device](https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-wipe?view=graph-rest-beta)

## Microsoft Documentation

### Graph API References
- [wipe action](https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-wipe?view=graph-rest-beta)
- [managedDevice resource type](https://learn.microsoft.com/en-us/graph/api/resources/intune-devices-manageddevice?view=graph-rest-beta)

### Intune Remote Actions Guides
- [Wipe devices - Windows](https://learn.microsoft.com/en-us/intune/intune-service/remote-actions/device-wipe?pivots=windows)
- [Wipe devices - iOS](https://learn.microsoft.com/en-us/intune/intune-service/remote-actions/device-wipe?pivots=ios)
- [Wipe devices - macOS](https://learn.microsoft.com/en-us/intune/intune-service/remote-actions/device-wipe?pivots=macos)
- [Wipe devices - Android](https://learn.microsoft.com/en-us/intune/intune-service/remote-actions/device-wipe?pivots=android)
- [Wipe devices - ChromeOS](https://learn.microsoft.com/en-us/intune/intune-service/remote-actions/device-wipe?pivots=chromeos)

## Microsoft Graph API Permissions

The following client `application` permissions are needed in order to use this action:

**Required:**
- `DeviceManagementConfiguration.ReadWrite.All`
- `DeviceManagementManagedDevices.ReadWrite.All`
- `DeviceManagementManagedDevices.PrivilegedOperations.All`

**Optional:**
- `None` `[N/A]`

## Version History

| Version | Status | Notes |
|---------|--------|-------|
| v0.33.0-alpha | Experimental | Initial release |
| v0.40.0-alpha | Experimental | Example fixes and refactored sync progress logic |

## Example Usage

```terraform
# Example 1: Wipe a single device (factory reset, removes all data) - Minimal
action "microsoft365_graph_beta_device_management_managed_device_wipe" "wipe_single" {
  config {
    device_ids = [
      "12345678-1234-1234-1234-123456789abc"
    ]

    timeouts = {
      invoke = "5m"
    }
  }
}

# Example 2: Wipe multiple devices with validation and failure handling
action "microsoft365_graph_beta_device_management_managed_device_wipe" "wipe_batch" {
  config {
    device_ids = [
      "12345678-1234-1234-1234-123456789abc",
      "87654321-4321-4321-4321-ba9876543210",
      "abcdef12-3456-7890-abcd-ef1234567890"
    ]

    validate_device_exists  = true # Validate devices before wiping
    ignore_partial_failures = true # Continue if some wipes fail

    timeouts = {
      invoke = "10m"
    }
  }
}

# Example 3: Selective wipe - keep user data, remove only company data
action "microsoft365_graph_beta_device_management_managed_device_wipe" "wipe_company_data_only" {
  config {
    device_ids = [
      "12345678-1234-1234-1234-123456789abc"
    ]

    keep_user_data = true

    timeouts = {
      invoke = "5m"
    }
  }
}

# Example 4: Wipe with enrollment data preserved (for automatic re-enrollment)
action "microsoft365_graph_beta_device_management_managed_device_wipe" "wipe_keep_enrollment" {
  config {
    device_ids = [
      "12345678-1234-1234-1234-123456789abc"
    ]

    keep_enrollment_data = true

    timeouts = {
      invoke = "5m"
    }
  }
}

# Example 5: Wipe macOS device with Activation Lock
action "microsoft365_graph_beta_device_management_managed_device_wipe" "wipe_macos" {
  config {
    device_ids = [
      "12345678-1234-1234-1234-123456789abc"
    ]

    macos_unlock_code = "123456" # 6-digit PIN for Activation Lock bypass

    timeouts = {
      invoke = "5m"
    }
  }
}

# Example 6: Wipe macOS with obliteration behavior control
action "microsoft365_graph_beta_device_management_managed_device_wipe" "wipe_macos_always_obliterate" {
  config {
    device_ids = [
      "12345678-1234-1234-1234-123456789abc"
    ]

    obliteration_behavior = "always" # Always obliterate on T2+ Macs

    timeouts = {
      invoke = "5m"
    }
  }
}

# Example 7: Wipe Windows device with protected wipe (preserves UEFI licenses)
action "microsoft365_graph_beta_device_management_managed_device_wipe" "wipe_windows" {
  config {
    device_ids = [
      "12345678-1234-1234-1234-123456789abc"
    ]

    use_protected_wipe = true

    timeouts = {
      invoke = "5m"
    }
  }
}

# Example 8: Wipe devices with eSIM, preserving data plan
action "microsoft365_graph_beta_device_management_managed_device_wipe" "wipe_keep_esim" {
  config {
    device_ids = [
      "12345678-1234-1234-1234-123456789abc"
    ]

    persist_esim_data_plan = true

    timeouts = {
      invoke = "5m"
    }
  }
}

# Example 9: Comprehensive wipe with multiple options and robust validation
action "microsoft365_graph_beta_device_management_managed_device_wipe" "wipe_comprehensive" {
  config {
    device_ids = [
      "12345678-1234-1234-1234-123456789abc"
    ]

    keep_enrollment_data    = true
    keep_user_data          = true
    persist_esim_data_plan  = true
    obliteration_behavior   = "doNotObliterate"
    validate_device_exists  = true
    ignore_partial_failures = false

    timeouts = {
      invoke = "5m"
    }
  }
}

# Example 10: Wipe non-compliant devices from data source
data "microsoft365_graph_beta_device_management_managed_device" "non_compliant_devices" {
  filter_type  = "odata"
  odata_filter = "complianceState eq 'noncompliant'"
}

action "microsoft365_graph_beta_device_management_managed_device_wipe" "wipe_non_compliant_devices" {
  config {
    device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.non_compliant_devices.items : device.id]

    # Wipe but keep enrollment data for automatic re-enrollment after compliance
    keep_enrollment_data    = true
    validate_device_exists  = true
    ignore_partial_failures = true # Continue even if some wipes fail

    timeouts = {
      invoke = "15m"
    }
  }
}

# Example 11: Wipe old devices by OS version
data "microsoft365_graph_beta_device_management_managed_device" "old_devices" {
  filter_type  = "odata"
  odata_filter = "operatingSystem eq 'Windows' and osVersion startsWith '10.0'"
}

action "microsoft365_graph_beta_device_management_managed_device_wipe" "wipe_old_windows_devices" {
  config {
    device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.old_devices.items : device.id]

    use_protected_wipe = true # Preserve Windows licenses

    timeouts = {
      invoke = "20m"
    }
  }
}

# Output examples
output "wiped_device_count" {
  value       = length(action.microsoft365_graph_beta_device_management_managed_device_wipe.wipe_batch.config.device_ids)
  description = "Number of devices wiped in batch operation"
}

output "non_compliant_devices_to_wipe" {
  value       = length(action.microsoft365_graph_beta_device_management_managed_device_wipe.wipe_non_compliant_devices.config.device_ids)
  description = "Number of non-compliant devices being wiped"
}
```

<!-- action schema generated by tfplugindocs -->
## Schema

### Required

- `device_ids` (List of String) List of managed device IDs to wipe. Each ID must be a valid GUID format. Multiple devices can be wiped in a single action. Example: `["12345678-1234-1234-1234-123456789abc", "87654321-4321-4321-4321-ba9876543210"]`

### Optional

- `ignore_partial_failures` (Boolean) When set to `true`, the action will complete successfully even if some devices fail to wipe. When `false` (default), the action will fail if any device wipe fails. Use this flag when wiping multiple devices and you want the action to succeed even if some wipes fail.
- `keep_enrollment_data` (Boolean) If `true`, maintains enrollment state data during wipe. This allows the device to automatically re-enroll after being wiped. Defaults to `false`.
- `keep_user_data` (Boolean) If `true`, preserves user data during the wipe operation. Only company data and managed apps are removed. **Note:** Not supported on all device types. Defaults to `false`.
- `macos_unlock_code` (String) The 6-digit PIN required to unlock macOS devices with Activation Lock enabled. Required for supervised macOS devices with Activation Lock. Format: 6-digit numeric string.
- `obliteration_behavior` (String) Specifies the obliteration behavior for macOS 12+ devices with Apple M1 chip or Apple T2 Security Chip. This controls fallback behavior when Erase All Content and Settings (EACS) cannot run.

Valid values:
- `default`: If EACS preflight fails, device responds with Error status and attempts to erase itself. If EACS preflight succeeds but EACS fails, the device attempts to erase itself.
- `doNotObliterate`: If EACS preflight fails, device responds with Error status and doesn't attempt to erase. If EACS preflight succeeds but EACS fails, the device doesn't attempt to erase itself.
- `obliterateWithWarning`: If EACS preflight fails, device responds with Acknowledged status and attempts to erase itself. If EACS preflight succeeds but EACS fails, the device attempts to erase itself.
- `always`: The system doesn't attempt EACS. T2 and later devices always obliterate.

**Note:** This setting only applies to Mac computers with Apple M1 chip or Apple T2 Security Chip running macOS 12 or later. It has no effect on machines prior to the T2 chip.

**Reference:** [obliterationBehavior enum](https://learn.microsoft.com/en-us/graph/api/resources/intune-devices-obliterationbehavior?view=graph-rest-beta)
- `persist_esim_data_plan` (Boolean) If `true`, preserves the eSIM data plan on the device during wipe. Only applicable to devices with eSIM support. Defaults to `false`.
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))
- `use_protected_wipe` (Boolean) If `true`, uses protected wipe for Windows 10/11 devices. Protected wipe maintains UEFI-embedded product keys and recovery partition. Only applicable to Windows devices. Defaults to `false`.
- `validate_device_exists` (Boolean) When set to `true` (default), the action will validate that all specified devices exist and are supported for wipe before attempting to wipe them. When `false`, device validation is skipped and the action will attempt to wipe devices directly. Disabling validation can improve performance but may result in errors if devices don't exist or are unsupported.

<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `invoke` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).

