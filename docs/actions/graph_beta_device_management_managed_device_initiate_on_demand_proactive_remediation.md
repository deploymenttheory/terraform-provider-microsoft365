---
page_title: "microsoft365_graph_beta_device_management_managed_device_initiate_on_demand_proactive_remediation Action - terraform-provider-microsoft365"
subcategory: "Device Management"

description: |-
  Initiates on-demand proactive remediation on managed Windows devices using the /deviceManagement/managedDevices/{managedDeviceId}/initiateOnDemandProactiveRemediation and /deviceManagement/comanagedDevices/{managedDeviceId}/initiateOnDemandProactiveRemediation endpoints. Proactive remediations (also called remediations or health scripts) are PowerShell scripts that detect and automatically fix common support issues on Windows devices. This action triggers immediate execution of a specified remediation script on selected devices, rather than waiting for the scheduled run. This is useful for urgent fixes, troubleshooting, or validating remediation effectiveness.
  Important Notes:
  Only works on Windows 10/11 devicesRequires script policy ID (remediation script GUID)Script executes immediately on device check-inRuns with SYSTEM privilegesResults available in Intune portal and reportsScript must be already deployed to the deviceDoes not create new script deployment
  Use Cases:
  Urgent issue remediation outside scheduled runsTroubleshooting and validationPost-incident recovery actionsAd-hoc compliance fixesTesting new remediation scriptsEnd-user requested fixes
  Platform Support:
  Windows: Windows 10/11 with Intune management extensionOther Platforms: Not supported (Windows-specific feature)
  Reference: Microsoft Graph API - Initiate On Demand Proactive Remediation https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-initiateondemandproactiveremediation?view=graph-rest-beta
---

# microsoft365_graph_beta_device_management_managed_device_initiate_on_demand_proactive_remediation (Action)

Initiates on-demand proactive remediation on managed Windows devices using the `/deviceManagement/managedDevices/{managedDeviceId}/initiateOnDemandProactiveRemediation` and `/deviceManagement/comanagedDevices/{managedDeviceId}/initiateOnDemandProactiveRemediation` endpoints. Proactive remediations (also called remediations or health scripts) are PowerShell scripts that detect and automatically fix common support issues on Windows devices. This action triggers immediate execution of a specified remediation script on selected devices, rather than waiting for the scheduled run. This is useful for urgent fixes, troubleshooting, or validating remediation effectiveness.

**Important Notes:**
- Only works on Windows 10/11 devices
- Requires script policy ID (remediation script GUID)
- Script executes immediately on device check-in
- Runs with SYSTEM privileges
- Results available in Intune portal and reports
- Script must be already deployed to the device
- Does not create new script deployment

**Use Cases:**
- Urgent issue remediation outside scheduled runs
- Troubleshooting and validation
- Post-incident recovery actions
- Ad-hoc compliance fixes
- Testing new remediation scripts
- End-user requested fixes

**Platform Support:**
- **Windows**: Windows 10/11 with Intune management extension
- **Other Platforms**: Not supported (Windows-specific feature)

**Reference:** [Microsoft Graph API - Initiate On Demand Proactive Remediation](https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-initiateondemandproactiveremediation?view=graph-rest-beta)

## Microsoft Documentation

### Graph API References
- [initiateOnDemandProactiveRemediation action](https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-initiateondemandproactiveremediation?view=graph-rest-beta)
- [managedDevice resource type](https://learn.microsoft.com/en-us/graph/api/resources/intune-devices-manageddevice?view=graph-rest-beta)

### Proactive Remediations Guides
- [Proactive remediations in Intune](https://learn.microsoft.com/en-us/mem/intune/fundamentals/remediations)
- [Create and run remediations scripts](https://learn.microsoft.com/en-us/mem/analytics/proactive-remediations)
- [Monitor remediations script results](https://learn.microsoft.com/en-us/mem/intune/fundamentals/remediations#monitor-your-scripts)

## API Permissions

The following API permissions are required in order to use this action.

### Microsoft Graph

- **Application**: `DeviceManagementConfiguration.Read.All`, `DeviceManagementManagedDevices.Read.All`
- **Delegated**: `DeviceManagementConfiguration.Read.All`, `DeviceManagementManagedDevices.Read.All`

## Version History

| Version | Status | Notes |
|---------|--------|-------|
| v0.33.0-alpha | Experimental | Initial release |

## Notes

### Platform Compatibility

| Platform | Support | Requirements |
|----------|---------|--------------|
| **Windows** | ✅ Full Support | Windows 10/11 with Intune Management Extension installed |
| **macOS** | ❌ Not Supported | Proactive remediations are Windows-specific |
| **iOS/iPadOS** | ❌ Not Supported | Proactive remediations are Windows-specific |
| **Android** | ❌ Not Supported | Proactive remediations are Windows-specific |

### What are Proactive Remediations?

**Proactive remediations** (also called **remediations** or **health scripts**) are:
- PowerShell scripts that run on Windows devices
- Designed to detect and automatically fix common support issues
- Consist of two scripts: detection script and remediation script
- Execute with SYSTEM-level privileges
- Scheduled to run at regular intervals (hourly, daily, weekly)
- Provide detailed execution results and status reporting

### Script Components

| Component | Purpose | Execution |
|-----------|---------|-----------|
| **Detection Script** | Identifies if issue exists | Always runs first |
| **Remediation Script** | Fixes the detected issue | Only runs if detection finds issue |

**Exit Codes:**
- **Detection**: 0 = No issue, 1 = Issue detected
- **Remediation**: 0 = Success, 1 = Failure

### On-Demand vs Scheduled Execution

| Aspect | Scheduled Execution | On-Demand Execution (This Action) |
|--------|-------------------|-----------------------------------|
| **Timing** | Per policy schedule (hourly/daily/weekly) | Immediate on next device check-in |
| **Initiation** | Automatic | Manual via API call |
| **Use Case** | Regular proactive maintenance | Urgent fixes, troubleshooting |
| **Deployment** | Policy must be assigned | Policy must already be deployed |

### How On-Demand Remediation Works

1. **Action Triggered**: API call initiates on-demand execution
2. **Device Check-in**: Device connects to Intune (typically within minutes)
3. **Script Download**: Intune Management Extension downloads scripts
4. **Detection Runs**: Detection script executes with SYSTEM privileges
5. **Conditional Remediation**: If issue detected, remediation script runs
6. **Result Upload**: Exit codes and output uploaded to Intune
7. **Portal Update**: Results visible in Intune portal

### Finding Script Policy IDs

**Method 1: Azure Portal**
1. Navigate to **Intune** → **Devices** → **Remediations**
2. Select the remediation script
3. Copy GUID from URL: `...remediations/{GUID}/...`

**Method 2: Graph Explorer**
```http
GET https://graph.microsoft.com/beta/deviceManagement/deviceHealthScripts
```

**Method 3: PowerShell**
```powershell
Connect-MgGraph
Get-MgDeviceManagementDeviceHealthScript | Select-Object Id, DisplayName
```

### When to Use On-Demand Remediation

- **Urgent Issues**: Critical problems requiring immediate attention
- **Troubleshooting**: Testing or validating remediation effectiveness
- **End-User Requests**: User-reported issues needing quick resolution
- **Post-Incident**: Recovery actions after security or stability incidents
- **Validation**: Verifying new scripts before scheduled deployment
- **Ad-Hoc Fixes**: One-time corrections outside regular schedule

### Script Execution Details

**Execution Context:**
- Runs as **SYSTEM** account (highest privileges)
- 32-bit PowerShell by default (configurable in script settings)
- Maximum execution time: 60 minutes (timeout)
- Network access available
- Registry and file system full access

**Execution Timing:**
- Initiated at next device check-in (typically 8 hours or less)
- Can be forced immediately via device sync action
- Scripts queue if multiple remediation triggered

### Important Considerations

✅ **Requirements:**
- Script policy must already be deployed to device
- Device must have Intune Management Extension installed
- Windows 10 version 1607 or later
- Device must be online and checking in
- Script must be in "Published" state

⚠️ **Limitations:**
- Cannot create new script deployments
- Only triggers execution of existing assignments
- Script must be previously deployed to device
- Cannot modify script content via this action
- Results may take minutes to appear in portal

### Troubleshooting

**Common Issues:**

1. **Script Not Found**
   - Error: Script policy ID not recognized
   - Solution: Verify script policy ID is correct
   - Check: Script is published and deployed

2. **Device Not Responding**
   - Issue: Script doesn't execute
   - Solution: Force device sync first
   - Check: Device last check-in time

3. **Script Fails to Execute**
   - Issue: Script errors during execution
   - Solution: Review script logs on device
   - Location: `C:\ProgramData\Microsoft\IntuneManagementExtension\Logs`

4. **No Results Visible**
   - Issue: Results not showing in portal
   - Solution: Wait 15-30 minutes for sync
   - Check: Device connectivity and IME service status

### Best Practices

**Operational:**
- ✅ Test scripts thoroughly before on-demand execution
- ✅ Document why on-demand execution is needed
- ✅ Monitor script execution results closely
- ✅ Use for urgent issues only (respect scheduled runs)
- ✅ Validate script assignment before triggering
- ✅ Force device sync if immediate execution needed

**Security:**
- ✅ Review script content before execution (SYSTEM privileges)
- ✅ Audit all on-demand executions
- ✅ Restrict access to script policy IDs
- ✅ Monitor for unexpected script modifications
- ✅ Validate script source and author
- ✅ Test in non-production first

**Efficiency:**
- ✅ Batch similar issues together
- ✅ Use dynamic configuration for multiple devices
- ✅ Track script execution success rates
- ✅ Document common remediation patterns
- ✅ Maintain script library documentation

### Viewing Results

**Intune Portal:**
1. Navigate to **Devices** → **Remediations**
2. Select the remediation script
3. View **Device status** tab
4. Check individual device results

**What's Reported:**
- Detection status (Issue detected or not)
- Remediation status (Ran, succeeded, failed)
- Script output (stdout/stderr)
- Exit codes
- Execution timestamp
- Pre-remediation/post-remediation snapshots (if configured)

## Example Usage

```terraform
# Example 1: Run remediation script on single device
action "microsoft365_graph_beta_device_management_managed_device_initiate_on_demand_proactive_remediation" "single_device" {
  managed_devices {
    device_id        = "12345678-1234-1234-1234-123456789abc"
    script_policy_id = "87654321-4321-4321-4321-ba9876543210"
  }

  timeouts = {
    invoke = "5m"
  }
}

# Example 2: Run same remediation on multiple devices
action "microsoft365_graph_beta_device_management_managed_device_initiate_on_demand_proactive_remediation" "multiple_same_script" {
  managed_devices {
    device_id        = "device1-1234-1234-1234-123456789abc"
    script_policy_id = "disk-cleanup-script-guid-here"
  }

  managed_devices {
    device_id        = "device2-1234-1234-1234-123456789abc"
    script_policy_id = "disk-cleanup-script-guid-here"
  }

  managed_devices {
    device_id        = "device3-1234-1234-1234-123456789abc"
    script_policy_id = "disk-cleanup-script-guid-here"
  }

  timeouts = {
    invoke = "10m"
  }
}

# Example 3: Run different scripts on different devices
action "microsoft365_graph_beta_device_management_managed_device_initiate_on_demand_proactive_remediation" "different_scripts" {
  managed_devices {
    device_id        = "device1-1234-1234-1234-123456789abc"
    script_policy_id = "disk-cleanup-script-guid"
  }

  managed_devices {
    device_id        = "device2-1234-1234-1234-123456789abc"
    script_policy_id = "network-fix-script-guid"
  }

  managed_devices {
    device_id        = "device3-1234-1234-1234-123456789abc"
    script_policy_id = "printer-repair-script-guid"
  }

  timeouts = {
    invoke = "15m"
  }
}

# Example 4: Urgent remediation from variable
variable "urgent_remediation" {
  description = "Devices requiring urgent remediation"
  type = map(object({
    device_id        = string
    script_policy_id = string
  }))
  default = {
    "critical1" = {
      device_id        = "aaaa1111-1111-1111-1111-111111111111"
      script_policy_id = "emergency-fix-script-guid"
    }
    "critical2" = {
      device_id        = "bbbb2222-2222-2222-2222-222222222222"
      script_policy_id = "emergency-fix-script-guid"
    }
  }
}

action "microsoft365_graph_beta_device_management_managed_device_initiate_on_demand_proactive_remediation" "urgent_fix" {
  dynamic "managed_devices" {
    for_each = var.urgent_remediation
    content {
      device_id        = managed_devices.value.device_id
      script_policy_id = managed_devices.value.script_policy_id
    }
  }

  timeouts = {
    invoke = "10m"
  }
}

# Example 5: Post-incident remediation
locals {
  incident_devices = [
    {
      device_id        = "incident1-1111-1111-1111-111111111111"
      script_policy_id = "security-hardening-script-guid"
    },
    {
      device_id        = "incident2-2222-2222-2222-222222222222"
      script_policy_id = "security-hardening-script-guid"
    }
  ]
}

action "microsoft365_graph_beta_device_management_managed_device_initiate_on_demand_proactive_remediation" "incident_remediation" {
  dynamic "managed_devices" {
    for_each = local.incident_devices
    content {
      device_id        = managed_devices.value.device_id
      script_policy_id = managed_devices.value.script_policy_id
    }
  }

  timeouts = {
    invoke = "15m"
  }
}

# Example 6: Co-managed device remediation
action "microsoft365_graph_beta_device_management_managed_device_initiate_on_demand_proactive_remediation" "comanaged" {
  comanaged_devices {
    device_id        = "comanaged-1234-1234-1234-123456789abc"
    script_policy_id = "sccm-integration-fix-script-guid"
  }

  timeouts = {
    invoke = "5m"
  }
}

# Example 7: Troubleshooting specific issue
action "microsoft365_graph_beta_device_management_managed_device_initiate_on_demand_proactive_remediation" "troubleshoot_vpn" {
  managed_devices {
    device_id        = "vpn-issue-device-1234-123456789abc"
    script_policy_id = "vpn-troubleshoot-script-guid"
  }

  timeouts = {
    invoke = "10m"
  }
}

# Example 8: Testing new remediation script
action "microsoft365_graph_beta_device_management_managed_device_initiate_on_demand_proactive_remediation" "test_new_script" {
  managed_devices {
    device_id        = "test-device-1234-1234-1234-123456789abc"
    script_policy_id = "new-remediation-test-script-guid"
  }

  timeouts = {
    invoke = "5m"
  }
}

# Output examples
output "remediation_summary" {
  value = {
    managed_count   = length([for d in action.multiple_same_script.managed_devices : d])
    comanaged_count = length([for d in action.comanaged.comanaged_devices : d])
  }
  description = "Count of devices with remediation initiated"
}
```

<!-- action schema generated by tfplugindocs -->
## Schema

### Optional

- `comanaged_devices` (Block List) List of co-managed devices to initiate proactive remediation for. These are devices managed by both Intune and Configuration Manager (SCCM).

**Note:** At least one of `managed_devices` or `comanaged_devices` must be provided. (see [below for nested schema](#nestedblock--comanaged_devices))
- `managed_devices` (Block List) List of managed devices to initiate proactive remediation for. Each entry specifies a device and the remediation script to run.

**Note:** At least one of `managed_devices` or `comanaged_devices` must be provided. Each device can have a different script policy executed.

**Important:** The script policy must already be deployed to the device. This action triggers immediate execution but does not create a new deployment. (see [below for nested schema](#nestedblock--managed_devices))
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

<a id="nestedblock--comanaged_devices"></a>
### Nested Schema for `comanaged_devices`

Required:

- `device_id` (String) The unique identifier (GUID) of the co-managed device to run the remediation script on.

**Example**: `"12345678-1234-1234-1234-123456789abc"`
- `script_policy_id` (String) The unique identifier (GUID) of the proactive remediation script policy to execute.

**Note**: The script must already be assigned/deployed to the device.

**Example**: `"87654321-4321-4321-4321-ba9876543210"`


<a id="nestedblock--managed_devices"></a>
### Nested Schema for `managed_devices`

Required:

- `device_id` (String) The unique identifier (GUID) of the managed device to run the remediation script on.

**Example**: `"12345678-1234-1234-1234-123456789abc"`
- `script_policy_id` (String) The unique identifier (GUID) of the proactive remediation script policy to execute.

**How to find**: Azure Portal → Intune → Devices → Remediations → Select script → Copy GUID from URL or Properties.

**Note**: The script must already be assigned/deployed to the device.

**Example**: `"87654321-4321-4321-4321-ba9876543210"`


<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).

