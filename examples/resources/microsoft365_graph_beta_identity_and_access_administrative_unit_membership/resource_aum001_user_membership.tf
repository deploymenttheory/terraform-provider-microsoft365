# AUM001: User Membership
# Adds multiple users as members of an administrative unit
resource "microsoft365_graph_beta_identity_and_access_administrative_unit_membership" "aum001_users" {
  administrative_unit_id = microsoft365_graph_beta_identity_and_access_administrative_unit.finance.id
  members = [
    microsoft365_graph_beta_users_user.finance_user1.id,
    microsoft365_graph_beta_users_user.finance_user2.id,
    microsoft365_graph_beta_users_user.finance_user3.id
  ]
}
