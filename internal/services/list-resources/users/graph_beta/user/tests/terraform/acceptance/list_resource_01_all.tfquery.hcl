provider "microsoft365" {}

# List all users after deployment
list "microsoft365_graph_beta_users_user" "all" {
  provider = microsoft365
  config {}
}
