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
| v0.40.0-alpha | Experimental | Example fixes and refactored sync progress logic |


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
# Example 1: Remove DFCI management from a single device - Minimal
action "microsoft365_graph_beta_device_management_managed_device_remove_device_firmware_configuration_interface_management" "remove_single" {
  config {
    managed_device_ids = [
      "12345678-1234-1234-1234-123456789abc"
    ]
  }
}

# Example 2: Remove DFCI management from multiple devices
action "microsoft365_graph_beta_device_management_managed_device_remove_device_firmware_configuration_interface_management" "remove_multiple" {
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

# Example 3: Remove with validation - Maximal
action "microsoft365_graph_beta_device_management_managed_device_remove_device_firmware_configuration_interface_management" "remove_maximal" {
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
      invoke = "5m"
    }
  }
}

# Example 4: Remove DFCI from all Surface devices
data "microsoft365_graph_beta_device_management_managed_device" "surface_devices" {
  filter_type  = "odata"
  odata_filter = "startsWith(model, 'Surface')"
}

action "microsoft365_graph_beta_device_management_managed_device_remove_device_firmware_configuration_interface_management" "remove_all_surface" {
  config {
    managed_device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.surface_devices.items : device.id]

    validate_device_exists  = true
    ignore_partial_failures = true

    timeouts = {
      invoke = "20m"
    }
  }
}

# Example 5: Remove DFCI from devices being decommissioned
data "microsoft365_graph_beta_device_management_managed_device" "decommission_devices" {
  filter_type  = "odata"
  odata_filter = "deviceCategoryDisplayName eq 'Decommission Queue'"
}

action "microsoft365_graph_beta_device_management_managed_device_remove_device_firmware_configuration_interface_management" "remove_decommission" {
  config {
    managed_device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.decommission_devices.items : device.id]

    timeouts = {
      invoke = "15m"
    }
  }
}
```

<!-- action schema generated by tfplugindocs -->
## Schema

### Optional

- `comanaged_device_ids` (List of String) List of co-managed device IDs (GUIDs) to remove from DFCI management. These are devices managed by both Intune and Configuration Manager (SCCM) that currently have DFCI management enabled.

**Note:** At least one of `managed_device_ids` or `comanaged_device_ids` must be provided.

Example: `["abcdef12-3456-7890-abcd-ef1234567890"]`
- `ignore_partial_failures` (Boolean) If set to `true`, the action will succeed even if some operations fail. Failed operations will be reported as warnings instead of errors. Default: `false` (action fails if any operation fails).
- `managed_device_ids` (List of String) List of managed device IDs (GUIDs) to remove from DFCI management. These are devices fully managed by Intune that currently have DFCI management enabled.

**Note:** At least one of `managed_device_ids` or `comanaged_device_ids` must be provided. You can provide both to remove DFCI management from different types of devices in one action.

**Important:** After removal, these devices will continue standard Intune MDM management but will no longer support remote UEFI/BIOS configuration through Intune.

Example: `["12345678-1234-1234-1234-123456789abc", "87654321-4321-4321-4321-ba9876543210"]`
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))
- `validate_device_exists` (Boolean) Whether to validate that devices exist and are Windows devices before attempting DFCI removal. Disabling this can speed up planning but may result in runtime errors for non-existent or non-Windows devices. Default: `true`.

<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `invoke` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).

