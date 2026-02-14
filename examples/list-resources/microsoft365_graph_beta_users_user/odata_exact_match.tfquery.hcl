# Find specific user by exact UPN match
list "microsoft365_graph_beta_users_user" "exact_upn" {
  provider = microsoft365
  config {
    odata_filter = "userPrincipalName eq 'admin@contoso.com'"
  }
}
