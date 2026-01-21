# Example 5: Look up group using custom OData query
# Use this for advanced filtering when standard attributes don't meet your needs

data "microsoft365_graph_beta_groups_group" "by_odata_query" {
  odata_query = "displayName eq 'My Group' and securityEnabled eq true"
}

output "group_id" {
  value = data.microsoft365_graph_beta_groups_group.by_odata_query.id
}

output "display_name" {
  value = data.microsoft365_graph_beta_groups_group.by_odata_query.display_name
}

# Example: More complex OData query
data "microsoft365_graph_beta_groups_group" "dynamic_group" {
  odata_query = "startswith(displayName, 'DYN-') and securityEnabled eq true and mailEnabled eq false"
}

output "dynamic_group_id" {
  value = data.microsoft365_graph_beta_groups_group.dynamic_group.id
}
