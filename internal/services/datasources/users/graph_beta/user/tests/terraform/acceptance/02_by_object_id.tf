# Test 02: Deploy a user, wait for propagation, then look it up by object_id.

resource "random_string" "test" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_users_user" "test" {
  display_name        = "acc-test-user-by-id-${random_string.test.result}"
  user_principal_name = "acc-test-user-by-id-${random_string.test.result}@deploymenttheory.com"
  mail_nickname       = "acctestuserbyid${random_string.test.result}"
  account_enabled     = true
  hard_delete         = true
  password_profile = {
    password                           = "SecureP@ssw0rd123!"
    force_change_password_next_sign_in = false
  }
}

resource "time_sleep" "wait_for_user" {
  depends_on      = [microsoft365_graph_beta_users_user.test]
  create_duration = "30s"
}

data "microsoft365_graph_beta_users_user" "test" {
  object_id  = microsoft365_graph_beta_users_user.test.id
  depends_on = [time_sleep.wait_for_user]
}
