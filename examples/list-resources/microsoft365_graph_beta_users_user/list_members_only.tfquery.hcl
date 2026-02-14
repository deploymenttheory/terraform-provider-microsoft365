# List only member users (excluding guests)
list "microsoft365_graph_beta_users_user" "members_only" {
  provider = microsoft365
  config {
    user_type_filter = "Member"
  }
}
