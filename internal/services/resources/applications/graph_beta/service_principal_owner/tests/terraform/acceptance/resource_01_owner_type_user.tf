# Acceptance test: Service Principal Owner with User owner type
# Full dependency chain: random_string -> application + user -> service_principal -> service_principal_owner

resource "random_string" "test_id_user" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_applications_application" "test_app_user" {
  display_name = "acc-test-sp-owner-user-${random_string.test_id_user.result}"
  description  = "Application for service principal user owner assignment acceptance test"
  hard_delete  = true
}

resource "microsoft365_graph_beta_applications_service_principal" "test_sp_user" {
  app_id      = microsoft365_graph_beta_applications_application.test_app_user.app_id
  hard_delete = true
}

resource "microsoft365_graph_beta_users_user" "test_owner_user" {
  display_name        = "acc-test-sp-owner-user-${random_string.test_id_user.result}"
  user_principal_name = "acc-test-sp-owner-user-${random_string.test_id_user.result}@deploymenttheory.com"
  mail_nickname       = "acc-test-sp-owner-user-${random_string.test_id_user.result}"
  account_enabled     = true
  password_profile = {
    password                           = "SecureP@ssw0rd123!"
    force_change_password_next_sign_in = false
  }
}

resource "time_sleep" "wait_for_sp_user" {
  depends_on      = [microsoft365_graph_beta_applications_service_principal.test_sp_user]
  create_duration = "15s"
}

resource "microsoft365_graph_beta_applications_service_principal_owner" "test_user" {
  service_principal_id = microsoft365_graph_beta_applications_service_principal.test_sp_user.id
  owner_id             = microsoft365_graph_beta_users_user.test_owner_user.id
  owner_object_type    = "User"

  depends_on = [time_sleep.wait_for_sp_user]
}
