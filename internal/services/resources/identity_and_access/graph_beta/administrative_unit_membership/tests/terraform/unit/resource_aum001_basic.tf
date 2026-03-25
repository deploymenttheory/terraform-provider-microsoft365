# AUM001: Basic membership with two users
resource "microsoft365_graph_beta_identity_and_access_administrative_unit_membership" "aum001_basic" {
  administrative_unit_id = "11111111-1111-1111-1111-111111111111"
  members = [
    "22222222-2222-2222-2222-222222222222",
    "33333333-3333-3333-3333-333333333333"
  ]
}
