# No policy is assigned to this role definition.
resource "microsoft365_graph_beta_identity_and_access_role_management_policy" "test" {
  role_definition_id = "99999999-9999-9999-9999-999999999999"

  expiration_rules = {
    Expiration_Admin_Assignment = {
      is_expiration_required = false
      maximum_duration       = "P90D"
    }
  }
}
