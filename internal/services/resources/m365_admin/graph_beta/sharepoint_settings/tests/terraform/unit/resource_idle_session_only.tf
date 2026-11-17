# Only is_enabled is managed inside idle_session_sign_out; the timeouts keep their tenant values.
resource "microsoft365_graph_beta_m365_admin_sharepoint_settings" "test" {
  idle_session_sign_out = {
    is_enabled = false
  }
}
