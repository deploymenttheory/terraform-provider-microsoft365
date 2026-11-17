# Manage every SharePoint and OneDrive tenant setting.
# Review each value before applying: sharing, legacy authentication and sync app
# restrictions affect every user in the tenant.
resource "microsoft365_graph_beta_m365_admin_sharepoint_settings" "example" {
  # OneDrive sync app
  allowed_domain_guids_for_sync_app           = ["bdd1ab9b-3fd0-4def-a761-ec8d7471732c"]
  excluded_file_extensions_for_sync_app       = ["*.pst", "*.mp3"]
  is_mac_sync_app_enabled                     = true
  is_sync_button_hidden_on_personal_site      = false
  is_unmanaged_sync_app_for_tenant_restricted = true

  # OneDrive storage and retention
  deleted_user_personal_site_retention_period_in_days = 365
  personal_site_default_storage_limit_in_mb           = 1048576

  # External sharing
  sharing_capability                                      = "externalUserSharingOnly"
  sharing_domain_restriction_mode                         = "allowList"
  sharing_allowed_domain_list                             = ["contoso.com", "fabrikam.com"]
  sharing_blocked_domain_list                             = []
  is_resharing_by_external_users_enabled                  = false
  is_require_accepting_user_to_match_invited_user_enabled = true

  # Sites and pages
  is_site_creation_enabled                  = true
  is_site_creation_ui_enabled               = true
  is_site_pages_creation_enabled            = true
  is_sites_storage_limit_automatic          = true
  site_creation_default_managed_path        = "/sites/"
  site_creation_default_storage_limit_in_mb = 26214400
  tenant_default_timezone                   = "(UTC) Coordinated Universal Time"
  is_commenting_on_site_pages_enabled       = true
  is_sharepoint_newsfeed_enabled            = false
  image_tagging_option                      = "basic"

  # Notifications, Loop and authentication
  is_file_activity_notification_enabled     = true
  is_sharepoint_mobile_notification_enabled = true
  is_loop_enabled                           = true
  is_legacy_auth_protocols_enabled          = false

  # Sign out users after one hour of inactivity, with a warning five minutes before.
  idle_session_sign_out = {
    is_enabled                = true
    warn_after_in_seconds     = 3300
    sign_out_after_in_seconds = 3600
  }

  timeouts = {
    create = "5m"
    read   = "5m"
    update = "5m"
    delete = "5m"
  }
}
