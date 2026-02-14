provider "microsoft365" {}

# List users filtered by user principal name
list "microsoft365_graph_beta_users_user" "by_upn" {
  provider = microsoft365
  config {
    user_principal_name_filter = "admin"
  }
}
