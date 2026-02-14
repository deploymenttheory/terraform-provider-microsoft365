provider "microsoft365" {}

# List users filtered by account enabled status
list "microsoft365_graph_beta_users_user" "by_account_enabled" {
  provider = microsoft365
  config {
    account_enabled_filter = true
  }
}
