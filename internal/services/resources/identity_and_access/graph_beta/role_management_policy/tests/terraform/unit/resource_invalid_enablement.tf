# Enablement rules accept MultiFactorAuthentication, Justification and Ticketing only.
resource "microsoft365_graph_beta_identity_and_access_role_management_policy" "test" {
  role_definition_id = "11111111-1111-1111-1111-111111111111"

  enablement_rules = {
    Enablement_EndUser_Assignment = {
      enabled_rules = ["Approval"]
    }
  }
}
