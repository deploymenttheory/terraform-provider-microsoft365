# AURA004: Full example — create AU, user, and scoped role assignment together
# Demonstrates creating all dependent resources inline and chaining them
# with depends_on for correct provisioning order.

resource "microsoft365_graph_beta_identity_and_access_administrative_unit" "emea_office" {
  display_name = "EMEA Office"
  description  = "Administrative unit for EMEA region users and devices"
  hard_delete  = true
}

resource "microsoft365_graph_beta_users_user" "emea_it_admin" {
  user_principal_name = "emea-it-admin@contoso.com"
  display_name        = "EMEA IT Administrator"
  mail_nickname       = "emea-it-admin"
  account_enabled     = true
  password_profile = {
    password                           = "ChangeMe123!"
    force_change_password_next_sign_in = true
  }
  hard_delete = true
}

resource "microsoft365_graph_beta_identity_and_access_administrative_unit_role_assignment" "emea_user_admin" {
  administrative_unit_id = microsoft365_graph_beta_identity_and_access_administrative_unit.emea_office.id
  # User Administrator role ID (well-known in Microsoft Entra ID)
  role_id        = "fe930be7-5e62-47db-91af-98c3a49a38b1"
  role_member_id = microsoft365_graph_beta_users_user.emea_it_admin.id

  depends_on = [
    microsoft365_graph_beta_identity_and_access_administrative_unit.emea_office,
    microsoft365_graph_beta_users_user.emea_it_admin,
  ]
}
