resource "microsoft365_graph_beta_identity_and_access_authentication_methods_policy" "test" {
  temporary_access_pass = {
    state                       = "enabled"
    default_lifetime_in_minutes = 120
    default_length              = 12
    minimum_lifetime_in_minutes = 60
    maximum_lifetime_in_minutes = 480
    is_usable_once              = true
    include_targets = [
      {
        id                       = "33333333-3333-3333-3333-333333333333"
        target_type              = "group"
        is_registration_required = false
      }
    ]
    exclude_targets = []
  }
}
