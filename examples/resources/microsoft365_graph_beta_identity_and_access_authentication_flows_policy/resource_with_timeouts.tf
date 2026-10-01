resource "microsoft365_graph_beta_identity_and_access_authentication_flows_policy" "example" {
  self_service_sign_up = {
    is_enabled = false
  }

  timeouts = {
    create = "5m"
    read   = "1m"
    update = "5m"
  }
}
