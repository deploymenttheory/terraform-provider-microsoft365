resource "microsoft365_graph_beta_identity_and_access_authentication_methods_policy" "test" {
  email = {
    state                              = "enabled"
    allow_external_id_to_use_email_otp = "disabled"
    include_targets = [
      {
        id          = "33333333-3333-3333-3333-333333333333"
        target_type = "group"
      }
    ]
    exclude_targets = []
  }
}
