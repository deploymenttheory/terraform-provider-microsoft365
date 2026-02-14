# List only guest users
list "microsoft365_graph_beta_users_user" "guests_only" {
  provider = microsoft365
  config {
    user_type_filter = "Guest"
  }
}
