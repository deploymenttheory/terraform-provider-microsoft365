# Invalid: included_groups is only valid when applies_to is selected.
resource "microsoft365_graph_beta_identity_and_access_mobility_management_policy" "invalid" {
  policy_type = "device"
  policy_id   = "0000000a-0000-0000-c000-000000000000"

  applies_to      = "all"
  included_groups = ["11111111-1111-1111-1111-111111111111"]
}
