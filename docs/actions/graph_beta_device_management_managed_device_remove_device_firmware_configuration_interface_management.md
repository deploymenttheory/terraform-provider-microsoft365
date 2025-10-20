---
page_title: "microsoft365_graph_beta_device_management_managed_device_remove_device_firmware_configuration_interface_management Action - terraform-provider-microsoft365"
subcategory: "Device Management"

description: |-
  Removes managed devices from Device Firmware Configuration Interface (DFCI) management using the /deviceManagement/managedDevices/{managedDeviceId}/removeDeviceFirmwareConfigurationInterfaceManagement and /deviceManagement/comanagedDevices/{managedDeviceId}/removeDeviceFirmwareConfigurationInterfaceManagement endpoints. DFCI enables Intune to manage UEFI (BIOS) settings on compatible Windows devices, providing low-level security controls. This action removes the DFCI management capability from devices, reverting them to standard Intune management without firmware-level control. After removal, the device's UEFI settings can no longer be managed remotely via Intune.
  Important Notes:
  Only works on Windows devices with DFCI-capable firmwareRequires devices enrolled with DFCI management enabledRemoves ability to manage UEFI/BIOS settings remotelyDoes not unenroll device from IntuneStandard MDM management continuesTypically used on Surface and compatible OEM devicesCannot be easily reversed
  Use Cases:
  Decommissioning devices from DFCI managementTransitioning to standard managementRemoving firmware-level security controlsPreparing devices for transfer or resaleTroubleshooting DFCI-related issuesDisabling low-level hardware management
  Platform Support:
  Windows: DFCI-capable devices only (Surface, select OEM devices)Other Platforms: Not supported (DFCI is Windows-specific)
  Reference: Microsoft Graph API - Remove Device Firmware Configuration Interface Management https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-removedevicefirmwareconfigurationinterfacemanagement?view=graph-rest-beta
---

# microsoft365_graph_beta_device_management_managed_device_remove_device_firmware_configuration_interface_management (Action)

Removes managed devices from Device Firmware Configuration Interface (DFCI) management using the `/deviceManagement/managedDevices/{managedDeviceId}/removeDeviceFirmwareConfigurationInterfaceManagement` and `/deviceManagement/comanagedDevices/{managedDeviceId}/removeDeviceFirmwareConfigurationInterfaceManagement` endpoints. DFCI enables Intune to manage UEFI (BIOS) settings on compatible Windows devices, providing low-level security controls. This action removes the DFCI management capability from devices, reverting them to standard Intune management without firmware-level control. After removal, the device's UEFI settings can no longer be managed remotely via Intune.

**Important Notes:**
- Only works on Windows devices with DFCI-capable firmware
- Requires devices enrolled with DFCI management enabled
- Removes ability to manage UEFI/BIOS settings remotely
- Does not unenroll device from Intune
- Standard MDM management continues
- Typically used on Surface and compatible OEM devices
- Cannot be easily reversed

**Use Cases:**
- Decommissioning devices from DFCI management
- Transitioning to standard management
- Removing firmware-level security controls
- Preparing devices for transfer or resale
- Troubleshooting DFCI-related issues
- Disabling low-level hardware management

**Platform Support:**
- **Windows**: DFCI-capable devices only (Surface, select OEM devices)
- **Other Platforms**: Not supported (DFCI is Windows-specific)

**Reference:** [Microsoft Graph API - Remove Device Firmware Configuration Interface Management](https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-removedevicefirmwareconfigurationinterfacemanagement?view=graph-rest-beta)

## Microsoft Documentation

### Graph API References
- [removeDeviceFirmwareConfigurationInterfaceManagement action](https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-removedevicefirmwareconfigurationinterfacemanagement?view=graph-rest-beta)
- [managedDevice resource type](https://learn.microsoft.com/en-us/graph/api/resources/intune-devices-manageddevice?view=graph-rest-beta)

### DFCI and Firmware Management Guides
- [Device Firmware Configuration Interface (DFCI) overview](https://learn.microsoft.com/en-us/mem/intune/configuration/device-firmware-configuration-interface-windows)
- [DFCI management with Intune](https://learn.microsoft.com/en-us/mem/intune/configuration/device-firmware-configuration-interface-windows-settings)
- [Surface DFCI management](https://learn.microsoft.com/en-us/surface/surface-manage-dfci-guide)

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
| **Windows** | ✅ Conditional | DFCI-capable firmware (Surface, select OEM devices) |
| **macOS** | ❌ Not Supported | DFCI is Windows-specific technology |
| **iOS/iPadOS** | ❌ Not Supported | DFCI is Windows-specific technology |
| **Android** | ❌ Not Supported | DFCI is Windows-specific technology |

### What is DFCI?

**Device Firmware Configuration Interface (DFCI)** is a technology that:
- Enables cloud-based management of UEFI (BIOS) settings on Windows devices
- Provides firmware-level security controls managed through Intune
- Works at a lower level than traditional MDM management
- Helps protect devices from unauthorized firmware changes
- Primarily available on Surface devices and select OEM devices
- Requires specific firmware support from the device manufacturer

### What is DFCI Management Removal?

Removing DFCI management is an action that:
- Disables Intune's ability to manage UEFI/BIOS settings remotely
- Removes firmware-level configuration control
- **Does not** unenroll the device from Intune
- **Does not** affect standard MDM management
- Reverts device to local firmware configuration only
- Typically cannot be easily reversed without physical device access
- Maintains all other Intune management capabilities

### DFCI Management States

| State | UEFI Remote Management | Standard MDM | Use Case |
|-------|------------------------|--------------|----------|
| **DFCI Enabled** | ✅ Yes | ✅ Yes | Full firmware + MDM control |
| **DFCI Removed** | ❌ No | ✅ Yes | Standard MDM only |
| **Unenrolled** | ❌ No | ❌ No | No Intune management |

### When to Remove DFCI Management

- **Decommissioning devices** from DFCI-enabled fleet
- **Transitioning to standard management** for devices that don't require firmware control
- **Preparing devices for transfer** to different organization or user
- **Troubleshooting DFCI-related issues** that prevent device operation
- **Removing firmware restrictions** for specific use cases
- **Preparing devices for resale** or repurposing

### What Happens When DFCI is Removed

1. **Immediate**: API request is processed by Microsoft Graph
2. **Device Contact**: Device receives DFCI removal command at next check-in
3. **Firmware Update**: Device firmware processes DFCI removal
4. **Management Change**: Remote UEFI configuration is disabled
5. **MDM Continues**: Standard Intune MDM management remains active
6. **User Impact**: Minimal - device continues normal operation

### Compatible Devices

DFCI management is primarily available on:

**Microsoft Surface Devices:**
- Surface Pro 7 and later
- Surface Laptop 3 and later
- Surface Book 3 and later
- Surface Go 2 and later
- Surface Studio 2+ and later

**Select OEM Devices:**
- Certain HP, Dell, Lenovo enterprise models with DFCI firmware support
- Check with device manufacturer for DFCI compatibility

### Important Considerations

⚠️ **Critical Points:**

1. **Irreversible Without Physical Access**: Once DFCI is removed remotely, re-enabling it may require physical device access
2. **Not the Same as Unenrollment**: Device remains enrolled in Intune with full MDM management
3. **Firmware Control Lost**: UEFI/BIOS settings can no longer be managed remotely
4. **Security Implications**: Consider impact on device security posture before removal
5. **Planning Required**: Ensure you understand why DFCI is being removed

### DFCI vs Standard MDM

| Feature | DFCI Management | Standard MDM |
|---------|----------------|--------------|
| **Operating System Settings** | ✅ Yes | ✅ Yes |
| **Applications** | ✅ Yes | ✅ Yes |
| **Policies** | ✅ Yes | ✅ Yes |
| **UEFI/BIOS Settings** | ✅ Yes | ❌ No |
| **Boot Configuration** | ✅ Yes | ❌ No |
| **Firmware Protection** | ✅ Yes | ❌ No |

### Re-enabling DFCI After Removal

To re-enable DFCI after removal:

1. **Physical Access Required**: Typically need physical device access
2. **UEFI Configuration**: May need to reset DFCI settings in UEFI
3. **Re-enrollment**: May require device re-enrollment or reset
4. **Manufacturer Support**: Consult device manufacturer documentation
5. **Complex Process**: Not a simple remote action

## Example Usage

```terraform
# Example 1: Remove DFCI management from single device
action "microsoft365_graph_beta_device_management_managed_device_remove_device_firmware_configuration_interface_management" "remove_dfci_single" {
  managed_device_ids = ["12345678-1234-1234-1234-123456789abc"]

  timeouts = {
    invoke = "5m"
  }
}

# Example 2: Remove DFCI from multiple Surface devices
action "microsoft365_graph_beta_device_management_managed_device_remove_device_firmware_configuration_interface_management" "remove_dfci_multiple_surface" {
  managed_device_ids = [
    "12345678-1234-1234-1234-123456789abc",
    "87654321-4321-4321-4321-ba9876543210",
    "abcdef12-3456-7890-abcd-ef1234567890"
  ]

  timeouts = {
    invoke = "10m"
  }
}

# Example 3: Remove DFCI from devices being decommissioned
variable "decommissioned_devices" {
  description = "Device IDs being decommissioned from DFCI management"
  type        = list(string)
  default = [
    "aaaa1111-1111-1111-1111-111111111111",
    "bbbb2222-2222-2222-2222-222222222222"
  ]
}

action "microsoft365_graph_beta_device_management_managed_device_remove_device_firmware_configuration_interface_management" "decommission_dfci" {
  managed_device_ids = var.decommissioned_devices

  timeouts = {
    invoke = "10m"
  }
}

# Example 4: Remove DFCI based on data source filter
data "microsoft365_graph_beta_device_management_managed_device" "dfci_devices_to_remove" {
  filter_type  = "odata"
  odata_filter = "model eq 'Surface Pro' and deviceCategoryDisplayName eq 'Remove DFCI'"
}

action "microsoft365_graph_beta_device_management_managed_device_remove_device_firmware_configuration_interface_management" "filtered_removal" {
  managed_device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.dfci_devices_to_remove.items : device.id]

  timeouts = {
    invoke = "20m"
  }
}

# Example 5: Transition to standard management
locals {
  transition_devices = [
    "11111111-1111-1111-1111-111111111111",
    "22222222-2222-2222-2222-222222222222",
    "33333333-3333-3333-3333-333333333333"
  ]
}

action "microsoft365_graph_beta_device_management_managed_device_remove_device_firmware_configuration_interface_management" "transition_standard" {
  managed_device_ids = local.transition_devices

  timeouts = {
    invoke = "15m"
  }
}

# Example 6: Remove DFCI from co-managed device
action "microsoft365_graph_beta_device_management_managed_device_remove_device_firmware_configuration_interface_management" "remove_comanaged_dfci" {
  comanaged_device_ids = ["abcdef12-3456-7890-abcd-ef1234567890"]

  timeouts = {
    invoke = "5m"
  }
}

# Example 7: Remove DFCI before device transfer
data "microsoft365_graph_beta_device_management_managed_device" "transfer_devices" {
  filter_type  = "odata"
  odata_filter = "deviceCategoryDisplayName eq 'Transfer'"
}

action "microsoft365_graph_beta_device_management_managed_device_remove_device_firmware_configuration_interface_management" "pre_transfer_removal" {
  managed_device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.transfer_devices.items : device.id]

  timeouts = {
    invoke = "30m"
  }
}

# Example 8: Remove DFCI from specific device models
locals {
  surface_models_map = {
    "surface_pro_7" = "11111111-1111-1111-1111-111111111111"
    "surface_pro_8" = "22222222-2222-2222-2222-222222222222"
  }
}

action "microsoft365_graph_beta_device_management_managed_device_remove_device_firmware_configuration_interface_management" "surface_models" {
  managed_device_ids = values(local.surface_models_map)

  timeouts = {
    invoke = "15m"
  }
}

# Output examples
output "dfci_removal_summary" {
  value = {
    managed   = length(action.remove_dfci_multiple_surface.managed_device_ids)
    comanaged = length(action.remove_comanaged_dfci.comanaged_device_ids)
  }
  description = "Count of devices with DFCI removed"
}
```

<!-- action schema generated by tfplugindocs -->
## Schema

### Optional

- `comanaged_device_ids` (List of String) List of co-managed device IDs (GUIDs) to remove from DFCI management. These are devices managed by both Intune and Configuration Manager (SCCM) that currently have DFCI management enabled.

**Note:** At least one of `managed_device_ids` or `comanaged_device_ids` must be provided.

Example: `["abcdef12-3456-7890-abcd-ef1234567890"]`
- `managed_device_ids` (List of String) List of managed device IDs (GUIDs) to remove from DFCI management. These are devices fully managed by Intune that currently have DFCI management enabled.

**Note:** At least one of `managed_device_ids` or `comanaged_device_ids` must be provided. You can provide both to remove DFCI management from different types of devices in one action.

**Important:** After removal, these devices will continue standard Intune MDM management but will no longer support remote UEFI/BIOS configuration through Intune.

Example: `["12345678-1234-1234-1234-123456789abc", "87654321-4321-4321-4321-ba9876543210"]`
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).

