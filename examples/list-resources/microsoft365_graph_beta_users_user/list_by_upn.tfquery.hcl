# List users with UPN starting with "admin"
list "microsoft365_graph_beta_users_user" "by_upn" {
  provider = microsoft365
  config {
    user_principal_name_filter = "admin"
  }
}
