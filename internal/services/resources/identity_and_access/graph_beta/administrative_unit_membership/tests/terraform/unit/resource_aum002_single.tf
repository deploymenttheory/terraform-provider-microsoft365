# AUM002: Single member
resource "microsoft365_graph_beta_identity_and_access_administrative_unit_membership" "aum002_single" {
  administrative_unit_id = "22222222-1111-1111-1111-111111111111"
  members = [
    "44444444-4444-4444-4444-444444444444"
  ]
}
