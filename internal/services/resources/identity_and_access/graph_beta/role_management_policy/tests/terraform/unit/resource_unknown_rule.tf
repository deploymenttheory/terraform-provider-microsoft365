# References a rule ID that does not exist on the policy.
resource "microsoft365_graph_beta_identity_and_access_role_management_policy" "test" {
  role_definition_id = "11111111-1111-1111-1111-111111111111"

  expiration_rules = {
    Expiration_EndUser_Eligibility = {
      is_expiration_required = true
      maximum_duration       = "P30D"
    }
  }
}
