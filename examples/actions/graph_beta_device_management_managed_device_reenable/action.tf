# Example 1: Re-enable a single device - Minimal
action "microsoft365_graph_beta_device_management_managed_device_reenable" "reenable_single" {
  config {
    managed_device_ids = [
      "12345678-1234-1234-1234-123456789abc"
    ]
  }
}

# Example 2: Re-enable multiple devices
action "microsoft365_graph_beta_device_management_managed_device_reenable" "reenable_multiple" {
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

# Example 3: Re-enable with validation - Maximal
action "microsoft365_graph_beta_device_management_managed_device_reenable" "reenable_with_validation" {
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

# Example 4: Re-enable devices after security investigation
variable "investigated_devices" {
  description = "Device IDs cleared from security investigation"
  type        = list(string)
  default = [
    "aaaa1111-1111-1111-1111-111111111111",
    "bbbb2222-2222-2222-2222-222222222222"
  ]
}

action "microsoft365_graph_beta_device_management_managed_device_reenable" "post_investigation" {
  config {
    managed_device_ids = var.investigated_devices

    validate_device_exists = true

    timeouts = {
      invoke = "10m"
    }
  }
}

# Example 5: Re-enable compliant devices
data "microsoft365_graph_beta_device_management_managed_device" "now_compliant" {
  filter_type  = "odata"
  odata_filter = "complianceState eq 'compliant'"
}

action "microsoft365_graph_beta_device_management_managed_device_reenable" "compliance_restored" {
  config {
    managed_device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.now_compliant.items : device.id]

    ignore_partial_failures = true

    timeouts = {
      invoke = "20m"
    }
  }
}

# Example 6: Re-enable co-managed device
action "microsoft365_graph_beta_device_management_managed_device_reenable" "reenable_comanaged" {
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
output "reenabled_devices_count" {
  value = {
    managed   = length(action.microsoft365_graph_beta_device_management_managed_device_reenable.reenable_multiple.config.managed_device_ids)
    comanaged = length(action.microsoft365_graph_beta_device_management_managed_device_reenable.reenable_comanaged.config.comanaged_device_ids)
  }
  description = "Count of devices re-enabled"
}
