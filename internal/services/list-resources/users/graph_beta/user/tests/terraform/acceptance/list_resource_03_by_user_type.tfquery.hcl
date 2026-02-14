provider "microsoft365" {}

# List users filtered by user type
list "microsoft365_graph_beta_users_user" "by_user_type" {
  provider = microsoft365
  config {
    user_type_filter = "Member"
  }
}
