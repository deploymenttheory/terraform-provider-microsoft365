# Example 4: Look up group by display_name with additional filters
# Use mail_enabled and security_enabled as additional filters to narrow results

data "microsoft365_graph_beta_groups_group" "security_group" {
  display_name     = "My Security Group"
  security_enabled = true
  mail_enabled     = false
}

output "group_id" {
  value = data.microsoft365_graph_beta_groups_group.security_group.id
}

output "group_types" {
  value = data.microsoft365_graph_beta_groups_group.security_group.group_types
}

output "members_count" {
  description = "Number of members in the group"
  value       = length(data.microsoft365_graph_beta_groups_group.security_group.members)
}
