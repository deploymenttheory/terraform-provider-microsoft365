# Example 1: Deprovision a single device - Minimal
action "microsoft365_graph_beta_device_management_managed_device_deprovision" "deprovision_single" {
  config {
    managed_devices = [
      {
        device_id          = "12345678-1234-1234-1234-123456789abc"
        deprovision_reason = "Device being transitioned to new management solution"
      }
    ]
  }
}

# Example 2: Deprovision multiple devices
action "microsoft365_graph_beta_device_management_managed_device_deprovision" "deprovision_multiple" {
  config {
    managed_devices = [
      {
        device_id          = "12345678-1234-1234-1234-123456789abc"
        deprovision_reason = "Device repurposing for different department"
      },
      {
        device_id          = "87654321-4321-4321-4321-ba9876543210"
        deprovision_reason = "Troubleshooting management issues"
      }
    ]

    timeouts = {
      invoke = "10m"
    }
  }
}

# Example 3: Maximal configuration with validation
action "microsoft365_graph_beta_device_management_managed_device_deprovision" "deprovision_maximal" {
  config {
    managed_devices = [
      {
        device_id          = "12345678-1234-1234-1234-123456789abc"
        deprovision_reason = "Transitioning to new management"
      },
      {
        device_id          = "87654321-4321-4321-4321-987654321cba"
        deprovision_reason = "Device repurposing"
      }
    ]

    comanaged_devices = [
      {
        device_id          = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
        deprovision_reason = "Removing co-management"
      }
    ]

    ignore_partial_failures = false
    validate_device_exists  = true

    timeouts = {
      invoke = "5m"
    }
  }
}

# Example 4: Deprovision devices by user
variable "departing_user_devices" {
  description = "Device IDs for departing user"
  type        = list(string)
  default = [
    "11111111-1111-1111-1111-111111111111",
    "22222222-2222-2222-2222-222222222222"
  ]
}

action "microsoft365_graph_beta_device_management_managed_device_deprovision" "user_departure" {
  config {
    managed_devices = [
      for device_id in var.departing_user_devices : {
        device_id          = device_id
        deprovision_reason = "User departure - removing management policies"
      }
    ]

    timeouts = {
      invoke = "15m"
    }
  }
}

# Example 5: Transition from Intune-only to co-management
data "microsoft365_graph_beta_device_management_managed_device" "transition_devices" {
  filter_type  = "odata"
  odata_filter = "deviceCategoryDisplayName eq 'Co-Management Transition'"
}

action "microsoft365_graph_beta_device_management_managed_device_deprovision" "comanagement_transition" {
  config {
    managed_devices = [
      for device in data.microsoft365_graph_beta_device_management_managed_device.transition_devices.items : {
        device_id          = device.id
        deprovision_reason = "Transitioning to co-management with Configuration Manager"
      }
    ]

    validate_device_exists = true

    timeouts = {
      invoke = "20m"
    }
  }
}

# Example 6: Deprovision co-managed device
action "microsoft365_graph_beta_device_management_managed_device_deprovision" "deprovision_comanaged" {
  config {
    comanaged_devices = [
      {
        device_id          = "abcdef12-3456-7890-abcd-ef1234567890"
        deprovision_reason = "Changing management authority to Configuration Manager only"
      }
    ]

    timeouts = {
      invoke = "5m"
    }
  }
}

# Example 7: Prepare devices for repurposing
data "microsoft365_graph_beta_device_management_managed_device" "repurpose_candidates" {
  filter_type  = "odata"
  odata_filter = "deviceCategoryDisplayName eq 'Repurpose Queue'"
}

action "microsoft365_graph_beta_device_management_managed_device_deprovision" "repurpose_prep" {
  config {
    managed_devices = [
      for device in data.microsoft365_graph_beta_device_management_managed_device.repurpose_candidates.items : {
        device_id          = device.id
        deprovision_reason = format("Repurposing device %s for new deployment", device.deviceName)
      }
    ]

    ignore_partial_failures = true

    timeouts = {
      invoke = "30m"
    }
  }
}

# Output examples
output "deprovision_summary" {
  value = {
    managed_count   = length(action.microsoft365_graph_beta_device_management_managed_device_deprovision.deprovision_multiple.config.managed_devices)
    comanaged_count = length(action.microsoft365_graph_beta_device_management_managed_device_deprovision.deprovision_comanaged.config.comanaged_devices)
  }
  description = "Count of devices deprovisioned"
}
