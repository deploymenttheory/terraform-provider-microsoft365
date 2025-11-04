# Example 1: Initiate MDM key recovery for single device
action "microsoft365_graph_beta_device_management_managed_device_initiate_mobile_device_management_key_recovery" "single_device" {
  managed_device_ids = ["12345678-1234-1234-1234-123456789abc"]

  timeouts = {
    invoke = "5m"
  }
}

# Example 2: Initiate key recovery for multiple devices
action "microsoft365_graph_beta_device_management_managed_device_initiate_mobile_device_management_key_recovery" "multiple_devices" {
  managed_device_ids = [
    "12345678-1234-1234-1234-123456789abc",
    "87654321-4321-4321-4321-ba9876543210",
    "abcdef12-3456-7890-abcd-ef1234567890"
  ]

  timeouts = {
    invoke = "10m"
  }
}

# Example 3: Compliance-driven key recovery
variable "compliance_devices" {
  description = "Device IDs requiring BitLocker key escrow for compliance"
  type        = list(string)
  default = [
    "aaaa1111-1111-1111-1111-111111111111",
    "bbbb2222-2222-2222-2222-222222222222"
  ]
}

action "microsoft365_graph_beta_device_management_managed_device_initiate_mobile_device_management_key_recovery" "compliance_escrow" {
  managed_device_ids = var.compliance_devices

  timeouts = {
    invoke = "10m"
  }
}

# Example 4: Periodic key rotation based on data source
data "microsoft365_graph_beta_device_management_managed_device" "encrypted_devices" {
  filter_type  = "odata"
  odata_filter = "operatingSystem eq 'Windows' and deviceCategoryDisplayName eq 'Encrypted'"
}

action "microsoft365_graph_beta_device_management_managed_device_initiate_mobile_device_management_key_recovery" "periodic_rotation" {
  managed_device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.encrypted_devices.items : device.id]

  timeouts = {
    invoke = "20m"
  }
}

# Example 5: New device enrollment key escrow
locals {
  new_enrollment_devices = [
    "11111111-1111-1111-1111-111111111111",
    "22222222-2222-2222-2222-222222222222",
    "33333333-3333-3333-3333-333333333333"
  ]
}

action "microsoft365_graph_beta_device_management_managed_device_initiate_mobile_device_management_key_recovery" "new_device_escrow" {
  managed_device_ids = local.new_enrollment_devices

  timeouts = {
    invoke = "15m"
  }
}

# Example 6: Co-managed device key recovery
action "microsoft365_graph_beta_device_management_managed_device_initiate_mobile_device_management_key_recovery" "comanaged_escrow" {
  comanaged_device_ids = ["abcdef12-3456-7890-abcd-ef1234567890"]

  timeouts = {
    invoke = "5m"
  }
}

# Example 7: Critical devices pre-deployment validation
data "microsoft365_graph_beta_device_management_managed_device" "critical_devices" {
  filter_type  = "odata"
  odata_filter = "deviceCategoryDisplayName eq 'Critical Infrastructure'"
}

action "microsoft365_graph_beta_device_management_managed_device_initiate_mobile_device_management_key_recovery" "critical_validation" {
  managed_device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.critical_devices.items : device.id]

  timeouts = {
    invoke = "30m"
  }
}

# Example 8: Department-specific key recovery
locals {
  finance_department_devices = {
    "finance_laptop_1"  = "11111111-1111-1111-1111-111111111111"
    "finance_laptop_2"  = "22222222-2222-2222-2222-222222222222"
    "finance_desktop_1" = "33333333-3333-3333-3333-333333333333"
  }
}

action "microsoft365_graph_beta_device_management_managed_device_initiate_mobile_device_management_key_recovery" "finance_department" {
  managed_device_ids = values(local.finance_department_devices)

  timeouts = {
    invoke = "15m"
  }
}

# Example 9: Scheduled quarterly key rotation
data "microsoft365_graph_beta_device_management_managed_device" "all_windows_devices" {
  filter_type  = "odata"
  odata_filter = "operatingSystem eq 'Windows'"
}

action "microsoft365_graph_beta_device_management_managed_device_initiate_mobile_device_management_key_recovery" "quarterly_rotation" {
  managed_device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.all_windows_devices.items : device.id]

  timeouts = {
    invoke = "60m"
  }
}

# Example 10: Audit preparation key escrow
locals {
  audit_scope_devices = [
    "audit01-1111-1111-1111-111111111111",
    "audit02-2222-2222-2222-222222222222"
  ]
}

action "microsoft365_graph_beta_device_management_managed_device_initiate_mobile_device_management_key_recovery" "audit_preparation" {
  managed_device_ids = local.audit_scope_devices

  timeouts = {
    invoke = "10m"
  }
}

# Output examples
output "key_recovery_summary" {
  value = {
    managed   = length(action.multiple_devices.managed_device_ids)
    comanaged = length(action.comanaged_escrow.comanaged_device_ids)
  }
  description = "Count of devices with MDM key recovery initiated"
}

