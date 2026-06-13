# Example 9: Look up users using a custom OData query
# Use this for advanced filtering when the standard lookup attributes don't fit.

data "microsoft365_graph_beta_users_user" "by_odata_query" {
  odata_query = "accountEnabled eq true and userType eq 'Member'"
}

output "enabled_members" {
  description = "All enabled member users matching the query"
  value = [
    for user in data.microsoft365_graph_beta_users_user.by_odata_query.items : {
      id                  = user.id
      display_name        = user.display_name
      user_principal_name = user.user_principal_name
    }
  ]
}

# Example: More complex OData query using a function
data "microsoft365_graph_beta_users_user" "by_odata_startswith" {
  odata_query = "startswith(displayName, 'A') and accountEnabled eq true"
}

output "users_starting_with_a_count" {
  value = length(data.microsoft365_graph_beta_users_user.by_odata_startswith.items)
}
