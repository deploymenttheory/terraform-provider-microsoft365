resource "microsoft365_graph_beta_users_agent_user" "invalid" {
  display_name    = "Invalid User"
  account_enabled = true
  # Missing required user_principal_name and mail_nickname
}
