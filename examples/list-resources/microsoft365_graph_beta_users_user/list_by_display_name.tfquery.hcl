# List users with display name starting with "John"
list "microsoft365_graph_beta_users_user" "by_display_name" {
  provider = microsoft365
  config {
    display_name_filter = "John"
  }
}
