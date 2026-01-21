---
page_title: "microsoft365_graph_beta_device_management_managed_device_disable Action - terraform-provider-microsoft365"
subcategory: "Device Management"

description: |-
  Disables managed devices in Microsoft Intune using the /deviceManagement/managedDevices/{managedDeviceId}/disable and /deviceManagement/comanagedDevices/{managedDeviceId}/disable endpoints. This action is used to disable devices from Intune management. This action disables a device's ability to interact with Intune services while maintaining its enrollment record. Disabled devices cannot receive policies, sync with Intune, or perform managed operations until re-enabled. This is useful for temporarily suspending device management without fully removing the device from Intune, such as during investigations, compliance violations, or security incidents.
  Important Notes:
  Device remains enrolled but cannot sync or receive policiesManagement operations are suspendedDevice can be re-enabled laterLess permanent than retire or wipeUseful for temporary suspensionsSecurity and compliance enforcement
  Use Cases:
  Security incident response (suspected compromise)Compliance violations requiring device suspensionTemporary device quarantineInvestigation of device issuesPreventing policy application during troubleshootingTemporary management suspension
  Platform Support:
  All Platforms: Windows, macOS, iOS/iPadOS, Android
  Reference: Microsoft Graph API - Disable https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-disable?view=graph-rest-beta
---

# microsoft365_graph_beta_device_management_managed_device_disable (Action)

Disables managed devices in Microsoft Intune using the `/deviceManagement/managedDevices/{managedDeviceId}/disable` and `/deviceManagement/comanagedDevices/{managedDeviceId}/disable` endpoints. This action is used to disable devices from Intune management. This action disables a device's ability to interact with Intune services while maintaining its enrollment record. Disabled devices cannot receive policies, sync with Intune, or perform managed operations until re-enabled. This is useful for temporarily suspending device management without fully removing the device from Intune, such as during investigations, compliance violations, or security incidents.

**Important Notes:**
- Device remains enrolled but cannot sync or receive policies
- Management operations are suspended
- Device can be re-enabled later
- Less permanent than retire or wipe
- Useful for temporary suspensions
- Security and compliance enforcement

**Use Cases:**
- Security incident response (suspected compromise)
- Compliance violations requiring device suspension
- Temporary device quarantine
- Investigation of device issues
- Preventing policy application during troubleshooting
- Temporary management suspension

**Platform Support:**
- **All Platforms**: Windows, macOS, iOS/iPadOS, Android

**Reference:** [Microsoft Graph API - Disable](https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-disable?view=graph-rest-beta)

## Microsoft Documentation

### Graph API References
- [disable action](https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-disable?view=graph-rest-beta)
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
# Example 1: Disable a single device - Minimal
action "microsoft365_graph_beta_device_management_managed_device_disable" "disable_single" {
  config {
    managed_device_ids = [
      "12345678-1234-1234-1234-123456789abc"
    ]
  }
}

# Example 2: Disable multiple devices
action "microsoft365_graph_beta_device_management_managed_device_disable" "disable_multiple" {
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

# Example 3: Disable with validation - Maximal
action "microsoft365_graph_beta_device_management_managed_device_disable" "disable_with_validation" {
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

# Example 4: Disable devices due to security incident
variable "compromised_devices" {
  description = "Device IDs suspected of compromise"
  type        = list(string)
  default = [
    "aaaa1111-1111-1111-1111-111111111111",
    "bbbb2222-2222-2222-2222-222222222222"
  ]
}

action "microsoft365_graph_beta_device_management_managed_device_disable" "security_incident" {
  config {
    managed_device_ids = var.compromised_devices

    validate_device_exists  = true
    ignore_partial_failures = false

    timeouts = {
      invoke = "10m"
    }
  }
}

# Example 5: Disable non-compliant devices
data "microsoft365_graph_beta_device_management_managed_device" "non_compliant" {
  filter_type  = "odata"
  odata_filter = "complianceState eq 'noncompliant'"
}

action "microsoft365_graph_beta_device_management_managed_device_disable" "compliance_enforcement" {
  config {
    managed_device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.non_compliant.items : device.id]

    ignore_partial_failures = true

    timeouts = {
      invoke = "20m"
    }
  }
}

# Example 6: Disable co-managed device
action "microsoft365_graph_beta_device_management_managed_device_disable" "disable_comanaged" {
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
output "disabled_devices_count" {
  value = {
    managed   = length(action.microsoft365_graph_beta_device_management_managed_device_disable.disable_multiple.config.managed_device_ids)
    comanaged = length(action.microsoft365_graph_beta_device_management_managed_device_disable.disable_comanaged.config.comanaged_device_ids)
  }
  description = "Count of devices disabled"
}
```

<!-- action schema generated by tfplugindocs -->
## Schema

### Optional

- `comanaged_device_ids` (List of String) List of co-managed device IDs (GUIDs) to disable. These are devices managed by both Intune and Configuration Manager (SCCM).

**Note:** At least one of `managed_device_ids` or `comanaged_device_ids` must be provided.

Example: `["abcdef12-3456-7890-abcd-ef1234567890"]`
- `ignore_partial_failures` (Boolean) If set to `true`, the action will succeed even if some operations fail. Failed operations will be reported as warnings instead of errors. Default: `false` (action fails if any operation fails).
- `managed_device_ids` (List of String) List of managed device IDs (GUIDs) to disable. These are devices fully managed by Intune only.

**Note:** At least one of `managed_device_ids` or `comanaged_device_ids` must be provided. You can provide both to disable different types of devices in one action.

**Important:** Disabled devices will not be able to sync with Intune or receive policy updates until they are re-enabled.

Example: `["12345678-1234-1234-1234-123456789abc", "87654321-4321-4321-4321-ba9876543210"]`
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))
- `validate_device_exists` (Boolean) Whether to validate that devices exist before attempting disable. Disabling this can speed up planning but may result in runtime errors for non-existent devices. Default: `true`.

<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `invoke` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).

