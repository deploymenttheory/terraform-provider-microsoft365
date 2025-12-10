resource "microsoft365_graph_beta_users_user" "minimal" {
  display_name        = "Minimal User"
  user_principal_name = "minimal.user@deploymenttheory.com"
  mail_nickname       = "minimal.user"
  account_enabled     = true
  hard_delete         = true
  password_profile = {
    password                           = "SecureP@ssw0rd123!"
    force_change_password_next_sign_in = false
  }
}

