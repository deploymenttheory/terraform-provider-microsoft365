resource "microsoft365_graph_beta_identity_and_access_network_threat_intelligence_policy" "example" {
  name           = "Example threat intelligence policy"
  description    = "Initial description"
  default_action = "allow"
}
