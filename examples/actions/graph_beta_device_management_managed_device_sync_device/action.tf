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

