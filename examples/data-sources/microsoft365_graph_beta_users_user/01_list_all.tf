# Example 1: List all users in the tenant
# Results are always returned in the `items` list.

data "microsoft365_graph_beta_users_user" "all" {
  list_all = true
}

# Total number of users returned
output "all_users_count" {
  description = "The total number of users in the tenant"
  value       = length(data.microsoft365_graph_beta_users_user.all.items)
}

# Basic projection of every user
output "all_users_basic_info" {
  description = "Basic information about all users"
  value = [
    for user in data.microsoft365_graph_beta_users_user.all.items : {
      id                  = user.id
      display_name        = user.display_name
      user_principal_name = user.user_principal_name
      mail                = user.mail
    }
  ]
}
