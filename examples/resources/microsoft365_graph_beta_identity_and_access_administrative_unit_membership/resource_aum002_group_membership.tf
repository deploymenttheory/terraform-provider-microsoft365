# AUM002: Group Membership
# Adds a security group as a member of an administrative unit
resource "microsoft365_graph_beta_identity_and_access_administrative_unit_membership" "aum002_group" {
  administrative_unit_id = microsoft365_graph_beta_identity_and_access_administrative_unit.it_department.id
  members = [
    microsoft365_graph_beta_groups_group.it_security_group.id
  ]
}
