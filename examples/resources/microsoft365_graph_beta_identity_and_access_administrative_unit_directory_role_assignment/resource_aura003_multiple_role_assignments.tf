# AURA003: Multiple role assignments scoped to the same Administrative Unit
# Assigns both User Administrator and Helpdesk Administrator roles to different users
# within the same administrative unit, enabling a tiered delegation model.

resource "microsoft365_graph_beta_identity_and_access_administrative_unit_directory_role_assignment" "au_user_admin" {
  administrative_unit_id = microsoft365_graph_beta_identity_and_access_administrative_unit.regional_office.id
  # User Administrator role ID (well-known in Microsoft Entra ID)
  directory_role_id = "fe930be7-5e62-47db-91af-98c3a49a38b1"
  role_member_id    = microsoft365_graph_beta_users_user.regional_it_manager.id
}

resource "microsoft365_graph_beta_identity_and_access_administrative_unit_directory_role_assignment" "au_helpdesk_admin" {
  administrative_unit_id = microsoft365_graph_beta_identity_and_access_administrative_unit.regional_office.id
  # Helpdesk Administrator role ID (well-known in Microsoft Entra ID)
  directory_role_id = "729827e3-9c14-49f7-bb1b-9608f156bbb8"
  role_member_id    = microsoft365_graph_beta_users_user.regional_helpdesk.id
}
