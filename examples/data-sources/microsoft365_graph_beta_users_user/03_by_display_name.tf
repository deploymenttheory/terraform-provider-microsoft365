# Example 3: Look up a user by display_name

data "microsoft365_graph_beta_users_user" "by_display_name" {
  display_name = "John Doe"
}

output "user_id" {
  value = data.microsoft365_graph_beta_users_user.by_display_name.items[0].id
}

output "user_principal_name" {
  value = data.microsoft365_graph_beta_users_user.by_display_name.items[0].user_principal_name
}
