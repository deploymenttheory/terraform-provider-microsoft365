---
page_title: "microsoft365_graph_beta_device_management_managed_device_clean_windows_device Action - terraform-provider-microsoft365"
subcategory: "Device Management"

description: |-
  Performs a clean operation on Windows managed and co-managed devices using the /deviceManagement/managedDevices/{managedDeviceId}/cleanWindowsDevice and /deviceManagement/comanagedDevices/{managedDeviceId}/cleanWindowsDevice endpoints. This action provides a lighter-weight alternative to full device wipe, allowing IT administrators to remove applications and settings while optionally preserving user data on each device independently.
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

Performs a clean operation on Windows managed and co-managed devices using the `/deviceManagement/managedDevices/{managedDeviceId}/cleanWindowsDevice` and `/deviceManagement/comanagedDevices/{managedDeviceId}/cleanWindowsDevice` endpoints. This action provides a lighter-weight alternative to full device wipe, allowing IT administrators to remove applications and settings while optionally preserving user data on each device independently.

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

| Platform | Support | Versions |
|----------|---------|----------|
| **Windows** | ✅ Full Support | Windows 10, Windows 11 |
| **macOS** | ❌ Not Supported | Windows-only action |
| **iOS/iPadOS** | ❌ Not Supported | Windows-only action |
| **Android** | ❌ Not Supported | Windows-only action |
| **ChromeOS** | ❌ Not Supported | Windows-only action |

### Clean vs Wipe vs Retire

| Action | What's Removed | User Data | Enrollment | Use Case |
|--------|----------------|-----------|------------|----------|
| **Clean** | Apps, settings | Optional | Preserved | Device refresh, app removal |
| **Wipe** | Everything | Removed | Lost (must re-enroll) | Full reset, security incidents |
| **Retire** | Company data only | Personal kept | Removed | BYOD unenrollment |

### What Gets Removed

#### When `keep_user_data = false` (Full Clean)
- Installed applications (except inbox Windows apps)
- User profiles
- User data (documents, desktop, downloads, etc.)
- Device configuration settings
- Company policies and profiles
- Personal customizations

#### When `keep_user_data = true` (App Clean)
- Installed applications (except inbox Windows apps)
- Device configuration settings
- Company policies and profiles

### What is Preserved

Regardless of `keep_user_data` setting:
- Windows OS installation
- Intune enrollment
- Device object in Intune
- Inbox Windows applications
- Device hardware configuration
- BitLocker encryption (if configured)

When `keep_user_data = true`:
- User profiles
- User documents and files
- Desktop, downloads, pictures, etc.
- User registry settings

### Common Use Cases

#### Full Clean (`keep_user_data = false`)
- Device reassignment to new user
- Malware removal requiring complete cleanup
- Compliance remediation
- Device refresh before new deployment
- Preparing device for return or disposal
- Removing previous user's data and settings

#### Preserve Data Clean (`keep_user_data = true`)
- Application troubleshooting
- Software bloat removal
- Performance optimization
- Configuration reset
- Maintaining user environment
- Fixing application conflicts
- Removing unwanted installed software

### Operation Details

#### Process Duration
- Typically 10-20 minutes
- Varies based on number of installed applications
- Device remains online during process
- Progress visible in Intune admin center

#### Network Requirements
- Device must be online
- Internet connectivity required
- Cannot be performed offline

#### User Experience
- Device becomes unavailable during clean
- User is automatically logged out
- Unsaved work is lost
- Device may restart
- Login prompt appears after completion

### Best Practices

- **Test First**: Always test with small device groups
- **Schedule Wisely**: Run during maintenance windows or off-hours
- **User Communication**: Notify users in advance of clean operations
- **Data Backup**: Ensure important data is backed up before cleaning
- **Choose Correctly**: Use `keep_user_data = true` for user devices when possible
- **Monitor Progress**: Check device status in Intune after cleaning
- **Sufficient Timeout**: Allow 15-20 minutes timeout for larger operations
- **Document Actions**: Log reason for clean in change management system

### Security Considerations

- **Full Clean**: Use for security incidents or device reassignment
- **Data Sanitization**: Full clean removes user data from device
- **Enrollment Persists**: Device remains in management (not like wipe)
- **BitLocker**: Encryption settings are maintained
- **Compliance**: Device will recheck compliance after clean
- **Policies Reapply**: All policies automatically reapply after clean

### Troubleshooting

| Issue | Possible Cause | Solution |
|-------|---------------|----------|
| Clean fails | Device offline | Ensure device is connected and online |
| Timeout | Large number of apps | Increase timeout value |
| User data remains | Wrong parameter | Verify `keep_user_data = false` |
| Apps remain | System apps | Inbox Windows apps cannot be removed |
| Device offline | Clean in progress | Wait for process to complete |

## Example Usage

```terraform
# ============================================================================
# Example 1: Clean single Windows device (remove user data)
# ============================================================================
# Use case: Full cleanup before reassignment
action "microsoft365_graph_beta_device_management_managed_device_clean_windows_device" "single_device_full_clean" {

  device_ids = ["12345678-1234-1234-1234-123456789abc"]
  
  keep_user_data = false

  timeouts = {
    invoke = "10m"
  }
}

# ============================================================================
# Example 2: Clean single Windows device (preserve user data)
# ============================================================================
# Use case: Application cleanup while keeping user profiles
action "microsoft365_graph_beta_device_management_managed_device_clean_windows_device" "single_device_preserve_data" {

  device_ids = ["12345678-1234-1234-1234-123456789abc"]
  
  keep_user_data = true

  timeouts = {
    invoke = "10m"
  }
}

# ============================================================================
# Example 3: Clean multiple Windows devices (remove user data)
# ============================================================================
# Use case: Bulk device refresh before new deployment
action "microsoft365_graph_beta_device_management_managed_device_clean_windows_device" "multiple_devices_full_clean" {

  device_ids = [
    "12345678-1234-1234-1234-123456789abc",
    "87654321-4321-4321-4321-ba9876543210",
    "abcdef12-3456-7890-abcd-ef1234567890"
  ]
  
  keep_user_data = false

  timeouts = {
    invoke = "20m"
  }
}

# ============================================================================
# Example 4: Clean non-compliant Windows devices (remove user data)
# ============================================================================
# Use case: Remediate non-compliant devices with full cleanup
data "microsoft365_graph_beta_device_management_managed_device" "windows_noncompliant" {
  filter_type  = "odata"
  odata_filter = "(operatingSystem eq 'Windows') and (complianceState eq 'noncompliant')"
}

action "microsoft365_graph_beta_device_management_managed_device_clean_windows_device" "clean_noncompliant" {

  device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.windows_noncompliant.items : device.id]
  
  keep_user_data = false

  timeouts = {
    invoke = "30m"
  }
}

# ============================================================================
# Example 5: Clean Windows 11 devices (preserve user data)
# ============================================================================
# Use case: Application troubleshooting while preserving user environment
data "microsoft365_graph_beta_device_management_managed_device" "windows11_devices" {
  filter_type  = "odata"
  odata_filter = "(operatingSystem eq 'Windows') and (startswith(osVersion, '10.0.22'))"
}

action "microsoft365_graph_beta_device_management_managed_device_clean_windows_device" "clean_win11_preserve_data" {

  device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.windows11_devices.items : device.id]
  
  keep_user_data = true

  timeouts = {
    invoke = "25m"
  }
}

# ============================================================================
# Example 6: Clean corporate Windows devices (remove user data)
# ============================================================================
# Use case: Company-owned device refresh
data "microsoft365_graph_beta_device_management_managed_device" "corporate_windows" {
  filter_type  = "odata"
  odata_filter = "(operatingSystem eq 'Windows') and (managedDeviceOwnerType eq 'company')"
}

action "microsoft365_graph_beta_device_management_managed_device_clean_windows_device" "clean_corporate_devices" {

  device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.corporate_windows.items : device.id]
  
  keep_user_data = false

  timeouts = {
    invoke = "30m"
  }
}

# ============================================================================
# IMPORTANT NOTES
# ============================================================================
# 
# Clean vs Wipe vs Retire:
# - Clean: Removes apps/settings, optionally keeps user data, device stays enrolled
# - Wipe: Factory reset, removes all data, device must re-enroll
# - Retire: Removes company data only, preserves personal data
#
# keep_user_data Parameter (REQUIRED):
# - false: Removes user profiles and data (full cleanup)
# - true: Preserves user profiles and data (app cleanup only)
# - Must be explicitly set - no default value
#
# Platform Support:
# - Windows 10: Fully supported
# - Windows 11: Fully supported
# - Other platforms: Not supported (Windows-only action)
#
# What Gets Removed (when keep_user_data = false):
# - Installed applications (except inbox Windows apps)
# - User profiles
# - User data (documents, desktop, etc.)
# - Device configuration settings
# - Company policies and profiles
#
# What Gets Removed (when keep_user_data = true):
# - Installed applications (except inbox Windows apps)
# - Device configuration settings
# - Company policies and profiles
# - User profiles and data are PRESERVED
#
# What is Preserved (Always):
# - Windows OS installation
# - Intune enrollment
# - Device in Intune management
# - Inbox Windows applications
# - Device hardware configuration
#
# Common Use Cases:
# - keep_user_data = false:
#   * Device reassignment to new user
#   * Malware removal requiring full cleanup
#   * Compliance remediation
#   * Device refresh before new deployment
#   * Preparing device for return/disposal
#
# - keep_user_data = true:
#   * Application troubleshooting
#   * Software bloat removal
#   * Configuration reset
#   * Maintaining user environment
#   * Performance optimization
#
# Best Practices:
# - Test with small device groups first
# - Schedule during maintenance windows
# - Notify users before cleaning their devices
# - Use keep_user_data=true for user devices when possible
# - Use keep_user_data=false for device reassignment
# - Document business justification
# - Monitor device status after clean
# - Allow sufficient timeout (devices may take 10-20 minutes)
#
# User Impact:
# - Device will be unavailable during clean process
# - Unsaved work will be lost
# - Process typically takes 10-20 minutes
# - Device remains enrolled (no re-enrollment needed)
# - User can log back in after clean completes
# - With keep_user_data=true, user profile is intact
# - With keep_user_data=false, user profile is removed
#
# Prerequisites:
# - Device must be Windows 10 or Windows 11
# - Device must be online
# - Device must be enrolled in Intune
# - Sufficient admin permissions required
#
# After Clean Operation:
# - Device remains in Intune
# - Policies will reapply automatically
# - Applications must be reinstalled
# - Check device status in Intune admin center
# - User may need to reconfigure personal settings
# - Company Portal remains available
```

<!-- action schema generated by tfplugindocs -->
## Schema

### Optional

- `comanaged_devices` (Block List) List of co-managed Windows devices to clean. These are devices managed by both Intune and Configuration Manager (SCCM). Each entry specifies a device ID and whether to preserve user data.

**Examples:**
```hcl
comanaged_devices = [
  {
    device_id       = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
    keep_user_data  = false
  }
]
```

**Note:** At least one of `managed_devices` or `comanaged_devices` must be provided. (see [below for nested schema](#nestedblock--comanaged_devices))
- `ignore_partial_failures` (Boolean) If set to `true`, the action will succeed even if some devices fail clean operation. Failed devices will be reported as warnings instead of errors. Default: `false` (action fails if any device fails).
- `managed_devices` (Block List) List of managed Windows devices to clean. These are devices fully managed by Intune only. Each entry specifies a device ID and whether to preserve user data.

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

**Note:** At least one of `managed_devices` or `comanaged_devices` must be provided. (see [below for nested schema](#nestedblock--managed_devices))
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))
- `validate_device_exists` (Boolean) Whether to validate that devices exist and are Windows devices before attempting clean. Disabling this can speed up planning but may result in runtime errors for non-existent or non-Windows devices. Default: `true`.

<a id="nestedblock--comanaged_devices"></a>
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


<a id="nestedblock--managed_devices"></a>
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

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).


