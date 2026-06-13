# Test 05: Deploy a user with a unique given_name, wait for propagation, then look it up by given_name.

resource "random_string" "test" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_users_user" "test" {
  display_name        = "acc-test-user-by-given-${random_string.test.result}"
  user_principal_name = "acc-test-user-by-given-${random_string.test.result}@deploymenttheory.com"
  mail_nickname       = "acctestuserbygiven${random_string.test.result}"
  account_enabled     = true
  given_name          = "AccTest${random_string.test.result}"
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
  given_name = microsoft365_graph_beta_users_user.test.given_name
  depends_on = [time_sleep.wait_for_user]
}
