# Example 1: Set device name for a single device - Minimal
action "microsoft365_graph_beta_device_management_managed_device_set_device_name" "set_name_single" {
  config {
    managed_devices = [
      {
        device_id   = "12345678-1234-1234-1234-123456789abc"
        device_name = "NYC-Marketing-Laptop-01"
      }
    ]
  }
}

# Example 2: Set device names for multiple devices
action "microsoft365_graph_beta_device_management_managed_device_set_device_name" "set_names_multiple" {
  config {
    managed_devices = [
      {
        device_id   = "12345678-1234-1234-1234-123456789abc"
        device_name = "NYC-Floor3-Conf-Room-01"
      },
      {
        device_id   = "87654321-4321-4321-4321-ba9876543210"
        device_name = "NYC-Floor3-Conf-Room-02"
      },
      {
        device_id   = "abcdef12-3456-7890-abcd-ef1234567890"
        device_name = "NYC-Floor3-Conf-Room-03"
      }
    ]

    timeouts = {
      invoke = "10m"
    }
  }
}

# Example 3: Maximal configuration with validation
action "microsoft365_graph_beta_device_management_managed_device_set_device_name" "set_names_maximal" {
  config {
    managed_devices = [
      {
        device_id   = "12345678-1234-1234-1234-123456789abc"
        device_name = "NYC-Marketing-Laptop-01"
      },
      {
        device_id   = "87654321-4321-4321-4321-ba9876543210"
        device_name = "NYC-IT-Desktop-05"
      }
    ]

    comanaged_devices = [
      {
        device_id   = "abcdef12-3456-7890-abcd-ef1234567890"
        device_name = "NYC-HR-Laptop-03"
      }
    ]

    ignore_partial_failures = true
    validate_device_exists  = true

    timeouts = {
      invoke = "5m"
    }
  }
}

# Example 4: Rename devices based on user assignment
data "microsoft365_graph_beta_device_management_managed_device" "user_devices" {
  filter_type  = "odata"
  odata_filter = "userPrincipalName eq 'john.doe@example.com'"
}

action "microsoft365_graph_beta_device_management_managed_device_set_device_name" "rename_user_devices" {
  config {
    managed_devices = [
      for device in data.microsoft365_graph_beta_device_management_managed_device.user_devices.items : {
        device_id   = device.id
        device_name = format("JohnDoe-%s-%s", device.operatingSystem, substr(device.serialNumber, 0, 8))
      }
    ]

    timeouts = {
      invoke = "15m"
    }
  }
}

# Example 5: Standardize naming for devices by department
data "microsoft365_graph_beta_device_management_managed_device" "it_devices" {
  filter_type  = "odata"
  odata_filter = "deviceCategoryDisplayName eq 'IT Department'"
}

action "microsoft365_graph_beta_device_management_managed_device_set_device_name" "standardize_it_devices" {
  config {
    managed_devices = [
      for device in data.microsoft365_graph_beta_device_management_managed_device.it_devices.items : {
        device_id   = device.id
        device_name = format("IT-DEPT-%s", substr(device.id, 0, 8))
      }
    ]

    validate_device_exists = true

    timeouts = {
      invoke = "20m"
    }
  }
}

# Example 6: Rename devices by location
locals {
  device_locations = {
    "12345678-1234-1234-1234-123456789abc" = "NYC-Office"
    "87654321-4321-4321-4321-ba9876543210" = "LA-Office"
    "abcdef12-3456-7890-abcd-ef1234567890" = "Chicago-Office"
  }
}

action "microsoft365_graph_beta_device_management_managed_device_set_device_name" "rename_by_location" {
  config {
    managed_devices = [
      for device_id, location in local.device_locations : {
        device_id   = device_id
        device_name = format("%s-Device", location)
      }
    ]

    timeouts = {
      invoke = "10m"
    }
  }
}

# Example 7: Set name for co-managed device
action "microsoft365_graph_beta_device_management_managed_device_set_device_name" "set_comanaged_name" {
  config {
    comanaged_devices = [
      {
        device_id   = "abcdef12-3456-7890-abcd-ef1234567890"
        device_name = "SCCM-Intune-Hybrid-01"
      }
    ]

    timeouts = {
      invoke = "5m"
    }
  }
}

# Example 8: Rename devices after asset reassignment
variable "reassigned_devices" {
  description = "Map of device IDs to new names after reassignment"
  type        = map(string)
  default = {
    "11111111-1111-1111-1111-111111111111" = "Finance-Laptop-A"
    "22222222-2222-2222-2222-222222222222" = "Finance-Laptop-B"
  }
}

action "microsoft365_graph_beta_device_management_managed_device_set_device_name" "reassign_devices" {
  config {
    managed_devices = [
      for device_id, device_name in var.reassigned_devices : {
        device_id   = device_id
        device_name = device_name
      }
    ]

    timeouts = {
      invoke = "10m"
    }
  }
}

# Output examples
output "devices_renamed_count" {
  value       = length(action.microsoft365_graph_beta_device_management_managed_device_set_device_name.set_names_multiple.config.managed_devices)
  description = "Number of devices that had their names set"
}

output "device_naming_info" {
  value = {
    managed   = length(action.microsoft365_graph_beta_device_management_managed_device_set_device_name.set_names_maximal.config.managed_devices)
    comanaged = length(action.microsoft365_graph_beta_device_management_managed_device_set_device_name.set_names_maximal.config.comanaged_devices)
  }
  description = "Count of renamed devices by type"
}
