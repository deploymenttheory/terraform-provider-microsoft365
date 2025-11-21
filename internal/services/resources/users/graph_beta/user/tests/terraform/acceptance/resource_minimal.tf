resource "random_string" "minimal_user_id" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_users_user" "minimal" {
  display_name        = "acc-test-user-minimal-${random_string.minimal_user_id.result}"
  user_principal_name = "acc-test-user-minimal-${random_string.minimal_user_id.result}@deploymenttheory.com"
  mail_nickname       = "acc-test-user-minimal-${random_string.minimal_user_id.result}"
  account_enabled     = true
  password_profile = {
    password                           = "SecureP@ssw0rd123!"
    force_change_password_next_sign_in = false
  }
}

