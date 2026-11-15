# Adopts the policy by its ID; the role definition and scope are computed.
resource "microsoft365_graph_beta_identity_and_access_role_management_policy" "test" {
  policy_id = "DirectoryRole_00000000-0000-0000-0000-000000000001_22222222-2222-2222-2222-222222222222"

  expiration_rules = {
    Expiration_Admin_Assignment = {
      is_expiration_required = false
      maximum_duration       = "P90D"
    }
  }
}
