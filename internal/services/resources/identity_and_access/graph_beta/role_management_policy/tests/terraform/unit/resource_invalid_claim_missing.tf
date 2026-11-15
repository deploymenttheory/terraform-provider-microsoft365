# An enabled authentication context rule needs a claim value.
resource "microsoft365_graph_beta_identity_and_access_role_management_policy" "test" {
  role_definition_id = "11111111-1111-1111-1111-111111111111"

  authentication_context_rules = {
    AuthenticationContext_EndUser_Assignment = {
      is_enabled = true
    }
  }
}
