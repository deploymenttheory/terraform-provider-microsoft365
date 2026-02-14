# Find enabled users with specific job title
list "microsoft365_graph_beta_users_user" "enabled_with_title" {
  provider = microsoft365
  config {
    odata_filter = "accountEnabled eq true and jobTitle eq 'Manager'"
  }
}
