# Example 1: Disable a single device - Minimal
action "microsoft365_graph_beta_device_management_managed_device_disable" "disable_single" {
  config {
    managed_device_ids = [
      "12345678-1234-1234-1234-123456789abc"
    ]
  }
}

# Example 2: Disable multiple devices
action "microsoft365_graph_beta_device_management_managed_device_disable" "disable_multiple" {
  config {
    managed_device_ids = [
      "12345678-1234-1234-1234-123456789abc",
      "87654321-4321-4321-4321-ba9876543210",
      "abcdef12-3456-7890-abcd-ef1234567890"
    ]

    timeouts = {
      invoke = "10m"
    }
  }
}

# Example 3: Disable with validation - Maximal
action "microsoft365_graph_beta_device_management_managed_device_disable" "disable_with_validation" {
  config {
    managed_device_ids = [
      "12345678-1234-1234-1234-123456789abc",
      "87654321-4321-4321-4321-ba9876543210"
    ]

    comanaged_device_ids = [
      "abcdef12-3456-7890-abcd-ef1234567890"
    ]

    ignore_partial_failures = true
    validate_device_exists  = true

    timeouts = {
      invoke = "5m"
    }
  }
}

# Example 4: Disable devices due to security incident
variable "compromised_devices" {
  description = "Device IDs suspected of compromise"
  type        = list(string)
  default = [
    "aaaa1111-1111-1111-1111-111111111111",
    "bbbb2222-2222-2222-2222-222222222222"
  ]
}

action "microsoft365_graph_beta_device_management_managed_device_disable" "security_incident" {
  config {
    managed_device_ids = var.compromised_devices

    validate_device_exists  = true
    ignore_partial_failures = false

    timeouts = {
      invoke = "10m"
    }
  }
}

# Example 5: Disable non-compliant devices
data "microsoft365_graph_beta_device_management_managed_device" "non_compliant" {
  filter_type  = "odata"
  odata_filter = "complianceState eq 'noncompliant'"
}

action "microsoft365_graph_beta_device_management_managed_device_disable" "compliance_enforcement" {
  config {
    managed_device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.non_compliant.items : device.id]

    ignore_partial_failures = true

    timeouts = {
      invoke = "20m"
    }
  }
}

# Example 6: Disable co-managed device
action "microsoft365_graph_beta_device_management_managed_device_disable" "disable_comanaged" {
  config {
    comanaged_device_ids = [
      "abcdef12-3456-7890-abcd-ef1234567890"
    ]

    timeouts = {
      invoke = "5m"
    }
  }
}

# Output examples
output "disabled_devices_count" {
  value = {
    managed   = length(action.microsoft365_graph_beta_device_management_managed_device_disable.disable_multiple.config.managed_device_ids)
    comanaged = length(action.microsoft365_graph_beta_device_management_managed_device_disable.disable_comanaged.config.comanaged_device_ids)
  }
  description = "Count of devices disabled"
}
