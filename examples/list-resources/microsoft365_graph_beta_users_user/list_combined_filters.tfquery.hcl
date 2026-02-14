# List enabled member users
list "microsoft365_graph_beta_users_user" "enabled_members" {
  provider = microsoft365
  config {
    account_enabled_filter = true
    user_type_filter       = "Member"
  }
}
