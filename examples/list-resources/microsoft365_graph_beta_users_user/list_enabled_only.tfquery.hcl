# List only enabled user accounts
list "microsoft365_graph_beta_users_user" "enabled_only" {
  provider = microsoft365
  config {
    account_enabled_filter = true
  }
}
