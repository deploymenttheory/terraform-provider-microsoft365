resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_groups_group" "join" {
  display_name     = "acc-test-device-registration-join-${random_string.suffix.result}"
  mail_nickname    = "acc-test-drp-join-${random_string.suffix.result}"
  mail_enabled     = false
  security_enabled = true
  description      = "Acceptance test group allowed to join devices to Microsoft Entra ID"
  hard_delete      = true
}

# Maximal configuration: every configurable setting is managed.
# azure_ad_registration is set to "all", the value Intune enforces in the test tenant.
resource "microsoft365_graph_beta_identity_and_access_device_registration_policy" "test" {
  user_device_quota               = 45
  multi_factor_auth_configuration = "notRequired"

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
      groups = [microsoft365_graph_beta_groups_group.join.id]
    }
    local_admins = {
      enable_global_admins = true
      registering_users = {
        type = "none"
      }
    }
  }
}
