---
page_title: "microsoft365_graph_beta_device_management_managed_device_set_device_name Action - terraform-provider-microsoft365"
subcategory: "Device Management"

description: |-
  Sets a custom device name for managed devices in Microsoft Intune using the /deviceManagement/managedDevices/{managedDeviceId}/setDeviceName and /deviceManagement/comanagedDevices/{managedDeviceId}/setDeviceName endpoints. This action is used to assign meaningful, custom names to devices for easier identification and management in the Intune console. Device names can be used to reflect location, user, function, or organizational naming conventions. This action supports setting names on multiple devices in a single operation with per-device name customization.
  Important Notes:
  Device name length and character restrictions vary by platformSome platforms may have specific naming conventions or limitationsDevice must be online to receive the name change commandName changes may take time to reflect after device check-inEach device can have its own unique custom name
  Use Cases:
  Implementing organizational naming conventionsIdentifying devices by location (e.g., 'NYC-Floor3-Conf-01')Associating devices with users or departmentsStandardizing device names across the organizationRenaming devices after reassignment or relocation
  Platform Support:
  Windows: Fully supported with various name length restrictionsiOS/iPadOS: Supported for supervised devicesmacOS: Supported for managed devicesAndroid: Support varies by management mode
  Reference: Microsoft Graph API - Set Device Name https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-setdevicename?view=graph-rest-beta
---

# microsoft365_graph_beta_device_management_managed_device_set_device_name (Action)

Sets a custom device name for managed devices in Microsoft Intune using the `/deviceManagement/managedDevices/{managedDeviceId}/setDeviceName` and `/deviceManagement/comanagedDevices/{managedDeviceId}/setDeviceName` endpoints. This action is used to assign meaningful, custom names to devices for easier identification and management in the Intune console. Device names can be used to reflect location, user, function, or organizational naming conventions. This action supports setting names on multiple devices in a single operation with per-device name customization.

**Important Notes:**
- Device name length and character restrictions vary by platform
- Some platforms may have specific naming conventions or limitations
- Device must be online to receive the name change command
- Name changes may take time to reflect after device check-in
- Each device can have its own unique custom name

**Use Cases:**
- Implementing organizational naming conventions
- Identifying devices by location (e.g., 'NYC-Floor3-Conf-01')
- Associating devices with users or departments
- Standardizing device names across the organization
- Renaming devices after reassignment or relocation

**Platform Support:**
- **Windows**: Fully supported with various name length restrictions
- **iOS/iPadOS**: Supported for supervised devices
- **macOS**: Supported for managed devices
- **Android**: Support varies by management mode

**Reference:** [Microsoft Graph API - Set Device Name](https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-setdevicename?view=graph-rest-beta)

## Microsoft Documentation

### Graph API References
- [setDeviceName action](https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-setdevicename?view=graph-rest-beta)
- [managedDevice resource type](https://learn.microsoft.com/en-us/graph/api/resources/intune-devices-manageddevice?view=graph-rest-beta)

### Intune Remote Actions Guides
- [Device rename - Windows](https://learn.microsoft.com/en-us/intune/intune-service/remote-actions/device-rename?pivots=windows)
- [Device rename - iOS](https://learn.microsoft.com/en-us/intune/intune-service/remote-actions/device-rename?pivots=ios)
- [Device rename - macOS](https://learn.microsoft.com/en-us/intune/intune-service/remote-actions/device-rename?pivots=macos)
- [Device rename - Android](https://learn.microsoft.com/en-us/intune/intune-service/remote-actions/device-rename?pivots=android)

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
# Example 1: Set device name for a single device - Minimal
action "microsoft365_graph_beta_device_management_managed_device_set_device_name" "set_name_single" {
  config {
    managed_devices = [
      {
        device_id   = "12345678-1234-1234-1234-123456789abc"
        device_name = "NYC-Marketing-Laptop-01"
      }
    ]
  }
}

# Example 2: Set device names for multiple devices
action "microsoft365_graph_beta_device_management_managed_device_set_device_name" "set_names_multiple" {
  config {
    managed_devices = [
      {
        device_id   = "12345678-1234-1234-1234-123456789abc"
        device_name = "NYC-Floor3-Conf-Room-01"
      },
      {
        device_id   = "87654321-4321-4321-4321-ba9876543210"
        device_name = "NYC-Floor3-Conf-Room-02"
      },
      {
        device_id   = "abcdef12-3456-7890-abcd-ef1234567890"
        device_name = "NYC-Floor3-Conf-Room-03"
      }
    ]

    timeouts = {
      invoke = "10m"
    }
  }
}

# Example 3: Maximal configuration with validation
action "microsoft365_graph_beta_device_management_managed_device_set_device_name" "set_names_maximal" {
  config {
    managed_devices = [
      {
        device_id   = "12345678-1234-1234-1234-123456789abc"
        device_name = "NYC-Marketing-Laptop-01"
      },
      {
        device_id   = "87654321-4321-4321-4321-ba9876543210"
        device_name = "NYC-IT-Desktop-05"
      }
    ]

    comanaged_devices = [
      {
        device_id   = "abcdef12-3456-7890-abcd-ef1234567890"
        device_name = "NYC-HR-Laptop-03"
      }
    ]

    ignore_partial_failures = true
    validate_device_exists  = true

    timeouts = {
      invoke = "5m"
    }
  }
}

# Example 4: Rename devices based on user assignment
data "microsoft365_graph_beta_device_management_managed_device" "user_devices" {
  filter_type  = "odata"
  odata_filter = "userPrincipalName eq 'john.doe@example.com'"
}

action "microsoft365_graph_beta_device_management_managed_device_set_device_name" "rename_user_devices" {
  config {
    managed_devices = [
      for device in data.microsoft365_graph_beta_device_management_managed_device.user_devices.items : {
        device_id   = device.id
        device_name = format("JohnDoe-%s-%s", device.operatingSystem, substr(device.serialNumber, 0, 8))
      }
    ]

    timeouts = {
      invoke = "15m"
    }
  }
}

# Example 5: Standardize naming for devices by department
data "microsoft365_graph_beta_device_management_managed_device" "it_devices" {
  filter_type  = "odata"
  odata_filter = "deviceCategoryDisplayName eq 'IT Department'"
}

action "microsoft365_graph_beta_device_management_managed_device_set_device_name" "standardize_it_devices" {
  config {
    managed_devices = [
      for device in data.microsoft365_graph_beta_device_management_managed_device.it_devices.items : {
        device_id   = device.id
        device_name = format("IT-DEPT-%s", substr(device.id, 0, 8))
      }
    ]

    validate_device_exists = true

    timeouts = {
      invoke = "20m"
    }
  }
}

# Example 6: Rename devices by location
locals {
  device_locations = {
    "12345678-1234-1234-1234-123456789abc" = "NYC-Office"
    "87654321-4321-4321-4321-ba9876543210" = "LA-Office"
    "abcdef12-3456-7890-abcd-ef1234567890" = "Chicago-Office"
  }
}

action "microsoft365_graph_beta_device_management_managed_device_set_device_name" "rename_by_location" {
  config {
    managed_devices = [
      for device_id, location in local.device_locations : {
        device_id   = device_id
        device_name = format("%s-Device", location)
      }
    ]

    timeouts = {
      invoke = "10m"
    }
  }
}

# Example 7: Set name for co-managed device
action "microsoft365_graph_beta_device_management_managed_device_set_device_name" "set_comanaged_name" {
  config {
    comanaged_devices = [
      {
        device_id   = "abcdef12-3456-7890-abcd-ef1234567890"
        device_name = "SCCM-Intune-Hybrid-01"
      }
    ]

    timeouts = {
      invoke = "5m"
    }
  }
}

# Example 8: Rename devices after asset reassignment
variable "reassigned_devices" {
  description = "Map of device IDs to new names after reassignment"
  type        = map(string)
  default = {
    "11111111-1111-1111-1111-111111111111" = "Finance-Laptop-A"
    "22222222-2222-2222-2222-222222222222" = "Finance-Laptop-B"
  }
}

action "microsoft365_graph_beta_device_management_managed_device_set_device_name" "reassign_devices" {
  config {
    managed_devices = [
      for device_id, device_name in var.reassigned_devices : {
        device_id   = device_id
        device_name = device_name
      }
    ]

    timeouts = {
      invoke = "10m"
    }
  }
}

# Output examples
output "devices_renamed_count" {
  value       = length(action.microsoft365_graph_beta_device_management_managed_device_set_device_name.set_names_multiple.config.managed_devices)
  description = "Number of devices that had their names set"
}

output "device_naming_info" {
  value = {
    managed   = length(action.microsoft365_graph_beta_device_management_managed_device_set_device_name.set_names_maximal.config.managed_devices)
    comanaged = length(action.microsoft365_graph_beta_device_management_managed_device_set_device_name.set_names_maximal.config.comanaged_devices)
  }
  description = "Count of renamed devices by type"
}
```

<!-- action schema generated by tfplugindocs -->
## Schema

### Optional

- `comanaged_devices` (Attributes List) List of co-managed devices to set custom names for. These are devices managed by both Intune and Configuration Manager (SCCM). Each entry specifies a device ID and the new name to assign.

**Note:** At least one of `managed_devices` or `comanaged_devices` must be provided. (see [below for nested schema](#nestedatt--comanaged_devices))
- `ignore_partial_failures` (Boolean) When set to `true`, the action will complete successfully even if some devices fail to be renamed. When `false` (default), the action will fail if any device rename operation fails. Use this flag when renaming multiple devices and you want the action to succeed even if some renames fail.
- `managed_devices` (Attributes List) List of managed devices to set custom names for. These are devices fully managed by Intune only. Each entry specifies a device ID and the new name to assign to that device.

**Note:** At least one of `managed_devices` or `comanaged_devices` must be provided. You can provide both to set names on different types of devices in one action. (see [below for nested schema](#nestedatt--managed_devices))
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))
- `validate_device_exists` (Boolean) When set to `true` (default), the action will validate that all specified devices exist before attempting to rename them. When `false`, device validation is skipped and the action will attempt to rename devices directly. Disabling validation can improve performance but may result in errors if devices don't exist.

<a id="nestedatt--comanaged_devices"></a>
### Nested Schema for `comanaged_devices`

Required:

- `device_id` (String) The unique identifier (GUID) of the co-managed device to rename. Example: `"12345678-1234-1234-1234-123456789abc"`
- `device_name` (String) The new name to assign to this co-managed device. Example: `"NYC-IT-Desktop-05"`


<a id="nestedatt--managed_devices"></a>
### Nested Schema for `managed_devices`

Required:

- `device_id` (String) The unique identifier (GUID) of the managed device to rename. Example: `"12345678-1234-1234-1234-123456789abc"`
- `device_name` (String) The new name to assign to this device. Device naming requirements vary by platform. Consult platform-specific documentation for character and length restrictions. Example: `"NYC-Marketing-Laptop-01"`


<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `invoke` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).

