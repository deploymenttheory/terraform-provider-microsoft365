resource "microsoft365_graph_beta_identity_and_access_b2b_management_policy_assignment" "target" {
  b2b_management_policy_id = "00000000-0000-0000-0000-000000000010"
  directory_object_id      = "00000000-0000-0000-0000-000000000020"
}
