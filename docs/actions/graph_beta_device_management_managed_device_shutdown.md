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

## Notes

### Platform Compatibility

| Platform | Support | Requirements |
|----------|---------|--------------|
| **Windows** | ✅ Full Support | All versions including Home |
| **macOS** | ✅ Supported | User-approved MDM or supervised |
| **iOS** | ⚠️ Limited | Supervised devices only, uncommon |
| **iPadOS** | ⚠️ Limited | Supervised devices only, uncommon |
| **Android** | ❌ Not Supported | Shutdown not available |

### ⚠️ Critical Warning

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
# Example 1: Shutdown a single device
action "microsoft365_graph_beta_device_management_managed_device_shutdown" "shutdown_single" {

  device_ids = [
    "12345678-1234-1234-1234-123456789abc"
  ]

  timeouts = {
    invoke = "5m"
  }
}

# Example 2: Shutdown multiple devices
action "microsoft365_graph_beta_device_management_managed_device_shutdown" "shutdown_batch" {

  device_ids = [
    "12345678-1234-1234-1234-123456789abc",
    "87654321-4321-4321-4321-ba9876543210",
    "abcdef12-3456-7890-abcd-ef1234567890"
  ]

  timeouts = {
    invoke = "10m"
  }
}

# Example 3: Shutdown lab devices for weekend energy conservation
data "microsoft365_graph_beta_device_management_managed_device" "lab_devices" {
  filter_type  = "device_name"
  filter_value = "LAB-"
}

action "microsoft365_graph_beta_device_management_managed_device_shutdown" "shutdown_lab_weekend" {

  device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.lab_devices.items : device.id]

  timeouts = {
    invoke = "10m"
  }
}

# Example 4: Emergency shutdown for specific device by ID
action "microsoft365_graph_beta_device_management_managed_device_shutdown" "emergency_shutdown" {

  device_ids = [
    "12345678-abcd-1234-abcd-123456789def" # Replace with actual compromised device ID
  ]

  timeouts = {
    invoke = "2m"
  }
}

# Example 5: Shutdown kiosk devices overnight
data "microsoft365_graph_beta_device_management_managed_device" "kiosk_devices" {
  filter_type  = "device_name"
  filter_value = "KIOSK-"
}

action "microsoft365_graph_beta_device_management_managed_device_shutdown" "shutdown_kiosks_overnight" {

  device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.kiosk_devices.items : device.id]

  timeouts = {
    invoke = "10m"
  }
}

# Example 6: Shutdown classroom devices for extended break
data "microsoft365_graph_beta_device_management_managed_device" "classroom_devices" {
  filter_type  = "device_name"
  filter_value = "CLASSROOM-"
}

action "microsoft365_graph_beta_device_management_managed_device_shutdown" "shutdown_classroom_break" {

  device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.classroom_devices.items : device.id]

  timeouts = {
    invoke = "15m"
  }
}

# Output examples
output "shutdown_device_count" {
  value       = length(action.shutdown_batch.device_ids)
  description = "Number of devices that received shutdown command"
}

output "lab_shutdown_count" {
  value       = length(action.shutdown_lab_weekend.device_ids)
  description = "Number of lab devices shut down for energy conservation"
}

# Important Notes:
#
# Critical Warning:
# SHUTDOWN POWERS DEVICES OFF COMPLETELY!
# - Devices will NOT restart automatically
# - Physical access required to power devices back on
# - Use with extreme caution
# - Consider using reboot action instead if devices need to come back online
#
# What is Shutdown?
# - Completely powers off devices
# - Requires manual power-on to restart
# - More disruptive than reboot
# - Used when devices need to remain offline
# - Typically for energy conservation or security incidents
#
# When to Use Shutdown (vs Reboot):
# - Energy conservation during extended non-use (weekends, holidays)
# - Security incident response (isolate compromised device)
# - Hardware maintenance requiring full power-off
# - Decommissioning devices before storage/shipment
# - Emergency response to prevent data exfiltration
# - Scheduled shutdowns for lab/classroom devices
# - Preparing devices for physical relocation
# - Extended maintenance periods
#
# Platform Support:
# - Windows: Fully supported (all versions)
# - macOS: Supported (user-approved MDM or supervised)
# - iOS/iPadOS: Limited (supervised only, very rare use case)
# - Android: Not supported
#
# User Impact - CRITICAL:
# - Users lose ALL unsaved work
# - Device becomes COMPLETELY unavailable
# - Physical access required to restart
# - May cause significant productivity loss
# - Users cannot access device remotely
# - Active sessions terminated immediately
# - No automatic recovery
#
# Shutdown vs Reboot Comparison:
#
# Shutdown:
# - Powers off completely
# - Manual restart required
# - Use for: Long-term offline, security, energy savings
# - Impact: Device offline until manually restarted
#
# Reboot:
# - Restarts automatically
# - Device comes back online
# - Use for: Updates, troubleshooting, config changes
# - Impact: Brief downtime (2-5 minutes)
#
# Best Practices:
# - ONLY use when devices must remain offline
# - Ensure physical access available to restart
# - Notify users well in advance
# - Schedule for end of day or weekends
# - Document reason in change management
# - Verify device location (ensure accessible)
# - Consider reboot instead if possible
# - Test with small groups first
# - Have rollback plan (manual power-on process)
#
# Common Use Cases:
#
# Case 1: Weekend Energy Conservation
# - Shutdown lab/classroom devices Friday evening
# - Manually power on Monday morning
# - Reduces energy costs
# - Environmental benefits
#
# Case 2: Security Incident Response
# - Immediately isolate compromised device
# - Prevent continued data exfiltration
# - Preserve forensic evidence (no auto-restart)
# - Part of incident response playbook
#
# Case 3: Extended Maintenance
# - Hardware upgrades requiring physical access
# - Facility maintenance (power/HVAC)
# - Building closures (holidays, renovations)
# - Device relocation projects
#
# Case 4: Device Decommissioning
# - Preparing devices for storage
# - Devices awaiting disposal/recycling
# - Inventory consolidation
# - Asset retirement process
#
# Scheduling Recommendations:
# - Lab devices: Friday 6pm (weekend shutdown)
# - Classroom devices: Before extended breaks
# - Kiosks: During facility closure
# - Corporate: Extended holidays only
# - Emergency: Immediate (security incidents)
#
# Prerequisites:
# - Physical access plan for restart
# - User notification completed
# - Business justification documented
# - Management approval (if required)
# - Backup power-on procedure
# - Contact information for on-site staff
#
# Validation Checklist:
# - [ ] Verified devices need to remain offline
# - [ ] Confirmed reboot won't suffice
# - [ ] Physical access available for restart
# - [ ] Users notified and approved
# - [ ] Change management ticket created
# - [ ] Rollback plan documented
# - [ ] Emergency contact designated
#
# Monitoring:
# - Track which devices were shut down
# - Monitor for unexpected shutdowns
# - Verify devices stay offline as expected
# - Alert if device comes back online unexpectedly
# - Document manual restart times
#
# Troubleshooting:
#
# Issue: Device won't shutdown
# Solution: Check connectivity, verify command received
#
# Issue: Device restarts automatically
# Solution: Check BIOS settings, Wake-on-LAN configuration
#
# Issue: Can't power device back on
# Solution: Check physical power, hardware issues
#
# Issue: Shutdown command fails
# Solution: Verify device online, check permissions
#
# Security Considerations:
# - Shutdown isolates device from network
# - Prevents remote attacks during incident
# - May affect monitoring/security agents
# - Consider forensic evidence preservation
# - Document in security audit logs
# - Part of incident response procedures
#
# Energy Conservation Benefits:
# - Reduces electricity costs
# - Environmental impact reduction
# - Extends hardware lifespan
# - Complies with sustainability policies
# - Reduces cooling/HVAC requirements
#
# Automation Examples:
# - Friday evening shutdown scripts
# - Holiday schedule automation
# - Security playbook integration
# - Facility management coordination
# - Environmental monitoring triggers
#
# Emergency Response Workflow:
# 1. Security incident detected
# 2. Identify affected device(s)
# 3. Issue immediate shutdown command
# 4. Isolate device from network
# 5. Notify security team
# 6. Preserve forensic evidence
# 7. Document incident timeline
# 8. Follow incident response plan
#
# Recovery Procedures:
# - Document manual power-on steps
# - Assign responsibility for restart
# - Verify device connectivity after restart
# - Confirm services restore properly
# - Monitor for issues post-restart
# - Update asset management systems
#
# Approval Requirements:
# - Management approval for bulk shutdowns
# - Security team approval for incidents
# - Facility coordination for physical access
# - User notification and consent
# - Change management board review
# - Business impact assessment
#
# Alternatives to Consider:
# - Reboot: If device needs to come back online
# - Sleep/Hibernate: For temporary offline
# - Network isolation: For security without full shutdown
# - Remote lock: To prevent use without full power-off
# - Lost mode: For iOS/iPadOS devices
#
# Reference:
# https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-shutdown?view=graph-rest-beta
```

<!-- action schema generated by tfplugindocs -->
## Schema

### Required

- `device_ids` (List of String) List of managed device IDs to shut down. Each ID must be a valid GUID format. Multiple devices can be shut down in a single action. Example: `["12345678-1234-1234-1234-123456789abc", "87654321-4321-4321-4321-ba9876543210"]`

**Critical Warning:** Devices will power off completely when they receive this command. Physical access will be required to power devices back on. Any unsaved work will be lost. Use this action only when devices need to remain powered off.

### Optional

- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).

