resource "microsoft365_graph_beta_identity_and_access_device_registration_policy" "test" {
  local_admin_password = {
    is_enabled = true
  }
}
