resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

data "microsoft365_graph_beta_identity_and_access_directory_role" "user_admin" {
  display_name = "User Administrator"
}

resource "microsoft365_graph_beta_users_user" "aura001_user" {
  user_principal_name = "aura001-user-${random_string.suffix.result}@deploymenttheory.com"
  display_name        = "AURA001 Test User"
  mail_nickname       = "aura001-user-${random_string.suffix.result}"
  account_enabled     = true
  password_profile = {
    password                           = "P@ssw0rd!${random_string.suffix.result}"
    force_change_password_next_sign_in = false
  }
  hard_delete = true
}

resource "microsoft365_graph_beta_identity_and_access_administrative_unit" "aura001_au" {
  display_name = "acc-test-aura001-au-${random_string.suffix.result}"
  description  = "Administrative unit for scoped role assignment test AURA001"
  hard_delete  = true
}

resource "time_sleep" "aura001_wait" {
  create_duration = "30s"

  depends_on = [
    microsoft365_graph_beta_users_user.aura001_user,
    microsoft365_graph_beta_identity_and_access_administrative_unit.aura001_au,
  ]
}

resource "microsoft365_graph_beta_identity_and_access_administrative_unit_role_assignment" "aura001_basic" {
  administrative_unit_id = microsoft365_graph_beta_identity_and_access_administrative_unit.aura001_au.id
  role_id                = data.microsoft365_graph_beta_identity_and_access_directory_role.user_admin.items[0].id
  role_member_id         = microsoft365_graph_beta_users_user.aura001_user.id

  depends_on = [time_sleep.aura001_wait]
}
