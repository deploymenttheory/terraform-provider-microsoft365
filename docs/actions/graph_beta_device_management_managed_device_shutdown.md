---
page_title: "microsoft365_graph_beta_device_management_managed_device_shutdown Action - terraform-provider-microsoft365"
subcategory: "Device Management"

description: |-
  Remotely shuts down managed devices using the /deviceManagement/managedDevices/{managedDeviceId}/shutDown endpoint. This action powers off devices completely, which is useful for energy conservation, maintenance operations, or security scenarios. Unlike reboot, shutdown powers the device off completely and requires manual intervention to power it back on.
  Important Notes:
  Device shuts down completely (powers off)Device requires manual power-on to restartAny unsaved work on the device will be lostUsers receive minimal or no warning before shutdownShutdown is forceful and does not wait for user interactionCommand is queued if device is offlineUse with extreme caution - device will be completely offline
  Use Cases:
  Energy conservation during extended non-use periodsSecurity incident response (isolate compromised device)Hardware maintenance requiring full power-offDecommissioning devices before storage or shipmentEmergency response to prevent data exfiltrationScheduled shutdowns for lab or classroom devicesReducing power consumption in device fleetsPreparing devices for physical relocation
  Platform Support:
  Windows: Fully supported (Windows 10/11, including Home edition)macOS: Supported (requires user-approved MDM or supervised)iOS/iPadOS: Limited support (supervised devices only, rare use case)Android: Not supported for shutdown action
  Shutdown vs Reboot:
  Shutdown: Device powers off completely, requires manual restartReboot: Device automatically restarts, comes back onlineUse shutdown for: Long-term offline, security incidents, energy savingsUse reboot for: Updates, troubleshooting, configuration changes
  Best Practices:
  Only use when device needs to remain offlineEnsure physical access is available to power device back onNotify users before shutdown (device will be offline)Schedule for end of business day or weekendsDocument reason for shutdown in change managementVerify device location before shutdown (ensure accessibility)Consider reboot instead if device needs to come back onlineTest with small groups before bulk operations
  User Impact:
  Users lose all unsaved workDevice becomes completely unavailableActive sessions are terminatedPhysical access required to power device back onMay cause significant disruption to user productivityUsers cannot access device remotely after shutdown
  Reference: Microsoft Graph API - Shutdown https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-shutdown?view=graph-rest-beta
---

# microsoft365_graph_beta_device_management_managed_device_shutdown (Action)

Remotely shuts down managed devices using the `/deviceManagement/managedDevices/{managedDeviceId}/shutDown` endpoint. This action powers off devices completely, which is useful for energy conservation, maintenance operations, or security scenarios. Unlike reboot, shutdown powers the device off completely and requires manual intervention to power it back on.

**Important Notes:**
- Device shuts down completely (powers off)
- Device requires manual power-on to restart
- Any unsaved work on the device will be lost
- Users receive minimal or no warning before shutdown
- Shutdown is forceful and does not wait for user interaction
- Command is queued if device is offline
- Use with extreme caution - device will be completely offline

**Use Cases:**
- Energy conservation during extended non-use periods
- Security incident response (isolate compromised device)
- Hardware maintenance requiring full power-off
- Decommissioning devices before storage or shipment
- Emergency response to prevent data exfiltration
- Scheduled shutdowns for lab or classroom devices
- Reducing power consumption in device fleets
- Preparing devices for physical relocation

**Platform Support:**
- **Windows**: Fully supported (Windows 10/11, including Home edition)
- **macOS**: Supported (requires user-approved MDM or supervised)
- **iOS/iPadOS**: Limited support (supervised devices only, rare use case)
- **Android**: Not supported for shutdown action

**Shutdown vs Reboot:**
- **Shutdown**: Device powers off completely, requires manual restart
- **Reboot**: Device automatically restarts, comes back online
- Use shutdown for: Long-term offline, security incidents, energy savings
- Use reboot for: Updates, troubleshooting, configuration changes

**Best Practices:**
- Only use when device needs to remain offline
- Ensure physical access is available to power device back on
- Notify users before shutdown (device will be offline)
- Schedule for end of business day or weekends
- Document reason for shutdown in change management
- Verify device location before shutdown (ensure accessibility)
- Consider reboot instead if device needs to come back online
- Test with small groups before bulk operations

**User Impact:**
- Users lose all unsaved work
- Device becomes completely unavailable
- Active sessions are terminated
- Physical access required to power device back on
- May cause significant disruption to user productivity
- Users cannot access device remotely after shutdown

**Reference:** [Microsoft Graph API - Shutdown](https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-shutdown?view=graph-rest-beta)

## Microsoft Documentation

### Graph API References
- [shutDown action](https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-shutdown?view=graph-rest-beta)
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
| v0.40.0-alpha | Experimental | Example fixes and refactored sync progress logic |


## Notes

### Platform Compatibility

| Platform | Support | Requirements |
|----------|---------|--------------|
| **Windows** | ✅ Full Support | All versions including Home |
| **macOS** | ✅ Supported | User-approved MDM or supervised |
| **iOS** | ⚠️ Note | Supervised devices only |
| **iPadOS** | ⚠️ Note | Supervised devices only |
| **Android** | ❌ Not Supported | Shutdown not available |

### ⚠️ Note

**SHUTDOWN POWERS DEVICES OFF COMPLETELY**
- Devices will NOT restart automatically
- Physical access required to power devices back on
- More disruptive than reboot
- Use with extreme caution
- Consider reboot action if devices need to come back online

### Shutdown vs Reboot

| Action | Result | Recovery | Use Case |
|--------|--------|----------|----------|
| **Shutdown** | Powers off | Manual restart required | Long-term offline, energy conservation |
| **Reboot** | Restarts | Automatic (2-5 min) | Updates, troubleshooting, config changes |

### Common Use Cases

- Energy conservation (weekends, holidays)
- Security incident response (device isolation)
- Hardware maintenance requiring full power-off
- Decommissioning devices before storage/shipment
- Emergency response to prevent data exfiltration
- Scheduled shutdowns for lab/classroom devices
- Extended maintenance periods
- Device preparation for physical relocation

### User Impact - CRITICAL

- Users lose ALL unsaved work
- Device becomes COMPLETELY unavailable
- Physical access required to restart
- Significant productivity loss possible
- Users cannot access device remotely
- Active sessions terminated immediately
- No automatic recovery
- Device remains offline indefinitely

### Best Practices

- ONLY use when devices must remain offline
- Ensure physical access available for restart
- Notify users well in advance
- Schedule for end of day or weekends
- Document reason in change management
- Verify device location (ensure accessible)
- **Consider reboot instead whenever possible**
- Test with small groups first
- Have rollback plan (manual power-on procedure)

### Prerequisites Before Shutdown

- Confirm physical access for power-on
- User notification completed
- Business justification documented
- Management approval (if required)
- Backup power-on procedure ready
- Contact information for on-site staff
- Emergency access plan

### Alternatives to Consider

- **Reboot**: If device needs to come back online
- **Sleep/Hibernate**: For temporary offline
- **Network isolation**: For security without full shutdown
- **Remote lock**: To prevent use without power-off
- **Lost mode**: For iOS/iPadOS devices

## Example Usage

```terraform
# Example 1: Shutdown a single device - Minimal
action "microsoft365_graph_beta_device_management_managed_device_shutdown" "shutdown_single" {
  config {
    device_ids = [
      "12345678-1234-1234-1234-123456789abc"
    ]
  }
}

# Example 2: Shutdown multiple devices
action "microsoft365_graph_beta_device_management_managed_device_shutdown" "shutdown_batch" {
  config {
    device_ids = [
      "12345678-1234-1234-1234-123456789abc",
      "87654321-4321-4321-4321-ba9876543210",
      "abcdef12-3456-7890-abcd-ef1234567890"
    ]

    timeouts = {
      invoke = "10m"
    }
  }
}

# Example 3: Shutdown with validation and failure handling - Maximal
action "microsoft365_graph_beta_device_management_managed_device_shutdown" "shutdown_with_validation" {
  config {
    device_ids = [
      "12345678-1234-1234-1234-123456789abc",
      "87654321-4321-4321-4321-ba9876543210",
      "abcdef12-3456-7890-abcd-ef1234567890"
    ]

    ignore_partial_failures = true
    validate_device_exists  = true

    timeouts = {
      invoke = "5m"
    }
  }
}

# Example 4: Shutdown lab devices for weekend energy conservation
data "microsoft365_graph_beta_device_management_managed_device" "lab_devices" {
  filter_type  = "odata"
  odata_filter = "startsWith(deviceName, 'LAB-')"
}

action "microsoft365_graph_beta_device_management_managed_device_shutdown" "shutdown_lab_weekend" {
  config {
    device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.lab_devices.items : device.id]

    timeouts = {
      invoke = "10m"
    }
  }
}

# Example 5: Emergency shutdown for specific device
action "microsoft365_graph_beta_device_management_managed_device_shutdown" "emergency_shutdown" {
  config {
    device_ids = [
      "12345678-abcd-1234-abcd-123456789def"
    ]

    timeouts = {
      invoke = "2m"
    }
  }
}

# Example 6: Shutdown kiosk devices overnight
data "microsoft365_graph_beta_device_management_managed_device" "kiosk_devices" {
  filter_type  = "odata"
  odata_filter = "startsWith(deviceName, 'KIOSK-')"
}

action "microsoft365_graph_beta_device_management_managed_device_shutdown" "shutdown_kiosks_overnight" {
  config {
    device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.kiosk_devices.items : device.id]

    ignore_partial_failures = true

    timeouts = {
      invoke = "10m"
    }
  }
}

# Output examples
output "shutdown_device_count" {
  value       = length(action.microsoft365_graph_beta_device_management_managed_device_shutdown.shutdown_batch.config.device_ids)
  description = "Number of devices that received shutdown command"
}

output "lab_shutdown_count" {
  value       = length(action.microsoft365_graph_beta_device_management_managed_device_shutdown.shutdown_lab_weekend.config.device_ids)
  description = "Number of lab devices shut down for energy conservation"
}

# Important Note: 
# Shutdown powers devices OFF completely and requires manual power-on to restart.
# Use with caution. Consider using reboot action if devices need to come back online automatically.
```

<!-- action schema generated by tfplugindocs -->
## Schema

### Required

- `device_ids` (List of String) List of managed device IDs to shut down. Each ID must be a valid GUID format. Multiple devices can be shut down in a single action. Example: `["12345678-1234-1234-1234-123456789abc", "87654321-4321-4321-4321-ba9876543210"]`

**Critical Warning:** Devices will power off completely when they receive this command. Physical access will be required to power devices back on. Any unsaved work will be lost. Use this action only when devices need to remain powered off.

### Optional

- `ignore_partial_failures` (Boolean) When set to `true`, the action will complete successfully even if some devices fail to shut down. When `false` (default), the action will fail if any device shutdown fails. Use this flag when shutting down multiple devices and you want the action to succeed even if some shutdowns fail.
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))
- `validate_device_exists` (Boolean) When set to `true` (default), the action will validate that all specified devices exist and support shutdown before attempting to shut them down. When `false`, device validation is skipped and the action will attempt to shut down devices directly. Disabling validation can improve performance but may result in errors if devices don't exist or are unsupported.

<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `invoke` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).

