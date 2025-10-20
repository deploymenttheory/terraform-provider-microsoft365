---
page_title: "microsoft365_graph_beta_device_management_managed_device_move_devices_to_ou Action - terraform-provider-microsoft365"
subcategory: "Device Management"

description: |-
  Moves hybrid Azure AD joined Windows devices to a specified Active Directory Organizational Unit (OU) using the /deviceManagement/managedDevices/moveDevicesToOU and /deviceManagement/comanagedDevices/moveDevicesToOU endpoints. This action updates the organizational unit placement of devices in on-premises Active Directory for hybrid-joined devices. The move operation is performed at the collection level, allowing multiple devices to be moved to the same OU in a single operation.
  Important Notes:
  Only works on Hybrid Azure AD joined Windows devicesRequires on-premises Active Directory connectivityRequires Azure AD Connect syncAll devices are moved to the same OU pathOU path must be valid in on-premises ADChanges reflect after next Azure AD Connect syncDoes not affect cloud-only or Workplace-joined devices
  Use Cases:
  Reorganizing device structure in Active DirectoryApplying different Group Policy Objects (GPOs)Moving devices between departments or locationsAligning device placement with organizational structureConsolidating devices for management purposesPreparing devices for different security policies
  Platform Support:
  Windows: Hybrid Azure AD joined devices onlyOther Platforms: Not supported (cloud-only management)
  Reference: Microsoft Graph API - Move Devices to OU https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-movedevicestoou?view=graph-rest-beta
---

# microsoft365_graph_beta_device_management_managed_device_move_devices_to_ou (Action)

Moves hybrid Azure AD joined Windows devices to a specified Active Directory Organizational Unit (OU) using the `/deviceManagement/managedDevices/moveDevicesToOU` and `/deviceManagement/comanagedDevices/moveDevicesToOU` endpoints. This action updates the organizational unit placement of devices in on-premises Active Directory for hybrid-joined devices. The move operation is performed at the collection level, allowing multiple devices to be moved to the same OU in a single operation.

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
| **Windows** | ✅ Full Support | Hybrid Azure AD joined devices with on-premises AD connectivity |
| **macOS** | ❌ Not Supported | Cloud-only management, no Active Directory integration |
| **iOS/iPadOS** | ❌ Not Supported | Cloud-only management |
| **Android** | ❌ Not Supported | Cloud-only management |

### What is Moving Devices to OU?

Moving devices to an Organizational Unit is an action that:
- Updates the Active Directory location of hybrid Azure AD joined Windows devices
- Changes the OU placement of computer objects in on-premises Active Directory
- Enables different Group Policy Objects (GPOs) to be applied based on new OU
- Reorganizes device structure to align with organizational hierarchy
- Does not affect Intune MDM management or device enrollment
- Reflects in Azure AD after next Azure AD Connect synchronization cycle

### Collection-Level Operation

This action operates at the **collection level**, which means:

| Aspect | Behavior |
|--------|----------|
| **API Call** | Single API call moves all specified devices |
| **Target OU** | All devices moved to the **same** OU path |
| **Efficiency** | More efficient than individual per-device moves |
| **Execution** | Batch operation processed together |
| **Result** | All devices succeed or fail as a group |

If you need to move devices to **different** OUs, you must create **separate actions** for each target OU.

### Requirements

**Active Directory Requirements:**
- On-premises Active Directory domain controller accessible
- Target OU must exist in Active Directory before the move
- Valid distinguished name (DN) format for OU path

**Azure AD Requirements:**
- Devices must be **hybrid Azure AD joined** (not cloud-only or workplace-joined)
- Azure AD Connect must be configured and actively syncing
- Sync schedule must be operational (default: 30 minutes)

**Permissions Requirements:**
- Azure AD Connect service account must have permissions to move computer objects
- Service account needs write access to target OU
- Service account needs read access to source OU

**Device Requirements:**
- Windows operating system only
- Hybrid Azure AD joined status
- Active enrollment in Intune

### When to Move Devices to OU

- Reorganizing device structure in Active Directory
- Applying different Group Policy Objects (GPOs) to device groups
- Moving devices between departments or organizational units
- Aligning device placement with organizational changes
- Implementing new security policies via GPO
- Consolidating devices for centralized management
- Responding to location or department transfers
- Preparing devices for different access control policies

### What Happens When Devices are Moved

1. **Immediate**: API request is processed by Microsoft Graph
2. **Active Directory**: Computer object is moved to new OU in on-premises AD
3. **Azure AD Connect**: Change is detected during next sync cycle (typically 30 minutes)
4. **Azure AD**: Device object is updated to reflect new OU placement
5. **Group Policy**: New OU's GPOs begin applying (computer refresh required for some settings)
6. **User Impact**: Minimal to none - device continues normal operation

### OU Path Format

The `organizational_unit_path` must be a valid Active Directory **distinguished name (DN)**:

**Format**: `OU=Name,OU=Parent,DC=domain,DC=com`

**Valid Examples**:
- `OU=Workstations,DC=contoso,DC=com`
- `OU=Laptops,OU=Mobile,DC=corp,DC=acme,DC=com`
- `OU=Finance,OU=Departments,DC=company,DC=local`
- `OU=Secure,OU=Security,OU=IT,DC=enterprise,DC=com`

**Important Notes**:
- Path is case-sensitive
- Must include full distinguished name from OU to DC components
- Cannot use LDAP shortcuts or abbreviations
- Must match exactly as it appears in Active Directory

### Group Policy Application

After moving devices to a new OU:

| GPO Type | Application Timing | Refresh Method |
|----------|-------------------|----------------|
| **Computer Settings** | Next computer startup or background refresh (90-120 min) | `gpupdate /force /boot` |
| **User Settings** | Next user logon or background refresh (90-120 min) | `gpupdate /force /logoff` |
| **Security Settings** | Requires restart for full application | Restart device |

### Azure AD Connect Sync

Changes reflect in Azure AD based on sync schedule:

| Sync Type | Default Interval | Scope |
|-----------|------------------|-------|
| **Delta Sync** | 30 minutes | Changed objects only |
| **Full Sync** | Varies (typically manual) | All objects |

Monitor sync status in:
- Azure AD Connect application on sync server
- Azure AD portal → Azure AD Connect → Sync Status
- Event Viewer on Azure AD Connect server

## Example Usage

```terraform
# Example 1: Move single device to workstations OU
action "microsoft365_graph_beta_device_management_managed_device_move_devices_to_ou" "move_to_workstations" {
  organizational_unit_path = "OU=Workstations,OU=Computers,DC=contoso,DC=com"
  managed_device_ids       = ["12345678-1234-1234-1234-123456789abc"]

  timeouts = {
    invoke = "5m"
  }
}

# Example 2: Move multiple devices to same OU
action "microsoft365_graph_beta_device_management_managed_device_move_devices_to_ou" "move_to_marketing" {
  organizational_unit_path = "OU=Marketing,OU=Departments,DC=contoso,DC=com"
  managed_device_ids = [
    "12345678-1234-1234-1234-123456789abc",
    "87654321-4321-4321-4321-ba9876543210",
    "abcdef12-3456-7890-abcd-ef1234567890"
  ]

  timeouts = {
    invoke = "10m"
  }
}

# Example 3: Move devices by department
variable "marketing_devices" {
  description = "Device IDs for marketing department"
  type        = list(string)
  default = [
    "aaaa1111-1111-1111-1111-111111111111",
    "bbbb2222-2222-2222-2222-222222222222"
  ]
}

action "microsoft365_graph_beta_device_management_managed_device_move_devices_to_ou" "marketing_org_unit" {
  organizational_unit_path = "OU=Marketing,OU=Departments,DC=corp,DC=example,DC=com"
  managed_device_ids       = var.marketing_devices

  timeouts = {
    invoke = "15m"
  }
}

# Example 4: Move devices based on data source filter
data "microsoft365_graph_beta_device_management_managed_device" "windows_devices" {
  filter_type  = "odata"
  odata_filter = "operatingSystem eq 'Windows' and deviceCategoryDisplayName eq 'Relocate'"
}

action "microsoft365_graph_beta_device_management_managed_device_move_devices_to_ou" "relocate_devices" {
  organizational_unit_path = "OU=NewLocation,OU=Offices,DC=contoso,DC=com"
  managed_device_ids       = [for device in data.microsoft365_graph_beta_device_management_managed_device.windows_devices.items : device.id]

  timeouts = {
    invoke = "20m"
  }
}

# Example 5: Move laptops to mobile OU
locals {
  laptop_devices = [
    "11111111-1111-1111-1111-111111111111",
    "22222222-2222-2222-2222-222222222222",
    "33333333-3333-3333-3333-333333333333"
  ]
}

action "microsoft365_graph_beta_device_management_managed_device_move_devices_to_ou" "laptops_to_mobile_ou" {
  organizational_unit_path = "OU=Laptops,OU=Mobile,OU=Devices,DC=corp,DC=acme,DC=com"
  managed_device_ids       = local.laptop_devices

  timeouts = {
    invoke = "15m"
  }
}

# Example 6: Move co-managed devices
action "microsoft365_graph_beta_device_management_managed_device_move_devices_to_ou" "move_comanaged" {
  organizational_unit_path = "OU=CoManaged,OU=SCCM,DC=contoso,DC=com"
  comanaged_device_ids     = ["abcdef12-3456-7890-abcd-ef1234567890"]

  timeouts = {
    invoke = "5m"
  }
}

# Example 7: Move devices to apply different GPOs
action "microsoft365_graph_beta_device_management_managed_device_move_devices_to_ou" "secure_workstations_gpo" {
  organizational_unit_path = "OU=SecureWorkstations,OU=Security,DC=contoso,DC=com"
  managed_device_ids = [
    "secure01-1111-1111-1111-111111111111",
    "secure02-2222-2222-2222-222222222222"
  ]

  timeouts = {
    invoke = "10m"
  }
}

# Example 8: Move by location/office
locals {
  office_locations = {
    "seattle" = {
      ou_path = "OU=Seattle,OU=Offices,DC=corp,DC=contoso,DC=com"
      devices = [
        "sea001-1111-1111-1111-111111111111",
        "sea002-2222-2222-2222-222222222222"
      ]
    }
    "portland" = {
      ou_path = "OU=Portland,OU=Offices,DC=corp,DC=contoso,DC=com"
      devices = [
        "pdx001-3333-3333-3333-333333333333",
        "pdx002-4444-4444-4444-444444444444"
      ]
    }
  }
}

action "microsoft365_graph_beta_device_management_managed_device_move_devices_to_ou" "seattle_office" {
  organizational_unit_path = local.office_locations["seattle"].ou_path
  managed_device_ids       = local.office_locations["seattle"].devices

  timeouts = {
    invoke = "15m"
  }
}

action "microsoft365_graph_beta_device_management_managed_device_move_devices_to_ou" "portland_office" {
  organizational_unit_path = local.office_locations["portland"].ou_path
  managed_device_ids       = local.office_locations["portland"].devices

  timeouts = {
    invoke = "15m"
  }
}

# Example 9: Move devices for compliance/security
data "microsoft365_graph_beta_device_management_managed_device" "quarantine_devices" {
  filter_type  = "odata"
  odata_filter = "deviceCategoryDisplayName eq 'Quarantine'"
}

action "microsoft365_graph_beta_device_management_managed_device_move_devices_to_ou" "quarantine_ou" {
  organizational_unit_path = "OU=Quarantine,OU=Security,DC=corp,DC=contoso,DC=com"
  managed_device_ids       = [for device in data.microsoft365_graph_beta_device_management_managed_device.quarantine_devices.items : device.id]

  timeouts = {
    invoke = "10m"
  }
}

# Example 10: Multi-tier OU structure
action "microsoft365_graph_beta_device_management_managed_device_move_devices_to_ou" "finance_restricted" {
  organizational_unit_path = "OU=Restricted,OU=Finance,OU=Departments,DC=corp,DC=contoso,DC=com"
  managed_device_ids = [
    "fin001-1111-1111-1111-111111111111",
    "fin002-2222-2222-2222-222222222222"
  ]

  timeouts = {
    invoke = "15m"
  }
}

# Output examples
output "moved_devices_summary" {
  value = {
    marketing = {
      ou_path      = action.move_to_marketing.organizational_unit_path
      device_count = length(action.move_to_marketing.managed_device_ids)
    }
    laptops = {
      ou_path      = action.laptops_to_mobile_ou.organizational_unit_path
      device_count = length(action.laptops_to_mobile_ou.managed_device_ids)
    }
  }
  description = "Summary of devices moved to OUs"
}
```

<!-- action schema generated by tfplugindocs -->
## Schema

### Required

- `organizational_unit_path` (String) The full distinguished name path of the target Organizational Unit in Active Directory. All specified devices will be moved to this OU.

**Format**: Must be a valid Active Directory OU distinguished name.

**Examples**:
- `"OU=Workstations,OU=Computers,DC=contoso,DC=com"`
- `"OU=Marketing,OU=Departments,DC=example,DC=local"`
- `"OU=Laptops,OU=Mobile,OU=Devices,DC=corp,DC=acme,DC=com"`

**Important**: The OU must exist in your on-premises Active Directory, and the Azure AD Connect sync account must have permissions to move computer objects to this OU.

### Optional

- `comanaged_device_ids` (List of String) List of co-managed device IDs (GUIDs) to move to the specified Organizational Unit. These are devices managed by both Intune and Configuration Manager (SCCM) that are hybrid Azure AD joined.

**Note:** At least one of `managed_device_ids` or `comanaged_device_ids` must be provided. All devices in this list will be moved to the same OU path.

Example: `["abcdef12-3456-7890-abcd-ef1234567890"]`
- `managed_device_ids` (List of String) List of managed device IDs (GUIDs) to move to the specified Organizational Unit. These are devices fully managed by Intune that are also hybrid Azure AD joined.

**Note:** At least one of `managed_device_ids` or `comanaged_device_ids` must be provided. All devices in this list will be moved to the same OU path specified in `organizational_unit_path`.

**Important:** Only hybrid Azure AD joined Windows devices can be moved. Cloud-only or workplace-joined devices will be ignored.

Example: `["12345678-1234-1234-1234-123456789abc", "87654321-4321-4321-4321-ba9876543210"]`
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).

