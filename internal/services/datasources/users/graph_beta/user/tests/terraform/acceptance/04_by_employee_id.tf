# Test 04: Deploy a user with an employee_id, wait for propagation, then look it up by employee_id.

resource "random_string" "test" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_users_user" "test" {
  display_name        = "acc-test-user-by-empid-${random_string.test.result}"
  user_principal_name = "acc-test-user-by-empid-${random_string.test.result}@deploymenttheory.com"
  mail_nickname       = "acctestuserbyempid${random_string.test.result}"
  account_enabled     = true
  employee_id         = "ACC-${random_string.test.result}"
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
  employee_id = microsoft365_graph_beta_users_user.test.employee_id
  depends_on  = [time_sleep.wait_for_user]
}
