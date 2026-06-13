# Example 4: Look up a user by employee_id

data "microsoft365_graph_beta_users_user" "by_employee_id" {
  employee_id = "100200"
}

output "user_id" {
  value = data.microsoft365_graph_beta_users_user.by_employee_id.items[0].id
}

output "display_name" {
  value = data.microsoft365_graph_beta_users_user.by_employee_id.items[0].display_name
}
