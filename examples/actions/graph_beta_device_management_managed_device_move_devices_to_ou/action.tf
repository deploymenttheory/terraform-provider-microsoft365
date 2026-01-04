# Example 1: Move devices to organizational unit - Minimal
action "microsoft365_graph_beta_device_management_managed_device_move_devices_to_ou" "move_single" {
  config {
    organizational_unit_path = "OU=Workstations,DC=contoso,DC=com"
    managed_device_ids       = ["12345678-1234-1234-1234-123456789abc"]
  }
}

# Example 2: Move multiple devices to OU
action "microsoft365_graph_beta_device_management_managed_device_move_devices_to_ou" "move_multiple" {
  config {
    organizational_unit_path = "OU=Finance,OU=Departments,DC=contoso,DC=com"
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

# Example 3: Move devices with validation - Maximal
action "microsoft365_graph_beta_device_management_managed_device_move_devices_to_ou" "move_maximal" {
  config {
    organizational_unit_path = "OU=IT,OU=Departments,DC=contoso,DC=com"
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

# Example 4: Move department devices to new OU
data "microsoft365_graph_beta_device_management_managed_device" "marketing_devices" {
  filter_type  = "odata"
  odata_filter = "deviceCategoryDisplayName eq 'Marketing'"
}

action "microsoft365_graph_beta_device_management_managed_device_move_devices_to_ou" "move_marketing" {
  config {
    organizational_unit_path = "OU=Marketing,OU=Departments,DC=contoso,DC=com"
    managed_device_ids       = [for device in data.microsoft365_graph_beta_device_management_managed_device.marketing_devices.items : device.id]

    validate_device_exists = true

    timeouts = {
      invoke = "15m"
    }
  }
}
