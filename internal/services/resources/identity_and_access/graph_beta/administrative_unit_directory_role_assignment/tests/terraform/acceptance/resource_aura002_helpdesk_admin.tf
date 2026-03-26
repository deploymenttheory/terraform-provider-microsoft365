resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

data "microsoft365_graph_beta_identity_and_access_directory_role" "helpdesk_admin" {
  display_name = "Helpdesk Administrator"
}

resource "microsoft365_graph_beta_users_user" "aura002_user" {
  user_principal_name = "aura002-user-${random_string.suffix.result}@deploymenttheory.com"
  display_name        = "AURA002 Test User"
  mail_nickname       = "aura002-user-${random_string.suffix.result}"
  account_enabled     = true
  password_profile = {
    password                           = "P@ssw0rd!${random_string.suffix.result}"
    force_change_password_next_sign_in = false
  }
  hard_delete = true
}

resource "microsoft365_graph_beta_identity_and_access_administrative_unit" "aura002_au" {
  display_name = "acc-test-aura002-au-${random_string.suffix.result}"
  description  = "Administrative unit for scoped role assignment test AURA002"
  hard_delete  = true
}

resource "time_sleep" "aura002_wait" {
  create_duration = "30s"

  depends_on = [
    microsoft365_graph_beta_users_user.aura002_user,
    microsoft365_graph_beta_identity_and_access_administrative_unit.aura002_au,
  ]
}

resource "microsoft365_graph_beta_identity_and_access_administrative_unit_directory_role_assignment" "aura002_helpdesk_admin" {
  administrative_unit_id = microsoft365_graph_beta_identity_and_access_administrative_unit.aura002_au.id
  directory_role_id      = data.microsoft365_graph_beta_identity_and_access_directory_role.helpdesk_admin.items[0].id
  role_member_id         = microsoft365_graph_beta_users_user.aura002_user.id

  depends_on = [time_sleep.aura002_wait]
}
