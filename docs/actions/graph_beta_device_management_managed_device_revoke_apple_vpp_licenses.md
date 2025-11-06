---
page_title: "Microsoft 365_microsoft365_graph_beta_device_management_managed_device_revoke_apple_vpp_licenses Action"
subcategory: "Device Management"
description: |-
  Revokes all Apple Volume Purchase Program (VPP) licenses from devices using the /deviceManagement/managedDevices/{managedDeviceId}/revokeAppleVppLicenses and /deviceManagement/comanagedDevices/{managedDeviceId}/revokeAppleVppLicenses endpoints. This action reclaims all VPP-purchased app licenses assigned to iOS/iPadOS devices, making them available for reassignment to other devices or users.
  What This Action Does:
  Revokes all VPP app licenses from deviceReturns licenses to available poolMakes licenses available for reassignmentRemoves apps from device (if enforced)Updates license inventoryAudits license revocation
  When to Use:
  Device retirement or decommissioningDevice lost or stolenUser departure from organizationLicense reallocation neededDevice platform changeLicense optimizationCompliance requirements
  Platform Support:
  iOS: Full support (VPP apps)iPadOS: Full support (VPP apps)Other platforms: Not applicable (no VPP)
  Important Considerations:
  Only affects VPP-purchased appsDevice may need to be onlineApps may be removed from deviceUser data typically preservedLicenses immediately availableCannot be undone easily
  Reference: Microsoft Graph API - Revoke Apple VPP Licenses https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-revokeapplevpplicenses?view=graph-rest-beta
---

# Microsoft 365_microsoft365_graph_beta_device_management_managed_device_revoke_apple_vpp_licenses (Action)

Revokes all Apple Volume Purchase Program (VPP) licenses from devices using the `/deviceManagement/managedDevices/{managedDeviceId}/revokeAppleVppLicenses` and `/deviceManagement/comanagedDevices/{managedDeviceId}/revokeAppleVppLicenses` endpoints. This action reclaims all VPP-purchased app licenses assigned to iOS/iPadOS devices, making them available for reassignment to other devices or users.

**What This Action Does:**
- Revokes all VPP app licenses from device
- Returns licenses to available pool
- Makes licenses available for reassignment
- Removes apps from device (if enforced)
- Updates license inventory
- Audits license revocation

**When to Use:**
- Device retirement or decommissioning
- Device lost or stolen
- User departure from organization
- License reallocation needed
- Device platform change
- License optimization
- Compliance requirements

**Platform Support:**
- **iOS**: Full support (VPP apps)
- **iPadOS**: Full support (VPP apps)
- **Other platforms**: Not applicable (no VPP)

**Important Considerations:**
- Only affects VPP-purchased apps
- Device may need to be online
- Apps may be removed from device
- User data typically preserved
- Licenses immediately available
- Cannot be undone easily

**Reference:** [Microsoft Graph API - Revoke Apple VPP Licenses](https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-revokeapplevpplicenses?view=graph-rest-beta)

## Use Cases

This action is essential for managing Apple Volume Purchase Program (VPP) licenses across iOS and iPadOS devices:

### License Management
- **Device Retirement**: Reclaim VPP licenses from devices being retired, decommissioned, or returned to inventory
- **User Departure**: Recover VPP app licenses when employees or students leave the organization
- **License Reallocation**: Free up licenses from inactive or underutilized devices for reassignment to active users
- **Device Platform Change**: Reclaim licenses when replacing iOS/iPadOS devices with other platforms
- **License Optimization**: Identify and revoke licenses from devices that no longer need specific VPP apps
- **Compliance Requirements**: Ensure license compliance by revoking licenses from non-compliant or unsupported devices
- **Emergency Scenarios**: Quickly revoke licenses from lost or stolen devices to prevent unauthorized app usage

### IT Operations
- **Bulk License Recovery**: Perform large-scale license reclamation across multiple devices simultaneously
- **Automated License Cleanup**: Schedule regular license audits to recover licenses from inactive devices
- **Cost Optimization**: Maximize ROI on VPP app purchases by ensuring licenses are actively used
- **Inventory Management**: Maintain accurate license inventory and availability tracking
- **Deployment Preparation**: Clear old licenses before deploying new VPP app assignments

### Educational Institutions
- **Semester End Cleanup**: Revoke licenses from graduating students' devices or devices returned at term end
- **Device Reassignment**: Prepare devices for new students by revoking previous VPP app licenses
- **Shared Device Programs**: Manage VPP licenses on shared iPads between different students or classes
- **Lab Equipment Management**: Reclaim licenses from lab devices when upgrading or retiring equipment

## API Documentation

- [Microsoft Graph API - Revoke Apple VPP Licenses](https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-revokeapplevpplicenses?view=graph-rest-beta)

## Permissions

The following Microsoft Graph API permissions are required to use this action:

| Permission Type | Permissions (Least Privileged) |
|:----------------|:------------------------------|
| Delegated (work or school account) | DeviceManagementManagedDevices.PrivilegedOperations.All |
| Delegated (personal Microsoft account) | Not supported |
| Application | DeviceManagementManagedDevices.PrivilegedOperations.All |

~> **Note:** This action requires privileged operations permissions as it directly affects app license assignments and availability.

## Related Documentation

- [Microsoft Intune Remote Actions - iOS/iPadOS](https://learn.microsoft.com/en-us/intune/intune-service/remote-actions/?tabs=ios-ipados)
- [Microsoft Intune Remote Actions - macOS](https://learn.microsoft.com/en-us/intune/intune-service/remote-actions/?tabs=macos)

## Version History

| Version | Status | Notes |
|---------|--------|-------|
| v0.33.0-alpha | Experimental | Initial release |

## Notes

### Platform Compatibility

This remote action is only available for iOS and iPadOS devices with VPP-purchased apps. The table below shows platform support:

| Platform | Supported | Notes |
|:---------|:----------|:------|
| **Windows** | ❌ | Not supported - VPP is Apple-specific |
| **macOS** | ❌ | Not supported - Uses different licensing model |
| **iOS** | ✅ | Fully supported for VPP apps |
| **iPadOS** | ✅ | Fully supported for VPP apps |
| **Android** | ❌ | Not supported - Uses Google Play licensing |
| **Android Enterprise** | ❌ | Not supported - Uses managed Google Play |
| **ChromeOS** | ❌ | Not supported - Uses Chrome Web Store licensing |


## Example Usage

```terraform
# REF: https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-revokeapplevpplicenses?view=graph-rest-beta

# Data source to find iOS devices with VPP apps assigned
data "microsoft365_graph_beta_device_management_managed_device" "ios_devices_with_vpp" {
  filter = "(operatingSystem eq 'iOS' or operatingSystem eq 'iPadOS')"
}

# Example 1: Revoke Apple VPP licenses from specific iOS managed devices
# Use this when reclaiming VPP licenses from specific devices that are being retired
action "microsoft365_graph_beta_device_management_managed_device_revoke_apple_vpp_licenses" "revoke_specific_devices" {
  managed_device_ids = [
    "12345678-1234-1234-1234-123456789abc",
    "87654321-4321-4321-4321-ba9876543210"
  ]
}

# Example 2: Revoke Apple VPP licenses from co-managed iOS devices
# Use this for devices managed by both Intune and Configuration Manager
action "microsoft365_graph_beta_device_management_managed_device_revoke_apple_vpp_licenses" "revoke_comanaged_devices" {
  comanaged_device_ids = [
    "11111111-1111-1111-1111-111111111111"
  ]
}

# Example 3: Revoke Apple VPP licenses from both managed and co-managed devices
# Use this for mixed device management scenarios
action "microsoft365_graph_beta_device_management_managed_device_revoke_apple_vpp_licenses" "revoke_mixed_devices" {
  managed_device_ids = [
    "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
    "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"
  ]

  comanaged_device_ids = [
    "cccccccc-cccc-cccc-cccc-cccccccccccc"
  ]
}

# Example 4: Revoke Apple VPP licenses from all iOS devices using data source
# Use this when performing bulk license reclamation for all iOS/iPadOS devices
action "microsoft365_graph_beta_device_management_managed_device_revoke_apple_vpp_licenses" "revoke_all_ios_devices" {
  managed_device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.ios_devices_with_vpp.managed_devices : device.id]

  timeouts = {
    invoke = "20m"
  }
}

# Example 5: Revoke licenses from specific devices by device name using filter
# Use this when you need to target devices by name for license recovery
data "microsoft365_graph_beta_device_management_managed_device" "specific_ios_devices" {
  filter = "startswith(deviceName, 'iPad-Retail-') and (operatingSystem eq 'iOS' or operatingSystem eq 'iPadOS')"
}

action "microsoft365_graph_beta_device_management_managed_device_revoke_apple_vpp_licenses" "revoke_by_name_pattern" {
  managed_device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.specific_ios_devices.managed_devices : device.id]
}

# Example 6: Revoke licenses from non-compliant iOS devices
# Use this to reclaim licenses from devices that are no longer compliant
data "microsoft365_graph_beta_device_management_managed_device" "non_compliant_ios" {
  filter = "complianceState eq 'noncompliant' and (operatingSystem eq 'iOS' or operatingSystem eq 'iPadOS')"
}

action "microsoft365_graph_beta_device_management_managed_device_revoke_apple_vpp_licenses" "revoke_noncompliant" {
  managed_device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.non_compliant_ios.managed_devices : device.id]

  timeouts = {
    invoke = "15m"
  }
}

# Example 7: Revoke licenses from inactive iOS devices (not synced in 30+ days)
# Use this to optimize license allocation by reclaiming from inactive devices
data "microsoft365_graph_beta_device_management_managed_device" "inactive_ios_devices" {
  filter = "lastSyncDateTime lt 2024-01-01T00:00:00Z and (operatingSystem eq 'iOS' or operatingSystem eq 'iPadOS')"
}

action "microsoft365_graph_beta_device_management_managed_device_revoke_apple_vpp_licenses" "revoke_inactive" {
  managed_device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.inactive_ios_devices.managed_devices : device.id]
}

# Example 8: Revoke licenses from supervised iOS devices
# Use this when you need to reclaim licenses specifically from supervised devices
data "microsoft365_graph_beta_device_management_managed_device" "supervised_ios" {
  filter = "isSupervised eq true and (operatingSystem eq 'iOS' or operatingSystem eq 'iPadOS')"
}

action "microsoft365_graph_beta_device_management_managed_device_revoke_apple_vpp_licenses" "revoke_supervised" {
  managed_device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.supervised_ios.managed_devices : device.id]
}

# Example 9: Revoke licenses with custom timeout
# Use this for large-scale license revocation operations
action "microsoft365_graph_beta_device_management_managed_device_revoke_apple_vpp_licenses" "revoke_with_timeout" {
  managed_device_ids = [
    "device-id-1",
    "device-id-2",
    "device-id-3"
  ]

  timeouts = {
    invoke = "30m"
  }
}
```

<!-- action schema generated by tfplugindocs -->
## Schema

### Optional

- `comanaged_device_ids` (List of String) List of co-managed device IDs to revoke Apple VPP licenses from. These are iOS/iPadOS devices managed by both Intune and Configuration Manager (SCCM). Each ID must be a valid GUID format. Example: `["12345678-1234-1234-1234-123456789abc"]`

**Note:** Co-management is rare for iOS/iPadOS devices but supported by this action. At least one of `managed_device_ids` or `comanaged_device_ids` must be provided.
- `managed_device_ids` (List of String) List of managed device IDs to revoke Apple VPP licenses from. These are iOS/iPadOS devices fully managed by Intune only. Each ID must be a valid GUID format. All VPP licenses will be revoked from these devices. Example: `["12345678-1234-1234-1234-123456789abc", "87654321-4321-4321-4321-ba9876543210"]`

**Note:** At least one of `managed_device_ids` or `comanaged_device_ids` must be provided. You can provide both to revoke licenses from different types of devices in one action.
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).

