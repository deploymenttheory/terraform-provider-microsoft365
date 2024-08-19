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