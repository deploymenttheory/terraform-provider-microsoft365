---
page_title: "microsoft365_graph_beta_device_management_managed_device_retire Action - terraform-provider-microsoft365"
subcategory: "Device Management"

description: |-
  Retires managed devices from Microsoft Intune using the /deviceManagement/managedDevices/{managedDeviceId}/retire endpoint. This action removes company data and managed apps from the device, while leaving personal data intact. The device is removed from Intune management and can no longer access company resources. This action supports retiring multiple devices in a single operation.
  Important Notes:
  For iOS/iPadOS devices, all data is removed except when enrolled via Device Enrollment Program (DEP) with User AffinityFor Windows devices, company data under %PROGRAMDATA%\Microsoft\MDM is removedFor Android devices, company data is removed and managed apps are uninstalledThis action cannot be reversed - devices must be re-enrolled to be managed again
  Reference: Microsoft Graph API - Retire Managed Device https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-retire?view=graph-rest-beta
---

# microsoft365_graph_beta_device_management_managed_device_retire (Action)

Retires managed devices from Microsoft Intune using the `/deviceManagement/managedDevices/{managedDeviceId}/retire` endpoint. This action removes company data and managed apps from the device, while leaving personal data intact. The device is removed from Intune management and can no longer access company resources. This action supports retiring multiple devices in a single operation.

**Important Notes:**
- For iOS/iPadOS devices, all data is removed except when enrolled via Device Enrollment Program (DEP) with User Affinity
- For Windows devices, company data under `%PROGRAMDATA%\Microsoft\MDM` is removed
- For Android devices, company data is removed and managed apps are uninstalled
- This action cannot be reversed - devices must be re-enrolled to be managed again

**Reference:** [Microsoft Graph API - Retire Managed Device](https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-retire?view=graph-rest-beta)

## Microsoft Documentation

### Graph API References
- [retire action](https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-retire?view=graph-rest-beta)
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

| Platform | Support | Data Removed | Personal Data Kept |
|----------|---------|--------------|-------------------|
| **Windows** | ✅ Full Support | Company apps, profiles, settings | User files, personal apps |
| **macOS** | ✅ Full Support | Management profiles, company apps | User data, personal apps |
| **iOS** | ✅ Full Support | Company apps, email profiles | Personal apps, photos |
| **iPadOS** | ✅ Full Support | Company apps, email profiles | Personal apps, photos |
| **Android** | ✅ Full Support | Work profile removed | Personal profile intact |

### Retire vs Wipe

| Action | Data Removed | Use Case |
|--------|--------------|----------|
| **Retire** | Company data only | BYOD devices, employee departures |
| **Wipe** | All data (factory reset) | Company-owned devices, security incidents |

### What Gets Removed

#### All Platforms
- Intune management enrollment
- Company email accounts
- Company apps and data
- VPN profiles
- Wi-Fi profiles
- Certificate profiles
- Configuration policies
- Compliance policies

#### Windows
- Company Portal app
- Microsoft 365 apps (if deployed)
- Company OneDrive data
- Windows Information Protection (WIP) data

#### iOS/iPadOS/macOS
- Managed apps and their data
- Email accounts configured by MDM
- Configuration profiles

#### Android
- Entire work profile (on BYOD)
- Company apps within work profile
- Work profile data

### What is Preserved

- Personal files and photos
- Personal apps
- Personal email accounts
- Device settings (wallpaper, etc.)
- Personal browsing history
- Personal contacts and calendar (non-corporate)

### Common Use Cases

- Employee leaving organization (BYOD)
- BYOD device unenrollment
- Transitioning device to personal use
- Removing corporate access gracefully
- Employee termination (selective wipe)
- Lost BYOD device (remove company data)
- Device ownership change
- End of MDM management

### Best Practices

- Preferred for BYOD devices
- Communicate with user before retiring
- Back up important company data first
- Document business justification
- Consider user data privacy
- Use wipe for company-owned devices instead
- Verify device ownership type
- Allow time for data sync/backup

### User Experience

- Device remains functional
- Personal data intact
- Company Portal removed/disabled
- Company apps removed
- Work profile deleted (Android)
- Device can be used personally
- Email/calendar may require reconfiguration

## Example Usage

```terraform
# Example 1: Retire a single managed device
action "microsoft365_graph_beta_device_management_managed_device_retire" "retire_single" {

  device_ids = [
    "12345678-1234-1234-1234-123456789abc"
  ]

  timeouts = {
    invoke = "5m"
  }
}

# Example 2: Retire multiple managed devices
action "microsoft365_graph_beta_device_management_managed_device_retire" "retire_batch" {

  device_ids = [
    "12345678-1234-1234-1234-123456789abc",
    "87654321-4321-4321-4321-ba9876543210",
    "abcdef12-3456-7890-abcd-ef1234567890"
  ]

  timeouts = {
    invoke = "10m"
  }
}

# Example 3: Retire devices from a data source query
# First, query for devices that meet certain criteria
data "microsoft365_graph_beta_device_management_managed_device" "non_compliant_devices" {
  filter_type  = "odata"
  odata_filter = "complianceState eq 'noncompliant'"
}

# Then retire those devices
action "microsoft365_graph_beta_device_management_managed_device_retire" "retire_non_compliant_devices" {

  # Extract device IDs from the data source
  device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.non_compliant_devices.items : device.id]

  timeouts = {
    invoke = "15m"
  }
}

# Example 4: Retire devices with specific operating system
data "microsoft365_graph_beta_device_management_managed_device" "old_ios_devices" {
  filter_type  = "odata"
  odata_filter = "operatingSystem eq 'iOS' and osVersion startsWith '14'"
}

action "microsoft365_graph_beta_device_management_managed_device_retire" "retire_old_ios_devices" {

  device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.old_ios_devices.items : device.id]

  timeouts = {
    invoke = "20m"
  }
}

# Output examples
output "retired_device_count" {
  value       = length(action.retire_batch.device_ids)
  description = "Number of devices retired in batch operation"
}

output "non_compliant_devices_to_retire" {
  value       = length(action.retire_non_compliant_devices.device_ids)
  description = "Number of non-compliant devices being retired"
}
```

<!-- action schema generated by tfplugindocs -->
## Schema

### Required

- `device_ids` (List of String) List of managed device IDs to retire from Intune management. Each ID must be a valid GUID format. Multiple devices can be retired in a single action. Example: `["12345678-1234-1234-1234-123456789abc", "87654321-4321-4321-4321-ba9876543210"]`

### Optional

- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).

