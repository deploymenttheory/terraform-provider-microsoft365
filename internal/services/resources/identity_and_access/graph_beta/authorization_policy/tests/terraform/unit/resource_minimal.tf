resource "microsoft365_graph_beta_identity_and_access_authorization_policy" "test" {
  allowed_to_sign_up_email_based_subscriptions    = true
  allowed_to_use_sspr                             = true
  allow_email_verified_users_to_join_organization = false
  allow_user_consent_for_risky_apps               = false
  block_msol_powershell                           = false
  allow_invites_from                              = "adminsAndGuestInviters"
  default_user_role_permissions = {
    allowed_to_create_agent_identity_blueprints     = true
    allowed_to_create_apps                          = true
    allowed_to_create_security_groups               = true
    allowed_to_create_tenants                       = true
    allowed_to_read_bitlocker_keys_for_owned_device = true
    allowed_to_read_other_users                     = true
  }
}
