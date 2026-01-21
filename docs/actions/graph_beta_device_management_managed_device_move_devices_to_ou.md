---
page_title: "microsoft365_graph_beta_device_management_managed_device_move_devices_to_ou Action - terraform-provider-microsoft365"
subcategory: "Device Management"

description: |-
  Moves hybrid Azure AD joined Windows devices to a specified Active Directory Organizational Unit (OU) in Microsoft Intune using the /deviceManagement/managedDevices/moveDevicesToOU and /deviceManagement/comanagedDevices/moveDevicesToOU endpoints. This action is used to update the organizational unit placement of devices in on-premises Active Directory for hybrid-joined devices. The move operation is performed at the collection level, allowing multiple devices to be moved to the same OU in a single operation.
  Important Notes:
  Only works on Hybrid Azure AD joined Windows devicesRequires on-premises Active Directory connectivityRequires Azure AD Connect syncAll devices are moved to the same OU pathOU path must be valid in on-premises ADChanges reflect after next Azure AD Connect syncDoes not affect cloud-only or Workplace-joined devices
  Use Cases:
  Reorganizing device structure in Active DirectoryApplying different Group Policy Objects (GPOs)Moving devices between departments or locationsAligning device placement with organizational structureConsolidating devices for management purposesPreparing devices for different security policies
  Platform Support:
  Windows: Hybrid Azure AD joined devices onlyOther Platforms: Not supported (cloud-only management)
  Reference: Microsoft Graph API - Move Devices to OU https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-movedevicestoou?view=graph-rest-beta
---

# microsoft365_graph_beta_device_management_managed_device_move_devices_to_ou (Action)

Moves hybrid Azure AD joined Windows devices to a specified Active Directory Organizational Unit (OU) in Microsoft Intune using the `/deviceManagement/managedDevices/moveDevicesToOU` and `/deviceManagement/comanagedDevices/moveDevicesToOU` endpoints. This action is used to update the organizational unit placement of devices in on-premises Active Directory for hybrid-joined devices. The move operation is performed at the collection level, allowing multiple devices to be moved to the same OU in a single operation.

**Important Notes:**
- Only works on **Hybrid Azure AD joined** Windows devices
- Requires on-premises Active Directory connectivity
- Requires Azure AD Connect sync
- All devices are moved to the **same** OU path
- OU path must be valid in on-premises AD
- Changes reflect after next Azure AD Connect sync
- Does not affect cloud-only or Workplace-joined devices

**Use Cases:**
- Reorganizing device structure in Active Directory
- Applying different Group Policy Objects (GPOs)
- Moving devices between departments or locations
- Aligning device placement with organizational structure
- Consolidating devices for management purposes
- Preparing devices for different security policies

**Platform Support:**
- **Windows**: Hybrid Azure AD joined devices only
- **Other Platforms**: Not supported (cloud-only management)

**Reference:** [Microsoft Graph API - Move Devices to OU](https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-movedevicestoou?view=graph-rest-beta)

## Microsoft Documentation

### Graph API References
- [moveDevicesToOU action](https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-movedevicestoou?view=graph-rest-beta)
- [managedDevice resource type](https://learn.microsoft.com/en-us/graph/api/resources/intune-devices-manageddevice?view=graph-rest-beta)

### Active Directory and Hybrid Identity Guides
- [Hybrid Azure AD join](https://learn.microsoft.com/en-us/azure/active-directory/devices/hybrid-azuread-join-plan)
- [Azure AD Connect sync](https://learn.microsoft.com/en-us/azure/active-directory/hybrid/how-to-connect-sync-whatis)
- [Group Policy Overview](https://learn.microsoft.com/en-us/previous-versions/windows/it-pro/windows-server-2012-R2-and-2012/hh831791(v=ws.11))

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
# Example 1: Move devices to organizational unit - Minimal
action "microsoft365_graph_beta_device_management_managed_device_move_devices_to_ou" "move_single" {
  config {
    organizational_unit_path = "OU=Workstations,DC=contoso,DC=com"
    managed_device_ids       = ["12345678-1234-1234-1234-123456789abc"]
  }
}

# Example 2: Move multiple devices to OU
action "microsoft365_graph_beta_device_management_managed_device_move_devices_to_ou" "move_multiple" {
  config {
    organizational_unit_path = "OU=Finance,OU=Departments,DC=contoso,DC=com"
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

# Example 3: Move devices with validation - Maximal
action "microsoft365_graph_beta_device_management_managed_device_move_devices_to_ou" "move_maximal" {
  config {
    organizational_unit_path = "OU=IT,OU=Departments,DC=contoso,DC=com"
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

# Example 4: Move department devices to new OU
data "microsoft365_graph_beta_device_management_managed_device" "marketing_devices" {
  filter_type  = "odata"
  odata_filter = "deviceCategoryDisplayName eq 'Marketing'"
}

action "microsoft365_graph_beta_device_management_managed_device_move_devices_to_ou" "move_marketing" {
  config {
    organizational_unit_path = "OU=Marketing,OU=Departments,DC=contoso,DC=com"
    managed_device_ids       = [for device in data.microsoft365_graph_beta_device_management_managed_device.marketing_devices.items : device.id]

    validate_device_exists = true

    timeouts = {
      invoke = "15m"
    }
  }
}
```

<!-- action schema generated by tfplugindocs -->
## Schema

### Required

- `organizational_unit_path` (String) The full distinguished name path of the target Organizational Unit in Active Directory. All specified devices will be moved to this OU.

**Important**: The OU must exist in your on-premises Active Directory, and the Azure AD Connect sync account must have permissions to move computer objects to this OU.

### Optional

- `comanaged_device_ids` (List of String) List of co-managed device IDs (GUIDs) to move to the specified Organizational Unit. These are devices managed by both Intune and Configuration Manager (SCCM) that are hybrid Azure AD joined.

**Note:** At least one of `managed_device_ids` or `comanaged_device_ids` must be provided. All devices in this list will be moved to the same OU path.

Example: `["abcdef12-3456-7890-abcd-ef1234567890"]`
- `ignore_partial_failures` (Boolean) If set to `true`, the action will succeed even if some operations fail. Failed operations will be reported as warnings instead of errors. Default: `false` (action fails if any operation fails).
- `managed_device_ids` (List of String) List of managed device IDs (GUIDs) to move to the specified Organizational Unit. These are devices fully managed by Intune that are also hybrid Azure AD joined.

**Note:** At least one of `managed_device_ids` or `comanaged_device_ids` must be provided. All devices in this list will be moved to the same OU path specified in `organizational_unit_path`.

**Important:** Only hybrid Azure AD joined Windows devices can be moved. Cloud-only or workplace-joined devices will be ignored.

Example: `["12345678-1234-1234-1234-123456789abc", "87654321-4321-4321-4321-ba9876543210"]`
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))
- `validate_device_exists` (Boolean) Whether to validate that devices exist and are hybrid Azure AD joined Windows devices before attempting to move them. Disabling this can speed up planning but may result in runtime errors for non-existent or unsupported devices. Default: `true`.

<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `invoke` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).

