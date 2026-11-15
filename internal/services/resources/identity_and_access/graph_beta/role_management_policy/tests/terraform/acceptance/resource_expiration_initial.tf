# Adopts the PIM policy of the built-in Usage Summary Reports Reader role and changes one expiration rule.
resource "microsoft365_graph_beta_identity_and_access_role_management_policy" "test" {
  role_definition_id = "75934031-6c7e-415a-99d7-48dbd49e875e"

  expiration_rules = {
    Expiration_Admin_Assignment = {
      is_expiration_required = false
      maximum_duration       = "P90D"
    }
  }
}
