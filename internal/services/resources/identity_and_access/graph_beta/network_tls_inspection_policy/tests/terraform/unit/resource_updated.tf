resource "microsoft365_graph_beta_identity_and_access_network_tls_inspection_policy" "test" {
  name = "tf-api-probe-policy-updated"

  default_action = "bypass"
}
