# AURA002: Scoped role assignment with a different member and role
resource "microsoft365_graph_beta_identity_and_access_administrative_unit_role_assignment" "aura002_different_member" {
  administrative_unit_id = "11111111-1111-1111-1111-111111111111"
  role_id                = "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"
  role_member_id         = "33333333-3333-3333-3333-333333333333"
}
