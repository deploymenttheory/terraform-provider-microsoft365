resource "microsoft365_graph_beta_identity_and_access_authorization_policy" "test" {
  display_name                                              = "Authorization Policy"
  description                                               = "Tenant authorization settings"
  allowed_to_sign_up_email_based_subscriptions              = false
  allowed_to_use_sspr                                       = true
  allow_email_verified_users_to_join_organization           = false
  allow_invites_from                                        = "none"
  allow_user_consent_for_risky_apps                         = false
  block_msol_powershell                                     = true
  enabled_preview_features                                  = ["test-preview"]
  guest_user_role_id                                        = "2af84b1e-32c8-42b7-82bc-daa82404023b"
  permission_grant_policy_ids_assigned_to_default_user_role = ["ManagePermissionGrantsForSelf.microsoft-user-default-low"]
  default_user_role_permissions = {
    allowed_to_create_agent_identity_blueprints     = false
    allowed_to_create_apps                          = false
    allowed_to_create_security_groups               = false
    allowed_to_create_tenants                       = false
    allowed_to_read_bitlocker_keys_for_owned_device = false
    allowed_to_read_other_users                     = true
  }
}
