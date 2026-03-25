# AUM003: Mixed Membership
# Adds both users and groups as members of an administrative unit
resource "microsoft365_graph_beta_identity_and_access_administrative_unit_membership" "aum003_mixed" {
  administrative_unit_id = microsoft365_graph_beta_identity_and_access_administrative_unit.executive.id
  members = [
    microsoft365_graph_beta_users_user.ceo.id,
    microsoft365_graph_beta_users_user.cfo.id,
    microsoft365_graph_beta_groups_group.executive_team.id
  ]
}
