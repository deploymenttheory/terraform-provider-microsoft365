# policy_id and role_definition_id cannot both be set.
resource "microsoft365_graph_beta_identity_and_access_role_management_policy" "test" {
  policy_id          = "DirectoryRole_00000000-0000-0000-0000-000000000001_22222222-2222-2222-2222-222222222222"
  role_definition_id = "11111111-1111-1111-1111-111111111111"
}
