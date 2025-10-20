# Example 1: Refresh machine policy on a single co-managed device
action "microsoft365_graph_beta_device_management_managed_device_trigger_configuration_manager_action" "refresh_single_policy" {

  comanaged_devices {
    device_id = "12345678-1234-1234-1234-123456789abc"
    action    = "refreshMachinePolicy"
  }

  timeouts = {
    invoke = "5m"
  }
}

# Example 2: Trigger app evaluation on multiple devices
action "microsoft365_graph_beta_device_management_managed_device_trigger_configuration_manager_action" "evaluate_apps" {

  comanaged_devices {
    device_id = "12345678-1234-1234-1234-123456789abc"
    action    = "appEvaluation"
  }

  comanaged_devices {
    device_id = "87654321-4321-4321-4321-ba9876543210"
    action    = "appEvaluation"
  }

  timeouts = {
    invoke = "10m"
  }
}

# Example 3: Different actions for different devices
action "microsoft365_graph_beta_device_management_managed_device_trigger_configuration_manager_action" "mixed_actions" {

  comanaged_devices {
    device_id = "aaaa1111-1111-1111-1111-111111111111"
    action    = "refreshMachinePolicy"
  }

  comanaged_devices {
    device_id = "bbbb2222-2222-2222-2222-222222222222"
    action    = "appEvaluation"
  }

  comanaged_devices {
    device_id = "cccc3333-3333-3333-3333-333333333333"
    action    = "quickScan"
  }

  timeouts = {
    invoke = "15m"
  }
}

# Example 4: Wake up all co-managed devices in a collection
variable "sccm_collection_devices" {
  description = "Device IDs from an SCCM collection"
  type        = list(string)
  default = [
    "11111111-1111-1111-1111-111111111111",
    "22222222-2222-2222-2222-222222222222",
    "33333333-3333-3333-3333-333333333333"
  ]
}

action "microsoft365_graph_beta_device_management_managed_device_trigger_configuration_manager_action" "wake_collection" {

  dynamic "comanaged_devices" {
    for_each = var.sccm_collection_devices
    content {
      device_id = comanaged_devices.value
      action    = "wakeUpClient"
    }
  }

  timeouts = {
    invoke = "10m"
  }
}

# Example 5: Trigger Windows Defender signature updates
data "microsoft365_graph_beta_device_management_managed_device" "comanaged_windows" {
  filter_type  = "odata"
  odata_filter = "operatingSystem eq 'Windows' and managementAgent eq 'configurationManagerClient'"
}

action "microsoft365_graph_beta_device_management_managed_device_trigger_configuration_manager_action" "update_defender" {

  dynamic "comanaged_devices" {
    for_each = data.microsoft365_graph_beta_device_management_managed_device.comanaged_windows.items
    content {
      device_id = comanaged_devices.value.id
      action    = "windowsDefenderUpdateSignatures"
    }
  }

  timeouts = {
    invoke = "20m"
  }
}

# Example 6: Refresh user policies after configuration changes
locals {
  user_policy_devices = [
    "aaaa1111-1111-1111-1111-111111111111",
    "bbbb2222-2222-2222-2222-222222222222"
  ]
}

action "microsoft365_graph_beta_device_management_managed_device_trigger_configuration_manager_action" "refresh_user_policies" {

  dynamic "comanaged_devices" {
    for_each = local.user_policy_devices
    content {
      device_id = comanaged_devices.value
      action    = "refreshUserPolicy"
    }
  }

  timeouts = {
    invoke = "10m"
  }
}

# Example 7: Perform full antivirus scan on specific devices
locals {
  security_scan_devices = {
    "device1" = "11111111-1111-1111-1111-111111111111"
    "device2" = "22222222-2222-2222-2222-222222222222"
  }
}

action "microsoft365_graph_beta_device_management_managed_device_trigger_configuration_manager_action" "full_scan" {

  dynamic "comanaged_devices" {
    for_each = local.security_scan_devices
    content {
      device_id = comanaged_devices.value
      action    = "fullScan"
    }
  }

  timeouts = {
    invoke = "30m"
  }
}

# Example 8: Scheduled maintenance - app evaluation for all devices
data "microsoft365_graph_beta_device_management_managed_device" "all_comanaged" {
  filter_type  = "odata"
  odata_filter = "managementAgent eq 'configurationManagerClient'"
}

action "microsoft365_graph_beta_device_management_managed_device_trigger_configuration_manager_action" "scheduled_app_eval" {

  dynamic "comanaged_devices" {
    for_each = { for device in data.microsoft365_graph_beta_device_management_managed_device.all_comanaged.items : device.id => device }
    content {
      device_id = comanaged_devices.key
      action    = "appEvaluation"
    }
  }

  timeouts = {
    invoke = "60m"
  }
}

# Output examples
output "action_summary" {
  value = {
    managed_devices   = length(action.mixed_actions.managed_devices)
    comanaged_devices = length(action.mixed_actions.comanaged_devices)
  }
  description = "Count of devices that had Configuration Manager actions triggered"
}

# Important Notes:
# Configuration Manager Action Features:
# - Triggers specific actions on Configuration Manager client
# - Supports co-managed scenarios (Intune + SCCM)
# - Different actions available for different purposes
# - Device must have Configuration Manager client installed
# - Actions execute on the client side
# - No response data (204 No Content)
#
# Available Actions:
# - refreshMachinePolicy: Refresh machine-level policies
#   - Updates device configuration
#   - Applies new policies from SCCM
#   - Useful after policy changes
#
# - refreshUserPolicy: Refresh user-level policies
#   - Updates user-specific configuration
#   - Applies user-targeted policies
#   - Use when user policies change
#
# - wakeUpClient: Wake up Configuration Manager client
#   - Activates client for immediate processing
#   - Useful before scheduled operations
#   - Ensures client is ready
#
# - appEvaluation: Trigger application evaluation
#   - Re-evaluates application deployments
#   - Checks for new app assignments
#   - Initiates installation if needed
#   - Use after deploying new apps
#
# - quickScan: Windows Defender quick scan
#   - Performs quick antivirus scan
#   - Scans common locations
#   - Faster than full scan
#   - Good for routine checks
#
# - fullScan: Windows Defender full scan
#   - Performs comprehensive antivirus scan
#   - Scans all files and folders
#   - Takes longer to complete
#   - Use for thorough security checks
#
# - windowsDefenderUpdateSignatures: Update antivirus signatures
#   - Downloads latest virus definitions
#   - Updates Windows Defender
#   - Critical for security
#   - Run before scans
#
# When to Use Configuration Manager Actions:
# - After policy or configuration changes
# - Before scheduled maintenance windows
# - After deploying new applications
# - During security incident response
# - For proactive maintenance
# - When troubleshooting device issues
# - To force immediate synchronization
#
# What Happens When Action is Triggered:
# - Command sent to device via Intune
# - Configuration Manager client receives trigger
# - Client executes requested action
# - Action runs based on client schedule
# - No direct response from action
# - Check device status in SCCM or Intune
#
# Platform Requirements:
# - Windows: Required (Configuration Manager is Windows-only)
# - Configuration Manager client must be installed
# - Device must be enrolled in Intune
# - Co-management enabled (recommended)
# - Device must be online
# - Proper permissions required
#
# Co-Management Context:
# - Intune and Configuration Manager together
# - Allows hybrid management scenarios
# - Workloads can be split between systems
# - Best of both platforms
# - Gradual migration path to Intune
# - Leverages existing SCCM infrastructure
#
# Best Practices:
# - Use appropriate action for the scenario
# - Consider device online status
# - Don't trigger too frequently
# - Monitor action success in logs
# - Coordinate with maintenance windows
# - Document action triggers
# - Test on pilot devices first
#
# Action Execution:
# - Asynchronous operation
# - No immediate feedback
# - Check client logs for results
# - Verify in Configuration Manager console
# - May take time to complete
# - Depends on client schedule
#
# Troubleshooting:
# - Verify device has SCCM client
# - Check device is online
# - Ensure co-management is configured
# - Review client logs (ccmexec.log)
# - Check Intune portal for status
# - Verify permissions are correct
# - Confirm network connectivity
#
# Common Use Cases:
# - Force policy refresh after changes
# - Immediate app deployment evaluation
# - Security scans on demand
# - Wake devices before operations
# - Update antivirus definitions
# - User policy updates after login
# - Troubleshooting configuration issues
#
# Action Timing:
# - Immediate trigger sent
# - Client processes based on schedule
# - Some actions run immediately
# - Others queue for next cycle
# - Network conditions affect timing
# - Client must be reachable
#
# Monitoring Results:
# - Check Configuration Manager console
# - Review Intune device details
# - Examine client logs on device
# - Monitor sync status
# - Verify expected outcomes
# - Check for error messages
#
# Limitations:
# - Windows devices only
# - Requires SCCM client
# - Co-management recommended
# - No direct response data
# - Asynchronous execution
# - Network dependent
# - Client schedule affects timing
#
# Related Actions:
# - sync_device: Intune device sync
# - reboot_now: Force device restart
# - Device compliance checks
# - Policy application monitoring
#
# Integration:
# - Part of co-management strategy
# - Complements Intune actions
# - Leverages SCCM infrastructure
# - Enables hybrid scenarios
# - Supports migration to cloud
#
# Reference:
# https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-triggerconfigurationmanageraction?view=graph-rest-beta

