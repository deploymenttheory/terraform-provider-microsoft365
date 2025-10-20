# Example 1: Move single device to workstations OU
action "microsoft365_graph_beta_device_management_managed_device_move_devices_to_ou" "move_to_workstations" {
  organizational_unit_path = "OU=Workstations,OU=Computers,DC=contoso,DC=com"
  managed_device_ids       = ["12345678-1234-1234-1234-123456789abc"]

  timeouts = {
    invoke = "5m"
  }
}

# Example 2: Move multiple devices to same OU
action "microsoft365_graph_beta_device_management_managed_device_move_devices_to_ou" "move_to_marketing" {
  organizational_unit_path = "OU=Marketing,OU=Departments,DC=contoso,DC=com"
  managed_device_ids = [
    "12345678-1234-1234-1234-123456789abc",
    "87654321-4321-4321-4321-ba9876543210",
    "abcdef12-3456-7890-abcd-ef1234567890"
  ]

  timeouts = {
    invoke = "10m"
  }
}

# Example 3: Move devices by department
variable "marketing_devices" {
  description = "Device IDs for marketing department"
  type        = list(string)
  default = [
    "aaaa1111-1111-1111-1111-111111111111",
    "bbbb2222-2222-2222-2222-222222222222"
  ]
}

action "microsoft365_graph_beta_device_management_managed_device_move_devices_to_ou" "marketing_org_unit" {
  organizational_unit_path = "OU=Marketing,OU=Departments,DC=corp,DC=example,DC=com"
  managed_device_ids       = var.marketing_devices

  timeouts = {
    invoke = "15m"
  }
}

# Example 4: Move devices based on data source filter
data "microsoft365_graph_beta_device_management_managed_device" "windows_devices" {
  filter_type  = "odata"
  odata_filter = "operatingSystem eq 'Windows' and deviceCategoryDisplayName eq 'Relocate'"
}

action "microsoft365_graph_beta_device_management_managed_device_move_devices_to_ou" "relocate_devices" {
  organizational_unit_path = "OU=NewLocation,OU=Offices,DC=contoso,DC=com"
  managed_device_ids       = [for device in data.microsoft365_graph_beta_device_management_managed_device.windows_devices.items : device.id]

  timeouts = {
    invoke = "20m"
  }
}

# Example 5: Move laptops to mobile OU
locals {
  laptop_devices = [
    "11111111-1111-1111-1111-111111111111",
    "22222222-2222-2222-2222-222222222222",
    "33333333-3333-3333-3333-333333333333"
  ]
}

action "microsoft365_graph_beta_device_management_managed_device_move_devices_to_ou" "laptops_to_mobile_ou" {
  organizational_unit_path = "OU=Laptops,OU=Mobile,OU=Devices,DC=corp,DC=acme,DC=com"
  managed_device_ids       = local.laptop_devices

  timeouts = {
    invoke = "15m"
  }
}

# Example 6: Move co-managed devices
action "microsoft365_graph_beta_device_management_managed_device_move_devices_to_ou" "move_comanaged" {
  organizational_unit_path = "OU=CoManaged,OU=SCCM,DC=contoso,DC=com"
  comanaged_device_ids     = ["abcdef12-3456-7890-abcd-ef1234567890"]

  timeouts = {
    invoke = "5m"
  }
}

# Example 7: Move devices to apply different GPOs
action "microsoft365_graph_beta_device_management_managed_device_move_devices_to_ou" "secure_workstations_gpo" {
  organizational_unit_path = "OU=SecureWorkstations,OU=Security,DC=contoso,DC=com"
  managed_device_ids = [
    "secure01-1111-1111-1111-111111111111",
    "secure02-2222-2222-2222-222222222222"
  ]

  timeouts = {
    invoke = "10m"
  }
}

# Example 8: Move by location/office
locals {
  office_locations = {
    "seattle" = {
      ou_path = "OU=Seattle,OU=Offices,DC=corp,DC=contoso,DC=com"
      devices = [
        "sea001-1111-1111-1111-111111111111",
        "sea002-2222-2222-2222-222222222222"
      ]
    }
    "portland" = {
      ou_path = "OU=Portland,OU=Offices,DC=corp,DC=contoso,DC=com"
      devices = [
        "pdx001-3333-3333-3333-333333333333",
        "pdx002-4444-4444-4444-444444444444"
      ]
    }
  }
}

action "microsoft365_graph_beta_device_management_managed_device_move_devices_to_ou" "seattle_office" {
  organizational_unit_path = local.office_locations["seattle"].ou_path
  managed_device_ids       = local.office_locations["seattle"].devices

  timeouts = {
    invoke = "15m"
  }
}

action "microsoft365_graph_beta_device_management_managed_device_move_devices_to_ou" "portland_office" {
  organizational_unit_path = local.office_locations["portland"].ou_path
  managed_device_ids       = local.office_locations["portland"].devices

  timeouts = {
    invoke = "15m"
  }
}

# Example 9: Move devices for compliance/security
data "microsoft365_graph_beta_device_management_managed_device" "quarantine_devices" {
  filter_type  = "odata"
  odata_filter = "deviceCategoryDisplayName eq 'Quarantine'"
}

action "microsoft365_graph_beta_device_management_managed_device_move_devices_to_ou" "quarantine_ou" {
  organizational_unit_path = "OU=Quarantine,OU=Security,DC=corp,DC=contoso,DC=com"
  managed_device_ids       = [for device in data.microsoft365_graph_beta_device_management_managed_device.quarantine_devices.items : device.id]

  timeouts = {
    invoke = "10m"
  }
}

# Example 10: Multi-tier OU structure
action "microsoft365_graph_beta_device_management_managed_device_move_devices_to_ou" "finance_restricted" {
  organizational_unit_path = "OU=Restricted,OU=Finance,OU=Departments,DC=corp,DC=contoso,DC=com"
  managed_device_ids = [
    "fin001-1111-1111-1111-111111111111",
    "fin002-2222-2222-2222-222222222222"
  ]

  timeouts = {
    invoke = "15m"
  }
}

# Output examples
output "moved_devices_summary" {
  value = {
    marketing = {
      ou_path      = action.move_to_marketing.organizational_unit_path
      device_count = length(action.move_to_marketing.managed_device_ids)
    }
    laptops = {
      ou_path      = action.laptops_to_mobile_ou.organizational_unit_path
      device_count = length(action.laptops_to_mobile_ou.managed_device_ids)
    }
  }
  description = "Summary of devices moved to OUs"
}