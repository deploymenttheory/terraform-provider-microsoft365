# Maximal configuration: every configurable setting is managed.
resource "microsoft365_graph_beta_m365_admin_sharepoint_settings" "test" {
  allowed_domain_guids_for_sync_app                       = ["bdd1ab9b-3fd0-4def-a761-ec8d7471732c"]
  deleted_user_personal_site_retention_period_in_days     = 365
  excluded_file_extensions_for_sync_app                   = ["*.pst", "*.mp3"]
  image_tagging_option                                    = "enhanced"
  is_commenting_on_site_pages_enabled                     = false
  is_file_activity_notification_enabled                   = false
  is_legacy_auth_protocols_enabled                        = false
  is_loop_enabled                                         = true
  is_mac_sync_app_enabled                                 = false
  is_require_accepting_user_to_match_invited_user_enabled = true
  is_resharing_by_external_users_enabled                  = false
  is_sharepoint_mobile_notification_enabled               = false
  is_sharepoint_newsfeed_enabled                          = true
  is_site_creation_enabled                                = false
  is_site_creation_ui_enabled                             = false
  is_site_pages_creation_enabled                          = false
  is_sites_storage_limit_automatic                        = false
  is_sync_button_hidden_on_personal_site                  = true
  is_unmanaged_sync_app_for_tenant_restricted             = true
  personal_site_default_storage_limit_in_mb               = 512000
  sharing_allowed_domain_list                             = ["contoso.com", "fabrikam.com"]
  sharing_blocked_domain_list                             = []
  sharing_capability                                      = "existingExternalUserSharingOnly"
  sharing_domain_restriction_mode                         = "allowList"
  site_creation_default_managed_path                      = "/teams/"
  site_creation_default_storage_limit_in_mb               = 1048576
  tenant_default_timezone                                 = "(UTC) Coordinated Universal Time"

  idle_session_sign_out = {
    is_enabled                = true
    warn_after_in_seconds     = 1800
    sign_out_after_in_seconds = 3600
  }
}
