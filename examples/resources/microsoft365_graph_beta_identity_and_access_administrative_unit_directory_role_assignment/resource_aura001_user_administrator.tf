# AURA001: User Administrator scoped to an Administrative Unit
# Assigns the User Administrator role to a user, scoped to a specific administrative unit.
# The role holder can manage users only within that administrative unit.

resource "microsoft365_graph_beta_identity_and_access_administrative_unit_directory_role_assignment" "user_admin" {
  administrative_unit_id = microsoft365_graph_beta_identity_and_access_administrative_unit.finance.id
  # User Administrator role ID (well-known in Microsoft Entra ID)
  directory_role_id = "fe930be7-5e62-47db-91af-98c3a49a38b1"
  role_member_id    = microsoft365_graph_beta_users_user.helpdesk_lead.id
}
