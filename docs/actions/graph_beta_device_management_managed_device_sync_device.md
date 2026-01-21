---
page_title: "microsoft365_graph_beta_device_management_managed_device_sync_device Action - terraform-provider-microsoft365"
subcategory: "Device Management"

description: |-
  Syncs managed and co-managed devices in Microsoft Intune using the /deviceManagement/managedDevices/{managedDeviceId}/syncDevice and /deviceManagement/comanagedDevices/{managedDeviceId}/syncDevice endpoints. This action is used to force devices to immediately check in with Intune, triggering an immediate synchronization that causes devices to apply the latest policies, configurations, and updates from Intune without waiting for the standard check-in interval.
  What This Action Does:
  Forces immediate check-in with IntuneApplies latest policies and configurationsDownloads pending applicationsReports updated device inventoryEnforces compliance evaluationProcesses queued remote actionsUpdates device status in console
  Managed vs Co-Managed Devices:
  Managed Devices: Fully managed by Intune onlyCo-Managed Devices: Managed by both Intune and Configuration Manager (SCCM)This action supports both types independently or together
  Platform Support:
  Windows: Full support (managed and co-managed)macOS: Full support (managed only)iOS/iPadOS: Full support (managed only)Android: Full support (managed only)ChromeOS: Limited support
  Common Use Cases:
  Apply new policies immediatelyForce app installation/updatesTrigger compliance re-evaluationUpdate device inventory quicklyVerify policy deploymentTroubleshoot deployment issuesEmergency configuration changes
  Check-In Behavior:
  Normal interval: Every 8 hours (varies by platform)This action: Immediate (within 1-5 minutes)Device must be online and powered onNetwork connectivity requiredResults visible in Intune admin center
  Important Considerations:
  Device must be online to receive commandCommand queued if device is offlineSync completes when device comes onlineMultiple syncs in short period may delay each otherNo user disruption (background operation)
  Reference: Microsoft Graph API - Sync Device https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-syncdevice?view=graph-rest-beta
---

# microsoft365_graph_beta_device_management_managed_device_sync_device (Action)

Syncs managed and co-managed devices in Microsoft Intune using the `/deviceManagement/managedDevices/{managedDeviceId}/syncDevice` and `/deviceManagement/comanagedDevices/{managedDeviceId}/syncDevice` endpoints. This action is used to force devices to immediately check in with Intune, triggering an immediate synchronization that causes devices to apply the latest policies, configurations, and updates from Intune without waiting for the standard check-in interval.

**What This Action Does:**
- Forces immediate check-in with Intune
- Applies latest policies and configurations
- Downloads pending applications
- Reports updated device inventory
- Enforces compliance evaluation
- Processes queued remote actions
- Updates device status in console

**Managed vs Co-Managed Devices:**
- **Managed Devices**: Fully managed by Intune only
- **Co-Managed Devices**: Managed by both Intune and Configuration Manager (SCCM)
- This action supports both types independently or together

**Platform Support:**
- **Windows**: Full support (managed and co-managed)
- **macOS**: Full support (managed only)
- **iOS/iPadOS**: Full support (managed only)
- **Android**: Full support (managed only)
- **ChromeOS**: Limited support

**Common Use Cases:**
- Apply new policies immediately
- Force app installation/updates
- Trigger compliance re-evaluation
- Update device inventory quickly
- Verify policy deployment
- Troubleshoot deployment issues
- Emergency configuration changes

**Check-In Behavior:**
- Normal interval: Every 8 hours (varies by platform)
- This action: Immediate (within 1-5 minutes)
- Device must be online and powered on
- Network connectivity required
- Results visible in Intune admin center

**Important Considerations:**
- Device must be online to receive command
- Command queued if device is offline
- Sync completes when device comes online
- Multiple syncs in short period may delay each other
- No user disruption (background operation)

**Reference:** [Microsoft Graph API - Sync Device](https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-syncdevice?view=graph-rest-beta)

## Microsoft Documentation

### Graph API References
- [syncDevice action](https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-syncdevice?view=graph-rest-beta)
- [managedDevice resource type](https://learn.microsoft.com/en-us/graph/api/resources/intune-devices-manageddevice?view=graph-rest-beta)

### Intune Remote Actions Guides
- [Device sync - Windows](https://learn.microsoft.com/en-us/intune/intune-service/remote-actions/device-sync?pivots=windows)
- [Device sync - iOS](https://learn.microsoft.com/en-us/intune/intune-service/remote-actions/device-sync?pivots=ios)
- [Device sync - macOS](https://learn.microsoft.com/en-us/intune/intune-service/remote-actions/device-sync?pivots=macos)
- [Device sync - Android](https://learn.microsoft.com/en-us/intune/intune-service/remote-actions/device-sync?pivots=android)

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
# Example 1: Sync a single managed device - Minimal
action "microsoft365_graph_beta_device_management_managed_device_sync_device" "sync_single" {
  config {
    managed_device_ids = [
      "12345678-1234-1234-1234-123456789abc"
    ]
  }
}

# Example 2: Sync multiple managed devices
action "microsoft365_graph_beta_device_management_managed_device_sync_device" "sync_managed_only" {
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

# Example 3: Sync co-managed devices only
action "microsoft365_graph_beta_device_management_managed_device_sync_device" "sync_comanaged_only" {
  config {
    comanaged_device_ids = [
      "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
      "11111111-2222-3333-4444-555555555555"
    ]

    timeouts = {
      invoke = "10m"
    }
  }
}

# Example 4: Sync both managed and co-managed devices - Maximal
action "microsoft365_graph_beta_device_management_managed_device_sync_device" "sync_mixed_devices" {
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
      invoke = "15m"
    }
  }
}

# Example 5: Sync all Windows devices using datasource
data "microsoft365_graph_beta_device_management_managed_device" "windows_devices" {
  filter_type  = "odata"
  odata_filter = "operatingSystem eq 'Windows'"
}

action "microsoft365_graph_beta_device_management_managed_device_sync_device" "sync_all_windows" {
  config {
    managed_device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.windows_devices.items : device.id]

    validate_device_exists  = true
    ignore_partial_failures = true

    timeouts = {
      invoke = "30m"
    }
  }
}

# Example 6: Sync non-compliant devices
data "microsoft365_graph_beta_device_management_managed_device" "non_compliant" {
  filter_type  = "odata"
  odata_filter = "complianceState eq 'noncompliant'"
}

action "microsoft365_graph_beta_device_management_managed_device_sync_device" "sync_non_compliant" {
  config {
    managed_device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.non_compliant.items : device.id]

    validate_device_exists  = true
    ignore_partial_failures = false

    timeouts = {
      invoke = "30m"
    }
  }
}

# Example 7: Sync iOS/iPadOS devices
data "microsoft365_graph_beta_device_management_managed_device" "ios_devices" {
  filter_type  = "odata"
  odata_filter = "(operatingSystem eq 'iOS') or (operatingSystem eq 'iPadOS')"
}

action "microsoft365_graph_beta_device_management_managed_device_sync_device" "sync_ios_devices" {
  config {
    managed_device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.ios_devices.items : device.id]

    timeouts = {
      invoke = "20m"
    }
  }
}

# Example 8: Emergency policy deployment
data "microsoft365_graph_beta_device_management_managed_device" "all_managed" {
  filter_type = "all"
}

data "microsoft365_graph_beta_device_management_managed_device" "all_comanaged" {
  filter_type  = "odata"
  odata_filter = "managementAgent eq 'configurationManagerClientMdm'"
}

action "microsoft365_graph_beta_device_management_managed_device_sync_device" "emergency_sync_all" {
  config {
    managed_device_ids   = [for device in data.microsoft365_graph_beta_device_management_managed_device.all_managed.items : device.id]
    comanaged_device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.all_comanaged.items : device.id]

    validate_device_exists  = true
    ignore_partial_failures = true

    timeouts = {
      invoke = "60m"
    }
  }
}
```

<!-- action schema generated by tfplugindocs -->
## Schema

### Optional

- `comanaged_device_ids` (List of String) List of co-managed device IDs to sync. These are devices managed by both Intune and Configuration Manager (SCCM). Each ID must be a valid GUID format. Example: `["12345678-1234-1234-1234-123456789abc"]`

**Co-Management Context:**
- Devices managed by both Intune and Configuration Manager
- Typically Windows 10/11 enterprise devices
- Workloads split between Intune and ConfigMgr
- Sync affects Intune-managed workloads only

**Note:** At least one of `managed_device_ids` or `comanaged_device_ids` must be provided.
- `ignore_partial_failures` (Boolean) When set to `true`, the action will complete successfully even if some devices fail to sync. When `false` (default), the action will fail if any device sync fails. Use this flag when syncing multiple devices and you want the action to succeed even if some syncs fail.
- `managed_device_ids` (List of String) List of managed device IDs to sync. These are devices fully managed by Intune only. Each ID must be a valid GUID format. Multiple devices can be synced in a single action. Example: `["12345678-1234-1234-1234-123456789abc", "87654321-4321-4321-4321-ba9876543210"]`

**Note:** At least one of `managed_device_ids` or `comanaged_device_ids` must be provided. You can provide both to sync different types of devices in one action.
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))
- `validate_device_exists` (Boolean) When set to `true` (default), the action will validate that all specified devices exist and support sync before attempting to sync them. When `false`, device validation is skipped and the action will attempt to sync devices directly. Disabling validation can improve performance but may result in errors if devices don't exist or are unsupported.

<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `invoke` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).


