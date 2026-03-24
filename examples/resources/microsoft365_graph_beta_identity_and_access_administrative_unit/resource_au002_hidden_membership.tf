# AU002: Hidden Membership Administrative Unit
# Creates an administrative unit with hidden membership where only members
# can see other members of the unit
resource "microsoft365_graph_beta_identity_and_access_administrative_unit" "au002_hidden" {
  display_name = "Executive Team"
  description  = "Administrative unit for executive team with hidden membership"
  visibility   = "HiddenMembership"
}
