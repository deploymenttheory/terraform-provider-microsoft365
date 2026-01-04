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
| v0.40.0-alpha | Experimental | Example fixes and refactored sync progress logic |


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
# Example 1: Revoke Apple VPP licenses from a single device - Minimal
action "microsoft365_graph_beta_device_management_managed_device_revoke_apple_vpp_licenses" "revoke_single" {
  config {
    managed_device_ids = [
      "12345678-1234-1234-1234-123456789abc"
    ]
  }
}

# Example 2: Revoke Apple VPP licenses from multiple devices
action "microsoft365_graph_beta_device_management_managed_device_revoke_apple_vpp_licenses" "revoke_multiple" {
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

# Example 3: Revoke Apple VPP licenses with validation - Maximal
action "microsoft365_graph_beta_device_management_managed_device_revoke_apple_vpp_licenses" "revoke_with_validation" {
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

# Example 4: Revoke VPP licenses from departing user's devices
data "microsoft365_graph_beta_device_management_managed_device" "departing_user_ios" {
  filter_type  = "odata"
  odata_filter = "(userPrincipalName eq 'departing.user@example.com') and ((operatingSystem eq 'iOS') or (operatingSystem eq 'iPadOS'))"
}

action "microsoft365_graph_beta_device_management_managed_device_revoke_apple_vpp_licenses" "revoke_departing_user" {
  config {
    managed_device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.departing_user_ios.items : device.id]

    validate_device_exists = true

    timeouts = {
      invoke = "10m"
    }
  }
}

# Example 5: Revoke VPP licenses from all iOS/iPadOS devices
data "microsoft365_graph_beta_device_management_managed_device" "all_apple_devices" {
  filter_type  = "odata"
  odata_filter = "(operatingSystem eq 'iOS') or (operatingSystem eq 'iPadOS')"
}

action "microsoft365_graph_beta_device_management_managed_device_revoke_apple_vpp_licenses" "revoke_all_ios" {
  config {
    managed_device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.all_apple_devices.items : device.id]

    ignore_partial_failures = true

    timeouts = {
      invoke = "30m"
    }
  }
}

# Example 6: Revoke VPP licenses for co-managed devices
action "microsoft365_graph_beta_device_management_managed_device_revoke_apple_vpp_licenses" "revoke_comanaged" {
  config {
    comanaged_device_ids = [
      "11111111-1111-1111-1111-111111111111",
      "22222222-2222-2222-2222-222222222222"
    ]

    timeouts = {
      invoke = "10m"
    }
  }
}

# Output examples
output "revoked_vpp_licenses_count" {
  value       = length(action.microsoft365_graph_beta_device_management_managed_device_revoke_apple_vpp_licenses.revoke_multiple.config.managed_device_ids)
  description = "Number of devices that had VPP licenses revoked"
}
```

<!-- action schema generated by tfplugindocs -->
## Schema

### Optional

- `comanaged_device_ids` (List of String) List of co-managed device IDs to revoke Apple VPP licenses from. These are iOS/iPadOS devices managed by both Intune and Configuration Manager (SCCM). Each ID must be a valid GUID format. Example: `["12345678-1234-1234-1234-123456789abc"]`

**Note:** Co-management is rare for iOS/iPadOS devices but supported by this action. At least one of `managed_device_ids` or `comanaged_device_ids` must be provided.
- `ignore_partial_failures` (Boolean) If set to `true`, the action will succeed even if some operations fail. Failed operations will be reported as warnings instead of errors. Default: `false` (action fails if any operation fails).
- `managed_device_ids` (List of String) List of managed device IDs to revoke Apple VPP licenses from. These are iOS/iPadOS devices fully managed by Intune only. Each ID must be a valid GUID format. All VPP licenses will be revoked from these devices. Example: `["12345678-1234-1234-1234-123456789abc", "87654321-4321-4321-4321-ba9876543210"]`

**Note:** At least one of `managed_device_ids` or `comanaged_device_ids` must be provided. You can provide both to revoke licenses from different types of devices in one action.
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))
- `validate_device_exists` (Boolean) Whether to validate that devices exist and are iOS/iPadOS devices before attempting to revoke licenses. Disabling this can speed up planning but may result in runtime errors for non-existent or non-Apple devices. Default: `true`.

<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `invoke` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).

