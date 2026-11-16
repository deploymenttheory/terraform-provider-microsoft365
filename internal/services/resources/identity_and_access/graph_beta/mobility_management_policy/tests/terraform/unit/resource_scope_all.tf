# Scope all: every user is in scope and no groups are included.
resource "microsoft365_graph_beta_identity_and_access_mobility_management_policy" "maximal" {
  policy_type = "device"
  policy_id   = "0000000a-0000-0000-c000-000000000000"

  applies_to = "all"
}
