# Example 1: Get all role definitions
data "microsoft365_graph_beta_identity_and_access_role_definitions" "all" {
  filter_type = "all"
}

# Example 2: Get a specific role definition by ID (Global Administrator)
data "microsoft365_graph_beta_identity_and_access_role_definitions" "by_id" {
  filter_type  = "id"
  filter_value = "62e90394-69f5-4237-9190-012177145e10"
}

# Example 3: Get role definitions by display name (partial match)
data "microsoft365_graph_beta_identity_and_access_role_definitions" "by_display_name" {
  filter_type  = "display_name"
  filter_value = "Security Administrator"
}

# Example 4: Get role definitions using OData filter (privileged roles only)
data "microsoft365_graph_beta_identity_and_access_role_definitions" "odata_filter" {
  filter_type  = "odata"
  odata_filter = "isPrivileged eq true"
}

# Example 5: Advanced OData query with filter, orderby, and select
data "microsoft365_graph_beta_identity_and_access_role_definitions" "odata_advanced" {
  filter_type   = "odata"
  odata_filter  = "isBuiltIn eq true"
  odata_orderby = "displayName"
  odata_select  = "id,displayName,description,isPrivileged"
}

# Example 6: Comprehensive OData query with filter, count, and orderby
data "microsoft365_graph_beta_identity_and_access_role_definitions" "odata_comprehensive" {
  filter_type   = "odata"
  odata_filter  = "isBuiltIn eq true"
  odata_count   = true
  odata_orderby = "displayName"
}

# Output examples
output "all_role_definitions_count" {
  value       = length(data.microsoft365_graph_beta_identity_and_access_role_definitions.all.items)
  description = "Total number of role definitions"
}

output "global_admin_role_id" {
  value       = data.microsoft365_graph_beta_identity_and_access_role_definitions.by_id.items[0].id
  description = "ID of the Global Administrator role"
}

output "privileged_roles" {
  value = [
    for role in data.microsoft365_graph_beta_identity_and_access_role_definitions.odata_filter.items :
    {
      id           = role.id
      display_name = role.display_name
      is_built_in  = role.is_built_in
    }
  ]
  description = "List of privileged roles"
}
