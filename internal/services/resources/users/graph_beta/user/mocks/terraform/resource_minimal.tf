resource "microsoft365_graph_beta_users_user" "minimal" {
  display_name        = "Minimal User"
  user_principal_name = "minimal.user@contoso.com"
  account_enabled     = true
  password_profile    = {
    password = "SecureP@ssw0rd123!"
    force_change_password_next_sign_in = false
  }
} 