resource "microsoft365_graph_beta_identity_and_access_b2b_management_policy_assignment" "service_principal" {
  b2b_management_policy_id = "00000000-0000-0000-0000-000000000011"
  directory_object_id      = "00000000-0000-0000-0000-000000000020"
}
