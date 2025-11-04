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

