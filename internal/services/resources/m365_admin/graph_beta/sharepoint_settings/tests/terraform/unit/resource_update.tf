# Update: several settings change in place and the sharing lists switch from allow to block.
resource "microsoft365_graph_beta_m365_admin_sharepoint_settings" "test" {
  allowed_domain_guids_for_sync_app                       = []
  deleted_user_personal_site_retention_period_in_days     = 90
  excluded_file_extensions_for_sync_app                   = ["*.pst"]
  image_tagging_option                                    = "disabled"
  is_commenting_on_site_pages_enabled                     = true
  is_file_activity_notification_enabled                   = true
  is_legacy_auth_protocols_enabled                        = false
  is_loop_enabled                                         = false
  is_mac_sync_app_enabled                                 = true
  is_require_accepting_user_to_match_invited_user_enabled = false
  is_resharing_by_external_users_enabled                  = true
  is_sharepoint_mobile_notification_enabled               = true
  is_sharepoint_newsfeed_enabled                          = false
  is_site_creation_enabled                                = true
  is_site_creation_ui_enabled                             = true
  is_site_pages_creation_enabled                          = true
  is_sites_storage_limit_automatic                        = true
  is_sync_button_hidden_on_personal_site                  = false
  is_unmanaged_sync_app_for_tenant_restricted             = false
  personal_site_default_storage_limit_in_mb               = 1048576
  sharing_allowed_domain_list                             = []
  sharing_blocked_domain_list                             = ["fabrikam.com"]
  sharing_capability                                      = "externalUserSharingOnly"
  sharing_domain_restriction_mode                         = "blockList"
  site_creation_default_managed_path                      = "/sites/"
  site_creation_default_storage_limit_in_mb               = 26214400
  tenant_default_timezone                                 = "(UTC-08:00) Pacific Time (US and Canada)"

  idle_session_sign_out = {
    is_enabled                = false
    warn_after_in_seconds     = 3300
    sign_out_after_in_seconds = 3600
  }
}
