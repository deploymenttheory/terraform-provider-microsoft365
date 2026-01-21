---
page_title: "microsoft365_graph_beta_device_management_managed_device_windows_defender_update_signatures Action - terraform-provider-microsoft365"
subcategory: "Device Management"

description: |-
  Updates Windows Defender (Microsoft Defender Antivirus) signatures on Windows devices in Microsoft Intune using the /deviceManagement/managedDevices/{managedDeviceId}/windowsDefenderUpdateSignatures and /deviceManagement/comanagedDevices/{managedDeviceId}/windowsDefenderUpdateSignatures endpoints. This action is used to force devices to immediately update antivirus definitions without waiting for the standard update schedule.
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

Updates Windows Defender (Microsoft Defender Antivirus) signatures on Windows devices in Microsoft Intune using the `/deviceManagement/managedDevices/{managedDeviceId}/windowsDefenderUpdateSignatures` and `/deviceManagement/comanagedDevices/{managedDeviceId}/windowsDefenderUpdateSignatures` endpoints. This action is used to force devices to immediately update antivirus definitions without waiting for the standard update schedule.

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


