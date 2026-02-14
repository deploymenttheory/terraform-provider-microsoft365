# Find users with display name starting with specific prefix
list "microsoft365_graph_beta_users_user" "name_prefix" {
  provider = microsoft365
  config {
    odata_filter = "startsWith(displayName, 'Adele')"
  }
}
