# Example 2: Look up group by display_name

data "microsoft365_graph_beta_groups_group" "by_display_name" {
  display_name = "My Group Name"
}

output "group_id" {
  value = data.microsoft365_graph_beta_groups_group.by_display_name.id
}

output "object_id" {
  value = data.microsoft365_graph_beta_groups_group.by_display_name.object_id
}
