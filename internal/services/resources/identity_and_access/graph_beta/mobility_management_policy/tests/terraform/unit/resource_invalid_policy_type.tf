# Invalid: policy_type must be device or app.
resource "microsoft365_graph_beta_identity_and_access_mobility_management_policy" "invalid" {
  policy_type = "windows"
  policy_id   = "0000000a-0000-0000-c000-000000000000"
}
