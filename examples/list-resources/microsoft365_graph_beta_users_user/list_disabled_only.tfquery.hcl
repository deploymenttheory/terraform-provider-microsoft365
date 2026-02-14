# List only disabled user accounts
list "microsoft365_graph_beta_users_user" "disabled_only" {
  provider = microsoft365
  config {
    account_enabled_filter = false
  }
}
