# Test 01: Deploy a user, wait for propagation, then list all users and attest the list is populated.

resource "random_string" "test" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_users_user" "test" {
  display_name        = "acc-test-user-list-all-${random_string.test.result}"
  user_principal_name = "acc-test-user-list-all-${random_string.test.result}@deploymenttheory.com"
  mail_nickname       = "acctestuserlistall${random_string.test.result}"
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
  list_all   = true
  depends_on = [time_sleep.wait_for_user]
}
