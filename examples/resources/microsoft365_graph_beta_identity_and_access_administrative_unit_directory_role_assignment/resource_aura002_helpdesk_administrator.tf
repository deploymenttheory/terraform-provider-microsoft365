# AURA002: Helpdesk Administrator scoped to an Administrative Unit
# Assigns the Helpdesk Administrator role to a user, scoped to a specific administrative unit.
# The role holder can reset passwords and manage service requests only for users within that unit.

resource "microsoft365_graph_beta_identity_and_access_administrative_unit_directory_role_assignment" "helpdesk_admin" {
  administrative_unit_id = microsoft365_graph_beta_identity_and_access_administrative_unit.it_department.id
  # Helpdesk Administrator role ID (well-known in Microsoft Entra ID)
  directory_role_id = "729827e3-9c14-49f7-bb1b-9608f156bbb8"
  role_member_id    = microsoft365_graph_beta_users_user.it_support_agent.id
}
