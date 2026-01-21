---
page_title: "microsoft365_graph_beta_device_management_managed_device_initiate_on_demand_proactive_remediation Action - terraform-provider-microsoft365"
subcategory: "Device Management"

description: |-
  Initiates on-demand proactive remediation on managed Windows devices in Microsoft Intune using the /deviceManagement/managedDevices/{managedDeviceId}/initiateOnDemandProactiveRemediation and /deviceManagement/comanagedDevices/{managedDeviceId}/initiateOnDemandProactiveRemediation endpoints. This action is used to trigger immediate execution of a specified remediation script on selected devices, rather than waiting for the scheduled run. Proactive remediations (also called remediations or health scripts) are PowerShell scripts that detect and automatically fix common support issues on Windows devices. This is useful for urgent fixes, troubleshooting, or validating remediation effectiveness.
  Important Notes:
  Only works on Windows 10/11 devicesRequires script policy ID (remediation script GUID)Script executes immediately on device check-inRuns with SYSTEM privilegesResults available in Intune portal and reportsScript must be already deployed to the deviceDoes not create new script deployment
  Use Cases:
  Urgent issue remediation outside scheduled runsTroubleshooting and validationPost-incident recovery actionsAd-hoc compliance fixesTesting new remediation scriptsEnd-user requested fixes
  Platform Support:
  Windows: Windows 10/11 with Intune management extensionOther Platforms: Not supported (Windows-specific feature)
  Reference: Microsoft Graph API - Initiate On Demand Proactive Remediation https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-initiateondemandproactiveremediation?view=graph-rest-beta
---

# microsoft365_graph_beta_device_management_managed_device_initiate_on_demand_proactive_remediation (Action)

Initiates on-demand proactive remediation on managed Windows devices in Microsoft Intune using the `/deviceManagement/managedDevices/{managedDeviceId}/initiateOnDemandProactiveRemediation` and `/deviceManagement/comanagedDevices/{managedDeviceId}/initiateOnDemandProactiveRemediation` endpoints. This action is used to trigger immediate execution of a specified remediation script on selected devices, rather than waiting for the scheduled run. Proactive remediations (also called remediations or health scripts) are PowerShell scripts that detect and automatically fix common support issues on Windows devices. This is useful for urgent fixes, troubleshooting, or validating remediation effectiveness.

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

## Microsoft Graph API Permissions

The following client `application` permissions are needed in order to use this action:

**Required:**
- `DeviceManagementConfiguration.Read.All`
- `DeviceManagementManagedDevices.Read.All`

**Optional:**
- `None` `[N/A]`

## Version History

| Version | Status | Notes |
|---------|--------|-------|
| v0.33.0-alpha | Experimental | Initial release |
| v0.40.0-alpha | Experimental | Example fixes and refactored sync progress logic |

## Example Usage

```terraform
# Example 1: Initiate on-demand proactive remediation on a single device - Minimal
action "microsoft365_graph_beta_device_management_managed_device_initiate_on_demand_proactive_remediation" "initiate_single" {
  config {
    managed_devices = [
      {
        device_id        = "12345678-1234-1234-1234-123456789abc"
        script_policy_id = "87654321-4321-4321-4321-ba9876543210"
      }
    ]
  }
}

# Example 2: Initiate proactive remediation on multiple devices
action "microsoft365_graph_beta_device_management_managed_device_initiate_on_demand_proactive_remediation" "initiate_multiple" {
  config {
    managed_devices = [
      {
        device_id        = "12345678-1234-1234-1234-123456789abc"
        script_policy_id = "87654321-4321-4321-4321-ba9876543210"
      },
      {
        device_id        = "abcdef12-3456-7890-abcd-ef1234567890"
        script_policy_id = "11111111-2222-3333-4444-555555555555"
      }
    ]

    timeouts = {
      invoke = "10m"
    }
  }
}

# Example 3: Initiate with validation - Maximal
action "microsoft365_graph_beta_device_management_managed_device_initiate_on_demand_proactive_remediation" "initiate_maximal" {
  config {
    managed_devices = [
      {
        device_id        = "12345678-1234-1234-1234-123456789abc"
        script_policy_id = "87654321-4321-4321-4321-ba9876543210"
      }
    ]

    comanaged_devices = [
      {
        device_id        = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
        script_policy_id = "bbbbbbbb-cccc-dddd-eeee-ffffffffffff"
      }
    ]

    ignore_partial_failures = true
    validate_device_exists  = true

    timeouts = {
      invoke = "5m"
    }
  }
}

# Example 4: Initiate remediation on all Windows devices with specific script
data "microsoft365_graph_beta_device_management_managed_device" "windows_devices" {
  filter_type  = "odata"
  odata_filter = "operatingSystem eq 'Windows'"
}

action "microsoft365_graph_beta_device_management_managed_device_initiate_on_demand_proactive_remediation" "initiate_all_windows" {
  config {
    managed_devices = [
      for device in data.microsoft365_graph_beta_device_management_managed_device.windows_devices.items : {
        device_id        = device.id
        script_policy_id = "87654321-4321-4321-4321-ba9876543210"
      }
    ]

    validate_device_exists = true

    timeouts = {
      invoke = "30m"
    }
  }
}
```

<!-- action schema generated by tfplugindocs -->
## Schema

### Optional

- `comanaged_devices` (Attributes List) List of co-managed devices to initiate proactive remediation for. These are devices managed by both Intune and Configuration Manager (SCCM).

**Note:** At least one of `managed_devices` or `comanaged_devices` must be provided. (see [below for nested schema](#nestedatt--comanaged_devices))
- `ignore_partial_failures` (Boolean) If set to `true`, the action will succeed even if some operations fail. Failed operations will be reported as warnings instead of errors. Default: `false` (action fails if any operation fails).
- `managed_devices` (Attributes List) List of managed devices to initiate proactive remediation for. Each entry specifies a device and the remediation script to run.

**Note:** At least one of `managed_devices` or `comanaged_devices` must be provided. Each device can have a different script policy executed.

**Important:** The script policy must already be deployed to the device. This action triggers immediate execution but does not create a new deployment. (see [below for nested schema](#nestedatt--managed_devices))
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))
- `validate_device_exists` (Boolean) Whether to validate that devices exist and are Windows devices before attempting remediation. Disabling this can speed up planning but may result in runtime errors for non-existent or unsupported devices. Default: `true`.

<a id="nestedatt--comanaged_devices"></a>
### Nested Schema for `comanaged_devices`

Required:

- `device_id` (String) The unique identifier (GUID) of the co-managed device to run the remediation script on.

**Example**: `"12345678-1234-1234-1234-123456789abc"`
- `script_policy_id` (String) The unique identifier (GUID) of the proactive remediation script policy to execute.

**Note**: The script must already be assigned/deployed to the device.

**Example**: `"87654321-4321-4321-4321-ba9876543210"`


<a id="nestedatt--managed_devices"></a>
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

- `invoke` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).

