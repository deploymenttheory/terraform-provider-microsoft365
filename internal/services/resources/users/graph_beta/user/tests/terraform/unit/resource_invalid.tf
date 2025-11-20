resource "microsoft365_graph_beta_users_user" "invalid" {
  display_name    = "Invalid User"
  account_enabled = true
  # Missing required user_principal_name and password_profile
}
