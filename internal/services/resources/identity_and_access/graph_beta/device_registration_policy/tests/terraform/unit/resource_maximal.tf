# Maximal configuration: every configurable setting is managed.
resource "microsoft365_graph_beta_identity_and_access_device_registration_policy" "test" {
  user_device_quota               = 20
  multi_factor_auth_configuration = "required"

  local_admin_password = {
    is_enabled = false
  }

  azure_ad_registration = {
    allowed_to_register = {
      type = "all"
    }
  }

  azure_ad_join = {
    allowed_to_join = {
      type   = "enumerated"
      users  = ["33333333-3333-3333-3333-333333333333"]
      groups = ["44444444-4444-4444-4444-444444444444"]
    }
    local_admins = {
      enable_global_admins = false
      registering_users = {
        type = "none"
      }
    }
  }
}
