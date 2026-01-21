---
page_title: "microsoft365_graph_beta_device_management_managed_device_windows_defender_scan Action - terraform-provider-microsoft365"
subcategory: "Device Management"

description: |-
  Triggers an antivirus scan on Windows devices using Windows Defender (Microsoft Defender Antivirus) in Microsoft Intune using the /deviceManagement/managedDevices/{managedDeviceId}/windowsDefenderScan and /deviceManagement/comanagedDevices/{managedDeviceId}/windowsDefenderScan endpoints. This action is used to initiate either a quick scan or full scan remotely on Windows devices managed by Intune.
  What This Action Does:
  Triggers immediate Windows Defender scanSupports both quick and full scan typesScans for viruses, malware, and threatsUpdates threat definitions before scanningReports results to IntuneCan be used for threat remediationWorks on managed and co-managed devices
  Scan Types:
  Quick Scan: Scans common threat locations (5-15 minutes)
  System folders and registry keysActive memory processesStartup locationsRecommended for routine scansFull Scan: Comprehensive scan of entire system (30+ minutes to hours)
  All files and foldersAll drives and partitionsArchive filesRecommended when threat detected or troubleshooting
  Platform Support:
  Windows 10/11: Full support (managed and co-managed)Windows Server: Full support (if Defender enabled)Other platforms: Not supported (Windows Defender only)
  Common Use Cases:
  Security incident responseThreat detection and remediationCompliance verificationPost-malware cleanupRoutine security checksAfter suspicious activityEmergency threat scanning
  Important Considerations:
  Device must be onlineFull scans can impact performanceScans run in backgroundResults reported to IntuneMay require user notificationCan be resource-intensive
  Reference: Microsoft Graph API - Windows Defender Scan https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-windowsdefenderscan?view=graph-rest-beta
---

# microsoft365_graph_beta_device_management_managed_device_windows_defender_scan (Action)

Triggers an antivirus scan on Windows devices using Windows Defender (Microsoft Defender Antivirus) in Microsoft Intune using the `/deviceManagement/managedDevices/{managedDeviceId}/windowsDefenderScan` and `/deviceManagement/comanagedDevices/{managedDeviceId}/windowsDefenderScan` endpoints. This action is used to initiate either a quick scan or full scan remotely on Windows devices managed by Intune.

**What This Action Does:**
- Triggers immediate Windows Defender scan
- Supports both quick and full scan types
- Scans for viruses, malware, and threats
- Updates threat definitions before scanning
- Reports results to Intune
- Can be used for threat remediation
- Works on managed and co-managed devices

**Scan Types:**
- **Quick Scan**: Scans common threat locations (5-15 minutes)
  - System folders and registry keys
  - Active memory processes
  - Startup locations
  - Recommended for routine scans
- **Full Scan**: Comprehensive scan of entire system (30+ minutes to hours)
  - All files and folders
  - All drives and partitions
  - Archive files
  - Recommended when threat detected or troubleshooting

**Platform Support:**
- **Windows 10/11**: Full support (managed and co-managed)
- **Windows Server**: Full support (if Defender enabled)
- **Other platforms**: Not supported (Windows Defender only)

**Common Use Cases:**
- Security incident response
- Threat detection and remediation
- Compliance verification
- Post-malware cleanup
- Routine security checks
- After suspicious activity
- Emergency threat scanning

**Important Considerations:**
- Device must be online
- Full scans can impact performance
- Scans run in background
- Results reported to Intune
- May require user notification
- Can be resource-intensive

**Reference:** [Microsoft Graph API - Windows Defender Scan](https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-windowsdefenderscan?view=graph-rest-beta)

## Microsoft Documentation

### Graph API References
- [windowsDefenderScan action](https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-windowsdefenderscan?view=graph-rest-beta)
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
# Example 1: Quick scan on a single Windows device - Minimal
action "microsoft365_graph_beta_device_management_managed_device_windows_defender_scan" "quick_scan_single" {
  config {
    managed_devices = [
      {
        device_id  = "12345678-1234-1234-1234-123456789abc"
        quick_scan = true
      }
    ]
  }
}

# Example 2: Full scan on a single Windows device
action "microsoft365_graph_beta_device_management_managed_device_windows_defender_scan" "full_scan_single" {
  config {
    managed_devices = [
      {
        device_id  = "12345678-1234-1234-1234-123456789abc"
        quick_scan = false
      }
    ]
  }
}

# Example 3: Mixed scans on multiple devices
action "microsoft365_graph_beta_device_management_managed_device_windows_defender_scan" "mixed_scans" {
  config {
    managed_devices = [
      {
        device_id  = "12345678-1234-1234-1234-123456789abc"
        quick_scan = true
      },
      {
        device_id  = "87654321-4321-4321-4321-ba9876543210"
        quick_scan = false
      }
    ]

    timeouts = {
      invoke = "20m"
    }
  }
}

# Example 4: Maximal configuration with validation
action "microsoft365_graph_beta_device_management_managed_device_windows_defender_scan" "scan_maximal" {
  config {
    managed_devices = [
      {
        device_id  = "12345678-1234-1234-1234-123456789abc"
        quick_scan = true
      },
      {
        device_id  = "87654321-4321-4321-4321-ba9876543210"
        quick_scan = false
      }
    ]

    comanaged_devices = [
      {
        device_id  = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
        quick_scan = true
      }
    ]

    ignore_partial_failures = true
    validate_device_exists  = true

    timeouts = {
      invoke = "5m"
    }
  }
}

# Example 5: Quick scan all Windows devices from data source
data "microsoft365_graph_beta_device_management_managed_device" "windows_devices" {
  filter_type  = "odata"
  odata_filter = "operatingSystem eq 'Windows'"
}

action "microsoft365_graph_beta_device_management_managed_device_windows_defender_scan" "quick_scan_all_windows" {
  config {
    managed_devices = [
      for device in data.microsoft365_graph_beta_device_management_managed_device.windows_devices.items : {
        device_id  = device.id
        quick_scan = true
      }
    ]

    validate_device_exists  = true
    ignore_partial_failures = true

    timeouts = {
      invoke = "30m"
    }
  }
}

# Example 6: Full scan on non-compliant devices
data "microsoft365_graph_beta_device_management_managed_device" "noncompliant_windows" {
  filter_type  = "odata"
  odata_filter = "(operatingSystem eq 'Windows') and (complianceState eq 'noncompliant')"
}

action "microsoft365_graph_beta_device_management_managed_device_windows_defender_scan" "full_scan_noncompliant" {
  config {
    managed_devices = [
      for device in data.microsoft365_graph_beta_device_management_managed_device.noncompliant_windows.items : {
        device_id  = device.id
        quick_scan = false
      }
    ]

    ignore_partial_failures = false

    timeouts = {
      invoke = "60m"
    }
  }
}

# Example 7: Scan co-managed devices
action "microsoft365_graph_beta_device_management_managed_device_windows_defender_scan" "scan_comanaged" {
  config {
    comanaged_devices = [
      {
        device_id  = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
        quick_scan = true
      },
      {
        device_id  = "bbbbbbbb-cccc-dddd-eeee-ffffffffffff"
        quick_scan = false
      }
    ]

    timeouts = {
      invoke = "15m"
    }
  }
}

# Output examples
output "scanned_devices_count" {
  value       = length(action.microsoft365_graph_beta_device_management_managed_device_windows_defender_scan.mixed_scans.config.managed_devices)
  description = "Number of devices that had scans initiated"
}
```

<!-- action schema generated by tfplugindocs -->
## Schema

### Optional

- `comanaged_devices` (Attributes List) List of co-managed Windows devices to scan with individual scan type configuration. These are devices managed by both Intune and Configuration Manager (SCCM). Each entry specifies a device ID and scan type.

**Co-Management Context:**
- Devices managed by both Intune and Configuration Manager
- Typically Windows 10/11 enterprise devices
- This action triggers Defender scan via Intune endpoint
- ConfigMgr can also trigger scans independently

**Note:** At least one of `managed_devices` or `comanaged_devices` must be provided. (see [below for nested schema](#nestedatt--comanaged_devices))
- `ignore_partial_failures` (Boolean) When set to `true`, the action will complete successfully even if some devices fail to scan. When `false` (default), the action will fail if any device scan fails. Use this flag when scanning multiple devices and you want the action to succeed even if some scans fail.
- `managed_devices` (Attributes List) List of managed Windows devices to scan with individual scan type configuration. Each entry specifies a device ID and whether to perform a quick scan or full scan. These are devices fully managed by Intune only.

Example:
```hcl
managed_devices = [
  {
    device_id  = "12345678-1234-1234-1234-123456789abc"
    quick_scan = true  # Quick scan (5-15 min)
  },
  {
    device_id  = "87654321-4321-4321-4321-ba9876543210"
    quick_scan = false # Full scan (30+ min)
  }
]
```

**Note:** At least one of `managed_devices` or `comanaged_devices` must be provided. (see [below for nested schema](#nestedatt--managed_devices))
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))
- `validate_device_exists` (Boolean) When set to `true` (default), the action will validate that all specified devices exist and are Windows devices before attempting to scan them. When `false`, device validation is skipped and the action will attempt to scan devices directly. Disabling validation can improve performance but may result in errors if devices don't exist or are not Windows devices.

<a id="nestedatt--comanaged_devices"></a>
### Nested Schema for `comanaged_devices`

Required:

- `device_id` (String) The co-managed device ID (GUID) of the Windows device to scan. Example: `12345678-1234-1234-1234-123456789abc`
- `quick_scan` (Boolean) Whether to perform a quick scan (`true`) or full scan (`false`). See managed_devices.quick_scan for detailed explanation of scan types.


<a id="nestedatt--managed_devices"></a>
### Nested Schema for `managed_devices`

Required:

- `device_id` (String) The managed device ID (GUID) of the Windows device to scan. Example: `12345678-1234-1234-1234-123456789abc`
- `quick_scan` (Boolean) Whether to perform a quick scan (`true`) or full scan (`false`).

- **Quick Scan (`true`)**: Fast scan of common threat locations (5-15 minutes)
  - Scans system folders, registry, memory, startup locations
  - Minimal impact on device performance
  - Recommended for routine/scheduled scans
  - Good for rapid security checks

- **Full Scan (`false`)**: Comprehensive scan of entire system (30+ minutes to hours)
  - Scans all files, folders, drives, archives
  - Higher impact on device performance
  - Recommended when threat detected
  - Thorough investigation of suspicious activity
  - Post-incident verification


<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `invoke` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).


