---
page_title: "microsoft365_graph_beta_device_management_managed_device_reenable Action - terraform-provider-microsoft365"
subcategory: "Device Management"

description: |-
  Re-enables previously disabled managed devices in Microsoft Intune using the /deviceManagement/managedDevices/{managedDeviceId}/reenable and /deviceManagement/comanagedDevices/{managedDeviceId}/reenable endpoints. This action is used to restore a disabled device's ability to interact with Intune services, allowing it to sync and receive policy updates again. Re-enabling is the counterpart to the disable action and restores full management capabilities to devices that were temporarily suspended. This is useful after resolving security incidents, compliance violations, or completing investigations that required temporary device suspension.
  Important Notes:
  Only works on previously disabled devicesRestores sync capability with IntuneRe-enables policy applicationMaintains existing enrollmentReverses the disable actionAll platforms supported
  Use Cases:
  Restoring devices after security investigation completionRe-enabling compliant devices after violations resolvedEnding temporary quarantine periodResuming management after troubleshootingRestoring devices after policy fixesCompleting incident response procedures
  Platform Support:
  All Platforms: Windows, macOS, iOS/iPadOS, Android
  Reference: Microsoft Graph API - Reenable https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-reenable?view=graph-rest-beta
---

# microsoft365_graph_beta_device_management_managed_device_reenable (Action)

Re-enables previously disabled managed devices in Microsoft Intune using the `/deviceManagement/managedDevices/{managedDeviceId}/reenable` and `/deviceManagement/comanagedDevices/{managedDeviceId}/reenable` endpoints. This action is used to restore a disabled device's ability to interact with Intune services, allowing it to sync and receive policy updates again. Re-enabling is the counterpart to the disable action and restores full management capabilities to devices that were temporarily suspended. This is useful after resolving security incidents, compliance violations, or completing investigations that required temporary device suspension.

**Important Notes:**
- Only works on previously disabled devices
- Restores sync capability with Intune
- Re-enables policy application
- Maintains existing enrollment
- Reverses the disable action
- All platforms supported

**Use Cases:**
- Restoring devices after security investigation completion
- Re-enabling compliant devices after violations resolved
- Ending temporary quarantine period
- Resuming management after troubleshooting
- Restoring devices after policy fixes
- Completing incident response procedures

**Platform Support:**
- **All Platforms**: Windows, macOS, iOS/iPadOS, Android

**Reference:** [Microsoft Graph API - Reenable](https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-reenable?view=graph-rest-beta)

## Microsoft Documentation

### Graph API References
- [reenable action](https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-reenable?view=graph-rest-beta)
- [managedDevice resource type](https://learn.microsoft.com/en-us/graph/api/resources/intune-devices-manageddevice?view=graph-rest-beta)

### Intune Device Management Guides
- [Remote actions in Microsoft Intune](https://learn.microsoft.com/en-us/mem/intune/remote-actions/device-management)
- [Device compliance in Intune](https://learn.microsoft.com/en-us/mem/intune/protect/device-compliance-get-started)

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
# Example 1: Re-enable a single device - Minimal
action "microsoft365_graph_beta_device_management_managed_device_reenable" "reenable_single" {
  config {
    managed_device_ids = [
      "12345678-1234-1234-1234-123456789abc"
    ]
  }
}

# Example 2: Re-enable multiple devices
action "microsoft365_graph_beta_device_management_managed_device_reenable" "reenable_multiple" {
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

# Example 3: Re-enable with validation - Maximal
action "microsoft365_graph_beta_device_management_managed_device_reenable" "reenable_with_validation" {
  config {
    managed_device_ids = [
      "12345678-1234-1234-1234-123456789abc",
      "87654321-4321-4321-4321-ba9876543210"
    ]

    comanaged_device_ids = [
      "abcdef12-3456-7890-abcd-ef1234567890"
    ]

    ignore_partial_failures = true
    validate_device_exists  = true

    timeouts = {
      invoke = "5m"
    }
  }
}

# Example 4: Re-enable devices after security investigation
variable "investigated_devices" {
  description = "Device IDs cleared from security investigation"
  type        = list(string)
  default = [
    "aaaa1111-1111-1111-1111-111111111111",
    "bbbb2222-2222-2222-2222-222222222222"
  ]
}

action "microsoft365_graph_beta_device_management_managed_device_reenable" "post_investigation" {
  config {
    managed_device_ids = var.investigated_devices

    validate_device_exists = true

    timeouts = {
      invoke = "10m"
    }
  }
}

# Example 5: Re-enable compliant devices
data "microsoft365_graph_beta_device_management_managed_device" "now_compliant" {
  filter_type  = "odata"
  odata_filter = "complianceState eq 'compliant'"
}

action "microsoft365_graph_beta_device_management_managed_device_reenable" "compliance_restored" {
  config {
    managed_device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.now_compliant.items : device.id]

    ignore_partial_failures = true

    timeouts = {
      invoke = "20m"
    }
  }
}

# Example 6: Re-enable co-managed device
action "microsoft365_graph_beta_device_management_managed_device_reenable" "reenable_comanaged" {
  config {
    comanaged_device_ids = [
      "abcdef12-3456-7890-abcd-ef1234567890"
    ]

    timeouts = {
      invoke = "5m"
    }
  }
}

# Output examples
output "reenabled_devices_count" {
  value = {
    managed   = length(action.microsoft365_graph_beta_device_management_managed_device_reenable.reenable_multiple.config.managed_device_ids)
    comanaged = length(action.microsoft365_graph_beta_device_management_managed_device_reenable.reenable_comanaged.config.comanaged_device_ids)
  }
  description = "Count of devices re-enabled"
}
```

<!-- action schema generated by tfplugindocs -->
## Schema

### Optional

- `comanaged_device_ids` (List of String) List of co-managed device IDs (GUIDs) to re-enable. These are devices managed by both Intune and Configuration Manager (SCCM) that were previously disabled.

**Note:** At least one of `managed_device_ids` or `comanaged_device_ids` must be provided.

Example: `["abcdef12-3456-7890-abcd-ef1234567890"]`
- `ignore_partial_failures` (Boolean) If set to `true`, the action will succeed even if some operations fail. Failed operations will be reported as warnings instead of errors. Default: `false` (action fails if any operation fails).
- `managed_device_ids` (List of String) List of managed device IDs (GUIDs) to re-enable. These are devices fully managed by Intune that were previously disabled.

**Note:** At least one of `managed_device_ids` or `comanaged_device_ids` must be provided. You can provide both to re-enable different types of devices in one action.

**Important:** Re-enabled devices will be able to sync with Intune and receive policy updates again.

Example: `["12345678-1234-1234-1234-123456789abc", "87654321-4321-4321-4321-ba9876543210"]`
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))
- `validate_device_exists` (Boolean) Whether to validate that devices exist before attempting to re-enable them. Disabling this can speed up planning but may result in runtime errors for non-existent devices. Default: `true`.

<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `invoke` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).

