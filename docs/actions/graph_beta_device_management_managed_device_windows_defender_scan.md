---
page_title: "microsoft365_graph_beta_device_management_managed_device_windows_defender_scan Action - terraform-provider-microsoft365"
subcategory: "Device Management"

description: |-
  Triggers an antivirus scan on Windows devices using Windows Defender (Microsoft Defender Antivirus) via the /deviceManagement/managedDevices/{managedDeviceId}/windowsDefenderScan and /deviceManagement/comanagedDevices/{managedDeviceId}/windowsDefenderScan endpoints. This action initiates either a quick scan or full scan remotely on Windows devices managed by Intune.
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

Triggers an antivirus scan on Windows devices using Windows Defender (Microsoft Defender Antivirus) via the `/deviceManagement/managedDevices/{managedDeviceId}/windowsDefenderScan` and `/deviceManagement/comanagedDevices/{managedDeviceId}/windowsDefenderScan` endpoints. This action initiates either a quick scan or full scan remotely on Windows devices managed by Intune.

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

| Platform | Managed Devices | Co-Managed Devices | Notes |
|----------|----------------|-------------------|-------|
| **Windows 10/11** | ✅ Full Support | ✅ Full Support | Primary platform |
| **Windows Server** | ✅ Full Support | ✅ Full Support | If Defender enabled |
| **macOS** | ❌ Not Supported | ❌ Not Applicable | Uses different antivirus |
| **iOS/iPadOS** | ❌ Not Supported | ❌ Not Applicable | No Defender available |
| **Android** | ❌ Not Supported | ❌ Not Applicable | Uses Defender for Endpoint (different action) |
| **ChromeOS** | ❌ Not Supported | ❌ Not Applicable | No Defender support |

### Scan Types Comparison

| Feature | Quick Scan | Full Scan |
|---------|-----------|-----------|
| **Duration** | 5-15 minutes | 30+ minutes to several hours |
| **Coverage** | Common threat locations | Entire system |
| **Performance Impact** | Minimal | Moderate to significant |
| **Recommended For** | Routine checks, daily scans | Security incidents, deep investigation |
| **Scans** | System folders, registry, memory, startup | All files, all drives, archives |
| **User Impact** | Minimal disruption | May slow down device |
| **Best Used** | During work hours | During off-hours |

### Quick Scan Details

**What Gets Scanned:**
- System folders (Windows, Program Files, ProgramData)
- Registry keys (startup, run keys)
- Active memory processes
- Boot sectors
- Startup locations
- Common malware locations

**Characteristics:**
- Fast completion (5-15 minutes)
- Minimal CPU/disk usage
- Users can work normally
- Detects most common threats
- Updates definitions before scanning
- Recommended for routine use

**Best For:**
- Daily/weekly scheduled scans
- Routine security checks
- Post-patch verification
- General maintenance
- During business hours
- BYOD devices

### Full Scan Details

**What Gets Scanned:**
- All files on all drives
- All folders and subfolders
- Archive files (ZIP, RAR, etc.)
- Boot sectors and firmware
- All registry keys
- Network drives (if mapped)
- Removable media
- Memory and processes

**Characteristics:**
- Thorough and comprehensive
- Long duration (30 min to hours)
- Higher CPU and disk usage
- May impact performance
- Detects hidden/deep threats
- More resource-intensive

**Best For:**
- Security incidents
- Malware suspected
- Post-infection cleanup
- Compliance requirements
- Off-hours scanning
- Monthly deep scans

### Scan Process

**Initialization**
1. Device receives scan command (if online)
2. Windows Defender service starts
3. Threat definitions updated
4. Scan begins immediately

**During Scan**
- Real-time protection continues
- Background process
- Can be paused by user
- Progress visible to user
- Threats quarantined automatically

**After Scan**
- Results reported to Intune
- Threats logged and quarantined
- Scan summary available
- User notification displayed
- Admin center updated

**Threat Handling**
- Threats automatically quarantined
- Malware removed or blocked
- Suspicious files isolated
- User and admin notified
- Remediation logged

### Performance Impact

| System Type | Quick Scan Impact | Full Scan Impact |
|-------------|------------------|------------------|
| **Modern Workstation** | Minimal (< 5% CPU) | Moderate (10-30% CPU) |
| **Laptop** | Minimal | Battery drain increase |
| **Server** | Minimal | May affect services |
| **Older Hardware** | Slight slowdown | Significant slowdown |
| **SSD Storage** | Negligible | Low |
| **HDD Storage** | Low | Moderate to High |

### Use Cases

| Scenario | Recommended Scan | Frequency | Timing |
|----------|-----------------|-----------|--------|
| **Routine Maintenance** | Quick | Daily/Weekly | Business hours |
| **Security Incident** | Full | Immediate | ASAP |
| **Malware Suspected** | Full | Immediate | ASAP |
| **Post-Patch** | Quick | After updates | Anytime |
| **Compliance Audit** | Full | Monthly/Quarterly | Off-hours |
| **New Device Enrollment** | Quick | Once | After enrollment |
| **User Report Issue** | Full | As needed | Off-hours preferred |
| **Threat Intel Alert** | Full | Emergency | Immediate |

### Best Practices

**Scan Scheduling**
- **Quick Scans**: During business hours (minimal impact)
- **Full Scans**: Off-hours, weekends, overnight
- Avoid scanning during critical business operations
- Consider time zones for global deployments
- Stagger scans to avoid network congestion

**Device Considerations**
- Check device uptime and availability
- Ensure devices are online
- Consider battery status for laptops
- Verify adequate disk space
- Check network bandwidth

**Resource Management**
- Batch devices into manageable groups
- Don't scan entire fleet simultaneously
- Monitor scan progress
- Plan for failed scans (offline devices)
- Set appropriate timeouts

**User Communication**
- Notify users before full scans
- Explain performance impact
- Provide estimated duration
- Advise saving work
- Set expectations

**Threat Response**
- Full scan immediately on threat detection
- Isolate affected devices
- Review scan results
- Follow up with remediation
- Document incidents

### Scan Results and Reporting

**Intune Admin Center**
1. Navigate to Devices → All devices
2. Select device
3. View "Defender Antivirus" section
4. Check "Last scan" timestamp
5. Review "Threats found"
6. Check "Scan status"

**Information Available**
- Last scan date/time
- Scan type (quick/full)
- Scan status (completed/failed/in progress)
- Threats found
- Threats quarantined
- Scan duration
- Definition version

**Time to Complete**
- Command delivery: < 1 minute
- Quick scan: 5-15 minutes
- Full scan: 30+ minutes to hours
- Results reporting: Real-time
- Console update: 5-15 minutes

### Threat Detection and Remediation

**Automatic Actions**
- Malware quarantined immediately
- Threats blocked from executing
- Suspicious files isolated
- Registry changes reverted
- System protected

**Threat Types Detected**
- Viruses and worms
- Trojans and spyware
- Ransomware
- Potentially unwanted applications (PUA)
- Rootkits
- Bootkits
- Exploits
- Backdoors

**Remediation Options**
- Automatic (recommended)
- User confirmation
- Custom actions
- Allow/block lists
- Exclusions

### Troubleshooting

| Issue | Cause | Solution |
|-------|-------|----------|
| Scan doesn't start | Device offline | Wait for device to come online |
| Scan fails | Defender disabled | Enable Defender in policy |
| Timeout | Full scan too long | Increase timeout setting |
| No results | Scan in progress | Wait for completion |
| Permission denied | Insufficient rights | Check API permissions |
| Device not found | Invalid device ID | Verify device exists |
| Non-Windows device | Wrong OS | Filter for Windows only |
| Scan cancelled | User intervention | Re-run scan, notify user |

### Limitations

**Technical Limits**
- Windows devices only
- Requires Windows Defender
- Device must be online
- Network connectivity required
- Definitions must be up-to-date

**Operational Limits**
- User can pause/cancel scan
- Performance impact during full scan
- Scan duration varies by device
- Offline devices won't scan
- Battery-powered devices may defer

**Platform Limits**
- Not available on non-Windows platforms
- Requires Windows Defender enabled
- Some devices may have third-party AV
- Co-managed devices require both systems healthy

### Windows Defender Requirements

**Minimum Requirements**
- Windows 10/11 or Windows Server
- Windows Defender Antivirus installed
- Windows Defender service running
- Real-time protection can be enabled or disabled
- Definitions can be outdated (updated before scan)

**Policy Requirements**
- Defender not disabled by group policy
- Scan permissions allowed
- Remote actions enabled
- Device enrolled in Intune

**Network Requirements**
- Internet connectivity
- Access to definition updates
- Access to Intune endpoints
- Firewall permits traffic

### Definition Updates

**Before Scan**
- Definitions automatically updated
- Latest signatures downloaded
- Ensures current threat detection
- May add 1-2 minutes to scan time

**Update Sources**
- Microsoft Update
- Windows Update
- WSUS
- Internal distribution
- Direct download

### Co-Management Context

**How Scanning Works**
- This action triggers scan via Intune endpoint
- ConfigMgr can also trigger scans independently
- Both systems can view results
- No conflict between systems
- Workload assignment doesn't affect scan capability

**Best Practices**
- Coordinate with ConfigMgr team
- Avoid duplicate scan scheduling
- Use one system for scan management
- Share scan results between teams
- Document which system manages scans

## Example Usage

```terraform
# ============================================================================
# Example 1: Quick scan on managed devices
# ============================================================================
# Use case: Routine security check on selected Windows devices
action "microsoft365_graph_beta_device_management_managed_device_windows_defender_scan" "quick_scan_managed" {

  managed_devices = [
    {
      device_id  = "12345678-1234-1234-1234-123456789abc"
      quick_scan = true
    },
    {
      device_id  = "87654321-4321-4321-4321-ba9876543210"
      quick_scan = true
    }
  ]

  timeouts = {
    invoke = "10m"
  }
}

# ============================================================================
# Example 2: Full scan on specific device
# ============================================================================
# Use case: Comprehensive scan after security incident
action "microsoft365_graph_beta_device_management_managed_device_windows_defender_scan" "full_scan_incident" {

  managed_devices = [
    {
      device_id  = "12345678-1234-1234-1234-123456789abc"
      quick_scan = false # Full comprehensive scan
    }
  ]

  timeouts = {
    invoke = "5m"
  }
}

# ============================================================================
# Example 3: Mixed scan types on different devices
# ============================================================================
# Use case: Quick scan for most, full scan for suspicious devices
action "microsoft365_graph_beta_device_management_managed_device_windows_defender_scan" "mixed_scan_types" {

  managed_devices = [
    {
      device_id  = "12345678-1234-1234-1234-123456789abc"
      quick_scan = true # Routine check
    },
    {
      device_id  = "87654321-4321-4321-4321-ba9876543210"
      quick_scan = false # Suspected malware - full scan
    },
    {
      device_id  = "abcdef12-3456-7890-abcd-ef1234567890"
      quick_scan = true # Routine check
    }
  ]

  timeouts = {
    invoke = "10m"
  }
}

# ============================================================================
# Example 4: Scan co-managed devices
# ============================================================================
# Use case: Scan Windows devices managed by both Intune and ConfigMgr
action "microsoft365_graph_beta_device_management_managed_device_windows_defender_scan" "scan_comanaged" {

  comanaged_devices = [
    {
      device_id  = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
      quick_scan = true
    },
    {
      device_id  = "11111111-2222-3333-4444-555555555555"
      quick_scan = true
    }
  ]

  timeouts = {
    invoke = "10m"
  }
}

# ============================================================================
# Example 5: Scan both managed and co-managed devices
# ============================================================================
# Use case: Mixed environment with different management types
action "microsoft365_graph_beta_device_management_managed_device_windows_defender_scan" "scan_mixed_management" {

  managed_devices = [
    {
      device_id  = "12345678-1234-1234-1234-123456789abc"
      quick_scan = true
    }
  ]

  comanaged_devices = [
    {
      device_id  = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
      quick_scan = true
    }
  ]

  timeouts = {
    invoke = "10m"
  }
}

# ============================================================================
# Example 6: Quick scan all Windows devices using datasource
# ============================================================================
# Use case: Routine security scan across entire Windows fleet
data "microsoft365_graph_beta_device_management_managed_device" "all_windows" {
  filter_type  = "odata"
  odata_filter = "operatingSystem eq 'Windows'"
}

action "microsoft365_graph_beta_device_management_managed_device_windows_defender_scan" "scan_all_windows" {

  managed_devices = [
    for device in data.microsoft365_graph_beta_device_management_managed_device.all_windows.items : {
      device_id  = device.id
      quick_scan = true
    }
  ]

  timeouts = {
    invoke = "30m"
  }
}

# ============================================================================
# Example 7: Full scan on non-compliant Windows devices
# ============================================================================
# Use case: Thorough scan on devices that failed compliance
data "microsoft365_graph_beta_device_management_managed_device" "non_compliant_windows" {
  filter_type  = "odata"
  odata_filter = "(operatingSystem eq 'Windows') and (complianceState eq 'noncompliant')"
}

action "microsoft365_graph_beta_device_management_managed_device_windows_defender_scan" "scan_non_compliant" {

  managed_devices = [
    for device in data.microsoft365_graph_beta_device_management_managed_device.non_compliant_windows.items : {
      device_id  = device.id
      quick_scan = false # Full scan for non-compliant devices
    }
  ]

  timeouts = {
    invoke = "20m"
  }
}

# ============================================================================
# Example 8: Scan Windows devices by naming convention
# ============================================================================
# Use case: Scan specific department or location devices
data "microsoft365_graph_beta_device_management_managed_device" "finance_windows" {
  filter_type  = "device_name"
  filter_value = "FIN-WS-"
}

action "microsoft365_graph_beta_device_management_managed_device_windows_defender_scan" "scan_finance_dept" {

  managed_devices = [
    for device in data.microsoft365_graph_beta_device_management_managed_device.finance_windows.items : {
      device_id  = device.id
      quick_scan = true
    }
  ]

  timeouts = {
    invoke = "15m"
  }
}

# ============================================================================
# Example 9: After-hours full scan on workstations
# ============================================================================
# Use case: Comprehensive scan during off-hours to avoid performance impact
data "microsoft365_graph_beta_device_management_managed_device" "workstations" {
  filter_type  = "device_name"
  filter_value = "WKSTN-"
}

action "microsoft365_graph_beta_device_management_managed_device_windows_defender_scan" "after_hours_full_scan" {

  managed_devices = [
    for device in data.microsoft365_graph_beta_device_management_managed_device.workstations.items : {
      device_id  = device.id
      quick_scan = false # Full scan during off-hours
    }
  ]

  timeouts = {
    invoke = "30m"
  }
}

# ============================================================================
# Example 10: Conditional scan based on last scan time
# ============================================================================
# Use case: Full scan on devices that haven't been scanned recently
locals {
  # Devices needing full scan (example logic)
  devices_need_full_scan = [
    "12345678-1234-1234-1234-123456789abc",
    "87654321-4321-4321-4321-ba9876543210"
  ]

  # Devices needing quick scan
  devices_need_quick_scan = [
    "abcdef12-3456-7890-abcd-ef1234567890",
    "fedcba98-7654-3210-fedc-ba9876543210"
  ]
}

action "microsoft365_graph_beta_device_management_managed_device_windows_defender_scan" "conditional_scan" {

  managed_devices = concat(
    [
      for device_id in local.devices_need_full_scan : {
        device_id  = device_id
        quick_scan = false
      }
    ],
    [
      for device_id in local.devices_need_quick_scan : {
        device_id  = device_id
        quick_scan = true
      }
    ]
  )

  timeouts = {
    invoke = "20m"
  }
}

# ============================================================================
# Example 11: Emergency threat response scan
# ============================================================================
# Use case: Immediate full scan after threat intel indicates new malware
data "microsoft365_graph_beta_device_management_managed_device" "all_windows_devices" {
  filter_type  = "odata"
  odata_filter = "operatingSystem eq 'Windows'"
}

action "microsoft365_graph_beta_device_management_managed_device_windows_defender_scan" "emergency_threat_scan" {

  managed_devices = [
    for device in data.microsoft365_graph_beta_device_management_managed_device.all_windows_devices.items : {
      device_id  = device.id
      quick_scan = false # Full scan for threat response
    }
  ]

  timeouts = {
    invoke = "60m"
  }
}

# ============================================================================
# Example 12: Scan Windows servers only
# ============================================================================
# Use case: Security scan on Windows Server infrastructure
data "microsoft365_graph_beta_device_management_managed_device" "windows_servers" {
  filter_type  = "odata"
  odata_filter = "(operatingSystem eq 'Windows') and (contains(deviceName, 'SRV'))"
}

action "microsoft365_graph_beta_device_management_managed_device_windows_defender_scan" "scan_servers" {

  managed_devices = [
    for device in data.microsoft365_graph_beta_device_management_managed_device.windows_servers.items : {
      device_id  = device.id
      quick_scan = true # Quick scan for servers to minimize impact
    }
  ]

  timeouts = {
    invoke = "20m"
  }
}

# ============================================================================
# Example 13: Scan by user assignment
# ============================================================================
# Use case: Scan all Windows devices assigned to specific user
data "microsoft365_graph_beta_device_management_managed_device" "user_devices" {
  filter_type  = "odata"
  odata_filter = "(operatingSystem eq 'Windows') and (userPrincipalName eq 'john.doe@company.com')"
}

action "microsoft365_graph_beta_device_management_managed_device_windows_defender_scan" "scan_user_devices" {

  managed_devices = [
    for device in data.microsoft365_graph_beta_device_management_managed_device.user_devices.items : {
      device_id  = device.id
      quick_scan = true
    }
  ]

  timeouts = {
    invoke = "10m"
  }
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

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).


