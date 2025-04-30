# Example 1: Lookup by ID
data "microsoft365_graph_beta_device_and_app_management_linux_platform_script" "by_id" {
  id = "0d6fee0f-d78d-4b00-87e8-cab65f31bb97"
  timeouts = {
    read = "30s"
  }
}

# Example 2: Lookup by Name
data "microsoft365_graph_beta_device_and_app_management_linux_platform_script" "by_name" {
  name = "Example Linux Script"
}

# Output all available fields from the data source
output "linux_script_by_id" {
  value = {
    id                  = data.microsoft365_graph_beta_device_and_app_management_linux_platform_script.by_id.id
    name                = data.microsoft365_graph_beta_device_and_app_management_linux_platform_script.by_id.name
    description         = data.microsoft365_graph_beta_device_and_app_management_linux_platform_script.by_id.description
    platforms           = data.microsoft365_graph_beta_device_and_app_management_linux_platform_script.by_id.platforms
    technologies        = data.microsoft365_graph_beta_device_and_app_management_linux_platform_script.by_id.technologies
    role_scope_tag_ids  = data.microsoft365_graph_beta_device_and_app_management_linux_platform_script.by_id.role_scope_tag_ids
    execution_context   = data.microsoft365_graph_beta_device_and_app_management_linux_platform_script.by_id.execution_context
    execution_frequency = data.microsoft365_graph_beta_device_and_app_management_linux_platform_script.by_id.execution_frequency
    execution_retries   = data.microsoft365_graph_beta_device_and_app_management_linux_platform_script.by_id.execution_retries
  }
}

output "linux_script_by_name" {
  value = {
    id                  = data.microsoft365_graph_beta_device_and_app_management_linux_platform_script.by_name.id
    name                = data.microsoft365_graph_beta_device_and_app_management_linux_platform_script.by_name.name
    description         = data.microsoft365_graph_beta_device_and_app_management_linux_platform_script.by_name.description
    platforms           = data.microsoft365_graph_beta_device_and_app_management_linux_platform_script.by_name.platforms
    technologies        = data.microsoft365_graph_beta_device_and_app_management_linux_platform_script.by_name.technologies
    role_scope_tag_ids  = data.microsoft365_graph_beta_device_and_app_management_linux_platform_script.by_name.role_scope_tag_ids
    execution_context   = data.microsoft365_graph_beta_device_and_app_management_linux_platform_script.by_name.execution_context
    execution_frequency = data.microsoft365_graph_beta_device_and_app_management_linux_platform_script.by_name.execution_frequency
    execution_retries   = data.microsoft365_graph_beta_device_and_app_management_linux_platform_script.by_name.execution_retries
  }
}
