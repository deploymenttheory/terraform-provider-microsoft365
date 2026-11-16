# Graph rejects an idle session warning that is not earlier than the sign-out.
resource "microsoft365_graph_beta_m365_admin_sharepoint_settings" "test" {
  idle_session_sign_out = {
    is_enabled                = true
    warn_after_in_seconds     = 3600
    sign_out_after_in_seconds = 1800
  }
}
