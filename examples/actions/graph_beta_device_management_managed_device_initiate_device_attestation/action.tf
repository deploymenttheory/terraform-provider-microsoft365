# Example 1: Initiate device attestation for single device
action "microsoft365_graph_beta_device_management_managed_device_initiate_device_attestation" "single_device" {
  managed_device_ids = ["12345678-1234-1234-1234-123456789abc"]

  timeouts = {
    invoke = "5m"
  }
}

# Example 2: Initiate attestation for multiple devices
action "microsoft365_graph_beta_device_management_managed_device_initiate_device_attestation" "multiple_devices" {
  managed_device_ids = [
    "12345678-1234-1234-1234-123456789abc",
    "87654321-4321-4321-4321-ba9876543210",
    "abcdef12-3456-7890-abcd-ef1234567890"
  ]

  timeouts = {
    invoke = "10m"
  }
}

# Example 3: Conditional access compliance check
variable "conditional_access_devices" {
  description = "Device IDs requiring attestation for conditional access"
  type        = list(string)
  default = [
    "aaaa1111-1111-1111-1111-111111111111",
    "bbbb2222-2222-2222-2222-222222222222"
  ]
}

action "microsoft365_graph_beta_device_management_managed_device_initiate_device_attestation" "conditional_access" {
  managed_device_ids = var.conditional_access_devices

  timeouts = {
    invoke = "10m"
  }
}

# Example 4: Zero Trust security validation
data "microsoft365_graph_beta_device_management_managed_device" "zero_trust_devices" {
  filter_type  = "odata"
  odata_filter = "operatingSystem eq 'Windows' and deviceCategoryDisplayName eq 'Zero Trust'"
}

action "microsoft365_graph_beta_device_management_managed_device_initiate_device_attestation" "zero_trust_validation" {
  managed_device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.zero_trust_devices.items : device.id]

  timeouts = {
    invoke = "20m"
  }
}

# Example 5: Periodic compliance verification
locals {
  compliance_scope_devices = [
    "11111111-1111-1111-1111-111111111111",
    "22222222-2222-2222-2222-222222222222",
    "33333333-3333-3333-3333-333333333333"
  ]
}

action "microsoft365_graph_beta_device_management_managed_device_initiate_device_attestation" "periodic_compliance" {
  managed_device_ids = local.compliance_scope_devices

  timeouts = {
    invoke = "15m"
  }
}

# Example 6: Co-managed device attestation
action "microsoft365_graph_beta_device_management_managed_device_initiate_device_attestation" "comanaged_attestation" {
  comanaged_device_ids = ["abcdef12-3456-7890-abcd-ef1234567890"]

  timeouts = {
    invoke = "5m"
  }
}

# Example 7: Post-incident device validation
data "microsoft365_graph_beta_device_management_managed_device" "incident_devices" {
  filter_type  = "odata"
  odata_filter = "deviceCategoryDisplayName eq 'Post-Incident Validation'"
}

action "microsoft365_graph_beta_device_management_managed_device_initiate_device_attestation" "incident_validation" {
  managed_device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.incident_devices.items : device.id]

  timeouts = {
    invoke = "30m"
  }
}

# Example 8: Secure workstation validation
locals {
  secure_workstations = {
    "secure_ws_1" = "11111111-1111-1111-1111-111111111111"
    "secure_ws_2" = "22222222-2222-2222-2222-222222222222"
    "secure_ws_3" = "33333333-3333-3333-3333-333333333333"
  }
}

action "microsoft365_graph_beta_device_management_managed_device_initiate_device_attestation" "secure_workstations" {
  managed_device_ids = values(local.secure_workstations)

  timeouts = {
    invoke = "15m"
  }
}

# Example 9: Pre-deployment security check
data "microsoft365_graph_beta_device_management_managed_device" "pre_deployment" {
  filter_type  = "odata"
  odata_filter = "deviceCategoryDisplayName eq 'Pre-Deployment'"
}

action "microsoft365_graph_beta_device_management_managed_device_initiate_device_attestation" "pre_deployment_check" {
  managed_device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.pre_deployment.items : device.id]

  timeouts = {
    invoke = "20m"
  }
}

# Example 10: VIP device attestation
locals {
  vip_devices = [
    "vip01-1111-1111-1111-111111111111",
    "vip02-2222-2222-2222-222222222222"
  ]
}

action "microsoft365_graph_beta_device_management_managed_device_initiate_device_attestation" "vip_attestation" {
  managed_device_ids = local.vip_devices

  timeouts = {
    invoke = "10m"
  }
}

# Output examples
output "attestation_summary" {
  value = {
    managed   = length(action.multiple_devices.managed_device_ids)
    comanaged = length(action.comanaged_attestation.comanaged_device_ids)
  }
  description = "Count of devices with attestation initiated"
}

