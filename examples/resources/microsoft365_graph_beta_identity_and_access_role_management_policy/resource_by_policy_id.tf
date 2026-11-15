# Adopt a policy by its ID instead of by role definition. role_definition_id, scope_id and scope_type
# are then computed from the policy and its assignment.
resource "microsoft365_graph_beta_identity_and_access_role_management_policy" "by_policy_id" {
  policy_id = "DirectoryRole_00000000-0000-0000-0000-000000000000_00000000-0000-0000-0000-000000000000"

  expiration_rules = {
    Expiration_Admin_Eligibility = {
      is_expiration_required = true
      maximum_duration       = "P180D"
    }
  }
}
