# Configures an enablement rule under expiration_rules.
resource "microsoft365_graph_beta_identity_and_access_role_management_policy" "test" {
  role_definition_id = "11111111-1111-1111-1111-111111111111"

  expiration_rules = {
    Enablement_Admin_Assignment = {
      is_expiration_required = true
      maximum_duration       = "P30D"
    }
  }
}
