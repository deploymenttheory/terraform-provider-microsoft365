# Complex query combining multiple conditions
list "microsoft365_graph_beta_users_user" "complex" {
  provider = microsoft365
  config {
    odata_filter = "(userType eq 'Member' and accountEnabled eq true) and (startsWith(userPrincipalName, 'admin') or startsWith(userPrincipalName, 'svc'))"
  }
}
