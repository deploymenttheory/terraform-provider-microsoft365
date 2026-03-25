# AURA001: Basic scoped role assignment - assigns User Administrator role to a user within an AU
resource "microsoft365_graph_beta_identity_and_access_administrative_unit_role_assignment" "aura001_basic" {
  administrative_unit_id = "11111111-1111-1111-1111-111111111111"
  role_id                = "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"
  role_member_id         = "22222222-2222-2222-2222-222222222222"
}
