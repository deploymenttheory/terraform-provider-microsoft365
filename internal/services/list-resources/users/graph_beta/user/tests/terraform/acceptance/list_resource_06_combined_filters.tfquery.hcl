provider "microsoft365" {}

# List users with combined filters
list "microsoft365_graph_beta_users_user" "combined" {
  provider = microsoft365
  config {
    account_enabled_filter = true
    user_type_filter       = "Member"
  }
}
