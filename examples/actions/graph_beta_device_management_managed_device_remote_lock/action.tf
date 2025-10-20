# Example 1: Remote lock a single device (lost device scenario)
action "microsoft365_graph_beta_device_management_managed_device_remote_lock" "lock_lost_device" {

  device_ids = [
    "12345678-1234-1234-1234-123456789abc"
  ]

  timeouts = {
    invoke = "5m"
  }
}

# Example 2: Remote lock multiple devices
action "microsoft365_graph_beta_device_management_managed_device_remote_lock" "lock_batch" {

  device_ids = [
    "12345678-1234-1234-1234-123456789abc",
    "87654321-4321-4321-4321-ba9876543210",
    "abcdef12-3456-7890-abcd-ef1234567890"
  ]

  timeouts = {
    invoke = "10m"
  }
}

# Example 3: Lock all devices for a specific user (security incident)
data "microsoft365_graph_beta_device_management_managed_device" "user_devices" {
  filter_type  = "user_id"
  filter_value = "compromised.user@example.com"
}

action "microsoft365_graph_beta_device_management_managed_device_remote_lock" "lock_compromised_user" {

  device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.user_devices.items : device.id]

  timeouts = {
    invoke = "10m"
  }
}

# Example 4: Lock non-compliant devices
data "microsoft365_graph_beta_device_management_managed_device" "non_compliant_devices" {
  filter_type  = "odata"
  odata_filter = "complianceState eq 'noncompliant'"
}

action "microsoft365_graph_beta_device_management_managed_device_remote_lock" "lock_non_compliant_devices" {

  device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.non_compliant_devices.items : device.id]

  timeouts = {
    invoke = "15m"
  }
}

# Example 5: Lock iOS devices reported as lost
data "microsoft365_graph_beta_device_management_managed_device" "ios_devices" {
  filter_type  = "odata"
  odata_filter = "operatingSystem eq 'iOS'"
}

action "microsoft365_graph_beta_device_management_managed_device_remote_lock" "lock_ios" {

  # In production, you would have additional filtering for "lost" status
  device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.ios_devices.items : device.id]

  timeouts = {
    invoke = "10m"
  }
}

# Example 6: Emergency lock all corporate Windows devices
data "microsoft365_graph_beta_device_management_managed_device" "corporate_windows" {
  filter_type  = "odata"
  odata_filter = "operatingSystem eq 'Windows' and managedDeviceOwnerType eq 'company'"
}

action "microsoft365_graph_beta_device_management_managed_device_remote_lock" "lock_corporate_windows" {

  device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.corporate_windows.items : device.id]

  timeouts = {
    invoke = "20m"
  }
}

# Example 7: Lock Android Enterprise devices
data "microsoft365_graph_beta_device_management_managed_device" "android_devices" {
  filter_type  = "odata"
  odata_filter = "operatingSystem eq 'Android'"
}

action "microsoft365_graph_beta_device_management_managed_device_remote_lock" "lock_android" {

  device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.android_devices.items : device.id]

  timeouts = {
    invoke = "15m"
  }
}

# Example 8: Lock devices by device name pattern (department-specific)
data "microsoft365_graph_beta_device_management_managed_device" "department_devices" {
  filter_type  = "device_name"
  filter_value = "SALES-"
}

action "microsoft365_graph_beta_device_management_managed_device_remote_lock" "lock_sales_devices" {

  device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.department_devices.items : device.id]

  timeouts = {
    invoke = "10m"
  }
}

# Output examples
output "locked_device_count" {
  value       = length(action.lock_batch.device_ids)
  description = "Number of devices that received remote lock command"
}

output "emergency_locked_count" {
  value       = length(action.lock_compromised_user.device_ids)
  description = "Number of devices locked in emergency scenario"
}

# Important Note: 
# - Devices lock IMMEDIATELY when they receive the command
# - Users must enter their existing passcode to unlock
# - For lost devices, follow up with locate/wipe if needed
# - Document the reason for locking devices for compliance/audit purposes

