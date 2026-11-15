# Changes the managed expiration rule and starts managing a second one.
resource "microsoft365_graph_beta_identity_and_access_role_management_policy" "test" {
  role_definition_id = "11111111-1111-1111-1111-111111111111"

  expiration_rules = {
    Expiration_Admin_Assignment = {
      is_expiration_required = true
      maximum_duration       = "P60D"
    }
    Expiration_EndUser_Assignment = {
      is_expiration_required = true
      maximum_duration       = "PT4H"
    }
  }
}
