---
page_title: "microsoft365_graph_beta_device_management_managed_device_reboot_now Action - terraform-provider-microsoft365"
subcategory: "Device Management"

description: |-
  Remotely reboots managed devices using the /deviceManagement/managedDevices/{managedDeviceId}/rebootNow endpoint. This action immediately restarts devices, which is essential for applying updates, troubleshooting system issues, or ensuring configuration changes take effect. The reboot command is sent to devices immediately if online, or queued for execution when the device next checks in with Intune.
  Important Notes:
  Device reboots immediately upon receiving command (if online)Any unsaved work on the device will be lostUsers receive minimal or no warning before rebootReboot is forceful and does not wait for user interactionCommand is queued if device is offlineUse with caution during business hours
  Use Cases:
  Applying Windows updates that require restartInstalling software that requires system rebootTroubleshooting devices with performance issuesForcing application of configuration profilesClearing temporary system issuesMaintenance windows for device refreshResolving frozen or unresponsive remote devicesCompleting BitLocker encryption setup
  Platform Support:
  Windows: Fully supported (Windows 10/11, including Home edition)macOS: Supported (requires user-approved MDM or supervised)iOS/iPadOS: Limited support (supervised devices only)Android: Not supported for reboot action
  Best Practices:
  Schedule reboots during maintenance windows or off-hoursNotify users in advance when possibleUse for non-interactive devices (kiosks, shared devices)Consider user impact before rebooting during business hoursTest with small device groups before bulk operationsDocument reason for reboot in change management systemCombine with compliance policies for automated maintenance
  User Impact:
  Users may lose unsaved workActive sessions are terminatedVideo calls and presentations are interruptedFile transfers may be interruptedUsers may not receive advance warningDevice is unavailable for 2-5 minutes during restart
  Reference: Microsoft Graph API - Reboot Now https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-rebootnow?view=graph-rest-beta
---

# microsoft365_graph_beta_device_management_managed_device_reboot_now (Action)

Remotely reboots managed devices using the `/deviceManagement/managedDevices/{managedDeviceId}/rebootNow` endpoint. This action immediately restarts devices, which is essential for applying updates, troubleshooting system issues, or ensuring configuration changes take effect. The reboot command is sent to devices immediately if online, or queued for execution when the device next checks in with Intune.

**Important Notes:**
- Device reboots immediately upon receiving command (if online)
- Any unsaved work on the device will be lost
- Users receive minimal or no warning before reboot
- Reboot is forceful and does not wait for user interaction
- Command is queued if device is offline
- Use with caution during business hours

**Use Cases:**
- Applying Windows updates that require restart
- Installing software that requires system reboot
- Troubleshooting devices with performance issues
- Forcing application of configuration profiles
- Clearing temporary system issues
- Maintenance windows for device refresh
- Resolving frozen or unresponsive remote devices
- Completing BitLocker encryption setup

**Platform Support:**
- **Windows**: Fully supported (Windows 10/11, including Home edition)
- **macOS**: Supported (requires user-approved MDM or supervised)
- **iOS/iPadOS**: Limited support (supervised devices only)
- **Android**: Not supported for reboot action

**Best Practices:**
- Schedule reboots during maintenance windows or off-hours
- Notify users in advance when possible
- Use for non-interactive devices (kiosks, shared devices)
- Consider user impact before rebooting during business hours
- Test with small device groups before bulk operations
- Document reason for reboot in change management system
- Combine with compliance policies for automated maintenance

**User Impact:**
- Users may lose unsaved work
- Active sessions are terminated
- Video calls and presentations are interrupted
- File transfers may be interrupted
- Users may not receive advance warning
- Device is unavailable for 2-5 minutes during restart

**Reference:** [Microsoft Graph API - Reboot Now](https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-rebootnow?view=graph-rest-beta)

## Microsoft Documentation

### Graph API References
- [rebootNow action](https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-rebootnow?view=graph-rest-beta)
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

| Platform | Support | Requirements |
|----------|---------|--------------|
| **Windows** | ✅ Full Support | All versions including Home |
| **macOS** | ✅ Supported | User-approved MDM or supervised |
| **iOS** | ⚠️ Limited | Supervised devices only |
| **iPadOS** | ⚠️ Limited | Supervised devices only |
| **Android** | ❌ Not Supported | Remote reboot not available |

### Reboot vs Shutdown

| Action | Result | Use Case |
|--------|--------|----------|
| **Reboot** | Device restarts automatically | Updates, troubleshooting, config changes |
| **Shutdown** | Device powers off completely | Long-term offline, energy conservation |

### User Impact

- **Unsaved Work**: Lost immediately
- **Active Sessions**: Terminated (calls, presentations, etc.)
- **Downtime**: Typically 2-5 minutes
- **Warning**: Minimal or no advance notification
- **File Transfers**: Interrupted
- **Network Sessions**: Disconnected

### Common Use Cases

- Applying Windows updates requiring restart
- Installing software that requires system reboot
- Troubleshooting performance issues
- Forcing application of configuration profiles
- Clearing temporary system issues
- Maintenance windows for device refresh
- Completing BitLocker encryption setup
- Resolving frozen or unresponsive devices

### Best Practices

- Schedule during maintenance windows or off-hours
- Notify users in advance when possible
- Use for unattended devices (kiosks, shared devices, labs)
- Test with small device groups before bulk operations
- Document reason for reboot in change management
- Consider user time zones for global deployments
- Avoid during peak business hours
- Monitor device recovery after reboot

## Example Usage

```terraform
# Example 1: Reboot a single device
action "microsoft365_graph_beta_device_management_managed_device_reboot_now" "reboot_single" {

  device_ids = [
    "12345678-1234-1234-1234-123456789abc"
  ]

  timeouts = {
    invoke = "5m"
  }
}

# Example 2: Reboot multiple devices
action "microsoft365_graph_beta_device_management_managed_device_reboot_now" "reboot_batch" {

  device_ids = [
    "12345678-1234-1234-1234-123456789abc",
    "87654321-4321-4321-4321-ba9876543210",
    "abcdef12-3456-7890-abcd-ef1234567890"
  ]

  timeouts = {
    invoke = "10m"
  }
}

# Example 3: Reboot Windows devices with non-compliant state
data "microsoft365_graph_beta_device_management_managed_device" "windows_noncompliant" {
  filter_type  = "odata"
  odata_filter = "(operatingSystem eq 'Windows') and (complianceState eq 'noncompliant')"
}

action "microsoft365_graph_beta_device_management_managed_device_reboot_now" "reboot_windows_noncompliant" {

  device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.windows_noncompliant.items : device.id]

  timeouts = {
    invoke = "15m"
  }
}

# Example 4: Reboot kiosk devices (minimal user impact)
data "microsoft365_graph_beta_device_management_managed_device" "kiosk_devices" {
  filter_type  = "device_name"
  filter_value = "KIOSK-"
}

action "microsoft365_graph_beta_device_management_managed_device_reboot_now" "reboot_kiosks" {

  device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.kiosk_devices.items : device.id]

  timeouts = {
    invoke = "10m"
  }
}

# Example 5: Reboot corporate Windows devices
data "microsoft365_graph_beta_device_management_managed_device" "corporate_windows" {
  filter_type  = "odata"
  odata_filter = "(operatingSystem eq 'Windows') and (managedDeviceOwnerType eq 'company')"
}

action "microsoft365_graph_beta_device_management_managed_device_reboot_now" "reboot_corporate_windows" {

  device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.corporate_windows.items : device.id]

  timeouts = {
    invoke = "20m"
  }
}

# Example 6: Scheduled maintenance reboot for lab devices
data "microsoft365_graph_beta_device_management_managed_device" "lab_devices" {
  filter_type  = "device_name"
  filter_value = "LAB-"
}

action "microsoft365_graph_beta_device_management_managed_device_reboot_now" "reboot_lab_maintenance" {

  device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.lab_devices.items : device.id]

  timeouts = {
    invoke = "10m"
  }
}

# Output examples
output "rebooted_device_count" {
  value       = length(action.reboot_batch.device_ids)
  description = "Number of devices that received reboot command"
}

output "windows_noncompliant_reboot_count" {
  value       = length(action.reboot_windows_noncompliant.device_ids)
  description = "Number of non-compliant Windows devices rebooted"
}
```

<!-- action schema generated by tfplugindocs -->
## Schema

### Required

- `device_ids` (List of String) List of managed device IDs to reboot. Each ID must be a valid GUID format. Multiple devices can be rebooted in a single action. Example: `["12345678-1234-1234-1234-123456789abc", "87654321-4321-4321-4321-ba9876543210"]`

**Important:** Devices will reboot immediately when they receive this command. Any unsaved work will be lost. Use with caution, especially during business hours.

### Optional

- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).

