# Basic lookup examples
# Look up by display name
data "microsoft365_graph_beta_device_and_app_management_macos_platform_script" "by_name" {
  display_name = "MacOS Shell Script"
}

# Look up by ID
data "microsoft365_graph_beta_device_and_app_management_macos_platform_script" "by_id" {
  id = "00000000-0000-0000-0000-000000000001"
}

# Example: Using the data source outputs
output "script_details" {
  value = {
    id                 = data.microsoft365_graph_beta_device_and_app_management_macos_platform_script.example.id
    display_name       = data.microsoft365_graph_beta_device_and_app_management_macos_platform_script.example.display_name
    description        = data.microsoft365_graph_beta_device_and_app_management_macos_platform_script.example.description
    created_date_time  = data.microsoft365_graph_beta_device_and_app_management_macos_platform_script.example.created_date_time
    last_modified_time = data.microsoft365_graph_beta_device_and_app_management_macos_platform_script.example.last_modified_date_time
    run_as_account     = data.microsoft365_graph_beta_device_and_app_management_macos_platform_script.example.run_as_account
    file_name          = data.microsoft365_graph_beta_device_and_app_management_macos_platform_script.example.file_name
  }
  description = "Details of the retrieved macOS platform script"
}

# Example: Accessing script assignments
output "script_assignments" {
  value       = data.microsoft365_graph_beta_device_and_app_management_macos_platform_script.example.assignments
  description = "Assignment configuration for the macOS platform script"
}