---
page_title: "microsoft365_graph_beta_device_management_managed_device_clean_windows_device Action - terraform-provider-microsoft365"
subcategory: "Device Management"

description: |-
  Performs a clean operation on Windows managed and co-managed devices in Microsoft Intune using the /deviceManagement/managedDevices/{managedDeviceId}/cleanWindowsDevice and /deviceManagement/comanagedDevices/{managedDeviceId}/cleanWindowsDevice endpoints. This action is used to remove applications and settings while optionally preserving user data, providing a lighter-weight alternative to full device wipe. IT administrators can remove applications and settings while optionally preserving user data on each device independently.
  What Clean Windows Device Does:
  Removes installed applications (except inbox Windows apps)Removes user profiles (unless keep_user_data is true for that device)Removes device configuration settingsRemoves company policies and profilesCan preserve user data per-device if specifiedDevice remains enrolled in IntuneLess destructive than full wipe
  Platform Support:
  Windows: Full support (Windows 10/11)Other platforms: Not supported (Windows-only action)
  Clean vs Wipe vs Retire:
  Clean: Removes apps/settings, optionally keeps user data, device stays enrolledWipe: Factory reset, removes all data, device must re-enrollRetire: Removes company data only, preserves personal data
  Common Use Cases:
  Device refresh without full rebuildRemoving malware/unwanted applicationsPreparing device for new user (keeping OS)Troubleshooting device issuesCompliance remediationSoftware bloat removalMaintaining device enrollment
  Important Considerations:
  Device must be online to receive commandUser will lose unsaved workInstalled applications will be removedProcess may take several minutesDevice remains in Intune (no re-enrollment needed)Each device can have different keep_user_data setting
  Reference: Microsoft Graph API - Clean Windows Device https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-cleanwindowsdevice?view=graph-rest-beta
---

# microsoft365_graph_beta_device_management_managed_device_clean_windows_device (Action)

Performs a clean operation on Windows managed and co-managed devices in Microsoft Intune using the `/deviceManagement/managedDevices/{managedDeviceId}/cleanWindowsDevice` and `/deviceManagement/comanagedDevices/{managedDeviceId}/cleanWindowsDevice` endpoints. This action is used to remove applications and settings while optionally preserving user data, providing a lighter-weight alternative to full device wipe. IT administrators can remove applications and settings while optionally preserving user data on each device independently.

**What Clean Windows Device Does:**
- Removes installed applications (except inbox Windows apps)
- Removes user profiles (unless `keep_user_data` is true for that device)
- Removes device configuration settings
- Removes company policies and profiles
- Can preserve user data per-device if specified
- Device remains enrolled in Intune
- Less destructive than full wipe

**Platform Support:**
- **Windows**: Full support (Windows 10/11)
- **Other platforms**: Not supported (Windows-only action)

**Clean vs Wipe vs Retire:**
- **Clean**: Removes apps/settings, optionally keeps user data, device stays enrolled
- **Wipe**: Factory reset, removes all data, device must re-enroll
- **Retire**: Removes company data only, preserves personal data

**Common Use Cases:**
- Device refresh without full rebuild
- Removing malware/unwanted applications
- Preparing device for new user (keeping OS)
- Troubleshooting device issues
- Compliance remediation
- Software bloat removal
- Maintaining device enrollment

**Important Considerations:**
- Device must be online to receive command
- User will lose unsaved work
- Installed applications will be removed
- Process may take several minutes
- Device remains in Intune (no re-enrollment needed)
- Each device can have different `keep_user_data` setting

**Reference:** [Microsoft Graph API - Clean Windows Device](https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-cleanwindowsdevice?view=graph-rest-beta)

## Microsoft Documentation

### Graph API References
- [cleanWindowsDevice action](https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-cleanwindowsdevice?view=graph-rest-beta)
- [managedDevice resource type](https://learn.microsoft.com/en-us/graph/api/resources/intune-devices-manageddevice?view=graph-rest-beta)

### Intune Remote Actions Guides
- [Windows Remote Actions](https://learn.microsoft.com/en-us/intune/intune-service/remote-actions/?tabs=windows)
- [iOS/iPadOS Remote Actions](https://learn.microsoft.com/en-us/intune/intune-service/remote-actions/?tabs=ios-ipados)
- [macOS Remote Actions](https://learn.microsoft.com/en-us/intune/intune-service/remote-actions/?tabs=macos)
- [Android Remote Actions](https://learn.microsoft.com/en-us/intune/intune-service/remote-actions/?tabs=android)
- [ChromeOS Remote Actions](https://learn.microsoft.com/en-us/intune/intune-service/remote-actions/?tabs=chromeos)

## Microsoft Graph API Permissions

The following client `application` permissions are needed in order to use this action:

**Required:**
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
# Example 1: Clean single Windows device (remove user data) - Minimal
action "microsoft365_graph_beta_device_management_managed_device_clean_windows_device" "single_device_full_clean" {
  config {
    managed_devices = [
      {
        device_id      = "12345678-1234-1234-1234-123456789abc"
        keep_user_data = false
      }
    ]

    timeouts = {
      invoke = "10m"
    }
  }
}

# Example 2: Clean single Windows device (preserve user data)
action "microsoft365_graph_beta_device_management_managed_device_clean_windows_device" "single_device_preserve_data" {
  config {
    managed_devices = [
      {
        device_id      = "12345678-1234-1234-1234-123456789abc"
        keep_user_data = true
      }
    ]

    timeouts = {
      invoke = "10m"
    }
  }
}

# Example 3: Clean multiple Windows devices with different options per device
action "microsoft365_graph_beta_device_management_managed_device_clean_windows_device" "multiple_devices_mixed_options" {
  config {
    managed_devices = [
      {
        device_id      = "12345678-1234-1234-1234-123456789abc"
        keep_user_data = false
      },
      {
        device_id      = "87654321-4321-4321-4321-ba9876543210"
        keep_user_data = true
      },
      {
        device_id      = "abcdef12-3456-7890-abcd-ef1234567890"
        keep_user_data = false
      }
    ]

    validate_device_exists  = true
    ignore_partial_failures = true

    timeouts = {
      invoke = "20m"
    }
  }
}

# Example 4: Clean co-managed Windows devices
action "microsoft365_graph_beta_device_management_managed_device_clean_windows_device" "comanaged_devices" {
  config {
    comanaged_devices = [
      {
        device_id      = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
        keep_user_data = false
      },
      {
        device_id      = "bbbbbbbb-cccc-dddd-eeee-ffffffffffff"
        keep_user_data = true
      }
    ]

    timeouts = {
      invoke = "15m"
    }
  }
}

# Example 5: Clean Windows devices from data source
data "microsoft365_graph_beta_device_management_managed_device" "windows_noncompliant" {
  filter_type  = "odata"
  odata_filter = "(operatingSystem eq 'Windows') and (complianceState eq 'noncompliant')"
}

action "microsoft365_graph_beta_device_management_managed_device_clean_windows_device" "clean_noncompliant" {
  config {
    managed_devices = [
      for device in data.microsoft365_graph_beta_device_management_managed_device.windows_noncompliant.items : {
        device_id      = device.id
        keep_user_data = false
      }
    ]

    validate_device_exists  = true
    ignore_partial_failures = false

    timeouts = {
      invoke = "30m"
    }
  }
}

# Example 6: Maximal configuration with both managed and co-managed devices
action "microsoft365_graph_beta_device_management_managed_device_clean_windows_device" "maximal_config" {
  config {
    managed_devices = [
      {
        device_id      = "12345678-1234-1234-1234-123456789abc"
        keep_user_data = false
      },
      {
        device_id      = "87654321-4321-4321-4321-987654321cba"
        keep_user_data = true
      }
    ]

    comanaged_devices = [
      {
        device_id      = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
        keep_user_data = false
      }
    ]

    ignore_partial_failures = false
    validate_device_exists  = true

    timeouts = {
      invoke = "5m"
    }
  }
}
```

<!-- action schema generated by tfplugindocs -->
## Schema

### Optional

- `comanaged_devices` (Attributes List) List of co-managed Windows devices to clean. These are devices managed by both Intune and Configuration Manager (SCCM). Each entry specifies a device ID and whether to preserve user data.

**Examples:**
```hcl
comanaged_devices = [
  {
    device_id       = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
    keep_user_data  = false
  }
]
```

**Note:** At least one of `managed_devices` or `comanaged_devices` must be provided. (see [below for nested schema](#nestedatt--comanaged_devices))
- `ignore_partial_failures` (Boolean) If set to `true`, the action will succeed even if some devices fail clean operation. Failed devices will be reported as warnings instead of errors. Default: `false` (action fails if any device fails).
- `managed_devices` (Attributes List) List of managed Windows devices to clean. These are devices fully managed by Intune only. Each entry specifies a device ID and whether to preserve user data.

**Examples:**
```hcl
managed_devices = [
  {
    device_id       = "12345678-1234-1234-1234-123456789abc"
    keep_user_data  = false
  },
  {
    device_id       = "87654321-4321-4321-4321-987654321cba"
    keep_user_data  = true
  }
]
```

**Note:** At least one of `managed_devices` or `comanaged_devices` must be provided. (see [below for nested schema](#nestedatt--managed_devices))
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))
- `validate_device_exists` (Boolean) Whether to validate that devices exist and are Windows devices before attempting clean. Disabling this can speed up planning but may result in runtime errors for non-existent or non-Windows devices. Default: `true`.

<a id="nestedatt--comanaged_devices"></a>
### Nested Schema for `comanaged_devices`

Required:

- `device_id` (String) The unique identifier (GUID) of the co-managed Windows device to clean. Example: `"12345678-1234-1234-1234-123456789abc"`
- `keep_user_data` (Boolean) Determines whether user data should be preserved for this device during the clean operation. **Required field** - must be explicitly set to `true` or `false`.

**When `false`:**
- User profiles removed
- User data deleted
- Applications removed
- Settings reset

**When `true`:**
- User profiles preserved
- User data kept (documents, desktop, etc.)
- Applications still removed
- Settings still reset


<a id="nestedatt--managed_devices"></a>
### Nested Schema for `managed_devices`

Required:

- `device_id` (String) The unique identifier (GUID) of the managed Windows device to clean. Device must be Windows 10 or Windows 11. Example: `"12345678-1234-1234-1234-123456789abc"`
- `keep_user_data` (Boolean) Determines whether user data should be preserved for this device during the clean operation. **Required field** - must be explicitly set to `true` or `false`.

**When `false`:**
- User profiles removed
- User data deleted
- Applications removed
- Settings reset

**When `true`:**
- User profiles preserved
- User data kept (documents, desktop, etc.)
- Applications still removed
- Settings still reset


<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `invoke` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).


