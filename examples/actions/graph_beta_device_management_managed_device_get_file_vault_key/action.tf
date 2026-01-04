# Example 1: Retrieve FileVault key for a single macOS device - Minimal
action "microsoft365_graph_beta_device_management_managed_device_get_file_vault_key" "retrieve_single" {
  config {
    managed_device_ids = [
      "12345678-1234-1234-1234-123456789abc"
    ]
  }
}

# Example 2: Retrieve FileVault keys for multiple macOS devices
action "microsoft365_graph_beta_device_management_managed_device_get_file_vault_key" "retrieve_multiple" {
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

# Example 3: Retrieve keys with validation - Maximal
action "microsoft365_graph_beta_device_management_managed_device_get_file_vault_key" "retrieve_with_validation" {
  config {
    managed_device_ids = [
      "12345678-1234-1234-1234-123456789abc",
      "87654321-4321-4321-4321-ba9876543210"
    ]

    comanaged_device_ids = [
      "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
    ]

    ignore_partial_failures = true
    validate_device_exists  = true

    timeouts = {
      invoke = "5m"
    }
  }
}

# Example 4: Emergency recovery for locked macOS devices
variable "locked_device_ids" {
  description = "Device IDs for locked macOS devices"
  type        = list(string)
  default = [
    "11111111-1111-1111-1111-111111111111",
    "22222222-2222-2222-2222-222222222222"
  ]
}

action "microsoft365_graph_beta_device_management_managed_device_get_file_vault_key" "emergency_recovery" {
  config {
    managed_device_ids = var.locked_device_ids

    validate_device_exists = true

    timeouts = {
      invoke = "5m"
    }
  }
}

# Example 5: Retrieve keys for departing employee's macOS devices
data "microsoft365_graph_beta_device_management_managed_device" "departing_employee" {
  filter_type  = "odata"
  odata_filter = "(userPrincipalName eq 'departing.employee@example.com') and (operatingSystem eq 'macOS')"
}

action "microsoft365_graph_beta_device_management_managed_device_get_file_vault_key" "departing_employee_recovery" {
  config {
    managed_device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.departing_employee.items : device.id]

    timeouts = {
      invoke = "10m"
    }
  }
}

# Example 6: Retrieve keys for co-managed macOS devices
action "microsoft365_graph_beta_device_management_managed_device_get_file_vault_key" "retrieve_comanaged" {
  config {
    comanaged_device_ids = [
      "abcdef12-3456-7890-abcd-ef1234567890"
    ]

    timeouts = {
      invoke = "5m"
    }
  }
}

# Example 7: Retrieve keys for all company macOS devices
data "microsoft365_graph_beta_device_management_managed_device" "all_macos" {
  filter_type  = "odata"
  odata_filter = "(operatingSystem eq 'macOS') and (managedDeviceOwnerType eq 'company')"
}

action "microsoft365_graph_beta_device_management_managed_device_get_file_vault_key" "all_company_macos" {
  config {
    managed_device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.all_macos.items : device.id]

    ignore_partial_failures = true

    timeouts = {
      invoke = "20m"
    }
  }
}

# Output examples
output "retrieved_keys_count" {
  value       = length(action.microsoft365_graph_beta_device_management_managed_device_get_file_vault_key.retrieve_multiple.config.managed_device_ids)
  description = "Number of devices for which FileVault keys were retrieved"
}
