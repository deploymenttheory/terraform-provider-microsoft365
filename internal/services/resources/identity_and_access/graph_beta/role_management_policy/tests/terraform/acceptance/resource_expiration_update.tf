# Changes the managed expiration rule again and starts managing the activation expiration rule.
resource "microsoft365_graph_beta_identity_and_access_role_management_policy" "test" {
  role_definition_id = "75934031-6c7e-415a-99d7-48dbd49e875e"

  expiration_rules = {
    Expiration_Admin_Assignment = {
      is_expiration_required = false
      maximum_duration       = "P60D"
    }
    Expiration_EndUser_Assignment = {
      is_expiration_required = true
      maximum_duration       = "PT4H"
    }
  }
}
