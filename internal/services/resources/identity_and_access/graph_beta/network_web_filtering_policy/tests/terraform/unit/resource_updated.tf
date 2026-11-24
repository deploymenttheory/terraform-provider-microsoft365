resource "microsoft365_graph_beta_identity_and_access_network_web_filtering_policy" "test" {
  name           = "unit-test-web-filtering-policy-updated"
  description    = "Updated unit test web filtering policy"
  default_action = "block"
}
