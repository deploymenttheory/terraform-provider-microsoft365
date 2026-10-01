resource "microsoft365_graph_beta_identity_and_access_authentication_flows_policy" "test" {
  self_service_sign_up = {
    is_enabled = false
  }
}
