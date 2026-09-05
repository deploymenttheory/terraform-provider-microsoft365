resource "microsoft365_graph_beta_identity_and_access_network_tls_inspection_policy" "example" {
  name           = "Example TLS inspection policy"
  description    = "Initial description"
  default_action = "inspect"
}
