---
page_title: "microsoft365_graph_beta_device_management_managed_device_deprovision Action - terraform-provider-microsoft365"
subcategory: "Device Management"

description: |-
  Deprovisions Windows managed devices from Intune management using the /deviceManagement/managedDevices/{managedDeviceId}/deprovision and /deviceManagement/comanagedDevices/{managedDeviceId}/deprovision endpoints. This action removes management capabilities from a device while allowing it to remain enrolled. Deprovisioning is less destructive than wiping or retiring a device, as it only removes management policies and profiles without deleting user data or removing the device entirely. This is useful when transitioning devices between management solutions or preparing devices for different management scenarios.
  Important Notes:
  Device remains enrolled in Intune after deprovisioningManagement policies and profiles are removedUser data is preserved on the deviceLess disruptive than wipe or retire actionsRequires a reason to be specified for auditingPrimarily used for Windows devices
  Use Cases:
  Transitioning devices to different management authorityPreparing devices for repurposingRemoving management overhead without data lossTroubleshooting management issuesMoving from co-management to different configuration
  Platform Support:
  Windows: Primary platform for deprovisioningOther Platforms: Limited or no support
  Reference: Microsoft Graph API - Deprovision https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-deprovision?view=graph-rest-beta
---

# microsoft365_graph_beta_device_management_managed_device_deprovision (Action)

Deprovisions Windows managed devices from Intune management using the `/deviceManagement/managedDevices/{managedDeviceId}/deprovision` and `/deviceManagement/comanagedDevices/{managedDeviceId}/deprovision` endpoints. This action removes management capabilities from a device while allowing it to remain enrolled. Deprovisioning is less destructive than wiping or retiring a device, as it only removes management policies and profiles without deleting user data or removing the device entirely. This is useful when transitioning devices between management solutions or preparing devices for different management scenarios.

**Important Notes:**
- Device remains enrolled in Intune after deprovisioning
- Management policies and profiles are removed
- User data is preserved on the device
- Less disruptive than wipe or retire actions
- Requires a reason to be specified for auditing
- Primarily used for Windows devices

**Use Cases:**
- Transitioning devices to different management authority
- Preparing devices for repurposing
- Removing management overhead without data loss
- Troubleshooting management issues
- Moving from co-management to different configuration

**Platform Support:**
- **Windows**: Primary platform for deprovisioning
- **Other Platforms**: Limited or no support

**Reference:** [Microsoft Graph API - Deprovision](https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-deprovision?view=graph-rest-beta)

## Microsoft Documentation

### Graph API References
- [deprovision action](https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-deprovision?view=graph-rest-beta)
- [managedDevice resource type](https://learn.microsoft.com/en-us/graph/api/resources/intune-devices-manageddevice?view=graph-rest-beta)

### Intune Device Management Guides
- [Deprovision devices](https://learn.microsoft.com/en-us/intune/intune-service/remote-actions/device-deprovision)

## API Permissions

The following API permissions are required in order to use this action.

### Microsoft Graph

- **Application**: `DeviceManagementConfiguration.Read.All`, `DeviceManagementManagedDevices.Read.All`
- **Delegated**: `DeviceManagementConfiguration.Read.All`, `DeviceManagementManagedDevices.Read.All`

## Version History

| Version | Status | Notes |
|---------|--------|-------|
| v0.33.0-alpha | Experimental | Initial release |
| v0.40.0-alpha | Experimental | Example fixes and refactored sync progress logic |


## Notes

### Platform Compatibility

| Platform | Support | Requirements |
|----------|---------|--------------|
| **Windows** | ✅ Primary Support | Windows 10/11 devices enrolled in Intune |
| **macOS** | ⚠️ Limited Support | Check API compatibility |
| **iOS/iPadOS** | ⚠️ Limited Support | Check API compatibility |
| **Android** | ⚠️ Limited Support | Check API compatibility |

### What is Device Deprovisioning?

Device Deprovisioning is an action that:
- Removes management policies and profiles from devices
- Maintains the device's enrollment status in Intune
- Preserves all user data on the device
- Provides a less destructive alternative to wipe or retire
- Requires a reason to be specified for auditing purposes
- Allows devices to be re-managed after deprovisioning

### Deprovision vs Retire vs Wipe

| Action | Management Removed | Enrollment Removed | User Data Preserved | Use Case |
|--------|-------------------|-------------------|---------------------|----------|
| **Deprovision** | ✅ Yes | ❌ No | ✅ Yes | Management transition, troubleshooting |
| **Retire** | ✅ Yes | ✅ Yes | ✅ Yes | Permanent device removal from management |
| **Wipe** | ✅ Yes | ✅ Yes | ❌ No | Factory reset, security incident response |

### When to Deprovision Devices

- Transitioning devices between management solutions
- Moving from Intune-only to co-management scenarios
- Troubleshooting persistent management or policy issues
- Preparing devices for repurposing or reassignment
- Removing management overhead without losing data
- Testing enrollment and management scenarios
- Changing management authority for devices

### What Happens When Device is Deprovisioned

- Intune sends deprovision command to the device
- Device removes management policies and profiles
- Configuration settings applied by Intune are removed
- Device enrollment status remains in Intune
- User data, files, and applications are preserved
- Device can be re-enrolled and managed again
- Deprovision reason is logged for auditing

## Example Usage

```terraform
# Example 1: Deprovision a single device - Minimal
action "microsoft365_graph_beta_device_management_managed_device_deprovision" "deprovision_single" {
  config {
    managed_devices = [
      {
        device_id          = "12345678-1234-1234-1234-123456789abc"
        deprovision_reason = "Device being transitioned to new management solution"
      }
    ]
  }
}

# Example 2: Deprovision multiple devices
action "microsoft365_graph_beta_device_management_managed_device_deprovision" "deprovision_multiple" {
  config {
    managed_devices = [
      {
        device_id          = "12345678-1234-1234-1234-123456789abc"
        deprovision_reason = "Device repurposing for different department"
      },
      {
        device_id          = "87654321-4321-4321-4321-ba9876543210"
        deprovision_reason = "Troubleshooting management issues"
      }
    ]

    timeouts = {
      invoke = "10m"
    }
  }
}

# Example 3: Maximal configuration with validation
action "microsoft365_graph_beta_device_management_managed_device_deprovision" "deprovision_maximal" {
  config {
    managed_devices = [
      {
        device_id          = "12345678-1234-1234-1234-123456789abc"
        deprovision_reason = "Transitioning to new management"
      },
      {
        device_id          = "87654321-4321-4321-4321-987654321cba"
        deprovision_reason = "Device repurposing"
      }
    ]

    comanaged_devices = [
      {
        device_id          = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
        deprovision_reason = "Removing co-management"
      }
    ]

    ignore_partial_failures = false
    validate_device_exists  = true

    timeouts = {
      invoke = "5m"
    }
  }
}

# Example 4: Deprovision devices by user
variable "departing_user_devices" {
  description = "Device IDs for departing user"
  type        = list(string)
  default = [
    "11111111-1111-1111-1111-111111111111",
    "22222222-2222-2222-2222-222222222222"
  ]
}

action "microsoft365_graph_beta_device_management_managed_device_deprovision" "user_departure" {
  config {
    managed_devices = [
      for device_id in var.departing_user_devices : {
        device_id          = device_id
        deprovision_reason = "User departure - removing management policies"
      }
    ]

    timeouts = {
      invoke = "15m"
    }
  }
}

# Example 5: Transition from Intune-only to co-management
data "microsoft365_graph_beta_device_management_managed_device" "transition_devices" {
  filter_type  = "odata"
  odata_filter = "deviceCategoryDisplayName eq 'Co-Management Transition'"
}

action "microsoft365_graph_beta_device_management_managed_device_deprovision" "comanagement_transition" {
  config {
    managed_devices = [
      for device in data.microsoft365_graph_beta_device_management_managed_device.transition_devices.items : {
        device_id          = device.id
        deprovision_reason = "Transitioning to co-management with Configuration Manager"
      }
    ]

    validate_device_exists = true

    timeouts = {
      invoke = "20m"
    }
  }
}

# Example 6: Deprovision co-managed device
action "microsoft365_graph_beta_device_management_managed_device_deprovision" "deprovision_comanaged" {
  config {
    comanaged_devices = [
      {
        device_id          = "abcdef12-3456-7890-abcd-ef1234567890"
        deprovision_reason = "Changing management authority to Configuration Manager only"
      }
    ]

    timeouts = {
      invoke = "5m"
    }
  }
}

# Example 7: Prepare devices for repurposing
data "microsoft365_graph_beta_device_management_managed_device" "repurpose_candidates" {
  filter_type  = "odata"
  odata_filter = "deviceCategoryDisplayName eq 'Repurpose Queue'"
}

action "microsoft365_graph_beta_device_management_managed_device_deprovision" "repurpose_prep" {
  config {
    managed_devices = [
      for device in data.microsoft365_graph_beta_device_management_managed_device.repurpose_candidates.items : {
        device_id          = device.id
        deprovision_reason = format("Repurposing device %s for new deployment", device.deviceName)
      }
    ]

    ignore_partial_failures = true

    timeouts = {
      invoke = "30m"
    }
  }
}

# Output examples
output "deprovision_summary" {
  value = {
    managed_count   = length(action.microsoft365_graph_beta_device_management_managed_device_deprovision.deprovision_multiple.config.managed_devices)
    comanaged_count = length(action.microsoft365_graph_beta_device_management_managed_device_deprovision.deprovision_comanaged.config.comanaged_devices)
  }
  description = "Count of devices deprovisioned"
}
```

<!-- action schema generated by tfplugindocs -->
## Schema

### Optional

- `comanaged_devices` (Attributes List) List of co-managed devices to deprovision. These are Windows devices managed by both Intune and Configuration Manager (SCCM).

**Note:** At least one of `managed_devices` or `comanaged_devices` must be provided. (see [below for nested schema](#nestedatt--comanaged_devices))
- `ignore_partial_failures` (Boolean) If set to `true`, the action will succeed even if some operations fail. Failed operations will be reported as warnings instead of errors. Default: `false` (action fails if any operation fails).
- `managed_devices` (Attributes List) List of managed devices to deprovision. These are Windows devices fully managed by Intune only.

**Note:** At least one of `managed_devices` or `comanaged_devices` must be provided. (see [below for nested schema](#nestedatt--managed_devices))
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))
- `validate_device_exists` (Boolean) Whether to validate that devices exist and are ChromeOS devices before attempting deprovision. Disabling this can speed up planning but may result in runtime errors for non-existent or unsupported devices. Default: `true`.

<a id="nestedatt--comanaged_devices"></a>
### Nested Schema for `comanaged_devices`

Required:

- `deprovision_reason` (String) The reason for deprovisioning this device. This is required for auditing and tracking purposes.

Examples:
- `"Transitioning from co-management"`
- `"Device repurposing"`
- `"Management authority change"`
- `device_id` (String) The unique identifier (GUID) of the co-managed device to deprovision.

Example: `"abcdef12-3456-7890-abcd-ef1234567890"`


<a id="nestedatt--managed_devices"></a>
### Nested Schema for `managed_devices`

Required:

- `deprovision_reason` (String) The reason for deprovisioning this device. This is required for auditing and tracking purposes.

Examples:
- `"Transitioning to new management solution"`
- `"Device repurposing"`
- `"Troubleshooting management issues"`
- `"User requested management removal"`
- `"Moving to co-management"`
- `device_id` (String) The unique identifier (GUID) of the managed device to deprovision.

Example: `"12345678-1234-1234-1234-123456789abc"`


<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `invoke` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).

