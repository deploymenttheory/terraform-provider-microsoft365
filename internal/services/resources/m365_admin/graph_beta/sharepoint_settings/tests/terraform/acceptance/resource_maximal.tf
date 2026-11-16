# Maximal configuration for the acceptance tenant: only low-impact settings are managed.
# Sharing, Loop, legacy authentication and unmanaged device settings are never changed by the acceptance tests.
resource "microsoft365_graph_beta_m365_admin_sharepoint_settings" "test" {
  deleted_user_personal_site_retention_period_in_days = 60
  excluded_file_extensions_for_sync_app               = ["*.acc-test-tmp"]
  is_mac_sync_app_enabled                             = false
  is_sync_button_hidden_on_personal_site              = true

  idle_session_sign_out = {
    is_enabled = false
  }
}
