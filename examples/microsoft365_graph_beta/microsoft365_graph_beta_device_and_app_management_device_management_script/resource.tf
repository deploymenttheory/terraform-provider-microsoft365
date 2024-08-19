# Example: Device Management Script Resource

# Data source for Azure AD group (assuming you have this data source available)
data "azuread_group" "example_group" {
  display_name = "Example Group"
}

resource "microsoft365_graph_beta_device_and_app_management_device_management_script" "example" {
  display_name = "Example Device Management Script"
  description  = "This is an example device management script"

  # The actual script content (PowerShell in this example)
  script_content = <<EOT
    # Your PowerShell script here
    Write-Host "Hello from device management script!"
    # Add more PowerShell commands as needed
  EOT

  run_as_account           = "system"  # Can be "system" or "user"
  enforce_signature_check  = false
  file_name                = "example_script.ps1"
  run_as_32_bit            = false

  role_scope_tag_ids = ["tag1", "tag2"]

  # Example assignment
  assignments {
    target {
      target_type                                     = "user"
      device_and_app_management_assignment_filter_id   = "filter-id-123"
      device_and_app_management_assignment_filter_type = "include"
      entra_object_id                                  = "user-object-id-456"
    }
  }

  # Example group assignment
  group_assignments {
    target_group_id = data.azuread_group.example_group.object_id
  }

  # Optionally specify timeouts
  timeouts {
    create = "30m"
    read   = "10m"
    update = "30m"
    delete = "30m"
  }
}

# Outputs
output "new_script_id" {
  value       = microsoft365_graph_beta_device_and_app_management_device_management_script.example.id
  description = "ID of the newly created Device Management Script"
}

output "existing_script_display_name" {
  value       = data.microsoft365_graph_beta_device_and_app_management_device_management_script.existing_script.display_name
  description = "Display name of the existing Device Management Script"
}

output "existing_script_last_modified" {
  value       = data.microsoft365_graph_beta_device_and_app_management_device_management_script.existing_script.last_modified_date_time
  description = "Last modified date of the existing Device Management Script"
}

output "existing_script_assignments" {
  value       = data.microsoft365_graph_beta_device_and_app_management_device_management_script.existing_script.assignments
  description = "Assignments of the existing Device Management Script"
  sensitive   = true
}

output "existing_script_group_assignments" {
  value       = data.microsoft365_graph_beta_device_and_app_management_device_management_script.existing_script.group_assignments
  description = "Group assignments of the existing Device Management Script"
  sensitive   = true
}