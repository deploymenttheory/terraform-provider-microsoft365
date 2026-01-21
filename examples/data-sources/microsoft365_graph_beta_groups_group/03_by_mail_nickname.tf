# Example 3: Look up group by mail_nickname

data "microsoft365_graph_beta_groups_group" "by_mail_nickname" {
  mail_nickname = "mygroup"
}

output "group_id" {
  value = data.microsoft365_graph_beta_groups_group.by_mail_nickname.id
}

output "display_name" {
  value = data.microsoft365_graph_beta_groups_group.by_mail_nickname.display_name
}
