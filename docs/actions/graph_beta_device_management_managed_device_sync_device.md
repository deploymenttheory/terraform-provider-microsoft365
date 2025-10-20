---
page_title: "microsoft365_graph_beta_device_management_managed_device_sync_device Action - terraform-provider-microsoft365"
subcategory: "Device Management"

description: |-
  Forces managed and co-managed devices to immediately check in with Intune using the /deviceManagement/managedDevices/{managedDeviceId}/syncDevice and /deviceManagement/comanagedDevices/{managedDeviceId}/syncDevice endpoints. This action triggers an immediate synchronization, causing devices to apply the latest policies, configurations, and updates from Intune without waiting for the standard check-in interval.
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

Forces managed and co-managed devices to immediately check in with Intune using the `/deviceManagement/managedDevices/{managedDeviceId}/syncDevice` and `/deviceManagement/comanagedDevices/{managedDeviceId}/syncDevice` endpoints. This action triggers an immediate synchronization, causing devices to apply the latest policies, configurations, and updates from Intune without waiting for the standard check-in interval.

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
- [Windows Remote Actions](https://learn.microsoft.com/en-us/intune/intune-service/remote-actions/?tabs=windows)
- [iOS/iPadOS Remote Actions](https://learn.microsoft.com/en-us/intune/intune-service/remote-actions/?tabs=ios-ipados)
- [macOS Remote Actions](https://learn.microsoft.com/en-us/intune/intune-service/remote-actions/?tabs=macos)
- [Android Remote Actions](https://learn.microsoft.com/en-us/intune/intune-service/remote-actions/?tabs=android)
- [ChromeOS Remote Actions](https://learn.microsoft.com/en-us/intune/intune-service/remote-actions/?tabs=chromeos)

## API Permissions

The following API permissions are required in order to use this action.

### Microsoft Graph

- **Application**: `DeviceManagementManagedDevices.PrivilegedOperations.All`
- **Delegated**: `DeviceManagementManagedDevices.PrivilegedOperations.All`

## Version History

| Version | Status | Notes |
|---------|--------|-------|
| v0.33.0-alpha | Experimental | Initial release |

## Notes

### Platform Compatibility

| Platform | Managed Devices | Co-Managed Devices | Notes |
|----------|----------------|-------------------|-------|
| **Windows** | ✅ Full Support | ✅ Full Support | Fastest sync response |
| **macOS** | ✅ Full Support | ❌ Not Applicable | Supervised devices recommended |
| **iOS** | ✅ Full Support | ❌ Not Applicable | Supervised devices recommended |
| **iPadOS** | ✅ Full Support | ❌ Not Applicable | Supervised devices recommended |
| **Android** | ✅ Full Support | ❌ Not Applicable | Includes Android Enterprise |
| **ChromeOS** | ⚠️ Limited Support | ❌ Not Applicable | Sync may be delayed |

### Managed vs Co-Managed Devices

| Device Type | Management | Sync Endpoint | Typical Use Case |
|------------|------------|---------------|------------------|
| **Managed** | Intune only | `/managedDevices/{id}/syncDevice` | BYOD, mobile devices, cloud-only management |
| **Co-Managed** | Intune + ConfigMgr | `/comanagedDevices/{id}/syncDevice` | Enterprise Windows devices, hybrid management |

### Co-Management Context

**What is Co-Management?**
- Windows devices managed by both Intune and Configuration Manager (SCCM)
- Workloads split between the two management platforms
- Provides gradual transition from ConfigMgr to Intune
- Typically Windows 10/11 enterprise devices

**How Sync Works with Co-Management**
- Sync affects only Intune-managed workloads
- ConfigMgr workloads use ConfigMgr sync mechanisms
- Device must be enrolled in both systems
- Sync status visible in both consoles

**Co-Management Workloads**
- Compliance policies → Intune or ConfigMgr
- Device configuration → Intune or ConfigMgr
- Endpoint Protection → Intune or ConfigMgr
- Resource access → Intune or ConfigMgr
- Windows Update policies → Intune or ConfigMgr
- Office Click-to-Run apps → Intune or ConfigMgr

### What Happens During Sync

**Immediate Actions**
1. Device receives sync command (if online)
2. Initiates check-in with Intune
3. Downloads latest policy assignments
4. Evaluates compliance state
5. Downloads assigned applications
6. Reports updated inventory
7. Processes queued remote actions
8. Updates device status in console

**Policy Application**
- New policies downloaded and applied
- Changed policies updated immediately
- Removed policies uninstalled
- Configuration profiles refreshed
- Compliance evaluated
- Conditional access re-checked

**Application Management**
- Required apps downloaded/installed
- Available apps catalog updated
- App assignments processed
- App configuration applied
- App protection policies refreshed

**Device Inventory**
- Hardware information updated
- Software inventory refreshed
- Installed apps reported
- Disk space usage updated
- Battery status (mobile devices)
- Network information

### Check-In Intervals

| Scenario | Check-In Frequency | Notes |
|----------|-------------------|-------|
| **Normal** | Every 8 hours | Default for most platforms |
| **iOS/iPadOS** | Every 6-8 hours | May vary by iOS version |
| **Android** | Every 8 hours | Enterprise devices |
| **Windows** | Every 8 hours | Can be customized via policy |
| **macOS** | Every 8 hours | Supervised devices |
| **Sync Action** | 1-5 minutes | Forces immediate check-in |

### Sync Behavior

**Online Devices**
- Receive command immediately
- Begin sync within 1-5 minutes
- Complete sync within 5-15 minutes
- Status updated in real-time

**Offline Devices**
- Sync command queued
- Executed when device comes online
- Queued for up to 7 days
- Auto-removed if expired

**Multiple Syncs**
- Avoid syncing same device repeatedly
- Wait 15+ minutes between syncs
- Multiple syncs may queue
- Can delay processing

**Network Requirements**
- Internet connectivity required
- Access to Intune endpoints
- Firewall rules must permit traffic
- VPN may impact sync speed

### Use Cases

| Scenario | Description | Expected Outcome |
|----------|-------------|------------------|
| **Policy Deployment** | Force new policies to apply immediately | Policies active within minutes |
| **App Installation** | Push required apps urgently | Apps install on next check-in |
| **Compliance Evaluation** | Re-check device compliance | Updated compliance state |
| **Troubleshooting** | Verify policy deployment | See real-time status |
| **Emergency Updates** | Critical security updates | Immediate application |
| **New Enrollments** | Ensure fresh devices get all policies | Complete configuration |
| **User Account Changes** | Apply new user-specific policies | Updated user context |
| **Inventory Update** | Refresh device information | Current device state |

### Best Practices

**When to Sync**
- After creating new policies
- After modifying existing policies
- When troubleshooting deployment issues
- For critical security updates
- After device enrollment
- When verifying policy application

**When NOT to Sync**
- Repeatedly within short periods (< 15 min)
- For large numbers of offline devices
- During normal operations (let auto-sync work)
- For policies that can wait for normal sync
- On metered/slow networks unnecessarily

**Planning Considerations**
- Sync during business hours for online devices
- Schedule bulk syncs during off-hours
- Consider network bandwidth impact
- Allow time for completion (1-5 min per device)
- Monitor sync status in Intune admin center
- Document reason for manual sync

**Performance Optimization**
- Batch related devices together
- Use datasources to target specific groups
- Avoid duplicate device IDs
- Set appropriate timeouts
- Monitor for failed syncs
- Re-sync failures individually

### Troubleshooting

| Issue | Cause | Solution |
|-------|-------|----------|
| Sync fails | Device offline | Wait for device to come online |
| Timeout | Too many devices | Reduce batch size or increase timeout |
| Device not found | Invalid device ID | Verify device exists in Intune |
| Permission denied | Insufficient permissions | Check API permissions |
| Sync queued | Device offline | Command will execute when online |
| No effect | Wrong device type | Check managed vs co-managed |
| Slow sync | Network issues | Check network connectivity |
| Repeated failures | Device enrollment issue | Re-enroll device |

### Monitoring Sync Status

**Intune Admin Center**
1. Navigate to Devices → All devices
2. Select specific device
3. View "Last check-in" timestamp
4. Check "Device sync status"
5. Review "Pending actions"

**Terraform**
- Action reports success/failure
- Check diagnostics for errors
- Review progress messages
- Monitor timeout settings

**Time to Complete**
- Initial command: < 1 minute
- Device check-in: 1-5 minutes
- Policy application: 5-15 minutes
- Complete refresh: 15-30 minutes

### Limitations

**Technical Limits**
- Devices must be online
- Requires network connectivity
- Subject to API rate limits
- Maximum concurrent syncs per tenant
- Queued commands expire after 7 days

**Platform Limits**
- ChromeOS sync may be delayed
- Some policies require reboot
- App installation depends on size
- Compliance evaluation takes time
- Inventory update is not instant

**Co-Management Limits**
- Only affects Intune workloads
- ConfigMgr workloads not triggered
- Workload assignment matters
- Both systems must be healthy

## Example Usage

```terraform
# ============================================================================
# Example 1: Sync managed devices only
# ============================================================================
# Use case: Force immediate policy application on fully Intune-managed devices
action "microsoft365_graph_beta_device_management_managed_device_sync_device" "sync_managed_only" {

  managed_device_ids = [
    "12345678-1234-1234-1234-123456789abc",
    "87654321-4321-4321-4321-ba9876543210",
    "abcdef12-3456-7890-abcd-ef1234567890"
  ]

  timeouts = {
    invoke = "10m"
  }
}

# ============================================================================
# Example 2: Sync co-managed devices only
# ============================================================================
# Use case: Force sync on devices managed by both Intune and Configuration Manager
action "microsoft365_graph_beta_device_management_managed_device_sync_device" "sync_comanaged_only" {

  comanaged_device_ids = [
    "12345678-1234-1234-1234-123456789abc",
    "87654321-4321-4321-4321-ba9876543210"
  ]

  timeouts = {
    invoke = "10m"
  }
}

# ============================================================================
# Example 3: Sync both managed and co-managed devices
# ============================================================================
# Use case: Mixed environment with both fully Intune-managed and co-managed devices
action "microsoft365_graph_beta_device_management_managed_device_sync_device" "sync_mixed_devices" {

  managed_device_ids = [
    "12345678-1234-1234-1234-123456789abc",
    "87654321-4321-4321-4321-ba9876543210"
  ]

  comanaged_device_ids = [
    "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
    "11111111-2222-3333-4444-555555555555"
  ]

  timeouts = {
    invoke = "15m"
  }
}

# ============================================================================
# Example 4: Sync all Windows devices using datasource
# ============================================================================
# Use case: Force immediate sync on all Windows devices after policy change
data "microsoft365_graph_beta_device_management_managed_device" "windows_devices" {
  filter_type  = "odata"
  odata_filter = "operatingSystem eq 'Windows'"
}

action "microsoft365_graph_beta_device_management_managed_device_sync_device" "sync_all_windows" {

  managed_device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.windows_devices.items : device.id]

  timeouts = {
    invoke = "30m"
  }
}

# ============================================================================
# Example 5: Sync devices by name pattern
# ============================================================================
# Use case: Sync specific group of devices based on naming convention
data "microsoft365_graph_beta_device_management_managed_device" "lab_devices" {
  filter_type  = "device_name"
  filter_value = "LAB-"
}

action "microsoft365_graph_beta_device_management_managed_device_sync_device" "sync_lab_devices" {

  managed_device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.lab_devices.items : device.id]

  timeouts = {
    invoke = "20m"
  }
}

# ============================================================================
# Example 6: Sync non-compliant devices
# ============================================================================
# Use case: Force compliance re-evaluation on non-compliant devices
data "microsoft365_graph_beta_device_management_managed_device" "non_compliant" {
  filter_type  = "odata"
  odata_filter = "complianceState eq 'noncompliant'"
}

action "microsoft365_graph_beta_device_management_managed_device_sync_device" "sync_non_compliant" {

  managed_device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.non_compliant.items : device.id]

  timeouts = {
    invoke = "30m"
  }
}

# ============================================================================
# Example 7: Sync iOS/iPadOS devices
# ============================================================================
# Use case: Force app installation on Apple mobile devices
data "microsoft365_graph_beta_device_management_managed_device" "ios_devices" {
  filter_type  = "odata"
  odata_filter = "(operatingSystem eq 'iOS') or (operatingSystem eq 'iPadOS')"
}

action "microsoft365_graph_beta_device_management_managed_device_sync_device" "sync_ios_devices" {

  managed_device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.ios_devices.items : device.id]

  timeouts = {
    invoke = "20m"
  }
}

# ============================================================================
# Example 8: Sync devices by user
# ============================================================================
# Use case: Sync all devices assigned to specific user after account changes
data "microsoft365_graph_beta_device_management_managed_device" "user_devices" {
  filter_type  = "odata"
  odata_filter = "userPrincipalName eq 'john.doe@company.com'"
}

action "microsoft365_graph_beta_device_management_managed_device_sync_device" "sync_user_devices" {

  managed_device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.user_devices.items : device.id]

  timeouts = {
    invoke = "10m"
  }
}

# ============================================================================
# Example 9: Sync recently enrolled devices
# ============================================================================
# Use case: Ensure new devices get all policies immediately after enrollment
locals {
  # Get current timestamp
  three_days_ago = timeadd(timestamp(), "-72h")
}

data "microsoft365_graph_beta_device_management_managed_device" "recent_enrollments" {
  filter_type  = "odata"
  odata_filter = "enrolledDateTime ge ${formatdate("YYYY-MM-DD'T'hh:mm:ss'Z'", local.three_days_ago)}"
}

action "microsoft365_graph_beta_device_management_managed_device_sync_device" "sync_recent_enrollments" {

  managed_device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.recent_enrollments.items : device.id]

  timeouts = {
    invoke = "20m"
  }
}

# ============================================================================
# Example 10: Emergency policy deployment
# ============================================================================
# Use case: Critical security update - sync all devices immediately
data "microsoft365_graph_beta_device_management_managed_device" "all_managed" {
  filter_type = "all"
}

data "microsoft365_graph_beta_device_management_managed_device" "all_comanaged" {
  filter_type  = "odata"
  odata_filter = "managementAgent eq 'configurationManagerClientMdm'"
}

action "microsoft365_graph_beta_device_management_managed_device_sync_device" "emergency_sync_all" {

  managed_device_ids   = [for device in data.microsoft365_graph_beta_device_management_managed_device.all_managed.items : device.id]
  comanaged_device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.all_comanaged.items : device.id]

  timeouts = {
    invoke = "60m"
  }
}

# ============================================================================
# Example 11: Sync specific macOS devices
# ============================================================================
# Use case: Force FileVault policy update on company MacBooks
data "microsoft365_graph_beta_device_management_managed_device" "macos_devices" {
  filter_type  = "odata"
  odata_filter = "operatingSystem eq 'macOS'"
}

action "microsoft365_graph_beta_device_management_managed_device_sync_device" "sync_macos" {

  managed_device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.macos_devices.items : device.id]

  timeouts = {
    invoke = "15m"
  }
}

# ============================================================================
# Example 12: Conditional sync based on device state
# ============================================================================
# Use case: Sync only online devices to avoid queuing for offline devices
data "microsoft365_graph_beta_device_management_managed_device" "all_devices" {
  filter_type = "all"
}

locals {
  # Get current timestamp
  one_day_ago = timeadd(timestamp(), "-24h")
  
  # Filter to devices that checked in within last 24 hours (likely online)
  online_devices = [
    for device in data.microsoft365_graph_beta_device_management_managed_device.all_devices.items :
    device.id
    if timecmp(device.last_sync_date_time, local.one_day_ago) > 0
  ]
}

action "microsoft365_graph_beta_device_management_managed_device_sync_device" "sync_online_only" {

  managed_device_ids = local.online_devices

  timeouts = {
    invoke = "30m"
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
- `managed_device_ids` (List of String) List of managed device IDs to sync. These are devices fully managed by Intune only. Each ID must be a valid GUID format. Multiple devices can be synced in a single action. Example: `["12345678-1234-1234-1234-123456789abc", "87654321-4321-4321-4321-ba9876543210"]`

**Note:** At least one of `managed_device_ids` or `comanaged_device_ids` must be provided. You can provide both to sync different types of devices in one action.
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).


