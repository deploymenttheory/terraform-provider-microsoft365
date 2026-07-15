resource "microsoft365_graph_beta_identity_and_access_network_content_policy" "test" {
  name           = "unit-test-content-policy-minimal"
  default_action = "allow"
}
