# Manage every device registration setting.
# When Microsoft Intune is enabled, azure_ad_registration cannot be changed and must match the
# current tenant value (usually "all"), or be left out.
resource "microsoft365_graph_beta_identity_and_access_device_registration_policy" "example" {
  user_device_quota               = 20
  multi_factor_auth_configuration = "required"

  local_admin_password = {
    is_enabled = true
  }

  azure_ad_registration = {
    allowed_to_register = {
      type = "all"
    }
  }

  azure_ad_join = {
    allowed_to_join = {
      type = "all"
    }
    local_admins = {
      enable_global_admins = true
      registering_users = {
        type = "none"
      }
    }
  }
}
