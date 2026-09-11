resource "microsoft365_graph_beta_identity_and_access_network_threat_intelligence_policy" "test" {
  name           = "tf-api-probe-policy"
  description    = "Initial description"
  default_action = "allow"
}
