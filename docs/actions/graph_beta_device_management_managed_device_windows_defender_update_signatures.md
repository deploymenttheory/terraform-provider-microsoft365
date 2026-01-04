---
page_title: "microsoft365_graph_beta_device_management_managed_device_windows_defender_update_signatures Action - terraform-provider-microsoft365"
subcategory: "Device Management"

description: |-
  Forces Windows devices to immediately update Windows Defender (Microsoft Defender Antivirus) signatures using the /deviceManagement/managedDevices/{managedDeviceId}/windowsDefenderUpdateSignatures and /deviceManagement/comanagedDevices/{managedDeviceId}/windowsDefenderUpdateSignatures endpoints. This action triggers an immediate update of antivirus definitions without waiting for the standard update schedule.
  What This Action Does:
  Forces immediate signature updateDownloads latest threat definitionsUpdates malware detection databaseEnsures current threat protectionWorks on managed and co-managed devicesNo device reboot requiredCompletes in 1-5 minutes
  When to Use:
  Zero-day threat emergenceCritical security updatesBefore antivirus scansAfter new threat intelCompliance requirementsOutdated definitions detectedEmergency response scenarios
  Platform Support:
  Windows 10/11: Full support (managed and co-managed)Windows Server: Full support (if Defender enabled)Other platforms: Not supported (Windows Defender only)
  Update Process:
  Device receives update commandConnects to Microsoft Update serversDownloads latest signaturesApplies updates automaticallyReports completion to IntuneNo user interaction required
  Important Considerations:
  Device must be onlineInternet connectivity requiredMinimal performance impactUpdates in backgroundNo device reboot neededAutomatic threat protection
  Reference: Microsoft Graph API - Windows Defender Update Signatures https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-windowsdefenderupdatesignatures?view=graph-rest-beta
---

# microsoft365_graph_beta_device_management_managed_device_windows_defender_update_signatures (Action)

Forces Windows devices to immediately update Windows Defender (Microsoft Defender Antivirus) signatures using the `/deviceManagement/managedDevices/{managedDeviceId}/windowsDefenderUpdateSignatures` and `/deviceManagement/comanagedDevices/{managedDeviceId}/windowsDefenderUpdateSignatures` endpoints. This action triggers an immediate update of antivirus definitions without waiting for the standard update schedule.

**What This Action Does:**
- Forces immediate signature update
- Downloads latest threat definitions
- Updates malware detection database
- Ensures current threat protection
- Works on managed and co-managed devices
- No device reboot required
- Completes in 1-5 minutes

**When to Use:**
- Zero-day threat emergence
- Critical security updates
- Before antivirus scans
- After new threat intel
- Compliance requirements
- Outdated definitions detected
- Emergency response scenarios

**Platform Support:**
- **Windows 10/11**: Full support (managed and co-managed)
- **Windows Server**: Full support (if Defender enabled)
- **Other platforms**: Not supported (Windows Defender only)

**Update Process:**
- Device receives update command
- Connects to Microsoft Update servers
- Downloads latest signatures
- Applies updates automatically
- Reports completion to Intune
- No user interaction required

**Important Considerations:**
- Device must be online
- Internet connectivity required
- Minimal performance impact
- Updates in background
- No device reboot needed
- Automatic threat protection

**Reference:** [Microsoft Graph API - Windows Defender Update Signatures](https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-windowsdefenderupdatesignatures?view=graph-rest-beta)

## Microsoft Documentation

### Graph API References
- [windowsDefenderUpdateSignatures action](https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-windowsdefenderupdatesignatures?view=graph-rest-beta)
- [managedDevice resource type](https://learn.microsoft.com/en-us/graph/api/resources/intune-devices-manageddevice?view=graph-rest-beta)

### Intune Remote Actions Guides
- [Manage protection updates for Microsoft Defender Antivirus](https://learn.microsoft.com/en-us/defender-endpoint/manage-protection-updates-microsoft-defender-antivirus)

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

| Platform | Managed Devices | Co-Managed Devices | Notes |
|----------|----------------|-------------------|-------|
| **Windows 10/11** | ✅ Full Support | ✅ Full Support | Primary platform |
| **Windows Server** | ✅ Full Support | ✅ Full Support | If Defender enabled |
| **macOS** | ❌ Not Supported | ❌ Not Applicable | Different antivirus system |
| **iOS/iPadOS** | ❌ Not Supported | ❌ Not Applicable | No Defender signatures |
| **Android** | ❌ Not Supported | ❌ Not Applicable | Defender for Endpoint (different mechanism) |
| **ChromeOS** | ❌ Not Supported | ❌ Not Applicable | No Defender support |

### What Are Signatures?

**Virus/Malware Definitions**
- Database of known threats
- Patterns to identify malware
- Heuristic detection rules
- Behavioral analysis patterns
- Exploit detection signatures
- Updated multiple times daily

**Types of Signatures**
- **Virus Definitions**: Identifies specific malware
- **Spyware Definitions**: Detects spyware and adware
- **Rootkit Signatures**: Finds hidden malware
- **Network Inspection**: Network-based threats
- **Behavior Monitoring**: Suspicious activity patterns

**Update Frequency**
- Microsoft releases: Multiple times per day
- Critical threats: Immediate updates
- Normal schedule: Every 1-4 hours
- Manual update: This action
- Automatic updates: Managed by policy

### Normal Update Schedule vs Manual Update

| Aspect | Automatic Updates | Manual Update (This Action) |
|--------|------------------|----------------------------|
| **Frequency** | Every 1-4 hours | Immediate (on demand) |
| **Timing** | Background schedule | User/admin initiated |
| **Use Case** | Normal operations | Emergency/before scan |
| **Network Load** | Distributed over time | Immediate spike |
| **User Impact** | None | None |
| **Duration** | 1-5 minutes | 1-5 minutes |

### Update Process

**Step-by-Step**
1. Device receives update command
2. Windows Defender service checks current version
3. Connects to Microsoft Update servers
4. Downloads latest signature package
5. Verifies signature integrity
6. Applies signatures to database
7. Reports completion to Intune
8. Real-time protection continues

**During Update**
- Antivirus protection continues
- Real-time scanning active
- Background process
- Minimal CPU/network usage
- No user notification
- No device reboot required
- Takes 1-5 minutes

**After Update**
- Latest threat protection active
- New malware detectable
- Enhanced security posture
- Scan effectiveness improved
- Compliance requirements met

### What Gets Updated

**Signature Components**
- Virus and malware definitions
- Spyware definitions
- Potentially unwanted application (PUA) definitions
- Rootkit detection signatures
- Network inspection system signatures
- Behavioral monitoring rules
- Exploit detection patterns
- Cloud-delivered protection metadata

**Update Package Contents**
- Full signature database (first install)
- Delta updates (incremental)
- Engine updates (occasionally)
- Platform updates
- Configuration updates

**Size and Bandwidth**
- Delta update: 5-20 MB
- Full update: 100-300 MB
- Compressed transfer
- Efficient delivery
- Minimal network impact

### Use Cases

| Scenario | Description | Timing | Priority |
|----------|-------------|--------|----------|
| **Zero-Day Threat** | New threat discovered | Immediate | Critical |
| **Before Scan** | Ensure latest definitions | Pre-scan | High |
| **Compliance Audit** | Verify current signatures | Pre-audit | High |
| **Security Incident** | Threat response | Immediate | Critical |
| **Outdated Definitions** | Devices haven't updated | As needed | Medium |
| **Policy Deployment** | New security policy | After policy | Medium |
| **New Device Enrollment** | Freshly enrolled device | Post-enrollment | Medium |
| **Monthly Maintenance** | Regular refresh | Scheduled | Low |

### Best Practices

**When to Use**
- Before running antivirus scans
- After zero-day threat announcement
- When compliance requires current definitions
- Before security audits
- After detecting outdated signatures
- In response to security incidents
- For devices offline during normal updates

**When NOT to Use**
- If automatic updates are working
- Repeatedly within short periods
- For devices updated within last hour
- During bandwidth-constrained periods
- For non-Windows devices
- When normal schedule is sufficient

**Planning Considerations**
- Check device online status
- Verify internet connectivity
- Consider network bandwidth
- Time for completion (1-5 minutes per device)
- Stagger large batches
- Monitor update status
- Allow time for propagation

**Batch Management**
- Group devices logically
- Stagger updates to reduce network load
- Update critical devices first
- Monitor for failures
- Retry failed updates
- Document update rationale

### Performance Impact

**Network Usage**
- Download size: 5-300 MB per device
- Compressed transfer
- Typical: 10-20 MB delta
- Concurrent downloads supported
- Microsoft CDN delivery
- Minimal impact on users

**Device Resources**
- CPU: < 5% during update
- Memory: Minimal
- Disk I/O: Low
- Network: Brief spike
- No reboot required
- Background processing

**User Experience**
- No visible impact
- Can continue working
- No performance degradation
- No notifications (unless configured)
- Seamless update
- No interruption

### Monitoring and Verification

**Intune Admin Center**
1. Navigate to Devices → All devices
2. Select specific device
3. View "Defender Antivirus" section
4. Check "Signature version"
5. Review "Last signature update"
6. Verify "Signature up to date"

**Information Available**
- Current signature version
- Last update timestamp
- Update status
- Definition age
- Update success/failure
- Next scheduled update

**Time to Complete**
- Command delivery: < 1 minute
- Download: 1-3 minutes
- Application: < 1 minute
- Total: 1-5 minutes
- Console update: 5-15 minutes

### Troubleshooting

| Issue | Cause | Solution |
|-------|-------|----------|
| Update fails | Device offline | Wait for device to come online |
| No effect | Defender disabled | Enable Defender in policy |
| Timeout | Network issues | Check connectivity |
| Already current | Recent update | No action needed |
| Permission denied | Insufficient rights | Check API permissions |
| Download fails | Firewall blocking | Configure firewall rules |
| Non-Windows device | Wrong OS | Filter for Windows only |
| Service unavailable | Update servers down | Retry later |

### Update Sources

**Microsoft Update**
- Primary signature source
- Global CDN delivery
- Multiple times daily updates
- Automatic failover
- Reliable delivery

**Alternative Sources**
- Windows Update
- WSUS (Windows Server Update Services)
- Configuration Manager
- Internal distribution points
- Direct download (this action)

**Network Requirements**
- Internet connectivity required
- Access to Microsoft Update URLs
- Firewall permits HTTPS (443)
- Proxy configuration (if applicable)
- Sufficient bandwidth

### Signature Versions

**Version Numbering**
- Format: 1.xxx.yyyy.z
- Updated multiple times daily
- Incremental version numbers
- Track-able and verifiable
- Documented by Microsoft

**Checking Version**
- PowerShell: `Get-MpComputerStatus`
- Command line: `"%ProgramFiles%\Windows Defender\MpCmdRun.exe" -SignatureUpdate`
- Intune admin center
- Device properties
- Defender app UI

**Current vs Latest**
- Check definition age
- Compare to published version
- Verify freshness
- Ensure current protection
- Meet compliance requirements

### Co-Management Context

**How Updates Work**
- This action triggers update via Intune endpoint
- ConfigMgr can also manage signature updates
- No conflict between systems
- Both can coexist
- Updates apply regardless of workload assignment

**Coordination**
- Communicate with ConfigMgr team
- Avoid duplicate scheduling
- Use one system primarily
- Document management responsibility
- Share update status

**Best Practices**
- Designate primary update system
- Configure fallback sources
- Monitor both consoles
- Coordinate emergency updates
- Document procedures

### Compliance and Audit

**Audit Requirements**
- Definition version verification
- Update frequency tracking
- Current signature validation
- Compliance reporting
- Documentation retention

**Compliance Checking**
- Signature age < 24 hours
- Version matches latest release
- Update history available
- No failed updates
- Timely remediation

**Reporting**
- Update success/failure rates
- Signature version distribution
- Outdated device identification
- Compliance status tracking
- Audit trail maintenance

### Security Benefits

**Immediate Protection**
- Latest threat detection
- Zero-day coverage
- Current malware signatures
- Enhanced security posture
- Reduced risk window

**Threat Response**
- Rapid deployment capability
- Emergency update delivery
- Incident response tool
- Proactive defense
- Security hygiene

**Operational Benefits**
- Compliance assurance
- Audit readiness
- Consistent protection
- Centralized management
- Automated deployment

### Limitations

**Technical Limits**
- Windows devices only
- Requires Windows Defender
- Device must be online
- Internet connectivity required
- Update server availability dependent

**Operational Limits**
- Manual process (not scheduled)
- Requires IT action
- Network bandwidth consumption
- Batch size considerations
- API rate limits

**Platform Limits**
- Not available for non-Windows
- Third-party AV not affected
- Defender-specific only
- Requires active Defender installation

### Alternative Update Methods

**Automatic Updates**
- Default mechanism
- Background schedule
- Policy-controlled
- Transparent to users
- Recommended for normal operations

**Windows Update**
- Part of Windows Update
- Monthly cumulative updates
- Security intelligence updates
- Controlled by update policies

**Configuration Manager**
- Software update deployment
- Scheduled distribution
- Bandwidth management
- Reporting capabilities
- Enterprise scale

**Manual (Local)**
- PowerShell: `Update-MpSignature`
- Command line: `MpCmdRun.exe -SignatureUpdate`
- Defender app UI
- Local administrator action

## Example Usage

```terraform
# Example 1: Update signatures on a single managed device - Minimal
action "microsoft365_graph_beta_device_management_managed_device_windows_defender_update_signatures" "update_single" {
  config {
    managed_device_ids = [
      "12345678-1234-1234-1234-123456789abc"
    ]
  }
}

# Example 2: Update signatures on multiple managed devices
action "microsoft365_graph_beta_device_management_managed_device_windows_defender_update_signatures" "update_multiple_managed" {
  config {
    managed_device_ids = [
      "12345678-1234-1234-1234-123456789abc",
      "87654321-4321-4321-4321-ba9876543210",
      "abcdef12-3456-7890-abcd-ef1234567890"
    ]

    timeouts = {
      invoke = "10m"
    }
  }
}

# Example 3: Update signatures on co-managed devices
action "microsoft365_graph_beta_device_management_managed_device_windows_defender_update_signatures" "update_comanaged" {
  config {
    comanaged_device_ids = [
      "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
      "11111111-2222-3333-4444-555555555555"
    ]

    timeouts = {
      invoke = "10m"
    }
  }
}

# Example 4: Update both managed and co-managed devices - Maximal
action "microsoft365_graph_beta_device_management_managed_device_windows_defender_update_signatures" "update_mixed_devices" {
  config {
    managed_device_ids = [
      "12345678-1234-1234-1234-123456789abc",
      "87654321-4321-4321-4321-ba9876543210"
    ]

    comanaged_device_ids = [
      "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
    ]

    ignore_partial_failures = true
    validate_device_exists  = true

    timeouts = {
      invoke = "15m"
    }
  }
}

# Example 5: Update all Windows devices using datasource
data "microsoft365_graph_beta_device_management_managed_device" "all_windows" {
  filter_type  = "odata"
  odata_filter = "operatingSystem eq 'Windows'"
}

action "microsoft365_graph_beta_device_management_managed_device_windows_defender_update_signatures" "update_all_windows" {
  config {
    managed_device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.all_windows.items : device.id]

    validate_device_exists  = true
    ignore_partial_failures = true

    timeouts = {
      invoke = "30m"
    }
  }
}

# Example 6: Update signatures before scheduled scan
data "microsoft365_graph_beta_device_management_managed_device" "workstations" {
  filter_type  = "odata"
  odata_filter = "startsWith(deviceName, 'WKSTN-')"
}

action "microsoft365_graph_beta_device_management_managed_device_windows_defender_update_signatures" "pre_scan_update" {
  config {
    managed_device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.workstations.items : device.id]

    timeouts = {
      invoke = "15m"
    }
  }
}

# Example 7: Update non-compliant devices
data "microsoft365_graph_beta_device_management_managed_device" "non_compliant_windows" {
  filter_type  = "odata"
  odata_filter = "(operatingSystem eq 'Windows') and (complianceState eq 'noncompliant')"
}

action "microsoft365_graph_beta_device_management_managed_device_windows_defender_update_signatures" "update_non_compliant" {
  config {
    managed_device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.non_compliant_windows.items : device.id]

    validate_device_exists  = true
    ignore_partial_failures = false

    timeouts = {
      invoke = "20m"
    }
  }
}

# Example 8: Emergency threat response across fleet
data "microsoft365_graph_beta_device_management_managed_device" "all_windows_devices" {
  filter_type  = "odata"
  odata_filter = "operatingSystem eq 'Windows'"
}

data "microsoft365_graph_beta_device_management_managed_device" "all_comanaged" {
  filter_type  = "odata"
  odata_filter = "(operatingSystem eq 'Windows') and (managementAgent eq 'configurationManagerClientMdm')"
}

action "microsoft365_graph_beta_device_management_managed_device_windows_defender_update_signatures" "emergency_threat_response" {
  config {
    managed_device_ids   = [for device in data.microsoft365_graph_beta_device_management_managed_device.all_windows_devices.items : device.id]
    comanaged_device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.all_comanaged.items : device.id]

    validate_device_exists  = true
    ignore_partial_failures = true

    timeouts = {
      invoke = "60m"
    }
  }
}
```

<!-- action schema generated by tfplugindocs -->
## Schema

### Optional

- `comanaged_device_ids` (List of String) List of co-managed device IDs to update Windows Defender signatures. These are devices managed by both Intune and Configuration Manager (SCCM). Each ID must be a valid GUID format. Example: `["12345678-1234-1234-1234-123456789abc"]`

**Co-Management Context:**
- Devices managed by both Intune and Configuration Manager
- Typically Windows 10/11 enterprise devices
- This action updates signatures via Intune endpoint
- ConfigMgr can also manage definition updates independently
- No conflict between systems

**Note:** At least one of `managed_device_ids` or `comanaged_device_ids` must be provided.
- `ignore_partial_failures` (Boolean) When set to `true`, the action will complete successfully even if some devices fail to update signatures. When `false` (default), the action will fail if any device update fails. Use this flag when updating multiple devices and you want the action to succeed even if some updates fail.
- `managed_device_ids` (List of String) List of managed device IDs to update Windows Defender signatures. These are devices fully managed by Intune only. Each ID must be a valid GUID format. Multiple devices can be updated in a single action. Example: `["12345678-1234-1234-1234-123456789abc", "87654321-4321-4321-4321-ba9876543210"]`

**Note:** At least one of `managed_device_ids` or `comanaged_device_ids` must be provided. You can provide both to update different types of devices in one action.
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))
- `validate_device_exists` (Boolean) When set to `true` (default), the action will validate that all specified devices exist and are Windows devices before attempting to update signatures. When `false`, device validation is skipped and the action will attempt to update signatures directly. Disabling validation can improve performance but may result in errors if devices don't exist or are not Windows devices.

<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `invoke` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).


