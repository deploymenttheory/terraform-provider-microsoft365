resource "microsoft365_graph_beta_identity_and_access_authentication_methods_policy" "test" {
  sms = {
    state = "disabled"
    include_targets = [
      {
        id                    = "all_users"
        target_type           = "group"
        is_usable_for_sign_in = true
      }
    ]
  }

  software_oath = {
    state = "disabled"
  }

  hardware_oath = {
    state = "enabled"
    include_targets = [
      {
        id          = "33333333-3333-3333-3333-333333333333"
        target_type = "group"
      }
    ]
    exclude_targets = [
      {
        id          = "44444444-4444-4444-4444-444444444444"
        target_type = "group"
      }
    ]
  }
}
