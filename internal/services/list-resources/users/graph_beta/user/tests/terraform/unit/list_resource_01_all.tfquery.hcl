provider "microsoft365" {}

# List all users
list "microsoft365_graph_beta_users_user" "all" {
  provider = microsoft365
  config {}
}
