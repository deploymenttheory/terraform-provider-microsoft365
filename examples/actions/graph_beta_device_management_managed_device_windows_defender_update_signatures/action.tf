# ============================================================================
# Example 1: Update signatures on managed devices only
# ============================================================================
# Use case: Force signature update on fully Intune-managed Windows devices
action "microsoft365_graph_beta_device_management_managed_device_windows_defender_update_signatures" "update_managed_only" {

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
# Example 2: Update signatures on co-managed devices only
# ============================================================================
# Use case: Update definitions on devices managed by both Intune and ConfigMgr
action "microsoft365_graph_beta_device_management_managed_device_windows_defender_update_signatures" "update_comanaged_only" {

  comanaged_device_ids = [
    "12345678-1234-1234-1234-123456789abc",
    "87654321-4321-4321-4321-ba9876543210"
  ]

  timeouts = {
    invoke = "10m"
  }
}

# ============================================================================
# Example 3: Update both managed and co-managed devices
# ============================================================================
# Use case: Mixed environment with both device types
action "microsoft365_graph_beta_device_management_managed_device_windows_defender_update_signatures" "update_mixed_devices" {

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
# Example 4: Update all Windows devices using datasource
# ============================================================================
# Use case: Emergency update after new threat discovered
data "microsoft365_graph_beta_device_management_managed_device" "all_windows" {
  filter_type  = "odata"
  odata_filter = "operatingSystem eq 'Windows'"
}

action "microsoft365_graph_beta_device_management_managed_device_windows_defender_update_signatures" "update_all_windows" {

  managed_device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.all_windows.items : device.id]

  timeouts = {
    invoke = "30m"
  }
}

# ============================================================================
# Example 5: Update signatures before scheduled scan
# ============================================================================
# Use case: Ensure latest definitions before running antivirus scans
data "microsoft365_graph_beta_device_management_managed_device" "workstations" {
  filter_type  = "device_name"
  filter_value = "WKSTN-"
}

# First, update signatures
action "microsoft365_graph_beta_device_management_managed_device_windows_defender_update_signatures" "pre_scan_update" {

  managed_device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.workstations.items : device.id]

  timeouts = {
    invoke = "15m"
  }
}

# Then, run full scan (would need to wait for signature update to complete)
action "microsoft365_graph_beta_device_management_managed_device_windows_defender_scan" "post_update_scan" {

  managed_devices = [
    for device in data.microsoft365_graph_beta_device_management_managed_device.workstations.items : {
      device_id  = device.id
      quick_scan = false
    }
  ]

  timeouts = {
    invoke = "20m"
  }

  # In practice, you'd want to ensure signature update completes first
  depends_on = [action.microsoft365_graph_beta_device_management_managed_device_windows_defender_update_signatures.pre_scan_update]
}

# ============================================================================
# Example 6: Update devices with outdated definitions
# ============================================================================
# Use case: Target devices that haven't updated recently
locals {
  # Example list of devices with outdated signatures
  devices_need_update = [
    "12345678-1234-1234-1234-123456789abc",
    "87654321-4321-4321-4321-ba9876543210",
    "abcdef12-3456-7890-abcd-ef1234567890"
  ]
}

action "microsoft365_graph_beta_device_management_managed_device_windows_defender_update_signatures" "update_outdated" {

  managed_device_ids = local.devices_need_update

  timeouts = {
    invoke = "15m"
  }
}

# ============================================================================
# Example 7: Department-specific update
# ============================================================================
# Use case: Update signatures for specific department or location
data "microsoft365_graph_beta_device_management_managed_device" "finance_devices" {
  filter_type  = "device_name"
  filter_value = "FIN-"
}

action "microsoft365_graph_beta_device_management_managed_device_windows_defender_update_signatures" "update_finance_dept" {

  managed_device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.finance_devices.items : device.id]

  timeouts = {
    invoke = "10m"
  }
}

# ============================================================================
# Example 8: Update after threat intelligence alert
# ============================================================================
# Use case: Zero-day threat response - immediate update across fleet
data "microsoft365_graph_beta_device_management_managed_device" "all_windows_devices" {
  filter_type  = "odata"
  odata_filter = "operatingSystem eq 'Windows'"
}

data "microsoft365_graph_beta_device_management_managed_device" "all_comanaged" {
  filter_type  = "odata"
  odata_filter = "(operatingSystem eq 'Windows') and (managementAgent eq 'configurationManagerClientMdm')"
}

action "microsoft365_graph_beta_device_management_managed_device_windows_defender_update_signatures" "emergency_threat_response" {

  managed_device_ids   = [for device in data.microsoft365_graph_beta_device_management_managed_device.all_windows_devices.items : device.id]
  comanaged_device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.all_comanaged.items : device.id]

  timeouts = {
    invoke = "60m"
  }
}

# ============================================================================
# Example 9: Update Windows Servers only
# ============================================================================
# Use case: Ensure server infrastructure has latest threat definitions
data "microsoft365_graph_beta_device_management_managed_device" "windows_servers" {
  filter_type  = "odata"
  odata_filter = "(operatingSystem eq 'Windows') and (contains(deviceName, 'SRV'))"
}

action "microsoft365_graph_beta_device_management_managed_device_windows_defender_update_signatures" "update_servers" {

  managed_device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.windows_servers.items : device.id]

  timeouts = {
    invoke = "15m"
  }
}

# ============================================================================
# Example 10: Update non-compliant devices
# ============================================================================
# Use case: Force update on non-compliant devices to help remediation
data "microsoft365_graph_beta_device_management_managed_device" "non_compliant_windows" {
  filter_type  = "odata"
  odata_filter = "(operatingSystem eq 'Windows') and (complianceState eq 'noncompliant')"
}

action "microsoft365_graph_beta_device_management_managed_device_windows_defender_update_signatures" "update_non_compliant" {

  managed_device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.non_compliant_windows.items : device.id]

  timeouts = {
    invoke = "20m"
  }
}

# ============================================================================
# Example 11: Update by user assignment
# ============================================================================
# Use case: Update all Windows devices for specific user
data "microsoft365_graph_beta_device_management_managed_device" "user_devices" {
  filter_type  = "odata"
  odata_filter = "(operatingSystem eq 'Windows') and (userPrincipalName eq 'john.doe@company.com')"
}

action "microsoft365_graph_beta_device_management_managed_device_windows_defender_update_signatures" "update_user_devices" {

  managed_device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.user_devices.items : device.id]

  timeouts = {
    invoke = "10m"
  }
}

# ============================================================================
# Example 12: Scheduled monthly update (using Terraform Cloud/Enterprise)
# ============================================================================
# Use case: Regular maintenance - ensure all devices have current definitions
data "microsoft365_graph_beta_device_management_managed_device" "all_managed_windows" {
  filter_type  = "odata"
  odata_filter = "operatingSystem eq 'Windows'"
}

action "microsoft365_graph_beta_device_management_managed_device_windows_defender_update_signatures" "monthly_signature_refresh" {

  managed_device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.all_managed_windows.items : device.id]

  timeouts = {
    invoke = "45m"
  }
}

# ============================================================================
# Example 13: Compliance preparation
# ============================================================================
# Use case: Update signatures before compliance audit
data "microsoft365_graph_beta_device_management_managed_device" "audit_scope_devices" {
  filter_type  = "device_name"
  filter_value = "AUDIT-"
}

action "microsoft365_graph_beta_device_management_managed_device_windows_defender_update_signatures" "pre_audit_update" {

  managed_device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.audit_scope_devices.items : device.id]

  timeouts = {
    invoke = "15m"
  }
}

