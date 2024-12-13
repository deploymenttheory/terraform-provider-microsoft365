# Output: Display information from the data source
output "existing_role_scope_tag_id" {
  value = data.microsoft365_graph_beta_device_and_app_management_role_scope_tag.existing_role_scope_tag.id
}

output "existing_role_scope_tag_description" {
  value = data.microsoft365_graph_beta_device_and_app_management_role_scope_tag.existing_role_scope_tag.description
}

output "existing_role_scope_tag_platform" {
  value = data.microsoft365_graph_beta_device_and_app_management_role_scope_tag.existing_role_scope_tag.platform
}

output "existing_role_scope_tag_rule" {
  value = data.microsoft365_graph_beta_device_and_app_management_role_scope_tag.existing_role_scope_tag.rule
}