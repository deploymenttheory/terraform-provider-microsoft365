# Minimal example with only required properties
resource "microsoft365_graph_beta_users_user" "minimal_example" {
  display_name        = "John Doe"
  account_enabled     = true
  user_principal_name = "john.doe@contoso.com"
  mail_nickname       = "johndoe"
  hard_delete         = true
  password_profile = {
    password                           = "SecurePassword123!"
    force_change_password_next_sign_in = true
  }
}