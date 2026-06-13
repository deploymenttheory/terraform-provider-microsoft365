# Example 6: Look up a user by user_principal_name (UPN)

data "microsoft365_graph_beta_users_user" "by_upn" {
  user_principal_name = "user@contoso.com"
}

output "user_id" {
  value = data.microsoft365_graph_beta_users_user.by_upn.items[0].id
}

output "display_name" {
  value = data.microsoft365_graph_beta_users_user.by_upn.items[0].display_name
}

output "mail" {
  value = data.microsoft365_graph_beta_users_user.by_upn.items[0].mail
}
