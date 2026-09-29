resource "microsoft365_graph_beta_identity_and_access_b2b_management_policy" "invalid" {
  display_name = "unit-test-b2b-management-policy-invalid"
  definition   = ["not json"]
}
