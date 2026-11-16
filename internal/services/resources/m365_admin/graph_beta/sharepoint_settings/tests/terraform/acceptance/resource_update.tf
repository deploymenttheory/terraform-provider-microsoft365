# Update: the low-impact settings change in place and idle session sign-out is enabled with long timeouts.
resource "microsoft365_graph_beta_m365_admin_sharepoint_settings" "test" {
  deleted_user_personal_site_retention_period_in_days = 90
  excluded_file_extensions_for_sync_app               = ["*.acc-test-tmp", "*.acc-test-bak"]
  is_mac_sync_app_enabled                             = true
  is_sync_button_hidden_on_personal_site              = false

  idle_session_sign_out = {
    is_enabled                = true
    warn_after_in_seconds     = 3300
    sign_out_after_in_seconds = 3600
  }
}
