# Example 5: Look up users by given_name (first name)
# A given name may match multiple users, all returned in `items`.

data "microsoft365_graph_beta_users_user" "by_given_name" {
  given_name = "John"
}

output "matched_users" {
  description = "All users matching the given name"
  value = [
    for user in data.microsoft365_graph_beta_users_user.by_given_name.items : {
      id                  = user.id
      display_name        = user.display_name
      user_principal_name = user.user_principal_name
    }
  ]
}
