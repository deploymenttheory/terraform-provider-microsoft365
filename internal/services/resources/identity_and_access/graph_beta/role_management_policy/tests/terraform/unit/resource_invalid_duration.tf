# Graph rejects week durations.
resource "microsoft365_graph_beta_identity_and_access_role_management_policy" "test" {
  role_definition_id = "11111111-1111-1111-1111-111111111111"

  expiration_rules = {
    Expiration_Admin_Assignment = {
      is_expiration_required = false
      maximum_duration       = "P2W"
    }
  }
}
